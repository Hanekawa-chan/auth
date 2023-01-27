package grpcserver

import (
	"github.com/kanji-team/auth/internal/app"
	"github.com/kanji-team/auth/proto/services"
	"github.com/kanji-team/grpc-server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net"
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
	lis, err := net.Listen("tcp", a.config.Address)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	services.RegisterHealthServer(a.server, a)
	log.Log().Msg("public server started")
	if err := a.server.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("listen server")
	}
	return nil
}

func (a *adapter) Shutdown() {
	a.server.GracefulStop()
	a.logger.Info().Msg("Server Exited Properly")
}

// New returns grpc server by config with middlewares
func newServer(cfg *Config, middlewares ...grpc.UnaryServerInterceptor) *grpc.Server {
	config := &grpc_server.Config{
		MaxConnectionIdle: cfg.MaxConnectionIdle,
		Timeout:           cfg.Timeout,
		MaxConnectionAge:  cfg.MaxConnectionAge,
	}

	return grpc_server.NewGRPCServer(config, middlewares...)
}
