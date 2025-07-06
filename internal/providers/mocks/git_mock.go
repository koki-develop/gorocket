package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockGitProvider struct {
	mock.Mock
}

func (m *MockGitProvider) GetCurrentVersion() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}