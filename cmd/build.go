package cmd

import (
	"github.com/koki-develop/gorocket/internal/gorocket"
	"github.com/spf13/cobra"
)

var (
	flagBuildClean bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build binaries for multiple platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder := gorocket.NewBuilder(".gorocket.yml")
		return builder.Build(gorocket.BuildParams{
			Clean: flagBuildClean,
		})
	},
}

func init() {
	buildCmd.Flags().BoolVar(&flagBuildClean, "clean", false, "Remove dist directory before building")
	rootCmd.AddCommand(buildCmd)
}
