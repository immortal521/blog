package repo

import (
	"blog-server/internal/entity"
	"gorm.io/gorm"
)

type PostRepo interface {
	GetAllPostsShort(db *gorm.DB) ([]entity.Post, error)
	GetPostIDs(db *gorm.DB) ([]uint, error)
	GetPostById(db *gorm.DB, id uint) (entity.Post, error)
}

type postRepo struct{}

func NewPostRepo() PostRepo {
	return &postRepo{}
}

func (r *postRepo) GetAllPostsShort(db *gorm.DB) ([]entity.Post, error) {
	var posts []entity.Post
	err := db.Select("id", "title", "cover").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostIDs(db *gorm.DB) ([]uint, error) {
	var posts []entity.Post
	err := db.Select("id").Find(&posts).Error
	if err != nil {
		return nil, err
	}
	var ids = make([]uint, len(posts))
	for i, post := range posts {
		ids[i] = post.ID
	}
	return ids, nil
}

func (r *postRepo) GetPostById(db *gorm.DB, id uint) (entity.Post, error) {
	var post entity.Post
	err := db.Where("id = ?", id).First(&post).Error
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}
