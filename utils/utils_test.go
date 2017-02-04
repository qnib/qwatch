package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/qnib/qwatch/types"
)

func TestParseImageName(t *testing.T) {
	// image
	got := ParseImageName("image")
	exp := qtypes.ImageName{Name: "image", Tag: "latest"}
	assert.Equal(t, exp, got)
	// image@sha256:abc
	got = ParseImageName("image@sha256:sha")
	exp = qtypes.ImageName{Name: "image", Sha256: "sha"}
	assert.Equal(t, exp, got)
	// repo/image@sha256:abc
	got = ParseImageName("repo/image@sha256:sha")
	exp = qtypes.ImageName{Repository: "repo", Name: "image", Sha256: "sha"}
	assert.Equal(t, exp, got)
	// registry/repo/image@sha256:sha
	got = ParseImageName("registry/repo/image@sha256:sha")
	exp = qtypes.ImageName{Registry: "registry", Repository: "repo", Name: "image", Sha256: "sha"}
	assert.Equal(t, exp, got)
	// registry/repo/image:tag
	got = ParseImageName("registry/repo/image:tag")
	exp = qtypes.ImageName{Registry: "registry", Repository: "repo", Name: "image", Tag: "tag"}
	assert.Equal(t, exp, got)

}
