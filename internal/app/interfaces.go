package app

import (
	"context"
	"github.com/google/uuid"
	"github.com/kanji-team/auth/proto/services"
)

type Service interface {
	Auth(ctx context.Context, req *services.AuthRequest) (*services.Session, error)
	SignUp(ctx context.Context, req *services.SignUpRequest) (*services.Session, error)
	Link(ctx context.Context, req *services.AuthRequest) error
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
}

type User interface {
	CreateUser(ctx context.Context, req *services.CreateUserRequest) (*services.CreateUserResponse, error)
}

type Api interface {
	GetUserInfoFromGoogleAPI(ctx context.Context, code string) (*GoogleAuthUser, error)
}
