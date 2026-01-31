package entity

import (
	"encoding/json"

	"blog-server/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRole represents the role of a user in the system
type UserRole int

const (
	// RoleReader represents a regular user with read-only access
	RoleReader UserRole = iota + 1
	// RoleAdmin represents an administrator with full access
	RoleAdmin
)

// MarshalJSON converts UserRole to string for JSON serialization
func (r UserRole) MarshalJSON() ([]byte, error) {
	var s string
	switch r {
	default:
		s = "unknown"
	case RoleReader:
		s = "reader"
	case RoleAdmin:
		s = "admin"
	}
	return json.Marshal(s)
}

// User represents a user entity in the system
type User struct {
	*gorm.Model

	UUID     uuid.UUID `gorm:"type:uuid;not null;unique"`           // Unique identifier for the user
	Avatar   *string   `gorm:"type:varchar(255)"`                   // Avatar image URL (optional)
	Email    string    `gorm:"type:varchar(100);not null;unique"`   // User email (unique)
	Password string    `gorm:"type:varchar(255);not null" json:"-"` // Hashed password (excluded from JSON)
	Role     UserRole  `gorm:"type:smallint;default:1;not null"`    // User role
	Username string    `gorm:"type:varchar(50);not null"`           // Display username

	Posts []Post `gorm:"foreignKey:UserID;references:ID"` // Posts created by this user

	// Comments []Comment `gorm:"foreignKey:UserID;references:ID" json:"comments,omitempty"` // Comments by this user (reserved)
}

func GenerateUsername() string {
	return "user_" + utils.RandomString(6, "abcdefghijklmnopqrstuvwxyz0123456789")
}

// TableName returns the table name for User model
func (User) TableName() string {
	return "users"
}
