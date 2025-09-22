package handler

import (
	"blog-server/internal/dto/request"
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
	SendCaptcha(c *fiber.Ctx) error
}

type AuthHandler struct {
	svc      service.IAuthService
	validate *validator.Validate
}

func NewAuthHandler(authService service.IAuthService) IAuthHandler {
	return &AuthHandler{svc: authService, validate: validator.New()}
}

func (a *AuthHandler) SendCaptcha(c *fiber.Ctx) error {
	request := new(request.GetCaptchaReq)

	if err := c.BodyParser(request); err != nil {
		return errs.BadRequest("invalid request")
	}

	if err := a.validate.Struct(request); err != nil {
		return errs.BadRequest("invalid")
	}

	if request.Type == "" {
		request.Type = string(service.Register)
	}

	err := a.svc.SendCaptchaMail(c.UserContext(), request.Email, service.CaptchaType(request.Type))
	if err != nil {
		return err
	}

	return c.JSON(response.Success(""))
}
