package qserver

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qwatch/inputs"
	"github.com/qnib/qwatch/outputs"
	"github.com/qnib/qwatch/types"
	"github.com/qnib/qwatch/utils"
)

// ServeQlog start daemon
func ServeQlog(ctx *cli.Context) error {
	//cfo := utils.NewQGraph()
	cfg := config.NewConfig([]config.Provider{})
	if _, err := os.Stat(ctx.String("config")); err == nil {
		log.Printf("[II] Use config file: %s", ctx.String("config"))
		cfg.Providers = append(cfg.Providers, config.NewYAMLFile(ctx.String("config")))
	}
	cfg.Providers = append(cfg.Providers, config.NewCLI(ctx, false))
	qC := utils.NewChannels()
	i, _ := cfg.Int("ticker.interval")
	interval := time.Duration(i) * time.Millisecond
	ticker := time.NewTicker(interval).C

	// create broadcasters
	go qC.Log.Broadcasting(0)       // accepts messages and broadcast it to all members
	go qC.Tick.Broadcasting(0)      // accepts messages and broadcast it to all members
	go qC.Inventory.Broadcasting(0) // accepts messages and broadcast it to all members
	// fetches interrupt and closes
	signal.Notify(qC.Done, os.Interrupt)

	// Inputs
	inputs, _ := cfg.String("inputs")
	for _, ins := range strings.Split(inputs, ",") {
		//cfo.AddInput(ins)
		var qw qtypes.QWorkers
		switch ins {
		case "docker-gelf":
			log.Println("Start the docker-gelf input")
			qw = qinput.NewDockerGelf(cfg, qC)
		case "docker-events":
			log.Println("Start the DockerEvents collector")
			qw = qinput.NewDockerEvents(cfg, qC)
		}
		go qw.Run()
	}

	// Filter
	/*
		  io := qoutput.NewInventoryOutput(cmd, qC)
			go io.Run()
	*/

	// Outputs
	outputs, _ := cfg.String("outputs")
	// TODO: iterate over keys in config output:
	for _, outs := range strings.Split(outputs, ",") {
		//cfo.AddOutput(h,[]string{"gelf"})
		var qw qtypes.QWorkers
		switch outs {
		case "log":
			log.Println("Start the log handler")
			qw = qoutput.NewLog(cfg, qC)
		case "neo4j":
			log.Println("Start the neo4j handler")
			qw = qoutput.NewNeo4j(cfg, qC)
		case "elasticsearch":
			log.Println("Start the elasticsearch handler")
			qw = qoutput.NewElasticsearch(cfg, qC)
		}
		go qw.Run()
	}
	//cfo.PrintGraph()
	// Inserts tick to get Inventory started
	qC.Tick.Send(0)
	for {
		select {
		case <-qC.Done:
			fmt.Printf("\nDone\n")
			return nil
		case <-ticker:
			qC.Tick.Send(0)
		}
	}
	return nil
}
