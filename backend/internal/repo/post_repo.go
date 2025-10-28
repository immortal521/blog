package repo

import (
	"blog-server/internal/entity"
	"blog-server/pkg/errs"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPostRepo interface {
	GetAllPosts(ctx context.Context, db *gorm.DB) ([]*entity.Post, error)
	GetAllPostsWithContent(ctx context.Context, db *gorm.DB) ([]*entity.Post, error)
	GetPostsMeta(ctx context.Context, db *gorm.DB) ([]*entity.Post, error)
	GetPostByID(ctx context.Context, db *gorm.DB, id uint) (*entity.Post, error)
	UpdateViewCounts(ctx context.Context, db *gorm.DB, updates map[uint]int64) error
}

type postRepo struct{}

func NewPostRepo() IPostRepo {
	return &postRepo{}
}

func (r *postRepo) GetAllPosts(ctx context.Context, db *gorm.DB) ([]*entity.Post, error) {
	posts, err := gorm.G[*entity.Post](db).
		Joins(clause.JoinTarget{Association: "User"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Select("username")
			return nil
		}).
		Select("posts.id", "posts.title", "posts.summary", "posts.cover", "posts.read_time_minutes", "posts.view_count", "posts.published_at", "posts.updated_at").
		Where("status = ?", entity.PostPublished).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetAllPostsWithContent(ctx context.Context, db *gorm.DB) ([]*entity.Post, error) {
	posts, err := gorm.G[*entity.Post](db).
		Joins(clause.JoinTarget{Association: "User"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Select("username")
			return nil
		}).
		Select("posts.id", "posts.title", "posts.summary", "posts.content", "posts.cover", "posts.read_time_minutes", "posts.view_count", "posts.published_at", "posts.updated_at").
		Where("status = ?", entity.PostPublished).
		Find(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostsMeta(ctx context.Context, db *gorm.DB) ([]*entity.Post, error) {
	posts, err := gorm.G[*entity.Post](db).Select("id", "updated_at").Find(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *postRepo) GetPostByID(ctx context.Context, db *gorm.DB, id uint) (*entity.Post, error) {
	post, err := gorm.G[*entity.Post](db).
		Joins(clause.JoinTarget{Association: "User"}, func(db gorm.JoinBuilder, joinTable clause.Table, curTable clause.Table) error {
			db.Select("username")
			return nil
		}).
		Select("posts.id", "posts.title", "posts.summary", "posts.content", "posts.cover", "posts.read_time_minutes", "posts.view_count", "posts.published_at", "posts.updated_at").
		Where("status = ?", entity.PostPublished).
		Where("posts.id = ?", id).
		First(ctx)
	if err == nil {
		return post, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &entity.Post{}, errs.ErrPostNotFound
	}
	return &entity.Post{}, err
}

func (r *postRepo) UpdateViewCounts(ctx context.Context, db *gorm.DB, updates map[uint]int64) error {
	if len(updates) == 0 {
		return nil
	}

	var caseBuilder strings.Builder
	var idArgs []any
	var ids []any
	caseBuilder.WriteString("view_count + CASE id ")
	for id, val := range updates {
		caseBuilder.WriteString("WHEN ? THEN ? ")
		idArgs = append(idArgs, id, val)
		ids = append(ids, id)
	}
	caseBuilder.WriteString("ELSE 0 END")
	return db.Model(&entity.Post{}).
		Where("id IN (?)", ids).
		UpdateColumn("view_count", gorm.Expr(caseBuilder.String(), idArgs...)).Error
}
