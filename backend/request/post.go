package request

import "blog-server/entity"

// CreatePostReq is the request body for creating a post.
// UserID is intentionally omitted — it is extracted from the JWT context
// by the handler, never accepted from the client.
type CreatePostReq struct {
	Title   string            `json:"title" validate:"required,min=3,max=100"`
	Summary *string           `json:"summary"`
	Cover   *string           `json:"cover"`
	Content string            `json:"content" validate:"required"`
	Status  entity.PostStatus `json:"status"`

	CategoryIDs []uint `json:"category_ids"`
	Tags        []uint `json:"tags"`
}

type PostPageReq struct {
	Page     int `json:"page" query:"page,default:1"`
	PageSize int `json:"pageSize" query:"pageSize,default:10"`
}
