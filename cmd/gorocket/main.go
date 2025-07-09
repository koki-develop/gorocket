package gorocket

import (
	"os"

	"github.com/koki-develop/gorocket/internal/gorocket"
	"github.com/spf13/cobra"
)

func Main() {
	rootCmd := &cobra.Command{
		Use:   "gorocket",
		Short: "Cross-platform Go binary builder",
	}

	// Command definitions
	rootCmd.AddCommand(
		newInitCommand(),
		newBuildCommand(),
		newReleaseCommand(),
		newVersionCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a default .gorocket.yml configuration file",
		RunE: func(cmd *cobra.Command, args []string) error {
			initer := gorocket.NewIniter(".gorocket.yml")
			return initer.Init()
		},
	}
	return cmd
}

func newBuildCommand() *cobra.Command {
	var clean bool

	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build binaries for multiple platforms",
		RunE: func(cmd *cobra.Command, args []string) error {
			builder := gorocket.NewBuilder(".gorocket.yml")
			return builder.Build(gorocket.BuildParams{
				Clean: clean,
			})
		},
	}

	cmd.Flags().BoolVar(&clean, "clean", false, "Remove dist directory before building")

	return cmd
}

func newReleaseCommand() *cobra.Command {
	var token string
	var draft bool

	cmd := &cobra.Command{
		Use:   "release",
		Short: "Create a GitHub release with built artifacts",
		RunE: func(cmd *cobra.Command, args []string) error {
			releaser, err := gorocket.NewReleaser(".gorocket.yml", token)
			if err != nil {
				return err
			}
			return releaser.Release(gorocket.ReleaseParams{
				Draft: draft,
			})
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "GitHub token (defaults to GITHUB_TOKEN env var)")
	cmd.Flags().BoolVar(&draft, "draft", false, "Create a draft release")

	return cmd
}

func newVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			version, err := gorocket.GetVersion()
			if err != nil {
				return err
			}
			cmd.Println(version)
			return nil
		},
	}
	return cmd
}
