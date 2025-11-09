package service

import (
	"blog-server/internal/cache"
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"blog-server/pkg/errs"
	"blog-server/pkg/logger"
	"errors"

	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type IPostService interface {
	GetPosts(ctx context.Context) ([]*entity.Post, error)
	GetPostsWithContent(ctx context.Context) ([]*entity.Post, error)
	GetPostsMeta(ctx context.Context) []*entity.Post
	GetPostByID(ctx context.Context, id uint) (*entity.Post, error)
	FlushViewCountToDB(ctx context.Context) error
}

type postService struct {
	log      logger.Logger
	db       database.DB
	rc       cache.CacheClient
	postRepo repo.IPostRepo
}

func NewPostService(log logger.Logger, db database.DB, rc cache.CacheClient, postRepo repo.IPostRepo) IPostService {
	return &postService{log: log, db: db, rc: rc, postRepo: postRepo}
}

func (p *postService) GetPosts(ctx context.Context) ([]*entity.Post, error) {
	return p.postRepo.GetAllPosts(ctx, p.db.Conn())
}

func (p *postService) GetPostsWithContent(ctx context.Context) ([]*entity.Post, error) {
	return p.postRepo.GetAllPostsWithContent(ctx, p.db.Conn())
}

func (p *postService) GetPostsMeta(ctx context.Context) []*entity.Post {
	posts, err := p.postRepo.GetPostsMeta(ctx, p.db.Conn())
	if err != nil {
		p.log.Error("get posts meta failed", logger.Error(err))
		return []*entity.Post{}
	}
	return posts
}

func (p *postService) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	post, err := p.postRepo.GetPostByID(ctx, p.db.Conn(), id)
	if err != nil {
		return nil, err
	}

	go func(postID uint) {
		if err := p.rc.Raw().Incr(ctx, fmt.Sprintf("blog:post:view_count:%d", post.ID)).Err(); err != nil {
			p.log.Error("incr post view count failed", logger.Error(err))
		}
	}(post.ID)

	return post, nil
}

func (p *postService) FlushViewCountToDB(ctx context.Context) error {
	var cursor uint64
	var allErrors []error

	for {
		keys, next, err := p.rc.Raw().Scan(ctx, cursor, "blog:post:view_count:*", 100).Result()
		if err != nil {
			return errs.New(errs.CodeCacheError, "Failed to scan Redis keys for post view count", err)
		}
		cursor = next

		if len(keys) == 0 && cursor == 0 {
			break
		}
		if len(keys) == 0 {
			continue
		}

		pipe := p.rc.Raw().Pipeline()
		cmds := make([]*redis.StringCmd, len(keys))

		for i, key := range keys {
			cmds[i] = pipe.GetDel(ctx, key)
		}

		if _, err = pipe.Exec(ctx); err != nil && !errors.Is(err, redis.Nil) {
			allErrors = append(allErrors, fmt.Errorf("pipeline exec failed: %w", err))
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
				allErrors = append(allErrors, fmt.Errorf("db update failed: %w", err))
			}
		}

		if cursor == 0 {
			break
		}
	}
	if len(allErrors) > 0 {
		return errs.New(errs.CodeInternalError, fmt.Sprintf("Flush completed with %d errors", len(allErrors)), fmt.Errorf("%v", allErrors))
	}
	return nil
}
