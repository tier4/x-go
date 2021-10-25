package bunx

import (
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

func (c *Client) DB() *bun.DB {
	return c.db
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
