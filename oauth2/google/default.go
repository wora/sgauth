// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"golang.org/x/net/context"
	"sgauth/oauth2"
)

// DefaultClient returns an HTTP Client that uses the
// DefaultTokenSource to obtain authentication credentials.
func DefaultClient(ctx context.Context, scope ...string) (*http.Client, error) {
	ts, err := DefaultTokenSource(ctx, scope...)
	if err != nil {
		return nil, err
	}
	return oauth2.NewClient(ctx, ts), nil
}

// Credentials holds Google credentials, including "Application Default Credentials".
// For more details, see:
// https://developers.google.com/accounts/docs/application-default-credentials
type Credentials struct {
	ProjectID   string // may be empty
	TokenSource oauth2.TokenSource

	// JSON contains the raw bytes from a JSON credentials file.
	// This field may be nil if authentication is provided by the
	// environment and not with a credentials file, e.g. when code is
	// running on Google Cloud Platform.
	JSON []byte
}

// DefaultTokenSource returns the token source for
// "Application Default Credentials".
// It is a shortcut for FindDefaultCredentials(ctx, scope).TokenSource.
func DefaultTokenSource(ctx context.Context, scope ...string) (oauth2.TokenSource, error) {
	creds, err := FindDefaultCredentials(ctx, scope)
	if err != nil {
		return nil, err
	}
	return creds.TokenSource, nil
}

func FindDefaultCredentials(ctx context.Context, scopes []string) (*Credentials, error) {
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

func CredentialsFromJSON(ctx context.Context, jsonData []byte, scopes []string) (*Credentials, error) {
	var f credentialsFile
	if err := json.Unmarshal(jsonData, &f); err != nil {
		return nil, err
	}
	ts, err := f.tokenSource(ctx, append([]string(nil), scopes...))
	if err != nil {
		return nil, err
	}
	return &Credentials{
		ProjectID:   f.ProjectID,
		TokenSource: ts,
		JSON:        jsonData,
	}, nil
}

func readCredentialsFile(ctx context.Context, filename string, scopes []string) (*Credentials, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return CredentialsFromJSON(ctx, b, scopes)
}