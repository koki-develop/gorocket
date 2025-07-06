package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSystemProvider_EnsureDistDir(t *testing.T) {
	tests := []struct {
		name       string
		clean      bool
		setup      func(t *testing.T, tmpDir string)
		expectErr  bool
		expectFile string
	}{
		{
			name:  "create dist directory when not exists",
			clean: false,
			setup: func(t *testing.T, tmpDir string) {
				// No setup needed - dist directory doesn't exist
			},
			expectErr: false,
		},
		{
			name:  "success when dist directory exists and empty",
			clean: false,
			setup: func(t *testing.T, tmpDir string) {
				err := os.MkdirAll(filepath.Join(tmpDir, "dist"), 0755)
				require.NoError(t, err)
			},
			expectErr: false,
		},
		{
			name:  "error when dist directory exists and not empty without clean",
			clean: false,
			setup: func(t *testing.T, tmpDir string) {
				distDir := filepath.Join(tmpDir, "dist")
				err := os.MkdirAll(distDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(distDir, "test.txt"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			expectErr: true,
		},
		{
			name:  "success when dist directory exists and not empty with clean",
			clean: true,
			setup: func(t *testing.T, tmpDir string) {
				distDir := filepath.Join(tmpDir, "dist")
				err := os.MkdirAll(distDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(distDir, "test.txt"), []byte("test"), 0644)
				require.NoError(t, err)
			},
			expectErr:  false,
			expectFile: "test.txt",
		},
		{
			name:  "create dist directory when not exists with clean",
			clean: true,
			setup: func(t *testing.T, tmpDir string) {
				// No setup needed - dist directory doesn't exist
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test
			tmpDir := t.TempDir()
			originalDir, err := os.Getwd()
			require.NoError(t, err)

			// Change to temporary directory
			err = os.Chdir(tmpDir)
			require.NoError(t, err)

			// Restore original directory after test
			defer func() {
				err := os.Chdir(originalDir)
				require.NoError(t, err)
			}()

			// Setup test environment
			tt.setup(t, tmpDir)

			// Create provider and test
			provider := NewFileSystemProvider()
			err = provider.EnsureDistDir(tt.clean)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Check that dist directory exists
			distDir := filepath.Join(tmpDir, "dist")
			stat, err := os.Stat(distDir)
			assert.NoError(t, err)
			assert.True(t, stat.IsDir())

			// If clean was used and there was an expected file, it should be gone
			if tt.clean && tt.expectFile != "" {
				filePath := filepath.Join(distDir, tt.expectFile)
				_, err := os.Stat(filePath)
				assert.True(t, os.IsNotExist(err), "Expected file %s to be removed", filePath)
			}
		})
	}
}

func TestFileSystemProvider_CalculateSHA256(t *testing.T) {
	tests := []struct {
		name        string
		fileContent string
		expectedSHA string
		expectErr   bool
	}{
		{
			name:        "calculate SHA256 for simple text",
			fileContent: "hello world",
			expectedSHA: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
			expectErr:   false,
		},
		{
			name:        "calculate SHA256 for empty file",
			fileContent: "",
			expectedSHA: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			expectErr:   false,
		},
		{
			name:        "calculate SHA256 for binary-like content",
			fileContent: "\x00\x01\x02\x03\xFF",
			expectedSHA: "ff5d8507b6a72bee2debce2c0054798deaccdc5d8a1b945b6280ce8aa9cba52e",
			expectErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use strings.Reader instead of temporary file
			reader := strings.NewReader(tt.fileContent)

			// Test SHA256 calculation
			provider := NewFileSystemProvider()
			result, err := provider.CalculateSHA256(reader)

			if tt.expectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSHA, result)
		})
	}
}
