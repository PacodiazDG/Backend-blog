package security

import (
	"crypto/rand"
	"encoding/base64"
)

// Generates random String
// Not secure for cryptographic and related implementations
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Generates random String
// Not secure for cryptographic and related implementations
func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}
