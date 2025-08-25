package entity

import "gorm.io/gorm"

type LinkCategory struct {
	*gorm.Model

	Name      string `gorm:"column:name;size:20;not null;unique;comment:链接分类排序"`
	SortOrder int    `gorm:"not null;default:0;comment:分类排序顺序"`
}

func (*LinkCategory) TableName() string {
	return "link_categories"
}
