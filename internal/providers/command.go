package providers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type CommandProvider interface {
	RunWithEnv(name string, env []string, args ...string) error
	BuildBinary(moduleName, version, osName, arch, ldflags string) (string, error)
}

type commandProvider struct{}

func NewCommandProvider() CommandProvider {
	return &commandProvider{}
}

func (c *commandProvider) RunWithEnv(name string, env []string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	return cmd.Run()
}

func (c *commandProvider) BuildBinary(moduleName, version, osName, arch, ldflags string) (string, error) {
	var extension string
	if osName == "windows" {
		extension = ".exe"
	}

	binaryName := fmt.Sprintf("%s_%s_%s_%s%s", moduleName, version, osName, arch, extension)
	binaryPath := filepath.Join("dist", binaryName)

	env := []string{
		fmt.Sprintf("GOOS=%s", osName),
		fmt.Sprintf("GOARCH=%s", arch),
	}

	args := []string{"build", "-o", binaryPath}
	if ldflags != "" {
		args = append(args, "-ldflags", ldflags)
	}
	args = append(args, ".")

	err := c.RunWithEnv("go", env, args...)
	if err != nil {
		return "", fmt.Errorf("failed to build for %s/%s: %w", osName, arch, err)
	}

	return binaryPath, nil
}
