package utils

import (
	"crypto/rand"
	"math/big"
)

func SecureRandomString(length int, alphabet string) (string, error) {
	bytes := make([]byte, length)
	for i := range length {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		bytes[i] = alphabet[num.Int64()]
	}
	return string(bytes), nil
}

func RandomString(length int, alphabet string) string {
	result := make([]byte, length)
	for i := range length {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			result[i] = 'a'
			continue
		}
		result[i] = alphabet[num.Int64()]
	}
	return string(result)
}
