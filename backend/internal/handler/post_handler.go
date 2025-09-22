package handler

import (
	"blog-server/internal/dto/response"
	"blog-server/internal/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type IPostHandler interface {
	GetPosts(c *fiber.Ctx) error
	GetPost(c *fiber.Ctx) error
	GetPostIds(c *fiber.Ctx) error
}

type postHandler struct {
	svc service.IPostService
}

func NewPostHandler(svc service.IPostService) IPostHandler {
	return &postHandler{svc: svc}
}

func (p *postHandler) GetPosts(c *fiber.Ctx) error {
	posts, err := p.svc.GetPosts(c.UserContext())
	if err != nil {
		return err
	}
	postDTOs := make([]response.PostListRes, len(posts))
	for i, post := range posts {
		postDTOs[i] = response.PostListRes{
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

	result := response.Success(postDTOs)
	return c.JSON(result)
}

func (p *postHandler) GetPostIds(c *fiber.Ctx) error {
	metas := p.svc.GetPostsMeta(c.UserContext())
	var metasDTO = make([]response.PostMetaRes, len(metas))
	for i, meta := range metas {
		metasDTO[i] = response.PostMetaRes{
			ID:        meta.ID,
			UpdatedAt: meta.UpdatedAt,
		}
	}
	result := response.Success(metasDTO)
	return c.JSON(result)
}

func (p *postHandler) GetPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	post, err := p.svc.GetPostByID(c.UserContext(), uint(id))
	if err != nil {
		return err
	}
	result := response.Success(response.PostRes{
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
