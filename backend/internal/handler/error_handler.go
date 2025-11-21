// Package handler provider web api handler
package handler

import (
	"net/http"

	"blog-server/internal/config"
	"blog-server/internal/response"
	"blog-server/pkg/errs"
	"blog-server/pkg/logger"

	"github.com/labstack/echo/v4"
)

func ErrorHandler(log logger.Logger, cfg *config.Config) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		appErr := errs.ToAppError(err)
		httpCode := errs.MapToHTTPStatus(appErr.Code)
		msg := appErr.Msg

		if httpCode == http.StatusInternalServerError {
			msg = "Internal Server Error"
		}

		if cfg.App.Environment == config.EnvDev {
			log.Error(appErr.FormatStack())
		} else {
			log.Error(msg, logger.Error(err))
		}

		// 返回 JSON 错误
		if !c.Response().Committed {
			if e := c.JSON(httpCode, response.Error(appErr.Code, msg)); e != nil {
				_ = c.String(httpCode, msg)
			}
		}
	}
}
