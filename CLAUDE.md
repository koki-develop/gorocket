# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gorocket is a Go CLI tool for cross-platform binary building and packaging. It automates the process of building Go applications for multiple OS/architecture combinations and packages them into appropriate archives (tar.gz for Unix, zip for Windows).

## Common Commands

### Development
```bash
# Run the application
go run main.go <command>

# Build the binary
go build -o gorocket .

# Run tests
go test ./...
go test ./internal/services/... -v    # Run specific package tests
go test -coverprofile=coverage.out ./...  # With coverage

# Initialize project configuration
go run main.go init

# Build cross-platform binaries (requires git tag)
git tag v1.0.0  # Create a version tag first
go run main.go build
git tag -d v1.0.0  # Clean up test tag
```

### Testing Specific Components
```bash
# Test individual packages
go test ./internal/models/... -v
go test ./internal/providers/... -v
go test ./internal/services/... -v

# Run tests with coverage
go test ./internal/services -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Architecture

### High-Level Structure
The codebase follows a clean architecture pattern with dependency injection:

```
cmd/          # CLI commands (uses services)
internal/
├── models/   # Data structures and types
├── providers/ # External resource access (files, git, commands)
└── services/ # Business logic (orchestrates providers)
```

### Dependency Flow
- **Commands** depend on **Services**
- **Services** depend on **Providers** 
- **Models** are used by all layers (no dependencies)
- All dependencies are injected via constructor functions

### Key Components

**Models (`internal/models/`)**
- `Config`: Configuration structure matching `.gorocket.yaml`
- `BuildTarget`, `BuildResult`, `ArchiveResult`: Build process data structures
- `BuildInfo`: Module name and version information

**Providers (`internal/providers/`)**
- `FileSystemProvider`: File operations abstraction
- `GitProvider`: Git operations (version retrieval)
- `CommandProvider`: External command execution (go build)
- `ConfigProvider`: Configuration file management

**Services (`internal/services/`)**
- `VersionService`: Aggregates module name and git version
- `BuilderService`: Cross-platform build orchestration
- `ArchiverService`: Archive creation (tar.gz/zip)
- `ConfigService`: Configuration management

### Command Pattern
Each CLI command follows this pattern:
```go
type ExampleCommand struct {
    service1 services.Service1
    service2 services.Service2
}

func NewExampleCommand() *cobra.Command {
    // Create providers
    provider1 := providers.NewProvider1()
    
    // Create services with providers
    service1 := services.NewService1(provider1)
    
    // Create command with services
    cmd := &ExampleCommand{service1: service1}
    
    return &cobra.Command{RunE: cmd.run}
}
```

## Configuration

The tool uses `.gorocket.yaml` for configuration:
- `build.targets`: Array of OS/architecture combinations to build
- Default config is embedded in `providers/config_default.yaml`
- Configuration is managed through `ConfigService`

## Testing Strategy

- **Unit tests** use mocks from `providers/mocks/` and `services/mocks/`
- **No file system side effects** in tests - all I/O is mocked
- **Test coverage** focuses on business logic rather than 100% coverage
- **Integration testing** through CLI commands with temporary git tags

## Build Process Requirements

- Requires a git tag on HEAD for version information
- Uses `git describe --tags --exact-match HEAD` to get version
- Builds to `dist/` directory (must be empty before build)
- Creates archives with consistent naming: `{module}_{version}_{os}_{arch}.{ext}`