package repo

import (
	"blog-server/internal/entity"

	"gorm.io/gorm"
)

type PostRepo interface {
	GetAllPosts(db *gorm.DB) ([]entity.Post, error)
	GetPostsMeta(db *gorm.DB) ([]entity.Post, error)
	GetPostById(db *gorm.DB, id uint) (entity.Post, error)
}

type postRepo struct{}

func NewPostRepo() PostRepo {
	return &postRepo{}
}

func (r *postRepo) GetAllPosts(db *gorm.DB) ([]entity.Post, error) {
	var posts []entity.Post
	// err := db.Joins("User", db.Select("username")).Select("posts.id", "posts.title", "posts.summary", "posts.cover", "posts.read_time_minutes", "posts.view_count", "posts.published_at", "posts.updated_at").Find(&posts).Error
	err := db.Preload("User").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostsMeta(db *gorm.DB) ([]entity.Post, error) {
	var posts []entity.Post
	err := db.Select("id", "updated_at").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostById(db *gorm.DB, id uint) (entity.Post, error) {
	var post entity.Post
	err := db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}
