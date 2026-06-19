package handler

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"blog-server/pkg/errx"
	"blog-server/response"
	"blog-server/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type ModelHandler interface {
	CreateSummarySession(c *echo.Context) error
	SummaryStream(c *echo.Context) error
}

type modelHandler struct {
	svn      service.ModelService
	sessions *sync.Map
}

func RegisterModelRoutes(r *echo.Group, h ModelHandler) {
	group := r.Group("/model")
	group.POST("/summarize", h.CreateSummarySession)
	group.GET("/summarize/:sessionId", h.SummaryStream)
}

type Session struct {
	Content  string
	TaskType string
	TextCh   <-chan string
	ErrCh    <-chan error
}

func (h *modelHandler) CreateSummarySession(c *echo.Context) error {
	req := new(struct {
		Content string `json:"content" validate:"required"`
	})
	if err := c.Bind(req); err != nil {
		return errx.New(errx.CodeInvalidParam, err)
	}

	sessionID := uuid.New().String()

	session := &Session{
		Content:  req.Content,
		TaskType: "summary",
		TextCh:   nil,
		ErrCh:    nil,
	}

	h.sessions.Store(sessionID, session)

	return response.OK(c, map[string]string{
		"sessionId": sessionID,
	})
}

func (h *modelHandler) SummaryStream(c *echo.Context) error {
	sessionID := c.Param("sessionId")
	return h.stream(c, sessionID)
}

func (h *modelHandler) stream(c *echo.Context, sessionID string) error {
	v, ok := h.sessions.Load(sessionID)
	if !ok {
		fmt.Println("session not found")
		return nil
	}
	session := v.(*Session)
	defer h.sessions.Delete(sessionID)

	if session.TextCh == nil && session.ErrCh == nil {
		switch session.TaskType {
		case "summary":
			session.TextCh, session.ErrCh = h.svn.GenerateSummary(c.Request().Context(), session.Content)
		default:
			return errx.New(errx.CodeInvalidParam, fmt.Errorf("unsupported task type"))
		}
		h.sessions.Store(sessionID, session)
	}

	w := c.Response()
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	rc := http.NewResponseController(w)

	for {
		select {
		case <-c.Request().Context().Done():
			return nil

		case err := <-session.ErrCh:
			if err != nil {
				if _, werr := fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error()); werr != nil {
					return werr
				}
				return rc.Flush()
			}

		case text := <-session.TextCh:
			if text == "[DONE]" {
				if _, werr := fmt.Fprintf(w, "event: done\ndata: %s\n\n", text); werr != nil {
					return werr
				}
				return rc.Flush()
			}
			if text == "" {
				continue
			}
			for _, line := range strings.Split(text, "\n") {
				if _, werr := fmt.Fprintf(w, "data: %s\n\n", line); werr != nil {
					return werr
				}
			}
			if err := rc.Flush(); err != nil {
				return err
			}
		}
	}
}

func NewModelHandler(svn service.ModelService) ModelHandler {
	return &modelHandler{
		svn:      svn,
		sessions: &sync.Map{},
	}
}
