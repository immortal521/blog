package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRole int

const (
	RoleReader UserRole = iota + 1
	RoleAdmin
)

type User struct {
	*gorm.Model

	UUID     uuid.UUID `gorm:"type:binary(16);not null;unique"`
	Avatar   *string   `gorm:"type:varchar(255)"`
	Email    string    `gorm:"type:varchar(100);not null;unique"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`
	Role     UserRole  `gorm:"type:tinyint;default:1;not null"`
	Username string    `gorm:"type:varchar(50);not null"`

	Posts []Post `gorm:"foreignKey:UserID;references:ID"`

	// 关联关系
	//Comments []Comment `gorm:"foreignKey:UserID;references:ID" json:"comments,omitempty"`
}

func (User) TableName() string {
	return "users"
}
