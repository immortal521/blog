package entity

import (
	"time"

	"github.com/google/uuid"
)

type ImageFolder struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	ParentID  *uuid.UUID `gorm:"type:uuid;index"`
	Name      string     `gorm:"type:text;not null"`
	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	DeletedAt *time.Time `gorm:"type:timestamptz;index"`
}

func (ImageFolder) TableName() string {
	return "image_folders"
}
