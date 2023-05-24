package dockertestx_test

import (
	"database/sql"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/tier4/x-go/dockertestx"
)

func TestPool_NewMysql(t *testing.T) {
	dsnRegex, err := regexp.Compile(`root:passw0rd@tcp\(localhost:\d{4,5}\)/test`)
	require.NoError(t, err)

	t.Run("no option", func(t *testing.T) {
		t.Parallel()

		p, err := dockertestx.New(dockertestx.PoolOption{})
		require.NoError(t, err)
		t.Cleanup(func() {
			require.NoError(t, p.Purge())
		})

		dsn, err := p.NewResource(new(dockertestx.MysqlFactory), dockertestx.ContainerOption{
			Tag: "8.0",
		})
		require.NoError(t, err)

		assert.Regexp(t, dsnRegex, dsn)
		db, err := sql.Open("mysql", dsn)
		require.NoError(t, err)
		assert.NoError(t, db.Ping())
		require.NoError(t, db.Close())

	})
}
