package handler

import (
	"blog-server/internal/dto"
	"blog-server/internal/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
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
	posts, err := p.svc.GetPosts(c.Context())
	if err != nil {
		return err
	}
	postDTOs := make([]dto.PostShortResponseDTO, len(posts))
	for i, post := range posts {
		postDTOs[i] = dto.PostShortResponseDTO{
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
	post, err := p.svc.GetPostByID(c.Context(), uint(id))
	if err != nil {
		return err
	}
	result := dto.Success(dto.PostResponseDTO{
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
