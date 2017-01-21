package qoutput

import (
	"fmt"

	"github.com/qnib/qwatch/types"
    "github.com/codegangsta/cli"
)

// RunLogOutput prints the logs to stdout
func RunLogOutput(ctx *cli.Context, qC qtypes.Channels) {
	bg := qC.Log.Join()
	for {
		select {
		case val := <-bg.In:
			log := val.(qtypes.Qmsg)
			ts := log.Time.Format("2006-01-02T15:04:05.999999-07:00")
			fmt.Printf("%-35v | Source:%-20s | Type:%-20s | Host:%-20s | %-50s  >> Log: %v\n", ts, log.Source, log.Type, log.Host, log.GetCntInfo(), log.Msg)
		}
	}
}
