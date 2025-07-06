package services

import (
	"github.com/hashicorp/go-multierror"
	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
)

type BuilderService interface {
	BuildTargets(buildInfo *models.BuildInfo, buildConfig models.BuildConfig) ([]models.BuildResult, error)
}

type builderService struct {
	commandProvider    providers.CommandProvider
	fileSystemProvider providers.FileSystemProvider
}

func NewBuilderService(commandProvider providers.CommandProvider, fileSystemProvider providers.FileSystemProvider) BuilderService {
	return &builderService{
		commandProvider:    commandProvider,
		fileSystemProvider: fileSystemProvider,
	}
}

func (b *builderService) BuildTargets(buildInfo *models.BuildInfo, buildConfig models.BuildConfig) ([]models.BuildResult, error) {
	var results []models.BuildResult
	var errGroup *multierror.Error

	for _, target := range buildConfig.Targets {
		for _, arch := range target.Arch {
			buildTarget := models.BuildTarget{
				OS:   target.OS,
				Arch: arch,
			}

			binaryPath, err := b.commandProvider.BuildBinary(buildInfo.ModuleName, buildInfo.Version, target.OS, arch, buildConfig.LdFlags)
			if err != nil {
				errGroup = multierror.Append(errGroup, err)
				continue
			}

			results = append(results, models.BuildResult{
				Target:     buildTarget,
				BinaryPath: binaryPath,
			})
		}
	}

	if err := errGroup.ErrorOrNil(); err != nil {
		return nil, err
	}

	return results, nil
}
