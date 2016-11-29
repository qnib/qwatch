package qserver

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/qnib/qwatch/collectors"
	"github.com/qnib/qwatch/output"
	"github.com/qnib/qwatch/utils"
)

// ServeQlog start daemon
func ServeQlog(cmd *cobra.Command, args []string) {
	qC := utils.NewChannels(cmd)
	// create broadcaster
	go qC.Group.Broadcasting(0) // accepts messages and broadcast it to all members
	// fetches interrupt and closes
	signal.Notify(qC.Done, os.Interrupt)

	// Log Collector
	go qcollect.RunDockerLogCollector(cmd, qC)
	go qcollect.RunDockerEventCollector(cmd, qC)

	// Handler
	go qoutput.RunLogOutput(cmd, qC)
	eo := qoutput.NewElasticsearchOutput(cmd, qC)
	go eo.RunElasticsearchOutput()
	for {
		<-qC.Done
		fmt.Printf("\nDone\n")
		return

	}
}
