package qinput

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
    "github.com/zpatrick/go-config"

	"github.com/qnib/qwatch/types"
	"github.com/qnib/qwatch/utils"
)

// DockerGelf is a simple qworker
type DockerGelf struct {
    qtypes.QWorker
}

// NewDockerGelf returns instance of DockerEventInput
func NewDockerGelf(cfg *config.Config, qC qtypes.Channels) DockerGelf {
    dg := DockerGelf{}
    dg.Cfg = cfg
    dg.QChan = qC
    return dg
}

// Run start a UDP server to listen for GELF messages (uncompressed)
func (dg DockerGelf) Run() {
	port, err := dg.Cfg.Int("input.docker-gelf.port")
    if err != nil {
        panic(err)
    }

	log.Printf("Start DockerGelf input listening on port %d\n", port)
	ServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	utils.CheckError(err)

	/* Now listen at selected port */
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	utils.CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	// Join the broadcast group
	bg := dg.QChan.Log.Join()

	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error: ", err)
		}
		dat := []byte(buf[0:n])
		msg := qtypes.GelfMsg{}
		json.Unmarshal(dat, &msg)
		qm := qtypes.NewQmsg("docker-gelf", msg.Msg, msg.Host)
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
