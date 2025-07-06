package cmd

import (
	"fmt"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
	"github.com/koki-develop/gorocket/internal/services"
	"github.com/spf13/cobra"
)

type InitCommand struct {
	configService services.ConfigService
}

func NewInitCommand() *cobra.Command {
	fsProvider := providers.NewFileSystemProvider()
	configService := services.NewConfigService(fsProvider)

	initCmd := &InitCommand{
		configService: configService,
	}

	return &cobra.Command{
		Use:   "init",
		Short: "Initialize gorocket configuration",
		Long:  fmt.Sprintf("Initialize gorocket configuration by creating a %s file in the current directory.", models.ConfigFileName),
		RunE:  initCmd.run,
	}
}

func (ic *InitCommand) run(cmd *cobra.Command, args []string) error {
	if ic.configService.ConfigExists() {
		return fmt.Errorf("%s already exists", models.ConfigFileName)
	}

	if err := ic.configService.CreateDefaultConfig(); err != nil {
		return fmt.Errorf("failed to create %s: %w", models.ConfigFileName, err)
	}

	fmt.Printf("Created %s configuration file\n", models.ConfigFileName)
	return nil
}

var initCmd = NewInitCommand()

func init() {
	rootCmd.AddCommand(initCmd)
}
