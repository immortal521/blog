package service

import (
	"blog-server/internal/config"
	"blog-server/pkg/errs"
	"blog-server/templates" // 假设这是包含 embed.FS 的包
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

// IMailService 接口变得更加通用
type IMailService interface {
	Send(to, subject, templateName string, data any) error
}

type mailService struct {
	dialer   *gomail.Dialer
	template *template.Template
	from     string // 将发件人地址存储起来
}

func NewEmailService() (IMailService, error) {
	mailCfg := config.Get().Email
	dialer := gomail.NewDialer(mailCfg.Host, mailCfg.Port, mailCfg.Username, mailCfg.Password)

	// 一次性解析所有模板
	t, err := template.ParseFS(templates.FS, "*.html")
	if err != nil {
		return nil, errs.New(errs.CodeInternalError, "failed to parse email templates", err)
	}

	return &mailService{
		dialer:   dialer,
		template: t,
		from:     mailCfg.From, // 从配置中获取发件人地址
	}, nil
}

// Send 成为了公开的、通用的发送方法
func (s *mailService) Send(to, subject, templateName string, data any) error {
	var body bytes.Buffer
	if err := s.template.ExecuteTemplate(&body, templateName, data); err != nil {
		return errs.New(errs.CodeInternalError, fmt.Sprintf("Failed to execute email template %s", templateName), err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())

	if err := s.dialer.DialAndSend(m); err != nil {
		return errs.New(errs.CodeExternalAPIError, fmt.Sprintf("Failed to send email to %s", to), err)
	}

	return nil
}
