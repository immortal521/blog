package response

import "blog-server/internal/entity"

type LoginRes struct {
	AccessToken  string          `json:"accessToken"`
	RefreshToken string          `json:"-"`
	UUID         string          `json:"uuid"`
	Avatar       *string         `json:"avatar"`
	Username     string          `json:"username"`
	Role         entity.UserRole `json:"role"`
}
