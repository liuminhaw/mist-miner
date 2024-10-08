package shelf

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// randomHex generates a random hex string with n bytes
func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("randomHex: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
