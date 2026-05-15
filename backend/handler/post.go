package handler

import (
	"fmt"
	"strconv"

	"blog-server/contextx"
	"blog-server/middleware"
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

func RegisterPostRoute(r fiber.Router, handler PostHandler, am *middleware.AuthMiddleware) {
	group := r.Group("/posts")
	group.Post("/", am.Handler(), handler.CreatePost)
	group.Get("/", handler.GetPosts)
	group.Get("/meta", handler.GetPostIds)
	group.Get("/:id", handler.GetPost)
}

func (h *postHandler) CreatePost(c fiber.Ctx) error {
	req := new(request.CreatePostReq)
	if err := c.Bind().Body(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	u, ok := contextx.GetUser(c.Context())
	if !ok {
		return errx.New(errx.CodeUnauthorized, fmt.Errorf("missing user in context"))
	}

	input := &service.CreatePostInput{
		Title:       req.Title,
		Summary:     req.Summary,
		Cover:       req.Cover,
		Content:     req.Content,
		UserID:      u.ID,
		Status:      req.Status,
		Tags:        req.Tags,
		CategoryIDs: req.CategoryIDs,
	}

	post, err := h.svc.CreatePost(c.Context(), u, input)
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
	query := new(request.PostPageReq)
	if err := c.Bind().Query(query); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	posts, total, err := h.svc.GetPosts(c.Context(), query.Page, query.PageSize)
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
	pagedRes := response.Page[response.PostListRes]{
		Total: total,
		List:  postDTOs,
	}

	return c.JSON(response.Success(pagedRes))
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
		return errx.New(errx.CodeInvalidParam, fmt.Errorf("invalid post id"))
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
	}

	return c.JSON(response.Success(res))
}
