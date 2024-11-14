package http

import (
	"goload/internal/logic"
	"io"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type GetDownloadTaskFileHandler interface {
	Handle(c echo.Context) error
}

type getDownloadTaskFileHandler struct {
	pocketbaseApp     core.App
	downloadTaskLogic logic.DownloadTask
}

func NewGetDownloadTaskFileHandler(
	pocketbaseApp core.App,
	downloadTaskLogic logic.DownloadTask,
) GetDownloadTaskFileHandler {
	return &getDownloadTaskFileHandler{
		pocketbaseApp:     pocketbaseApp,
		downloadTaskLogic: downloadTaskLogic,
	}
}

func (h getDownloadTaskFileHandler) Handle(c echo.Context) error {
	authInfo := apis.RequestInfo(c)
	if authInfo.AuthRecord == nil {
		h.pocketbaseApp.Logger().Error("user is not logged in")
		c.Response().WriteHeader(http.StatusUnauthorized)
		return nil
	}

	accountID := authInfo.AuthRecord.Id
	downloadTaskID := c.PathParam("id")
	fileReader, err := h.downloadTaskLogic.GetDownloadTaskFile(c.Request().Context(), accountID, downloadTaskID)
	if err != nil {
		h.pocketbaseApp.Logger().With("err", err).Error("failed to get download task file")
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	_, err = io.Copy(c.Response().Writer, fileReader)
	if err != nil {
		h.pocketbaseApp.Logger().With("err", err).Error("failed to write download task file")
	}

	return nil
}
