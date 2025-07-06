package mocks

import (
	"io"
	"os"
)

type MockFileSystemProvider struct {
	ReadFileFunc     func(path string) ([]byte, error)
	WriteFileFunc    func(path string, data []byte, perm os.FileMode) error
	OpenFunc         func(path string) (io.ReadCloser, error)
	CreateFunc       func(path string) (io.WriteCloser, error)
	StatFunc         func(path string) (os.FileInfo, error)
	MkdirAllFunc     func(path string, perm os.FileMode) error
	RemoveFunc       func(path string) error
	RemoveAllFunc    func(path string) error
	GetModuleNameFunc func() (string, error)
	EnsureDistDirFunc func() error
}

func (m *MockFileSystemProvider) ReadFile(path string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystemProvider) WriteFile(path string, data []byte, perm os.FileMode) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, data, perm)
	}
	return nil
}

func (m *MockFileSystemProvider) Open(path string) (io.ReadCloser, error) {
	if m.OpenFunc != nil {
		return m.OpenFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystemProvider) Create(path string) (io.WriteCloser, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystemProvider) Stat(path string) (os.FileInfo, error) {
	if m.StatFunc != nil {
		return m.StatFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystemProvider) MkdirAll(path string, perm os.FileMode) error {
	if m.MkdirAllFunc != nil {
		return m.MkdirAllFunc(path, perm)
	}
	return nil
}

func (m *MockFileSystemProvider) Remove(path string) error {
	if m.RemoveFunc != nil {
		return m.RemoveFunc(path)
	}
	return nil
}

func (m *MockFileSystemProvider) RemoveAll(path string) error {
	if m.RemoveAllFunc != nil {
		return m.RemoveAllFunc(path)
	}
	return nil
}

func (m *MockFileSystemProvider) GetModuleName() (string, error) {
	if m.GetModuleNameFunc != nil {
		return m.GetModuleNameFunc()
	}
	return "test-module", nil
}

func (m *MockFileSystemProvider) EnsureDistDir() error {
	if m.EnsureDistDirFunc != nil {
		return m.EnsureDistDirFunc()
	}
	return nil
}