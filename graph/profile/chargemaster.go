package profile

import "gitlab.slade360emr.com/go/base"

const (
	// ChargeMasterHostEnvVarName is the name of an environment variable that
	//points at the API root e.g "https://base.chargemaster.slade360emr.com/v1"
	ChargeMasterHostEnvVarName = "CHARGE_MASTER_API_HOST"

	// ChargeMasterAPISchemeEnvVarName points at an environment variable that
	// indicates whether the API is "http" or "https". It is used when our code
	// needs to construct custom API paths from scratch.
	ChargeMasterAPISchemeEnvVarName = "CHARGE_MASTER_API_SCHEME"

	// ChargeMasterTokenURLEnvVarName is an environment variable that contains
	// the path to the OAuth 2 token URL for the charge master base. This URL
	// could be the same as that used by other Slade 360 products e.g EDI.
	// It could also be different.
	ChargeMasterTokenURLEnvVarName = "CHARGE_MASTER_TOKEN_URL"

	// ChargeMasterClientIDEnvVarName is the name of an environment variable that holds
	// the OAuth2 client ID for a charge master API application.
	ChargeMasterClientIDEnvVarName = "CHARGE_MASTER_CLIENT_ID"

	// ChargeMasterClientSecretEnvVarName is the name of an environment variable that holds
	// the OAuth2 client secret for a charge master API application.
	ChargeMasterClientSecretEnvVarName = "CHARGE_MASTER_CLIENT_SECRET"

	// ChargeMasterUsernameEnvVarName is the name of an environment variable that holds the
	// username of a charge master API user.
	ChargeMasterUsernameEnvVarName = "CHARGE_MASTER_USERNAME"

	// ChargeMasterPasswordEnvVarName is the name of an environment variable that holds the
	// password of the charge master API user referred to by `ChargeMasterUsernameEnvVarName`.
	ChargeMasterPasswordEnvVarName = "CHARGE_MASTER_PASSWORD"

	// ChargeMasterGrantTypeEnvVarName should be "password" i.e the only type of OAuth 2
	// "application" that will work for this client is a confidential one that supports
	// password authentication.
	ChargeMasterGrantTypeEnvVarName = "CHARGE_MASTER_GRANT_TYPE"
)

// NewChargeMasterClient initializes a new charge master client from the environment.
// It assumes that the environment variables were confirmed to be present during
// server initialization. For that reason, it will panic if an environment variable is
// unexpectedly absent.
func NewChargeMasterClient() (*base.ServerClient, error) {
	clientID := base.MustGetEnvVar(ChargeMasterClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(ChargeMasterClientSecretEnvVarName)
	apiTokenURL := base.MustGetEnvVar(ChargeMasterTokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(ChargeMasterHostEnvVarName)
	apiScheme := base.MustGetEnvVar(ChargeMasterAPISchemeEnvVarName)
	grantType := base.MustGetEnvVar(ChargeMasterGrantTypeEnvVarName)
	username := base.MustGetEnvVar(ChargeMasterUsernameEnvVarName)
	password := base.MustGetEnvVar(ChargeMasterPasswordEnvVarName)
	extraHeaders := make(map[string]string)
	return base.NewServerClient(
		clientID, clientSecret, apiTokenURL, apiHost, apiScheme, grantType, username, password, extraHeaders)
}
