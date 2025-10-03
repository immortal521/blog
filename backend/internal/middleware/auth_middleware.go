package middleware

import (
	"blog-server/internal/entity"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	ContextUserIDKey = "user_id"
)

type AuthMiddleware struct {
	jwtService service.IJwtService
}

func NewAuthMiddleware(jwtService service.IJwtService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

func (a *AuthMiddleware) Handler(roles ...entity.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errs.Unauthorized("missing authorization header")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return errs.Unauthorized("invalid authorization header")
		}
		token := parts[1]
		if token == "" {
			return errs.Unauthorized("empty token")
		}

		claims, err := a.jwtService.ParseToken(token)
		if err != nil {
			return err // ParseToken 已返回 Unauthorized
		}

		// 验证 UUID 格式
		if _, err := uuid.Parse(claims.UserID); err != nil {
			return errs.Unauthorized("invalid user id in token")
		}

		c.Locals(ContextUserIDKey, claims.UserID)

		c.Locals(ContextUserIDKey)

		return c.Next()
	}
}

func GetUserUUID(c *fiber.Ctx) (string, bool) {
	value := c.Locals(ContextUserIDKey)
	uuid, ok := value.(string)
	return uuid, ok
}
