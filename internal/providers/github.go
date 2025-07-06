package providers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v50/github"
	"github.com/koki-develop/gorocket/internal/models"
	"golang.org/x/oauth2"
)

type GitHubProvider interface {
	ReleaseExists(ctx context.Context, repo *models.GitHubRepository, tagName string) (bool, error)
	CreateRelease(ctx context.Context, repo *models.GitHubRepository, tagName string) (*github.RepositoryRelease, error)
	UploadAssets(ctx context.Context, repo *models.GitHubRepository, release *github.RepositoryRelease, assets []models.ReleaseAsset) error
	UpdateTapRepository(ctx context.Context, tapRepo *models.Repository, formula string, moduleName, version string) error
}

type gitHubProvider struct {
	client     *github.Client
	fsProvider FileSystemProvider
}

func NewGitHubProvider(token string, fsProvider FileSystemProvider) GitHubProvider {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	return &gitHubProvider{
		client:     github.NewClient(tc),
		fsProvider: fsProvider,
	}
}

func (g *gitHubProvider) ReleaseExists(ctx context.Context, repo *models.GitHubRepository, tagName string) (bool, error) {
	_, resp, err := g.client.Repositories.GetReleaseByTag(ctx, repo.Owner, repo.Name, tagName)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check release existence: %w", err)
	}
	return true, nil
}

func (g *gitHubProvider) CreateRelease(ctx context.Context, repo *models.GitHubRepository, tagName string) (*github.RepositoryRelease, error) {
	release := &github.RepositoryRelease{
		TagName: github.String(tagName),
		Name:    github.String(tagName),
		Draft:   github.Bool(false),
	}

	createdRelease, _, err := g.client.Repositories.CreateRelease(ctx, repo.Owner, repo.Name, release)
	if err != nil {
		return nil, fmt.Errorf("failed to create release: %w", err)
	}

	return createdRelease, nil
}

func (g *gitHubProvider) UploadAssets(ctx context.Context, repo *models.GitHubRepository, release *github.RepositoryRelease, assets []models.ReleaseAsset) error {
	for _, asset := range assets {
		file, err := os.Open(asset.Path)
		if err != nil {
			return fmt.Errorf("failed to open asset file %s: %w", asset.Path, err)
		}
		defer func() { _ = file.Close() }()

		uploadOptions := &github.UploadOptions{
			Name: asset.Name,
		}

		_, _, err = g.client.Repositories.UploadReleaseAsset(ctx, repo.Owner, repo.Name, *release.ID, uploadOptions, file)
		if err != nil {
			return fmt.Errorf("failed to upload asset %s: %w", asset.Name, err)
		}
	}

	return nil
}

func (g *gitHubProvider) UpdateTapRepository(ctx context.Context, tapRepo *models.Repository, formula string, moduleName, version string) error {
	formulaPath := fmt.Sprintf("Formula/%s.rb", moduleName)

	existingFile, _, resp, err := g.client.Repositories.GetContents(ctx, tapRepo.Owner, tapRepo.Name, formulaPath, nil)

	var sha *string
	if err == nil && existingFile != nil {
		sha = existingFile.SHA
	} else if resp != nil && resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to check existing formula: %w", err)
	}

	commitMessage := fmt.Sprintf("Update %s to %s", moduleName, version)

	repositoryContentFileOptions := &github.RepositoryContentFileOptions{
		Message: github.String(commitMessage),
		Content: []byte(formula),
		SHA:     sha,
	}

	_, _, err = g.client.Repositories.CreateFile(ctx, tapRepo.Owner, tapRepo.Name, formulaPath, repositoryContentFileOptions)
	if err != nil {
		return fmt.Errorf("failed to update tap repository: %w", err)
	}

	return nil
}
