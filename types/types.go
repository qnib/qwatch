package qtypes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/grafov/bcast"
)

// Channels holds the communication channels
type Channels struct {
	Log       *bcast.Group
	Inventory *bcast.Group
	Tick      *bcast.Group
	Done      chan os.Signal
}

// ContainerInfo holds information when the message was emmited by a container
type ContainerInfo struct {
	Command       string `json:"_command"`
	ContainerID   string `json:"_container_id"`
	ContainerName string `json:"_container_name"`
	Created       string `json:"_created"`
	ImageID       string `json:"_image_id"`
	ImageName     string `json:"_image_name"`
	Tag           string `json:"_tag"`
}

// GetCntInfo returns string representation for logging
func (ci *ContainerInfo) GetCntInfo() string {
	info := fmt.Sprintf("Image:%-15s # Name:%-15s", ci.ImageName, ci.ContainerName)
	return strings.TrimRight(info, " ")
}

// Qmsg holds a GELF message from a container
type Qmsg struct {
	Version     string        `json:"version"`
	Type        string        `json:"type"`
	Source      string        `json:"source"`
	Host        string        `json:"host"`
	Msg         string        `json:"short_message"`
	Time        time.Time     `json:"time"`
	Level       int           `json:"level"`
	IsContainer bool          `json:"is_container"`
	Container   ContainerInfo `json:"container"`
}

// GetCntInfo returns container info if a container is involved, an empty string otherwise
func (qm *Qmsg) GetCntInfo() string {
	if qm.IsContainer {
		return qm.Container.GetCntInfo()
	}
	return ""
}

// GelfMsg holds a GELF message from a container
type GelfMsg struct {
	Version       string  `json:"version"`
	Host          string  `json:"host"`
	Msg           string  `json:"short_message"`
	Timestamp     float32 `json:"timestamp"`
	Level         int     `json:"level"`
	Command       string  `json:"_command"`
	ContainerID   string  `json:"_container_id"`
	ContainerName string  `json:"_container_name"`
	Created       string  `json:"_created"`
	ImageID       string  `json:"_image_id"`
	ImageName     string  `json:"_image_name"`
	Tag           string  `json:"_tag"`
}

// NewQmsg create an instance of Log
func NewQmsg(source, msg, host string) Qmsg {
	return Qmsg{
		Version:     "1.1",
		Source:      source,
		Host:        host,
		Msg:         msg,
		IsContainer: false,
		Time:        time.Now(),
	}
}

// SetContainer sets information about the container
func (qm *Qmsg) SetContainer(cnt ContainerInfo) Qmsg {
	qm.IsContainer = true
	qm.Container = cnt
	return *qm
}

// SetTimestamp updates the Timestamp attribute
func (qm *Qmsg) SetTimestamp(ts time.Time) Qmsg {
	qm.Time = ts
	return *qm
}
