package gorocket

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// BuildResult represents a build result
type BuildResult struct {
	Binary  string
	OS      string
	Arch    string
	Version string
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

// executeGoBuild executes go build command (for compatibility)
func executeGoBuild(env map[string]string, args []string) error {
	cmd := exec.Command("go", args...)

	// Set environment variables
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w\nstderr: %s", err, stderr.String())
	}

	return nil
}

// getModuleName retrieves module name from go.mod
func getModuleName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer file.Close()

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
