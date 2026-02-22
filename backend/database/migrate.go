package database

import (
	"blog-server/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AutoMigrate handles the database schema migration for all defined entities.
func AutoMigrate(db *gorm.DB) error {
	return db.Session(&gorm.Session{
		Logger: db.Logger.LogMode(logger.Silent),
	}).AutoMigrate(
		&entity.User{},
		&entity.Link{},
		&entity.LinkCategory{},
		&entity.Post{},
		&entity.PostCategory{},
		&entity.PostTag{},
		&entity.Image{},
		&entity.ImageFolder{},
	)
}
