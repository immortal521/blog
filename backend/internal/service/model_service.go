package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
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
	textCh := make(chan string, 10) // 加缓冲
	errCh := make(chan error, 1)

	go func() {
		defer func() {
			textCh <- "__EOF__"
			close(textCh)
		}()
		defer close(errCh)

		payload := map[string]any{
			"model":      "gpt-4.1-mini",
			"messages":   []map[string]string{{"role": "user", "content": "请帮我总结这篇文章：" + content}},
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
				// textCh <- "[DONE]"
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
