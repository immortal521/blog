package middleware

import (
	"blog-server/internal/entity"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"errors"

	"github.com/gofiber/fiber/v2"
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

func (r *RoleMiddleware) Handler(roles ...entity.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userUUID, ok := GetUserUUID(c)
		if !ok || userUUID == "" {
			return errs.New(errs.CodeUnauthorized, "missing user uuid", nil)
		}
		ok, err := r.authService.HasRole(c.UserContext(), userUUID)
		if errors.Is(err, errs.ErrUserNotFound) || !ok {
			return errs.New(errs.CodeUnauthorized, "user not found", err)
		}

		return c.Next()
	}

}
