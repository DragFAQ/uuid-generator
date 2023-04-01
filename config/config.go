package config

import (
	logger "github.com/DragFAQ/uuid-generator/logger"
	"github.com/kelseyhightower/envconfig"
)

type AppParams struct {
	Environment string `envconfig:"APP_ENV" required:"true"`
	Name        string `envconfig:"APP_NAME" required:"true"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
}

type ServerParams struct {
	HttpPort string `envconfig:"HTTP_PORT" default:"8080"`
	GrpcPort string `envconfig:"GRPC_PORT" default:"8090"`
}

type SettingsParams struct {
	HashTTLSeconds int `envconfig:"HASH_TTL_SECONDS" default:"300"`
}

type Config struct {
	App      AppParams
	Server   ServerParams
	Settings SettingsParams
	Logger   logger.Configuration
}

func NewConfig() Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	var isJSONFormat bool
	if cfg.App.Environment != "local" {
		isJSONFormat = true
	}

	return Config{
		App:      cfg.App,
		Server:   cfg.Server,
		Settings: cfg.Settings,
		Logger: logger.Configuration{
			EnableConsole:     true,
			ConsoleJSONFormat: isJSONFormat,
			ConsoleLevel:      cfg.App.LogLevel,
		},
	}
}
