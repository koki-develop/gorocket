package mocks

import (
	"github.com/stretchr/testify/mock"
)

type MockCommandProvider struct {
	mock.Mock
}

func (m *MockCommandProvider) Run(name string, args ...string) (string, error) {
	arguments := []interface{}{name}
	for _, arg := range args {
		arguments = append(arguments, arg)
	}
	callArgs := m.Called(arguments...)
	return callArgs.String(0), callArgs.Error(1)
}

func (m *MockCommandProvider) RunWithEnv(name string, env []string, args ...string) error {
	arguments := []interface{}{name, env}
	for _, arg := range args {
		arguments = append(arguments, arg)
	}
	callArgs := m.Called(arguments...)
	return callArgs.Error(0)
}

func (m *MockCommandProvider) BuildBinary(moduleName, version, osName, arch string) (string, error) {
	args := m.Called(moduleName, version, osName, arch)
	return args.String(0), args.Error(1)
}