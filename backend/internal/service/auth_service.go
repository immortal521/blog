// Package service
package service

import (
	"blog-server/internal/database"
	"blog-server/internal/dto/request"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"blog-server/pkg/errs"
	"blog-server/pkg/util"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	Register(ctx context.Context, dto *request.RegisterReq) (accessToken, refreshToken string, err error)
	Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error)
}

type AuthService struct {
	db          database.DB
	rdb         *redis.Client
	userRepo    repo.IUserRepo
	jwtService  IJwtService
	mailService IMailService
}

func NewAuthService(db database.DB, rdb *redis.Client, userRepo repo.IUserRepo, jwtService IJwtService, mailService IMailService) IAuthService {
	return &AuthService{
		db:          db,
		rdb:         rdb,
		userRepo:    userRepo,
		jwtService:  jwtService,
		mailService: mailService,
	}
}

func (s *AuthService) Register(ctx context.Context, dto *request.RegisterReq) (accessToken, refreshToken string, err error) {
	email := dto.Email

	cachedCaptcha := s.rdb.Get(ctx, fmt.Sprintf("Register:%s", email)).Val()

	if !strings.EqualFold(cachedCaptcha, dto.Captcha) {
		return "", "", errs.ErrInvalidCaptcha
	}

	hashPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return "", "", err
	}

	user := &entity.User{
		UUID:     uuid.New(),
		Email:    dto.Email,
		Password: hashPassword,
	}

	err = s.db.Trans(func(txCtx *database.TxContext) error {
		return s.userRepo.CreateUser(ctx, txCtx.GetTx(), user)
	})
	if err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err = s.jwtService.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (accessToken, refreshToken string, err error) {
	return "", "", nil
}

func (s *AuthService) SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error {
	userIsExist, err := s.userRepo.ExistsByEmail(ctx, s.db.Conn(), to)
	if err != nil {
		return fmt.Errorf("check user existence: %w", err)
	}
	if userIsExist {
		return errs.ErrUserExists
	}

	if captchaType == "" {
		captchaType = Register
	}

	data, err := getCaptchaEmailMeta(captchaType)
	if err != nil {
		return err
	}

	captcha := util.GenerateCaptcha()
	data.Captcha = captcha

	if err = s.rdb.Set(ctx, fmt.Sprintf("%s:%s", captchaType, to), captcha, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("set captcha in redis: %w", err)
	}

	templateName := "captcha.html" // 指定要使用的模板文件
	err = s.mailService.Send(to, data.Subject, templateName, data)
	if err != nil {
		return err
	}

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
