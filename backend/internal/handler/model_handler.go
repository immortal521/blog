package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"blog-server/internal/request"
	"blog-server/internal/service"
	"blog-server/pkg/errs"
	"blog-server/pkg/validatorx"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Session struct {
	Content  string
	TaskType string // "summary", "translation", etc.
	TextCh   <-chan string
	ErrCh    <-chan error
}

type IModelHandler interface {
	SSE(c echo.Context) error
	ProcessContent(c echo.Context) error
}

type modelHandler struct {
	svn      service.IModelService
	validate validatorx.Validator
	sessions sync.Map
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

	session := &Session{
		Content:  req.Content,
		TaskType: "summary",
		TextCh:   nil, // 在 SSE 中初始化
		ErrCh:    nil,
	}

	h.sessions.Store(sessionID, session)

	return c.JSON(http.StatusOK, echo.Map{
		"sseUrl":    fmt.Sprintf("/sse?sessionId=%s", sessionID),
		"sessionId": sessionID,
	})
}

func (h *modelHandler) SSE(c echo.Context) error {
	sessionID := c.QueryParam("sessionId")
	v, ok := h.sessions.Load(sessionID)
	if !ok {
		return errs.New(errs.CodeInvalidParam, "invalid session id", nil)
	}
	session := v.(*Session)

	if session.TextCh == nil && session.ErrCh == nil {
		session.TextCh, session.ErrCh = h.svn.GenerateSummary(c.Request().Context(), session.Content)
		h.sessions.Store(sessionID, session)
	}

	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")
	c.Response().WriteHeader(http.StatusOK)

	w := c.Response().Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported")
	}

	defer h.sessions.Delete(sessionID)

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case err := <-session.ErrCh:
			if err != nil && errors.Is(err, context.Canceled) {
				_, _ = fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
				flusher.Flush()
			}

		case text := <-session.TextCh:
			if text == "" {
				continue
			}
			lines := strings.Split(text, "\n")
			for _, line := range lines {
				_, _ = fmt.Fprintf(w, "data: %s\n\n", line)
			}
			_, _ = fmt.Fprint(w, "\n")
			flusher.Flush()
		}
	}
}

func NewModelHandler(modelService service.IModelService, v validatorx.Validator) IModelHandler {
	return &modelHandler{
		svn:      modelService,
		validate: v,
		sessions: sync.Map{},
	}
}
