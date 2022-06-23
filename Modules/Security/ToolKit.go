package Security

import (
	"crypto/rand"
	"encoding/base64"
	"unicode"
)

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func IsSafePassword(Password string) bool {
	var Uppercase bool
	var lowercase bool
	var Number bool
	var Digit bool

	for _, v := range Password {
		if unicode.IsUpper(v) {
			Uppercase = true
		}
		if unicode.IsLower(v) {
			lowercase = true
		}
		if unicode.IsNumber(v) {
			Number = true
		}

	}

	return Uppercase && lowercase && Number && Digit
}
