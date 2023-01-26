package config

import (
	"github.com/kanji-team/auth/internal/app"
	"github.com/kanji-team/auth/internal/database"
	"github.com/kanji-team/auth/internal/grpcserver"
	"github.com/kanji-team/auth/internal/user"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Logger     *LoggerConfig
	GRPCServer *grpcserver.Config
	Auth       *app.Config
	DB         *database.Config
	User       *user.Config
}

type LoggerConfig struct {
	LogLevel string `default:"debug"`
}

func Parse() (*Config, error) {
	cfg := Config{}
	auth := app.Config{}
	logger := LoggerConfig{}
	db := database.Config{}
	grpc := grpcserver.Config{}
	userConfig := user.Config{}
	project := "KANJI_AUTH"

	err := envconfig.Process(project, &logger)
	if err != nil {
		log.Err(err).Msg("logger config error")
		return nil, err
	}

	err = envconfig.Process(project, &db)
	if err != nil {
		log.Err(err).Msg("db config error")
		return nil, err
	}

	err = envconfig.Process(project, &auth)
	if err != nil {
		log.Err(err).Msg("auth config error")
		return nil, err
	}

	err = envconfig.Process(project, &grpc)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	err = envconfig.Process(project, &userConfig)
	if err != nil {
		log.Err(err).Msg("http config error")
		return nil, err
	}

	cfg.Auth = &auth
	cfg.DB = &db
	cfg.Logger = &logger
	cfg.GRPCServer = &grpc
	cfg.User = &userConfig

	return &cfg, nil
}
