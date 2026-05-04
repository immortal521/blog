package jwt

import (
	"blog-server/entity"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID   uint            `json:"id"`
	Role entity.UserRole `json:"role"`
	jwt.RegisteredClaims
}
