package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"tutorial-auth/internal/config"
	"tutorial-auth/internal/mongodb/models"
	"tutorial-auth/internal/services"
)

var LoginRequiredError = fmt.Errorf("login required")
var PasswordRequiredError = fmt.Errorf("password required")

type AuthRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (r *AuthRequest) Validate() (bool, error) {
	if len(r.Login) == 0 {
		return false, LoginRequiredError
	}

	if len(r.Password) == 0 {
		return false, PasswordRequiredError
	}

	return true, nil
}

type AuthResponseError struct {
	OK    bool   `json:"ok"`
	Cause string `json:"cause"`
}

type AuthResponseOK struct {
	OK           bool        `json:"ok"`
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	User         models.User `json:"user"`
}

type AuthController struct {
	cfg         *config.AppConfig
	logger      *zap.Logger
	userService *services.UserService
}

func NewAuthController(cfg *config.AppConfig, logger *zap.Logger, userService *services.UserService) *AuthController {
	return &AuthController{
		cfg:         cfg,
		logger:      logger,
		userService: userService,
	}
}

func (c *AuthController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method: "POST", Path: "/",
			Handler: c.authHandler(),
		},
	}
}

func (c *AuthController) GetGroup() string {
	return "/auth"
}

func (c *AuthController) authHandler() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Accepts("application/json")
		var req AuthRequest
		if err := c.BodyParser(&req); err != nil {
			return err
		}
		if valid, err := (*AuthRequest).Validate(&req); !valid {
			return c.JSON(AuthResponseError{
				OK:    false,
				Cause: err.Error(),
			})
		}
		return c.JSON(AuthResponseOK{
			OK:           true,
			Token:        "123",
			RefreshToken: "132",
			User: models.User{
				Login: req.Login,
			},
		})
	}
}
