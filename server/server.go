package qserver

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/qnib/qwatch/collectors"
	"github.com/qnib/qwatch/output"
	"github.com/qnib/qwatch/utils"
)

// ServeQlog start daemon
func ServeQlog(cmd *cobra.Command, args []string) {
	qC := utils.NewChannels(cmd)
	i, _ := strconv.Atoi(cmd.Flag("ticker-interval").Value.String())
	interval := time.Duration(i) * time.Millisecond
	ticker := time.NewTicker(interval).C

	// create broadcasters
	go qC.Log.Broadcasting(0)       // accepts messages and broadcast it to all members
	go qC.Tick.Broadcasting(0)      // accepts messages and broadcast it to all members
	go qC.Inventory.Broadcasting(0) // accepts messages and broadcast it to all members
	// fetches interrupt and closes
	signal.Notify(qC.Done, os.Interrupt)

	// Inventory Collector
	dc := qcollect.NewDockerInventoryCollector(cmd, qC)
	go dc.RunDockerInventoryCollector()

	// Log Collector
	go qcollect.RunDockerLogCollector(cmd, qC)
	go qcollect.RunDockerEventCollector(cmd, qC)

	// Handler
	go qoutput.RunLogOutput(cmd, qC)
	eo := qoutput.NewElasticsearchOutput(cmd, qC)
	go eo.RunElasticsearchOutput()
	io := qoutput.NewInventoryOutput(cmd, qC)
	go io.Run()
	// Inserts tick to get Inventory started
	qC.Tick.Send(0)
	for {
		select {
		case <-qC.Done:
			fmt.Printf("\nDone\n")
			return
		case <-ticker:
			qC.Tick.Send(0)
		}
	}
}
