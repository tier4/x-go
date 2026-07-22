package random_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tier4/x-go/random"
)

// wantDefaultLength mirrors the unexported defaultLength in random.go, used when a negative length is given.
const wantDefaultLength = 32

func TestGenerateAlphabets(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		length int
		regexp string
	}{
		"positive length": {length: 50, regexp: "^[a-zA-Z]{50}$"},
		"negative length": {length: -1, regexp: fmt.Sprintf("^[a-zA-Z]{%d}$", wantDefaultLength)},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Regexp(t, tt.regexp, random.GenerateAlphabets(tt.length))
		})
	}
}

func TestGenerateAlphabetsLowerCase(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		length int
		regexp string
	}{
		"positive length": {length: 50, regexp: "^[a-z]{50}$"},
		"negative length": {length: -1, regexp: fmt.Sprintf("^[a-z]{%d}$", wantDefaultLength)},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Regexp(t, tt.regexp, random.GenerateAlphabetsLowerCase(tt.length))
		})
	}
}

func TestGenerateBase58(t *testing.T) {
	t.Parallel()

	const charset = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
	tests := map[string]struct {
		length int
		regexp string
	}{
		"positive length": {length: 50, regexp: "^[" + charset + "]{50}$"},
		"negative length": {length: -1, regexp: fmt.Sprintf("^[%s]{%d}$", charset, wantDefaultLength)},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Regexp(t, tt.regexp, random.GenerateBase58(tt.length))
		})
	}
}
