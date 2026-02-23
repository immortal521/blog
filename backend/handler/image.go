package handler

import (
	"os"

	"blog-server/errs"
	"blog-server/request"
	"blog-server/response"
	"blog-server/service"
	"blog-server/validatorx"

	"github.com/gofiber/fiber/v2"
)

type IImageHandler interface {
	Upload(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Download(c *fiber.Ctx) error
}

type imageHandler struct {
	svc      service.IImageService
	validate validatorx.Validator
}

// Delete implements IImageHandler.
func (i *imageHandler) Delete(c *fiber.Ctx) error {
	panic("unimplemented")
}

// Download implements IImageHandler.
func (i *imageHandler) Download(c *fiber.Ctx) error {
	panic("unimplemented")
}

func (i *imageHandler) Upload(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	req := request.UploadReq{File: fileHeader}
	if err := i.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return errs.New(errs.CodeInvalidParam, "Failed to open file", err)
	}
	defer func() {
		_ = file.Close()
	}()

	tmp, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return errs.New(errs.CodeInternalError, "Failed to create temp file", err)
	}
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	filename := fileHeader.Filename
	image, err := i.svc.Upload(c.UserContext(), nil, filename, file, fileHeader.Header.Get("Content-Type"), fileHeader.Size)
	if err != nil {
		return err
	}

	res := response.Success(image)
	return c.JSON(res)
}

func NewImageHandler(svc service.IImageService, validate validatorx.Validator) IImageHandler {
	return &imageHandler{svc, validate}
}
