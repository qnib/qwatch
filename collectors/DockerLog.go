package qcollect

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/qnib/qwatch/types"
	"github.com/qnib/qwatch/utils"
)

// RunDockerLogCollector start a UDP server to listen for GELF messages (uncompressed)
func RunDockerLogCollector(cmd *cobra.Command, qChan qtypes.Channels) {
	port, _ := strconv.Atoi(cmd.Flag("gelf-port").Value.String())

	fmt.Printf("Start DockerLog collector listening on port %d\n", port)
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	utils.CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	utils.CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	// Join the broadcast group
	bg := qChan.Group.Join()

	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		dat := []byte(buf[0:n])
		msg := qtypes.GelfMsg{}
		json.Unmarshal(dat, &msg)
		qm := qtypes.NewQmsg("DockerLog", msg.Msg, msg.Host)
		cnt := qtypes.ContainerInfo{
			ContainerID:   msg.ContainerID,
			ContainerName: msg.ContainerName,
			Created:       msg.Created,
			ImageID:       msg.ImageID,
			ImageName:     msg.ImageName,
			Command:       msg.Command,
			Tag:           msg.Tag,
		}
		qm.SetContainer(cnt)
		qm.Type = "container.stdout"
		bg.Send(qm)
	}
}
