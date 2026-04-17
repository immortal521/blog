package jwt

import (
	"fmt"
	"time"

	"blog-server/config"
	"blog-server/pkg/errx"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type jwtx struct {
	cfg config.JWTConfig
}

func (j *jwtx) GenerateAccessToken(userUUID string) (string, error) {
	return j.generateToken(userUUID, j.cfg.AccessExpiration)
}

func (j *jwtx) GenerateRefreshToken(userUUID string) (string, error) {
	return j.generateToken(userUUID, j.cfg.RefreshExpiration)
}

func (j *jwtx) GenerateAllTokens(userUUID string) (string, string, error) {
	access, err := j.GenerateAccessToken(userUUID)
	if err != nil {
		return "", "", err
	}

	refresh, err := j.GenerateRefreshToken(userUUID)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (j *jwtx) Parse(tokenStr string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(token *jwtv5.Token) (any, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errx.New(
				errx.CodeInvalidSignature,
				fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]),
				nil,
			)
		}
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		return nil, errx.New(errx.CodeUnauthorized, "invalid or expired token", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errx.New(errx.CodeUnauthorized, "invalid token", nil)
	}

	return claims, nil
}

func (j *jwtx) Validate(tokenStr string) (bool, error) {
	_, err := jwtv5.Parse(tokenStr, func(token *jwtv5.Token) (any, error) {
		if _, ok := token.Method.(*jwtv5.SigningMethodHMAC); !ok {
			return nil, errx.New(
				errx.CodeInvalidSignature,
				fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]),
				nil,
			)
		}
		return []byte(j.cfg.Secret), nil
	})
	if err != nil {
		return false, errx.New(errx.CodeUnauthorized, "token invalid", err)
	}

	return true, nil
}

func (j *jwtx) generateToken(userUUID string, expires time.Duration) (string, error) {
	now := time.Now()

	claims := Claims{
		UserID: userUUID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(now.Add(expires)),
			IssuedAt:  jwtv5.NewNumericDate(now),
			NotBefore: jwtv5.NewNumericDate(now),
			Issuer:    j.cfg.Issuer,
			Subject:   "user token",
		},
	}

	token := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(j.cfg.Secret))
	if err != nil {
		return "", errx.New(errx.CodeInternalError, "failed to sign token", err)
	}

	return signed, nil
}

func New(cfg config.JWTConfig) Jwt {
	return &jwtx{cfg: cfg}
}
