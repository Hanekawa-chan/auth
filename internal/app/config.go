package app

type Config struct {
	GoogleRedirectURL  string   `envconfig:"GOOGLE_REDIRECT_URL"`
	GoogleClientID     string   `envconfig:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string   `envconfig:"GOOGLE_CLIENT_SECRET"`
	GoogleOAuthURL     string   `envconfig:"GOOGLE_OAUTH_URL"`
	GoogleScopes       []string `envconfig:"GOOGLE_SCOPES"`
	JWTSecretKey       string   `envconfig:"JWT_SECRET_KEY"`
}
