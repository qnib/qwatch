package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/grafov/bcast"
	"github.com/qnib/qwatch/types"
	"github.com/spf13/cobra"
)

// CheckError A Simple function to verify error
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// NewChannels create an instance of Channels
func NewChannels(cmd *cobra.Command) qtypes.Channels {
	i, _ := strconv.Atoi(cmd.Flag("ticker-interval").Value.String())
	interval := time.Duration(i) * time.Millisecond
	return qtypes.Channels{
		Tick:  time.NewTicker(interval).C,
		Log:   make(chan qtypes.Qmsg),
		Done:  make(chan os.Signal, 1),
		Group: bcast.NewGroup(), // create broadcast group
	}
}
