package httpserver

import (
	"context"
	"kanji-auth/internal/services/models"
	"net/http"
)

func (a *adapter) Auth(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.AuthRequest{}
	_, err := a.service.Auth(ctx, &req)
	if err != nil {
		return
	}
}

func (a *adapter) Signup(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	req := models.SignupRequest{}
	_, err := a.service.Signup(ctx, &req)
	if err != nil {
		return
	}
}
