// Package service
package service

import (
	"context"
	"fmt"
)

type AuthMailData struct {
	Title   string
	Type    string
	Content string
	Subject string
	Captcha string
}

type CaptchaType string

const (
	Register      CaptchaType = "Register"
	PasswordReset CaptchaType = "PasswordReset"
	ChangeEmail   CaptchaType = "ChangeEmail" // 更改邮箱
)

type IAuthService interface {
	SendCaptchaMail(ctx context.Context, to string, captcha string, captchaType CaptchaType) error
}

type AuthService struct {
	mailService IMailService
}

func NewAuthService(mailService IMailService) IAuthService {
	return &AuthService{
		mailService: mailService,
	}
}

func (s *AuthService) SendCaptchaMail(ctx context.Context, to string, captcha string, captchaType CaptchaType) error {
	templateName := "captcha.html" // 指定要使用的模板文件

	data, err := getCaptchaEmailMeta(captchaType)
	if err != nil {
		return err
	}
	data.Captcha = captcha

	return s.mailService.Send(to, data.Subject, templateName, data)
}

func getCaptchaEmailMeta(t CaptchaType) (AuthMailData, error) {
	switch t {
	case Register:
		return AuthMailData{
			Subject: "【Immortal's Blog】邮箱验证",
			Title:   "请验证您的注册邮箱",
			Type:    "注册",
			Content: "请使用此验证码完成注册操作",
		}, nil
	case PasswordReset:
		return AuthMailData{
			Subject: "【Immortal's Blog】重置密码",
			Title:   "请验证您的密码重置请求",
			Type:    "重置密码",
			Content: "您正在进行密码重置操作，请使用此验证码完成验证",
		}, nil
	case ChangeEmail:
		return AuthMailData{
			Subject: "【Immortal's Blog】更改邮箱",
			Title:   "请验证您的新邮箱地址",
			Type:    "更改邮箱",
			Content: "您正在进行更改邮箱操作，请使用此验证码验证您的新邮箱",
		}, nil
	default:
		// 对于未知的类型，返回错误，保证程序的健壮性
		return AuthMailData{}, fmt.Errorf("unknown captcha type: %s", t)
	}
}
