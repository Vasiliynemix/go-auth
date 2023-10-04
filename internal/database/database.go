package database

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"tutorial-auth/internal/config"
)

func NewConnectionDB(cfg *config.DBConnectionConfig) (*sqlx.DB, error) {
	var db *sqlx.DB
	var err error

	switch cfg.Type {
	case "postgres":
		db, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database))
	case "mysql":
		db, err = sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	default:
		db, err = nil, fmt.Errorf("unsupported db type: %s", cfg.Type)
	}

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(cfg.Pool.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Pool.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.Pool.IdleTimeout)

	return db, nil
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func ApplyMigration(logger *zap.Logger, DBType string, db *sqlx.DB) error {
	const op = "database.ApplyMigration"

	goose.SetBaseFS(embedMigrations)

	var dialect string
	switch DBType {
	case "postgres":
		dialect = "postgres"
	case "mysql":
		dialect = "mysql"
	default:
		logger.Panic("unsupported db type", zap.String("op", op), zap.String("dbType", DBType))
	}

	if err := goose.SetDialect(dialect); err != nil {
		return err
	}

	if err := goose.Up(db.DB, "migrations"); err != nil {
		return err
	}
	return nil
}
