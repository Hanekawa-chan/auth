package httpserver

import (
	"context"
	"errors"
	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
	"kanji-auth/internal/app"
	"net/http"
)

type adapter struct {
	logger  *zerolog.Logger
	config  *app.Config
	server  *http.Server
	service app.Service
	limiter ratelimit.Limiter
}

func NewAdapter(logger *zerolog.Logger, config *app.Config, service app.Service) app.HTTPServer {
	a := &adapter{
		logger:  logger,
		config:  config,
		service: service,
		limiter: ratelimit.New(config.HTTPServer.RateLimit),
	}

	r := chi.NewRouter()

	a.InitMiddlewares(r)

	r.Group(func(r chi.Router) {
		r.Use(a.authMiddleware)
	})

	r.MethodFunc(http.MethodPost, "/api/v1/auth", a.Auth)
	r.MethodFunc(http.MethodPost, "/api/v1/signup", a.Signup)

	r.Handle("/metrics", promhttp.Handler())
	r.MethodFunc(http.MethodGet, "/health-check", a.HealthCheck)

	a.server = &http.Server{
		Addr:    config.HTTPServer.Address,
		Handler: r,
	}

	return a
}

func (a adapter) ListenAndServe() error {
	a.logger.Info().Msgf("Listening and serving HTTP requests on: %v", a.config.HTTPServer.Address)

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error().Err(err).Msg("Error listening and serving HTTP requests.")
		return err
	}

	return nil
}

func (a adapter) Shutdown(ctx context.Context) error {
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Error().Err(err).Msg("Error shutting down HTTP adapter!")
		return err
	}

	return nil
}

func (a adapter) HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		return
	}
}
