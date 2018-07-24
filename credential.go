package sgauth

type Credentials struct {
	ServiceAccount *ServiceAccount
	APIKey string
	AccessToken string
}

type ServiceAccount struct {
	EnableOAuth bool
	ServiceName string
	APIName string
	Scopes []string
	JSONFile string
}