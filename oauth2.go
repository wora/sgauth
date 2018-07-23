package sgauth

import (
	"net/http"
	"github.com/shinfan/sgauth/internal"
	"github.com/shinfan/sgauth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"crypto/x509"
	grpccredentials "google.golang.org/grpc/credentials"

	"fmt"
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
func DefaultGrpcConn(ctx context.Context, host string, port string, scope ...string) (*grpc.ClientConn, error) {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := grpccredentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := jwt.NewGrpcApplicationDefault(ctx, scope...)
	return grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
}

// NewClient creates an *http.Client from a TokenSource.
// The returned client is not valid beyond the lifetime of the context.
func JWTGrpcConn(ctx context.Context, host string, port string, aud string) (*grpc.ClientConn, error) {
	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := grpccredentials.NewClientTLSFromCert(pool, "")
	perRPC, _ := jwt.NewGrpcJWT(ctx, aud)
	return grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
}

// DefaultClient returns an HTTP Client that uses the
// DefaultTokenSource to obtain authentication credentials.
func DefaultClient(ctx context.Context, scope ...string) (*http.Client, error) {
	ts, err := jwt.DefaultTokenSource(ctx, scope...)
	if err != nil {
		return nil, err
	}
	return NewClient(ts), nil
}

// DefaultClient returns an HTTP Client that uses the
// DefaultTokenSource to obtain authentication credentials.
func JWTClient(ctx context.Context, aud string, scope ...string) (*http.Client, error) {
	creds, err := jwt.FindDefaultCredentials(ctx, scope)
	if creds != nil {
		ts, err := jwt.JWTAccessTokenSourceFromJSON(creds.JSON, aud)
		if (err != nil) {
			return nil, err
		}
		return NewClient(ts), nil
	}
	return nil, err
}

