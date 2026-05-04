package handler

import (
	"fmt"

	"blog-server/pkg/errx"
	"blog-server/pkg/validatorx"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

// LinkHandler defines the interface for link HTTP handlers
type LinkHandler interface {
	GetLinks(c fiber.Ctx) error
	ApplyForALinks(c fiber.Ctx) error
}

type linkHandler struct {
	svc      service.LinkService
	validate validatorx.Validator
}

func RegisterLinkRoutes(r fiber.Router, h LinkHandler) {
	group := r.Group("/links")
	group.Get("/", h.GetLinks)
	group.Post("/apply-link", h.ApplyForALinks)
}

// NewLinkHandler creates a new link handler instance
func NewLinkHandler(svc service.LinkService, validate validatorx.Validator) LinkHandler {
	return &linkHandler{svc: svc, validate: validate}
}

// GetLinks retrieves all enabled links
func (h *linkHandler) GetLinks(c fiber.Ctx) error {
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
func (h *linkHandler) ApplyForALinks(c fiber.Ctx) error {
	req := new(request.CreateLinkReq)
	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	if (len(req.URL) == 0) || (len(req.Name) == 0) {
		return errx.New(errx.CodeInvalidParam, fmt.Errorf("URL or name cannot be empty"))
	}

	input := &service.CreateLinkInput{
		Name:        req.Name,
		Description: &req.Description,
		Avatar:      &req.Avatar,
		URL:         req.URL,
	}
	if err := h.validate.Struct(input); err != nil {
		return err
	}

	err := h.svc.CreateLink(
		c.Context(),
		input,
	)
	if err != nil {
		return err
	}

	return c.JSON(response.Success(""))
}
