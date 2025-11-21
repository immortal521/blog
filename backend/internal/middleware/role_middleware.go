package middleware

import (
	"blog-server/internal/entity"
	"blog-server/internal/service"
	"blog-server/pkg/errs"

	"github.com/labstack/echo/v4"
)

const (
	ContextUserRoleKey = "user_role"
)

type RoleMiddleware struct {
	authService service.IAuthService
}

func NewRoleMiddleware(authSerevice service.IAuthService) *RoleMiddleware {
	return &RoleMiddleware{
		authService: authSerevice,
	}
}

func (r *RoleMiddleware) Handler(roles ...entity.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从上下文获取用户 UUID
			userUUID, ok := GetUserUUID(c)
			if !ok || userUUID == "" {
				return errs.New(errs.CodeUnauthorized, "missing user uuid", nil)
			}

			// 调用服务层验证用户角色
			hasRole, err := r.authService.HasRole(c.Request().Context(), userUUID)
			if err != nil || !hasRole {
				return errs.New(errs.CodeUnauthorized, "user not found or role not allowed", err)
			}

			return next(c)
		}
	}
}
