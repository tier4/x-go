package dockertestx

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gobuffalo/pop/v5"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"

	"github.com/tier4/x-go/runtimex"
)

type PurgeFunc func() error

// NewPostgres is to create PostgreSQL container and to return its connection and close function
func NewPostgres(tag string) (*pop.Connection, PurgeFunc, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Could not connect to docker")
	}

	dbUser := "dockertest"
	dbPassword := "passw0rd"
	dbName := "test"

	resource, err := pool.Run(
		"postgres",
		tag,
		[]string{
			fmt.Sprintf("POSTGRES_USER=%s", dbUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
		})
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Could not start resource")
	}

	getDSN := func() string {
		return fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			resource.GetPort("5432/tcp"),
			dbName,
		)
	}

	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", getDSN())
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return nil, nil, errors.WithMessage(err, "Could not connect to docker")
	}

	conn, err := pop.NewConnection(&pop.ConnectionDetails{
		URL:             getDSN(),
		Pool:            runtimex.MaxParallelism() * 2,
		IdlePool:        runtimex.MaxParallelism(),
		ConnMaxLifetime: time.Duration(0),
	})
	if err != nil {
		return nil, nil, errors.WithMessage(err, "Could not connect resource")
	}
	if err := conn.Open(); err != nil {
		return nil, nil, errors.WithMessage(err, "Could not open connection")
	}

	var purgeFunc PurgeFunc = func() error {
		if conn.Store != nil {
			if err := conn.Close(); err != nil {
				return errors.WithMessage(err, "Could not close connection")
			}
		}
		if err := pool.Purge(resource); err != nil {
			return errors.WithMessage(err, "Could not purge resource")
		}
		return nil
	}

	return conn, purgeFunc, nil
}
