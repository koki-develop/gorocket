package cmd

import (
	"github.com/koki-develop/gorocket/internal/gorocket"
	"github.com/spf13/cobra"
)

var (
	flagReleaseDraft       bool   // --draft
	flagReleaseGitHubToken string // --github-token

	// TODO: implement
	// flagReleaseClean bool // --clean
)

var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Create a GitHub release with built artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		releaser, err := gorocket.NewReleaser(".gorocket.yml", flagReleaseGitHubToken)
		if err != nil {
			return err
		}
		return releaser.Release(gorocket.ReleaseParams{
			Draft: flagReleaseDraft,
		})
	},
}

func init() {
	rootCmd.AddCommand(releaseCmd)
	releaseCmd.Flags().StringVar(&flagReleaseGitHubToken, "github-token", "", "GitHub token (defaults to GITHUB_TOKEN env var)")
	releaseCmd.Flags().BoolVar(&flagReleaseDraft, "draft", false, "Create a draft release")
}
