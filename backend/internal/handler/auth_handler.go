package handler

import (
	"net/http"
	"time"

	"blog-server/internal/config"
	"blog-server/internal/request"
	"blog-server/internal/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"

	"github.com/labstack/echo/v4"
)

type IAuthHandler interface {
	SendCaptcha(c echo.Context) error
	Register(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	Refresh(c echo.Context) error
}

type AuthHandler struct {
	svc      service.IAuthService
	validate validatorx.Validator
}

func NewAuthHandler(authService service.IAuthService, validate validatorx.Validator) IAuthHandler {
	return &AuthHandler{svc: authService, validate: validate}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(request.RegisterReq)

	if err := c.Bind(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	result, err := h.svc.Register(c.Request().Context(), req)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, result.RefreshToken)

	res := response.Success(result)
	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Logout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	return c.JSON(http.StatusOK, response.SuccessWithMsg("logout success", "logout success"))
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(request.LoginReq)

	if err := c.Bind(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	result, err := h.svc.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, result.RefreshToken)

	res := response.Success(result)
	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Refresh(c echo.Context) error {
	token, err := c.Cookie("refreshToken")
	if err != nil || token.Value == "" {
		return errs.New(errs.CodeUnauthorized, "Missing refresh token", nil)
	}
	newTokens, err := h.svc.RefreshAccessToken(c.Request().Context(), token.Value)
	if err != nil {
		return err
	}

	setRefreshTokenCookie(c, newTokens.RefreshToken)

	return c.JSON(http.StatusOK, response.Success(newTokens))
}

func (h *AuthHandler) SendCaptcha(c echo.Context) error {
	req := new(request.GetCaptchaReq)

	if err := c.Bind(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	// 默认类型
	if req.Type == "" {
		req.Type = string(service.Register)
	}

	// ===== 调用 Service 发送验证码 =====
	err := h.svc.SendCaptchaMail(c.Request().Context(), req.Email, service.CaptchaType(req.Type))
	if err != nil {
		return err
	}

	// ===== 成功响应 =====
	return c.JSON(http.StatusOK, response.SuccessWithMsg("Captcha sent successfully", "Captcha sent successfully"))
}

func setRefreshTokenCookie(c echo.Context, value string) {
	maxAge := config.Get().JWT.RefreshExpiration

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    value,
		Expires:  time.Now().Add(maxAge),
		HttpOnly: true,
		Secure:   true,
	}

	c.SetCookie(cookie)
}
