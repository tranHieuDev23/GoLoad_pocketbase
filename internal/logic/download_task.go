package logic

import (
	"context"
	"errors"
	"os"
)

type ExecuteDownloadTaskParams struct {
	ID           string
	DownloadType int
	URL          string
}

type DownloadTask interface {
	ExecuteDownloadTask(ctx context.Context, params ExecuteDownloadTaskParams) error
}

type downloadTask struct{}

func NewDownloadTask() DownloadTask {
	return &downloadTask{}
}

// ExecuteDownloadTask implements DownloadTask.
func (d *downloadTask) ExecuteDownloadTask(ctx context.Context, params ExecuteDownloadTaskParams) error {
	var downloader Downloader
	switch params.DownloadType {
	case 1: // HTTP
		downloader = NewHTTPDownloader(params.URL)

	default:
		return errors.New("unsupported download type")
	}

	file, err := os.Create(params.ID)
	if err != nil {
		return err
	}

	return downloader.Download(ctx, file)
}
