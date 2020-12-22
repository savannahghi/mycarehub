package profile

import "gitlab.slade360emr.com/go/base"

// Default login client settings (env var names)
const (
	CoreEDIClientIDEnvVarName     = "CORE_CLIENT_ID"
	CoreEDIClientSecretEnvVarName = "CORE_CLIENT_SECRET"
	CoreEDIUsernameEnvVarName     = "CORE_USERNAME"
	CoreEDIPasswordEnvVarName     = "CORE_PASSWORD"
	CoreEDIGrantTypeEnvVarName    = "CORE_GRANT_TYPE"
	CoreEDIAPISchemeEnvVarName    = "CORE_API_SCHEME"
	CoreEDITokenURLEnvVarName     = "CORE_TOKEN_URL"
	CoreEDIAPIHostEnvVarName      = "CORE_HOST"
)

// LoginClient returns an API client that is logged in with the supplied username and password
// to EDI Core authserver
func LoginClient(username string, password string) (base.Client, error) {
	clientID, clientIDErr := base.GetEnvVar(CoreEDIClientIDEnvVarName)
	if clientIDErr != nil {
		return nil, clientIDErr
	}

	clientSecret, clientSecretErr := base.GetEnvVar(CoreEDIClientSecretEnvVarName)
	if clientSecretErr != nil {
		return nil, clientSecretErr
	}

	apiTokenURL, apiTokenURLErr := base.GetEnvVar(CoreEDITokenURLEnvVarName)
	if apiTokenURLErr != nil {
		return nil, apiTokenURLErr
	}

	apiHost, apiHostErr := base.GetEnvVar(CoreEDIAPIHostEnvVarName)
	if apiHostErr != nil {
		return nil, apiHostErr
	}

	apiScheme, apiSchemeErr := base.GetEnvVar(CoreEDIAPISchemeEnvVarName)
	if apiSchemeErr != nil {
		return nil, apiSchemeErr
	}

	grantType, grantTypeErr := base.GetEnvVar(CoreEDIGrantTypeEnvVarName)
	if grantTypeErr != nil {
		return nil, grantTypeErr
	}
	extraHeaders := make(map[string]string)
	return base.NewServerClient(
		clientID, clientSecret, apiTokenURL, apiHost, apiScheme, grantType, username, password, extraHeaders)
}
