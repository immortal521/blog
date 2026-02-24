package response

import "github.com/google/uuid"

type ImageFolderRes struct {
	ID                *uuid.UUID `json:"id"`
	ParentID          *uuid.UUID `json:"parentId"`
	Name              string     `json:"name"`
	ChildFoldersCount int64      `json:"childFoldersCount"`
	ImagesCount       int64      `json:"imagesCount"`
	CreatedAt         string     `json:"createdAt"`
	UpdatedAt         string     `json:"updatedAt"`
}
