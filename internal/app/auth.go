package app

import (
	"auth/proto/services"
	"context"
	"github.com/golang-jwt/jwt/v4"
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

	accessToken, refreshToken, err := a.generateTokens(ctx, existUser.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &services.Session{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}

func (a *service) ValidateSession(ctx context.Context, session *services.Session) (*services.ValidSession, error) {
	id, iat, err := a.parseAccessToken(session.AccessToken)
	if err == jwt.ErrTokenExpired {
		user, err := a.db.GetUserByID(ctx, id)
		if err != nil {
			return nil, err
		}

		if user.IssuedAt > iat {
			return nil, ErrInvalidated
		}

		session, err = a.refresh(context.WithValue(ctx, "user_id", id), session.RefreshToken)
		if err != nil {
			return nil, err
		}

		id, iat, err = a.parseAccessToken(session.AccessToken)
		if err != nil {
			return nil, err
		}

		return &services.ValidSession{
			UserId:       id.String(),
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
		}, nil
	} else if err != nil {
		a.logger.Err(err).Msg("token parse")
		return nil, err
	}

	user, err := a.db.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.IssuedAt > iat {
		return nil, ErrInvalidated
	}

	return &services.ValidSession{
		UserId:       id.String(),
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
	}, nil
}

func (a *service) refresh(ctx context.Context, refreshToken string) (*services.Session, error) {
	id := ctx.Value("user_id").(uuid.UUID)

	err := a.parseRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := a.generateTokens(ctx, id)
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
