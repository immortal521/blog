// Package service
package service

import (
	"blog-server/internal/cache"
	"blog-server/internal/config"
	"blog-server/internal/database"
	"blog-server/internal/dto/request"
	"blog-server/internal/dto/response"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"blog-server/pkg/errs"
	"blog-server/pkg/util"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// AuthMailData represents the data required to send a captcha email.
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
	ChangeEmail   CaptchaType = "ChangeEmail"
)

// IAuthService is an interface for authentication services.
type IAuthService interface {
	SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error
	Register(ctx context.Context, dto *request.RegisterReq) (*response.LoginRes, error)
	Login(ctx context.Context, dto *request.LoginReq) (*response.LoginRes, error)
}

// AuthService implements the IAuthService interface.
type AuthService struct {
	db          database.DB
	rc          cache.RedisClient
	cfg         *config.Config
	userRepo    repo.IUserRepo
	jwtService  IJwtService
	mailService IMailService
}

// NewAuthService creates and returns a new AuthService instance
func NewAuthService(db database.DB, rc cache.RedisClient, userRepo repo.IUserRepo, jwtService IJwtService, mailService IMailService) IAuthService {
	return &AuthService{
		db:          db,
		rc:          rc,
		cfg:         config.Get(),
		userRepo:    userRepo,
		jwtService:  jwtService,
		mailService: mailService,
	}
}

// Register registers a new user and generated access/refresh tokens
func (s *AuthService) Register(ctx context.Context, dto *request.RegisterReq) (*response.LoginRes, error) {
	email := dto.Email

	cachedCaptcha := s.rc.Raw().Get(ctx, fmt.Sprintf("Register:%s", email)).Val()

	if !strings.EqualFold(cachedCaptcha, dto.Captcha) {
		return nil, errs.ErrInvalidCaptcha
	}

	hashPassword, err := util.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		UUID:     uuid.New(),
		Email:    dto.Email,
		Password: hashPassword,
		Username: util.GenerateUsername(),
	}
	var newUser *entity.User

	err = s.db.Trans(func(txCtx *database.TxContext) error {
		newUser, err = s.userRepo.CreateUser(ctx, txCtx.GetTx(), user)
		return err
	})
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.jwtService.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return nil, err
	}

	if err := s.cacheRefreshToken(ctx, user.UUID.String(), refreshToken); err != nil {
		return nil, err
	}

	result := &response.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UUID:         newUser.UUID.String(),
		Username:     newUser.Username,
		Avatar:       newUser.Avatar,
		Role:         newUser.Role,
	}

	return result, nil
}

// Login logs in a user and generates access/refresh tokens
func (s *AuthService) Login(ctx context.Context, dto *request.LoginReq) (*response.LoginRes, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, s.db.Conn(), dto.Email)
	if err != nil {
		return nil, err
	}
	if !util.VerifyPassword(dto.Password, user.Password) {
		return nil, errs.ErrPasswordWrong
	}

	accessToken, refreshToken, err := s.jwtService.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return nil, errs.ErrTokenGeneration
	}

	if err := s.cacheRefreshToken(ctx, user.UUID.String(), refreshToken); err != nil {
		return nil, err
	}

	result := &response.LoginRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UUID:         user.UUID.String(),
		Username:     user.Username,
		Avatar:       user.Avatar,
		Role:         user.Role,
	}

	return result, nil
}

// SendCaptchaMail generates a captcha, stores it in Redis, and sends an email.
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

	if err = s.rc.Raw().Set(ctx, fmt.Sprintf("%s:%s", captchaType, to), captcha, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("set captcha in redis: %w", err)
	}

	templateName := "captcha.html" // 指定要使用的模板文件
	err = s.mailService.Send(to, data.Subject, templateName, data)
	if err != nil {
		return err
	}

	return nil
}

// cacheRefreshToken stores the refresh token in Redis.
func (s *AuthService) cacheRefreshToken(ctx context.Context, userUUID string, refreshToken string) error {
	key := fmt.Sprintf("RefreshToken:%s", userUUID)
	return s.rc.Raw().Set(ctx, key, refreshToken, s.cfg.JWT.RefreshExpiration).Err()
}

// getCaptchaEmailMeta returns email metadata based on captcha type.
func getCaptchaEmailMeta(t CaptchaType) (AuthMailData, error) {
	data, ok := captchaMetaMap[t]
	if !ok {
		return AuthMailData{}, fmt.Errorf("unknown captcha type: %s", t)
	}
	return data, nil
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
