package dockertestx_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_NewPostgres(t *testing.T) {
	dsnRegex, err := regexp.Compile(`postgres://dockertestx:passw0rd@localhost:\d{4,5}/test\?sslmode=disable`)
	require.NoError(t, err)

	t.Run("no option", func(t *testing.T) {
		t.Parallel()

		p, err := dockertestx.New(dockertestx.PoolOption{})
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, p.Purge())
		})

		dsn, err := p.NewResource(new(dockertestx.PostgresFactory), dockertestx.ContainerOption{
			Tag: "alpine",
		})
		require.NoError(t, err)

		assert.Regexp(t, dsnRegex, dsn)
		db, err := sql.Open("postgres", dsn)
		require.NoError(t, err)
		assert.NoError(t, db.Ping())
		require.NoError(t, db.Close())

	})
}
