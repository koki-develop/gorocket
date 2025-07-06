package providers

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/koki-develop/gorocket/internal/models"
)

var (
	httpsPattern = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	sshPattern   = regexp.MustCompile(`^git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`)
)

type GitProvider interface {
	GetCurrentVersion() (string, error)
	GetGitHubRepository() (*models.GitHubRepository, error)
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

func (g *gitProvider) GetGitHubRepository() (*models.GitHubRepository, error) {
	// Check GITHUB_REPOSITORY environment variable first
	if repo := os.Getenv("GITHUB_REPOSITORY"); repo != "" {
		parts := strings.SplitN(repo, "/", 2)
		if len(parts) == 2 {
			return &models.GitHubRepository{
				Owner: parts[0],
				Name:  parts[1],
			}, nil
		}
	}

	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git remote origin: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))

	var matches []string
	if httpsPattern.MatchString(remoteURL) {
		matches = httpsPattern.FindStringSubmatch(remoteURL)
	} else if sshPattern.MatchString(remoteURL) {
		matches = sshPattern.FindStringSubmatch(remoteURL)
	}

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid GitHub repository URL: %s", remoteURL)
	}

	return &models.GitHubRepository{
		Owner: matches[1],
		Name:  matches[2],
	}, nil
}
