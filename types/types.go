package qtypes

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	"github.com/grafov/bcast"
)

// Channels holds the communication channels
type Channels struct {
	Log       *bcast.Group
	Inventory *bcast.Group
	Tick      *bcast.Group
	Done      chan os.Signal
}

// ImageInfo holds information about a docker image
type ImageInfo struct {
	Name string `json:"_name"`
	ID   string `json:"_image_id"`
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
	Action      string        `json:"action"`
	Source      string        `json:"source"`
	Host        string        `json:"host"`
	Msg         string        `json:"short_message"`
	Time        time.Time     `json:"time"`
	TimeNano    int64         `json:"time_nano"`
	Level       int           `json:"level"`
	IsContainer bool          `json:"is_container"`
	Container   ContainerInfo `json:"container"`
	Image       ImageInfo     `json:"container"`
	EngineID    string        `json:"engine_id"`
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

// SetImage sets information about the image
func (qm *Qmsg) SetImage(img ImageInfo) Qmsg {
	qm.Image = img
	return *qm
}

// SetTimestamp updates the Timestamp attribute
func (qm *Qmsg) SetTimestamp(ts time.Time) Qmsg {
	qm.Time = ts
	return *qm
}

// DockerNode is a superset of swarm.Node, which passes along the ID of cli.Info
type DockerNode struct {
	swarm.Node
	EngineID string
}

// ImageName holds information about a docker image
type ImageName struct {
	Registry   string
	Repository string
	Name       string
	Tag        string
	Sha256     string
}

func (in *ImageName) String() string {
	l := []string{in.Registry, in.Repository, in.Name}
	nl := l[:0]
	for _, x := range l {
		if x != "" {
			nl = append(nl, x)
		}
	}
	res := strings.Join(nl, "/")
	if in.Tag != "" {
		res = fmt.Sprintf("%s:%s", res, in.Tag)
	}
	if in.Sha256 != "" {
		res = fmt.Sprintf("%s@sha256:%s", res, in.Sha256)
	}
	return res
}

// DockerImageSummary is a superset of swarm.Node, which passes along the ID of cli.Info
type DockerImageSummary struct {
	types.ImageSummary
	EngineID string
}

func (di DockerImageSummary) String() string {
	var name string
	if len(di.RepoTags) > 0 {
		name = di.RepoTags[0]
	}
	return fmt.Sprintf("%-20s ID:%s", name, string(di.ID[7:19]))
}
