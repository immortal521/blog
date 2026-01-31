package entity

import "gorm.io/gorm"

type PostTag struct {
	*gorm.Model

	Name string `gorm:"type:varchar(100);unique;not null" json:"name"`
	Slug string `gorm:"type:varchar(100);unique;not null" json:"slug"`

	Posts []Post `gorm:"many2many:post_tag_relations;"`
}

func (PostTag) TableName() string {
	return "post_tags"
}
