package providers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileSystemProvider interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte, perm os.FileMode) error
	Open(path string) (io.ReadCloser, error)
	Create(path string) (io.WriteCloser, error)
	Stat(path string) (os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	Remove(path string) error
	RemoveAll(path string) error
	GetModuleName() (string, error)
	EnsureDistDir(clean bool) error
}

type fileSystemProvider struct{}

func NewFileSystemProvider() FileSystemProvider {
	return &fileSystemProvider{}
}

func (f *fileSystemProvider) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (f *fileSystemProvider) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (f *fileSystemProvider) Open(path string) (io.ReadCloser, error) {
	return os.Open(path)
}

func (f *fileSystemProvider) Create(path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (f *fileSystemProvider) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (f *fileSystemProvider) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (f *fileSystemProvider) Remove(path string) error {
	return os.Remove(path)
}

func (f *fileSystemProvider) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (f *fileSystemProvider) GetModuleName() (string, error) {
	file, err := f.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("failed to open go.mod: %w", err)
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "module ") {
			moduleName := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			return filepath.Base(moduleName), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	return "", fmt.Errorf("module name not found in go.mod")
}

func (f *fileSystemProvider) EnsureDistDir(clean bool) error {
	distDir := "dist"
	if stat, err := f.Stat(distDir); err == nil {
		if stat.IsDir() {
			if clean {
				if err := f.RemoveAll(distDir); err != nil {
					return fmt.Errorf("failed to remove dist directory: %w", err)
				}
			} else {
				entries, err := os.ReadDir(distDir)
				if err != nil {
					return fmt.Errorf("failed to read dist directory: %w", err)
				}
				if len(entries) > 0 {
					return fmt.Errorf("dist directory is not empty. Please clean it first")
				}
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check dist directory: %w", err)
	}

	if err := f.MkdirAll(distDir, 0755); err != nil {
		return fmt.Errorf("failed to create dist directory: %w", err)
	}

	return nil
}
