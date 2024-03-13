package random_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tier4/x-go/random"
)

func TestGenerateAlphabets(t *testing.T) {
	t.Parallel()
	assert.Regexp(t, "^[a-zA-Z]{50}$", random.GenerateAlphabets(50))
}

func TestGenerateAlphabetsLowerCase(t *testing.T) {
	t.Parallel()
	assert.Regexp(t, "^[a-z]{50}$", random.GenerateAlphabetsLowerCase(50))
}

func TestGenerateBase58(t *testing.T) {
	t.Parallel()
	assert.Regexp(t, "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{50}$", random.GenerateBase58(50))
}
