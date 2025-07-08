package github

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v66/github"
	"golang.org/x/oauth2"
)

// Client is a GitHub API client
type Client struct {
	client *github.Client
}

// GetReleaseByTagParams represents parameters for GetReleaseByTag
type GetReleaseByTagParams struct {
	Owner string
	Repo  string
	Tag   string
}

// CreateReleaseParams represents parameters for CreateRelease
type CreateReleaseParams struct {
	Owner string
	Repo  string
	Tag   string
	Name  string
	Draft bool
}

// UploadAssetParams represents parameters for UploadAsset
type UploadAssetParams struct {
	Owner     string
	Repo      string
	ReleaseID int64
	Asset     Asset
}

// UpdateFileParams represents parameters for UpdateFile
type UpdateFileParams struct {
	Owner         string
	Repo          string
	Path          string
	Content       string
	CommitMessage string
}

// Asset represents a release asset
type Asset struct {
	Name string
	Path string
}

// New creates a new GitHub client
func New(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
	}
}

// GetReleaseByTag retrieves a release by tag name
func (c *Client) GetReleaseByTag(params GetReleaseByTagParams) (*github.RepositoryRelease, error) {
	ctx := context.Background()
	release, resp, err := c.client.Repositories.GetReleaseByTag(ctx, params.Owner, params.Repo, params.Tag)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	return release, nil
}

// CreateRelease creates a new release
func (c *Client) CreateRelease(params CreateReleaseParams) (*github.RepositoryRelease, error) {
	ctx := context.Background()

	githubRelease := &github.RepositoryRelease{
		TagName: github.String(params.Tag),
		Name:    github.String(params.Name),
		Draft:   github.Bool(params.Draft),
	}

	created, _, err := c.client.Repositories.CreateRelease(ctx, params.Owner, params.Repo, githubRelease)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	return created, nil
}

// UploadAsset uploads an asset to a release
func (c *Client) UploadAsset(params UploadAssetParams) error {
	ctx := context.Background()

	file, err := os.Open(params.Asset.Path)
	if err != nil {
		return fmt.Errorf("failed to open asset file %s: %w", params.Asset.Path, err)
	}
	defer func() { _ = file.Close() }()

	uploadOptions := &github.UploadOptions{
		Name: params.Asset.Name,
	}

	_, _, err = c.client.Repositories.UploadReleaseAsset(ctx, params.Owner, params.Repo, params.ReleaseID, uploadOptions, file)
	if err != nil {
		return fmt.Errorf("failed to upload asset %s: %w", params.Asset.Name, err)
	}

	return nil
}

// UpdateFile creates or updates a file in the repository
func (c *Client) UpdateFile(params UpdateFileParams) error {
	ctx := context.Background()

	// Get SHA of existing file
	existingFile, _, resp, err := c.client.Repositories.GetContents(ctx, params.Owner, params.Repo, params.Path, nil)

	var sha *string
	if err == nil && existingFile != nil {
		sha = existingFile.SHA
	} else if resp != nil && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to check existing file: %w", err)
	}

	// Create or update file
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(params.CommitMessage),
		Content: []byte(params.Content),
		SHA:     sha,
	}

	_, _, err = c.client.Repositories.CreateFile(ctx, params.Owner, params.Repo, params.Path, opts)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}
