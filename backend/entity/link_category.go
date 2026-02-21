package entity

import "time"

type LinkCategory struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`

	Name      string `gorm:"column:name;size:20;not null;unique"`
	SortOrder int    `gorm:"not null;default:0"`
}

func (*LinkCategory) TableName() string {
	return "link_categories"
}
