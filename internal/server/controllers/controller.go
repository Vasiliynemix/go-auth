package controllers

import "github.com/gofiber/fiber/v2"

type GroupController interface {
	GetGroup() string
	GetHandlers() []ControllerHandler
}

type ControllerHandler interface {
	GetMethod() string
	GetHandler() func(c *fiber.Ctx) error
	GetPath() string
}

type Handler struct {
	Method  string
	Path    string
	Handler func(c *fiber.Ctx) error
}

func (h *Handler) GetPath() string {
	return h.Path
}

func (h *Handler) GetHandler() func(c *fiber.Ctx) error {
	return h.Handler
}

func (h *Handler) GetMethod() string {
	return h.Method
}
