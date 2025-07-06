package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
	"github.com/koki-develop/gorocket/internal/services"
	"github.com/spf13/cobra"
)

type ReleaseCommand struct {
	buildCommand   *BuildCommand
	gitProvider    providers.GitProvider
	githubProvider providers.GitHubProvider
	flagClean      bool
}

func NewReleaseCommand() *cobra.Command {
	var flagClean bool

	cobraCmd := &cobra.Command{
		Use:   "release",
		Short: "Create GitHub Release and upload assets",
		Long:  "Create GitHub Release, upload built assets, and optionally update Homebrew tap repository",
	}

	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		cobraCmd.RunE = func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("GITHUB_TOKEN environment variable is required")
		}
	} else {
		gitProvider := providers.NewGitProvider()
		fsProvider := providers.NewFileSystemProvider()
		commandProvider := providers.NewCommandProvider()

		githubProvider := providers.NewGitHubProvider(githubToken, fsProvider)

		versionService := services.NewVersionService(gitProvider, fsProvider)
		builderService := services.NewBuilderService(commandProvider, fsProvider)
		archiverService := services.NewArchiverService(fsProvider)
		configService := services.NewConfigService(fsProvider)
		formulaService := services.NewFormulaService(fsProvider, gitProvider)

		buildCmd := &BuildCommand{
			versionService:  versionService,
			builderService:  builderService,
			archiverService: archiverService,
			configService:   configService,
			formulaService:  formulaService,
			fsProvider:      fsProvider,
		}

		releaseCmd := &ReleaseCommand{
			buildCommand:   buildCmd,
			gitProvider:    gitProvider,
			githubProvider: githubProvider,
		}

		cobraCmd.RunE = func(cmd *cobra.Command, args []string) error {
			releaseCmd.flagClean = flagClean
			return releaseCmd.run()
		}
	}

	cobraCmd.Flags().BoolVar(&flagClean, "clean", false, "Clean dist directory before building")

	return cobraCmd
}

func (rc *ReleaseCommand) run() error {
	ctx := context.Background()

	// Set clean flag and run build
	rc.buildCommand.flagClean = rc.flagClean
	fmt.Println("Running build...")
	buildResult, err := rc.buildCommand.RunBuildWithResults()
	if err != nil {
		return fmt.Errorf("failed to run build: %w", err)
	}

	githubRepo, err := rc.gitProvider.GetGitHubRepository()
	if err != nil {
		return fmt.Errorf("failed to get GitHub repository info: %w", err)
	}

	// Check if release already exists
	fmt.Println("Checking for existing release...")
	existingRelease, err := rc.githubProvider.GetRelease(ctx, githubRepo, buildResult.BuildInfo.Version)
	if err != nil {
		return fmt.Errorf("failed to get release: %w", err)
	}

	var releaseURL string
	var assets []models.ReleaseAsset

	// Collect archive files from build results
	for _, result := range buildResult.ArchiveResults {
		assets = append(assets, models.ReleaseAsset{
			Name: filepath.Base(result.ArchivePath),
			Path: result.ArchivePath,
		})
	}

	if existingRelease == nil {
		fmt.Println("Creating GitHub Release...")
		githubRelease, err := rc.githubProvider.CreateRelease(ctx, githubRepo, buildResult.BuildInfo.Version)
		if err != nil {
			return fmt.Errorf("failed to create GitHub release: %w", err)
		}
		releaseURL = *githubRelease.HTMLURL

		fmt.Println("Uploading assets...")
		if err := rc.githubProvider.UploadAssets(ctx, githubRepo, githubRelease, assets); err != nil {
			return fmt.Errorf("failed to upload assets: %w", err)
		}
	} else {
		releaseURL = *existingRelease.HTMLURL
		fmt.Printf("Release %s already exists\n", buildResult.BuildInfo.Version)
	}

	// Update Homebrew tap repository if configured
	if buildResult.Config.Brew != nil {
		fmt.Println("Updating Homebrew tap repository...")
		formulaPath := filepath.Join("dist", buildResult.BuildInfo.ModuleName+".rb")
		formulaContent, err := os.ReadFile(formulaPath)
		if err != nil {
			return fmt.Errorf("failed to read formula file: %w", err)
		}

		if err := rc.githubProvider.UpdateTapRepository(ctx, &buildResult.Config.Brew.Repository, string(formulaContent), buildResult.BuildInfo.ModuleName, buildResult.BuildInfo.Version); err != nil {
			return fmt.Errorf("failed to update tap repository: %w", err)
		}
	}

	fmt.Printf("Release %s created successfully!\n", buildResult.BuildInfo.Version)
	fmt.Printf("Release URL: %s\n", releaseURL)

	if len(assets) > 0 {
		fmt.Println("\nUploaded assets:")
		for _, asset := range assets {
			fmt.Printf("  - %s\n", asset.Name)
		}
	}

	fmt.Println("\nRelease completed successfully!")
	return nil
}

var releaseCmd = NewReleaseCommand()

func init() {
	rootCmd.AddCommand(releaseCmd)
}
