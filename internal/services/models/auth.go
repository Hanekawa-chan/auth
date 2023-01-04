package models

type GoogleAuthUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	VerifiedEmail bool   `json:"verified_email"`
}

type AuthRequest struct {
}

type SignupRequest struct {
	AuthHash string
	Username string
	Country  string
}

type Session struct {
	Token    string
	AuthHash string
}

type GoogleAuth struct {
	Code string
}

type PairAuth struct {
	Email    string
	Password string
}
