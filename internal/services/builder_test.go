package services

import (
	"errors"
	"testing"

	"github.com/koki-develop/gorocket/internal/models"
	"github.com/koki-develop/gorocket/internal/providers/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuilderService_BuildTargets(t *testing.T) {
	tests := []struct {
		name            string
		buildInfo       *models.BuildInfo
		targets         []models.Target
		buildBinaryErr  error
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
			mockCommand := mocks.NewMockCommandProvider(t)
			if tt.buildBinaryErr != nil {
				mockCommand.EXPECT().BuildBinary(mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("", tt.buildBinaryErr).Times(tt.expectedResults)
			} else {
				mockCommand.EXPECT().BuildBinary(mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return("dist/binary", nil).Times(tt.expectedResults)
			}

			mockFS := mocks.NewMockFileSystemProvider(t)

			service := NewBuilderService(mockCommand, mockFS)
			results, err := service.BuildTargets(tt.buildInfo, tt.targets)

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, results, tt.expectedResults)

			if tt.buildBinaryErr != nil {
				for _, result := range results {
					assert.Error(t, result.Error)
				}
			} else {
				for _, result := range results {
					assert.NoError(t, result.Error)
					assert.NotEmpty(t, result.BinaryPath)
				}
			}
		})
	}
}
