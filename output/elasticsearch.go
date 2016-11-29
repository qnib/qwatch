package qoutput

import (
	//"fmt"
	"net/url"
	"time"

	"github.com/qnib/qwatch/types"
	"github.com/spf13/cobra"

	"github.com/OwnLocal/goes"
)

var (
	// EsHost to connect
	EsHost = "localhost"
	// EsPort to connect
	EsPort = "9200"
)

// ElasticsearchOutput holds a buffer and the initial information from the server
type ElasticsearchOutput struct {
	buffer chan qtypes.Qmsg
	cmd    *cobra.Command
	qChan  qtypes.Channels
}

// NewElasticsearchOutput returns an initial instance
func NewElasticsearchOutput(cmd *cobra.Command, qC qtypes.Channels) ElasticsearchOutput {
	return ElasticsearchOutput{
		buffer: make(chan qtypes.Qmsg, 1000),
		cmd:    cmd,
		qChan:  qC,
	}
}

// Takes log from framework and buffers it in elasticsearch buffer
func (eo *ElasticsearchOutput) pushToBuffer() {
	bg := eo.qChan.Group.Join()
	for {
		select {
		case val := <-bg.In:
			log := val.(qtypes.Qmsg)
			eo.buffer <- log
		}
	}
}

func createESClient() (conn *goes.Connection) {
	conn = goes.NewConnection("localhost", "9200")
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

// RunElasticsearchOutput pushes the logs to elasticsearch
func (eo *ElasticsearchOutput) RunElasticsearchOutput() {
	go eo.pushToBuffer()
	conn := createESClient()
	createIndex(conn)
	for {
		log := <-eo.buffer
		indexLog(conn, log)
	}
}
