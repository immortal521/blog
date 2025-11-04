package handler

import (
	"blog-server/internal/config"
	"blog-server/internal/dto/request"
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"
	"time"

	"github.com/gofiber/fiber/v2"
)

type IAuthHandler interface {
	SendCaptcha(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Login(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
}

type AuthHandler struct {
	svc      service.IAuthService
	validate validatorx.Validator
}

func NewAuthHandler(authService service.IAuthService, validate validatorx.Validator) IAuthHandler {
	return &AuthHandler{svc: authService, validate: validate}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	req := new(request.RegisterReq)

	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	result, err := h.svc.Register(c.UserContext(), req)
	if err != nil {
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

func (h *AuthHandler) Refresh(c *fiber.Ctx) error {
	token := c.Cookies("refreshToken")
	if token == "" {
		return errs.New(errs.CodeUnauthorized, "Missing refresh token", nil)
	}
	newTokens, err := h.svc.RefreshAccessToken(c.UserContext(), token)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, newTokens.RefreshToken)

	return c.JSON(response.Success(newTokens))
}

func (h *AuthHandler) SendCaptcha(c *fiber.Ctx) error {
	req := new(request.GetCaptchaReq)

	// ===== 请求体解析 =====
	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	// ===== 参数校验 =====
	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	// 默认类型
	if req.Type == "" {
		req.Type = string(service.Register)
	}

	// ===== 调用 Service 发送验证码 =====
	err := h.svc.SendCaptchaMail(c.UserContext(), req.Email, service.CaptchaType(req.Type))
	if err != nil {
		return err
	}

	// ===== 成功响应 =====
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
