package handler

import (
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

type RssHandler interface {
	Subscript(c fiber.Ctx) error
}

type rssHandler struct {
	svc service.RssService
}

func (r *rssHandler) Subscript(c fiber.Ctx) error {
	data, err := r.svc.GenerateRSSFeedXML(c.Context())
	if err != nil {
		return err
	}

	c.Type("xml")
	return c.Send(data)
}

func NewRssHandler(svc service.RssService) RssHandler {
	return &rssHandler{svc: svc}
}

func RegisterRssRoute(r fiber.Router, handler RssHandler) {
	group := r.Group("/rss")
	group.Get("/", handler.Subscript)
}
