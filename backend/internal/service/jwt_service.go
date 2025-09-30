package service

import (
	"blog-server/internal/config"
	"blog-server/pkg/errs"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	GenerateAllTokens(userUUID string) (accessToken, refreshToken string, err error)
	GenerateAccessToken(userUUID string) (string, error)
	GenerateRefreshToken(userUUID string) (string, error)
	ValidateToken(token string) (bool, error)
	ParseToken(tokenString string) (*Claims, error)
}

type JwtService struct {
	cfg config.JWTConfig
}

type Claims struct {
	UserID string `json:"userID"`
	jwt.RegisteredClaims
}

func NewJwtService() IJwtService {
	return &JwtService{
		cfg: config.Get().JWT,
	}
}

// GenerateToken implements IJwtService.
func (j *JwtService) generateTokenWithExpires(userUUID string, expires time.Duration) (string, error) {
	claims := Claims{
		UserID: userUUID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.Issuer,
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return "", errs.ErrTokenGeneration
	}
	return signed, nil
}

func (j *JwtService) GenerateAccessToken(userUUID string) (string, error) {
	return j.generateTokenWithExpires(userUUID, j.cfg.AccessExpiration)
}

func (j *JwtService) GenerateRefreshToken(userUUID string) (string, error) {
	return j.generateTokenWithExpires(userUUID, j.cfg.RefreshExpiration)
}

func (j *JwtService) GenerateAllTokens(userUUID string) (accessToken, refreshToken string, err error) {
	accessToken, err = j.GenerateAccessToken(userUUID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = j.GenerateRefreshToken(userUUID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken implements IJwtService.
func (j *JwtService) ValidateToken(tokenStr string) (bool, error) {
	_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.cfg.Secret), nil
	})

	if err != nil {
		return false, err
	}
	return true, nil
}

func (j *JwtService) ParseToken(tokenString string) (*Claims, error) {
	// 用 HMAC 校验签名
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		return nil, errs.Unauthorized("invalid or expired token")
	}

	// 提取 claims
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errs.Unauthorized("invalid token")
}
