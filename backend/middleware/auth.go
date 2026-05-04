package middleware

import (
	"fmt"
	"strings"

	"blog-server/authz"
	"blog-server/config"
	"blog-server/contextx"
	"blog-server/pkg/errx"
	"blog-server/pkg/jwt"

	"github.com/gofiber/fiber/v3"
)

// AuthMiddleware handles JWT authentication for protected routes
type AuthMiddleware struct {
	jwt jwt.Jwt
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(cfg *config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt.New(cfg.JWT),
	}
}

// Handler returns a Fiber handler that validates JWT tokens
// Optional roles can be specified for role-based access control
func (m *AuthMiddleware) Handler() fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errx.New(errx.CodeUnauthorized, fmt.Errorf("missing authorization header"))
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return errx.New(errx.CodeUnauthorized, fmt.Errorf("invalid authorization header"))
		}

		tokenStr := parts[1]
		if tokenStr == "" {
			return errx.New(errx.CodeUnauthorized, fmt.Errorf("empty token"))
		}

		claims, err := m.jwt.Parse(tokenStr)
		if err != nil {
			return err
		}

		user := contextx.User{
			ID:   claims.ID,
			Role: authz.FromEntityRole(claims.Role),
		}

		c.SetContext(contextx.SetUser(c.Context(), user))

		return c.Next()
	}
}
