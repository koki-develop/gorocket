package mocks

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type MockWriteCloser struct {
	mock.Mock
}

type _m_WriteCloserRecorder struct {
	mock *MockWriteCloser
}

func NewMockWriteCloser(t *testing.T) *MockWriteCloser {
	m := &MockWriteCloser{}
	m.Mock.Test(t)
	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}

func (m *MockWriteCloser) EXPECT() *_m_WriteCloserRecorder {
	return &_m_WriteCloserRecorder{mock: m}
}

func (m *MockWriteCloser) Write(p []byte) (n int, err error) {
	args := m.Called(p)
	return args.Int(0), args.Error(1)
}

func (m *MockWriteCloser) Close() error {
	args := m.Called()
	return args.Error(0)
}

type WriteCall struct {
	*mock.Call
}

func (r *_m_WriteCloserRecorder) Write(p interface{}) *WriteCall {
	call := r.mock.On("Write", p)
	return &WriteCall{Call: call}
}

func (c *WriteCall) Return(n int, err error) *WriteCall {
	c.Call.Return(n, err)
	return c
}

type CloseCall struct {
	*mock.Call
}

func (r *_m_WriteCloserRecorder) Close() *CloseCall {
	call := r.mock.On("Close")
	return &CloseCall{Call: call}
}

func (c *CloseCall) Return(err error) *CloseCall {
	c.Call.Return(err)
	return c
}
