package services

import (
	"errors"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers/mocks"
)

func TestBuilderService_BuildTargets(t *testing.T) {
	tests := []struct {
		name           string
		buildInfo      *models.BuildInfo
		targets        []models.Target
		buildBinaryErr error
		expectedResults int
		expectedError   bool
	}{
		{
			name: "successful build",
			buildInfo: &models.BuildInfo{
				ModuleName: "test-module",
				Version:    "v1.0.0",
			},
			targets: []models.Target{
				{OS: "linux", Arch: []string{"amd64", "arm64"}},
				{OS: "windows", Arch: []string{"amd64"}},
			},
			buildBinaryErr:  nil,
			expectedResults: 3,
			expectedError:   false,
		},
		{
			name: "build with error",
			buildInfo: &models.BuildInfo{
				ModuleName: "test-module",
				Version:    "v1.0.0",
			},
			targets: []models.Target{
				{OS: "linux", Arch: []string{"amd64"}},
			},
			buildBinaryErr:  errors.New("build error"),
			expectedResults: 1,
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCommand := &mocks.MockCommandProvider{
				BuildBinaryFunc: func(moduleName, version, osName, arch string) (string, error) {
					if tt.buildBinaryErr != nil {
						return "", tt.buildBinaryErr
					}
					return "dist/binary", nil
				},
			}

			mockFS := &mocks.MockFileSystemProvider{}

			service := NewBuilderService(mockCommand, mockFS)
			results, err := service.BuildTargets(tt.buildInfo, tt.targets)

			if tt.expectedError {
				if err == nil {
					t.Errorf("expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(results) != tt.expectedResults {
				t.Errorf("expected %d results, got %d", tt.expectedResults, len(results))
			}

			if tt.buildBinaryErr != nil {
				for _, result := range results {
					if result.Error == nil {
						t.Errorf("expected error in result, but got nil")
					}
				}
			} else {
				for _, result := range results {
					if result.Error != nil {
						t.Errorf("unexpected error in result: %v", result.Error)
					}
					if result.BinaryPath == "" {
						t.Errorf("expected binary path, but got empty string")
					}
				}
			}
		})
	}
}