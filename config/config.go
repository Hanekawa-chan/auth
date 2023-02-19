package config

import (
	"auth/internal/database"
	"auth/internal/grpcserver"
	"auth/internal/user"
	"auth/pkg/api"
	"auth/pkg/jwtgenerator"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Logger     *LoggerConfig
	Api        *api.Config
	GRPCServer *grpcserver.Config
	Auth       *jwtgenerator.Config
	DB         *database.Config
	User       *user.Config
}

type LoggerConfig struct {
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
}

func Parse() (*Config, error) {
	cfg := Config{}
	apiConfig := api.Config{}
	auth := jwtgenerator.Config{}
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

	err = envconfig.Process(project, &apiConfig)
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

	cfg.Api = &apiConfig
	cfg.Auth = &auth
	cfg.DB = &db
	cfg.Logger = &logger
	cfg.GRPCServer = &grpc
	cfg.User = &userConfig

	return &cfg, nil
}
