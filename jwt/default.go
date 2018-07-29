package jwt

import (
	"google.golang.org/grpc/credentials"
	"context"
	"github.com/shinfan/sgauth/internal"
	"os"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

// DefaultTokenSource returns the token source for
// "Application Default Credentials".
// It is a shortcut for FindDefaultCredentials(ctx, scope).TokenSource.
func DefaultTokenSource(ctx context.Context, scope ...string) (internal.TokenSource, error) {
	creds, err := applicationDefaultCredentials(ctx, scope)
	if err != nil {
		return nil, err
	}
	return creds.TokenSource, nil
}

func OAuthJSONTokenSource(ctx context.Context, json string, scopes []string) (internal.TokenSource, error) {
	creds, err := findJSONCredentials(ctx, json, scopes)
	if err != nil {
		return nil, err
	}
	return creds.TokenSource, nil

}

func JWTTokenSource(ctx context.Context, json string, aud string, scopes []string) (internal.TokenSource, error) {
	creds, err := findJSONCredentials(ctx, json, scopes)
	if err != nil {
		return nil, err
	}
	ts, err := JWTAccessTokenSourceFromJSON(creds.JSON, aud)
	return ts, err
}

// NewApplicationDefault returns "Application Default Credentials". For more
// detail, see https://developers.google.com/accounts/docs/application-default-credentials.
func NewGrpcApplicationDefault(ctx context.Context, scope ...string) (credentials.PerRPCCredentials, error) {
	t, err := DefaultTokenSource(ctx, scope...)
	if err != nil {
		return nil, err
	}
	return internal.GrpcTokenSource{t}, nil
}

// NewApplicationDefault returns "Application Default Credentials". For more
// detail, see https://developers.google.com/accounts/docs/application-default-credentials.
func NewGrpcJWT(ctx context.Context, aud string) (credentials.PerRPCCredentials, error) {
	creds, err := applicationDefaultCredentials(ctx, []string{})
	if creds != nil {
		ts, err := JWTAccessTokenSourceFromJSON(creds.JSON, aud)
		if (err != nil) {
			return nil, err
		}
		return internal.GrpcTokenSource{ts}, nil
	}
	return nil, err
}

func findJSONCredentials(ctx context.Context, json string, scopes[]string) (*internal.Credentials, error) {
	if json != "" {
		return credentialsFromJSON(ctx, []byte(json), scopes)

	} else {
		return applicationDefaultCredentials(ctx, scopes)

	}
}

func applicationDefaultCredentials(ctx context.Context, scopes []string) (*internal.Credentials, error) {
	const envVar = "GOOGLE_APPLICATION_CREDENTIALS"
	if filename := os.Getenv(envVar); filename != "" {
		creds, err := readCredentialsFile(ctx, filename, scopes)
		if err != nil {
			return nil, fmt.Errorf("google: error getting credentials using %v environment variable: %v", envVar, err)
		}
		return creds, nil
	}
	// None are found; return helpful error.
	const url = "https://developers.google.com/accounts/docs/application-default-credentials"
	return nil, fmt.Errorf("google: could not find default credentials. See %v for more information.", url)
}

func readCredentialsFile(ctx context.Context, filename string, scopes []string) (*internal.Credentials, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return credentialsFromJSON(ctx, b, scopes)
}

func credentialsFromJSON(ctx context.Context, jsonData []byte, scopes []string) (*internal.Credentials, error) {
	var f CredentialsFile
	if err := json.Unmarshal(jsonData, &f); err != nil {
		return nil, err
	}
	ts, err := f.tokenSource(ctx, append([]string(nil), scopes...))
	if err != nil {
		return nil, err
	}
	return &internal.Credentials{
		ProjectID:   f.ProjectID,
		TokenSource: ts,
		JSON:        jsonData,
	}, nil
}
