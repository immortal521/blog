package handler

import (
	"blog-server/internal/dto"
	"blog-server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type IPostHandler interface {
	GetPosts(c *fiber.Ctx) error
	GetPost(c *fiber.Ctx) error
}

type PostHandler struct {
	svc service.IPostService
}

func NewPostHandler(svc service.IPostService) *PostHandler {
	return &PostHandler{svc: svc}
}

func (p PostHandler) GetPosts(c *fiber.Ctx) error {
	posts, err := p.svc.GetPosts(c.UserContext())
	if err != nil {
		return err
	}
	postDTOs := make([]dto.PostShortRes, len(posts))
	for i, post := range posts {
		postDTOs[i] = dto.PostShortRes{
			ID:    post.ID,
			Title: post.Title,
			Cover: post.Cover,
		}
	}
	result := dto.Success(postDTOs)
	return c.JSON(result)
}

func (p PostHandler) GetPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	post, err := p.svc.GetPostByID(c.UserContext(), uint(id))
	if err != nil {
		return err
	}
	result := dto.Success(dto.PostRes{
		ID:              post.ID,
		Title:           post.Title,
		Summary:         post.Summary,
		Content:         post.Content,
		Cover:           post.Cover,
		ReadTimeMinutes: post.ReadTimeMinutes,
		ViewCount:       post.ViewCount,
		PublishedAt:     post.PublishedAt,
	})
	return c.JSON(result)
}
