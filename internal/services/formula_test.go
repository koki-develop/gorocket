package services

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFormulaService_GenerateFormula(t *testing.T) {
	tests := []struct {
		name           string
		buildInfo      models.BuildInfo
		archiveResults []models.ArchiveResult
		brewConfig     models.BrewConfig
		setupMocks     func(*mocks.MockFileSystemProvider, *mocks.MockGitProvider)
		wantErr        bool
	}{
		{
			name: "successful formula generation",
			buildInfo: models.BuildInfo{
				ModuleName: "gorocket",
				Version:    "v1.0.0",
			},
			archiveResults: []models.ArchiveResult{
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
				},
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "arm64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_arm64.tar.gz",
				},
				{
					Target:      models.BuildTarget{OS: "linux", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_linux_amd64.tar.gz",
				},
				{
					Target:      models.BuildTarget{OS: "linux", Arch: "arm64"},
					ArchivePath: "dist/gorocket_v1.0.0_linux_arm64.tar.gz",
				},
			},
			brewConfig: models.BrewConfig{
				Repository: models.Repository{
					Owner: "koki-develop",
					Name:  "gorocket",
				},
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider, mockGit *mocks.MockGitProvider) {
				mockGit.EXPECT().GetGitHubRepository().Return(&models.GitHubRepository{
					Owner: "koki-develop",
					Name:  "gorocket",
				}, nil).Once()
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_amd64.tar.gz").Return(io.NopCloser(strings.NewReader("content1")), nil).Once()
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("darwin_amd64_sha256", nil).Once()
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_arm64.tar.gz").Return(io.NopCloser(strings.NewReader("content2")), nil).Once()
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("darwin_arm64_sha256", nil).Once()
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_linux_amd64.tar.gz").Return(io.NopCloser(strings.NewReader("content3")), nil).Once()
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("linux_amd64_sha256", nil).Once()
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_linux_arm64.tar.gz").Return(io.NopCloser(strings.NewReader("content4")), nil).Once()
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("linux_arm64_sha256", nil).Once()
				mockFS.EXPECT().WriteFile(
					filepath.Join("dist", "gorocket.rb"),
					mock.MatchedBy(func(content []byte) bool {
						contentStr := string(content)
						return strings.Contains(contentStr, "# typed: strict") &&
							strings.Contains(contentStr, "# frozen_string_literal: true") &&
							strings.Contains(contentStr, "class Gorocket < Formula") &&
							strings.Contains(contentStr, `version "1.0.0"`) &&
							strings.Contains(contentStr, "darwin_amd64_sha256") &&
							strings.Contains(contentStr, "darwin_arm64_sha256") &&
							strings.Contains(contentStr, "linux_amd64_sha256") &&
							strings.Contains(contentStr, "linux_arm64_sha256") &&
							strings.Contains(contentStr, "https://github.com/koki-develop/gorocket/releases/download/v1.0.0/gorocket_v1.0.0_darwin_amd64.tar.gz") &&
							strings.Contains(contentStr, "https://github.com/koki-develop/gorocket/releases/download/v1.0.0/gorocket_v1.0.0_darwin_arm64.tar.gz") &&
							strings.Contains(contentStr, "https://github.com/koki-develop/gorocket/releases/download/v1.0.0/gorocket_v1.0.0_linux_amd64.tar.gz") &&
							strings.Contains(contentStr, "https://github.com/koki-develop/gorocket/releases/download/v1.0.0/gorocket_v1.0.0_linux_arm64.tar.gz")
					}),
					mock.AnythingOfType("fs.FileMode"),
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "successful formula generation with multiple platforms",
			buildInfo: models.BuildInfo{
				ModuleName: "gorocket",
				Version:    "v1.0.0",
			},
			archiveResults: []models.ArchiveResult{
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
				},
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "arm64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_arm64.tar.gz",
				},
			},
			brewConfig: models.BrewConfig{
				Repository: models.Repository{
					Owner: "koki-develop",
					Name:  "gorocket",
				},
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider, mockGit *mocks.MockGitProvider) {
				mockGit.EXPECT().GetGitHubRepository().Return(&models.GitHubRepository{
					Owner: "koki-develop",
					Name:  "gorocket",
				}, nil)
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_amd64.tar.gz").Return(io.NopCloser(strings.NewReader("content1")), nil)
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("darwin_amd64_sha256", nil)
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_arm64.tar.gz").Return(io.NopCloser(strings.NewReader("content2")), nil)
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("darwin_arm64_sha256", nil)
				mockFS.EXPECT().WriteFile(
					filepath.Join("dist", "gorocket.rb"),
					mock.MatchedBy(func(content []byte) bool {
						contentStr := string(content)
						return strings.Contains(contentStr, "class Gorocket < Formula") &&
							strings.Contains(contentStr, "darwin_amd64_sha256")
					}),
					mock.AnythingOfType("fs.FileMode"),
				).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "sha256 calculation error",
			buildInfo: models.BuildInfo{
				ModuleName: "gorocket",
				Version:    "v1.0.0",
			},
			archiveResults: []models.ArchiveResult{
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
				},
			},
			brewConfig: models.BrewConfig{
				Repository: models.Repository{
					Owner: "koki-develop",
					Name:  "gorocket",
				},
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider, mockGit *mocks.MockGitProvider) {
				mockGit.EXPECT().GetGitHubRepository().Return(&models.GitHubRepository{
					Owner: "koki-develop",
					Name:  "gorocket",
				}, nil)
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_amd64.tar.gz").Return(nil, fmt.Errorf("file not found"))
			},
			wantErr: true,
		},
		{
			name: "write file error",
			buildInfo: models.BuildInfo{
				ModuleName: "gorocket",
				Version:    "v1.0.0",
			},
			archiveResults: []models.ArchiveResult{
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
				},
			},
			brewConfig: models.BrewConfig{
				Repository: models.Repository{
					Owner: "koki-develop",
					Name:  "gorocket",
				},
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider, mockGit *mocks.MockGitProvider) {
				mockGit.EXPECT().GetGitHubRepository().Return(&models.GitHubRepository{
					Owner: "koki-develop",
					Name:  "gorocket",
				}, nil)
				mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_amd64.tar.gz").Return(io.NopCloser(strings.NewReader("content")), nil)
				mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("darwin_amd64_sha256", nil)
				mockFS.EXPECT().WriteFile(
					mock.AnythingOfType("string"),
					mock.AnythingOfType("[]uint8"),
					mock.AnythingOfType("fs.FileMode"),
				).Return(fmt.Errorf("write failed"))
			},
			wantErr: true,
		},
		{
			name: "git repository error",
			buildInfo: models.BuildInfo{
				ModuleName: "gorocket",
				Version:    "v1.0.0",
			},
			archiveResults: []models.ArchiveResult{
				{
					Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
					ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
				},
			},
			brewConfig: models.BrewConfig{
				Repository: models.Repository{
					Owner: "koki-develop",
					Name:  "gorocket",
				},
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider, mockGit *mocks.MockGitProvider) {
				mockGit.EXPECT().GetGitHubRepository().Return(nil, fmt.Errorf("git repository not found"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			mockGit := mocks.NewMockGitProvider(t)
			tt.setupMocks(mockFS, mockGit)

			service := NewFormulaService(mockFS, mockGit)
			err := service.GenerateFormula(tt.buildInfo, tt.archiveResults, tt.brewConfig)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFormulaService_buildFormulaInfo(t *testing.T) {
	mockFS := mocks.NewMockFileSystemProvider(t)
	mockGit := mocks.NewMockGitProvider(t)
	service := &formulaService{fsProvider: mockFS, gitProvider: mockGit}

	buildInfo := models.BuildInfo{
		ModuleName: "gorocket",
		Version:    "v1.0.0",
	}

	archiveResults := []models.ArchiveResult{
		{
			Target:      models.BuildTarget{OS: "darwin", Arch: "amd64"},
			ArchivePath: "dist/gorocket_v1.0.0_darwin_amd64.tar.gz",
		},
		{
			Target:      models.BuildTarget{OS: "linux", Arch: "arm64"},
			ArchivePath: "dist/gorocket_v1.0.0_linux_arm64.tar.gz",
		},
	}

	brewConfig := models.BrewConfig{
		Repository: models.Repository{
			Owner: "koki-develop",
			Name:  "gorocket",
		},
	}

	mockGit.EXPECT().GetGitHubRepository().Return(&models.GitHubRepository{
		Owner: "koki-develop",
		Name:  "gorocket",
	}, nil).Once()
	mockFS.EXPECT().Open("dist/gorocket_v1.0.0_darwin_amd64.tar.gz").Return(io.NopCloser(strings.NewReader("content1")), nil).Once()
	mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("darwin_amd64_sha256", nil).Once()
	mockFS.EXPECT().Open("dist/gorocket_v1.0.0_linux_arm64.tar.gz").Return(io.NopCloser(strings.NewReader("content2")), nil).Once()
	mockFS.EXPECT().CalculateSHA256(mock.Anything).Return("linux_arm64_sha256", nil).Once()

	formulaInfo, err := service.buildFormulaInfo(buildInfo, archiveResults, brewConfig)

	assert.NoError(t, err)
	assert.Equal(t, "gorocket", formulaInfo.ModuleName)
	assert.Equal(t, "1.0.0", formulaInfo.Version)
	assert.Equal(t, "koki-develop", formulaInfo.Repository.Owner)
	assert.Equal(t, "gorocket", formulaInfo.Repository.Name)

	assert.Equal(t, "https://github.com/koki-develop/gorocket/releases/download/v1.0.0/gorocket_v1.0.0_darwin_amd64.tar.gz", formulaInfo.PlatformURLs["darwin"]["amd64"].URL)
	assert.Equal(t, "darwin_amd64_sha256", formulaInfo.PlatformURLs["darwin"]["amd64"].SHA256)

	assert.Equal(t, "https://github.com/koki-develop/gorocket/releases/download/v1.0.0/gorocket_v1.0.0_linux_arm64.tar.gz", formulaInfo.PlatformURLs["linux"]["arm64"].URL)
	assert.Equal(t, "linux_arm64_sha256", formulaInfo.PlatformURLs["linux"]["arm64"].SHA256)
}

func TestFormulaService_buildTemplateData(t *testing.T) {
	service := &formulaService{}

	formulaInfo := models.FormulaInfo{
		ModuleName: "gorocket",
		Version:    "1.0.0",
		Repository: models.Repository{
			Owner: "koki-develop",
			Name:  "gorocket",
		},
		PlatformURLs: map[string]map[string]models.FormulaURL{
			"darwin": {
				"amd64": {URL: "https://example.com/darwin-amd64.tar.gz", SHA256: "darwin_amd64_sha256"},
				"arm64": {URL: "https://example.com/darwin-arm64.tar.gz", SHA256: "darwin_arm64_sha256"},
			},
			"linux": {
				"amd64": {URL: "https://example.com/linux-amd64.tar.gz", SHA256: "linux_amd64_sha256"},
				"arm64": {URL: "https://example.com/linux-arm64.tar.gz", SHA256: "linux_arm64_sha256"},
			},
		},
	}

	templateData := service.buildTemplateData(formulaInfo)

	assert.Equal(t, "Gorocket", templateData.ClassName)
	assert.Equal(t, "1.0.0", templateData.Version)
	assert.Equal(t, "gorocket", templateData.ModuleName)
	assert.Equal(t, "https://example.com/darwin-amd64.tar.gz", templateData.MacOSAMD64URL)
	assert.Equal(t, "darwin_amd64_sha256", templateData.MacOSAMD64SHA256)
	assert.Equal(t, "https://example.com/darwin-arm64.tar.gz", templateData.MacOSARM64URL)
	assert.Equal(t, "darwin_arm64_sha256", templateData.MacOSARM64SHA256)
	assert.Equal(t, "https://example.com/linux-amd64.tar.gz", templateData.LinuxAMD64URL)
	assert.Equal(t, "linux_amd64_sha256", templateData.LinuxAMD64SHA256)
	assert.Equal(t, "https://example.com/linux-arm64.tar.gz", templateData.LinuxARM64URL)
	assert.Equal(t, "linux_arm64_sha256", templateData.LinuxARM64SHA256)
}

func TestFormulaService_getURL(t *testing.T) {
	service := &formulaService{}

	platformURLs := map[string]map[string]models.FormulaURL{
		"darwin": {
			"amd64": {URL: "https://example.com/darwin-amd64.tar.gz", SHA256: "sha256"},
		},
	}

	tests := []struct {
		name     string
		os       string
		arch     string
		expected string
	}{
		{
			name:     "existing platform and arch",
			os:       "darwin",
			arch:     "amd64",
			expected: "https://example.com/darwin-amd64.tar.gz",
		},
		{
			name:     "existing platform, non-existing arch",
			os:       "darwin",
			arch:     "arm64",
			expected: "",
		},
		{
			name:     "non-existing platform",
			os:       "windows",
			arch:     "amd64",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getURL(platformURLs, tt.os, tt.arch)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormulaService_getSHA256(t *testing.T) {
	service := &formulaService{}

	platformURLs := map[string]map[string]models.FormulaURL{
		"darwin": {
			"amd64": {URL: "https://example.com/darwin-amd64.tar.gz", SHA256: "darwin_amd64_sha256"},
		},
	}

	tests := []struct {
		name     string
		os       string
		arch     string
		expected string
	}{
		{
			name:     "existing platform and arch",
			os:       "darwin",
			arch:     "amd64",
			expected: "darwin_amd64_sha256",
		},
		{
			name:     "existing platform, non-existing arch",
			os:       "darwin",
			arch:     "arm64",
			expected: "",
		},
		{
			name:     "non-existing platform",
			os:       "windows",
			arch:     "amd64",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.getSHA256(platformURLs, tt.os, tt.arch)
			assert.Equal(t, tt.expected, result)
		})
	}
}
