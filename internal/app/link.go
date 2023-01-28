package app

import (
	"context"
	"github.com/google/uuid"
	"github.com/kanji-team/auth/proto/services"
	jwt "github.com/kanji-team/jwt"
)

func (a *service) Link(ctx context.Context, req *services.AuthRequest) error {
	id, err := jwt.GetUserId(ctx, a.jwtGenerator)
	if err != nil {
		return err
	}
	switch req.AuthType.(type) {
	case *services.AuthRequest_Google:
		v := req.GetGoogle()
		cred, err := a.linkGoogle(ctx, v, id)
		if err != nil {
			return err
		}
		err = a.db.CreateGoogle(ctx, cred)
		if err != nil {
			return err
		}
		return nil
	}
	return ErrType
}

func (a *service) linkGoogle(ctx context.Context, req *services.GoogleAuth, id uuid.UUID) (*Google, error) {
	googleUser, err := a.api.GetUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &Google{Id: id, Email: googleUser.Email, GoogleId: googleUser.ID}, nil
}
