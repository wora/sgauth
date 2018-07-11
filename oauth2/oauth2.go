package oauth2

import (
	"net/http"
	"sgauth/oauth2/internal"
	"golang.org/x/net/context"
	"sgauth/oauth2/credentials"
	"google.golang.org/grpc"
	"crypto/x509"
	grpccredentials "google.golang.org/grpc/credentials"

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

// NewClient creates an *http.Client from a TokenSource.
// The returned client is not valid beyond the lifetime of the context.
func DefaultGrpcConn(ctx context.Context, scope ...string) (*grpc.ClientConn, error) {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := grpccredentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := credentials.NewApplicationDefault(ctx, scope...)
	return grpc.Dial(
		"servicemanagement.googleapis.com:443",
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
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

