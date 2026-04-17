package jwt

type Jwt interface {
	GenerateAccessToken(userID string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	GenerateAllTokens(userID string) (accessToken, refreshToken string, err error)

	Validate(token string) (bool, error)
	Parse(token string) (*Claims, error)
}
