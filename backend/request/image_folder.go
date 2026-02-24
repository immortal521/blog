package request

import "github.com/google/uuid"

type CreateImageFolderReq struct {
	Name     string     `json:"name" validate:"required"`
	ParentID *uuid.UUID `json:"parentId"`
}

type RenameImageFolderReq struct {
	Name string    `json:"name" validate:"required"`
	ID   uuid.UUID `json:"id" validate:"required"`
}

type MoveImageFolderReq struct {
	ID             uuid.UUID  `json:"id" validate:"required"`
	TargetParentID *uuid.UUID `json:"targetParentId"`
}
