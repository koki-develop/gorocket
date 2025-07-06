package services

import (
	"errors"
	"os"
	"testing"

	"github.com/koki-develop/gorocket/internal/providers/mocks"
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
			mockFS := &mocks.MockFileSystemProvider{
				StatFunc: func(path string) (os.FileInfo, error) {
					return nil, tt.statErr
				},
			}

			service := NewConfigService(mockFS)
			result := service.ConfigExists()

			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
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
			mockFS := &mocks.MockFileSystemProvider{
				StatFunc: func(path string) (os.FileInfo, error) {
					if tt.configExists {
						return nil, nil
					}
					return nil, os.ErrNotExist
				},
				WriteFileFunc: func(path string, data []byte, perm os.FileMode) error {
					return tt.writeFileErr
				},
			}

			service := NewConfigService(mockFS)
			err := service.CreateDefaultConfig()

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
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
			mockFS := &mocks.MockFileSystemProvider{
				ReadFileFunc: func(path string) ([]byte, error) {
					return tt.fileContent, tt.readFileErr
				},
			}

			service := NewConfigService(mockFS)
			config, err := service.LoadConfig()

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if config == nil {
					t.Errorf("expected config, but got nil")
				}
			}
		})
	}
}

func TestConfigService_GetDefaultConfigData(t *testing.T) {
	mockFS := &mocks.MockFileSystemProvider{}
	service := NewConfigService(mockFS)
	data := service.GetDefaultConfigData()

	if len(data) == 0 {
		t.Errorf("expected default config data, but got empty")
	}
}