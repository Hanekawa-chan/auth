package main

import (
	"auth/config"
	"auth/internal/app"
	"auth/internal/database"
	"auth/internal/grpcserver"
	"auth/internal/user"
	"auth/internal/version"
	"auth/pkg/api"
	"auth/pkg/jwtgenerator"
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Print version and commit sha
	log.Println("Loading Mailing - v", version.Version, "| Commit:", version.Commit)

	// Parse all configs form env
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Parse log level
	level, err := zerolog.ParseLevel(cfg.Logger.LogLevel)
	_ = err
	if err != nil {
		log.Fatal(err)
	}

	// Initializations
	logger := zerolog.New(os.Stdout).Level(level)

	zl := &logger

	db, err := database.NewAdapter(zl, cfg.DB)
	if err != nil {
		zl.Fatal().Err(err).Msg("Database init")
	}

	jwt, err := jwtgenerator.NewAdapter(cfg.Auth)
	if err != nil {
		zl.Fatal().Err(err).Msg("JWT init")
	}

	userClient := user.NewUserClient(zl, cfg.User)

	apiClient := api.NewAdapter(zl, cfg.Api)

	service := app.NewService(zl, userClient, apiClient, jwt, db)
	grpcServer := grpcserver.NewAdapter(zl, cfg.GRPCServer, service)
	zl.Info().Msg("initialized everything")

	// Channels for errors and os signals
	stop := make(chan error, 1)
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)

	// Receive errors form start bot func into error channel
	go func(stop chan<- error) {
		stop <- grpcServer.ListenAndServe()
	}(stop)

	// Blocking select
	select {
	case sig := <-osSig:
		zl.Info().Msgf("Received os syscall signal %v", sig)
	case err := <-stop:
		zl.Error().Err(err).Msg("Received Error signal")
	}

	// Shutdown code
	zl.Info().Msg("Shutting down...")

	grpcServer.Shutdown()

	zl.Info().Msg("Shutdown - success")
}
