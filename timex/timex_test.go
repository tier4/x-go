package timex_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tier4/x-go/timex"
)

func TestSetTimeNow(t *testing.T) {
	timex.SetTimeNow(func() time.Time { return time.Time{} })
	assert.Equal(t, time.Time{}, timex.Now())
}
