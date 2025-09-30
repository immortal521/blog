// Package service
package service

import (
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"blog-server/pkg/util"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
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
	SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error
	Register(ctx context.Context, user *entity.User) error
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
}

type AuthService struct {
	rdb         *redis.Client
	userRepo    repo.IUserRepo
	mailService IMailService
	jwtService  IJwtService
}

func NewAuthService(rdb *redis.Client, jwtService IJwtService, mailService IMailService) IAuthService {
	return &AuthService{
		rdb:         rdb,
		jwtService:  jwtService,
		mailService: mailService,
	}
}

func (s *AuthService) Register(ctx context.Context, user *entity.User) error {

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	return "", "", nil
}

func (s *AuthService) SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error {

	if captchaType == "" {
		captchaType = Register
	}
	templateName := "captcha.html" // 指定要使用的模板文件

	data, err := getCaptchaEmailMeta(captchaType)
	if err != nil {
		return err
	}

	captcha := util.GenerateCaptcha()

	data.Captcha = captcha

	err = s.mailService.Send(to, data.Subject, templateName, data)
	if err != nil {
		return err
	}

	s.rdb.Set(ctx, fmt.Sprintf("%s:%s", captchaType, to), captcha, 5*time.Minute)

	return nil
}

var captchaMetaMap = map[CaptchaType]AuthMailData{
	Register: {
		Subject: "【Immortal's Blog】邮箱验证",
		Title:   "请验证您的注册邮箱",
		Type:    "注册",
		Content: "请使用此验证码完成注册操作",
	},
	PasswordReset: {
		Subject: "【Immortal's Blog】重置密码",
		Title:   "请验证您的密码重置请求",
		Type:    "重置密码",
		Content: "您正在进行密码重置操作，请使用此验证码完成验证",
	},
	ChangeEmail: {
		Subject: "【Immortal's Blog】更改邮箱",
		Title:   "请验证您的新邮箱地址",
		Type:    "更改邮箱",
		Content: "您正在进行更改邮箱操作，请使用此验证码验证您的新邮箱",
	},
}

func getCaptchaEmailMeta(t CaptchaType) (AuthMailData, error) {
	data, ok := captchaMetaMap[t]
	if !ok {
		return AuthMailData{}, fmt.Errorf("unknown captcha type: %s", t)
	}
	return data, nil
}
