package service

import (
	"blog-server/internal/cache"
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"errors"

	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type IPostService interface {
	GetPosts(ctx context.Context) ([]*entity.Post, error)
	GetPostsMeta(ctx context.Context) []*entity.Post
	GetPostByID(ctx context.Context, id uint) (*entity.Post, error)
	FlushViewCountToDB(ctx context.Context) error
}

type postService struct {
	db       database.DB
	rc       cache.RedisClient
	postRepo repo.IPostRepo
}

func NewPostService(db database.DB, rc cache.RedisClient, postRepo repo.IPostRepo) IPostService {
	return &postService{db: db, rc: rc, postRepo: postRepo}
}

func (p *postService) GetPosts(ctx context.Context) ([]*entity.Post, error) {
	posts, err := p.postRepo.GetAllPosts(ctx, p.db.Conn())
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p *postService) GetPostsMeta(ctx context.Context) []*entity.Post {
	posts, err := p.postRepo.GetPostsMeta(ctx, p.db.Conn())
	if err != nil {
		return []*entity.Post{}
	}
	return posts
}

func (p *postService) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	post, err := p.postRepo.GetPostByID(ctx, p.db.Conn(), id)
	if err != nil {
		return &entity.Post{}, err
	}

	p.rc.Raw().Incr(ctx, fmt.Sprintf("blog:post:view_count:%d", post.ID))
	return post, nil
}

func (p *postService) FlushViewCountToDB(ctx context.Context) error {
	const batchSize = 1000
	var cursor uint64
	var errs []error

	for {
		keys, next, err := p.rc.Raw().Scan(ctx, cursor, "blog:post:view_count:*", 100).Result()
		if err != nil {
			return err
		}
		cursor = next
		if len(keys) == 0 {
			if cursor == 0 {
				break
			}
			continue
		}

		pipe := p.rc.Raw().Pipeline()
		cmds := make([]*redis.StringCmd, len(keys))

		for i, key := range keys {
			cmds[i] = pipe.GetDel(ctx, key)
		}

		if _, err = pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
			errs = append(errs, fmt.Errorf("pipeline exec failed: %w", err))
			continue
		}

		updates := make(map[uint]int64)
		for i, key := range keys {
			valStr, err := cmds[i].Result()
			if err != nil && !errors.Is(err, redis.Nil) {
				continue
			}

			count, err := strconv.ParseInt(valStr, 10, 64)
			if err != nil || count == 0 {
				continue
			}

			parts := strings.Split(key, ":")
			if len(parts) < 4 {
				continue
			}
			postID, err := strconv.ParseUint(parts[len(parts)-1], 10, 64)
			if err != nil {
				continue
			}

			updates[uint(postID)] += count
		}

		if len(updates) > 0 {
			if err := p.postRepo.UpdateViewCounts(ctx, p.db.Conn(), updates); err != nil {
				errs = append(errs, fmt.Errorf("db update failed: %w", err))
			}
		}

		if cursor == 0 {
			break
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("flush completed with %d errors: %v", len(errs), errs)
	}
	return nil
}
