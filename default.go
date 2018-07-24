package sgauth

import (
	"github.com/shinfan/sgauth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"crypto/x509"
	grpccredentials "google.golang.org/grpc/credentials"

	"fmt"
)

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
