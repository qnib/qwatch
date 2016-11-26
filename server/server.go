package qserver

import (
	"fmt"
	"os"
	"os/signal"

    "github.com/spf13/cobra"

	"github.com/qnib/qwatch/collectors"
    "github.com/qnib/qwatch/handlers"
	//"github.com/qnib/qwatch/types"
	"github.com/qnib/qwatch/utils"
)

// ServeQlog start daemon
func ServeQlog(cmd *cobra.Command, args []string) {
	qC := utils.NewChannels(cmd)
	// fetches interrupt and closes
	signal.Notify(qC.Done, os.Interrupt)

	// Log Collector
	go qcollect.RunDockerLogCollector(cmd, qC)
	go qcollect.RunDockerEventCollector(cmd, qC)

	// Handler
    go qhandler.RunLogHandler(cmd, qC)
	for {
	    <-qC.Done
		fmt.Printf("\nDone\n")
		return

	}
}
