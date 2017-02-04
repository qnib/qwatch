package qinput

import (
	"log"
	"math"

	"golang.org/x/net/context"

	"github.com/docker/docker/api/types"

	"github.com/docker/docker/client"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qwatch/types"
)

// DockerAPI is a simple qworker
type DockerAPI struct {
	qtypes.QWorker
	cli  *client.Client
	info types.Info
}

// NewDockerAPI returns instance of DockerEventInput
func NewDockerAPI(cfg *config.Config, qC qtypes.Channels) DockerAPI {
	de := DockerAPI{}
	de.Cfg = cfg
	de.QChan = qC
	return de
}

// Run subscribes connects to the API
func (de DockerAPI) Run() {
	var err error
	de.cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	de.info, err = de.cli.Info(context.Background())
	if err != nil {
		log.Printf("[EE] Error during Info(): %v >err> %s", de.info, err)
	}

	tick := de.QChan.Tick.Join()
	for {
		select {
		case t := <-tick.In:
			de.querySwarm(t)
			de.queryImages(t)
		}
	}
}

func (de DockerAPI) querySwarm(t interface{}) {
	swarmMod, _ := de.Cfg.Int("input.docker-api.swarm.tick")
	tick := float64(t.(int64))
	if tick != 0 && math.Mod(tick, float64(swarmMod)) != 0 {
		return
	}
	info, err := de.cli.Info(context.Background())
	if err != nil {
		log.Printf("[EE] Error during Info(): %v >err> %s", info, err)
	} else {
		de.QChan.Inventory.Send(info)
	}

	nodes, err := de.cli.NodeList(context.Background(), types.NodeListOptions{})
	if err != nil {
		log.Printf("[EE] Error during NodeList(): ", err)
	} else {
		for _, node := range nodes {
			qnode := new(qtypes.DockerNode)
			qnode.Node = node
			qnode.EngineID = info.ID
			de.QChan.Inventory.Send(*qnode)
		}
	}
}

func (de DockerAPI) queryImages(t interface{}) {
	imgTick, _ := de.Cfg.Int("input.docker-api.images.tick")
	tick := float64(t.(int64))
	if !(tick == 0 || math.Mod(tick, float64(imgTick)) == 0) {
		return
	}
	images, err := de.cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		log.Printf("[EE] Error during NodeList(): ", err)
	} else {
		for _, image := range images {
			qimg := new(qtypes.DockerImageSummary)
			qimg.ImageSummary = image
			qimg.EngineID = de.info.ID
			de.QChan.Inventory.Send(*qimg)
		}
	}
}