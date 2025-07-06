package services

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

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
			config, err := service.LoadConfig()

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
