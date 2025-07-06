package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
	"github.com/koki-develop/gorocket/internal/services"
	"github.com/spf13/cobra"
)

type BuildCommand struct {
	versionService  services.VersionService
	builderService  services.BuilderService
	archiverService services.ArchiverService
	configService   services.ConfigService
	formulaService  services.FormulaService
	fsProvider      providers.FileSystemProvider
	flagClean       bool
}

func NewBuildCommand() *cobra.Command {
	gitProvider := providers.NewGitProvider()
	fsProvider := providers.NewFileSystemProvider()
	commandProvider := providers.NewCommandProvider()

	versionService := services.NewVersionService(gitProvider, fsProvider)
	builderService := services.NewBuilderService(commandProvider, fsProvider)
	archiverService := services.NewArchiverService(fsProvider)
	configService := services.NewConfigService(fsProvider)
	formulaService := services.NewFormulaService(fsProvider)

	buildCmd := &BuildCommand{
		versionService:  versionService,
		builderService:  builderService,
		archiverService: archiverService,
		configService:   configService,
		formulaService:  formulaService,
		fsProvider:      fsProvider,
	}

	cobraCmd := &cobra.Command{
		Use:   "build",
		Short: "Build binaries for multiple platforms",
		Long:  "Build binaries for multiple platforms based on the configuration in .gorocket.yaml",
		RunE:  buildCmd.run,
	}

	cobraCmd.Flags().BoolVar(&buildCmd.flagClean, "clean", false, "Clean dist directory before building")

	return cobraCmd
}

func (bc *BuildCommand) run(cmd *cobra.Command, args []string) error {
	if !bc.configService.ConfigExists() {
		return fmt.Errorf("%s not found. Run 'gorocket init' first", models.ConfigFileName)
	}

	cfg, err := bc.configService.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	buildInfo, err := bc.versionService.GetBuildInfo()
	if err != nil {
		return fmt.Errorf("failed to get build info: %w", err)
	}

	if err := bc.fsProvider.EnsureDistDir(bc.flagClean); err != nil {
		return fmt.Errorf("failed to prepare dist directory: %w", err)
	}

	fmt.Printf("Building %s version %s\n", buildInfo.ModuleName, buildInfo.Version)

	buildResults, err := bc.builderService.BuildTargets(buildInfo, cfg.Build.Targets)
	if err != nil {
		return fmt.Errorf("failed to build targets: %w", err)
	}

	for _, result := range buildResults {
		if result.Error != nil {
			return fmt.Errorf("failed to build target %s/%s: %w", result.Target.OS, result.Target.Arch, result.Error)
		}
		fmt.Printf("Building for %s/%s... Done\n", result.Target.OS, result.Target.Arch)
	}

	archiveResults, err := bc.archiverService.CreateArchives(buildInfo, buildResults)
	if err != nil {
		return fmt.Errorf("failed to create archives: %w", err)
	}

	for _, result := range archiveResults {
		if result.Error != nil {
			return fmt.Errorf("failed to create archive for %s/%s: %w", result.Target.OS, result.Target.Arch, result.Error)
		}
		fmt.Printf("Created %s\n", filepath.Base(result.ArchivePath))
	}

	for _, buildResult := range buildResults {
		if buildResult.Error == nil {
			if err := bc.fsProvider.Remove(buildResult.BinaryPath); err != nil {
				return fmt.Errorf("failed to remove binary file: %w", err)
			}
		}
	}

	if cfg.Brew != nil {
		fmt.Println("Generating Homebrew Formula...")
		if err := bc.formulaService.GenerateFormula(*buildInfo, archiveResults, *cfg.Brew); err != nil {
			return fmt.Errorf("failed to generate formula: %w", err)
		}
		fmt.Printf("Created %s.rb\n", buildInfo.ModuleName)
	}

	fmt.Println("Build completed successfully!")
	return nil
}

var buildCmd = NewBuildCommand()

func init() {
	rootCmd.AddCommand(buildCmd)
}
