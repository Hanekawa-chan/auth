package httpserver

import (
	"context"
	"encoding/json"
	"github.com/Hanekawa-chan/kanji-auth/internal/services/models"
	"net/http"
)

func (a *adapter) googleAuth(w http.ResponseWriter, r *http.Request) error {
	req := models.GoogleAuth{}

	return a.auth(w, r, &req)
}

func (a *adapter) pairAuth(w http.ResponseWriter, r *http.Request) error {
	req := models.PairAuth{}

	return a.auth(w, r, &req)
}

func (a *adapter) auth(w http.ResponseWriter, r *http.Request, req models.AuthType) error {
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	resp, err := a.service.Auth(ctx, &models.AuthRequest{AuthType: req})
	if err != nil {
		return err
	}

	err = sendResponse(w, resp)
	if err != nil {
		return err
	}

	return err
}

func (a *adapter) signup(w http.ResponseWriter, r *http.Request) error {
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

func (a *adapter) linkGoogle(w http.ResponseWriter, r *http.Request) error {
	req := models.GoogleAuth{}

	return a.link(w, r, &req)
}

func (a *adapter) linkPair(w http.ResponseWriter, r *http.Request) error {
	req := models.PairAuth{}

	return a.link(w, r, &req)
}

func (a *adapter) link(w http.ResponseWriter, r *http.Request, req models.AuthType) error {
	ctx := context.Background()

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	err = a.service.Link(ctx, &models.AuthRequest{AuthType: req})
	if err != nil {
		return err
	}

	return err
}
