package service

import (
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"context"
)

type IPostService interface {
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetPostsMeta(ctx context.Context) []entity.Post
	GetPostByID(ctx context.Context, id uint) (entity.Post, error)
}

type postService struct {
	db       database.DB
	postRepo repo.PostRepo
}

func NewPostService(db database.DB, postRepo repo.PostRepo) IPostService {
	return &postService{db: db, postRepo: postRepo}
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
	return post, nil
}
