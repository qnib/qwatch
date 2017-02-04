package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/grafov/bcast"
	"github.com/qnib/qwatch/types"
)

// CheckError A Simple function to verify error
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

// NewChannels create an instance of Channels
func NewChannels() qtypes.Channels {
	//i, _ := strconv.Atoi(cmd.Flag("ticker-interval").Value.String())
	//interval := time.Duration(i) * time.Millisecond
	return qtypes.Channels{
		Log:       bcast.NewGroup(), // create broadcast group
		Inventory: bcast.NewGroup(), // create broadcast group
		Tick:      bcast.NewGroup(), // create broadcast group
		Done:      make(chan os.Signal, 1),
	}
}

//ParseImageName splits Docker image names
func ParseImageName(i string) qtypes.ImageName {
	// Should use regex, but let's start silly
	in := qtypes.ImageName{}
	sha := strings.Split(i, "@")
	if len(sha) != 1 {
		in.Sha256 = sha[1][7:]
		i = sha[0]
	}
	sTag := strings.Split(i, ":")
	if len(sTag) != 1 {
		in.Tag = sTag[1]
		i = sTag[0]
	} else if in.Sha256 == "" {
		in.Tag = "latest"
	}
	l := strings.Split(i, "/")
	in.Name, l = l[len(l)-1], l[:len(l)-1]
	if len(l) > 0 {
		in.Repository, l = l[len(l)-1], l[:len(l)-1]
	}
	if len(l) > 0 {
		in.Registry = l[len(l)-1]
	}
	return in
}
