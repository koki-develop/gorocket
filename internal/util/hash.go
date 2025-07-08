package util

import (
	"crypto/sha256"
	"fmt"
	"io"
)

// CalculateSHA256 calculates SHA256 hash from an io.Reader
func CalculateSHA256(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
