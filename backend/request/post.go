package request

import "blog-server/entity"

type CreatePostReq struct {
	Title   string            `json:"title" validate:"required,min=3,max=100"`
	Summary *string           `json:"summary"`
	Cover   *string           `json:"cover"`
	Content string            `json:"content" validate:"required"`
	Status  entity.PostStatus `json:"status"`

	UserID uint `json:"userId"`

	CategoryIDs []uint `json:"category_ids"`
	Tags        []uint `json:"tags"`
}
