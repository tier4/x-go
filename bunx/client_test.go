package bunx_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
)

func TestClient_Transaction(t *testing.T) {
	_ = migrator.Reset(context.Background())
	require.NoError(t, migrator.Migrate(context.Background()))

	t.Run("succeed", func(t *testing.T) {
		ctx := context.Background()
		user1 := &User{Email: "Transaction_01@example.com"}
		user2 := &User{Email: "Transaction_02@example.com"}
		err := client.Transaction(ctx, func(ctx context.Context, tx bun.Tx) error {
			if _, err := client.DB(ctx).NewInsert().Model(user1).Exec(ctx); err != nil {
				return err
			}
			if _, err := client.DB(ctx).NewInsert().Model(user2).Exec(ctx); err != nil {
				return err
			}
			return nil
		})
		require.NoError(t, err)

		{
			ctx := context.Background()
			found := new(User)
			require.NoError(t, client.DB(ctx).NewSelect().Model(found).Where("id = ?", user1.ID).Scan(ctx))
			assert.EqualValues(t, user1, found)
		}
		{
			ctx := context.Background()
			found := new(User)
			require.NoError(t, client.DB(ctx).NewSelect().Model(found).Where("id = ?", user2.ID).Scan(ctx))
			assert.EqualValues(t, user2, found)
		}
	})
	t.Run("fail", func(t *testing.T) {
		ctx := context.Background()
		user1 := &User{Email: "Transaction_03@example.com"}
		user2 := &User{Email: "Transaction_04@example.com"}
		err := client.Transaction(ctx, func(ctx context.Context, tx bun.Tx) error {
			if _, err := client.DB(ctx).NewInsert().Model(user1).Exec(ctx); err != nil {
				return err
			}
			if _, err := client.DB(ctx).NewInsert().Model(user2).Exec(ctx); err != nil {
				return err
			}
			return errors.New("something error occurred")
		})
		require.Error(t, err)

		{
			ctx := context.Background()
			found := new(User)
			assert.ErrorIs(t, client.DB(ctx).NewSelect().Model(found).Where("id = ?", user1.ID).Scan(ctx), sql.ErrNoRows)
		}
		{
			ctx := context.Background()
			found := new(User)
			assert.ErrorIs(t, client.DB(ctx).NewSelect().Model(found).Where("id = ?", user2.ID).Scan(ctx), sql.ErrNoRows)
		}
	})
}
