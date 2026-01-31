package handler

import (
	"time"

	"blog-server/config"
	"blog-server/errs"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"
	"blog-server/validatorx"

	"github.com/gofiber/fiber/v2"
)

// IAuthHandler defines the interface for authentication HTTP handlers
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

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService service.IAuthService, validate validatorx.Validator) IAuthHandler {
	return &AuthHandler{svc: authService, validate: validate}
}

// Register handles user registration
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

// Logout handles user logout by clearing the refresh token cookie
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

// Login handles user login
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	req := new(request.LoginReq)

	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	result, err := h.svc.Login(c.UserContext(), req)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, result.RefreshToken)

	res := response.Success(result)
	return c.JSON(res)
}

// Refresh handles access token refresh using refresh token from cookie
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

// SendCaptcha handles sending captcha email for verification
func (h *AuthHandler) SendCaptcha(c *fiber.Ctx) error {
	req := new(request.GetCaptchaReq)

	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	if req.Type == "" {
		req.Type = string(service.Register)
	}

	err := h.svc.SendCaptchaMail(c.UserContext(), req.Email, service.CaptchaType(req.Type))
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessWithMsg("Captcha sent successfully", "Captcha sent successfully"))
}

// setRefreshTokenCookie sets the refresh token as an HTTP-only cookie
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
