package handler

import (
	"blog-server/pkg/errx"
	"blog-server/pkg/validatorx"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

// ILinkHandler defines the interface for link HTTP handlers
type ILinkHandler interface {
	GetLinks(c fiber.Ctx) error
	ApplyForALinks(c fiber.Ctx) error
}

type LinkHandler struct {
	svc      service.LinkService
	validate validatorx.Validator
}

func RegisterLinkRoutes(r fiber.Router, h LinkHandler) {
	group := r.Group("/links")
	group.Get("/", h.GetLinks)
	group.Post("/apply-link", h.ApplyForALinks)
}

// NewLinkHandler creates a new link handler instance
func NewLinkHandler(svc service.LinkService, validate validatorx.Validator) ILinkHandler {
	return &LinkHandler{svc: svc, validate: validate}
}

// GetLinks retrieves all enabled links
func (h *LinkHandler) GetLinks(c fiber.Ctx) error {
	links, err := h.svc.GetLinks(c.Context())
	if err != nil {
		return err
	}

	linkDTOs := make([]response.LinkResponse, len(links))
	for i, link := range links {
		linkDTOs[i] = response.LinkResponse{
			ID:          link.ID,
			Name:        link.Name,
			URL:         link.URL,
			Description: *link.Description,
			Avatar:      *link.Avatar,
		}
	}

	result := response.Success(linkDTOs)
	return c.JSON(result)
}

// ApplyForALinks creates a new link application
func (h *LinkHandler) ApplyForALinks(c fiber.Ctx) error {
	req := new(request.CreateLinkReq)
	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, "Failed to parse request body", err)
	}

	if (len(req.URL) == 0) || (len(req.Name) == 0) {
		return errx.New(errx.CodeInvalidParam, "URL or name is empty", nil)
	}

	err := h.svc.CreateLink(
		c.Context(),
		req.Name,
		req.Description,
		req.Avatar,
		req.URL,
	)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(""))
}
