package sgauth

import (
	"github.com/shinfan/sgauth/jwt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"crypto/x509"
	"google.golang.org/grpc/credentials"

	"fmt"
)

func NewGrpcConn(ctx context.Context, settings *Settings, host string, port string) (*grpc.ClientConn, error) {
	if settings == nil {
		settings = &Settings {
			Scope: defaultScope,
		}
	}

	pool, _ := x509.SystemCertPool()
	// error handling omitted
	creds := credentials.NewClientTLSFromCert(pool, "")
	var perRPC credentials.PerRPCCredentials
	if settings.Scope != "" {
		perRPC, _ = jwt.NewGrpcApplicationDefault(ctx, settings.Scope)
	} else {

		perRPC, _ = jwt.NewGrpcJWT(ctx, settings.Audience)
	}
	return grpc.Dial(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(creds),
		grpc.WithPerRPCCredentials(perRPC),
	)
}
