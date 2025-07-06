package providers

import (
	"fmt"
	"os/exec"
	"strings"
)

type GitProvider interface {
	GetCurrentVersion() (string, error)
}

type gitProvider struct{}

func NewGitProvider() GitProvider {
	return &gitProvider{}
}

func (g *gitProvider) GetCurrentVersion() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version from git tag: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}