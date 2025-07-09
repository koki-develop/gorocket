package cmd

import (
	"github.com/koki-develop/gorocket/internal/gorocket"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a default .gorocket.yml configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		initer := gorocket.NewIniter(".gorocket.yml")
		return initer.Init()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
