package handler

import (
	"blog-server/internal/service"

	"github.com/labstack/echo/v4"
)

type IRssHandler interface {
	Subscribe(c echo.Context) error
}

type rssHandler struct {
	svc service.IRssService
}

func (r *rssHandler) Subscribe(c echo.Context) error {
	data, err := r.svc.GenerateRSSFeedXML(c.Request().Context())
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, "application/xml")
	c.Response().Header().Set(echo.HeaderContentDisposition, `attachment; filename="rss.xml"`)

	// 写入数据
	_, err = c.Response().Write(data)
	return err
}

func NewRssHandler(svc service.IRssService) IRssHandler {
	return &rssHandler{svc: svc}
}
