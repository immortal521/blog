package entity

import "gorm.io/gorm"

type PostCategory struct {
	*gorm.Model

	Name string `gorm:"unique;not null;size:100"`
	Slug string `gorm:"unique;not null;size:100"`

	Posts []Post `gorm:"many2many:post_category_relations;"`
}

func (PostCategory) TableName() string {
	return "post_categories"
}
