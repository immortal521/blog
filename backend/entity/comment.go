package entity

import (
	"time"
)

type Comment struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`
}

func (c Comment) TableName() string {
	return "comments"
}
