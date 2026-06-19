package handler

import (
	"net/http"

	"blog-server/entity"
	"blog-server/pkg/errx"
	"blog-server/request"
	"blog-server/service"

	"github.com/labstack/echo/v5"
)

type RssHandler interface {
	Subscript(c *echo.Context) error
	Complete(c *echo.Context) error
}

type rssHandler struct {
	svc service.RssService
}

func NewRssHandler(svc service.RssService) RssHandler {
	return &rssHandler{svc: svc}
}

func (r *rssHandler) Subscript(c *echo.Context) error {
	p := new(request.RssSubscriptReq)

	if err := c.Bind(p); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	var data *entity.RSS
	var err error

	if p.Page <= 0 {
		data, err = r.svc.GenerateRSSFeed(c.Request().Context())
	} else {
		defaultPageSize := 10
		data, err = r.svc.GeneratePagedFeed(c.Request().Context(), p.Page, defaultPageSize)
	}

	if err != nil {
		return err
	}

	return c.XML(http.StatusOK, data)
}

func (r *rssHandler) Complete(c *echo.Context) error {
	data, err := r.svc.GenerateCompleteFeed(c.Request().Context())
	if err != nil {
		return err
	}

	return c.XML(http.StatusOK, data)
}

// RegisterRssRoute 路由注册
func RegisterRssRoute(r *echo.Group, handler RssHandler) {
	group := r.Group("/rss")

	group.GET("", handler.Subscript)

	group.GET("/complete", handler.Complete)
}
