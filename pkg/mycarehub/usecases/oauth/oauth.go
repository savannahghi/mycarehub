package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/oauth/storage"
	"github.com/savannahghi/serverutils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	secret = serverutils.MustGetEnvVar("FOSITE_SECRET")
	// debug, _     = serverutils.GetEnvVar(serverutils.DebugEnvVarName)
	clientID     = serverutils.MustGetEnvVar("MYCAREHUB_CLIENT_ID")
	clientSecret = serverutils.MustGetEnvVar("MYCAREHUB_CLIENT_SECRET")
	tokenURL     = serverutils.MustGetEnvVar("MYCAREHUB_TOKEN_URL")
)

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesOauth interface {
	CreateOauthClient(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error)
	FositeProvider() fosite.OAuth2Provider
	GenerateUserAuthTokens(ctx context.Context, userID string) (*AuthTokens, error)
	RefreshAutToken(ctx context.Context, refreshToken string) (*AuthTokens, error)
}

// UseCasesOauthImpl represents oauth implementation
type UseCasesOauthImpl struct {
	update   infrastructure.Update
	query    infrastructure.Query
	create   infrastructure.Create
	delete   infrastructure.Delete
	provider fosite.OAuth2Provider
}

// NewUseCasesOauthImplementation initializes an implementation of the fosite storage
func NewUseCasesOauthImplementation(create infrastructure.Create, update infrastructure.Update, query infrastructure.Query, delete infrastructure.Delete) UseCasesOauthImpl {
	// var debugEnv bool
	// debugEnv, err := strconv.ParseBool(debug)
	// if err != nil {
	// 	debugEnv = true
	// }

	conf := &fosite.Config{
		GlobalSecret: []byte(secret),

		AccessTokenLifespan: 1 * time.Hour,

		RefreshTokenLifespan: 24 * time.Hour,
		RefreshTokenScopes:   []string{},

		AuthorizeCodeLifespan: 5 * time.Minute,

		SendDebugMessagesToClients: true,
	}

	storage := storage.NewFositeStorage(create, update, query, delete)

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
		OAuth2InternalGrantFactory,
	)

	return UseCasesOauthImpl{
		update:   update,
		query:    query,
		create:   create,
		delete:   delete,
		provider: provider,
	}
}

func (u UseCasesOauthImpl) FositeProvider() fosite.OAuth2Provider {
	return u.provider
}

type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (u UseCasesOauthImpl) GenerateUserAuthTokens(ctx context.Context, userID string) (*AuthTokens, error) {
	client, err := u.getOrCreateInternalCLient(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	extraDetails := map[string]interface{}{
		"user_id": *user.ID,
	}

	session := domain.NewSession(ctx, client.ID, *user.ID, user.Username, user.Name, extraDetails)
	request := fosite.NewAccessRequest(session)
	request.GrantTypes = []string{"internal"}
	request.Client = client

	response, err := u.provider.NewAccessResponse(ctx, request)
	if err != nil {
		return nil, err
	}

	resp := response.ToMap()

	expires := resp["expires_in"].(int64)

	tokens := AuthTokens{
		RefreshToken: resp["refresh_token"].(string),
		AccessToken:  resp["access_token"].(string),
		ExpiresIn:    int(expires),
	}

	return &tokens, nil
}

func (u UseCasesOauthImpl) getOrCreateInternalCLient(ctx context.Context) (*domain.OauthClient, error) {
	client, err := u.query.GetOauthClient(ctx, clientID)
	if err == nil {
		return client, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	input := dto.OauthClientInput{
		Name:          "Mycarehub",
		Secret:        clientSecret,
		RedirectURIs:  []string{},
		ResponseTypes: []string{"token"},
		Grants:        []string{"internal", "refresh_token"},
	}

	secret, err := bcrypt.GenerateFromPassword([]byte(input.Secret), fosite.DefaultBCryptWorkFactor)
	if err != nil {
		return nil, err
	}

	defaultClient := &domain.OauthClient{
		ID:                      clientID,
		Name:                    input.Name,
		Secret:                  string(secret),
		RedirectURIs:            input.RedirectURIs,
		Active:                  true,
		Grants:                  input.Grants,
		ResponseTypes:           input.ResponseTypes,
		TokenEndpointAuthMethod: "client_secret_basic",
	}

	err = u.create.CreateOauthClient(ctx, defaultClient)
	if err != nil {
		return nil, err
	}

	return defaultClient, nil
}

// CreateOauthClient is the resolver for the createOauthClient field.
func (u UseCasesOauthImpl) CreateOauthClient(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error) {
	secret, err := bcrypt.GenerateFromPassword([]byte(input.Secret), fosite.DefaultBCryptWorkFactor)
	if err != nil {
		return nil, err
	}

	client := &domain.OauthClient{
		Name:                    input.Name,
		Secret:                  string(secret),
		RedirectURIs:            input.RedirectURIs,
		Active:                  true,
		Grants:                  input.Grants,
		ResponseTypes:           input.ResponseTypes,
		TokenEndpointAuthMethod: "client_secret_basic",
	}

	err = u.create.CreateOauthClient(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// ListOauthClients is the resolver for the listOauthClients field.
func (u UseCasesOauthImpl) ListOauthClients(ctx context.Context) ([]*domain.OauthClient, error) {
	return nil, nil
}

// RefreshAutToken is the resolver for the listOauthClients field.
func (u UseCasesOauthImpl) RefreshAutToken(ctx context.Context, refreshToken string) (*AuthTokens, error) {
	formData := url.Values{
		"refresh_token": []string{refreshToken},
		"grant_type":    []string{"refresh_token"},
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to refresh token")
		return nil, err
	}

	var tokens AuthTokens

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bs, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}
