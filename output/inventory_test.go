package qoutput

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/stretchr/testify/assert"

	"github.com/qnib/qwatch/types"
)

func TestDeleteEmpty(t *testing.T) {
	in := []string{"a", "", "b"}
	exp := []string{"a", "b"}
	got := DeleteItem(in, "")
	assert.Equal(t, exp, got)
}

func TestExtractMounts(t *testing.T) {
	in := types.Container{
		Mounts: []types.MountPoint{
			types.MountPoint{
				Source:      "/srv/",
				Destination: "/server/",
				Type:        "bind",
				RW:          true,
			},
			types.MountPoint{
				Source:      "/home/",
				Destination: "/home/",
				Type:        "bind",
				RW:          true,
			},
		},
		Names: []string{"/test-container"},
	}
	_ = in
	fs := make(map[string]qtypes.DiscoveryFilesystem)
	fs["/server/"] = qtypes.DiscoveryFilesystem{
		Device:       "/srv/",
		Mountpoint:   "/server/",
		FsType:       "bind",
		FsOpts:       make(map[string]string),
		Mountoptions: []string{"rw"},
	}
	fs["/home/"] = qtypes.DiscoveryFilesystem{
		Device:       "/home/",
		Mountpoint:   "/home/",
		FsType:       "bind",
		FsOpts:       make(map[string]string),
		Mountoptions: []string{"rw"},
	}
	exp := qtypes.DiscoveryFS{
		Discovery: fs,
		Timestamp: time.Now().Unix(),
	}
	fname, got := extractMounts(in)
	assert.Equal(t, exp, got)
	assert.Equal(t, "test-container_filesystems.json", fname)
}
