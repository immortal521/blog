package middleware

import (
	"strings"

	"blog-server/internal/entity"
	"blog-server/internal/service"
	"blog-server/pkg/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	// ContextUserIDKey is the key used to store user ID in Fiber context
	ContextUserIDKey = "user_id"
)

// AuthMiddleware handles JWT authentication for protected routes
type AuthMiddleware struct {
	jwtService service.IJwtService
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware(jwtService service.IJwtService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
	}
}

// Handler returns a Fiber handler that validates JWT tokens
// Optional roles can be specified for role-based access control
func (a *AuthMiddleware) Handler(roles ...entity.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
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

		if _, err := uuid.Parse(claims.UserID); err != nil {
			return errs.New(errs.CodeUnauthorized, "invalid user id in token", nil)
		}

		c.Locals(ContextUserIDKey, claims.UserID)

		return c.Next()
	}
}

// GetUserUUID retrieves the user UUID from the Fiber context
func GetUserUUID(c *fiber.Ctx) (string, bool) {
	value := c.Locals(ContextUserIDKey)
	uuid, ok := value.(string)
	return uuid, ok
}
