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
	"blog-server/datastore"
	"blog-server/entity"
	"blog-server/pkg/errx"
	"blog-server/pkg/jwt"
	"blog-server/repository"
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
type AuthService interface {
	SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error
	Register(ctx context.Context, email, password, captcha string) (string, string, error)
	Login(ctx context.Context, email, password string) (string, string, error)
	HasRole(ctx context.Context, uuid string, roles ...entity.UserRole) (bool, error)
	RefreshAccessToken(context.Context, string) (string, string, error)
}

// AuthService implements the IAuthService interface
type authService struct {
	ds          *datastore.DataStore
	rc          cache.CacheClient
	cfg         *config.Config
	userRepo    repository.UserRepo
	mailService IMailService
}

// NewAuthService creates and returns a new AuthService instance
func NewAuthService(ds *datastore.DataStore, rc cache.CacheClient, userRepo repository.UserRepo, mailService IMailService) AuthService {
	return &authService{
		ds:          ds,
		rc:          rc,
		cfg:         config.Get(),
		userRepo:    userRepo,
		mailService: mailService,
	}
}

// Register registers a new user and generates access/refresh tokens
func (s *authService) Register(ctx context.Context, email, password, captcha string) (string, string, error) {
	// Verify captcha
	cachedCaptcha, err := s.rc.Get(ctx, fmt.Sprintf("Register:%s", email))
	if !strings.EqualFold(strings.TrimSpace(cachedCaptcha), strings.TrimSpace(captcha)) {
		return "", "", errx.New(errx.CodeInvalidParam, "Invalid captcha", nil)
	}

	hashPassword, err := utils.HashPassword(password)
	if err != nil {
		return "", "", errx.New(errx.CodeInternalError, "Hash password failed", err)
	}

	user := &entity.User{
		UUID:     uuid.New(),
		Email:    email,
		Username: entity.GenerateUsername(),
	}

	err = s.ds.WithTx(ctx, func(ctx context.Context) error {
		_, err = s.userRepo.Create(ctx, user, hashPassword)
		if err != nil {
			return errx.New(errx.CodeDatabaseError, "Create user failed", err)
		}
		return nil
	})
	if err != nil {
		return "", "", errx.New(errx.CodeInternalError, "Register transaction failed", err)
	}

	j := jwt.New(s.cfg.JWT)

	accessToken, refreshToken, err := j.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return "", "", err
	}

	if err := s.cacheRefreshToken(ctx, user.UUID.String(), refreshToken); err != nil {
		return "", "", errx.New(errx.CodeCacheError, "Cache refresh token failed", err)
	}

	return accessToken, refreshToken, nil
}

// Login logs in a user and generates access/refresh tokens
func (s *authService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetAuthByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	if !utils.VerifyPassword(password, user.Password) {
		return "", "", errx.New(errx.CodeInvalidParam, "Invalid password", nil)
	}

	j := jwt.New(s.cfg.JWT)

	accessToken, refreshToken, err := j.GenerateAllTokens(user.UUID.String())
	if err != nil {
		return "", "", err
	}

	if err := s.cacheRefreshToken(ctx, user.UUID.String(), refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// SendCaptchaMail generates a captcha, stores it in Redis, and sends an email
func (s *authService) SendCaptchaMail(ctx context.Context, to string, captchaType CaptchaType) error {
	exists, err := s.userRepo.ExistsByEmail(ctx, to)
	if err != nil {
		return err
	}
	if exists {
		return errx.New(errx.CodeConflict, fmt.Sprintf("Sent captcha failed to %s: user exists", to), nil)
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
		return errx.New(errx.CodeCacheError, "Set captcha in Redis failed", err)
	}

	templateName := "captcha.html"
	if err := s.mailService.Send(to, data.Subject, templateName, data); err != nil {
		return err
	}

	return nil
}

// RefreshAccessToken refreshes the access token using a valid refresh token
func (s *authService) RefreshAccessToken(ctx context.Context, token string) (string, string, error) {
	j := jwt.New(s.cfg.JWT)
	claims, err := j.Parse(token)
	if err != nil {
		return "", "", err
	}

	uuid := claims.UserID

	existed, err := s.userRepo.ExistsByUUID(ctx, uuid)
	if err != nil {
		return "", "", err
	}
	if !existed {
		return "", "", errx.New(errx.CodeUserNotFound, "User not found", nil)
	}

	cacheRefreshToken, err := s.rc.Get(ctx, fmt.Sprintf("RefreshToken:%s", uuid))
	if err != nil {
		return "", "", errx.New(errx.CodeCacheError, "Get refresh token from Redis failed", err)
	}

	if !strings.EqualFold(token, cacheRefreshToken) {
		return "", "", errx.New(errx.CodeUnauthorized, "Invalid token", nil)
	}

	accessToken, refreshToken, err := j.GenerateAllTokens(uuid)
	if err != nil {
		return "", "", err
	}

	if err := s.cacheRefreshToken(ctx, uuid, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// HasRole checks if a user has any of the specified roles
func (s *authService) HasRole(ctx context.Context, uuid string, roles ...entity.UserRole) (bool, error) {
	user, err := s.userRepo.GetAuthByUUID(ctx, uuid)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, errx.New(errx.CodeUserNotFound, "User not found", nil)
	}

	if slices.Contains(roles, user.Role) {
		return true, nil
	}
	return false, nil
}

// cacheRefreshToken stores the refresh token in Redis
func (s *authService) cacheRefreshToken(ctx context.Context, userUUID string, refreshToken string) error {
	key := fmt.Sprintf("RefreshToken:%s", userUUID)
	if err := s.rc.Set(ctx, key, refreshToken, s.cfg.JWT.RefreshExpiration); err != nil {
		return errx.New(
			errx.CodeCacheError,
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
		return nil, errx.New(
			errx.CodeInvalidParam,
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
