package handler

import (
	"fmt"
	"net/http"

	"blog-server/internal/request"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type IModelHandler interface {
	SSE(c echo.Context) error
	ProcessContent(c echo.Context) error
}

type modelHandler struct {
	svn        service.IModelService
	validate   validatorx.Validator
	sseClients map[string]struct {
		TextCh <-chan string
		ErrCh  <-chan error
	}
}

func (h *modelHandler) ProcessContent(c echo.Context) error {
	req := new(request.SummarizeReq)
	if err := c.Bind(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Failed to parse request body", err)
	}

	if err := h.validate.Struct(req); err != nil {
		return errs.New(errs.CodeValidationFailed, "Validation failed", err)
	}

	sessionID := uuid.New().String()

	textCh, errCh := h.svn.GenerateSummary(c.Request().Context(), req.Content)

	// 把 channel 存入 map，以 sessionId 为 key
	h.sseClients[sessionID] = struct {
		TextCh <-chan string
		ErrCh  <-chan error
	}{
		TextCh: textCh,
		ErrCh:  errCh,
	}

	return c.JSON(http.StatusOK, echo.Map{
		"sseUrl":    fmt.Sprintf("/sse?sessionId=%s", sessionID),
		"sessionId": sessionID,
	})
}

func (h *modelHandler) SSE(c echo.Context) error {
	// sessionID := c.QueryParam("sessionId")
	// client, ok := h.sseClients[sessionID]
	// if !ok {
	// 	return errs.New(errs.CodeInvalidParam, "Invalid session id", nil)
	// }

	// c.Set("Content-Type", "text/event-stream")
	// c.Set("Cache-Control", "no-cache")
	// c.Set("Connection", "keep-alive")
	//
	// c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
	// 	for {
	// 		select {
	// 		case <-c.Context().Done():
	// 			delete(h.sseClients, sessionID)
	// 			return
	// 		case err, ok := <-client.ErrCh:
	// 			if ok && err != nil {
	// 				fmt.Fprintf(w, "data: %s\n\n", err.Error())
	// 			}
	// 			delete(h.sseClients, sessionID)
	// 			return
	// 		case text, ok := <-client.TextCh:
	// 			if !ok {
	// 				w.Write([]byte("data: [DONE]\n\n"))
	// 				delete(h.sseClients, sessionID)
	// 				return
	// 			}
	// 			fmt.Fprintf(w, "data: %s\n\n", text)
	// 			w.Flush()
	// 		}
	// 	}
	// }))
	//
	return nil
}

func NewModelHandler(modelService service.IModelService, v validatorx.Validator) IModelHandler {
	return &modelHandler{
		svn:      modelService,
		validate: v,
		sseClients: make(map[string]struct {
			TextCh <-chan string
			ErrCh  <-chan error
		}),
	}
}
