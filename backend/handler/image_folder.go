package handler

import "github.com/gofiber/fiber/v2"

type IImageFolderHandler interface {
	Create(c *fiber.Ctx) error
	Rename(c *fiber.Ctx) error
	Move(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Tree(c *fiber.Ctx) error
	ListFolders(c *fiber.Ctx) error
	GetFolder(c *fiber.Ctx) error
}

type imageFolderHandler struct{}

// Create implements [IImageFolderHandler].
func (i *imageFolderHandler) Create(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Delete implements [IImageFolderHandler].
func (i *imageFolderHandler) Delete(c *fiber.Ctx) error {
	panic("unimplemented")
}

// GetFolder implements [IImageFolderHandler].
func (i *imageFolderHandler) GetFolder(c *fiber.Ctx) error {
	panic("unimplemented")
}

// ListFolders implements [IImageFolderHandler].
func (i *imageFolderHandler) ListFolders(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Move implements [IImageFolderHandler].
func (i *imageFolderHandler) Move(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Rename implements [IImageFolderHandler].
func (i *imageFolderHandler) Rename(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Tree implements [IImageFolderHandler].
func (i *imageFolderHandler) Tree(c *fiber.Ctx) error {
	panic("unimplemented")
}

func NewImageFolderHandler() IImageFolderHandler {
	return &imageFolderHandler{}
}
