package handler

import (
	"blog-server/internal/service"

	"github.com/gofiber/fiber/v2"
)

type IRssHandler interface {
	Subscribe(c *fiber.Ctx) error
}

type rssHandler struct {
	svc service.IRssService
}

func (r *rssHandler) Subscribe(c *fiber.Ctx) error {
	data, err := r.svc.GenerateRSSFeedXML(c.UserContext())
	if err != nil {
		return err
	}

	c.Type("xml")
	c.Append("Content-Disposition", "attachment; filename=rss.xml")
	return c.Send(data)
}

func NewRssHandler(svc service.IRssService) IRssHandler {
	return &rssHandler{svc: svc}
}
