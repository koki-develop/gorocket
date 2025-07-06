package models

import (
	"errors"
	"testing"
)

func TestBuildTarget(t *testing.T) {
	target := BuildTarget{
		OS:   "linux",
		Arch: "amd64",
	}

	if target.OS != "linux" {
		t.Errorf("expected OS 'linux', got '%s'", target.OS)
	}

	if target.Arch != "amd64" {
		t.Errorf("expected Arch 'amd64', got '%s'", target.Arch)
	}
}

func TestBuildResult(t *testing.T) {
	target := BuildTarget{OS: "linux", Arch: "amd64"}
	result := BuildResult{
		Target:     target,
		BinaryPath: "/path/to/binary",
		Error:      errors.New("test error"),
	}

	if result.Target.OS != "linux" {
		t.Errorf("expected target OS 'linux', got '%s'", result.Target.OS)
	}

	if result.BinaryPath != "/path/to/binary" {
		t.Errorf("expected binary path '/path/to/binary', got '%s'", result.BinaryPath)
	}

	if result.Error == nil {
		t.Errorf("expected error, but got nil")
	}
}

func TestArchiveResult(t *testing.T) {
	target := BuildTarget{OS: "linux", Arch: "amd64"}
	result := ArchiveResult{
		Target:      target,
		ArchivePath: "/path/to/archive.tar.gz",
		Error:       nil,
	}

	if result.Target.OS != "linux" {
		t.Errorf("expected target OS 'linux', got '%s'", result.Target.OS)
	}

	if result.ArchivePath != "/path/to/archive.tar.gz" {
		t.Errorf("expected archive path '/path/to/archive.tar.gz', got '%s'", result.ArchivePath)
	}

	if result.Error != nil {
		t.Errorf("expected no error, but got %v", result.Error)
	}
}

func TestBuildInfo(t *testing.T) {
	info := BuildInfo{
		ModuleName: "test-module",
		Version:    "v1.0.0",
	}

	if info.ModuleName != "test-module" {
		t.Errorf("expected module name 'test-module', got '%s'", info.ModuleName)
	}

	if info.Version != "v1.0.0" {
		t.Errorf("expected version 'v1.0.0', got '%s'", info.Version)
	}
}