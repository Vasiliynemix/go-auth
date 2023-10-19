package server

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"reflect"
	"time"
	"tutorial-auth/internal/config"
	"tutorial-auth/internal/mongodb"
	"tutorial-auth/internal/server/controllers"
	"tutorial-auth/internal/services"
)

type WebServer struct {
	log    *zap.Logger
	cfg    *config.WebServerConfig
	client *fiber.App
}

func NewWebServer(logger *zap.Logger, cfg *config.WebServerConfig) *WebServer {
	return &WebServer{
		log: logger,
		cfg: cfg,
		client: fiber.New(fiber.Config{
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  60 * time.Second,
			AppName:      "My App v1.0.0",
		}),
	}
}

func (ws *WebServer) RegisterRoutes(routes []controllers.GroupController) {
	for _, route := range routes {
		group := ws.client.Group(route.GetGroup())
		for _, handler := range route.GetHandlers() {
			switch handler.GetMethod() {
			case "GET":
				group.Get(handler.GetPath(), handler.GetHandler())
			case "POST":
				group.Post(handler.GetPath(), handler.GetHandler())
			case "PUT":
				group.Put(handler.GetPath(), handler.GetHandler())
			case "PATCH":
				group.Patch(handler.GetPath(), handler.GetHandler())
			case "DELETE":
				group.Delete(handler.GetPath(), handler.GetHandler())
			default:
				ws.log.Error(
					"unsupported HTTP method",
					zap.String("controller", reflect.TypeOf(route).Elem().Name()),
					zap.String("path", handler.GetPath()),
					zap.String("method", handler.GetMethod()),
				)
			}
		}
	}
}

func (ws *WebServer) Run(cfg *config.AppConfig, logger *zap.Logger, mongoClient *mongodb.MongoDB, db *sqlx.DB) {
	ws.RegisterRoutes([]controllers.GroupController{
		controllers.NewAuthController(
			cfg,
			logger,
			services.NewAuthService(cfg, logger, services.NewUserService(logger, mongoClient, db)),
		),
	})
	ws.client.Listen(fmt.Sprintf(":%d", ws.cfg.Port))
}
