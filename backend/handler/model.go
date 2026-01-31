package handler

import (
	"bufio"
	"fmt"
	"strings"
	"sync"

	"blog-server/errs"
	"blog-server/service"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IModelHandler interface {
	CreateSummarySession(c *fiber.Ctx) error
	SummaryStream(c *fiber.Ctx) error
}

type modelHandler struct {
	svn      service.IModelService
	sessions *sync.Map
}

type Session struct {
	Content  string
	TaskType string
	TextCh   <-chan string
	ErrCh    <-chan error
}

func (h *modelHandler) CreateSummarySession(c *fiber.Ctx) error {
	req := new(struct {
		Content string `json:"content" validate:"required"`
	})
	if err := c.BodyParser(req); err != nil {
		return errs.New(errs.CodeInvalidParam, "Invalid request body", err)
	}

	sessionID := uuid.New().String()

	session := &Session{
		Content:  req.Content,
		TaskType: "summary",
		TextCh:   nil,
		ErrCh:    nil,
	}

	h.sessions.Store(sessionID, session)

	return c.JSON(fiber.Map{
		"sessionId": sessionID,
	})
}

func (h *modelHandler) SummaryStream(c *fiber.Ctx) error {
	sessionID := c.Params("sessionId")
	return h.stream(c, sessionID)
}

func (h *modelHandler) stream(c *fiber.Ctx, sessionID string) error {
	v, ok := h.sessions.Load(sessionID)
	if !ok {
		fmt.Println("session not found")
		return nil
	}
	session := v.(*Session)

	if session.TextCh == nil && session.ErrCh == nil {
		switch session.TaskType {
		case "summary":
			session.TextCh, session.ErrCh = h.svn.GenerateSummary(c.UserContext(), session.Content)
		// case "translation":
		// 	session.TextCh, session.ErrCh = h.svn.GenerateTranslation(c.UserContext(), session.Content)
		default:
			return errs.New(errs.CodeInvalidParam, "unsupported task type", nil)
		}
		h.sessions.Store(sessionID, session)
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			_ = w.Flush()
		}()

		for {
			select {
			case err := <-session.ErrCh:
				if err != nil {
					_, _ = fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
					_ = w.Flush()
					return
				}
			case text := <-session.TextCh:
				if text == "[DONE]" {
					_, _ = fmt.Fprintf(w, "event: done\ndata: %s\n\n", text)
					return
				}
				if text == "" {
					continue
				}

				lines := strings.Split(text, "\n")
				for _, line := range lines {
					_, _ = fmt.Fprintf(w, "data: %s\n\n", line)
				}
				_ = w.Flush()
			}
		}
	})

	return nil
}

func NewModelHandler(svn service.IModelService) IModelHandler {
	return &modelHandler{
		svn:      svn,
		sessions: &sync.Map{},
	}
}
