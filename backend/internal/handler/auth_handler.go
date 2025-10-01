package handler

import (
	"blog-server/internal/config"
	"blog-server/internal/dto/request"
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
	SendCaptcha(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
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

	result, err := h.svc.Register(c.UserContext(), req)
	switch {
	case errors.Is(err, errs.ErrInvalidCaptcha):
		return errs.BadRequest(err.Error())
	case errors.Is(err, errs.ErrUserExists):
		return errs.Conflict(err.Error())
	case errors.Is(err, errs.ErrTokenGeneration):
		// 用户已注册成功，但登录失败
		return c.JSON(response.SuccessWithMsg("Register success", "Register success"))
	case err != nil:
		return err
	}

	setRefreshTokenCookie(c, result.RefreshToken)

	res := response.Success(result)
	return c.JSON(res)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
	})

	return c.JSON(response.SuccessWithMsg("logout success", "logout success"))
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(request.LoginReq)

	if err := c.BodyParser(req); err != nil {
		return errs.BadRequest("invalid request")
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.BadRequest("parameter validation failed")
	}

	result, err := h.svc.Login(c.UserContext(), req)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, result.RefreshToken)

	res := response.Success(result)
	return c.JSON(res)
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

func setRefreshTokenCookie(c *fiber.Ctx, value string) {
	maxAge := config.Get().JWT.RefreshExpiration
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    value,
		Expires:  time.Now().Add(maxAge),
		HTTPOnly: true,
		Secure:   true,
	})
}
