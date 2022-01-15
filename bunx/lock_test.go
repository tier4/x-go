package bunx_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"

	"github.com/tier4/x-go/bunx"
)

func TestClient_TransactionWithTryAdvisoryLock(t *testing.T) {
	ctx := context.Background()
	user1 := &User{Email: "TransactionWithTryAdvisoryLock_01@example.com"}
	user2 := &User{Email: "TransactionWithTryAdvisoryLock_02@example.com"}

	key := "test"
	var eg errgroup.Group
	eg.Go(func() error {
		return client.TransactionWithTryAdvisoryLock(ctx, key, func(ctx context.Context, tx bun.Tx) error {
			time.Sleep(100 * time.Millisecond)
			_, err := tx.NewInsert().Model(user1).Exec(ctx)
			return err
		})
	})
	// to ensure to start 1st transaction
	time.Sleep(10 * time.Millisecond)
	eg.Go(func() error {
		return client.TransactionWithTryAdvisoryLock(ctx, key, func(ctx context.Context, tx bun.Tx) error {
			_, err := tx.NewInsert().Model(user2).Exec(ctx)
			return err
		})
	})

	err := eg.Wait()
	require.Error(t, err)
	assert.ErrorIs(t, err, bunx.ErrDataLockTaken)

	var (
		found1 User
		found2 User
	)

	require.NoError(t, client.DB(ctx).NewSelect().Model(&found1).Where("id = ?", user1.ID).Scan(ctx))
	assert.Equal(t, *user1, found1)

	err = client.DB(ctx).NewSelect().Model(&found2).Where("id = ?", user2.ID).Scan(ctx)
	assert.True(t, errors.Is(err, sql.ErrNoRows))
}
