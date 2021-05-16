package dockertestx

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/pkg/errors"
)

type PostgresFactory struct{}

func (f *PostgresFactory) repository() string {
	return "postgres"
}

func (f *PostgresFactory) create(p *Pool, opt ContainerOption) (*state, error) {
	dbUser := "dockertestx"
	dbPassword := "passw0rd"
	dbName := "test"

	rOpt := &dockertest.RunOptions{
		Name:       opt.Name,
		Repository: f.repository(),
		Tag:        opt.Tag,
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", dbUser),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", dbPassword),
			fmt.Sprintf("POSTGRES_DB=%s", dbName),
		},
	}
	resource, err := p.Pool.RunWithOptions(rOpt)
	if err != nil {
		return nil, errors.WithMessage(err, "Could not start resource")
	}
	return &state{
		ContainerName: opt.Name,
		Repository:    f.repository(),
		Tag:           opt.Tag,
		Env:           rOpt.Env,
		DSN: fmt.Sprintf(
			"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
			dbUser,
			dbPassword,
			resource.GetPort("5432/tcp"),
			dbName,
		),
		r: resource,
	}, nil
}

func (f *PostgresFactory) ready(p *Pool, s *state) error {
	return p.Retry(func() error {
		db, err := sql.Open("postgres", s.DSN)
		if err != nil {
			return err
		}
		return db.Ping()
	})
}
