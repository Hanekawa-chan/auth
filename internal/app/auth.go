package app

import (
	"auth/proto/services"
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *service) Auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error) {
	authUser, err := a.getAuthUser(ctx, req)
	if err != nil {
		return nil, err
	}

	var existUser *Credentials
	switch req.AuthType.(type) {
	case *services.AuthRequest_Google:
		user, err := a.db.GetUserByGoogleEmail(ctx, authUser.Email)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		existUser, err = a.db.GetUserByID(ctx, user.Id)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
	case *services.AuthRequest_Pair:
		existUser, err = a.db.GetUserByEmail(ctx, authUser.Email)
		if err != nil {
			return nil, status.Error(codes.NotFound, err.Error())
		}
	}

	accessToken, err := a.generateAccessToken(existUser.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, err := a.generateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &services.Session{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *service) Refresh(ctx context.Context, req *services.RefreshRequest) (*services.Session, error) {
	id := ctx.Value("user_id").(uuid.UUID)

	err := a.parseRefreshToken(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := a.generateAccessToken(id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	refreshToken, err := a.generateRefreshToken()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &services.Session{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *service) getAuthUser(ctx context.Context, req *services.AuthRequest) (*Credentials, error) {
	switch req.AuthType.(type) {
	case *services.AuthRequest_Google:
		v := req.GetGoogle()
		return a.getUserByGoogle(ctx, v)
	case *services.AuthRequest_Pair:
		v := req.GetPair()
		return a.getUserByPair(ctx, v.Email, v.Password)
	}
	return nil, ErrType
}

func (a *service) getUserByGoogle(ctx context.Context, req *services.GoogleAuth) (*Credentials, error) {
	googleUser, err := a.api.GetUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &Credentials{Email: googleUser.Email}, nil
}

func (a *service) getUserByPair(ctx context.Context, login string, password string) (*Credentials, error) {
	err := a.validatePair(login, password)
	if err != nil {
		return nil, err
	}

	hash, err := a.hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := a.db.GetUserByEmail(ctx, login)
	if err == ErrNotFound {
		return &Credentials{
			Email:    login,
			Password: string(hash),
		}, nil
	}

	if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), hash); err != nil {
		return nil, ErrValidation
	}
	return user, nil
}
