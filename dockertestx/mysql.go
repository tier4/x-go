package dockertestx

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ory/dockertest/v3"
)

type MysqlFactory struct{}

func (f *MysqlFactory) repository() string {
	return "mysql"
}

func (f *MysqlFactory) create(p *Pool, opt ContainerOption) (*state, error) {
	dbUser := "root"
	dbPassword := "passw0rd"
	dbName := "test"

	rOpt := &dockertest.RunOptions{
		Name:       opt.Name,
		Repository: f.repository(),
		Tag:        opt.Tag,
		Env: []string{
			fmt.Sprintf("MYSQL_ROOT_PASSWORD=%s", dbPassword),
			fmt.Sprintf("MYSQL_DATABASE=%s", dbName),
		},
	}
	resource, err := p.Pool.RunWithOptions(rOpt)
	if err != nil {
		return nil, fmt.Errorf("could not start resource: %w", err)
	}
	return &state{
		ContainerName: opt.Name,
		Repository:    f.repository(),
		Tag:           opt.Tag,
		Env:           rOpt.Env,
		DSN: fmt.Sprintf(
			"%s:%s@tcp(localhost:%s)/%s",
			dbUser,
			dbPassword,
			resource.GetPort("3306/tcp"),
			dbName,
		),
		r: resource,
	}, nil
}

func (f *MysqlFactory) ready(p *Pool, s *state) error {
	return p.Retry(func() error {
		db, err := sql.Open("mysql", s.DSN)
		if err != nil {
			return err
		}
		return db.Ping()
	})
}
