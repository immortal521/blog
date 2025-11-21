package handler

import (
	"net/http"
	"strconv"

	"blog-server/internal/response"
	"blog-server/internal/service"
	"blog-server/pkg/errs"

	"github.com/labstack/echo/v4"
)

type IPostHandler interface {
	GetPosts(c echo.Context) error
	GetPost(c echo.Context) error
	GetPostIds(c echo.Context) error
}

type postHandler struct {
	svc service.IPostService
}

func NewPostHandler(svc service.IPostService) IPostHandler {
	return &postHandler{svc: svc}
}

func (h *postHandler) GetPosts(c echo.Context) error {
	posts, err := h.svc.GetPosts(c.Request().Context())
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
	return c.JSON(http.StatusOK, result)
}

func (h *postHandler) GetPostIds(c echo.Context) error {
	metas := h.svc.GetPostsMeta(c.Request().Context())
	metasDTO := make([]response.PostMetaRes, len(metas))
	for i, meta := range metas {
		metasDTO[i] = response.PostMetaRes{
			ID:        meta.ID,
			UpdatedAt: meta.UpdatedAt,
		}
	}
	result := response.Success(metasDTO)
	return c.JSON(http.StatusOK, result)
}

func (h *postHandler) GetPost(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid post ID", err)
	}

	post, err := h.svc.GetPostByID(c.Request().Context(), uint(id))
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

	return c.JSON(http.StatusOK, result)
}
