package handler

import (
	"fmt"
	"strconv"

	"blog-server/contextx"
	"blog-server/entity"
	"blog-server/middleware"
	"blog-server/pkg/errx"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"

	"github.com/labstack/echo/v5"
)

// PostHandler defines the interface for post HTTP handlers
type PostHandler interface {
	GetPosts(c *echo.Context) error
	GetPost(c *echo.Context) error
	GetPostIds(c *echo.Context) error
	CreatePost(c *echo.Context) error
}

type postHandler struct {
	svc service.PostService
}

// NewPostHandler creates a new post handler instance
func NewPostHandler(svc service.PostService) PostHandler {
	return &postHandler{svc: svc}
}

func RegisterPostRoute(r *echo.Group, handler PostHandler, am *middleware.AuthMiddleware) {
	group := r.Group("/posts")
	group.POST("", handler.CreatePost, am.Handler())
	group.GET("", handler.GetPosts)
	group.GET("/meta", handler.GetPostIds)
	group.GET("/:id", handler.GetPost)
}

// toPostRes maps a domain Post to the detail response DTO.
func toPostRes(p *entity.Post) response.PostRes {
	return response.PostRes{
		ID:              p.ID,
		Title:           p.Title,
		Summary:         p.Summary,
		Content:         p.Content,
		Cover:           p.Cover,
		ReadTimeMinutes: p.ReadTimeMinutes,
		ViewCount:       p.ViewCount,
		PublishedAt:     p.PublishedAt,
		UpdatedAt:       p.UpdatedAt,
		Author:          p.User,
	}
}

// toPostListRes maps a domain Post to the list response DTO (no content).
func toPostListRes(p *entity.Post) response.PostListRes {
	return response.PostListRes{
		ID:              p.ID,
		Title:           p.Title,
		Summary:         p.Summary,
		Cover:           p.Cover,
		ReadTimeMinutes: p.ReadTimeMinutes,
		ViewCount:       p.ViewCount,
		PublishedAt:     p.PublishedAt,
		UpdatedAt:       p.UpdatedAt,
		Author:          p.User,
	}
}

func (h *postHandler) CreatePost(c *echo.Context) error {
	req := new(request.CreatePostReq)
	if err := c.Bind(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	u, ok := contextx.GetUser(c.Request().Context())
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

	post, err := h.svc.CreatePost(c.Request().Context(), u, input)
	if err != nil {
		return err
	}

	return response.OK(c, response.Success(toPostRes(post)))
}

// GetPosts retrieves all published posts
func (h *postHandler) GetPosts(c *echo.Context) error {
	query := new(request.PostPageReq)
	if err := c.Bind(query); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	posts, total, err := h.svc.GetPosts(c.Request().Context(), query.Page, query.PageSize)
	if err != nil {
		return err
	}

	postDTOs := make([]response.PostListRes, len(posts))
	for i, post := range posts {
		postDTOs[i] = toPostListRes(post)
	}

	return response.OK(c, response.Success(response.Page[response.PostListRes]{
		Total: total,
		List:  postDTOs,
	}))
}

// GetPostIds retrieves metadata (id and updated_at) for all posts
func (h *postHandler) GetPostIds(c *echo.Context) error {
	metas := h.svc.GetPostsMeta(c.Request().Context())

	metasDTO := make([]response.PostMetaRes, len(metas))
	for i, meta := range metas {
		metasDTO[i] = response.PostMetaRes{
			ID:        meta.ID,
			UpdatedAt: meta.UpdatedAt,
		}
	}

	return response.OK(c, response.Success(metasDTO))
}

// GetPost retrieves a single published post by ID
func (h *postHandler) GetPost(c *echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	post, err := h.svc.GetPostByID(c.Request().Context(), uint(id))
	if err != nil {
		return err
	}

	return response.OK(c, response.Success(toPostRes(post)))
}
