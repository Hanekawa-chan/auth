package app

import (
	"auth/proto/services"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type Service interface {
	Auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error)
	SignUp(ctx context.Context, req *services.SignUpRequest) (*services.Session, error)
	Link(ctx context.Context, req *services.AuthRequest) error
	ValidateSession(ctx context.Context, session *services.Session) (*services.ValidSession, error)
}

type GRPCServer interface {
	ListenAndServe() error
	Shutdown()
}

type Database interface {
	CreateUser(ctx context.Context, user *Credentials) error
	GetUserByEmail(ctx context.Context, login string) (*Credentials, error)
	GetUserByGoogleEmail(ctx context.Context, email string) (*Google, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*Credentials, error)
	CreateGoogle(ctx context.Context, creds *Google) error
	UpdateIssuedAt(ctx context.Context, id uuid.UUID, issuedAt int64) error
}

type User interface {
	CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.CreateUserResponse, error)
}

type Api interface {
	GetUserInfoFromGoogleAPI(ctx context.Context, code string) (*GoogleAuthUser, error)
}

type JWT interface {
	Generate(claims map[string]interface{}) (string, error)
	ParseToken(token string) (jwt.MapClaims, error)
}
