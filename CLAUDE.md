# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gorocket is a Go CLI tool for cross-platform binary building and packaging. It automates the process of building Go applications for multiple OS/architecture combinations and packages them into appropriate archives (tar.gz for Unix, zip for Windows).

## Common Commands

### Task-based Development (Recommended)
```bash
# Build the binary
task build

# Run tests
task test
task test-coverage    # With coverage report

# Code quality
task format    # Format with goimports
task lint      # Run golangci-lint
task formula-lint    # Validate generated Homebrew formula
```

### Direct Go Commands
```bash
# Run the application
go run main.go <command>

# Build the binary
go build -o gorocket .

# Run tests
go test ./...
go test ./internal/gorocket/... -v    # Run specific package tests
go test -coverprofile=coverage.out ./...  # With coverage

# Initialize project configuration
go run main.go init

# Build cross-platform binaries (requires git tag)
git tag v1.0.0  # Create a version tag first
go run main.go build
go run main.go build --clean  # Clean dist directory before building
go run main.go release --clean  # Clean dist directory before building and releasing
git tag -d v1.0.0  # Clean up test tag
```

### Testing Specific Components
```bash
# Test individual packages
go test ./internal/gorocket/... -v
go test ./internal/git/... -v
go test ./internal/github/... -v
go test ./internal/formula/... -v
go test ./internal/util/... -v

# Run tests with coverage
go test ./internal/gorocket -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run a single test
go test -run TestSpecificFunction ./internal/gorocket/
```

## Architecture

### High-Level Structure
The codebase follows a simplified architecture with direct implementations:

```
cmd/               # CLI commands (build.go, init.go, release.go, root.go, version.go)
internal/
├── config/        # Configuration management
├── gorocket/      # Core application logic
├── git/           # Git operations
├── github/        # GitHub API client
├── formula/       # Homebrew Formula generation
└── util/          # Utility functions
```

### Dependency Flow
- **CLI** (`cmd/`) → **Core** (`internal/gorocket/`)
- **Core** → **Config** (`internal/config/`) + **External** (`git/`, `github/`, `formula/`)
- Minimal abstractions - direct implementation preferred
- External dependencies only where necessary

### Key Components

**Config (`internal/config/`)**
- `config.go`: Configuration file management with Go template support
- Contains `Config` and `Target` structs, `LoadConfig` function

**Core (`internal/gorocket/`)** - Split into responsibility-focused structs:
- `initer.go`: `Initer` struct handles configuration initialization
- `builder.go`: `Builder` struct handles cross-platform builds, archive creation (tar.gz/zip), and Formula generation
- `releaser.go`: `Releaser` struct handles GitHub releases and asset uploads
- `config.go`: Contains embedded default configuration YAML
- `config_default.yaml`: Default configuration template file

**Git (`internal/git/`)**
- `git.go`: Git operations (tag retrieval, repository info)
- Uses direct `exec.Command` calls

**GitHub (`internal/github/`)**
- `client.go`: GitHub API operations using go-github v66
- Uses parameter structs for all methods (GetReleaseByTagParams, CreateReleaseParams, UploadAssetParams, UpdateFileParams)
- Release creation and asset upload with owner/repo per-method specification

**Formula (`internal/formula/`)**
- `formula.go`: Homebrew Formula generation

**Utility (`internal/util/`)**
- `hash.go`: SHA256 calculation for archives and files
- `hash_test.go`: Test cases using testify framework

### Command Pattern
Each CLI command is defined as a package-level variable with global flags:
```go
var (
    flagBuildClean bool
)

var buildCmd = &cobra.Command{
    Use:   "build",
    Short: "Build binaries for multiple platforms",
    RunE: func(cmd *cobra.Command, args []string) error {
        builder := gorocket.NewBuilder(".gorocket.yml")
        return builder.Build(gorocket.BuildParams{
            Clean: flagBuildClean,
        })
    },
}

func init() {
    buildCmd.Flags().BoolVar(&flagBuildClean, "clean", false, "Remove dist directory before building")
    rootCmd.AddCommand(buildCmd)
}
```

For the release command, token validation happens during initialization:
```go
var (
    flagReleaseDraft       bool
    flagReleaseGitHubToken string
    flagReleaseClean       bool
)

var releaseCmd = &cobra.Command{
    Use:   "release",
    Short: "Create a GitHub release with built artifacts",
    RunE: func(cmd *cobra.Command, args []string) error {
        releaser, err := gorocket.NewReleaser(".gorocket.yml", flagReleaseGitHubToken)
        if err != nil {
            return err
        }
        return releaser.Release(gorocket.ReleaseParams{
            Draft: flagReleaseDraft,
            Clean: flagReleaseClean,
        })
    },
}

func init() {
    rootCmd.AddCommand(releaseCmd)
    releaseCmd.Flags().StringVar(&flagReleaseGitHubToken, "github-token", "", "GitHub token (defaults to GITHUB_TOKEN env var)")
    releaseCmd.Flags().BoolVar(&flagReleaseDraft, "draft", false, "Create a draft release")
    releaseCmd.Flags().BoolVar(&flagReleaseClean, "clean", false, "Remove dist directory before building")
}
```

## Configuration

The tool uses `.gorocket.yml` for configuration:
- `build.targets`: Array of OS/architecture combinations
- `build.ldflags`: Optional linker flags
- `brew.repository`: Optional Homebrew tap repository (owner/name)
- Supports Go templates for dynamic values: `{{.Version}}`, `{{.Module}}`
- Config loading accepts `map[string]any` for template data
- Configuration management is handled by `internal/config` package
- Default configuration YAML is embedded in `internal/gorocket/config.go`

## Build Process Requirements

- **Requires a git tag on HEAD** for version information
- Uses `git describe --tags --exact-match HEAD` to get version
- Builds to `dist/` directory
- Creates archives with naming: `{module}_{version}_{os}_{arch}.{ext}`
- Generates Homebrew Formula when `brew.repository` is configured

## Available Commands

- `init`: Create a default .gorocket.yml configuration file
- `build`: Build binaries for multiple platforms (requires git tag)
  - `--clean`: Remove dist directory before building
- `release`: Create GitHub release with built artifacts
  - `--token`: GitHub token (defaults to GITHUB_TOKEN env var)
  - `--draft`: Create a draft release
  - `--clean`: Remove dist directory before building
- `version`: Display version information

## Code Quality

- **goimports** for formatting (`task format`)
- **golangci-lint** for static analysis (`task lint`)
  - Configured with `.golangci.yml` to enable `unparam` linter for unused parameter detection
  - Includes default linters: `errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`
- **testify** for test assertions and better test readability
- Direct implementations without excessive abstractions
- Comments in English

## Commit Convention

The project follows conventional commit format:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks (build, dependencies, tooling)
- `refactor:` - Code restructuring without functional changes
- `improve:` - Enhancements to existing features
- `test:` - Adding or modifying tests

## Task Management

The project uses go-task for common development workflows:
- Prefer `task <command>` over direct go commands
- Use `task --list` to see all available tasks
- Essential tasks: `build`, `test`, `test-coverage`, `format`, `lint`
- Additional tasks: `mocks` (generate test mocks), `formula-lint` (validate Homebrew formulas)

## Release Workflow

1. Create a git tag: `git tag v1.0.0`
2. Run release command: `gorocket release` or `gorocket release --token <GITHUB_TOKEN>` or `gorocket release --clean`
3. This will:
   - Build all configured targets
   - Create GitHub release
   - Upload artifacts
   - Update Homebrew tap (if configured)

## Environment Variables

- `GITHUB_TOKEN`: Required for release command
- `GITHUB_REPOSITORY`: Optional, overrides git remote detection

## Testing Framework

The project uses testify for test assertions:
- `github.com/stretchr/testify/assert` for assertions
- Table-driven test pattern for comprehensive coverage
- Test files follow `*_test.go` naming convention

Example test pattern:
```go
func Test_CalculateSHA256(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        // test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            reader := strings.NewReader(tt.input)
            result, err := CalculateSHA256(reader)
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

## Development Environment

The project uses `mise` for tool version management:
- Go 1.24.3
- golangci-lint 2.1.6
- goreleaser 2.9.0
- goimports 0.34.0

CI/CD is handled by GitHub Actions with automated testing, building, and linting.

## Architecture Documentation

For detailed architecture design, refer to `ARCHITECTURE.md` (written in Japanese) which contains comprehensive design principles, package structure, and data flow documentation.
