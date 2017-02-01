package qtypes

import (
	"github.com/deckarep/golang-set"
	"github.com/zpatrick/go-config"
)

// QWorker
type QWorker struct {
	Cfg   *config.Config
	QChan Channels
	Subs  mapset.Set
}

// AddSub puts subscription in set
func (qw *QWorker) AddSub(sub string) {
	qw.Subs.Add(sub)
}

// RmSub removes subscription from set
func (qw *QWorker) RmSub(sub string) {
	qw.Subs.Remove(sub)
}

// QWorkers provides the basic function each worker must implement
type QWorkers interface {
	Run()
}
