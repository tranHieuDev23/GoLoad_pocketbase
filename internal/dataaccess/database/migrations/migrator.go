package migrations

import (
	"github.com/pocketbase/pocketbase/core"
)

type Migrator interface {
	Migrate() error
}

type migrator struct {
	pocketbaseApp core.App
}

func NewMigrator(
	pocketbaseApp core.App,
) Migrator {
	return &migrator{
		pocketbaseApp: pocketbaseApp,
	}
}

func (m migrator) Migrate() error {
	if err := Register0001Migration(m.pocketbaseApp); err != nil {
		return err
	}

	return nil
}
