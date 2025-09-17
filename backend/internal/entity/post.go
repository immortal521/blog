package entity

import (
	"time"

	"gorm.io/gorm"
)

// PostStatus 文章状态枚举
type PostStatus int

const (
	PostDraft PostStatus = iota + 1
	PostPublished
	PostArchived
)

type Post struct {
	*gorm.Model

	Title           string     `gorm:"type:varchar(255);not null"`
	Summary         *string    `gorm:"type:varchar(500)"`
	Content         string     `gorm:"type:longtext;not null"`
	Cover           *string    `gorm:"type:varchar(255)"`
	ReadTimeMinutes uint       `gorm:"not null"`
	ViewCount       uint       `gorm:"not null"`
	Status          PostStatus `gorm:"type:tinyint;default:1;not null"`

	UserID uint `gorm:"not null;index"`
	User   User `gorm:"foreignKey:UserID;references:ID"`

	PublishedAt *time.Time `gorm:"type:datetime(6)"`

	Categories []PostCategory `gorm:"many2many:post_category_relations;"`
	Tags       []PostTag      `gorm:"many2many:post_tag_relations;"`
}

func (Post) TableName() string {
	return "posts"
}
