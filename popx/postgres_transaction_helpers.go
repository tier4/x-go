package popx

import (
	"context"
	"errors"
	"fmt"

	"github.com/gobuffalo/pop/v5"
	"github.com/jmoiron/sqlx"
)

var (
	ErrDataLockTaken = errors.New("data lock taken")
)

type transactionContextKey int

const transactionKey transactionContextKey = 0

func WithTransaction(ctx context.Context, tx *pop.Connection) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

func (c *Client) Transaction(ctx context.Context, callback func(ctx context.Context, connection *pop.Connection) error) error {
	txCtx := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := txCtx.(*pop.Connection); ok {
			return callback(ctx, conn.WithContext(ctx))
		}
	}

	return c.c.WithContext(ctx).Transaction(func(tx *pop.Connection) error {
		return callback(WithTransaction(ctx, tx), tx)
	})
}

// TransactionWithTryAdvisoryLock is Transaction with pg_try_advisory_xact_lock
// if cannot take lock, returns error immediately
func (c *Client) TransactionWithTryAdvisoryLock(ctx context.Context, key string, callback func(ctx context.Context, connection *pop.Connection) error) error {
	txCtx := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := txCtx.(*pop.Connection); ok {
			return callback(ctx, conn)
		}
	}

	return c.c.Transaction(func(tx *pop.Connection) error {
		if err := tryTakeAdvisoryLock(tx, key); err != nil {
			return err
		}
		return callback(WithTransaction(ctx, tx), tx)
	})
}

func tryTakeAdvisoryLock(tx *pop.Connection, key string) error {
	rows, err := tx.Store.(sqlx.QueryerContext).
		QueryxContext(tx.Context(), `select pg_try_advisory_xact_lock(hashtext($1))`, key)
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

func (c *Client) GetConnection(ctx context.Context) *pop.Connection {
	txCtx := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := txCtx.(*pop.Connection); ok {
			return conn.WithContext(ctx)
		}
	}
	// WithContext() returns connection with store which incompatible with sqlx.QueryerContext
	return c.c
}

// GetSqlxQueryer returns sqlx.QueryerContext wrapped by pop
// This is useful for join query
func (c *Client) GetSqlxQueryer(ctx context.Context) sqlx.QueryerContext {
	return c.GetConnection(ctx).Store.(sqlx.QueryerContext)
}

// GetSqlxExecer returns sqlx.ExecerContext wrapped by pop
func (c *Client) GetSqlxExecer(ctx context.Context) sqlx.ExecerContext {
	return c.GetConnection(ctx).Store.(sqlx.ExecerContext)
}
