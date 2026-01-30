// Package entity defines data models for the blog system
package entity

import (
	"time"

	"gorm.io/gorm"
)

// PostStatus represents the status of a post
type PostStatus int

const (
	// PostDraft indicates the post is in draft state
	PostDraft PostStatus = iota + 1
	// PostPublished indicates the post is published and visible
	PostPublished
	// PostArchived indicates the post is archived
	PostArchived
)

// Post represents a blog post entity
type Post struct {
	*gorm.Model

	Title           string     `gorm:"type:varchar(255);not null"`           // Post title
	Summary         *string    `gorm:"type:varchar(500)"`                    // Post summary (optional)
	Content         string     `gorm:"type:text;not null"`                    // Post content in markdown or HTML
	Cover           *string    `gorm:"type:varchar(255)"`                     // Cover image URL (optional)
	ReadTimeMinutes uint       `gorm:"not null"`                              // Estimated reading time in minutes
	ViewCount       uint       `gorm:"not null"`                              // Total view count
	Status          PostStatus `gorm:"type:smallint;default:1;not null"`      // Post status

	UserID uint `gorm:"not null;index"` // Foreign key to User
	User   User `gorm:"foreignKey:UserID;references:ID"` // Author of the post

	PublishedAt *time.Time `gorm:"type:timestamp(6)"` // Publication timestamp

	Categories []PostCategory `gorm:"many2many:post_category_relations;"` // Many-to-many relationship with categories
	Tags       []PostTag      `gorm:"many2many:post_tag_relations;"`      // Many-to-many relationship with tags
}

// TableName returns the table name for Post model
func (Post) TableName() string {
	return "posts"
}
