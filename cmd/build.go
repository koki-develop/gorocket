package cmd

import (
	"github.com/koki-develop/gorocket/internal/gorocket"
	"github.com/spf13/cobra"
)

var (
	flagBuildClean      bool
	flagBuildAllowDirty bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build binaries for multiple platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		builder := gorocket.NewBuilder(".gorocket.yml")
		return builder.Build(gorocket.BuildParams{
			Clean:      flagBuildClean,
			AllowDirty: flagBuildAllowDirty,
		})
	},
}

func init() {
	buildCmd.Flags().BoolVar(&flagBuildClean, "clean", false, "Remove dist directory before building")
	buildCmd.Flags().BoolVar(&flagBuildAllowDirty, "allow-dirty", false, "Allow building without git tag (uses v0.0.0-dev as version)")
	rootCmd.AddCommand(buildCmd)
}
