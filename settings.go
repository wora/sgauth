package sgauth

var MethodOAuth = "oauth"
var MethodJWT = "jwt"
var MethodAPIKey = "apikey"

// An extensible structure that holds the credentials for
// Google API authentication.
type Settings struct {
	// The JSON credentials content downloaded from Google Cloud Console.
	CredentialsJSON string
	// If specified, use OAuth. Otherwise, JWT.
	Scope string
	// The audience field for JWT auth
	Audience string
	// The Google API key
	APIKey string
	// This is only used for domain-wide delegation.
	User string
	// This name is confusing now. Since we have quotaUser and userProject.
	// We should have named them as quotaUser and quotaProject.
	QuotaUser string
	QuotaProject string
}

func (s Settings)AuthMethod() string {
	if s.APIKey != "" {
		return MethodAPIKey
	} else if s.Scope != "" {
		return MethodOAuth
	}
	return MethodJWT
}

