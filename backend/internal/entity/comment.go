package entity

import "gorm.io/gorm"

type Comment struct {
	*gorm.Model
}

func (c Comment) TableName() string {
	return "comments"
}
