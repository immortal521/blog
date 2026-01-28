package handler

import (
	"blog-server/internal/request"
	"blog-server/internal/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"

	"github.com/gofiber/fiber/v2"
)

type ILinkHandler interface {
	GetLinks(c *fiber.Ctx) error
	ApplyForALinks(c *fiber.Ctx) error
	GetLinksOverview(c *fiber.Ctx) error
}

type LinkHandler struct {
	svc      service.ILinkService
	validate validatorx.Validator
}

func NewLinkHandler(svc service.ILinkService, validate validatorx.Validator) ILinkHandler {
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

func (h *LinkHandler) GetLinksOverview(c *fiber.Ctx) error {
	overview, err := h.svc.GetOverview(c.UserContext())
	if err != nil {
		return err
	}

	return c.JSON(response.Success(overview))
}

func (h *LinkHandler) ApplyForALinks(c *fiber.Ctx) error {
	req := new(request.CreateLinkReq)
	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Failed to parse request body", err)
	}

	if (len(req.URL) == 0) || (len(req.Name) == 0) {
		return errs.New(errs.CodeInvalidParam, "URL or name is empty", nil)
	}

	err := h.svc.CreateLink(c.UserContext(), req)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(""))
}
