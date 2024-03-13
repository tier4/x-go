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

func generateRandomString(letters string, length int) string {
	b := make([]byte, length)
	_, _ = io.ReadFull(rand.Reader, b)

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}
	return result
}
