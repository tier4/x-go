package bunx

import (
	"context"
	"errors"
	"fmt"

	"github.com/uptrace/bun"
)

var (
	ErrDataLockTaken = errors.New("data lock taken")
)

// TransactionWithTryAdvisoryLock is Transaction with pg_try_advisory_xact_lock
// if a lock has already taken, returns error immediately
func (c *Client) TransactionWithTryAdvisoryLock(ctx context.Context, key string, callback func(ctx context.Context, tx bun.Tx) error) error {
	return c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := tryTakeAdvisoryLock(ctx, tx, key); err != nil {
			return err
		}
		return callback(ctx, tx)
	})
}

func tryTakeAdvisoryLock(ctx context.Context, tx bun.Tx, key string) error {
	rows, err := tx.QueryContext(ctx, `select pg_try_advisory_xact_lock(hashtext(?))`, key)
	if err != nil {
		return err
	}
	if !rows.Next() {
		return errors.New("unexpected error: try to take advisory lock but no rows returned")
	}

	var result bool
	defer rows.Close()
	if err := rows.Scan(&result); err != nil {
		return err
	}
	if !result {
		return fmt.Errorf("data lock taken at the key %s: %w", key, ErrDataLockTaken)
	}
	return nil
}
