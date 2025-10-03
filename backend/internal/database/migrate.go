package database

import (
	"blog-server/internal/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// AutoMigrate 自动迁移数据库表结构。
// 它会使用 GORM 的 AutoMigrate 方法来创建或更新以下实体对应的表：
//   - User
//   - Link
//   - LinkCategory
//   - Post
//   - PostCategory
//   - PostTag
//
// 如果迁移过程中发生错误，则会返回该错误。
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
	)
}
