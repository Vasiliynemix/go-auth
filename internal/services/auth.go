package services

import (
	"fmt"
	"go.uber.org/zap"
	"time"
	"tutorial-auth/internal/config"
	"tutorial-auth/internal/mongodb/models"
	"tutorial-auth/pkg/authToken"
)

var LoginOrPasswordInvalid = fmt.Errorf("login or password invalid")
var RefreshTokenExpired = fmt.Errorf("refresh token expired")

type AuthService struct {
	cfg         *config.AppConfig
	logger      *zap.Logger
	userService *UserService
}

type AuthResult struct {
	Token        string       `json:"authToken"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
	Err          error        `json:"error"`
}

func NewAuthService(cfg *config.AppConfig, logger *zap.Logger, userService *UserService) *AuthService {
	return &AuthService{
		cfg:         cfg,
		logger:      logger,
		userService: userService,
	}
}

func (as *AuthService) Login(login string, password string) *AuthResult {
	user, err := as.userService.GetByLogin(login)
	if err != nil {
		return &AuthResult{Err: err}
	}

	userPassword, err := as.userService.GetPassword(user.GUID)
	if err != nil {
		return &AuthResult{Err: err}
	}

	valid := as.userService.CheckPasswordHash(password, userPassword)
	if !valid {
		return &AuthResult{Err: LoginOrPasswordInvalid}
	}

	token, refreshToken, err := as.generateTokens(user)
	if err != nil {
		return &AuthResult{Err: err}
	}

	user.LastLoginAt = time.Now()
	return &AuthResult{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}
}

func (as *AuthService) Refresh(guid string, rt string) *AuthResult {
	user, err := as.userService.GetByGuid(guid)
	if err != nil {
		return &AuthResult{Err: err}
	}

	_, valid := authToken.VerifyToken(as.cfg.TokenSecret, rt)
	if !valid {
		return &AuthResult{Err: RefreshTokenExpired}
	}

	token, refreshToken, err := as.generateTokens(user)
	if err != nil {
		return &AuthResult{Err: err}
	}

	user.LastLoginAt = time.Now()
	return &AuthResult{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}
}

func (as *AuthService) generateTokens(user *models.User) (string, string, error) {
	token, err := authToken.NewToken(as.cfg.TokenSecret, as.cfg.TokenExpirationTimeMinutes, &authToken.UserTokenInfo{
		ID:    user.GUID,
		Login: user.Login,
		Name:  user.Name,
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := authToken.NewToken(as.cfg.TokenSecret, as.cfg.RefreshTokenExpirationTimeMinutes, nil)
	if err != nil {
		return "", "", err
	}

	if err = as.userService.UpdateRefreshTokenAndLastLoginAt(user.GUID, refreshToken); err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}
