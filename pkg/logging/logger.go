package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"tutorial-auth/internal/config"
)

func NewLogger(cfg *config.LoggingConfig, name string) *zap.Logger {
	_ = os.Mkdir(cfg.Path, os.ModePerm)

	logger := zap.New(configure(cfg, name))
	return logger
}

func configure(cfg *config.LoggingConfig, name string) zapcore.Core {
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, name),
		MaxSize:    cfg.MaxSize, // megabytes
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge, // days
	})

	var priority zap.LevelEnablerFunc

	switch cfg.Level {
	case "debug":
		priority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zap.DebugLevel
		})
	case "warn":
		priority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zap.WarnLevel
		})
	default:
		priority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zap.InfoLevel
		})
	}

	consoleWriter := zapcore.Lock(os.Stdout)
	jsonEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, priority),
		zapcore.NewCore(jsonEncoder, fileWriter, priority),
	)
}
