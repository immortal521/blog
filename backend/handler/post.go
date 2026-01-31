package handler

import (
	"strconv"

	"blog-server/errs"
	"blog-server/response"
	"blog-server/service"

	"github.com/gofiber/fiber/v2"
)

// IPostHandler defines the interface for post HTTP handlers
type IPostHandler interface {
	GetPosts(c *fiber.Ctx) error
	GetPost(c *fiber.Ctx) error
	GetPostIds(c *fiber.Ctx) error
}

type postHandler struct {
	svc service.IPostService
}

// NewPostHandler creates a new post handler instance
func NewPostHandler(svc service.IPostService) IPostHandler {
	return &postHandler{svc: svc}
}

// GetPosts retrieves all published posts
func (h *postHandler) GetPosts(c *fiber.Ctx) error {
	posts, err := h.svc.GetPosts(c.UserContext())
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

// GetPostIds retrieves metadata (id and updated_at) for all posts
func (h *postHandler) GetPostIds(c *fiber.Ctx) error {
	metas := h.svc.GetPostsMeta(c.UserContext())
	metasDTO := make([]response.PostMetaRes, len(metas))
	for i, meta := range metas {
		metasDTO[i] = response.PostMetaRes{
			ID:        meta.ID,
			UpdatedAt: meta.UpdatedAt,
		}
	}
	result := response.Success(metasDTO)
	return c.JSON(result)
}

// GetPost retrieves a single published post by ID
func (h *postHandler) GetPost(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid post ID", err)
	}
	post, err := h.svc.GetPostByID(c.UserContext(), uint(id))
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
