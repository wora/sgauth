package sgauth

import (
	"github.com/shinfan/sgauth/internal"
	"net/http"
	"github.com/shinfan/sgauth/jwt"
	"golang.org/x/net/context"
	"fmt"
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
func NewHTTPClient(ctx context.Context, credentials *Credentials) (*http.Client, error) {
	var ts internal.TokenSource
	var err error
	if (credentials != nil) {
		if (credentials.ServiceAccount != nil) {
			ts, err = serviceAccountTokenSource(ctx, credentials.ServiceAccount)
		}
	} else {
		ts, err = jwt.DefaultTokenSource(ctx, defaultScope)
		if err != nil {
			return nil, err
		}
	}
	return createClient(ts), nil
}

func serviceAccountTokenSource(ctx context.Context, account *ServiceAccount)(internal.TokenSource, error) {
	if (account.EnableOAuth) {
		return jwt.OAuthJSONTokenSource(ctx, account.JSONFile, account.Scopes)
	} else {
		aud := fmt.Sprintf("https://%s/%s", account.ServiceName, account.APIName)
		return jwt.JWTTokenSource(ctx, account.JSONFile, aud, account.Scopes)
	}
}
