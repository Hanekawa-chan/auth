package grpcserver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hanekawa-chan/kanji-auth/internal/app"
	"github.com/Hanekawa-chan/kanji-auth/internal/app/config"
	"github.com/go-chi/chi"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/http"
	"time"
)

type adapter struct {
	logger  *zerolog.Logger
	config  *config.Config
	server  *grpc.Server
	service app.Service
}

func NewAdapter(logger *zerolog.Logger, config *config.Config, service app.Service) app.HTTPServer {
	a := &adapter{
		logger:  logger,
		config:  config,
		service: service,
	}

	r := chi.NewRouter()

	a.initMiddlewares(r)

	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health-check", wrap(a.HealthCheck))

	a.server = &http.Server{
		Addr:    config.HTTPServer.Address,
		Handler: r,
	}

	return a
}

func recoveryHandler(ctx context.Context, p interface{}) (err error) {
	log.Error().
		Str("panic", fmt.Sprintf("%+v", p)).
		Str("ctx", fmt.Sprintf("%+v", ctx)).
		Msg("PANIC")
	if err, ok := p.(error); ok {
		return fmt.Errorf("panic: %w", err)
	}
	return fmt.Errorf("panic: %+v", p)
}

// New returns grpc server by config with middlewares
func New(conf *config.Config, middlewares ...grpc.UnaryServerInterceptor) *grpc.Server {
	interceptors := []grpc.UnaryServerInterceptor{
		grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandlerContext(recoveryHandler)),
		logIncomingRequestsMiddleware,
		grpcPrometheus.UnaryServerInterceptor,
		grpcValidator.UnaryServerInterceptor(),
	}

	interceptors = append(interceptors, middlewares...)
	return grpc.NewServer(grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(interceptors...)))
}

func logIncomingRequestsMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	requestJSON, _ := json.Marshal(req)
	result, err := handler(ctx, req)
	responseJSON, _ := json.Marshal(result)
	var logEvent *zerolog.Event
	if err != nil {
		logEvent = log.Error().Str("error", fmt.Sprintf("%+v", err))
	} else {
		logEvent = log.Info()
	}
	logEvent.
		Dur("duration", time.Since(start)).
		RawJSON("json_response", responseJSON).
		RawJSON("json_request", requestJSON).
		Str("url", info.FullMethod).
		Str("ctx", fmt.Sprintf("%+v", ctx)).
		Msg("complete")

	return result, err
}
