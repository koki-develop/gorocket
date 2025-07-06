package mocks

type MockGitProvider struct {
	GetCurrentVersionFunc func() (string, error)
}

func (m *MockGitProvider) GetCurrentVersion() (string, error) {
	if m.GetCurrentVersionFunc != nil {
		return m.GetCurrentVersionFunc()
	}
	return "v1.0.0", nil
}