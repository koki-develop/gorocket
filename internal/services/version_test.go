package services

import (
	"errors"
	"testing"

	"github.com/koki-develop/gorocket/internal/providers/mocks"
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
			mockGit := &mocks.MockGitProvider{
				GetCurrentVersionFunc: func() (string, error) {
					return tt.version, tt.versionErr
				},
			}

			mockFS := &mocks.MockFileSystemProvider{
				GetModuleNameFunc: func() (string, error) {
					return tt.moduleName, tt.moduleNameErr
				},
			}

			service := NewVersionService(mockGit, mockFS)
			buildInfo, err := service.GetBuildInfo()

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

			if buildInfo.ModuleName != tt.expectedModule {
				t.Errorf("expected module name %s, got %s", tt.expectedModule, buildInfo.ModuleName)
			}

			if buildInfo.Version != tt.expectedVersion {
				t.Errorf("expected version %s, got %s", tt.expectedVersion, buildInfo.Version)
			}
		})
	}
}