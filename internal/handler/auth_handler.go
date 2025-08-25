package handler

import (
	"blog-server/internal/dto"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"fmt"
	"math/rand"
	"time"

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
	request := new(dto.CaptchaReq)

	if err := c.BodyParser(request); err != nil {
		return errs.BadRequest("invalid request")
	}

	if err := a.validate.Struct(request); err != nil {
		return errs.BadRequest("invalid")
	}

	if request.Type == "" {
		request.Type = "Register" // 如果用户未提供 Type，则默认为 Register
	}

	captcha := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	err := a.svc.SendCaptchaMail(c.UserContext(), request.Email, captcha, service.CaptchaType(request.Type))
	if err != nil {
		return err
	}

	// TODO: 在 Redis 中存储验证码和邮箱的映射关系，并设置过期时间

	return c.JSON(dto.Success(""))
}
