package formula

import (
	"fmt"
	"strings"
	"text/template"
)

const formulaTemplate = `# typed: strict
# frozen_string_literal: true

# {{.ClassName}} formula
class {{.ClassName}} < Formula
  version "{{.Version}}"

  on_macos do
    if Hardware::CPU.arm?
      url "{{.MacOSARM64URL}}"
      sha256 "{{.MacOSARM64SHA256}}"
    else
      url "{{.MacOSAMD64URL}}"
      sha256 "{{.MacOSAMD64SHA256}}"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "{{.LinuxARM64URL}}"
      sha256 "{{.LinuxARM64SHA256}}"
    else
      url "{{.LinuxAMD64URL}}"
      sha256 "{{.LinuxAMD64SHA256}}"
    end
  end

  def install
    bin.install "{{.ModuleName}}"
  end
end
`

// Formula holds information needed to generate Homebrew Formula
type Formula struct {
	Name        string
	Version     string
	Description string
	Homepage    string
	Artifacts   []Artifact
}

// Artifact represents downloadable artifact information
type Artifact struct {
	OS     string
	Arch   string
	URL    string
	SHA256 string
}

// Generate generates Homebrew Formula content
func Generate(formula *Formula) (string, error) {
	// Generate class name (e.g. gorocket -> Gorocket)
	className := toClassName(formula.Name)

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

// UpdateTapRepository updates Homebrew tap repository
func UpdateTapRepository(client interface {
	UpdateFile(path, content, message string) error
}, formula string, moduleName, version string) error {
	formulaPath := fmt.Sprintf("Formula/%s.rb", moduleName)
	commitMessage := fmt.Sprintf("Update %s to %s", moduleName, version)

	return client.UpdateFile(formulaPath, formula, commitMessage)
}

// toClassName generates class name from module name
func toClassName(moduleName string) string {
	// Get the last part of the path
	parts := strings.Split(moduleName, "/")
	name := parts[len(parts)-1]

	// Remove hyphens and underscores
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")

	// Capitalize the first letter
	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	return name
}
