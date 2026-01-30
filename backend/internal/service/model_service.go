package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"blog-server/internal/config"
	"blog-server/pkg/errs"
	"blog-server/pkg/logger"
)

type IModelService interface {
	GenerateSummary(ctx context.Context, content string) (<-chan string, <-chan error)
}

type modelService struct {
	cfg *config.Config
	log logger.Logger
}

func (s *modelService) GenerateSummary(ctx context.Context, content string) (<-chan string, <-chan error) {
	textCh := make(chan string, 10)
	errCh := make(chan error, 1)

	go func() {
		defer close(textCh)
		defer close(errCh)

		payload := map[string]any{
			"model": "gpt-4.1-mini",
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": "Please summarize this article for me (return in plain text format): " + content,
				},
			},
			"max_tokens": 8000,
			"stream":     true,
		}

		body, _ := json.Marshal(payload)
		req, err := http.NewRequestWithContext(ctx, "POST",
			"https://models.inference.ai.azure.com/chat/completions", bytes.NewBuffer(body))
		if err != nil {
			errCh <- errs.New(errs.CodeInternalError, "llm request error", err)
			return
		}
		req.Header.Set("Authorization", "Bearer "+s.cfg.LLM.APIKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			errCh <- errs.New(errs.CodeInternalError, "llm request error", err)
			return
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			err = errs.New(errs.CodeInternalError, "llm request error", fmt.Errorf("status code: %d", resp.StatusCode))
			errCh <- err
			s.log.Errorf("llm request error %v", err)
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 0, 1024*1024)
		scanner.Buffer(buf, 1024*1024)

		doneRegexp := regexp.MustCompile(`data:\s*\[DONE\]`)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				s.log.Info("[GenerateSummary] context canceled")
				return
			default:
			}

			line := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(line, "data:") {
				continue
			}
			if doneRegexp.MatchString(line) {
				s.log.Info("[GenerateSummary] stream finished")
				textCh <- "[DONE]"
				return
			}

			jsonStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if jsonStr == "" {
				continue
			}

			var respChunk map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &respChunk); err != nil {
				s.log.Infof("[GenerateSummary] json parse error: %v, data: %s", err, jsonStr)
				continue
			}

			choices, ok := respChunk["choices"].([]any)
			if !ok || len(choices) == 0 {
				continue
			}
			choice, ok := choices[0].(map[string]any)
			if !ok {
				continue
			}
			delta, ok := choice["delta"].(map[string]any)
			if !ok {
				continue
			}
			content, ok := delta["content"].(string)
			if !ok || content == "" {
				continue
			}

			select {
			case <-ctx.Done():
				s.log.Info("[GenerateSummary] context canceled during sending")
				return
			case textCh <- content:
			}
		}

		if err := scanner.Err(); err != nil {
			s.log.Errorf("[GenerateSummary] scanner error: %v", err)
			errCh <- err
		}
	}()

	return textCh, errCh
}

func NewModelService(cfg *config.Config, log logger.Logger) IModelService {
	return &modelService{cfg: cfg, log: log}
}
