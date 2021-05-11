package popx

import (
	"context"
	"io"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

type Client struct {
	c  *pop.Connection
	mb *pop.MigrationBox
}

func New(conn *pop.Connection, box *pop.MigrationBox) (*Client, error) {
	return &Client{
		c:  conn,
		mb: box,
	}, nil
}

// MigrationStatus returns migration status
func (c *Client) MigrationStatus(_ context.Context, w io.Writer) error {
	return c.mb.Status(w)
}

// MigrateDown rollbacks given steps
func (c *Client) MigrateDown(_ context.Context, steps int) error {
	return c.mb.Down(steps)
}

// MigrateUp migrates all of un-executed
func (c *Client) MigrateUp(_ context.Context) error {
	return c.mb.Up()
}

func (c *Client) Close(ctx context.Context) error {
	return errors.WithStack(c.GetConnection(ctx).Close())
}

func (c *Client) Ping() error {
	type pinger interface {
		Ping() error
	}
	// This can not be contextualized because of some gobuffalo/pop limitations.
	return errors.WithStack(c.c.Store.(pinger).Ping())
}
