package bunx

import (
	"context"
	"embed"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type Migrator struct {
	migrator *migrate.Migrator
	logger   Logger
}

func NewMigrator(db *bun.DB, migrationFS embed.FS, logger Logger) (*Migrator, error) {
	migrations := migrate.NewMigrations()
	if err := migrations.Discover(migrationFS); err != nil {
		return nil, err
	}

	return &Migrator{
		migrator: migrate.NewMigrator(db, migrations),
		logger:   logger,
	}, nil
}

func (m *Migrator) Migrate(ctx context.Context) error {
	if err := m.migrator.Init(ctx); err != nil {
		return err
	}
	group, err := m.migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	m.logger.Info(fmt.Sprintf("migrated to %s", group))
	return nil
}

func (m *Migrator) Rollback(ctx context.Context) error {
	group, err := m.migrator.Rollback(ctx)
	if err != nil {
		return err
	}
	m.logger.Info(fmt.Sprintf("rolled back %s", group))
	return nil
}

func (m *Migrator) Reset(ctx context.Context) error {
	for {
		group, err := m.migrator.Rollback(ctx)
		if err != nil {
			return err
		}
		if group.IsZero() {
			return m.migrator.Reset(ctx)
		}
		m.logger.Info(fmt.Sprintf("rolled back %s", group))
	}
}

func (m *Migrator) Status(ctx context.Context) error {
	ms, err := m.migrator.MigrationsWithStatus(ctx)
	if err != nil {
		return err
	}
	m.logger.Info(fmt.Sprintf("migrations: %s", ms))
	m.logger.Info(fmt.Sprintf("un-applied migrations: %s", ms.Unapplied()))
	m.logger.Info(fmt.Sprintf("last migration group: %s", ms.LastGroup()))

	return nil
}
