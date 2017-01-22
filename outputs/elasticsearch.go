package qoutput

import (
	"net/url"
	"time"
    "github.com/zpatrick/go-config"
    "github.com/deckarep/golang-set"

	"github.com/qnib/qwatch/types"
	"github.com/OwnLocal/goes"
)

// Elasticsearch holds a buffer and the initial information from the server
type Elasticsearch struct {
    qtypes.QWorker
    buffer chan qtypes.Qmsg
}

// NewElasticsearch returns an initial instance
func NewElasticsearch(cfg *config.Config, qC qtypes.Channels) Elasticsearch {
    subs := mapset.NewSet()
    es := Elasticsearch{
		buffer: make(chan qtypes.Qmsg, 1000),
    }
	es.Cfg = cfg
    es.QChan = qC
    es.Subs = subs
    return es
}

// Takes log from framework and buffers it in elasticsearch buffer
func (eo *Elasticsearch) pushToBuffer() {
	bg := eo.QChan.Log.Join()
	for {
		select {
		case val := <-bg.In:
			log := val.(qtypes.Qmsg)
			eo.buffer <- log
		}
	}
}

func (eo *Elasticsearch) createESClient() (conn *goes.Connection) {
    host, err := eo.Cfg.String("output.elasticsearch.host")
    if err != nil {
        panic(err)
    }
    port, err := eo.Cfg.String("output.elasticsearch.port")
    if err != nil {
        panic(err)
    }
    conn = goes.NewConnection(host, port)
	return
}

func createIndex(conn *goes.Connection) error {
	// Create an index
	mapping := map[string]interface{}{
		"settings": map[string]interface{}{
			"index.number_of_shards":   1,
			"index.number_of_replicas": 0,
		},
		"mappings": map[string]interface{}{
			"_default_": map[string]interface{}{
				"_source": map[string]interface{}{
					"enabled": true,
				},
				"_all": map[string]interface{}{
					"enabled": false,
				},
			},
		},
	}

	resp, err := conn.CreateIndex("logstash-2016-11-27", mapping)
	_ = resp
	//fmt.Printf("%v\n", resp)
	return err
}

func indexLog(conn *goes.Connection, log qtypes.Qmsg) error {
	now := time.Now()
	d := goes.Document{
		Index: "logstash-2016-11-27",
		Type:  "log",
		Fields: map[string]interface{}{
			"Timestamp": now.Format("2006-01-02T15:04:05.999999-07:00"),
			"msg":       log.Msg,
			"source":    log.Source,
			"type":      log.Type,
			"host":      log.Host,
		},
	}
	extraArgs := make(url.Values, 1)
	//extraArgs.Set("ttl", "86400000")
	response, err := conn.Index(d, extraArgs)

	_ = response
	//fmt.Printf("%s | %s\n", d, response.Error)
	return err
}

// Run pushes the logs to elasticsearch
func (eo Elasticsearch) Run() {
	go eo.pushToBuffer()
	conn := eo.createESClient()
	createIndex(conn)
	for {
		log := <-eo.buffer
		indexLog(conn, log)
	}
}
