package qoutput

import (
	"fmt"
    "github.com/deckarep/golang-set"
    "github.com/zpatrick/go-config"

	"github.com/qnib/qwatch/types"
)
// Log holds a buffer and the initial information from the server
type Log struct {
    qtypes.QWorker
}
// NewLog returns an initial instance
func NewLog(cfg *config.Config, qC qtypes.Channels) Log {
    l := Log{}
    l.Cfg = cfg
    l.QChan = qC
	l.Subs = mapset.NewSet()
	return l
}

// Run prints the logs to stdout
func (l Log) Run() {
	bg := l.QChan.Log.Join()
	for {
		select {
		case val := <-bg.In:
			log := val.(qtypes.Qmsg)
			ts := log.Time.Format("2006-01-02T15:04:05.999999-07:00")
			fmt.Printf("%-35v | Source:%-20s | Type:%-20s | Host:%-20s | %-50s  >> Log: %v\n", ts, log.Source, log.Type, log.Host, log.GetCntInfo(), log.Msg)
		}
	}
}
