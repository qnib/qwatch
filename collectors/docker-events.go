package qcollect

import (
	"fmt"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/qnib/qwatch/types"
)

// RunDockerEventCollector subscribes to messages and events from the docker-engine
func RunDockerEventCollector(cmd *cobra.Command, qChan qtypes.Channels) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	msgs, errs := cli.Events(context.Background(), types.EventsOptions{})
	bg := qChan.Group.Join()
	for {
		select {
		case dMsg := <-msgs:
			bg.Send(parseMessage(dMsg))
		case dErr := <-errs:
			if dErr != nil {
				qChan.Log <- qtypes.Qmsg{
					Msg: fmt.Sprintf("%s", dErr),
				}
			}
		}
	}
}

func parseMessage(msg events.Message) qtypes.Qmsg {
	host := os.Getenv("DOCKER_HOST")
	message := fmt.Sprintf("%s.%s", msg.Type, msg.Action)
	cnt := qtypes.ContainerInfo{
		ImageName:     msg.Actor.Attributes["image"],
		ContainerID:   msg.ID,
		ContainerName: msg.Actor.Attributes["name"],
	}
	qm := qtypes.Qmsg{
		Version:     "1.1",
		Source:      "DockerEvents",
		Host:        host,
		Msg:         message,
		IsContainer: false,
		Time:        time.Unix(0, msg.TimeNano),
	}
	qm.SetContainer(cnt)
	qm.Type = fmt.Sprintf("%s.%s", msg.Type, msg.Action)
	return qm
}
