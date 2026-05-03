package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"blog-server/cache"
	"blog-server/datastore"
	"blog-server/entity"
	"blog-server/logger"
	"blog-server/repository"
)

type PostService interface {
	GetPosts(ctx context.Context) ([]*entity.Post, error)
	GetPostsWithContent(ctx context.Context) ([]*entity.Post, error)
	GetPostsMeta(ctx context.Context) []*entity.Post
	GetPostByID(ctx context.Context, id uint) (*entity.Post, error)
	FlushViewCountToDB(ctx context.Context) error
	CreatePost(ctx context.Context, input *CreatePostInput) (*entity.Post, error)
}

type postService struct {
	log       logger.Logger
	rc        cache.CacheClient
	postRepo  repository.PostRepo
	datastore *datastore.DataStore
}

type CreatePostInput struct {
	Title           string
	Summary         *string
	Content         string
	Cover           *string
	ReadTimeMinutes uint
	Status          entity.PostStatus

	UserID uint

	CategoryIDs []uint
	Tags        []uint
}

func (p *postService) CreatePost(ctx context.Context, input *CreatePostInput) (*entity.Post, error) {
	var post *entity.Post
	var err error

	err = p.datastore.WithTx(ctx, func(ctx context.Context) error {
		post, err = p.postRepo.Create(ctx, &entity.Post{
			Title:           input.Title,
			Summary:         input.Summary,
			Cover:           input.Cover,
			Content:         input.Content,
			UserID:          input.UserID,
			ReadTimeMinutes: input.ReadTimeMinutes,
			Status:          input.Status,
		})
		if err != nil {
			return err
		}

		if err = p.postRepo.ReplaceTags(ctx, post.ID, input.Tags); err != nil {
			return err
		}

		if err = p.postRepo.ReplaceCategories(ctx, post.ID, input.CategoryIDs); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (p *postService) FlushViewCountToDB(ctx context.Context) error {
	var cursor uint64
	var allErrors []error
	p.log.Info("Flushing post view count to DB...")

	for {
		keys, next, err := p.rc.Scan(ctx, "blog:post:view_count:*", cursor, 100)
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
			if err := p.postRepo.UpdateViewCounts(ctx, updates); err != nil {
				allErrors = append(allErrors, fmt.Errorf("db update failed: %w", err))
			}
		}

		if cursor == 0 {
			break
		}
	}
	if len(allErrors) > 0 {
		return fmt.Errorf(fmt.Sprintf("Flush completed with %d errors", len(allErrors)), fmt.Errorf("%v", allErrors))
	}
	return nil
}

func NewPostService(
	log logger.Logger,
	rc cache.CacheClient,
	postRepo repository.PostRepo,
	datastore *datastore.DataStore,
) PostService {
	return &postService{
		log:       log,
		rc:        rc,
		datastore: datastore,
		postRepo:  postRepo,
	}
}

// GetPosts retrieves all published posts
func (p *postService) GetPosts(ctx context.Context) ([]*entity.Post, error) {
	return p.postRepo.GetAllPublished(ctx)
}

// GetPostsWithContent retrieves all posts with content
// WARNING:  The current repo does not support it yet, so reuse it for now.
func (p *postService) GetPostsWithContent(ctx context.Context) ([]*entity.Post, error) {
	return p.postRepo.GetAllPublished(ctx)
}

// GetPostsMeta retrieves metadata
func (p *postService) GetPostsMeta(ctx context.Context) []*entity.Post {
	posts, err := p.postRepo.GetAllPublished(ctx)
	if err != nil {
		p.log.Error("get posts meta failed", logger.Err(err))
		return []*entity.Post{}
	}
	return posts
}

// GetPostByID retrieves a single post and async
func (p *postService) GetPostByID(ctx context.Context, id uint) (*entity.Post, error) {
	post, err := p.postRepo.GetPublishedByID(ctx, id)
	if err != nil {
		return nil, err
	}

	go func(postID uint) {
		if _, err := p.rc.Incr(ctx, fmt.Sprintf("blog:post:view_count:%d", postID)); err != nil {
			p.log.Error("incr post view count failed", logger.Err(err))
		}
	}(post.ID)

	return post, nil
}
