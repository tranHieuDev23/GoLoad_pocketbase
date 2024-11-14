package http

import (
	"context"
	"goload/internal/dataaccess/database/migrations"
	"goload/internal/handler/events"
	"goload/internal/logic"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type Server interface {
	Start() error
}

type server struct {
	pocketbaseApp *pocketbase.PocketBase
}

func NewServer() Server {
	pocketbaseApp := pocketbase.New()
	databaseMigrator := migrations.NewMigrator(pocketbaseApp)
	downloadTaskLogic := logic.NewDownloadTask(pocketbaseApp)
	downloadTaskCreatedEventHandler := events.NewDownloadTaskCreatedHandler(downloadTaskLogic)
	getDownloadTaskFileHandler := NewGetDownloadTaskFileHandler(pocketbaseApp, downloadTaskLogic)

	pocketbaseApp.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		if err := databaseMigrator.Migrate(context.Background()); err != nil {
			return err
		}

		pocketbaseApp.OnRecordAfterCreateRequest("download_tasks").Add(func(e *core.RecordCreateEvent) error {
			return downloadTaskCreatedEventHandler.Handle(context.Background(), e.Record)
		})

		e.Router.GET("/api/download-tasks/:id/file", getDownloadTaskFileHandler.Handle)
		return nil
	})

	return &server{
		pocketbaseApp: pocketbaseApp,
	}
}

func (s server) Start() error {
	return s.pocketbaseApp.Start()
}
