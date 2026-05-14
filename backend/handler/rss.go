package handler

import (
	"blog-server/pkg/errx"
	"blog-server/request"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

type RssHandler interface {
	Subscript(c fiber.Ctx) error
	Complete(c fiber.Ctx) error
}

type rssHandler struct {
	svc service.RssService
}

func NewRssHandler(svc service.RssService) RssHandler {
	return &rssHandler{svc: svc}
}

func (r *rssHandler) Subscript(c fiber.Ctx) error {
	p := new(request.RssSubscriptReq)

	if err := c.Bind().Query(p); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	var data []byte
	var err error

	if p.Page <= 0 {
		data, err = r.svc.GenerateRSSFeedXML(c.Context())
	} else {
		defaultPageSize := 10
		data, err = r.svc.GeneratePagedFeedXML(c.Context(), p.Page, defaultPageSize)
	}

	if err != nil {
		return err
	}

	c.Type("xml")
	return c.Send(data)
}

func (r *rssHandler) Complete(c fiber.Ctx) error {
	data, err := r.svc.GenerateCompleteFeedXML(c.Context())
	if err != nil {
		return err
	}

	c.Type("xml")
	return c.Send(data)
}

// RegisterRssRoute 路由注册
func RegisterRssRoute(r fiber.Router, handler RssHandler) {
	group := r.Group("/rss")

	group.Get("/", handler.Subscript)

	group.Get("/complete", handler.Complete)
}
