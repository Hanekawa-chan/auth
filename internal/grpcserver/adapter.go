package grpcserver

import (
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/kanji-team/auth/internal/app"
	"github.com/kanji-team/auth/proto/services"
	"github.com/kanji-team/grpc-server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

type adapter struct {
	logger  *zerolog.Logger
	config  *Config
	server  *grpc.Server
	service app.Service
}

func NewAdapter(logger *zerolog.Logger, config *Config, service app.Service) app.GRPCServer {
	a := &adapter{
		logger:  logger,
		config:  config,
		service: service,
	}

	server := newServer(config, nil)
	a.server = server

	return a
}

func (a *adapter) ListenAndServe() error {
	services.RegisterAuthServiceServer(a.server, a)
	services.RegisterHealthServer(a.server, a)
	grpc_prometheus.Register(a.server)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		err := http.ListenAndServe(":6089", nil)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to serve metrics")
		}
	}()
	lis, err := net.Listen("tcp", a.config.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}

	log.Log().Msg("public server started on " + a.config.Address)
	if err := a.server.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("listen server")
	}
	return nil
}

func (a *adapter) Shutdown() {
	a.server.GracefulStop()
	a.logger.Info().Msg("Server Exited Properly")
}

func newServer(cfg *Config, middlewares ...grpc.UnaryServerInterceptor) *grpc.Server {
	config := &grpc_server.Config{
		MaxConnectionIdle: cfg.MaxConnectionIdle,
		Timeout:           cfg.Timeout,
		MaxConnectionAge:  cfg.MaxConnectionAge,
	}

	return grpc_server.NewGRPCServer(config, middlewares...)
}
