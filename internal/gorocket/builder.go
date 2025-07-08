package gorocket

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/koki-develop/gorocket/internal/formula"
	"github.com/koki-develop/gorocket/internal/git"
	"github.com/koki-develop/gorocket/internal/util"
)

// BuildOptions contains options for the build command
type BuildOptions struct {
	Clean bool
}

// BuildInfo holds build information
type BuildInfo struct {
	Module  string
	Version string
}

// BuildResult represents a build result
type BuildResult struct {
	Binary  string
	OS      string
	Arch    string
	Version string
}

// Builder provides build functionality
type Builder struct {
	configPath string
	git        *git.Client
	formula    *formula.Client
}

// NewBuilder creates a new Builder instance
func NewBuilder(configPath string) *Builder {
	return &Builder{
		configPath: configPath,
		git:        git.New(),
		formula:    formula.New(),
	}
}

// Build executes cross-platform builds
func (b *Builder) Build(opts BuildOptions) error {
	// Get build info
	buildInfo, err := b.getBuildInfo()
	if err != nil {
		return err
	}

	// Load config file
	config, err := LoadConfig(b.configPath, map[string]interface{}{
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
		if err := b.generateFormula(config, buildInfo, results); err != nil {
			return fmt.Errorf("failed to generate formula: %w", err)
		}
	}

	return nil
}

// getBuildInfo retrieves module name and version
func (b *Builder) getBuildInfo() (*BuildInfo, error) {
	module, err := getModuleName()
	if err != nil {
		return nil, fmt.Errorf("failed to get module name: %w", err)
	}

	version, err := b.git.GetHeadTag()
	if err != nil {
		return nil, fmt.Errorf("failed to get version: %w", err)
	}

	return &BuildInfo{
		Module:  module,
		Version: version,
	}, nil
}

// generateFormula generates Homebrew Formula
func (b *Builder) generateFormula(config *Config, buildInfo *BuildInfo, results []*BuildResult) error {
	fmt.Println("Generating Homebrew Formula...")

	// Get repository info
	repo, err := b.git.GetRepository()
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
		file, err := os.Open(archivePath)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", archivePath, err)
		}
		sha256, err := util.CalculateSHA256(file)
		_ = file.Close()
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

	content, err := b.formula.Generate(f)
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

// buildBinary builds a single binary
func buildBinary(module, version string, target Target, ldflags string) (*BuildResult, error) {
	// Determine output file name
	binaryName := filepath.Base(module)
	if target.OS == "windows" {
		binaryName += ".exe"
	}
	binaryPath := filepath.Join("dist", binaryName)

	// Build command
	args := []string{"build", "-o", binaryPath}

	// Add ldflags if specified
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}

	// Execute command
	cmd := exec.Command("go", args...)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", target.OS),
		fmt.Sprintf("GOARCH=%s", target.Arch),
	)

	// Capture error output
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("go build failed: %w\nstderr: %s", err, stderr.String())
	}

	return &BuildResult{
		Binary:  binaryPath,
		OS:      target.OS,
		Arch:    target.Arch,
		Version: version,
	}, nil
}

// getModuleName retrieves module name from go.mod
func getModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			module := strings.TrimPrefix(line, "module ")
			return strings.TrimSpace(module), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	return "", fmt.Errorf("module name not found in go.mod")
}
