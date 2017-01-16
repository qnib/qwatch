package qcollect

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/qnib/qwatch/types"
	"github.com/qnib/qwatch/utils"
	"github.com/urfave/cli"
)

// RunDockerLogCollector start a UDP server to listen for GELF messages (uncompressed)
func RunDockerLogCollector(ctx *cli.Context, qChan qtypes.Channels) {
	port := ctx.Int("gelf-port")

	log.Printf("Start DockerLog collector listening on port %d\n", port)
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	utils.CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	utils.CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	// Join the broadcast group
	bg := qChan.Log.Join()

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
