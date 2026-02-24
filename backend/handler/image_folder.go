package handler

import (
	"time"

	"blog-server/errs"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"
	"blog-server/validatorx"

	"github.com/gofiber/fiber/v2"
)

type IImageFolderHandler interface {
	Create(c *fiber.Ctx) error
	Rename(c *fiber.Ctx) error
	Move(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Tree(c *fiber.Ctx) error
	ListFolders(c *fiber.Ctx) error
	GetFolder(c *fiber.Ctx) error
}

type imageFolderHandler struct {
	svc      service.IImageFolderService
	validate validatorx.Validator
}

func (i *imageFolderHandler) Create(c *fiber.Ctx) error {
	req := new(request.CreateImageFolderReq)
	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Failed to parse request body", err)
	}

	if err := i.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	folder, err := i.svc.Create(c.UserContext(), req.Name, req.ParentID)
	if err != nil {
		return err
	}

	out := response.ImageFolderRes{
		ID:                &folder.ID,
		ParentID:          folder.ParentID,
		Name:              folder.Name,
		ChildFoldersCount: 0,
		ImagesCount:       0,
		CreatedAt:         folder.CreatedAt.Format(time.RFC3339Nano),
		UpdatedAt:         folder.UpdatedAt.Format(time.RFC3339Nano),
	}
	return c.JSON(response.Success(out))
}

func (i *imageFolderHandler) Delete(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (i *imageFolderHandler) GetFolder(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (i *imageFolderHandler) ListFolders(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (i *imageFolderHandler) Move(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (i *imageFolderHandler) Rename(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (i *imageFolderHandler) Tree(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewImageFolderHandler() IImageFolderHandler {
	return &imageFolderHandler{}
}
