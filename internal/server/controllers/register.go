package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"tutorial-auth/internal/mongodb/models"
	"tutorial-auth/internal/services"
)

var PasswordNotEqualError = fmt.Errorf("password not equal")

type RegisterRequest struct {
	Login           string `json:"login"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (r *RegisterRequest) Validate() (bool, error) {
	if len(r.Login) == 0 {
		return false, LoginRequiredError
	}

	if len(r.Password) == 0 {
		return false, PasswordRequiredError
	}

	if len(r.ConfirmPassword) == 0 {
		return false, PasswordRequiredError
	}

	if r.Password != r.ConfirmPassword {
		return false, PasswordNotEqualError
	}

	return true, nil
}

type RegisterResponseError struct {
	OK    bool   `json:"ok"`
	Cause string `json:"cause"`
}

type RegisterResponseOK struct {
	OK   bool         `json:"ok"`
	User *models.User `json:"user"`
}

type RegisterController struct {
	logger      *zap.Logger
	userService *services.UserService
}

func NewRegisterController(logger *zap.Logger, userService *services.UserService) *RegisterController {
	return &RegisterController{
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
		if valid, err := (*RegisterRequest).Validate(&req); !valid {
			return fc.JSON(RegisterResponseError{
				OK:    false,
				Cause: err.Error(),
			})
		}

		user, err := c.userService.Register(req.Login, req.Password)
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
