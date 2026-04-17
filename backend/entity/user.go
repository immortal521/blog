package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleReader UserRole = "reader"
	UserRoleAdmin  UserRole = "admin"
)

func (UserRole) Values() []string {
	return []string{
		string(UserRoleReader),
		string(UserRoleAdmin),
	}
}

type User struct {
	ID   uint
	UUID uuid.UUID

	Username string
	Email    string
	Avatar   *string

	Role UserRole

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAuth struct {
	UUID     uuid.UUID
	Password string
	Role     UserRole
}
