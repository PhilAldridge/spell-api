package utils

import "crypto/rand"

const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
const length = 6

func RandomString() (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}