package services

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfigService_ConfigExists(t *testing.T) {
	tests := []struct {
		name     string
		statErr  error
		expected bool
	}{
		{
			name:     "config exists",
			statErr:  nil,
			expected: true,
		},
		{
			name:     "config does not exist",
			statErr:  os.ErrNotExist,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, tt.statErr)

			service := NewConfigService(mockFS)
			result := service.ConfigExists()

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfigService_CreateDefaultConfig(t *testing.T) {
	tests := []struct {
		name          string
		configExists  bool
		writeFileErr  error
		expectedError bool
	}{
		{
			name:          "successful creation",
			configExists:  false,
			writeFileErr:  nil,
			expectedError: false,
		},
		{
			name:          "config already exists",
			configExists:  true,
			writeFileErr:  nil,
			expectedError: true,
		},
		{
			name:          "write file error",
			configExists:  false,
			writeFileErr:  errors.New("write error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			if tt.configExists {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, nil)
			} else {
				mockFS.EXPECT().Stat(".gorocket.yml").Return(nil, os.ErrNotExist)
				mockFS.EXPECT().WriteFile(".gorocket.yml", mock.AnythingOfType("[]uint8"), os.FileMode(0644)).Return(tt.writeFileErr)
			}

			service := NewConfigService(mockFS)
			err := service.CreateDefaultConfig()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfigService_LoadConfig(t *testing.T) {
	validYAML := `build:
  targets:
    - os: linux
      arch: [amd64]`

	tests := []struct {
		name          string
		fileContent   []byte
		readFileErr   error
		expectedError bool
	}{
		{
			name:          "successful load",
			fileContent:   []byte(validYAML),
			readFileErr:   nil,
			expectedError: false,
		},
		{
			name:          "read file error",
			fileContent:   nil,
			readFileErr:   errors.New("read error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			if tt.readFileErr != nil {
				mockFS.EXPECT().Open(".gorocket.yml").Return(nil, tt.readFileErr)
			} else {
				mockFS.EXPECT().Open(".gorocket.yml").Return(io.NopCloser(strings.NewReader(string(tt.fileContent))), nil)
			}

			service := NewConfigService(mockFS)
			config, err := service.LoadConfig(nil)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
			}
		})
	}
}

func TestConfigService_GetDefaultConfigData(t *testing.T) {
	mockFS := mocks.NewMockFileSystemProvider(t)
	service := NewConfigService(mockFS)
	data := service.GetDefaultConfigData()

	assert.NotEmpty(t, data)
}

func TestConfigService_LoadConfig_WithTemplate(t *testing.T) {
	validTemplateYAML := `build:
  ldflags: "-X main.version={{.Version}}"
  targets:
    - os: linux
      arch: [amd64]`

	tests := []struct {
		name          string
		fileContent   []byte
		templateData  *models.TemplateData
		readFileErr   error
		expectedError bool
	}{
		{
			name:        "successful template processing",
			fileContent: []byte(validTemplateYAML),
			templateData: &models.TemplateData{
				Version: "v1.0.0",
				Module:  "test-module",
			},
			readFileErr:   nil,
			expectedError: false,
		},
		{
			name:          "file read error",
			fileContent:   nil,
			templateData:  &models.TemplateData{},
			readFileErr:   errors.New("read error"),
			expectedError: true,
		},
		{
			name: "invalid template syntax",
			fileContent: []byte(`build:
  ldflags: "{{.InvalidSyntax}"`),
			templateData:  &models.TemplateData{},
			readFileErr:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFS := mocks.NewMockFileSystemProvider(t)
			if tt.readFileErr != nil {
				mockFS.EXPECT().Open(".gorocket.yml").Return(nil, tt.readFileErr)
			} else {
				mockFS.EXPECT().Open(".gorocket.yml").Return(io.NopCloser(strings.NewReader(string(tt.fileContent))), nil)
			}

			service := NewConfigService(mockFS)
			config, err := service.LoadConfig(tt.templateData)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, config)
				// Verify template was processed
				if tt.name == "successful template processing" {
					assert.Equal(t, "-X main.version=v1.0.0", config.Build.LdFlags)
				}
			}
		})
	}
}
