package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"blog-server/internal/cache"
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repository"
	"blog-server/pkg/errs"
	"blog-server/pkg/logger"
)

// IPostService defines the interface for post business logic operations
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
	postRepo repository.IPostRepo
}

// NewPostService creates a new post service instance
func NewPostService(log logger.Logger, db database.DB, rc cache.CacheClient, postRepo repository.IPostRepo) IPostService {
	return &postService{log: log, db: db, rc: rc, postRepo: postRepo}
}

// GetPosts retrieves all published posts without content
func (p *postService) GetPosts(ctx context.Context) ([]*entity.Post, error) {
	return p.postRepo.GetAllPosts(ctx, p.db.Conn())
}

// GetPostsWithContent retrieves all published posts with full content
func (p *postService) GetPostsWithContent(ctx context.Context) ([]*entity.Post, error) {
	return p.postRepo.GetAllPostsWithContent(ctx, p.db.Conn())
}

// GetPostsMeta retrieves metadata (id and updated_at) for all posts
func (p *postService) GetPostsMeta(ctx context.Context) []*entity.Post {
	posts, err := p.postRepo.GetPostsMeta(ctx, p.db.Conn())
	if err != nil {
		p.log.Error("get posts meta failed", logger.Error(err))
		return []*entity.Post{}
	}
	return posts
}

// GetPostByID retrieves a single published post by ID and increments the view count in Redis
func (p *postService) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	post, err := p.postRepo.GetPostByID(ctx, p.db.Conn(), id)
	if err != nil {
		return nil, err
	}

	// Increment view count asynchronously
	go func(postID uint) {
		if _, err := p.rc.Incr(ctx, fmt.Sprintf("blog:post:view_count:%d", post.ID)); err != nil {
			p.log.Error("incr post view count failed", logger.Error(err))
		}
	}(post.ID)

	return post, nil
}

// FlushViewCountToDB flushes cached view counts from Redis to the database
func (p *postService) FlushViewCountToDB(ctx context.Context) error {
	var cursor uint64
	var allErrors []error

	for {
		keys, next, err := p.rc.Scan(ctx, "blog:post:view_count:*", cursor, 100)
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

		data, err := p.rc.PopBatch(ctx, keys)
		if err != nil {
			allErrors = append(allErrors, fmt.Errorf("PopBatch failed: %w", err))
		}

		updates := make(map[uint]int64)
		for key, valStr := range data {
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
