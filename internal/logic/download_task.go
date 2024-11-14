package logic

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

const (
	downloadTypeHTTP = 1
)

type DownloadTask interface {
	ExecuteDownloadTask(ctx context.Context, id string) error
	GetDownloadTaskFile(ctx context.Context, accountID, downloadTaskID string) (io.Reader, error)
}

type downloadTask struct {
	pocketbaseApp core.App
}

func NewDownloadTask(pocketbaseApp core.App) DownloadTask {
	return &downloadTask{
		pocketbaseApp: pocketbaseApp,
	}
}

func (d downloadTask) ExecuteDownloadTask(ctx context.Context, id string) error {
	logger := d.pocketbaseApp.Logger().With(slog.String("id", id))
	downloadTask, err := d.pocketbaseApp.Dao().FindRecordById("download_tasks", id)
	if err != nil {
		logger.With("err", err).Error("failed to get download task")
		return err
	}

	downloadType := downloadTask.GetInt("download_type")
	if downloadType != downloadTypeHTTP {
		logger.With(slog.Int("download_type", downloadType)).Error("unsupported download type")
		return errors.New("unsupported download type")
	}

	url := downloadTask.GetString("url")
	file, err := filesystem.NewFileFromUrl(ctx, url)
	if err != nil {
		logger.With("err", err).Error("failed to download file")
		return err
	}

	fs, err := d.pocketbaseApp.NewFilesystem()
	if err != nil {
		logger.With("err", err).Error("failed to get new pocketbase filesystem")
		return err
	}

	defer fs.Close()

	err = fs.UploadFile(file, id)
	if err != nil {
		logger.With("err", err).Error("failed to upload file")
		return err
	}

	logger.Info("download task completed")
	return nil
}

func (d downloadTask) GetDownloadTaskFile(_ context.Context, accountID, downloadTaskID string) (io.Reader, error) {
	logger := d.pocketbaseApp.Logger().
		With(slog.String("account_id", accountID)).
		With(slog.String("download_task_id", downloadTaskID))

	downloadTask, err := d.pocketbaseApp.Dao().FindRecordById("download_tasks", downloadTaskID)
	if err != nil {
		logger.With("err", err).Error("failed to get download task")
		return nil, err
	}

	if accountID != downloadTask.GetStringSlice("of_account_id")[0] {
		logger.With("err", err).Error("account does not have access to this file")
		return nil, errors.New("account does not have access to this file")
	}

	fs, err := d.pocketbaseApp.NewFilesystem()
	if err != nil {
		logger.With("err", err).Error("failed to get new pocketbase filesystem")
		return nil, err
	}

	defer fs.Close()

	reader, err := fs.GetFile(downloadTaskID)
	if err != nil {
		logger.With("err", err).Error("failed to get file reader")
		return nil, err
	}

	return reader, nil
}
