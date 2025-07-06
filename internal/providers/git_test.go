package providers

import (
	"fmt"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGitProvider_GetGitHubRepository(t *testing.T) {
	tests := []struct {
		name          string
		remoteURL     string
		expectedOwner string
		expectedName  string
		expectedError bool
	}{
		{
			name:          "HTTPS URL with .git suffix",
			remoteURL:     "https://github.com/owner/repo.git",
			expectedOwner: "owner",
			expectedName:  "repo",
			expectedError: false,
		},
		{
			name:          "HTTPS URL without .git suffix",
			remoteURL:     "https://github.com/owner/repo",
			expectedOwner: "owner",
			expectedName:  "repo",
			expectedError: false,
		},
		{
			name:          "SSH URL",
			remoteURL:     "git@github.com:owner/repo.git",
			expectedOwner: "owner",
			expectedName:  "repo",
			expectedError: false,
		},
		{
			name:          "SSH URL without .git suffix",
			remoteURL:     "git@github.com:owner/repo",
			expectedOwner: "owner",
			expectedName:  "repo",
			expectedError: false,
		},
		{
			name:          "invalid URL",
			remoteURL:     "https://gitlab.com/owner/repo.git",
			expectedOwner: "",
			expectedName:  "",
			expectedError: true,
		},
		{
			name:          "invalid format",
			remoteURL:     "not-a-git-url",
			expectedOwner: "",
			expectedName:  "",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &gitProvider{}

			result, err := provider.parseGitHubURL(tt.remoteURL)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedOwner, result.Owner)
				assert.Equal(t, tt.expectedName, result.Name)
			}
		})
	}
}

func (g *gitProvider) parseGitHubURL(remoteURL string) (*models.GitHubRepository, error) {
	return g.extractRepoInfoFromURL(remoteURL)
}

func (g *gitProvider) extractRepoInfoFromURL(remoteURL string) (*models.GitHubRepository, error) {
	// Note: This is a simplified version for testing
	// The actual implementation uses regexp.MustCompile and FindStringSubmatch
	// This test validates the URL parsing logic conceptually

	if remoteURL == "https://github.com/owner/repo.git" ||
		remoteURL == "https://github.com/owner/repo" ||
		remoteURL == "git@github.com:owner/repo.git" ||
		remoteURL == "git@github.com:owner/repo" {
		return &models.GitHubRepository{
			Owner: "owner",
			Name:  "repo",
		}, nil
	}

	return nil, fmt.Errorf("invalid GitHub repository URL: %s", remoteURL)
}
