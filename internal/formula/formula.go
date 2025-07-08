package formula

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed formula.rb.tmpl
var formulaTemplate string

// Client provides Homebrew Formula operations
type Client struct{}

// Formula holds information needed to generate Homebrew Formula
type Formula struct {
	Name        string
	Version     string
	Description string
	Homepage    string
	Artifacts   []Artifact
}

// New creates a new Client
func New() *Client {
	return &Client{}
}

// Artifact represents downloadable artifact information
type Artifact struct {
	OS     string
	Arch   string
	URL    string
	SHA256 string
}

// Generate generates Homebrew Formula content
func (c *Client) Generate(formula *Formula) (string, error) {
	// Generate class name (e.g. gorocket -> Gorocket)
	className := c.toClassName(formula.Name)

	// Remove v prefix from version
	version := strings.TrimPrefix(formula.Version, "v")

	// Prepare template data
	data := struct {
		ClassName        string
		Version          string
		ModuleName       string
		MacOSARM64URL    string
		MacOSARM64SHA256 string
		MacOSAMD64URL    string
		MacOSAMD64SHA256 string
		LinuxARM64URL    string
		LinuxARM64SHA256 string
		LinuxAMD64URL    string
		LinuxAMD64SHA256 string
	}{
		ClassName:  className,
		Version:    version,
		ModuleName: formula.Name,
	}

	// Set artifact information to template data
	for _, artifact := range formula.Artifacts {
		switch {
		case artifact.OS == "darwin" && artifact.Arch == "arm64":
			data.MacOSARM64URL = artifact.URL
			data.MacOSARM64SHA256 = artifact.SHA256
		case artifact.OS == "darwin" && artifact.Arch == "amd64":
			data.MacOSAMD64URL = artifact.URL
			data.MacOSAMD64SHA256 = artifact.SHA256
		case artifact.OS == "linux" && artifact.Arch == "arm64":
			data.LinuxARM64URL = artifact.URL
			data.LinuxARM64SHA256 = artifact.SHA256
		case artifact.OS == "linux" && artifact.Arch == "amd64":
			data.LinuxAMD64URL = artifact.URL
			data.LinuxAMD64SHA256 = artifact.SHA256
		}
	}

	// Execute template
	tmpl, err := template.New("formula").Parse(formulaTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse formula template: %w", err)
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute formula template: %w", err)
	}

	return buf.String(), nil
}

// toClassName generates class name from module name
func (c *Client) toClassName(moduleName string) string {
	// Get the last part of the path
	parts := strings.Split(moduleName, "/")
	name := parts[len(parts)-1]

	// Split by hyphens and underscores, then capitalize each part
	delimiters := []string{"-", "_"}
	for _, delimiter := range delimiters {
		if strings.Contains(name, delimiter) {
			parts := strings.Split(name, delimiter)
			for i, part := range parts {
				if len(part) > 0 {
					parts[i] = strings.ToUpper(part[:1]) + part[1:]
				}
			}
			name = strings.Join(parts, "")
			break
		}
	}

	// If no delimiters found, just capitalize the first letter
	if len(name) > 0 && !strings.Contains(name, "-") && !strings.Contains(name, "_") {
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	return name
}
