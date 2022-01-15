package bunx_test

import (
	"context"
	"database/sql"
	"embed"
	"log"
	"testing"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	"github.com/tier4/x-go/bunx"
	"github.com/tier4/x-go/dockertestx"
)

var (
	client   *bunx.Client
	migrator *bunx.Migrator
)

//go:embed testdata/migrations/*.sql
var migrationFS embed.FS

// SA3000 os.Exit is not necessary from Go 1.15
func TestMain(m *testing.M) {
	p, err := dockertestx.New(dockertestx.PoolOption{})
	if err != nil {
		log.Fatalln(err)
	}
	defer func(p *dockertestx.Pool) {
		err := p.Purge()
		if err != nil {
			log.Fatalln(err)
		}
	}(p)

	dsn, err := p.NewResource(new(dockertestx.PostgresFactory), dockertestx.ContainerOption{
		Tag: "alpine",
	})
	if err != nil {
		log.Fatalln(err)
	}

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqlDB, pgdialect.New())

	client, err = bunx.NewClient(db)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(client *bunx.Client) {
		err := client.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(client)

	migrator, err = bunx.NewMigrator(db, migrationFS, bunx.NewNoopLogger())
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()
	_ = migrator.Reset(ctx)

	if err := migrator.Migrate(ctx); err != nil {
		log.Fatalln(err)
	}

	m.Run()
}

type User struct {
	ID    int64  `db:"id,pk"`
	Email string `db:"email"`
}
