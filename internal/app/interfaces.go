package app

import (
	"context"
	"github.com/google/uuid"
	"github.com/kanji-team/auth/proto/services"
)

type Service interface {
	Auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error)
	Signup(ctx context.Context, req *services.SignUpRequest) (*services.Session, error)
	Link(ctx context.Context, req *services.AuthRequest) error
}

type GRPCServer interface {
	ListenAndServe() error
	Shutdown()
}

type Database interface {
	GetUserByAuthHash(ctx context.Context, hash string) (*Credentials, error)
	UpdateId(ctx context.Context, id uuid.UUID, hash string) error
	RemoveAuthHash(ctx context.Context, id uuid.UUID) error
	CreateUser(ctx context.Context, user *Credentials) error
	GetUserByEmail(ctx context.Context, login string) (*Credentials, error)
	GetUserByGoogleEmail(ctx context.Context, email string) (*Credentials, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*Credentials, error)
	CreateGoogle(ctx context.Context, creds *Google) error
}

type JWTGenerator interface {
	Generate(claims map[string]interface{}) (string, error)
}

type User interface {
	CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.CreateUserResponse, error)
}
