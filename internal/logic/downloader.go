package logic

import (
	"context"
	"io"
	"net/http"
)

type Downloader interface {
	Download(ctx context.Context, writer io.Writer) error
}

type HTTPDownloader struct {
	url string
}

func NewHTTPDownloader(url string) Downloader {
	return &HTTPDownloader{
		url: url,
	}
}

func (h HTTPDownloader) Download(ctx context.Context, writer io.Writer) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, h.url, http.NoBody)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	_, err = io.Copy(writer, response.Body)
	return err
}
