package popx_test

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"testing"
	"time"

	"github.com/gobuffalo/pop/v6"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/tier4/x-go/dockertestx"
	"github.com/tier4/x-go/popx"
)

type User struct {
	ID    int    `db:"id"`
	Email string `db:"email"`
}

//go:embed testdata/migrations
var migrationFS embed.FS

func TestClient_TransactionWithTryAdvisoryLock(t *testing.T) {
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

	conn, err := pop.NewConnection(&pop.ConnectionDetails{
		URL: dsn,
	})
	require.NoError(t, err)
	mb, err := pop.NewMigrationBox(migrationFS, conn)
	require.NoError(t, err)

	cl, err := popx.New(conn, &mb)
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, cl.MigrateUp(ctx))

	tx1 := &User{
		Email: "example01@example.com",
	}
	tx2 := &User{
		Email: "example02@example.com",
	}

	key := "test"
	var eg errgroup.Group
	eg.Go(func() error {
		return cl.TransactionWithTryAdvisoryLock(ctx, key, func(ctx context.Context, conn *pop.Connection) error {
			time.Sleep(100 * time.Millisecond)
			return cl.GetConnection(ctx).Save(tx1)
		})
	})
	// to ensure to execute 1st transaction
	time.Sleep(10 * time.Millisecond)
	eg.Go(func() error {
		return cl.TransactionWithTryAdvisoryLock(ctx, key, func(ctx context.Context, conn *pop.Connection) error {
			return cl.GetConnection(ctx).Save(tx2)
		})
	})

	err = eg.Wait()
	require.Error(t, err)
	assert.True(t, errors.Is(err, popx.ErrDataLockTaken), err)

	var (
		found1 User
		found2 User
	)
	require.NoError(t, cl.GetConnection(ctx).Find(&found1, tx1.ID))
	assert.Equal(t, *tx1, found1)

	err = cl.GetConnection(ctx).Find(&found2, tx2.ID)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}
