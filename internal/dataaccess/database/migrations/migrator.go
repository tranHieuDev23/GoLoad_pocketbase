package migrations

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type Migrator interface {
	Migrate(ctx context.Context) error
}

type migrator struct {
	pocketbaseApp core.App
}

func NewMigrator(pocketbaseApp core.App) Migrator {
	return &migrator{
		pocketbaseApp: pocketbaseApp,
	}
}

func (m migrator) Migrate(_ context.Context) error {
	if err := register0001Migration(m.pocketbaseApp); err != nil {
		return err
	}

	return nil
}
