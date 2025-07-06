package models

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildTarget(t *testing.T) {
	target := BuildTarget{
		OS:   "linux",
		Arch: "amd64",
	}

	assert.Equal(t, "linux", target.OS)
	assert.Equal(t, "amd64", target.Arch)
}

func TestBuildResult(t *testing.T) {
	target := BuildTarget{OS: "linux", Arch: "amd64"}
	result := BuildResult{
		Target:     target,
		BinaryPath: "/path/to/binary",
		Error:      errors.New("test error"),
	}

	assert.Equal(t, "linux", result.Target.OS)
	assert.Equal(t, "/path/to/binary", result.BinaryPath)
	assert.Error(t, result.Error)
}

func TestArchiveResult(t *testing.T) {
	target := BuildTarget{OS: "linux", Arch: "amd64"}
	result := ArchiveResult{
		Target:      target,
		ArchivePath: "/path/to/archive.tar.gz",
		Error:       nil,
	}

	assert.Equal(t, "linux", result.Target.OS)
	assert.Equal(t, "/path/to/archive.tar.gz", result.ArchivePath)
	assert.NoError(t, result.Error)
}

func TestBuildInfo(t *testing.T) {
	info := BuildInfo{
		ModuleName: "test-module",
		Version:    "v1.0.0",
	}

	assert.Equal(t, "test-module", info.ModuleName)
	assert.Equal(t, "v1.0.0", info.Version)
}
