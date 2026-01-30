package entity

import "gorm.io/gorm"

type LinkCategory struct {
	*gorm.Model

	Name      string `gorm:"column:name;size:20;not null;unique"`
	SortOrder int    `gorm:"not null;default:0"`
}

func (*LinkCategory) TableName() string {
	return "link_categories"
}
