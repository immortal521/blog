package handler

import (
	"fmt"
	"time"

	"blog-server/config"
	"blog-server/pkg/errx"
	"blog-server/pkg/validatorx"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

// AuthHandler defines the interface for authentication HTTP handlers
type AuthHandler interface {
	SendCaptcha(c fiber.Ctx) error
	Register(c fiber.Ctx) error
	Login(c fiber.Ctx) error
	Logout(c fiber.Ctx) error
	Refresh(c fiber.Ctx) error
}

type authHandler struct {
	svc      service.AuthService
	validate validatorx.Validator
}

func RegisterAuthRoutes(r fiber.Router, h AuthHandler) {
	group := r.Group("/auth")
	group.Post("/captcha", h.SendCaptcha)
	group.Post("/register", h.Register)
	group.Post("/login", h.Login)
	group.Post("/logout", h.Logout)
	group.Post("/refresh", h.Refresh)
}

// NewAuthHandler creates a new auth handler instance
func NewAuthHandler(authService service.AuthService, validate validatorx.Validator) AuthHandler {
	return &authHandler{svc: authService, validate: validate}
}

// Register handles user registration
func (h *authHandler) Register(c fiber.Ctx) error {
	req := new(request.RegisterReq)

	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errx.New(errx.CodeValidationFailed, err)
	}

	accessToken, refreshToken, err := h.svc.Register(
		c.Context(),
		req.Email,
		req.Password,
		req.Captcha,
	)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, refreshToken)

	res := response.Success(&response.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	return c.JSON(res)
}

// Logout handles user logout by clearing the refresh token cookie
func (h *authHandler) Logout(c fiber.Ctx) error {
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
func (h *authHandler) Login(c fiber.Ctx) error {
	req := new(request.LoginReq)

	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errx.New(errx.CodeValidationFailed, err)
	}

	accessToken, refreshToken, err := h.svc.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, refreshToken)

	res := response.Success(&response.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
	return c.JSON(res)
}

// Refresh handles access token refresh using refresh token from cookie
func (h *authHandler) Refresh(c fiber.Ctx) error {
	token := c.Cookies("refreshToken")
	if token == "" {
		return errx.New(errx.CodeUnauthorized, fmt.Errorf("Missing refresh token"))
	}
	accessToken, refreshToken, err := h.svc.RefreshAccessToken(c.Context(), token)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, refreshToken)

	return c.JSON(response.Success(&response.RefreshRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}))
}

// SendCaptcha handles sending captcha email for verification
func (h *authHandler) SendCaptcha(c fiber.Ctx) error {
	req := new(request.GetCaptchaReq)

	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errx.New(errx.CodeValidationFailed, err)
	}

	if req.Type == "" {
		req.Type = string(service.Register)
	}

	err := h.svc.SendCaptchaMail(c.Context(), req.Email, service.CaptchaType(req.Type))
	if err != nil {
		return err
	}

	return c.JSON(response.SuccessWithMsg("Captcha sent successfully", "Captcha sent successfully"))
}

// setRefreshTokenCookie sets the refresh token as an HTTP-only cookie
func setRefreshTokenCookie(c fiber.Ctx, value string) {
	maxAge := config.Get().JWT.RefreshExpiration
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    value,
		Expires:  time.Now().Add(maxAge),
		HTTPOnly: true,
		Secure:   true,
	})
}
