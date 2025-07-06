package providers

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/goccy/go-yaml"
	"github.com/koki-develop/gorocket/internal/models"
)

//go:embed config_default.yaml
var defaultConfigYAML []byte

type ConfigProvider interface {
	ConfigExists() bool
	CreateDefaultConfig() error
	LoadConfig(templateData *models.TemplateData) (*models.Config, error)
	GetDefaultConfigData() []byte
}

type configProvider struct {
	fsProvider FileSystemProvider
}

func NewConfigProvider(fsProvider FileSystemProvider) ConfigProvider {
	return &configProvider{
		fsProvider: fsProvider,
	}
}

func (c *configProvider) ConfigExists() bool {
	_, err := c.fsProvider.Stat(models.ConfigFileName)
	return err == nil
}

func (c *configProvider) CreateDefaultConfig() error {
	if c.ConfigExists() {
		return fmt.Errorf("%s already exists", models.ConfigFileName)
	}

	return c.fsProvider.WriteFile(models.ConfigFileName, defaultConfigYAML, 0644)
}

func (c *configProvider) LoadConfig(templateData *models.TemplateData) (*models.Config, error) {
	file, err := c.fsProvider.Open(models.ConfigFileName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// Read file content
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(file); err != nil {
		return nil, err
	}

	// Process template
	tmpl, err := template.New("config").Parse(buf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to parse config template: %w", err)
	}

	var processedBuf bytes.Buffer
	if err := tmpl.Execute(&processedBuf, templateData); err != nil {
		return nil, fmt.Errorf("failed to execute config template: %w", err)
	}

	// Parse YAML
	var config models.Config
	decoder := yaml.NewDecoder(&processedBuf)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config YAML: %w", err)
	}

	return &config, nil
}

func (c *configProvider) GetDefaultConfigData() []byte {
	return defaultConfigYAML
}

