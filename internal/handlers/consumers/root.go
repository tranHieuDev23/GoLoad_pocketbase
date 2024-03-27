package consumers

import (
	"context"

	"github.com/pocketbase/pocketbase/core"
)

type Root interface {
	RegisterHandlers(app core.App)
}

type root struct {
	downloadTaskCreated DownloadTaskCreated
}

func NewRoot(
	downloadTaskCreated DownloadTaskCreated,
) Root {
	return &root{
		downloadTaskCreated: downloadTaskCreated,
	}
}

// RegisterHandlers implements Root.
func (r root) RegisterHandlers(app core.App) {
	app.OnRecordAfterCreateRequest("download_tasks").Add(func(e *core.RecordCreateEvent) error {
		return r.downloadTaskCreated.Handle(context.Background(), e)
	})
}
