package config

import (
	"bytes"
	"fmt"
	"os"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Config defines the configuration file structure
type Config struct {
	Build struct {
		Targets []Target `yaml:"targets"`
		Ldflags string   `yaml:"ldflags"`
	} `yaml:"build"`

	Brew struct {
		Repository string `yaml:"repository"`
	} `yaml:"brew"`
}

// Target represents a build target
type Target struct {
	OS   string `yaml:"os"`
	Arch string `yaml:"arch"`
}

// LoadConfig loads the configuration file
func LoadConfig(path string, data map[string]any) (*Config, error) {
	// Read file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Process template if data is provided
	if data != nil {
		tmpl, err := template.New("config").Parse(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("failed to execute template: %w", err)
		}
		content = buf.Bytes()
	}

	// Parse YAML
	var config Config
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}
