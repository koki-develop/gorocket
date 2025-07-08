package util

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CalculateSHA256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "test data",
			input:    "test data for sha256",
			expected: "2284de595a972f7093c632e953e387ce70b7e7bbab1e65a60c3b64b9b70c9db6",
		},
		{
			name:     "multiline text",
			input:    "line1\nline2\nline3",
			expected: "6bb6a5ad9b9c43a7cb535e636578716b64ac42edea814a4cad102ba404946837",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := CalculateSHA256(reader)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
