package grpcserver

import (
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcValidator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/kanji-team/auth/internal/app"
	"github.com/kanji-team/auth/proto/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"
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

	server := new(config, nil)
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
func new(cfg *Config, middlewares ...grpc.UnaryServerInterceptor) *grpc.Server {
	interceptors := []grpc.UnaryServerInterceptor{
		grpcRecovery.UnaryServerInterceptor(grpcRecovery.WithRecoveryHandlerContext(recoveryHandler)),
		logIncomingRequestsMiddleware,
		grpcPrometheus.UnaryServerInterceptor,
		grpcValidator.UnaryServerInterceptor(),
	}

	interceptors = append(interceptors, middlewares...)
	return grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: cfg.MaxConnectionIdle * time.Minute,
		Timeout:           cfg.Timeout * time.Second,
		MaxConnectionAge:  cfg.MaxConnectionAge * time.Minute,
		Time:              cfg.Timeout * time.Minute,
	}), grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(interceptors...)))
}
