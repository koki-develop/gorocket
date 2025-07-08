package gorocket

import (
	"fmt"
	"os"
)

// Initer provides initialization functionality
type Initer struct {
	configPath string
}

// NewIniter creates a new Initer instance
func NewIniter(configPath string) *Initer {
	return &Initer{
		configPath: configPath,
	}
}

// Init creates a default configuration file
func (i *Initer) Init() error {
	// Error if config file already exists
	if _, err := os.Stat(i.configPath); err == nil {
		return fmt.Errorf("config file already exists: %s", i.configPath)
	}

	// Create default config file
	if err := SaveDefaultConfig(i.configPath); err != nil {
		return fmt.Errorf("failed to save default config: %w", err)
	}

	fmt.Printf("Created %s\n", i.configPath)
	return nil
}
