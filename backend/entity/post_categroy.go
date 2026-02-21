package entity

import "time"

type PostCategory struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`

	Name string `gorm:"unique;not null;size:100"`
	Slug string `gorm:"unique;not null;size:100"`

	Posts []Post `gorm:"many2many:post_category_relations;"`
}

func (PostCategory) TableName() string {
	return "post_categories"
}
