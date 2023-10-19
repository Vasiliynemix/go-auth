package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"time"
)

type WebServerConfig struct {
	Port int
}

type DBConnectionConfig struct {
	Type     string
	Host     string
	Port     int
	Database string
	Username string `json:"-"`
	Password string `json:"-"`
	Pool     struct {
		MaxIdleConns int           `mapstructure:"max_idle_conns"`
		MaxOpenConns int           `mapstructure:"max_open_conns"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	}
}

type MongoDbConnectionConfig struct {
	Host     string
	Port     int
	Database string
	Username string `json:"-"`
	Password string `json:"-"`
}

type AppConfig struct {
	Name                              string
	PasswordLifeTime                  int    `mapstructure:"password_life_time"` // in hours
	PasswordMinLength                 int    `mapstructure:"password_min_length"`
	TokenExpirationTimeMinutes        int    `mapstructure:"token_expiration_time_minutes"`         // in minutes
	RefreshTokenExpirationTimeMinutes int    `mapstructure:"refresh_token_expiration_time_minutes"` // in minutes
	TokenSecret                       string `mapstructure:"token_secret"`
}

type LoggingConfig struct {
	Level      string
	Path       string
	MaxSize    int `mapstructure:"max_size"`
	MaxBackups int `mapstructure:"max_backups"`
	MaxAge     int `mapstructure:"max_age"`
}

type Config struct {
	Debug   bool                    `mapstructure:"debug"`
	App     *AppConfig              `mapstructure:"app"`
	Logging LoggingConfig           `mapstructure:"logging"`
	Mongo   MongoDbConnectionConfig `mapstructure:"mongo"`
	Web     WebServerConfig         `mapstructure:"web"`
	Db      DBConnectionConfig      `mapstructure:"db"`
}

var C = new(Config)

func InitConfiguration() *Config {
	initConfig()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		initConfig()
	})
	return C
}

func initConfig() {
	LoadDefault()
	LoadFile()
	LoadEnv()

	if viper.GetBool("debug") {
		viper.SetDefault("logging.level", "debug")
	}

	viper.Unmarshal(C)
}

func LoadDefault() {
	viper.SetDefault("debug", false)
	viper.SetDefault("app.name", "tutorial-auth")
	viper.SetDefault("app.password_life_time", 1)
	viper.SetDefault("app.password_min_length", 8)
	viper.SetDefault("app.token_expiration_time_minutes", 5)
	viper.SetDefault("app.refresh_token_expiration_time_hours", 60)
	viper.SetDefault("app.token_secret", "secret")

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.path", "logs")
	viper.SetDefault("logging.max_size", 500)
	viper.SetDefault("logging.max_backups", 3)
	viper.SetDefault("logging.max_age", 30)

	viper.SetDefault("mongo.host", "localhost")
	viper.SetDefault("mongo.port", 27018)
	viper.SetDefault("mongo.database", "mongo-auth-db")
	viper.SetDefault("mongo.username", "root")
	viper.SetDefault("mongo.password", "root")

	viper.SetDefault("web.port", 8080)

	viper.SetDefault("db.type", "postgres")
	viper.SetDefault("db.host", "localhost")
	viper.SetDefault("db.port", 5433)
	viper.SetDefault("db.database", "postgres-auth-db")
	viper.SetDefault("db.username", "user")
	viper.SetDefault("db.password", "password")
	viper.SetDefault("db.pool.max_idle_conns", 1)
	viper.SetDefault("db.pool.max_open_conns", 10)
	viper.SetDefault("db.pool.idle_timeout", 300*time.Second)
}

func LoadEnv() {
	viper.SetEnvPrefix("auth")
	viper.AutomaticEnv()
}

func LoadFile() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./config/")
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", viper.GetString("app.name")))
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", viper.GetString("app.name")))
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/config/", viper.GetString("app.name")))

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to read config file, %s. Use default values.", err))
	}
}
