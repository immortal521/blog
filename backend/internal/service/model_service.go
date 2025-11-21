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
)

type IModelService interface {
	GenerateSummary(ctx context.Context, content string) (<-chan string, <-chan error)
}

type modelService struct {
	cfg *config.Config
}

func NewModelService(cfg *config.Config) IModelService {
	return &modelService{cfg: cfg}
}

func (s *modelService) GenerateSummary(ctx context.Context, content string) (<-chan string, <-chan error) {
	textCh := make(chan string) // 加缓冲避免阻塞
	errCh := make(chan error, 1)

	go func() {
		defer close(textCh)
		defer close(errCh)

		payload := map[string]any{
			"model":      "gpt-4.1-mini",
			"messages":   []map[string]string{{"role": "user", "content": "请帮我总结这篇文章：" + content}},
			"max_tokens": 8000,
			"stream":     true,
		}

		apiKey := s.cfg.LLM.APIKey

		body, _ := json.Marshal(payload)
		req, err := http.NewRequestWithContext(ctx, "POST", "https://models.inference.ai.azure.com/chat/completions", bytes.NewBuffer(body))
		if err != nil {
			errCh <- err
			return
		}
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 0} // SSE 不超时
		resp, err := client.Do(req)
		if err != nil {
			errCh <- err
			return
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		scanner := bufio.NewScanner(resp.Body)
		buf := make([]byte, 0, 1024*1024)
		scanner.Buffer(buf, 1024*1024)

		scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
			if atEOF && len(data) == 0 {
				return 0, nil, nil
			}
			if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
				return i + 2, data[:i], nil
			}
			if atEOF {
				return len(data), data, nil
			}
			return 0, nil, nil
		})

		doneRegexp := regexp.MustCompile(`data:\s*\[DONE\]`)

		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
			}

			chunk := strings.TrimSpace(scanner.Text())
			if !strings.HasPrefix(chunk, "data:") {
				// 忽略非 data 行 (如果 API 返回了其他行)
				continue
			}

			jsonStr := strings.TrimSpace(strings.TrimPrefix(chunk, "data:"))

			if doneRegexp.MatchString(chunk) {
				return
			}

			var respChunk map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &respChunk); err != nil {
				// 可以在这里通过 errCh 返回 JSON 解析错误
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
				return
			case textCh <- content:
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- err
		}
	}()

	return textCh, errCh
}
