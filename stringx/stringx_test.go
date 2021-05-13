package stringx_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tier4/x-go/stringx"
)

func TestRef(t *testing.T) {
	s := "test"
	assert.Equal(t, stringx.Ref(s), &s)
}

func TestDeref(t *testing.T) {
	s := "test"
	assert.Equal(t, stringx.Deref(&s), s)
}

func TestFormatDateTime(t *testing.T) {
	v := time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)
	assert.Equal(t, "2006-01-02T15:04:05Z", stringx.FormatDateTime(v))
}
