package bunx_test

import (
	"context"
	"database/sql"
	"embed"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"golang.org/x/sync/errgroup"

	"github.com/tier4/x-go/bunx"
	"github.com/tier4/x-go/dockertestx"
)

type User struct {
	ID    int64  `db:"id,pk"`
	Email string `db:"email"`
}

//go:embed testdata/migrations/*.sql
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

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())

	cl, err := bunx.NewClient(db)
	require.NoError(t, err)
	migrator, err := bunx.NewMigrator(db, migrationFS, bunx.NewNoopLogger())
	require.NoError(t, err)

	ctx := context.Background()

	_ = migrator.Reset(ctx)
	require.NoError(t, migrator.Migrate(ctx))

	tx1 := &User{
		Email: "example01@example.com",
	}
	tx2 := &User{
		Email: "example02@example.com",
	}

	key := "test"
	var eg errgroup.Group
	eg.Go(func() error {
		return cl.TransactionWithTryAdvisoryLock(ctx, key, func(ctx context.Context, tx bun.Tx) error {
			time.Sleep(100 * time.Millisecond)
			_, err := tx.NewInsert().Model(tx1).Exec(ctx)
			return err
		})
	})
	// to ensure to start 1st transaction
	time.Sleep(10 * time.Millisecond)
	eg.Go(func() error {
		return cl.TransactionWithTryAdvisoryLock(ctx, key, func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().Model(tx2).Exec(ctx)
			return err
		})
	})

	err = eg.Wait()
	require.Error(t, err)
	assert.ErrorIs(t, err, bunx.ErrDataLockTaken)

	var (
		found1 User
		found2 User
	)

	require.NoError(t, cl.DB().NewSelect().Model(&found1).Where("id = ?", tx1.ID).Scan(ctx))
	assert.Equal(t, *tx1, found1)

	err = cl.DB().NewSelect().Model(&found2).Where("id = ?", tx2.ID).Scan(ctx)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}
