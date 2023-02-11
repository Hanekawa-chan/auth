package grpcserver

import (
	"auth/proto/services"
	"context"
	"time"
)

func (a *adapter) Auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error) {
	session, err := a.service.Auth(ctx, req)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (a *adapter) SignUp(ctx context.Context, req *services.SignUpRequest) (*services.Session, error) {
	a.logger.Debug().Msg("got sign up grpc")
	session, err := a.service.SignUp(ctx, req)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (a *adapter) Link(ctx context.Context, req *services.AuthRequest) (*services.Empty, error) {
	err := a.service.Link(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, err
}

func (a *adapter) Refresh(ctx context.Context, req *services.RefreshRequest) (*services.Session, error) {
	session, err := a.service.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (a *adapter) Check(context.Context, *services.HealthCheckRequest) (*services.HealthCheckResponse, error) {
	return &services.HealthCheckResponse{Status: services.HealthCheckResponse_SERVING}, nil
}

func (a *adapter) Watch(_ *services.HealthCheckRequest, server services.Health_WatchServer) error {
	var err error
	for {
		time.Sleep(a.config.HealthCheckRate)
		err = server.Send(&services.HealthCheckResponse{Status: services.HealthCheckResponse_SERVING})
		if err != nil {
			break
		}
	}
	return err
}
