package sgauth

// An extensible structure that holds the credentials for
// Google API authentication.
type Credentials struct {
	// Service account authentication method
	ServiceAccount *ServiceAccount
	// The Google API key.
	// https://cloud.google.com/docs/authentication/api-keys
	APIKey string
	// OAuth2.0 access token.
	AccessToken string
}

// A structure that holds information required for Google
// API service account authentication flow.
// For more information please visit:
// https://developers.google.com/identity/protocols/OAuth2ServiceAccount
type ServiceAccount struct {
	// Returns true if OAuth2.0 flow is enabled. Otherwise
	// client-signed JWT token flow will be used.
	// Default to false.
	EnableOAuth bool
	// The full host name of the API. This field is used to construct the
	// audience field for the client-signed JWT token.
	ServiceName string
	// The full API name. This field is used to construct the audience field
	// for the client-signed JWT token.
	APIName string
	// The required scopes of the API requests. If not set, the default
	// scope will be used.
	Scopes []string
	// The full path of the downloaded service account JSON file.
	// If not set, the application default file path will be used.
	JSONFile string
}
