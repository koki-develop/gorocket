package mocks

type MockCommandProvider struct {
	RunFunc           func(name string, args ...string) (string, error)
	RunWithEnvFunc    func(name string, env []string, args ...string) error
	BuildBinaryFunc   func(moduleName, version, osName, arch string) (string, error)
}

func (m *MockCommandProvider) Run(name string, args ...string) (string, error) {
	if m.RunFunc != nil {
		return m.RunFunc(name, args...)
	}
	return "", nil
}

func (m *MockCommandProvider) RunWithEnv(name string, env []string, args ...string) error {
	if m.RunWithEnvFunc != nil {
		return m.RunWithEnvFunc(name, env, args...)
	}
	return nil
}

func (m *MockCommandProvider) BuildBinary(moduleName, version, osName, arch string) (string, error) {
	if m.BuildBinaryFunc != nil {
		return m.BuildBinaryFunc(moduleName, version, osName, arch)
	}
	return "dist/test-binary", nil
}