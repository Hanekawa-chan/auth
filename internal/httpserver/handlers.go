package httpserver

import (
	"context"
	"encoding/json"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	"net/http"
)

func (a *adapter) Auth(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	req := models.AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	resp, err := a.service.Auth(ctx, &req)
	if err != nil {
		return err
	}

	err = sendResponse(w, resp)
	if err != nil {
		return err
	}

	return err
}

func (a *adapter) Signup(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	req := models.SignupRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	resp, err := a.service.Signup(ctx, &req)
	if err != nil {
		return err
	}

	err = sendResponse(w, resp)
	if err != nil {
		return err
	}

	return err
}

func (a *adapter) Link(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()
	req := models.AuthRequest{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = a.service.Link(ctx, &req)
	if err != nil {
		return err
	}
	return err
}
