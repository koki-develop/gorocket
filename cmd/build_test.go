package cmd

import (
	"fmt"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	providerMocks "github.com/koki-develop/gorocket/internal/providers/mocks"
	serviceMocks "github.com/koki-develop/gorocket/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommand_FormulaGeneration(t *testing.T) {
	tests := []struct {
		name               string
		brewConfig         *models.BrewConfig
		expectFormulaCall  bool
		formulaGenerateErr error
		expectErr          bool
	}{
		{
			name:               "formula generation when brew config exists",
			brewConfig:         &models.BrewConfig{Repository: models.Repository{Owner: "koki-develop", Name: "gorocket"}},
			expectFormulaCall:  true,
			formulaGenerateErr: nil,
			expectErr:          false,
		},
		{
			name:               "no formula generation when brew config is nil",
			brewConfig:         nil,
			expectFormulaCall:  false,
			formulaGenerateErr: nil,
			expectErr:          false,
		},
		{
			name:               "formula generation error",
			brewConfig:         &models.BrewConfig{Repository: models.Repository{Owner: "koki-develop", Name: "gorocket"}},
			expectFormulaCall:  true,
			formulaGenerateErr: fmt.Errorf("formula generation failed"),
			expectErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockFS := providerMocks.NewMockFileSystemProvider(t)
			mockVersionService := serviceMocks.NewMockVersionService(t)
			mockBuilderService := serviceMocks.NewMockBuilderService(t)
			mockArchiverService := serviceMocks.NewMockArchiverService(t)
			mockConfigService := serviceMocks.NewMockConfigService(t)
			mockFormulaService := serviceMocks.NewMockFormulaService(t)

			// Setup common expectations
			config := &models.Config{
				Build: models.BuildConfig{
					Targets: []models.Target{
						{OS: "darwin", Arch: []string{"amd64"}},
					},
				},
				Brew: tt.brewConfig,
			}

			buildInfo := &models.BuildInfo{
				ModuleName: "gorocket",
				Version:    "v1.0.0",
			}

			buildResults := []models.BuildResult{
				{
					Target:     models.BuildTarget{OS: "darwin", Arch: "amd64"},
					BinaryPath: "dist/gorocket_darwin_amd64",
				},
			}

			archiveResults := []models.ArchiveResult{
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
				},
			}

			mockConfigService.EXPECT().ConfigExists().Return(true)
			mockConfigService.EXPECT().LoadConfig().Return(config, nil)
			mockVersionService.EXPECT().GetBuildInfo().Return(buildInfo, nil)
			mockFS.EXPECT().EnsureDistDir(false).Return(nil)
			mockBuilderService.EXPECT().BuildTargets(buildInfo, config.Build).Return(buildResults, nil)
			mockArchiverService.EXPECT().CreateArchives(buildInfo, buildResults).Return(archiveResults, nil)
			mockFS.EXPECT().Remove("dist/gorocket_darwin_amd64").Return(nil)

			// Setup formula service expectation based on test case
			if tt.expectFormulaCall {
				mockFormulaService.EXPECT().GenerateFormula(*buildInfo, archiveResults, *tt.brewConfig).Return(tt.formulaGenerateErr)
			}

			// Create command
			buildCmd := &BuildCommand{
				versionService:  mockVersionService,
				builderService:  mockBuilderService,
				archiverService: mockArchiverService,
				configService:   mockConfigService,
				formulaService:  mockFormulaService,
				fsProvider:      mockFS,
				flagClean:       false,
			}

			// Execute
			err := buildCmd.run()

			// Assert
			if tt.expectErr {
				assert.Error(t, err)
				if tt.formulaGenerateErr != nil {
					assert.Contains(t, err.Error(), "failed to generate formula")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildCommand_FullWorkflow(t *testing.T) {
	// Create mocks
	mockFS := providerMocks.NewMockFileSystemProvider(t)
	mockVersionService := serviceMocks.NewMockVersionService(t)
	mockBuilderService := serviceMocks.NewMockBuilderService(t)
	mockArchiverService := serviceMocks.NewMockArchiverService(t)
	mockConfigService := serviceMocks.NewMockConfigService(t)
	mockFormulaService := serviceMocks.NewMockFormulaService(t)

	// Setup test data
	config := &models.Config{
		Build: models.BuildConfig{
			Targets: []models.Target{
				{OS: "darwin", Arch: []string{"amd64", "arm64"}},
				{OS: "linux", Arch: []string{"amd64"}},
			},
		},
		Brew: &models.BrewConfig{
			Repository: models.Repository{Owner: "koki-develop", Name: "gorocket"},
		},
	}

	buildInfo := &models.BuildInfo{
		ModuleName: "gorocket",
		Version:    "v1.0.0",
	}

	buildResults := []models.BuildResult{
		{Target: models.BuildTarget{OS: "darwin", Arch: "amd64"}, BinaryPath: "dist/gorocket_darwin_amd64"},
		{Target: models.BuildTarget{OS: "darwin", Arch: "arm64"}, BinaryPath: "dist/gorocket_darwin_arm64"},
		{Target: models.BuildTarget{OS: "linux", Arch: "amd64"}, BinaryPath: "dist/gorocket_linux_amd64"},
	}

	archiveResults := []models.ArchiveResult{
		{Target: models.BuildTarget{OS: "darwin", Arch: "amd64"}, ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz"},
		{Target: models.BuildTarget{OS: "darwin", Arch: "arm64"}, ArchivePath: "dist/gorocket_v1.0.0_darwin_arm64.tar.gz"},
		{Target: models.BuildTarget{OS: "linux", Arch: "amd64"}, ArchivePath: "dist/gorocket_v1.0.0_linux_amd64.tar.gz"},
	}

	// Setup expectations in the order they should be called
	mockConfigService.EXPECT().ConfigExists().Return(true)
	mockConfigService.EXPECT().LoadConfig().Return(config, nil)
	mockVersionService.EXPECT().GetBuildInfo().Return(buildInfo, nil)
	mockFS.EXPECT().EnsureDistDir(true).Return(nil)
	mockBuilderService.EXPECT().BuildTargets(buildInfo, config.Build).Return(buildResults, nil)
	mockArchiverService.EXPECT().CreateArchives(buildInfo, buildResults).Return(archiveResults, nil)

	// Expect removal of binary files
	for _, result := range buildResults {
		mockFS.EXPECT().Remove(result.BinaryPath).Return(nil)
	}

	// Expect formula generation
	mockFormulaService.EXPECT().GenerateFormula(*buildInfo, archiveResults, *config.Brew).Return(nil)

	// Create command
	buildCmd := &BuildCommand{
		versionService:  mockVersionService,
		builderService:  mockBuilderService,
		archiverService: mockArchiverService,
		configService:   mockConfigService,
		formulaService:  mockFormulaService,
		fsProvider:      mockFS,
		flagClean:       true,
	}

	// Execute
	err := buildCmd.run()

	// Assert
	assert.NoError(t, err)
}
