// Package service provides business logic layer for the blog system
package service

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"time"

	"blog-server/cache"
	"blog-server/config"
	"blog-server/database"
	"blog-server/entity"
	"blog-server/errs"
	"blog-server/repository"
	"blog-server/request"
	"blog-server/response"
	"blog-server/utils"

	"github.com/google/uuid"
)

// AuthMailData represents the data required to send a captcha email
type AuthMailData struct {
	Title   string
	Type    string
	Content string
	Subject string
	Captcha string
}

// CaptchaType represents different types of captcha operations
type CaptchaType string

const (
	// Register is for user registration captcha
	Register CaptchaType = "Register"
	// PasswordReset is for password reset captcha
	PasswordReset CaptchaType = "PasswordReset"
	// ChangeEmail is for email change captcha
	ChangeEmail CaptchaType = "ChangeEmail"
)

// IAuthService defines the interface for authentication services
type IAuthService interface {
	SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error
	Register(ctx context.Context, dto *request.RegisterReq) (*response.LoginRes, error)
	Login(ctx context.Context, dto *request.LoginReq) (*response.LoginRes, error)
	HasRole(ctx context.Context, uuid string, roles ...entity.UserRole) (bool, error)
	RefreshAccessToken(context context.Context, token string) (*response.RefreshRes, error)
}

// AuthService implements the IAuthService interface
type AuthService struct {
	db          database.DB
	rc          cache.CacheClient
	cfg         *config.Config
	userRepo    repository.IUserRepo
	jwtService  IJwtService
	mailService IMailService
}

// NewAuthService creates and returns a new AuthService instance
func NewAuthService(db database.DB, rc cache.CacheClient, userRepo repository.IUserRepo, jwtService IJwtService, mailService IMailService) IAuthService {
	return &AuthService{
		db:          db,
		rc:          rc,
		cfg:         config.Get(),
		userRepo:    userRepo,
		jwtService:  jwtService,
		mailService: mailService,
	}
}

// Register registers a new user and generates access/refresh tokens
func (s *AuthService) Register(ctx context.Context, dto *request.RegisterReq) (*response.LoginRes, error) {
	email := dto.Email

	// Verify captcha
	cachedCaptcha, err := s.rc.Get(ctx, fmt.Sprintf("Register:%s", email))
	if !strings.EqualFold(strings.TrimSpace(cachedCaptcha), strings.TrimSpace(dto.Captcha)) {
		return nil, errs.New(errs.CodeInvalidParam, "Invalid captcha", nil)
	}

	hashPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		return nil, errs.New(errs.CodeInternalError, "Hash password failed", err)
	}

	user := &entity.User{
		UUID:     uuid.New(),
		Email:    dto.Email,
		Password: hashPassword,
		Username: entity.GenerateUsername(),
	}

	var newUser *entity.User

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

	accessToken, refreshToken, err := s.jwtService.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return nil, err
	}

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

// SendCaptchaMail generates a captcha, stores it in Redis, and sends an email
func (s *AuthService) SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error {
	userExists, err := s.userRepo.ExistsByEmail(ctx, s.db.Conn(), to)
	if err != nil {
		return err
	}
	if userExists {
		return errs.New(errs.CodeConflict, fmt.Sprintf("Sent captcha failed to %s: user exists", to), nil)
	}

	if captchaType == "" {
		captchaType = Register
	}

	data, err := getCaptchaEmailMeta(captchaType)
	if err != nil {
		return err
	}

	captcha := generateCaptcha()
	data.Captcha = captcha

	if err := s.rc.Set(ctx, fmt.Sprintf("%s:%s", captchaType, to), captcha, 5*time.Minute); err != nil {
		return errs.New(errs.CodeCacheError, "Set captcha in Redis failed", err)
	}

	templateName := "captcha.html"
	if err := s.mailService.Send(to, data.Subject, templateName, data); err != nil {
		return err
	}

	return nil
}

// RefreshAccessToken refreshes the access token using a valid refresh token
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

	cacheRefreshToken, err := s.rc.Get(ctx, fmt.Sprintf("RefreshToken:%s", uuid))
	if err != nil {
		return nil, errs.New(errs.CodeCacheError, "Get refresh token from Redis failed", err)
	}

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

// HasRole checks if a user has any of the specified roles
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

// cacheRefreshToken stores the refresh token in Redis
func (s *AuthService) cacheRefreshToken(ctx context.Context, userUUID string, refreshToken string) error {
	key := fmt.Sprintf("RefreshToken:%s", userUUID)
	if err := s.rc.Set(ctx, key, refreshToken, s.cfg.JWT.RefreshExpiration); err != nil {
		return errs.New(
			errs.CodeCacheError,
			fmt.Sprintf("Failed to cache refresh token for user %s", userUUID),
			err,
		)
	}
	return nil
}

// getCaptchaEmailMeta returns email metadata based on captcha type
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
		Subject: "[Immortal's Blog] Email Verification",
		Title:   "Verify Your Registration",
		Type:    "Registration",
		Content: "Please use the following verification code to complete your registration.",
	},
	PasswordReset: {
		Subject: "[Immortal's Blog] Password Reset",
		Title:   "Verify Your Password Reset Request",
		Type:    "Password Reset",
		Content: "You are attempting to reset your password. Please use this verification code to proceed.",
	},
	ChangeEmail: {
		Subject: "[Immortal's Blog] Change Email Address",
		Title:   "Verify Your New Email Address",
		Type:    "Email Change",
		Content: "You are attempting to change your email address. Please use this verification code to verify your new email.",
	},
}

func generateCaptcha() string {
	return utils.RandomString(6, "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}
