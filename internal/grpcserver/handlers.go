package grpcserver

import (
	"context"
	"github.com/Hanekawa-chan/kanji-auth/proto/services"
)

func (a *adapter) auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error) {
	session, err := a.service.Auth(ctx, req)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (a *adapter) signup(ctx context.Context, req *services.SignUpRequest) (*services.Session, error) {
	session, err := a.service.Signup(ctx, &req)
	if err != nil {
		return nil, err
	}

	return session, err
}

func (a *adapter) link(ctx context.Context, req *services.AuthRequest) error {
	err := a.service.Link(ctx, req)
	if err != nil {
		return err
	}

	return err
}
