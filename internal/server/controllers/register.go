package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"tutorial-auth/internal/config"
	"tutorial-auth/internal/mongodb/models"
	"tutorial-auth/internal/services"
)

var PasswordNotEqualError = fmt.Errorf("password not equal")

type RegisterRequestValidationConfig struct {
	LoginRequired     bool `json:"login_required"`
	PasswordRequired  bool `json:"password_required"`
	PasswordEqual     bool `json:"password_not_equal"`
	PasswordMinLength int  `json:"password_min_length"`
}

type RegisterRequest struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Name            string `json:"name"`
	LastName        string `json:"last_name,omitempty"`
}

func (r *RegisterRequest) Validate(cfg RegisterRequestValidationConfig) (bool, error) {
	if cfg.LoginRequired && len(r.Login) == 0 {
		return false, LoginRequiredError
	}

	if cfg.PasswordRequired {
		if len(r.Password) == 0 {
			return false, PasswordRequiredError
		}

		if len(r.ConfirmPassword) == 0 {
			return false, PasswordRequiredError
		}
	}

	if cfg.PasswordEqual && r.Password != r.ConfirmPassword {
		return false, PasswordNotEqualError
	}

	if cfg.PasswordMinLength > 0 && len(r.Password) < cfg.PasswordMinLength {
		errMsg := fmt.Sprintf("Password is too short. Min length is %d", cfg.PasswordMinLength)
		return false, errors.New(errMsg)
	}

	return true, nil
}

type RegisterResponseError struct {
	OK    bool   `json:"ok"`
	Cause string `json:"cause"`
}

type RegisterResponseOK struct {
	OK   bool         `json:"ok"`
	User *models.User `json:"user,omitempty"`
}

type RegisterController struct {
	cfg         *config.AppConfig
	logger      *zap.Logger
	userService *services.UserService
}

func NewRegisterController(cfg *config.AppConfig, logger *zap.Logger, userService *services.UserService) *RegisterController {
	return &RegisterController{
		cfg:         cfg,
		logger:      logger,
		userService: userService,
	}
}

func (c *RegisterController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method: "POST", Path: "/register",
			Handler: c.RegisterHandler(),
		},
	}
}

func (c *RegisterController) GetGroup() string {
	return "/auth"
}

func (c *RegisterController) RegisterHandler() func(*fiber.Ctx) error {
	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var req RegisterRequest
		if err := fc.BodyParser(&req); err != nil {
			return err
		}
		if valid, err := (*RegisterRequest).Validate(&req, RegisterRequestValidationConfig{
			LoginRequired:     true,
			PasswordRequired:  true,
			PasswordEqual:     true,
			PasswordMinLength: c.cfg.PasswordMinLength,
		}); !valid {
			return fc.JSON(RegisterResponseError{
				OK:    false,
				Cause: err.Error(),
			})
		}

		ctx := context.WithValue(context.Background(), "cfg", c.cfg)

		user, err := c.userService.Register(ctx, &services.NewUser{
			Login:    req.Login,
			Password: req.Password,
			Name:     req.Name,
			LastName: req.LastName,
		})
		if err != nil {
			return fc.JSON(RegisterResponseError{
				OK:    false,
				Cause: err.Error(),
			})
		}

		return fc.JSON(RegisterResponseOK{
			OK:   true,
			User: user,
		})
	}
}
