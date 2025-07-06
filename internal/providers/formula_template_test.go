package providers

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateFormula(t *testing.T) {
	tests := []struct {
		name     string
		data     FormulaTemplateData
		expected []string
		wantErr  bool
	}{
		{
			name: "generate complete formula",
			data: FormulaTemplateData{
				ClassName:        "Gorocket",
				Version:          "1.0.0",
				ModuleName:       "gorocket",
				MacOSARM64URL:    "https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_darwin_arm64.tar.gz",
				MacOSARM64SHA256: "darwin_arm64_sha256",
				MacOSAMD64URL:    "https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_darwin_amd64.tar.gz",
				MacOSAMD64SHA256: "darwin_amd64_sha256",
				LinuxARM64URL:    "https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_linux_arm64.tar.gz",
				LinuxARM64SHA256: "linux_arm64_sha256",
				LinuxAMD64URL:    "https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_linux_amd64.tar.gz",
				LinuxAMD64SHA256: "linux_amd64_sha256",
			},
			expected: []string{
				"# typed: strict",
				"# frozen_string_literal: true",
				"# Gorocket formula",
				"class Gorocket < Formula",
				`version "1.0.0"`,
				"on_macos do",
				"if Hardware::CPU.arm?",
				"https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_darwin_arm64.tar.gz",
				"darwin_arm64_sha256",
				"https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_darwin_amd64.tar.gz",
				"darwin_amd64_sha256",
				"on_linux do",
				"https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_linux_arm64.tar.gz",
				"linux_arm64_sha256",
				"https://github.com/owner/repo/releases/download/v1.0.0/gorocket_v1.0.0_linux_amd64.tar.gz",
				"linux_amd64_sha256",
				"def install",
				`bin.install "gorocket"`,
				"end",
			},
			wantErr: false,
		},
		{
			name: "generate formula with empty values",
			data: FormulaTemplateData{
				ClassName:        "TestApp",
				Version:          "0.1.0",
				ModuleName:       "test-app",
				MacOSARM64URL:    "",
				MacOSARM64SHA256: "",
				MacOSAMD64URL:    "",
				MacOSAMD64SHA256: "",
				LinuxARM64URL:    "",
				LinuxARM64SHA256: "",
				LinuxAMD64URL:    "",
				LinuxAMD64SHA256: "",
			},
			expected: []string{
				"# typed: strict",
				"# frozen_string_literal: true",
				"# TestApp formula",
				"class TestApp < Formula",
				`version "0.1.0"`,
				"on_macos do",
				"on_linux do",
				"def install",
				`bin.install "test-app"`,
				"end",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GenerateFormula(tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, result)

			// Check that all expected strings are present
			for _, expected := range tt.expected {
				assert.Contains(t, result, expected, "Formula should contain: %s", expected)
			}

			// Check Ruby syntax elements
			assert.Contains(t, result, "# typed: strict")
			assert.Contains(t, result, "# frozen_string_literal: true")
			assert.Contains(t, result, "class "+tt.data.ClassName+" < Formula")
			assert.Contains(t, result, `version "`+tt.data.Version+`"`)
			assert.Contains(t, result, "def install")
			assert.Contains(t, result, `bin.install "`+tt.data.ModuleName+`"`)
		})
	}
}

func TestToClassName(t *testing.T) {
	tests := []struct {
		name       string
		moduleName string
		expected   string
	}{
		{
			name:       "simple name",
			moduleName: "gorocket",
			expected:   "Gorocket",
		},
		{
			name:       "name with hyphen",
			moduleName: "go-rocket",
			expected:   "Gorocket",
		},
		{
			name:       "name with underscore",
			moduleName: "go_rocket",
			expected:   "Gorocket",
		},
		{
			name:       "name with mixed separators",
			moduleName: "go-rocket_tool",
			expected:   "Gorockettool",
		},
		{
			name:       "github path",
			moduleName: "github.com/koki-develop/gorocket",
			expected:   "Gorocket",
		},
		{
			name:       "github path with hyphen",
			moduleName: "github.com/koki-develop/go-rocket",
			expected:   "Gorocket",
		},
		{
			name:       "single character",
			moduleName: "a",
			expected:   "A",
		},
		{
			name:       "empty string",
			moduleName: "",
			expected:   "",
		},
		{
			name:       "multiple path segments",
			moduleName: "example.com/user/my-awesome-tool",
			expected:   "Myawesometool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToClassName(tt.moduleName)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormulaTemplateStructure(t *testing.T) {
	data := FormulaTemplateData{
		ClassName:        "TestFormula",
		Version:          "1.0.0",
		ModuleName:       "test",
		MacOSARM64URL:    "https://example.com/macos-arm64.tar.gz",
		MacOSARM64SHA256: "macos_arm64_sha",
		MacOSAMD64URL:    "https://example.com/macos-amd64.tar.gz",
		MacOSAMD64SHA256: "macos_amd64_sha",
		LinuxARM64URL:    "https://example.com/linux-arm64.tar.gz",
		LinuxARM64SHA256: "linux_arm64_sha",
		LinuxAMD64URL:    "https://example.com/linux-amd64.tar.gz",
		LinuxAMD64SHA256: "linux_amd64_sha",
	}

	result, err := GenerateFormula(data)
	require.NoError(t, err)

	lines := strings.Split(result, "\n")

	// Check that the structure follows the expected order
	var foundSections []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# typed:") {
			foundSections = append(foundSections, "sorbet_sigil")
		} else if strings.HasPrefix(line, "# frozen_string_literal:") {
			foundSections = append(foundSections, "frozen_string")
		} else if strings.HasPrefix(line, "class") {
			foundSections = append(foundSections, "class_declaration")
		} else if strings.HasPrefix(line, "on_macos") {
			foundSections = append(foundSections, "macos_section")
		} else if strings.HasPrefix(line, "on_linux") {
			foundSections = append(foundSections, "linux_section")
		} else if strings.HasPrefix(line, "def install") {
			foundSections = append(foundSections, "install_method")
		}
	}

	// Verify the sections appear in the correct order
	expectedOrder := []string{
		"sorbet_sigil",
		"frozen_string",
		"class_declaration",
		"macos_section",
		"linux_section",
		"install_method",
	}

	assert.Equal(t, expectedOrder, foundSections, "Formula sections should appear in the correct order")
}

func TestFormulaTemplateIndentation(t *testing.T) {
	data := FormulaTemplateData{
		ClassName:        "TestFormula",
		Version:          "1.0.0",
		ModuleName:       "test",
		MacOSARM64URL:    "https://example.com/macos-arm64.tar.gz",
		MacOSARM64SHA256: "macos_arm64_sha",
		MacOSAMD64URL:    "https://example.com/macos-amd64.tar.gz",
		MacOSAMD64SHA256: "macos_amd64_sha",
		LinuxARM64URL:    "https://example.com/linux-arm64.tar.gz",
		LinuxARM64SHA256: "linux_arm64_sha",
		LinuxAMD64URL:    "https://example.com/linux-amd64.tar.gz",
		LinuxAMD64SHA256: "linux_amd64_sha",
	}

	result, err := GenerateFormula(data)
	require.NoError(t, err)

	lines := strings.Split(result, "\n")

	// Check specific indentation patterns
	for i, line := range lines {
		if strings.Contains(line, "version ") {
			assert.True(t, strings.HasPrefix(line, "  "), "Line %d should have 2-space indentation: %s", i+1, line)
		}
		if strings.Contains(line, "on_macos") || strings.Contains(line, "on_linux") {
			assert.True(t, strings.HasPrefix(line, "  "), "Line %d should have 2-space indentation: %s", i+1, line)
		}
		if strings.Contains(line, "if Hardware::CPU.arm?") {
			assert.True(t, strings.HasPrefix(line, "    "), "Line %d should have 4-space indentation: %s", i+1, line)
		}
		if strings.Contains(line, "def install") {
			assert.True(t, strings.HasPrefix(line, "  "), "Line %d should have 2-space indentation: %s", i+1, line)
		}
		if strings.Contains(line, "bin.install") {
			assert.True(t, strings.HasPrefix(line, "    "), "Line %d should have 4-space indentation: %s", i+1, line)
		}
	}
}
