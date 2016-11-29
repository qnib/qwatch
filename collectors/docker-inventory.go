package qcollect

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/grafov/bcast"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"

	"github.com/qnib/qwatch/types"
)

// DockerInventoryCollector holds the object
type DockerInventoryCollector struct {
	DockerClient *client.Client
	RunCmd       *cobra.Command
	Groups       qtypes.Channels
}

// NewDockerInventoryCollector returns instance
func NewDockerInventoryCollector(cmd *cobra.Command, qGroups qtypes.Channels) DockerInventoryCollector {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return DockerInventoryCollector{
		DockerClient: cli,
		RunCmd:       cmd,
		Groups:       qGroups,
	}
}

// RunDockerInventoryCollector is triggered by the Ticker, fetches the /container/json and container/<id>/json end-point
func (dic *DockerInventoryCollector) RunDockerInventoryCollector() {
	tg := dic.Groups.Tick.Join()
	ig := dic.Groups.Inventory.Join()
	for {
		_ = tg.Recv()
		fmt.Println("Inventory: Tick received")
		dic.assembleInventory(ig)
	}
}

func (dic *DockerInventoryCollector) assembleInventory(ig *bcast.Member) {
	cfilter := filters.NewArgs()
	//cfilter.Add("id", task.ContainerID)
	containers, err := dic.DockerClient.ContainerList(context.Background(), types.ContainerListOptions{Filters: cfilter})
	if err != nil {
		fmt.Printf("ListContainers() failed: %v\n", err)
	}
	for _, container := range containers {
		ig.Send(container)
	}
}
