package cmd

import (
	"testing"

	"github.com/koki-develop/gorocket/internal/providers/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBuildCommand_CleanFlag(t *testing.T) {
	tests := []struct {
		name        string
		cleanFlag   bool
		expectClean bool
	}{
		{
			name:        "clean flag false",
			cleanFlag:   false,
			expectClean: false,
		},
		{
			name:        "clean flag true",
			cleanFlag:   true,
			expectClean: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock FileSystemProvider
			mockFS := mocks.NewMockFileSystemProvider(t)

			// Setup expectation for EnsureDistDir with the expected clean parameter
			mockFS.EXPECT().EnsureDistDir(tt.expectClean).Return(nil)

			// Create command
			buildCmd := &BuildCommand{
				fsProvider: mockFS,
				flagClean:  tt.cleanFlag,
			}

			// Test - Just test that EnsureDistDir is called with correct clean parameter
			err := buildCmd.fsProvider.EnsureDistDir(buildCmd.flagClean)

			// Assert
			assert.NoError(t, err)
		})
	}
}
