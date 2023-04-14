package oauth

import (
	"strconv"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/savannahghi/serverutils"
)

var (
	secret   = serverutils.MustGetEnvVar("FOSITE_SECRET")
	debug, _ = serverutils.GetEnvVar(serverutils.DebugEnvVarName)
)

func (u UseCasesOauthImpl) FositeProvider() fosite.OAuth2Provider {
	var debugEnv bool
	debugEnv, err := strconv.ParseBool(debug)
	if err != nil {
		debugEnv = false
	}

	conf := &fosite.Config{
		GlobalSecret: []byte(secret),

		AccessTokenLifespan:   1 * time.Hour,
		RefreshTokenLifespan:  24 * time.Hour,
		AuthorizeCodeLifespan: 5 * time.Minute,

		SendDebugMessagesToClients: debugEnv,
	}

	storage := u

	provider := compose.Compose(
		conf,
		storage,
		compose.NewOAuth2HMACStrategy(conf),
		compose.OAuth2AuthorizeExplicitFactory,
		compose.OAuth2AuthorizeImplicitFactory,
		compose.OAuth2ClientCredentialsGrantFactory,
		compose.OAuth2RefreshTokenGrantFactory,
		compose.OAuth2TokenIntrospectionFactory,
		compose.OAuth2TokenRevocationFactory,
	)

	return provider
}
