package repo

import (
	"blog-server/internal/entity"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PostRepo interface {
	GetAllPosts(ctx context.Context, db *gorm.DB) ([]entity.Post, error)
	GetPostsMeta(ctx context.Context, db *gorm.DB) ([]entity.Post, error)
	GetPostById(ctx context.Context, db *gorm.DB, id uint) (entity.Post, error)
}

type postRepo struct{}

func NewPostRepo() PostRepo {
	return &postRepo{}
}

func (r *postRepo) GetAllPosts(ctx context.Context, db *gorm.DB) ([]entity.Post, error) {
	posts, err := gorm.G[entity.Post](db).Joins(clause.JoinTarget{Association: "User"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
		db.Select("username")
		return nil
	}).Select("posts.id", "posts.title", "posts.summary", "posts.cover", "posts.read_time_minutes", "posts.view_count", "posts.published_at", "posts.updated_at").Find(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostsMeta(ctx context.Context, db *gorm.DB) ([]entity.Post, error) {
	posts, err := gorm.G[entity.Post](db).Select("id", "updated_at").Find(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostById(ctx context.Context, db *gorm.DB, id uint) (entity.Post, error) {
	var post entity.Post
	post, err := gorm.G[entity.Post](db).Where("id = ?", id).First(ctx)
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}
