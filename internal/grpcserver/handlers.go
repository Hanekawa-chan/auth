package grpcserver

import (
	"context"
	"github.com/kanji-team/auth/proto/services"
)

func (a *adapter) Auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error) {
	session, err := a.service.Auth(ctx, req)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (a *adapter) SignUp(ctx context.Context, req *services.SignUpRequest) (*services.Session, error) {
	session, err := a.service.Signup(ctx, req)
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
