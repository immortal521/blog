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
	GetPostIds(c *fiber.Ctx) error
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
	postDTOs := make([]dto.PostListRes, len(posts))
	for i, post := range posts {
		postDTOs[i] = dto.PostListRes{
			ID:              post.ID,
			Title:           post.Title,
			Cover:           post.Cover,
			Summary:         post.Summary,
			PublishedAt:     post.PublishedAt,
			UpdatedAt:       post.UpdatedAt,
			ReadTimeMinutes: post.ReadTimeMinutes,
			ViewCount:       post.ViewCount,
			Author:          post.User.Username,
		}
	}
	result := dto.Success(postDTOs)
	return c.JSON(result)
}

func (p PostHandler) GetPostIds(c *fiber.Ctx) error {
	metas := p.svc.GetPostsMeta(c.UserContext())
	var metasDTO = make([]dto.PostMetaRes, len(metas))
	for i, meta := range metas {
		metasDTO[i] = dto.PostMetaRes{
			ID:        meta.ID,
			UpdatedAt: meta.UpdatedAt,
		}
	}
	result := dto.Success(metasDTO)
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
		Author:          post.User.Username,
	})
	return c.JSON(result)
}
