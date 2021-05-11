package dockertestx_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x/dockertestx"
)

func TestNewPostgres(t *testing.T) {
	t.Parallel()

	conn, purge, err := dockertestx.NewPostgres("13.2-alpine")
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, purge())
	})

	assert.Regexp(t, `postgres:\/\/dockertest:passw0rd@localhost:\d{4,5}/test\?sslmode=disable`, conn.String())

	type pinger interface {
		Ping() error
	}
	assert.NoError(t, conn.Store.(pinger).Ping())
}
