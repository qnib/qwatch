package qtypes

import (
	"testing"
    "time"

	"github.com/stretchr/testify/assert"
)

func TestContainerInfo(t *testing.T) {
    ci := ContainerInfo{
        Command: "bash",
    	ContainerID: "123456",
    	ContainerName: "container01",
    	Created: "yesterday",
    	ImageID: "654321",
    	ImageName: "image01",
    	Tag: "",
    }
    assert.Equal(t, "Image:image01         # Name:container01", ci.GetCntInfo(), "GetCntInfo() should return a string representation of the container")
}

func TestQMSGGetContainerInfoEmpty(t *testing.T) {
    qm := NewQmsg("Source", "Message", "host1")
    assert.Equal(t, "", qm.GetCntInfo(), "GetCntInfo() should return empty string, as it is not a container")
}

func TestQMSGGetContainerInfo(t *testing.T) {
    qm := NewQmsg("Source", "Message", "host1")
    ci := ContainerInfo{
        Command: "bash",
    	ContainerID: "123456",
    	ContainerName: "container01",
    	Created: "yesterday",
    	ImageID: "654321",
    	ImageName: "image01",
    	Tag: "",
    }
    qm.SetContainer(ci)
    assert.Equal(t, "Image:image01         # Name:container01", qm.GetCntInfo(), "GetCntInfo() should return empty string, as it is not a container")
}

func TestQMSGSetTimestamp(t *testing.T) {
    qm := NewQmsg("Source", "Message", "host1")
    now := time.Now()
    qm.SetTimestamp(now)
    assert.Equal(t, now, qm.Time)
}
