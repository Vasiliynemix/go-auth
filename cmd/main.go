package main

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"os/signal"
	"time"
	"tutorial-auth/internal/config"
	"tutorial-auth/internal/database"
	"tutorial-auth/internal/mongodb"
	"tutorial-auth/internal/server"
	"tutorial-auth/internal/server/controllers"
	"tutorial-auth/internal/services"
	"tutorial-auth/pkg/logging"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	const op = "cmd.main.main"
	cfg := config.InitConfiguration()

	logger := logging.NewLogger(&cfg.Logging, "auth.log")
	logger.Debug("run this configuration", zap.Any("config", cfg), zap.String("op", op))

	mongoClient := mongodb.NewMongoDB(logger, &cfg.Mongo)
	defer mongoClient.Disconnect()
	err := mongoClient.Connect()
	if err != nil {
		logger.Fatal("failed to connect to mongodb", zap.String("op", op), zap.Error(err))
	}

	for {
		if mongoClient.IsConnected() {
			logger.Info("connected to mongodb", zap.String("op", op))
			break
		}
		time.Sleep(1 * time.Second)
	}

	db, err := database.NewConnectionDB(&cfg.Db)
	defer db.Close()
	if err != nil {
		logger.Fatal("failed to connect to postgresql", zap.String("op", op), zap.Error(err))
	}
	logger.Info("connected to postgresql", zap.String("op", op))

	err = database.ApplyMigration(logger, cfg.Db.Type, db)
	if err != nil {
		logger.Fatal("failed to apply migrations", zap.String("op", op), zap.Error(err))
	}
	logger.Info("applied migrations", zap.String("op", op))

	wApp := server.NewWebServer(logger, &cfg.Web)
	registerRoutes(cfg.App, logger, mongoClient, db, wApp)
	go wApp.Run(cfg.App, logger, mongoClient, db)

	// Graceful shutdown
	SigCh := make(chan os.Signal)
	signal.Notify(SigCh, os.Interrupt, os.Kill)
	takeSig := <-SigCh
	logger.Info("received signal", zap.String("signal", takeSig.String()))
}

func registerRoutes(cfg *config.AppConfig, logger *zap.Logger, mongoClient *mongodb.MongoDB, db *sqlx.DB, wApp *server.WebServer) {
	const op = "cmd.main.registerRoutes"
	logger.Info("registering routes", zap.String("op", op))

	userService := services.NewUserService(logger, mongoClient, db)

	wApp.RegisterRoutes([]controllers.GroupController{
		controllers.NewAuthController(cfg, logger, userService),
		controllers.NewRegisterController(cfg, logger, userService),
	})
}
