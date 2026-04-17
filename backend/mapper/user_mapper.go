package mapper

import (
	"blog-server/ent"
	"blog-server/entity"
)

func ToUser(u *ent.User) *entity.User {
	if u == nil {
		return &entity.User{}
	}

	return &entity.User{
		ID:        u.ID,
		UUID:      u.UUID,
		Username:  u.Username,
		Email:     u.Email,
		Avatar:    u.Avatar,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
