package handler

import (
	"blog-server/internal/dto/request"
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
	SendCaptcha(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Layout(c *fiber.Ctx) error
}

type AuthHandler struct {
	svc      service.IAuthService
	validate *validator.Validate
}

func NewAuthHandler(authService service.IAuthService) IAuthHandler {
	return &AuthHandler{svc: authService, validate: validatorx.Get()}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(request.RegisterReq)

	if err := c.BodyParser(req); err != nil {
		return errs.BadRequest("invalid request")
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.BadRequest("invalid")
	}

	err := h.svc.Register(c.UserContext(), req)
	if errors.Is(err, errs.ErrInvalidCaptcha) {
		return errs.BadRequest(err.Error())
	}
	if errors.Is(err, errs.ErrUserExists) {
		return errs.Conflict(err.Error())
	}
	if err != nil {
		return err
	}

	return nil
}

func (h *AuthHandler) Layout(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(request.LoginReq)

	if err := c.BodyParser(req); err != nil {
		return errs.BadRequest("invalid request")
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.BadRequest("invalid")
	}

	return nil
}

func (h *AuthHandler) SendCaptcha(c *fiber.Ctx) error {
	req := new(request.GetCaptchaReq)

	if err := c.BodyParser(req); err != nil {
		return errs.BadRequest("invalid request")
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.BadRequest("invalid")
	}

	if req.Type == "" {
		req.Type = string(service.Register)
	}

	err := h.svc.SendCaptchaMail(c.UserContext(), req.Email, service.CaptchaType(req.Type))
	if errors.Is(err, errs.ErrUserExists) {
		return errs.Conflict(err.Error())
	}
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessWithMsg("Captcha sent successfully", "Captcha sent successfully"))
}
