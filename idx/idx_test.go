package idx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tier4/x-go/idx"
)

func TestIsUUIDString(t *testing.T) {
	for name, c := range map[string]struct {
		arg      string
		expected bool
	}{
		"valid": {
			arg:      "f792d50f-8e11-44ae-8adc-d7b8aeaeecf3",
			expected: true,
		},
		"valid uppercase": {
			arg:      "F792D50F-8E11-44AE-8ADC-D7B8AEAEECF3",
			expected: true,
		},
		"invalid character": {
			arg:      "z792d50f-8e11-44ae-8adc-d7b8aeaeecf3",
			expected: false,
		},
		"invalid length": {
			arg:      "f792d50f-8e11-44ae-8adc-d7b8aeaeecf",
			expected: false,
		},
		"invalid version": {
			arg:      "f792d50f-8e11-94ae-8adc-d7b8aeaeecf3",
			expected: false,
		},
		"invalid variant": {
			arg:      "f792d50f-8e11-44ae-cadc-d7b8aeaeecf3",
			expected: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, c.expected, idx.IsUUIDString(c.arg))
		})
	}
}

func TestNewUUID(t *testing.T) {
	assert.True(t, idx.IsUUIDString(idx.NewUUID()))
}

func TestShortID(t *testing.T) {
	assert.Len(t, idx.ShortID(), 8)
}
