package gorocket

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/koki-develop/gorocket/internal/formula"
	"github.com/koki-develop/gorocket/internal/git"
	"github.com/koki-develop/gorocket/internal/github"
)

// GoRocket provides the main application logic
type GoRocket struct {
	configPath string
	git        *git.Client
	github     *github.Client
}

// BuildOptions contains options for the build command
type BuildOptions struct {
	Clean bool
}

// ReleaseOptions contains options for the release command
type ReleaseOptions struct {
	Token string
	Draft bool
}

// New creates a new GoRocket instance
func New() *GoRocket {
	return &GoRocket{
		configPath: ".gorocket.yml",
		git:        git.New(),
	}
}

// Init creates a default configuration file
func (g *GoRocket) Init() error {
	// Error if config file already exists
	if _, err := os.Stat(g.configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", g.configPath)
	}

	// Create default config file
	if err := SaveDefaultConfig(g.configPath); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}

	fmt.Printf("Created %s\n", g.configPath)
	return nil
}

// Build executes cross-platform builds
func (g *GoRocket) Build(opts BuildOptions) error {
	// Get build info
	buildInfo, err := g.getBuildInfo()
	if err != nil {
		return err
	}

	// Load config file
	config, err := LoadConfig(g.configPath, map[string]interface{}{
		"Version": buildInfo.Version,
		"Module":  buildInfo.Module,
	})
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Prepare dist directory
	distDir := "dist"
	if opts.Clean {
		if err := os.RemoveAll(distDir); err != nil {
			return fmt.Errorf("failed to clean dist directory: %w", err)
		}
	}

	if err := os.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("failed to create dist directory: %w", err)
	}

	// Build each target
	var results []*BuildResult
	for _, target := range config.Build.Targets {
		fmt.Printf("Building %s/%s...\n", target.OS, target.Arch)

		result, err := buildBinary(buildInfo.Module, buildInfo.Version, target, config.Build.Ldflags)
		if err != nil {
			return fmt.Errorf("failed to build %s/%s: %w", target.OS, target.Arch, err)
		}

		// Create archive
		archivePath, err := createArchive(result)
		if err != nil {
			return fmt.Errorf("failed to create archive: %w", err)
		}

		fmt.Printf("Created %s\n", archivePath)

		// Remove binary file
		if err := os.Remove(result.Binary); err != nil {
			return fmt.Errorf("failed to remove binary: %w", err)
		}

		results = append(results, result)
	}

	// Generate Homebrew Formula if configured
	if config.Brew.Repository != "" {
		if err := g.generateFormula(config, buildInfo, results); err != nil {
			return fmt.Errorf("failed to generate formula: %w", err)
		}
	}

	return nil
}

// Release creates a GitHub release
func (g *GoRocket) Release(opts ReleaseOptions) error {
	// First build the binaries
	if err := g.Build(BuildOptions{Clean: false}); err != nil {
		return fmt.Errorf("failed to build: %w", err)
	}

	// Get GitHub token
	token := opts.Token
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
		if token == "" {
			return fmt.Errorf("GitHub token is required (use --token or GITHUB_TOKEN env var)")
		}
	}

	// Get repository info
	repo, err := g.git.GetRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository info: %w", err)
	}

	// Initialize GitHub client
	g.github = github.New(token, repo.Owner, repo.Name)

	// Get version tag
	tag, err := g.git.GetHeadTag()
	if err != nil {
		return fmt.Errorf("failed to get git tag: %w", err)
	}

	// Check existing release
	release, err := g.github.GetRelease(tag)
	if err == nil && release != nil {
		fmt.Printf("Release %s already exists\n", tag)
		return nil
	}

	// Create new release
	release = &github.Release{
		Tag:   tag,
		Name:  tag,
		Body:  fmt.Sprintf("Release %s", tag),
		Draft: opts.Draft,
	}

	if err := g.github.CreateRelease(release); err != nil {
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
		if err := g.github.UploadAsset(release.ID, asset); err != nil {
			return fmt.Errorf("failed to upload asset %s: %w", asset.Name, err)
		}
	}

	// Update Homebrew tap repository if configured
	config, err := LoadConfig(g.configPath, nil)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if config.Brew.Repository != "" {
		if err := g.updateTapRepository(config.Brew.Repository); err != nil {
			return fmt.Errorf("failed to update tap repository: %w", err)
		}
	}

	return nil
}

// Version returns version information
func (g *GoRocket) Version() (string, error) {
	// Return embedded version if available
	if version != "" && version != "dev" {
		return version, nil
	}

	// Try to get version from git tag in development
	tag, err := g.git.GetHeadTag()
	if err != nil {
		return "dev", nil
	}
	return tag, nil
}

// version is set by -ldflags during build
var version = "dev"

// BuildInfo holds build information
type BuildInfo struct {
	Module  string
	Version string
}

// getBuildInfo retrieves module name and version
func (g *GoRocket) getBuildInfo() (*BuildInfo, error) {
	module, err := getModuleName()
	if err != nil {
		return nil, fmt.Errorf("failed to get module name: %w", err)
	}

	version, err := g.git.GetHeadTag()
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return &BuildInfo{
		Module:  module,
		Version: version,
	}, nil
}

// generateFormula generates Homebrew Formula
func (g *GoRocket) generateFormula(config *Config, buildInfo *BuildInfo, results []*BuildResult) error {
	fmt.Println("Generating Homebrew Formula...")

	// Get repository info
	repo, err := g.git.GetRepository()
	if err != nil {
		return fmt.Errorf("failed to get repository info: %w", err)
	}

	// Create artifact information
	var artifacts []formula.Artifact
	for _, result := range results {
		// Determine archive name
		var archiveName string
		if result.OS == "windows" {
			archiveName = fmt.Sprintf("%s_%s_%s_%s.zip", buildInfo.Module, result.Version, result.OS, result.Arch)
		} else {
			archiveName = fmt.Sprintf("%s_%s_%s_%s.tar.gz", buildInfo.Module, result.Version, result.OS, result.Arch)
		}
		archivePath := filepath.Join("dist", archiveName)

		// Calculate SHA256
		sha256, err := formula.CalculateSHA256(archivePath)
		if err != nil {
			return fmt.Errorf("failed to calculate SHA256 for %s: %w", archivePath, err)
		}

		// Build URL
		url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s",
			repo.Owner, repo.Name, result.Version, archiveName)

		artifacts = append(artifacts, formula.Artifact{
			OS:     result.OS,
			Arch:   result.Arch,
			URL:    url,
			SHA256: sha256,
		})
	}

	// Generate Formula
	f := &formula.Formula{
		Name:      filepath.Base(buildInfo.Module),
		Version:   buildInfo.Version,
		Artifacts: artifacts,
	}

	content, err := formula.Generate(f)
	if err != nil {
		return fmt.Errorf("failed to generate formula: %w", err)
	}

	// Save to file
	formulaPath := filepath.Join("dist", fmt.Sprintf("%s.rb", f.Name))
	if err := os.WriteFile(formulaPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write formula file: %w", err)
	}

	fmt.Printf("Created %s\n", formulaPath)
	return nil
}

// updateTapRepository updates Homebrew tap repository
func (g *GoRocket) updateTapRepository(repository string) error {
	fmt.Printf("Updating tap repository %s...\n", repository)

	// Split tap repository owner and name
	parts := strings.Split(repository, "/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid tap repository format: %s", repository)
	}
	tapOwner, tapRepo := parts[0], parts[1]

	// Create GitHub client for tap repository
	// Use current token (passed from ReleaseOptions)
	token := os.Getenv("GITHUB_TOKEN")
	tapClient := github.New(token, tapOwner, tapRepo)

	// Get build info
	buildInfo, err := g.getBuildInfo()
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
	if err := formula.UpdateTapRepository(tapClient, string(content), moduleName, buildInfo.Version); err != nil {
		return fmt.Errorf("failed to update tap repository: %w", err)
	}

	fmt.Printf("Updated Formula in %s\n", repository)
	return nil
}
