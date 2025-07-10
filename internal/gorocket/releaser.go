package gorocket

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/koki-develop/gorocket/internal/config"
	"github.com/koki-develop/gorocket/internal/git"
	"github.com/koki-develop/gorocket/internal/github"
)

// ReleaseParams contains options for the release command
type ReleaseParams struct {
	Draft bool
	Clean bool
}

// Releaser provides release functionality
type Releaser struct {
	configPath string
	git        *git.Client
	github     *github.Client
	builder    *Builder
}

// NewReleaser creates a new Releaser instance
func NewReleaser(configPath string, token string) (*Releaser, error) {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
		if token == "" {
			return nil, fmt.Errorf("GitHub token is required (use --token or GITHUB_TOKEN env var)")
		}
	}

	return &Releaser{
		configPath: configPath,
		git:        git.New(),
		github:     github.New(token),
		builder:    NewBuilder(configPath),
	}, nil
}

// Release creates a GitHub release
func (r *Releaser) Release(params ReleaseParams) error {
	// First build the binaries
	if err := r.builder.Build(BuildParams{Clean: params.Clean}); err != nil {
		return fmt.Errorf("failed to build: %w", err)
	}

	// Get repository info
	repo, err := r.git.GetRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository info: %w", err)
	}

	// Get version tag
	tag, err := r.git.GetHeadTag()
	if err != nil {
		return fmt.Errorf("failed to get git tag: %w", err)
	}

	// Check existing release
	release, err := r.github.GetReleaseByTag(github.GetReleaseByTagParams{
		Owner: repo.Owner,
		Repo:  repo.Name,
		Tag:   tag,
	})
	if err == nil && release != nil {
		fmt.Printf("Release %s already exists\n", tag)
		return nil
	}

	// Create new release
	release, err = r.github.CreateRelease(github.CreateReleaseParams{
		Owner: repo.Owner,
		Repo:  repo.Name,
		Tag:   tag,
		Name:  tag,
		Draft: params.Draft,
	})
	if err != nil {
		return fmt.Errorf("failed to create release: %w", err)
	}

	fmt.Printf("Created release %s\n", tag)

	// Upload assets
	distFiles, err := filepath.Glob("dist/*")
	if err != nil {
		return fmt.Errorf("failed to list dist files: %w", err)
	}

	for _, file := range distFiles {
		// Skip .rb files (Homebrew Formula)
		if filepath.Ext(file) == ".rb" {
			continue
		}

		asset := github.Asset{
			Name: filepath.Base(file),
			Path: file,
		}

		fmt.Printf("Uploading %s...\n", asset.Name)
		if err := r.github.UploadAsset(github.UploadAssetParams{
			Owner:     repo.Owner,
			Repo:      repo.Name,
			ReleaseID: release.GetID(),
			Asset:     asset,
		}); err != nil {
			return fmt.Errorf("failed to upload asset %s: %w", asset.Name, err)
		}
	}

	// Update Homebrew tap repository if configured
	cfg, err := config.LoadConfig(r.configPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if cfg.Brew.Repository != "" {
		if err := r.updateTapRepository(cfg.Brew.Repository); err != nil {
			return fmt.Errorf("failed to update tap repository: %w", err)
		}
	}

	return nil
}

// updateTapRepository updates Homebrew tap repository
func (r *Releaser) updateTapRepository(repository string) error {
	fmt.Printf("Updating tap repository %s...\n", repository)

	// Split tap repository owner and name
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid tap repository format: %s", repository)
	}
	tapOwner, tapRepo := parts[0], parts[1]

	// Get build info
	buildInfo, err := r.builder.getBuildInfo()
	if err != nil {
		return fmt.Errorf("failed to get build info: %w", err)
	}

	// Read Formula file
	moduleName := filepath.Base(buildInfo.Module)
	formulaPath := filepath.Join("dist", fmt.Sprintf("%s.rb", moduleName))
	content, err := os.ReadFile(formulaPath)
	if err != nil {
		return fmt.Errorf("failed to read formula file: %w", err)
	}

	// Update tap repository
	tapFormulaPath := fmt.Sprintf("Formula/%s.rb", moduleName)
	tapCommitMessage := fmt.Sprintf("Update %s to %s", moduleName, buildInfo.Version)

	if err := r.github.UpdateFile(github.UpdateFileParams{
		Owner:         tapOwner,
		Repo:          tapRepo,
		Path:          tapFormulaPath,
		Content:       string(content),
		CommitMessage: tapCommitMessage,
	}); err != nil {
		return fmt.Errorf("failed to update tap repository: %w", err)
	}

	fmt.Printf("Updated Formula in %s\n", repository)
	return nil
}
