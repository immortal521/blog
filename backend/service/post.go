package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"blog-server/authz"
	"blog-server/cache"
	"blog-server/contextx"
	"blog-server/entity"
	"blog-server/logger"
	"blog-server/pkg/txmgr"
	"blog-server/repository"
)

// PostService defines the interface for post business logic operations.
type PostService interface {
	GetPosts(ctx context.Context, page, pageSize int) ([]*entity.Post, int, error)
	GetPostsWithContent(ctx context.Context) ([]*entity.Post, error)
	GetPostsMeta(ctx context.Context) []*entity.Post
	GetPostByID(ctx context.Context, id uint) (*entity.Post, error)
	CreatePost(ctx context.Context, user contextx.User, input *CreatePostInput) (*entity.Post, error)
	FlushViewCountToDB(ctx context.Context) error
}

// CreatePostInput groups all parameters for creating a post.
type CreatePostInput struct {
	Title   string
	Summary *string
	Content string
	Cover   *string
	Status  entity.PostStatus

	UserID uint

	CategoryIDs []uint
	Tags        []uint
}

// postService implements the PostService interface.
type postService struct {
	tx    txmgr.TxManager
	log   logger.Logger
	rc    cache.CacheClient
	pr    repository.PostRepo
	authz *authz.Authorizer
}

// NewPostService creates and returns a new PostService instance.
func NewPostService(
	tx txmgr.TxManager,
	log logger.Logger,
	pr repository.PostRepo,
	rc cache.CacheClient,
	authz *authz.Authorizer,
) PostService {
	return &postService{
		log:   log,
		rc:    rc,
		tx:    tx,
		pr:    pr,
		authz: authz,
	}
}

// GetPosts retrieves all published posts with pagination.
func (s *postService) GetPosts(ctx context.Context, page, pageSize int) ([]*entity.Post, int, error) {
	count, err := s.pr.CountPublished(ctx)
	if err != nil {
		return nil, 0, err
	}
	ps, err := s.pr.ListPublished(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	return ps, count, nil
}

// GetPostsWithContent retrieves all posts with content.
// WARNING: The current repo does not support it yet, so reuse it for now.
func (s *postService) GetPostsWithContent(ctx context.Context) ([]*entity.Post, error) {
	return nil, nil
}

// GetPostsMeta retrieves metadata for all published posts.
func (s *postService) GetPostsMeta(ctx context.Context) []*entity.Post {
	posts, err := s.pr.ListPublishedForSitemap(ctx)
	if err != nil {
		s.log.Error("get posts meta failed", logger.Err(err))
		return []*entity.Post{}
	}
	return posts
}

// GetPostByID retrieves a single post and asynchronously increments view count.
func (s *postService) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	post, err := s.pr.GetPublishedByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func(postID uint) {
		bgCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if _, err := s.rc.Incr(bgCtx, fmt.Sprintf("blog:post:view_count:%d", postID)); err != nil {
			s.log.Error("incr post view count failed", logger.Err(err))
		}
	}(post.ID)

	return post, nil
}

// CreatePost creates a new post with tags and categories.
func (s *postService) CreatePost(ctx context.Context, user contextx.User, input *CreatePostInput) (*entity.Post, error) {
	if err := s.authz.Authorize(ctx, user.ID, user.Role, authz.ResourcePost, authz.ActionCreate, nil); err != nil {
		return nil, err
	}
	contentLength := utf8.RuneCountInString(input.Content)

	readTime := uint(max(1, (contentLength+199)/200))

	var post *entity.Post
	var err error

	err = s.tx.WithTx(ctx, func(ctx context.Context) error {
		post, err = s.pr.Create(ctx, &entity.Post{
			Title:           input.Title,
			Summary:         input.Summary,
			Cover:           input.Cover,
			Content:         input.Content,
			ReadTimeMinutes: readTime,
			UserID:          input.UserID,
			Status:          input.Status,
		})
		if err != nil {
			return err
		}

		if err = s.pr.SetTags(ctx, post.ID, input.Tags); err != nil {
			return err
		}

		if err = s.pr.SetCategories(ctx, post.ID, input.CategoryIDs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

// FlushViewCountToDB flushes accumulated view counts from Redis to the database.
func (s *postService) FlushViewCountToDB(ctx context.Context) error {
	var cursor uint64
	var allErrors []error
	s.log.Info("Flushing post view count to DB...")

	for {
		keys, next, err := s.rc.Scan(ctx, "blog:post:view_count:*", cursor, 100)
		if err != nil {
			return fmt.Errorf("failed to scan Redis keys for post view count %w", err)
		}
		cursor = next

		if len(keys) == 0 && cursor == 0 {
			break
		}
		if len(keys) == 0 {
			continue
		}

		data, err := s.rc.PopBatch(ctx, keys)
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
			if err := s.pr.BatchIncrViewCounts(ctx, updates); err != nil {
				allErrors = append(allErrors, fmt.Errorf("db update failed: %w", err))
			}
		}

		if cursor == 0 {
			break
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf("flush completed with %d errors: %v", len(allErrors), allErrors)
	}
	return nil
}
