package qoutput

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/qnib/qwatch/types"
	"github.com/spf13/cobra"
)

// InventoryOutput holds the struct
type InventoryOutput struct {
	buffer chan qtypes.Qmsg
	cmd    *cobra.Command
	qChan  qtypes.Channels
	cli    *client.Client
}

// NewInventoryOutput returns an initial instance
func NewInventoryOutput(cmd *cobra.Command, qC qtypes.Channels) InventoryOutput {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	return InventoryOutput{
		buffer: make(chan qtypes.Qmsg, 1000),
		cmd:    cmd,
		qChan:  qC,
		cli:    cli,
	}
}

// Run creates json files from the inventory
func (io *InventoryOutput) Run() {
	fmt.Println("[START] InventoryOutput")
	ig := io.qChan.Inventory.Join()
	for {
		select {
		case val := <-ig.In:
			switch val.(type) {
			case types.Container:
				cnt := val.(types.Container)
				writeMounts(cnt)
			}
		}
	}
}

func writeMounts(cnt types.Container) {
	fname, fs := extractMounts(cnt)
	writeJSON(fname, fs)
}

// DeleteItem removes a given item from the array
func DeleteItem(in []string, selector string) []string {
	var r []string
	for _, str := range in {
		if str != selector {
			r = append(r, str)
		}
	}
	return r
}

func logit(app, msg string) {
	ts := time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	fmt.Printf("%-35s [%-10s] %s\n", ts, app, msg)
}

func writeJSON(fname string, j interface{}) {
	b, _ := json.Marshal(j)
	logit("InvWRITER", fmt.Sprintf("%s > %s", fname, string(b)))
	err := ioutil.WriteFile(fname, b, 0644)
	if err != nil {
		logit("InvWRITER", fmt.Sprintf("%s", err))
	} else {
		logit("InvWRITER", fmt.Sprintf("Wrote to %s", fname))
	}
}

func extractMounts(cnt types.Container) (string, qtypes.DiscoveryFS) {
	fs := make(map[string]qtypes.DiscoveryFilesystem)
	for _, mnt := range cnt.Mounts {
		mopts := []string{"ro"}
		if mnt.RW {
			mopts = []string{"rw"}
		}
		df := qtypes.DiscoveryFilesystem{
			Device:       mnt.Source,
			Mountpoint:   mnt.Destination,
			FsType:       string(mnt.Type),
			FsOpts:       make(map[string]string),
			Mountoptions: mopts,
		}
		fs[mnt.Destination] = df
	}
	fname := fmt.Sprintf("%s_filesystems.json", strings.Trim(cnt.Names[0], "/"))
	dfs := qtypes.DiscoveryFS{
		Discovery: fs,
		Timestamp: time.Now().Unix(),
	}
	return fname, dfs
}
