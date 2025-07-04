package cmd

import (
	"fmt"

	"github.com/koki-develop/gorocket/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gorocket configuration",
	Long:  fmt.Sprintf("Initialize gorocket configuration by creating a %s file in the current directory.", config.ConfigFileName),
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.ConfigExists() {
			return fmt.Errorf("%s already exists", config.ConfigFileName)
		}

		if err := config.CreateDefaultConfig(); err != nil {
			return fmt.Errorf("failed to create %s: %w", config.ConfigFileName, err)
		}

		fmt.Printf("Created %s configuration file\n", config.ConfigFileName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
