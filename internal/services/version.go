package services

import (
	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers"
)

type VersionService interface {
	GetBuildInfo() (*models.BuildInfo, error)
}

type versionService struct {
	gitProvider        providers.GitProvider
	fileSystemProvider providers.FileSystemProvider
}

func NewVersionService(gitProvider providers.GitProvider, fileSystemProvider providers.FileSystemProvider) VersionService {
	return &versionService{
		gitProvider:        gitProvider,
		fileSystemProvider: fileSystemProvider,
	}
}

func (v *versionService) GetBuildInfo() (*models.BuildInfo, error) {
	moduleName, err := v.fileSystemProvider.GetModuleName()
	if err != nil {
		return nil, err
	}

	version, err := v.gitProvider.GetCurrentVersion()
	if err != nil {
		return nil, err
	}

	return &models.BuildInfo{
		ModuleName: moduleName,
		Version:    version,
	}, nil
}
