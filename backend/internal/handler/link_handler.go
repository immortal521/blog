package handler

import (
	"net/http"

	"blog-server/internal/request"
	"blog-server/internal/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"

	"github.com/labstack/echo/v4"
)

type ILinkHandler interface {
	GetLinks(c echo.Context) error
	ApplyForALinks(c echo.Context) error
}

type LinkHandler struct {
	svc      service.ILinkService
	validate validatorx.Validator
}

func NewLinkHandler(svc service.ILinkService, validate validatorx.Validator) ILinkHandler {
	return &LinkHandler{svc: svc, validate: validate}
}

func (h *LinkHandler) GetLinks(c echo.Context) error {
	links, err := h.svc.GetLinks(c.Request().Context())
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
	return c.JSON(http.StatusOK, result)
}

func (h *LinkHandler) ApplyForALinks(c echo.Context) error {
	req := new(request.CreateLinkReq)
	if err := c.Bind(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Failed to parse request body", err)
	}

	if (len(req.URL) == 0) || (len(req.Name) == 0) {
		return errs.New(errs.CodeInvalidParam, "URL or name is empty", nil)
	}

	err := h.svc.CreateLink(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.Success(""))
}
