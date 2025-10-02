package handler

import (
	"blog-server/internal/dto/request"
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"errors"

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

func NewLinkHandler(svc service.ILinkService, validate *validator.Validate) ILinkHandler {
	return &LinkHandler{svc: svc, validate: validate}
}

func (h *LinkHandler) GetLinks(c *fiber.Ctx) error {
	links, err := h.svc.GetLinks(c.UserContext())
	if err != nil {
		return err
	}
	linkDTOs := make([]response.LinkResponse, len(links))
	for i, link := range links {
		linkDTOs[i] = response.LinkResponse{
			ID:          link.ID,
			Name:        link.Name,
			URL:         link.URL,
			Description: link.Description,
			Avatar:      link.Avatar,
			SortOrder:   link.SortOrder,
		}
	}
	result := response.Success(linkDTOs)
	return c.JSON(result)
}

func (h *LinkHandler) ApplyForALinks(c *fiber.Ctx) error {
	req := new(request.CreateLinkReq)
	if err := c.BodyParser(req); err != nil {
		return err
	}

	if (len(req.URL) == 0) || (len(req.Name) == 0) {
		return errs.BadRequest("url or name is empty")
	}

	err := h.svc.CreateLink(c.UserContext(), req)
	if err == nil {
		return c.JSON(response.Success(""))
	}
	if errors.Is(err, errs.ErrDuplicateURL) {
		return errs.Conflict(err.Error())
	}
	return err
}
