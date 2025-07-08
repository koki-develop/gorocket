package git

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	httpsPattern = regexp.MustCompile(`^https://github\.com/([^/]+)/([^/]+?)(?:\.git)?/?$`)
	sshPattern   = regexp.MustCompile(`^git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`)
)

// Repository represents a GitHub repository
type Repository struct {
	Owner string
	Name  string
}

// Client provides Git operations
type Client struct{}

// New creates a new Git client
func New() *Client {
	return &Client{}
}

// GetHeadTag retrieves the current HEAD tag
func (c *Client) GetHeadTag() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--exact-match", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git tag: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetRepository retrieves GitHub repository information
func (c *Client) GetRepository() (*Repository, error) {
	// Prefer environment variable
	if env := os.Getenv("GITHUB_REPOSITORY"); env != "" {
		parts := strings.SplitN(env, "/", 2)
		if len(parts) == 2 {
			return &Repository{
				Owner: parts[0],
				Name:  parts[1],
			}, nil
		}
		return nil, fmt.Errorf("invalid GITHUB_REPOSITORY format: %s (expected owner/repo)", env)
	}

	// Get from git remote
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git remote origin: %w", err)
	}

	remoteURL := strings.TrimSpace(string(output))

	// Parse URL
	var matches []string
	if httpsPattern.MatchString(remoteURL) {
		matches = httpsPattern.FindStringSubmatch(remoteURL)
	} else if sshPattern.MatchString(remoteURL) {
		matches = sshPattern.FindStringSubmatch(remoteURL)
	}

	if len(matches) != 3 {
		return nil, fmt.Errorf("invalid GitHub repository URL: %s", remoteURL)
	}

	return &Repository{
		Owner: matches[1],
		Name:  matches[2],
	}, nil
}
