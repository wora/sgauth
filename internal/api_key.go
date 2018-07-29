package internal

import "golang.org/x/net/context"

// TokenSource supplies PerRPCCredentials from an oauth2.TokenSource.
type GrpcApiKey struct {
	Value string
}

// GetRequestMetadata gets the request metadata as a map from a TokenSource.
func (key GrpcApiKey) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"x-goog-api-key": key.Value,
	}, nil
}

// RequireTransportSecurity indicates whether the credentials requires transport security.
func (key GrpcApiKey) RequireTransportSecurity() bool {
	return true
}
