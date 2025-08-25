package service

import (
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"context"
)

type IPostService interface {
	GetPosts(ctx context.Context) ([]entity.Post, error)
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
	posts, err := p.postRepo.GetAllPostsShort(p.db.Conn(ctx))
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (p postService) GetPostByID(ctx context.Context, id uint) (entity.Post, error) {
	post, err := p.postRepo.GetPostById(p.db.Conn(ctx), id)
	if err != nil {
		return entity.Post{}, err
	}
	return post, nil
}
