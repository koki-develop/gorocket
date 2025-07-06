package services

import (
	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
)

type BuilderService interface {
	BuildTargets(buildInfo *models.BuildInfo, targets []models.Target) ([]models.BuildResult, error)
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

func (b *builderService) BuildTargets(buildInfo *models.BuildInfo, targets []models.Target) ([]models.BuildResult, error) {
	var results []models.BuildResult

	for _, target := range targets {
		for _, arch := range target.Arch {
			buildTarget := models.BuildTarget{
				OS:   target.OS,
				Arch: arch,
			}

			binaryPath, err := b.commandProvider.BuildBinary(buildInfo.ModuleName, buildInfo.Version, target.OS, arch)
			
			results = append(results, models.BuildResult{
				Target:     buildTarget,
				BinaryPath: binaryPath,
				Error:      err,
			})
		}
	}

	return results, nil
}