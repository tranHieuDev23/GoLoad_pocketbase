package http

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/tranHieuDev23/GoLoad_pocketbase/internal/dataaccess/database/migrations"
	"github.com/tranHieuDev23/GoLoad_pocketbase/internal/handlers/consumers"
	"github.com/tranHieuDev23/GoLoad_pocketbase/internal/logic"
)

type Server interface {
	Start() error
}

type server struct {
	pocketbaseApp *pocketbase.PocketBase
}

func NewServer() Server {
	pocketbaseApp := pocketbase.New()
	pocketbaseApp.OnAfterBootstrap().Add(func(e *core.BootstrapEvent) error {
		databaseMigrator := migrations.NewMigrator(e.App)
		if err := databaseMigrator.Migrate(); err != nil {
			return err
		}

		downloadTaskLogic := logic.NewDownloadTask()
		downloadTaskCreatedConsumer := consumers.NewDownloadTaskCreated(downloadTaskLogic)
		rootConsumer := consumers.NewRoot(downloadTaskCreatedConsumer)
		rootConsumer.RegisterHandlers(e.App)

		return nil
	})

	return &server{
		pocketbaseApp: pocketbaseApp,
	}
}

func (s *server) Start() error {
	return s.pocketbaseApp.Start()
}
