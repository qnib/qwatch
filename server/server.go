package qserver

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/qnib/qwatch/collectors"
	"github.com/qnib/qwatch/output"
	"github.com/qnib/qwatch/utils"

	"github.com/urfave/cli"
)

// ServeQlog start daemon
func ServeQlog(ctx *cli.Context) error {
	qC := utils.NewChannels()
	i := ctx.Int("ticker-interval")
	interval := time.Duration(i) * time.Millisecond
	ticker := time.NewTicker(interval).C

	// create broadcasters
	go qC.Log.Broadcasting(0)       // accepts messages and broadcast it to all members
	go qC.Tick.Broadcasting(0)      // accepts messages and broadcast it to all members
	go qC.Inventory.Broadcasting(0) // accepts messages and broadcast it to all members
	// fetches interrupt and closes
	signal.Notify(qC.Done, os.Interrupt)

	// Collectors
	for _, c := range strings.Split(ctx.String("collectors"), ",") {
		switch c {
		case "gelf":
			log.Println("Start the GELF DockerLog collector")
			go qcollect.RunDockerLogCollector(ctx, qC)
		case "docker-events":
			log.Println("Start the DockerEvents collector")
			go qcollect.RunDockerEventCollector(ctx, qC)
		case "docker-inventory":
			log.Println("Start the DockerInventory collector")
			//dc := qcollect.NewDockerInventoryCollector(ctx, qC)
			//go dc.RunDockerInventoryCollector()
		}
	}

	// Filter
	/*
		  io := qoutput.NewInventoryOutput(cmd, qC)
			go io.Run()
	*/

	// Handler
	for _, c := range strings.Split(ctx.String("handlers"), ",") {
		switch c {
		case "log":
			log.Println("Start the log handler")
			go qoutput.RunLogOutput(ctx, qC)
		case "elasticsearch":
			log.Println("Start the elasticsearch handler")
			//eo := qoutput.NewElasticsearchOutput(ctx, qC)
			//go eo.RunElasticsearchOutput()
		}
	}
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
