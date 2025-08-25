package util

import (
	"blog-server/internal/config"
	"blog-server/pkg/errs"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateTokenByExpires(userId string, expires time.Duration) (string, error) {
	jwtCfg := config.Get().JWT
	claims := Claims{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expires)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),              // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),              // 生效时间
			Issuer:    jwtCfg.Issuer,                               // 签发人
			Subject:   "user token",                                // 主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtCfg.Secret)
}

func GenerateAccessToken(userId string) (string, error) {
	jwtCfg := config.Get().JWT
	token, err := GenerateTokenByExpires(userId, jwtCfg.AccessExpiration)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateRefreshToken(userId string) (string, error) {
	jwtCfg := config.Get().JWT
	token, err := GenerateTokenByExpires(userId, jwtCfg.RefreshExpiration)
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	jwtCfg := config.Get().JWT

	// 用 HMAC 校验签名
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtCfg.Secret), nil
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
