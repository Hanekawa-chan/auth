package app

import (
	"context"
	"github.com/google/uuid"
	"github.com/kanji-team/auth/proto/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *service) SignUp(ctx context.Context, req *services.SignUpRequest) (*services.Session, error) {
	var err error
	id := uuid.UUID{}

	switch req.AuthType.(type) {
	case *services.SignUpRequest_Google:
		id, err = a.signUpGoogle(ctx, req.GetGoogle())
		if err != nil {
			return nil, err
		}

	case *services.SignUpRequest_Pair:
		id, err = a.signUpPair(ctx, req.GetPair())
		if err != nil {
			return nil, err
		}
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

func (a *service) signUpGoogle(ctx context.Context, req *services.GoogleAuth) (uuid.UUID, error) {
	googleUser, err := a.api.GetUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return uuid.UUID{}, err
	}

	res, err := a.user.CreateUser(ctx, &services.CreateUserRequest{
		Name:  googleUser.Name,
		Email: googleUser.Email,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(res.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = a.db.CreateUser(ctx, &Credentials{
		Id:            id,
		Email:         googleUser.Email,
		VerifiedEmail: true,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	err = a.db.CreateGoogle(ctx, &Google{
		Id:       id,
		Email:    googleUser.Email,
		GoogleId: googleUser.ID,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, err
}

func (a *service) signUpPair(ctx context.Context, req *services.PairSignUp) (uuid.UUID, error) {
	err := a.validatePair(req.Pair.Email, req.Pair.Password)
	if err != nil {
		return uuid.UUID{}, err
	}

	hash, err := a.hashPassword(req.Pair.Password)
	if err != nil {
		return uuid.UUID{}, err
	}

	res, err := a.user.CreateUser(ctx, &services.CreateUserRequest{
		Name:  req.Name,
		Email: req.Pair.Email,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	id, err := uuid.Parse(res.UserId)
	if err != nil {
		return uuid.UUID{}, err
	}

	err = a.db.CreateUser(ctx, &Credentials{
		Id:            id,
		Email:         req.Pair.Email,
		Password:      string(hash),
		VerifiedEmail: false,
	})
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, err
}
