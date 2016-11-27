package qoutput

import (
    "fmt"

    "github.com/qnib/qwatch/types"
    "github.com/spf13/cobra"
)
// RunLogOutput prints the logs to stdout
func RunLogOutput(cmd *cobra.Command, qC qtypes.Channels) {
    for {
        log := <-qC.Log
        ts := log.Time.Format("2006-01-02T15:04:05.999999-07:00")
        fmt.Printf("%-35v | Source:%-20s | Type:%-20s | Host:%-20s | %-50s  >> Log: %v\n", ts, log.Source, log.Type, log.Host, log.GetCntInfo(), log.Msg)
    }
}
