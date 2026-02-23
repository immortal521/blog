package entity

import (
	"time"

	"github.com/google/uuid"
)

type Image struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey"`
	FolderID   *uuid.UUID `gorm:"type:uuid;index"`
	StorageKey string     `gorm:"type:text;not null;index"`
	OriginName string     `gorm:"type:text;not null"`
	Mime       string     `gorm:"type:text;not null"`
	Size       int64      `gorm:"not null"`
	Width      *int
	Height     *int
	Sha256     *string    `gorm:"type:text;index"`
	CreatedAt  time.Time  `gorm:"type:timestamptz;not null"`
	UpdatedAt  time.Time  `gorm:"type:timestamptz;not null"`
	DeletedAt  *time.Time `gorm:"type:timestamptz;index"`
}

func (Image) TableName() string {
	return "images"
}
