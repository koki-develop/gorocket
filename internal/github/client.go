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
	owner  string
	repo   string
}

// Release represents a GitHub release
type Release struct {
	ID    int64
	Tag   string
	Name  string
	Body  string
	Draft bool
}

// Asset represents a release asset
type Asset struct {
	Name string
	Path string
}

// New creates a new GitHub client
func New(token, owner, repo string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}
}

// GetReleaseByTag retrieves a release by tag name
func (c *Client) GetReleaseByTag(tag string) (*Release, error) {
	ctx := context.Background()
	release, resp, err := c.client.Repositories.GetReleaseByTag(ctx, c.owner, c.repo, tag)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get release: %w", err)
	}

	return &Release{
		ID:    release.GetID(),
		Tag:   release.GetTagName(),
		Name:  release.GetName(),
		Body:  release.GetBody(),
		Draft: release.GetDraft(),
	}, nil
}

// CreateRelease creates a new release
func (c *Client) CreateRelease(release *Release) error {
	ctx := context.Background()

	githubRelease := &github.RepositoryRelease{
		TagName: github.String(release.Tag),
		Name:    github.String(release.Name),
		Body:    github.String(release.Body),
		Draft:   github.Bool(release.Draft),
	}

	created, _, err := c.client.Repositories.CreateRelease(ctx, c.owner, c.repo, githubRelease)
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	// Set the created release ID
	release.ID = created.GetID()

	return nil
}

// UploadAsset uploads an asset to a release
func (c *Client) UploadAsset(releaseID int64, asset Asset) error {
	ctx := context.Background()

	file, err := os.Open(asset.Path)
	if err != nil {
		return fmt.Errorf("failed to open asset file %s: %w", asset.Path, err)
	}
	defer func() { _ = file.Close() }()

	uploadOptions := &github.UploadOptions{
		Name: asset.Name,
	}

	_, _, err = c.client.Repositories.UploadReleaseAsset(ctx, c.owner, c.repo, releaseID, uploadOptions, file)
	if err != nil {
		return fmt.Errorf("failed to upload asset %s: %w", asset.Name, err)
	}

	return nil
}

// UpdateFile creates or updates a file in the repository
func (c *Client) UpdateFile(path, content, message string) error {
	ctx := context.Background()

	// Get SHA of existing file
	existingFile, _, resp, err := c.client.Repositories.GetContents(ctx, c.owner, c.repo, path, nil)

	var sha *string
	if err == nil && existingFile != nil {
		sha = existingFile.SHA
	} else if resp != nil && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to check existing file: %w", err)
	}

	// Create or update file
	opts := &github.RepositoryContentFileOptions{
		Message: github.String(message),
		Content: []byte(content),
		SHA:     sha,
	}

	_, _, err = c.client.Repositories.CreateFile(ctx, c.owner, c.repo, path, opts)
	if err != nil {
		return fmt.Errorf("failed to update file: %w", err)
	}

	return nil
}
