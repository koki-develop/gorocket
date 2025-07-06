package services

import (
	"errors"
	"testing"

	"github.com/koki-develop/gorocket/internal/providers/mocks"
	"github.com/stretchr/testify/assert"
)

func TestVersionService_GetBuildInfo(t *testing.T) {
	tests := []struct {
		name           string
		moduleName     string
		moduleNameErr  error
		version        string
		versionErr     error
		expectedModule string
		expectedVersion string
		expectedError  bool
	}{
		{
			name:            "successful build info retrieval",
			moduleName:      "test-module",
			moduleNameErr:   nil,
			version:         "v1.0.0",
			versionErr:      nil,
			expectedModule:  "test-module",
			expectedVersion: "v1.0.0",
			expectedError:   false,
		},
		{
			name:          "module name error",
			moduleName:    "",
			moduleNameErr: errors.New("module name error"),
			version:       "v1.0.0",
			versionErr:    nil,
			expectedError: true,
		},
		{
			name:          "version error",
			moduleName:    "test-module",
			moduleNameErr: nil,
			version:       "",
			versionErr:    errors.New("version error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGit := &mocks.MockGitProvider{}
			mockGit.On("GetCurrentVersion").Return(tt.version, tt.versionErr)

			mockFS := &mocks.MockFileSystemProvider{}
			mockFS.On("GetModuleName").Return(tt.moduleName, tt.moduleNameErr)

			service := NewVersionService(mockGit, mockFS)
			buildInfo, err := service.GetBuildInfo()

			if tt.expectedError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedModule, buildInfo.ModuleName)
			assert.Equal(t, tt.expectedVersion, buildInfo.Version)
		})
	}
}