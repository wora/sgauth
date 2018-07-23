package jwt

import (
	"errors"
	"fmt"
	"github.com/shinfan/sgauth/internal"
	"golang.org/x/net/context"
)

// DefaultTokenURL is Google's OAuth 2.0 token URL to use with the JWT flow.
const DefaultTokenURL = "https://accounts.google.com/o/oauth2/token"

// JSON key file types.
const (
	ServiceAccountKey  = "service_account"
	UserCredentialsKey = "authorized_user"
)

// credentialsFile is the unmarshalled representation of a credentials file.
type CredentialsFile struct {
	Type string `json:"type"` // serviceAccountKey or userCredentialsKey

	// Service Account fields
	ClientEmail  string `json:"client_email"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	TokenURL     string `json:"token_uri"`
	ProjectID    string `json:"project_id"`

	// User Credential fields
	// (These typically come from gcloud auth.)
	ClientSecret string `json:"client_secret"`
	ClientID     string `json:"client_id"`
	RefreshToken string `json:"refresh_token"`
}

func (f *CredentialsFile) jwtConfig(scopes []string) *JWTConfig {
	cfg := &JWTConfig{
		Email:        f.ClientEmail,
		PrivateKey:   []byte(f.PrivateKey),
		PrivateKeyID: f.PrivateKeyID,
		Scopes:       scopes,
		TokenURL:     f.TokenURL,
	}
	if cfg.TokenURL == "" {
		cfg.TokenURL = DefaultTokenURL
	}
	return cfg
}

func (f *CredentialsFile) tokenSource(ctx context.Context, scopes []string) (internal.TokenSource, error) {
	switch f.Type {
	case ServiceAccountKey:
		cfg := f.jwtConfig(scopes)
		return cfg.TokenSource(ctx), nil
	case "":
		return nil, errors.New("missing 'type' field in credentials")
	default:
		return nil, fmt.Errorf("unknown credential type: %q", f.Type)
	}
}
