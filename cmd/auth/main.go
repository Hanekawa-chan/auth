package main

import (
	"github.com/kanji-team/auth/internal/app"
	"github.com/kanji-team/auth/internal/app/api"
	"github.com/kanji-team/auth/internal/app/config"
	"github.com/kanji-team/auth/internal/database"
	"github.com/kanji-team/auth/internal/grpcserver"
	"github.com/kanji-team/auth/internal/user"
	"github.com/kanji-team/auth/internal/version"
	jwt "github.com/kanji-team/jwt"
	"github.com/rs/zerolog"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	jwtGenerator, err := jwt.New(cfg.Auth.JWTSecretKey)
	if err != nil {
		zl.Fatal().Err(err).Msg("JWT init")
	}

	userClient := user.NewUserClient(zl, cfg.User)

	apiClient := api.NewAdapter(zl, cfg.Api)

	service := app.NewService(zl, cfg.Auth, userClient, apiClient, jwtGenerator, db)
	grpcServer := grpcserver.NewAdapter(zl, cfg.GRPCServer, service)

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

	time.Sleep(time.Second * 2)

	zl.Info().Msg("Shutdown - success")
}
