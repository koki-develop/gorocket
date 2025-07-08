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
git tag -d v1.0.0  # Clean up test tag
```

### Testing Specific Components
```bash
# Test individual packages
go test ./internal/gorocket/... -v
go test ./internal/git/... -v
go test ./internal/github/... -v
go test ./internal/formula/... -v

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
cmd/gorocket/      # CLI entry point
internal/
├── gorocket/      # Core application logic
├── git/           # Git operations
├── github/        # GitHub API client
└── formula/       # Homebrew Formula generation
```

### Dependency Flow
- **CLI** (`cmd/gorocket/`) → **Core** (`internal/gorocket/`)
- **Core** → **External** (`git/`, `github/`, `formula/`)
- Minimal abstractions - direct implementation preferred
- External dependencies only where necessary

### Key Components

**Core (`internal/gorocket/`)**
- `gorocket.go`: Main application logic, orchestrates all operations
- `config.go`: Configuration management with Go template support
- `build.go`: Cross-platform build logic
- `archive.go`: Archive creation (tar.gz/zip)

**Git (`internal/git/`)**
- `git.go`: Git operations (tag retrieval, repository info)
- Uses direct `exec.Command` calls

**GitHub (`internal/github/`)**
- `client.go`: GitHub API operations using go-github v66
- Release creation and asset upload

**Formula (`internal/formula/`)**
- `formula.go`: Homebrew Formula generation
- SHA256 calculation for archives

### Command Pattern
Each CLI command follows this simplified pattern:
```go
func newBuildCommand(app *gorocket.GoRocket) *cobra.Command {
    var clean bool
    
    cmd := &cobra.Command{
        Use:   "build",
        Short: "Build binaries for multiple platforms",
        RunE: func(cmd *cobra.Command, args []string) error {
            return app.Build(gorocket.BuildOptions{
                Clean: clean,
            })
        },
    }
    
    cmd.Flags().BoolVar(&clean, "clean", false, "Remove dist directory before building")
    return cmd
}
```

## Configuration

The tool uses `.gorocket.yml` for configuration:
- `build.targets`: Array of OS/architecture combinations
- `build.ldflags`: Optional linker flags
- `brew.repository`: Optional Homebrew tap repository (owner/name)
- Supports Go templates for dynamic values: `{{.Version}}`, `{{.Module}}`

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
- `version`: Display version information

## Code Quality

- **goimports** for formatting (`task format`)
- **golangci-lint** for static analysis (`task lint`)
- Direct implementations without excessive abstractions
- Comments in English

## Task Management

The project uses go-task for common development workflows:
- Prefer `task <command>` over direct go commands
- Use `task --list` to see all available tasks
- Essential tasks: `build`, `test`, `test-coverage`, `format`, `lint`

## Release Workflow

1. Create a git tag: `git tag v1.0.0`
2. Run release command: `gorocket release`
3. This will:
   - Build all configured targets
   - Create GitHub release
   - Upload artifacts
   - Update Homebrew tap (if configured)

## Environment Variables

- `GITHUB_TOKEN`: Required for release command
- `GITHUB_REPOSITORY`: Optional, overrides git remote detection