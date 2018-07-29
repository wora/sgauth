package sgauth

import (
	"github.com/shinfan/sgauth/internal"
	"net/http"
	"github.com/shinfan/sgauth/jwt"
	"golang.org/x/net/context"
	"strings"
)

// createClient creates an *http.Client from a TokenSource.
// The returned client is not valid beyond the lifetime of the context.
func createClient(src internal.TokenSource) *http.Client {
	if src == nil {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: &Transport{
			Base:   http.DefaultClient.Transport,
			Source: internal.ReuseTokenSource(nil, src),
		},
	}
}

var defaultScope = "https://www.googleapis.com/auth/cloud-platform"

// DefaultClient returns an HTTP Client that uses the
// DefaultTokenSource to obtain authentication credentials.
func NewHTTPClient(ctx context.Context, settings *Settings) (*http.Client, error) {
	var ts internal.TokenSource
	var err error
	if settings != nil {
		if settings.Scope != "" {
			ts, err = jwt.OAuthJSONTokenSource(ctx, settings.CredentialsJSON, strings.Split(settings.Scope, " "))
		} else {
			ts, err = jwt.JWTTokenSource(
				ctx, settings.CredentialsJSON, settings.Audience, []string{})
		}
	} else {
		ts, err = jwt.DefaultTokenSource(ctx, defaultScope)
		if err != nil {
			return nil, err
		}
	}
	return createClient(ts), nil
}
