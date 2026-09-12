package handler

import (
	"blog-server/pkg/errx"
	"blog-server/pkg/validatorx"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"

	"github.com/labstack/echo/v5"
)

type SetupHandler interface {
	Status(c *echo.Context) error
	Initialize(c *echo.Context) error
}

type setupHandler struct {
	svc      service.SetupService
	validate validatorx.Validator
}

func NewSetupHandler(svc service.SetupService, validate validatorx.Validator) SetupHandler {
	return &setupHandler{
		svc, validate,
	}
}

func (h *setupHandler) Initialize(c *echo.Context) error {
	req := new(request.InitializeReq)
	if err := c.Bind(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	input := &service.InitializeInput{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		SiteName:    req.SiteName,
		Logo:        req.Logo,
		Greeting:    req.Greeting,
		Description: req.Description,
	}

	err := h.svc.Initialize(c.Request().Context(), input)
	if err != nil {
		return err
	}

	return response.OK(c, struct {
		Initialized bool `json:"initialized"`
	}{
		Initialized: true,
	})
}

func (h *setupHandler) Status(c *echo.Context) error {
	initialized, err := h.svc.Status(c.Request().Context())
	if err != nil {
		return err
	}

	return response.OK(c, struct {
		Initialized bool `json:"initialized"`
	}{
		Initialized: initialized,
	})
}

func RegisterSetupRouter(r *echo.Group, h SetupHandler) {
	group := r.Group("/setup")
	group.GET("", h.Status)
	group.POST("", h.Initialize)
}
