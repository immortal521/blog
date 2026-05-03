package request

import "blog-server/entity"

type CreatePostReq struct {
	Title   string            `json:"title"`
	Summary *string           `json:"summary"`
	Cover   *string           `json:"cover"`
	Content string            `json:"content"`
	Status  entity.PostStatus `json:"status"`

	UserID uint `json:"userId"`

	CategoryIDs []uint `json:"category_ids"`
	Tags        []uint `json:"tags"`
}
