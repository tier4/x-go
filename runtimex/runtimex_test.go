package runtimex_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tier4/x-go/runtimex"
)

func TestMaxParallelism(t *testing.T) {
	assert.Greater(t, runtimex.MaxParallelism(), 0)
}
