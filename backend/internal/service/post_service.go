package service

import (
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type IPostService interface {
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPostsMeta(ctx context.Context) []entity.Post
	GetPostByID(ctx context.Context, id uint) (entity.Post, error)
	FlushViewCountToDB(ctx context.Context) error
}

type postService struct {
	db       database.DB
	rdb      *redis.Client
	postRepo repo.PostRepo
}

func NewPostService(db database.DB, rdb *redis.Client, postRepo repo.PostRepo) IPostService {
	return &postService{db: db, rdb: rdb, postRepo: postRepo}
}

func (p postService) GetPosts(ctx context.Context) ([]entity.Post, error) {
	posts, err := p.postRepo.GetAllPosts(ctx, p.db.Conn())
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p postService) GetPostsMeta(ctx context.Context) []entity.Post {
	posts, err := p.postRepo.GetPostsMeta(ctx, p.db.Conn())
	if err != nil {
		return []entity.Post{}
	}
	return posts
}

func (p postService) GetPostByID(ctx context.Context, id uint) (entity.Post, error) {
	post, err := p.postRepo.GetPostById(ctx, p.db.Conn(), id)
	if err != nil {
		return entity.Post{}, err
	}

	p.rdb.Incr(ctx, fmt.Sprintf("blog:post:view_count:%d", post.ID))
	return post, nil
}

func (p postService) FlushViewCountToDB(ctx context.Context) error {
	var cursor uint64
	for {
		keys, next, err := p.rdb.Scan(ctx, cursor, "blog:post:view_count:*", 100).Result()
		if err != nil {
			return err
		}
		cursor = next

		for _, key := range keys {
			val, err := p.rdb.GetDel(ctx, key).Int64()
			if err != nil && err != redis.Nil {
				return err
			}
			if val == 0 {
				continue
			}
			parts := strings.Split(key, ":")
			if len(parts) < 4 {
				continue
			}
			postID := parts[3]

			// 使用 gorm 表达式累加浏览数
			if err := p.db.Conn().Model(&entity.Post{}).
				Where("id = ?", postID).
				UpdateColumn("view_count", gorm.Expr("view_count + ?", val)).Error; err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
