package utils

import (
	"crypto/rand"
	"math/big"
)

const base62Chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateShortCode generates a cryptographically secure random base62 string of the given length
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		length = 6
	}

	result := make([]byte, length)
	charsLen := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, charsLen)
		if err != nil {
			return "", err
		}
		result[i] = base62Chars[num.Int64()]
	}

	return string(result), nil
}
