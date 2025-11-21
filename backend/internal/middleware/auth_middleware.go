package middleware

import (
	"strings"

	"blog-server/internal/entity"
	"blog-server/internal/service"
	"blog-server/pkg/errs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
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

// Handler 返回 Echo 中间件
func (a *AuthMiddleware) Handler(roles ...entity.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return errs.New(errs.CodeUnauthorized, "missing authorization header", nil)
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				return errs.New(errs.CodeUnauthorized, "invalid authorization header", nil)
			}

			token := parts[1]
			if token == "" {
				return errs.New(errs.CodeUnauthorized, "empty token", nil)
			}

			claims, err := a.jwtService.ParseToken(token)
			if err != nil {
				return err
			}

			// 验证 UUID 格式
			if _, err := uuid.Parse(claims.UserID); err != nil {
				return errs.New(errs.CodeUnauthorized, "invalid user id in token", nil)
			}

			// 保存用户 ID 到 Echo Context
			c.Set(ContextUserIDKey, claims.UserID)

			return next(c)
		}
	}
}

func GetUserUUID(c echo.Context) (string, bool) {
	value := c.Get(ContextUserIDKey)
	userID, ok := value.(string)
	return userID, ok
}
