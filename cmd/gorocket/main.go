package gorocket

import (
	"os"

	"github.com/koki-develop/gorocket/internal/gorocket"
	"github.com/spf13/cobra"
)

func Main() {
	app := gorocket.New()

	rootCmd := &cobra.Command{
		Use:   "gorocket",
		Short: "Cross-platform Go binary builder",
	}

	// Command definitions
	rootCmd.AddCommand(
		newInitCommand(app),
		newBuildCommand(app),
		newReleaseCommand(app),
		newVersionCommand(app),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newInitCommand(app *gorocket.GoRocket) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a default .gorocket.yml configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Init()
		},
	}
	return cmd
}

func newBuildCommand(app *gorocket.GoRocket) *cobra.Command {
	var clean bool

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build binaries for multiple platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Build(gorocket.BuildOptions{
				Clean: clean,
			})
		},
	}

	cmd.Flags().BoolVar(&clean, "clean", false, "Remove dist directory before building")

	return cmd
}

func newReleaseCommand(app *gorocket.GoRocket) *cobra.Command {
	var token string
	var draft bool

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Create a GitHub release with built artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.Release(gorocket.ReleaseOptions{
				Token: token,
				Draft: draft,
			})
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "GitHub token (defaults to GITHUB_TOKEN env var)")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create a draft release")

	return cmd
}

func newVersionCommand(app *gorocket.GoRocket) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := app.Version()
			if err != nil {
				return err
			}
			cmd.Println(version)
			return nil
		},
	}
	return cmd
}
