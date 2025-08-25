package database

import (
	"blog-server/internal/entity"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Link{},
		&entity.LinkCategory{},
		&entity.Post{},
		&entity.PostCategory{},
		&entity.PostTag{},
	)
}
