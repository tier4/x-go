package random

import (
	"crypto/rand"
	"io"
)

func GenerateAlphabets(length int) string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	return generateRandomString(letters, length)
}

func GenerateAlphabetsLowerCase(length int) string {
	const letters = "abcdefghijkmnopqrstuvwxyz"
	return generateRandomString(letters, length)
}

func GenerateBase58(length int) string {
	const letters = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	return generateRandomString(letters, length)
}

const defaultLength = 32

func generateRandomString(letters string, length int) string {
	if length < 0 {
		length = defaultLength
	}

	b := make([]byte, length)
	_, _ = io.ReadFull(rand.Reader, b)

	for i, v := range b {
		b[i] = letters[int(v)%len(letters)]
	}
	return string(b)
}
