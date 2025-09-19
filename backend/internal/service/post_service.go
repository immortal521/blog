package service

import (
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
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
	rdb      *redis.Client
	postRepo repo.PostRepo
}

func NewPostService(db database.DB, rdb *redis.Client, postRepo repo.PostRepo) IPostService {
	return &postService{db: db, rdb: rdb, postRepo: postRepo}
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

	p.rdb.Incr(ctx, fmt.Sprintf("blog:post:view_count:%d", post.ID))
	return post, nil
}

func (p *postService) FlushViewCountToDB(ctx context.Context) error {
	var cursor uint64
	updates := make(map[uint]int64)

	for {
		keys, next, err := p.rdb.Scan(ctx, cursor, "blog:post:view_count:*", 100).Result()
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

		pipe := p.rdb.Pipeline()
		cmds := make([]*redis.StringCmd, 0, len(keys))

		for _, key := range keys {
			cmds = append(cmds, pipe.GetDel(ctx, key))
		}

		_, _ = pipe.Exec(ctx)

		for i, key := range keys {
			val, err := cmds[i].Int64()
			if err != nil && val == 0 {
				continue
			}

			parts := strings.Split(key, ":")
			postID, err := strconv.ParseUint(parts[len(parts)-1], 10, 64)
			if err != nil {
				continue
			}

			updates[uint(postID)] += val
		}

		if cursor == 0 {
			break
		}
	}
	if len(updates) == 0 {
		return nil
	}
	return p.postRepo.UpdateViewCounts(ctx, p.db.Conn(), updates)
}
