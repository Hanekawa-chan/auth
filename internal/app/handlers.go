package app

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/kanji-team/auth/proto/services"
	jwt "github.com/kanji-team/jwt"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io"
	"net/http"
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
		if err == ErrNotFound {
			uuID, err := uuid.NewUUID()
			if err != nil {
				return nil, err
			}
			authHash, err := generateAuthHash()
			if err != nil {
				return nil, err
			}

			user := &Credentials{
				Id:       uuID,
				Email:    authUser.Email,
				AuthHash: authHash,
			}

			if err = a.db.CreateUser(ctx, user); err != nil {
				return nil, err
			}

			err = a.Link(ctx, req)
			if err != nil {
				return nil, err
			}

			return &services.Session{SessionResponse: &services.Session_AuthHash{AuthHash: authHash}}, nil
		} else if err != nil {
			return nil, err
		}
		existUser, err = a.db.GetUserByID(ctx, user.Id)
		if err != nil {
			return nil, err
		}
	case *services.AuthRequest_Pair:
		existUser, err = a.db.GetUserByEmail(ctx, authUser.Email)
		if err == ErrNotFound {
			uuID, err := uuid.NewUUID()
			if err != nil {
				return nil, err
			}
			authHash, err := generateAuthHash()
			if err != nil {
				return nil, err
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(authUser.Password), 12)
			if err != nil {
				return nil, err
			}

			user := &Credentials{
				Id:       uuID,
				Email:    authUser.Email,
				Password: string(hash),
				AuthHash: authHash,
			}

			if err := a.db.CreateUser(ctx, user); err != nil {
				return nil, err
			}

			return &services.Session{SessionResponse: &services.Session_AuthHash{AuthHash: authHash}}, nil
		} else if err != nil {
			return nil, err
		}
	}

	if len(existUser.AuthHash) > 0 {
		return &services.Session{SessionResponse: &services.Session_AuthHash{AuthHash: existUser.AuthHash}}, nil
	}

	token, err := a.generateJWT(existUser.Id)
	if err != nil {
		return nil, err
	}

	return &services.Session{SessionResponse: &services.Session_Token{Token: token}}, nil
}

func (a *service) Signup(ctx context.Context, req *services.SignUpRequest) (*services.Session, error) {
	if req.AuthHash == "" {
		return nil, ErrEmptyRequired
	}

	_, err := a.db.GetUserByAuthHash(ctx, req.AuthHash)
	if err != nil {
		return nil, err
	}

	res, err := a.user.CreateUser(ctx, &services.CreateUserRequest{
		Email:   req.Email,
		Country: req.Country,
	})
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(res.UserId)
	if err != nil {
		return nil, err
	}

	err = a.db.UpdateId(ctx, id, req.AuthHash)
	if err != nil {
		return nil, err
	}

	token, err := a.generateJWT(id)
	if err != nil {
		return nil, err
	}

	err = a.db.RemoveAuthHash(ctx, id)
	if err != nil {
		return nil, err
	}

	return &services.Session{SessionResponse: &services.Session_Token{Token: token}}, nil
}

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
	googleUser, err := a.getUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &Credentials{Email: googleUser.Email}, nil
}

func (a *service) getUserByPair(ctx context.Context, login string, password string) (*Credentials, error) {
	err := a.validateEmail(login)
	if err != nil {
		return nil, ErrValidation
	}

	err = a.validatePassword(password)
	if err != nil {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
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

func (a *service) linkGoogle(ctx context.Context, req *services.GoogleAuth, id uuid.UUID) (*Google, error) {
	googleUser, err := a.getUserInfoFromGoogleAPI(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &Google{Id: id, Email: googleUser.Email, GoogleId: googleUser.ID}, nil
}

func (a *service) getUserInfoFromGoogleAPI(ctx context.Context, code string) (*GoogleAuthUser, error) {
	var userInfo GoogleAuthUser

	configGoogleAPI := &oauth2.Config{
		RedirectURL:  a.config.GoogleRedirectURL,
		ClientID:     a.config.GoogleClientID,
		ClientSecret: a.config.GoogleClientSecret,
		Scopes:       a.config.GoogleScopes,
		Endpoint:     google.Endpoint,
	}

	token, err := configGoogleAPI.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.config.GoogleOAuthURL+token.AccessToken, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(content, &userInfo)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil
}
