package models

type AuthRequest struct {
	AuthType
}

type SignupRequest struct {
	AuthHash string `json:"auth_hash"`
	Email    string `json:"email"`
	Country  string `json:"country"`
}

type Session struct {
	Token    string `json:"token"`
	AuthHash string `json:"auth_hash"`
}

type GoogleAuth struct {
	Code string `json:"code"`
}

type PairAuth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthType interface {
	isAuthRequestAuthType()
}

func (*PairAuth) isAuthRequestAuthType()   {}
func (*GoogleAuth) isAuthRequestAuthType() {}
