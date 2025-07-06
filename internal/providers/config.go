package providers

import (
	_ "embed"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/koki-develop/gorocket/internal/models"
)

//go:embed config_default.yaml
var defaultConfigYAML []byte

type ConfigProvider interface {
	ConfigExists() bool
	CreateDefaultConfig() error
	LoadConfig() (*models.Config, error)
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

	file, err := c.fsProvider.Create(models.ConfigFileName)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	_, err = file.Write(defaultConfigYAML)
	return err
}

func (c *configProvider) LoadConfig() (*models.Config, error) {
	file, err := c.fsProvider.Open(models.ConfigFileName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var config models.Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *configProvider) GetDefaultConfigData() []byte {
	return defaultConfigYAML
}
