package entity

import "time"

type PostTag struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`

	Name string `gorm:"type:varchar(100);unique;not null" json:"name"`
	Slug string `gorm:"type:varchar(100);unique;not null" json:"slug"`

	Posts []Post `gorm:"many2many:post_tag_relations;"`
}

func (PostTag) TableName() string {
	return "post_tags"
}
