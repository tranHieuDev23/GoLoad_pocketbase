package events

import (
	"context"
	"goload/internal/logic"

	"github.com/pocketbase/pocketbase/models"
)

type DownloadTaskCreatedHandler interface {
	Handle(ctx context.Context, downloadTask models.Model) error
}

type downloadTaskCreatedHandler struct {
	downloadTaskLogic logic.DownloadTask
}

func NewDownloadTaskCreatedHandler(
	downloadTaskLogic logic.DownloadTask,
) DownloadTaskCreatedHandler {
	return &downloadTaskCreatedHandler{
		downloadTaskLogic: downloadTaskLogic,
	}
}

func (d downloadTaskCreatedHandler) Handle(ctx context.Context, downloadTask models.Model) error {
	return d.downloadTaskLogic.ExecuteDownloadTask(ctx, downloadTask.GetId())
}
