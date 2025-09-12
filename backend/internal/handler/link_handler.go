package handler

import (
	"blog-server/internal/dto"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"github.com/go-playground/validator/v10"

	"github.com/gofiber/fiber/v2"
)

type ILinkHandler interface {
	GetLinks(c *fiber.Ctx) error
	ApplyForALinks(c *fiber.Ctx) error
}
type LinkHandler struct {
	svc      service.ILinkService
	validate *validator.Validate
}

func NewLinkHandler(svc service.ILinkService) ILinkHandler {
	return &LinkHandler{svc: svc, validate: validator.New()}
}

func (l *LinkHandler) GetLinks(c *fiber.Ctx) error {
	links, err := l.svc.GetLinks(c.UserContext())
	if err != nil {
		return err
	}
	linkDTOs := make([]dto.LinkRes, len(links))
	for i, link := range links {
		linkDTOs[i] = dto.LinkRes{
			ID:          link.ID,
			Name:        link.Name,
			Url:         link.URL,
			Description: link.Description,
			Avatar:      link.Avatar,
			SortOrder:   link.SortOrder,
		}
	}
	result := dto.Success(linkDTOs)
	return c.JSON(result)
}

func (l *LinkHandler) ApplyForALinks(c *fiber.Ctx) error {
	request := new(dto.LinkCreateReq)
	if err := c.BodyParser(request); err != nil {
		return err
	}

	if (len(request.Url) == 0) || (len(request.Name) == 0) {
		return errs.BadRequest("url or name is empty")
	}

	err := l.svc.CreateLink(c.UserContext(), request)
	if err != nil {
		return err
	}

	return c.JSON(dto.Success(""))
}
