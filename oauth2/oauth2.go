package oauth2

import (
	"net/http"
	"sgauth/oauth2/internal"
	"golang.org/x/net/context"
	"sgauth/oauth2/credentials"
)

// NewClient creates an *http.Client from a TokenSource.
// The returned client is not valid beyond the lifetime of the context.
func NewClient(src internal.TokenSource) *http.Client {
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

// DefaultClient returns an HTTP Client that uses the
// DefaultTokenSource to obtain authentication credentials.
func DefaultClient(ctx context.Context, scope ...string) (*http.Client, error) {
	ts, err := credentials.DefaultTokenSource(ctx, scope...)
	if err != nil {
		return nil, err
	}
	return NewClient(ts), nil
}

