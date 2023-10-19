package controllers

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"tutorial-auth/internal/config"
	"tutorial-auth/internal/mongodb/models"
	"tutorial-auth/internal/services"
)

var validate = validator.New()

var LoginRequiredError = fmt.Errorf("login required")
var PasswordRequiredError = fmt.Errorf("password required")

type ValidationError struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

type RefreshAuthRequest struct {
	ID           string `json:"id" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (ra *RefreshAuthRequest) Validate() []*ValidationError {
	var validationErrors []*ValidationError

	errs := validate.Struct(ra)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ValidationError

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, &elem)
		}
	}

	return validationErrors
}

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
	OK     bool               `json:"ok"`
	Cause  string             `json:"cause"`
	Errors []*ValidationError `json:"errors,omitempty"`
}

type AuthResponseOK struct {
	OK           bool         `json:"ok"`
	Token        string       `json:"authToken,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	User         *models.User `json:"user,omitempty"`
}

type AuthController struct {
	cfg         *config.AppConfig
	logger      *zap.Logger
	authService *services.AuthService
}

func NewAuthController(cfg *config.AppConfig, logger *zap.Logger, authService *services.AuthService) *AuthController {
	return &AuthController{
		cfg:         cfg,
		logger:      logger,
		authService: authService,
	}
}

func (c *AuthController) GetGroup() string {
	return "/auth"
}

func (c *AuthController) GetHandlers() []ControllerHandler {
	return []ControllerHandler{
		&Handler{
			Method: "POST", Path: "/login",
			Handler: c.authHandler(),
		},
		&Handler{
			Method: "POST", Path: "/refresh",
			Handler: c.refreshHandler(),
		},
	}
}

func (c *AuthController) authHandler() func(fc *fiber.Ctx) error {
	const op = "internal.server.controllers.auth.authHandler"

	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var req AuthRequest
		if err := fc.BodyParser(&req); err != nil {
			return err
		}
		if valid, err := (*AuthRequest).Validate(&req); !valid {
			c.logger.Error("Validation error", zap.String("op", op), zap.String("error", err.Error()))
			return fc.JSON(AuthResponseError{
				OK:    false,
				Cause: err.Error(),
			})
		}

		authResult := c.authService.Login(req.Login, req.Password)
		if authResult.Err != nil {
			c.logger.Error("Login error", zap.String("op", op), zap.String("error", authResult.Err.Error()))
			return fc.JSON(AuthResponseError{
				OK:    false,
				Cause: authResult.Err.Error(),
			})
		}

		return fc.JSON(AuthResponseOK{
			OK:           true,
			Token:        authResult.Token,
			RefreshToken: authResult.RefreshToken,
			User:         authResult.User,
		})
	}
}

func (c *AuthController) refreshHandler() func(fc *fiber.Ctx) error {
	const op = "internal.server.controllers.auth.refreshHandler"

	return func(fc *fiber.Ctx) error {
		fc.Accepts("application/json")
		var req RefreshAuthRequest
		if err := fc.BodyParser(&req); err != nil {
			return err
		}

		validationErrors := (*RefreshAuthRequest).Validate(&req)
		if len(validationErrors) > 0 {
			for _, validationError := range validationErrors {
				c.logger.Error("Validation error", zap.String("op", op), zap.String("error", validationError.FailedField))
			}
			return fc.JSON(AuthResponseError{
				OK:     false,
				Cause:  "validation errors",
				Errors: validationErrors,
			})
		}

		authResult := c.authService.Refresh(req.ID, req.RefreshToken)
		if authResult.Err != nil {
			c.logger.Error("Refresh error", zap.String("op", op), zap.String("error", authResult.Err.Error()))
			return fc.JSON(AuthResponseError{
				OK:    false,
				Cause: authResult.Err.Error(),
			})
		}

		return fc.JSON(AuthResponseOK{
			OK:           true,
			Token:        authResult.Token,
			RefreshToken: authResult.RefreshToken,
			User:         authResult.User,
		})
	}
}
