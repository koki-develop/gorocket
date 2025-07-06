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

# Generate mocks
task mocks

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
go test ./internal/services/... -v    # Run specific package tests
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
go test ./internal/models/... -v
go test ./internal/providers/... -v
go test ./internal/services/... -v

# Run tests with coverage
go test ./internal/services -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# Run a single test
go test -run TestBuilderService_BuildTargets ./internal/services/
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
- `FormulaService`: Homebrew Formula generation for tap repositories

### Command Pattern
Each CLI command follows this pattern:
```go
type ExampleCommand struct {
    service1 services.Service1
    service2 services.Service2
    flagExample bool  // Command flags use flag prefix
}

func NewExampleCommand() *cobra.Command {
    // Create providers
    provider1 := providers.NewProvider1()
    
    // Create services with providers
    service1 := services.NewService1(provider1)
    
    // Create command with services
    cmd := &ExampleCommand{service1: service1}
    
    cobraCmd := &cobra.Command{RunE: cmd.run}
    cobraCmd.Flags().BoolVar(&cmd.flagExample, "example", false, "Example flag")
    
    return cobraCmd
}
```

## Configuration

The tool uses `.gorocket.yaml` for configuration:
- `build.targets`: Array of OS/architecture combinations to build
- `brew.repository`: Optional Homebrew tap repository configuration (owner/name)
- Default config is embedded in `providers/config_default.yaml`
- Configuration is managed through `ConfigService`

## Testing Strategy

- **Unit tests** use auto-generated mocks from `internal/providers/mocks/` and `internal/services/mocks/`
- **Mock generation** via `mockery` command using `.mockery.yml` configuration
- **testify/assert** for assertions and **testify/mock** for mocking
- **No file system side effects** in tests - all I/O is mocked
- **Test coverage** focuses on business logic rather than 100% coverage
- **Integration testing** through CLI commands with temporary git tags

### Mock Management
```bash
# Regenerate mocks when interfaces change
task mocks    # Preferred
mockery       # Direct command

# Generated mocks use modern EXPECT() API:
mockFS := mocks.NewMockFileSystemProvider(t)
mockFS.EXPECT().ReadFile(".gorocket.yaml").Return(data, nil)
```

## Build Process Requirements

- Requires a git tag on HEAD for version information
- Uses `git describe --tags --exact-match HEAD` to get version
- Builds to `dist/` directory (must be empty before build unless --clean flag is used)
- Creates archives with consistent naming: `{module}_{version}_{os}_{arch}.{ext}`
- Generates Homebrew Formula (.rb file) when `brew.repository` is configured

## Available Commands

- `init`: Create a default .gorocket.yaml configuration file
- `build`: Build binaries for multiple platforms (requires git tag)
  - `--clean`: Remove dist directory before building
- `version`: Display version information

## Code Quality

- **goimports** for import formatting and organization (`task format`)
- **golangci-lint** for static analysis (`task lint`)
- **testify v1.10.0** for modern testing patterns
- **mockery v3.5.0** for automated mock generation (`task mocks`)
- All Provider interfaces have auto-generated mocks for testing

## Task Management

The project uses go-task for common development workflows. All tasks are defined in `Taskfile.yml`:
- Prefer `task <command>` over direct go commands
- Use `task --list` to see all available tasks
- Essential tasks: `build`, `test`, `test-coverage`, `mocks`, `format`, `lint`, `formula-lint`

## Homebrew Formula Generation

The tool can generate Homebrew Formula files for tap repositories:

- **Configuration**: Add `brew.repository` section to `.gorocket.yaml` with owner/name
- **Formula Generation**: Automatically creates `.rb` file during build when brew config exists
- **Quality Assurance**: Use `task formula-lint` to validate generated Formula files
- **Ruby Compliance**: Generated formulas include Sorbet typing and follow Ruby style guidelines
- **Cross-platform**: Supports macOS (Intel/Apple Silicon) and Linux architectures