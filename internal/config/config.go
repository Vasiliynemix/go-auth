package config

import (
	"fmt"
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
	Username string
	Password string
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
	Username string
	Password string
}

type Config struct {
	Mongo MongoDbConnectionConfig `mapstructure:"mongo"`
	Web   WebServerConfig         `mapstructure:"web"`
	Db    DBConnectionConfig      `mapstructure:"db"`
}

func InitConfiguration() *Config {
	var C = new(Config)

	LoadDefault()
	LoadFile()

	err := viper.Unmarshal(C)
	if err != nil {
		panic(err)
	}
	return C
}

func LoadDefault() {
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

func LoadFile() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(fmt.Sprintf("failed to read config file, %s", err))
	}
}
