// Package service
package service

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"blog-server/internal/cache"
	"blog-server/internal/config"
	"blog-server/internal/database"
	"blog-server/internal/entity"
	"blog-server/internal/repo"
	"blog-server/internal/request"
	"blog-server/internal/response"
	"blog-server/pkg/errs"
	"blog-server/pkg/utils"

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
	HasRole(ctx context.Context, uuid string, roles ...entity.UserRole) (bool, error)
	RefreshAccessToken(context context.Context, token string) (*response.RefreshRes, error)
}

// AuthService implements the IAuthService interface.
type AuthService struct {
	db          database.DB
	rc          cache.CacheClient
	cfg         *config.Config
	userRepo    repo.IUserRepo
	jwtService  IJwtService
	mailService IMailService
}

// NewAuthService creates and returns a new AuthService instance
func NewAuthService(db database.DB, rc cache.CacheClient, userRepo repo.IUserRepo, jwtService IJwtService, mailService IMailService) IAuthService {
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

	// ===== 验证码检查 =====
	cachedCaptcha := s.rc.Raw().Get(ctx, fmt.Sprintf("Register:%s", email)).Val()
	if !strings.EqualFold(strings.TrimSpace(cachedCaptcha), strings.TrimSpace(dto.Captcha)) {
		return nil, errs.New(errs.CodeInvalidParam, "Invalid captcha", nil)
	}

	// ===== 密码加密 =====
	hashPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, errs.New(errs.CodeInternalError, "Hash password failed", err)
	}

	user := &entity.User{
		UUID:     uuid.New(),
		Email:    dto.Email,
		Password: hashPassword,
		Username: utils.GenerateUsername(),
	}

	var newUser *entity.User

	// ===== 创建用户事务 =====
	err = s.db.Trans(func(txCtx *database.TxContext) error {
		newUser, err = s.userRepo.CreateUser(ctx, txCtx.GetTx(), user)
		if err != nil {
			return errs.New(errs.CodeDatabaseError, "Create user failed", err)
		}
		return nil
	})
	if err != nil {
		return nil, errs.New(errs.CodeInternalError, "Register transaction failed", err)
	}

	// ===== 生成 JWT Token =====
	accessToken, refreshToken, err := s.jwtService.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return nil, err
	}

	// ===== 缓存 Refresh Token =====
	if err := s.cacheRefreshToken(ctx, user.UUID.String(), refreshToken); err != nil {
		return nil, errs.New(errs.CodeCacheError, "Cache refresh token failed", err)
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
	if !utils.VerifyPassword(dto.Password, user.Password) {
		return nil, errs.New(errs.CodeInvalidParam, "Invalid password", nil)
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
		UUID:         user.UUID.String(),
		Username:     user.Username,
		Avatar:       user.Avatar,
		Role:         user.Role,
	}

	return result, nil
}

// SendCaptchaMail generates a captcha, stores it in Redis, and sends an email.
func (s *AuthService) SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error {
	// ===== 检查用户是否存在 =====
	userExists, err := s.userRepo.ExistsByEmail(ctx, s.db.Conn(), to)
	if err != nil {
		return err
	}
	if userExists {
		return errs.New(errs.CodeConflict, fmt.Sprintf("Sent captcha failed to %s: user exists", to), nil)
	}

	// 默认类型
	if captchaType == "" {
		captchaType = Register
	}

	// ===== 获取邮件模板数据 =====
	data, err := getCaptchaEmailMeta(captchaType)
	if err != nil {
		return err
	}

	// ===== 生成验证码并缓存 =====
	captcha := utils.GenerateCaptcha()
	data.Captcha = captcha

	if err := s.rc.Raw().Set(ctx, fmt.Sprintf("%s:%s", captchaType, to), captcha, 5*time.Minute).Err(); err != nil {
		return errs.New(errs.CodeCacheError, "Set captcha in Redis failed", err)
	}

	// ===== 发送邮件 =====
	templateName := "captcha.html"
	if err := s.mailService.Send(to, data.Subject, templateName, data); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) RefreshAccessToken(ctx context.Context, token string) (*response.RefreshRes, error) {
	claims, err := s.jwtService.ParseToken(token)
	if err != nil {
		return nil, err
	}
	uuid := claims.UserID

	existed, err := s.userRepo.ExistsByUUID(ctx, s.db.Conn(), uuid)
	if err != nil {
		return nil, err
	}
	if !existed {
		return nil, errs.New(errs.CodeUserNotFound, "User not found", nil)
	}

	cacheRefreshToken := s.rc.Raw().Get(ctx, fmt.Sprintf("RefreshToken:%s", uuid)).Val()

	if !strings.EqualFold(token, cacheRefreshToken) {
		return nil, errs.New(errs.CodeUnauthorized, "Invalid token", nil)
	}

	accessToken, refreshToken, err := s.jwtService.GenerateAllTokens(uuid)
	if err != nil {
		return nil, err
	}

	if err := s.cacheRefreshToken(ctx, uuid, refreshToken); err != nil {
		return nil, err
	}

	return &response.RefreshRes{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) HasRole(ctx context.Context, uuid string, roles ...entity.UserRole) (bool, error) {
	userRole, err := s.userRepo.GetRoleByUUID(ctx, s.db.Conn(), uuid)
	if err != nil {
		return false, err
	}

	if userRole == nil {
		return false, errs.New(errs.CodeUserNotFound, "User not found", nil)
	}

	if slices.Contains(roles, *userRole) {
		return true, nil
	}
	return false, nil
}

// cacheRefreshToken stores the refresh token in Redis.
func (s *AuthService) cacheRefreshToken(ctx context.Context, userUUID string, refreshToken string) error {
	key := fmt.Sprintf("RefreshToken:%s", userUUID)
	if err := s.rc.Raw().Set(ctx, key, refreshToken, s.cfg.JWT.RefreshExpiration).Err(); err != nil {
		return errs.New(
			errs.CodeCacheError,
			fmt.Sprintf("Failed to cache refresh token for user %s", userUUID),
			err,
		)
	}
	return nil
}

// getCaptchaEmailMeta returns email metadata based on captcha type.
func getCaptchaEmailMeta(t CaptchaType) (*AuthMailData, error) {
	data, ok := captchaMetaMap[t]
	if !ok {
		return nil, errs.New(
			errs.CodeInvalidParam,
			fmt.Sprintf("Unknown captcha type: %s", t),
			nil,
		)
	}
	copyData := *data
	return &copyData, nil
}

var captchaMetaMap = map[CaptchaType]*AuthMailData{
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
