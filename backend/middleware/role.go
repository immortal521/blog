package middleware

import (
	"blog-server/entity"
	"blog-server/pkg/errx"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

const (
	ContextUserRoleKey = "user_role"
)

type RoleMiddleware struct {
	authService service.AuthService
}

func NewRoleMiddleware(authSerevice service.AuthService) *RoleMiddleware {
	return &RoleMiddleware{
		authService: authSerevice,
	}
}

func (r *RoleMiddleware) Handler(roles ...entity.UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		userUUID, ok := GetUserUUID(c)
		if !ok || userUUID == "" {
			return errx.New(errx.CodeUnauthorized, "missing user uuid", nil)
		}
		ok, err := r.authService.HasRole(c.Context(), userUUID)
		if err != nil || !ok {
			return errx.New(errx.CodeUnauthorized, "user not found", err)
		}

		return c.Next()
	}
}
