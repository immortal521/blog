package handler

import (
	"fmt"
	"strconv"
	"unicode/utf8"

	"blog-server/pkg/errx"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"

	"github.com/gofiber/fiber/v3"
)

// PostHandler defines the interface for post HTTP handlers
type PostHandler interface {
	GetPosts(c fiber.Ctx) error
	GetPost(c fiber.Ctx) error
	GetPostIds(c fiber.Ctx) error
	CreatePost(c fiber.Ctx) error
}

type postHandler struct {
	svc service.PostService
}

// NewPostHandler creates a new post handler instance
func NewPostHandler(svc service.PostService) PostHandler {
	return &postHandler{svc: svc}
}

func RegisterPostRoute(r fiber.Router, handler PostHandler) {
	group := r.Group("/posts")
	group.Post("/", handler.CreatePost)
	group.Get("/", handler.GetPosts)
	group.Get("/meta", handler.GetPostIds)
	group.Get("/:id", handler.GetPost)
}

func (h *postHandler) CreatePost(c fiber.Ctx) error {
	req := new(request.CreatePostReq)
	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, "Failed to parse request body", err)
	}

	contentLength := utf8.RuneCountInString(req.Content)

	readTime := uint(max(1, (contentLength+199)/200))

	input := &service.CreatePostInput{
		Title:           req.Title,
		Summary:         req.Summary,
		Cover:           req.Cover,
		Content:         req.Content,
		UserID:          req.UserID,
		ReadTimeMinutes: readTime,
		Status:          req.Status,
		Tags:            req.Tags,
		CategoryIDs:     req.CategoryIDs,
	}

	post, err := h.svc.CreatePost(c.Context(), input)
	if err != nil {
		return err
	}

	res := &response.PostRes{
		ID:              post.ID,
		Title:           post.Title,
		Summary:         post.Summary,
		Content:         post.Content,
		Cover:           post.Cover,
		ReadTimeMinutes: post.ReadTimeMinutes,
		ViewCount:       post.ViewCount,
		PublishedAt:     post.PublishedAt,
		UpdatedAt:       post.UpdatedAt,
	}

	return c.JSON(response.Success(res))
}

// GetPosts retrieves all published posts
func (h *postHandler) GetPosts(c fiber.Ctx) error {
	posts, err := h.svc.GetPosts(c.Context())
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
		}
	}

	return c.JSON(response.Success(postDTOs))
}

// GetPostIds retrieves metadata (id and updated_at) for all posts
func (h *postHandler) GetPostIds(c fiber.Ctx) error {
	metas := h.svc.GetPostsMeta(c.Context())

	metasDTO := make([]response.PostMetaRes, len(metas))
	for i, meta := range metas {
		metasDTO[i] = response.PostMetaRes{
			ID:        meta.ID,
			UpdatedAt: meta.UpdatedAt,
		}
	}

	return c.JSON(response.Success(metasDTO))
}

// GetPost retrieves a single published post by ID
func (h *postHandler) GetPost(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fmt.Errorf("invalid post ID: %w", err)
	}

	post, err := h.svc.GetPostByID(c.Context(), uint(id))
	if err != nil {
		return err
	}

	res := response.PostRes{
		ID:              post.ID,
		Title:           post.Title,
		Summary:         post.Summary,
		Content:         post.Content,
		Cover:           post.Cover,
		ReadTimeMinutes: post.ReadTimeMinutes,
		ViewCount:       post.ViewCount,
		PublishedAt:     post.PublishedAt,
		// Author:          safeUsername(post),
	}

	return c.JSON(response.Success(res))
}
