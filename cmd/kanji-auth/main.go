package main

import (
	"context"
	"github.com/Hanekawa-chan/kanji-auth/internal/app"
	"github.com/Hanekawa-chan/kanji-auth/internal/app/config"
	"github.com/Hanekawa-chan/kanji-auth/internal/database"
	"github.com/Hanekawa-chan/kanji-auth/internal/grpcserver"
	"github.com/Hanekawa-chan/kanji-auth/internal/user"
	"github.com/Hanekawa-chan/kanji-auth/internal/version"
	kanjiJwt "github.com/Hanekawa-chan/kanji-jwt"
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

	db, err := database.NewAdapter(zl, cfg)
	if err != nil {
		zl.Fatal().Err(err).Msg("Database init")
	}

	jwtGenerator, err := kanjiJwt.New(cfg.Auth.JWTSecretKey)
	if err != nil {
		zl.Fatal().Err(err).Msg("JWT init")
	}

	userClient := user.NewUserClient(zl, cfg.User)

	service := app.NewService(zl, cfg, userClient, jwtGenerator, db)
	httpServerAdapter := grpcserver.NewAdapter(zl, cfg, service)

	// Channels for errors and os signals
	stop := make(chan error, 1)
	osSig := make(chan os.Signal, 1)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)

	// Receive errors form start bot func into error channel
	go func(stop chan<- error) {
		stop <- httpServerAdapter.ListenAndServe()
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := httpServerAdapter.Shutdown(ctx); err != nil {
		zl.Error().Err(err).Msg("Error shutting down the HTTP server!")
	}

	time.Sleep(time.Second * 2)

	zl.Info().Msg("Shutdown - success")
}
