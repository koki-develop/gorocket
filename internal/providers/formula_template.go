package providers

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

type FormulaTemplateData struct {
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
}

func GenerateFormula(data FormulaTemplateData) (string, error) {
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

func ToClassName(moduleName string) string {
	parts := strings.Split(moduleName, "/")
	name := parts[len(parts)-1]

	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "_", "")

	if len(name) > 0 {
		name = strings.ToUpper(name[:1]) + name[1:]
	}

	return name
}
