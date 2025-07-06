package models

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

const ConfigFileName = ".gorocket.yaml"