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
func createAuthTokenClient(src internal.TokenSource) *http.Client {
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

// createClient creates an *http.Client from a TokenSource.
// The returned client is not valid beyond the lifetime of the context.
func createAPIKeyClient(key string) *http.Client {
	if key == "" {
		return http.DefaultClient
	}
	return &http.Client{
		Transport: &Transport{
			Base:   http.DefaultClient.Transport,
			APIKey: key,
		},
	}
}

var defaultScope = "https://www.googleapis.com/auth/cloud-platform"

// DefaultClient returns an HTTP Client that uses the
// DefaultTokenSource to obtain authentication credentials.
func NewHTTPClient(ctx context.Context, settings *Settings) (*http.Client, error) {
	if settings.APIKey != "" {
		return createAPIKeyClient(settings.APIKey), nil
	} else {
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
		return createAuthTokenClient(ts), nil
	}
}
