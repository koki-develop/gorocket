package config

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

const ConfigFileName = ".gorocket.yaml"

type Config struct {
	Build BuildConfig `yaml:"build"`
	Brew  *BrewConfig `yaml:"brew,omitempty"`
}

type BuildConfig struct {
	Targets []Target `yaml:"targets"`
}

type Target struct {
	OS   string   `yaml:"os"`
	Arch []string `yaml:"arch"`
}

type BrewConfig struct {
	Repository Repository `yaml:"repository"`
}

type Repository struct {
	Owner string `yaml:"owner"`
	Name  string `yaml:"name"`
}

//go:embed default.yaml
var defaultConfigYAML []byte

func ConfigExists() bool {
	_, err := os.Stat(ConfigFileName)
	return !os.IsNotExist(err)
}

func CreateDefaultConfig() error {
	if ConfigExists() {
		return fmt.Errorf("%s already exists", ConfigFileName)
	}

	return os.WriteFile(ConfigFileName, defaultConfigYAML, 0644)
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile(ConfigFileName)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
