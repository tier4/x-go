package bunx

import (
	"context"

	"github.com/uptrace/bun"
)

type Client struct {
	db *bun.DB
}

func NewClient(db *bun.DB) (*Client, error) {
	return &Client{
		db: db,
	}, nil
}

func (c *Client) Ping() error {
	if err := c.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() error {
	return c.db.Close()
}

type transactionContextKey int

const transactionKey transactionContextKey = 0

type DB interface {
	bun.IDB
	bun.IConn
}

func (c *Client) DB(ctx context.Context) DB {
	txCtx := ctx.Value(transactionKey)
	if c != nil {
		if conn, ok := txCtx.(bun.Tx); ok {
			return conn
		}
	}
	return c.db
}

func WithTransaction(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, transactionKey, tx)
}

func (c *Client) Transaction(ctx context.Context, callback func(ctx context.Context, tx bun.Tx) error) error {
	txCtx := ctx.Value(transactionKey)
	if tx, ok := txCtx.(bun.Tx); ok {
		return callback(ctx, tx)
	}

	return c.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		return callback(WithTransaction(ctx, tx), tx)
	})
}
