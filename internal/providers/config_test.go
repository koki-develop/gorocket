package providers

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigProvider_LoadConfig_WithTemplate(t *testing.T) {
	tests := []struct {
		name          string
		configContent string
		templateData  *models.TemplateData
		expected      *models.Config
		setupMocks    func(mockFS *mocks.MockFileSystemProvider)
		expectedError string
	}{
		{
			name: "successful template processing",
			configContent: `build:
  ldflags: "-X main.version={{.Version}} -X main.module={{.Module}}"
  targets:
    - os: linux
      arch: [amd64, arm64]
brew:
  repository:
    owner: "{{.Module}}"
    name: "homebrew-tap"`,
			templateData: &models.TemplateData{
				Version: "v1.0.0",
				Module:  "github.com/example/project",
			},
			expected: &models.Config{
				Build: models.BuildConfig{
					LdFlags: "-X main.version=v1.0.0 -X main.module=github.com/example/project",
					Targets: []models.Target{
						{
							OS:   "linux",
							Arch: []string{"amd64", "arm64"},
						},
					},
				},
				Brew: &models.BrewConfig{
					Repository: models.Repository{
						Owner: "github.com/example/project",
						Name:  "homebrew-tap",
					},
				},
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				reader := io.NopCloser(strings.NewReader(`build:
  ldflags: "-X main.version={{.Version}} -X main.module={{.Module}}"
  targets:
    - os: linux
      arch: [amd64, arm64]
brew:
  repository:
    owner: "{{.Module}}"
    name: "homebrew-tap"`))
				mockFS.EXPECT().Open(".gorocket.yml").Return(reader, nil)
			},
		},
		{
			name: "template parse error",
			configContent: `build:
  ldflags: "{{.InvalidSyntax}"`,
			templateData: &models.TemplateData{
				Version: "v1.0.0",
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				reader := io.NopCloser(strings.NewReader(`build:
  ldflags: "{{.InvalidSyntax}"`))
				mockFS.EXPECT().Open(".gorocket.yml").Return(reader, nil)
			},
			expectedError: "failed to parse config template",
		},
		{
			name: "template variable does not exist",
			configContent: `build:
  ldflags: "{{.NonExistentField}}"`,
			templateData: &models.TemplateData{
				Version: "v1.0.0",
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				reader := io.NopCloser(strings.NewReader(`build:
  ldflags: "{{.NonExistentField}}"`))
				mockFS.EXPECT().Open(".gorocket.yml").Return(reader, nil)
			},
			expectedError: "failed to execute config template",
		},
		{
			name: "invalid YAML after template processing",
			configContent: `build:
  ldflags: "{{.Version}}"
  targets:
    invalid yaml content`,
			templateData: &models.TemplateData{
				Version: "v1.0.0",
			},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				reader := io.NopCloser(strings.NewReader(`build:
  ldflags: "{{.Version}}"
  targets:
    invalid yaml content`))
				mockFS.EXPECT().Open(".gorocket.yml").Return(reader, nil)
			},
			expectedError: "failed to decode config YAML",
		},
		{
			name:         "file open error",
			templateData: &models.TemplateData{},
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Open(".gorocket.yml").Return(nil, errors.New("file not found"))
			},
			expectedError: "file not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			tt.setupMocks(mockFS)

			provider := NewConfigProvider(mockFS)
			result, err := provider.LoadConfig(tt.templateData)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				if tt.expected != nil {
					assert.Equal(t, tt.expected, result)
				}
			}

			mockFS.AssertExpectations(t)
		})
	}
}

func TestConfigProvider_LoadConfig(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(mockFS *mocks.MockFileSystemProvider)
		expected      *models.Config
		expectedError string
	}{
		{
			name: "successful config load",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				configYAML := `build:
  targets:
    - os: linux
      arch: [amd64]`
				reader := io.NopCloser(strings.NewReader(configYAML))
				mockFS.EXPECT().Open(".gorocket.yml").Return(reader, nil)
			},
			expected: &models.Config{
				Build: models.BuildConfig{
					Targets: []models.Target{
						{
							OS:   "linux",
							Arch: []string{"amd64"},
						},
					},
				},
			},
		},
		{
			name: "file open error",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Open(".gorocket.yml").Return(nil, errors.New("file not found"))
			},
			expectedError: "file not found",
		},
		{
			name: "invalid YAML",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				invalidYAML := `build:
  targets:
    invalid yaml`
				reader := io.NopCloser(strings.NewReader(invalidYAML))
				mockFS.EXPECT().Open(".gorocket.yml").Return(reader, nil)
			},
			expectedError: "yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			tt.setupMocks(mockFS)

			provider := NewConfigProvider(mockFS)
			result, err := provider.LoadConfig(nil)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			mockFS.AssertExpectations(t)
		})
	}
}

func TestConfigProvider_ConfigExists(t *testing.T) {
	tests := []struct {
		name       string
		setupMocks func(mockFS *mocks.MockFileSystemProvider)
		expected   bool
	}{
		{
			name: "config exists",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, nil)
			},
			expected: true,
		},
		{
			name: "config does not exist",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, errors.New("file not found"))
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			tt.setupMocks(mockFS)

			provider := NewConfigProvider(mockFS)
			result := provider.ConfigExists()

			assert.Equal(t, tt.expected, result)
			mockFS.AssertExpectations(t)
		})
	}
}

func TestConfigProvider_CreateDefaultConfig(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(mockFS *mocks.MockFileSystemProvider)
		expectedError string
	}{
		{
			name: "successful creation",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, errors.New("not found"))
				mockFS.EXPECT().WriteFile(".gorocket.yml", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "config already exists",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, nil)
			},
			expectedError: ".gorocket.yml already exists",
		},
		{
			name: "write error",
			setupMocks: func(mockFS *mocks.MockFileSystemProvider) {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, errors.New("not found"))
				mockFS.EXPECT().WriteFile(".gorocket.yml", mock.Anything, mock.Anything).Return(errors.New("permission denied"))
			},
			expectedError: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			tt.setupMocks(mockFS)

			provider := NewConfigProvider(mockFS)
			err := provider.CreateDefaultConfig()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockFS.AssertExpectations(t)
		})
	}
}
