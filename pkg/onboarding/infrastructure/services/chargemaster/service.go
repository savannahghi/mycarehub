package chargemaster

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

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

	// ChargeMasterBusinessPartnerPath endpoint for business partners on charge master
	ChargeMasterBusinessPartnerPath = "/v1/business_partners/"
)

// ServiceChargeMaster represents logic required to communicate with chargemaster
type ServiceChargeMaster interface {
	FetchChargeMasterClient() *base.ServerClient
	FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BusinessPartnerFilterInput,
		sort []*resources.BusinessPartnerSortInput) (*resources.BusinessPartnerConnection, error)
	FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BranchFilterInput,
		sort []*resources.BranchSortInput) (*resources.BranchConnection, error)
	FetchProviderByID(ctx context.Context, id string) (*domain.BusinessPartner, error)
}

// ServiceChargeMasterImpl ..
type ServiceChargeMasterImpl struct {
	ChargeMasterClient *base.ServerClient
}

// NewChargeMasterUseCasesImpl ...
func NewChargeMasterUseCasesImpl(r repository.OnboardingRepository) ServiceChargeMaster {

	clientID := base.MustGetEnvVar(ChargeMasterClientIDEnvVarName)
	clientSecret := base.MustGetEnvVar(ChargeMasterClientSecretEnvVarName)
	apiTokenURL := base.MustGetEnvVar(ChargeMasterTokenURLEnvVarName)
	apiHost := base.MustGetEnvVar(ChargeMasterHostEnvVarName)
	apiScheme := base.MustGetEnvVar(ChargeMasterAPISchemeEnvVarName)
	grantType := base.MustGetEnvVar(ChargeMasterGrantTypeEnvVarName)
	username := base.MustGetEnvVar(ChargeMasterUsernameEnvVarName)
	password := base.MustGetEnvVar(ChargeMasterPasswordEnvVarName)
	extraHeaders := make(map[string]string)
	client, err := base.NewServerClient(
		clientID, clientSecret, apiTokenURL, apiHost, apiScheme, grantType, username, password, extraHeaders)
	if err != nil {
		log.Panicf("unable to initialize Chargemaster client for profile service: %s", err)
		os.Exit(1)
	}

	return &ServiceChargeMasterImpl{ChargeMasterClient: client}
}

// FetchChargeMasterClient ...
func (chr ServiceChargeMasterImpl) FetchChargeMasterClient() *base.ServerClient {
	return chr.ChargeMasterClient
}

// FindProvider search for a provider in chargemaster using their name
//
// Example https://base.chargemaster.slade360emr.com/v1/business_partners/?bp_type=PROVIDER&search={name}
func (chr ServiceChargeMasterImpl) FindProvider(ctx context.Context, pagination *base.PaginationInput,
	filter []*resources.BusinessPartnerFilterInput, sort []*resources.BusinessPartnerSortInput) (*resources.BusinessPartnerConnection, error) {

	paginationParams, err := base.GetAPIPaginationParams(pagination)
	if err != nil {
		return nil, err
	}

	defaultParams := url.Values{}
	defaultParams.Add("fields", "id,name,slade_code,parent")
	defaultParams.Add("is_active", "True")
	defaultParams.Add("bp_type", "PROVIDER")

	queryParams := []url.Values{defaultParams, paginationParams}
	for _, fp := range filter {
		queryParams = append(queryParams, fp.ToURLValues())
	}
	for _, fp := range sort {
		queryParams = append(queryParams, fp.ToURLValues())
	}

	mergedParams := base.MergeURLValues(queryParams...)
	queryFragment := mergedParams.Encode()

	type apiResp struct {
		base.SladeAPIListRespBase

		Results []*domain.BusinessPartner `json:"results,omitempty"`
	}

	r := apiResp{}
	err = base.ReadRequestToTarget(chr.FetchChargeMasterClient(), "GET", ChargeMasterBusinessPartnerPath, queryFragment, nil, &r)
	if err != nil {
		return nil, err
	}

	startOffset := base.CreateAndEncodeCursor(r.StartIndex)
	endOffset := base.CreateAndEncodeCursor(r.EndIndex)
	hasNextPage := r.Next != ""
	hasPreviousPage := r.Previous != ""

	edges := []*resources.BusinessPartnerEdge{}
	for pos, org := range r.Results {
		edge := &resources.BusinessPartnerEdge{
			Node: &domain.BusinessPartner{
				ID:        org.ID,
				Name:      org.Name,
				SladeCode: org.SladeCode,
				Parent:    org.Parent,
			},
			Cursor: base.CreateAndEncodeCursor(pos + 1),
		}
		edges = append(edges, edge)
	}
	pageInfo := &base.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		StartCursor:     startOffset,
		EndCursor:       endOffset,
	}
	connection := &resources.BusinessPartnerConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
	return connection, nil
}

// FetchProviderByID returns details of a specific provider given the ID
func (chr ServiceChargeMasterImpl) FetchProviderByID(ctx context.Context, id string) (*domain.BusinessPartner, error) {

	partner := &domain.BusinessPartner{}
	BusinessPartnerPath := fmt.Sprintf("%v%v/", ChargeMasterBusinessPartnerPath, id)

	err := base.ReadRequestToTarget(chr.FetchChargeMasterClient(), "GET", BusinessPartnerPath, "", nil, partner)
	if err != nil {
		return nil, err
	}

	// check if partner is empty
	if (*partner == domain.BusinessPartner{}) {
		return nil, fmt.Errorf("business partner not found. invalid id: %v", id)
	}

	return partner, nil
}

// FindBranch lists all locations known to Slade 360 Charge Master
// Example URL: https://base.chargemaster.slade360emr.com/v1/business_partners/?format=json&page_size=100&parent=6ba48d97-93d2-4815-a447-f51240cbcab8&fields=id,name,slade_code
func (chr ServiceChargeMasterImpl) FindBranch(ctx context.Context, pagination *base.PaginationInput,
	filter []*resources.BranchFilterInput, sort []*resources.BranchSortInput) (*resources.BranchConnection, error) {

	paginationParams, err := base.GetAPIPaginationParams(pagination)
	if err != nil {
		return nil, err
	}
	defaultParams := url.Values{}
	defaultParams.Add("fields", "id,name,slade_code")
	defaultParams.Add("is_active", "True")
	defaultParams.Add("is_branch", "True")

	queryParams := []url.Values{defaultParams, paginationParams}
	for _, fp := range filter {
		queryParams = append(queryParams, fp.ToURLValues())
	}
	for _, fp := range sort {
		queryParams = append(queryParams, fp.ToURLValues())
	}
	mergedParams := base.MergeURLValues(queryParams...)
	queryFragment := mergedParams.Encode()

	type apiResp struct {
		base.SladeAPIListRespBase

		Results []*domain.BusinessPartner `json:"results,omitempty"`
	}

	r := apiResp{}
	err = base.ReadRequestToTarget(chr.FetchChargeMasterClient(), "GET", ChargeMasterBusinessPartnerPath, queryFragment, nil, &r)
	if err != nil {
		return nil, err
	}
	startOffset := base.CreateAndEncodeCursor(r.StartIndex)
	endOffset := base.CreateAndEncodeCursor(r.EndIndex)
	hasNextPage := r.Next != ""
	hasPreviousPage := r.Previous != ""

	edges := []*resources.BranchEdge{}
	for pos, branch := range r.Results {
		orgSladeCode, err := parentOrgSladeCodeFromBranch(branch)
		if err != nil {
			return nil, err
		}

		edge := &resources.BranchEdge{
			Node: &domain.Branch{
				ID:                    branch.ID,
				Name:                  branch.Name,
				BranchSladeCode:       branch.SladeCode,
				OrganizationSladeCode: orgSladeCode,
			},
			Cursor: base.CreateAndEncodeCursor(pos + 1),
		}
		edges = append(edges, edge)
	}
	pageInfo := &base.PageInfo{
		HasNextPage:     hasNextPage,
		HasPreviousPage: hasPreviousPage,
		StartCursor:     startOffset,
		EndCursor:       endOffset,
	}
	connection := &resources.BranchConnection{
		Edges:    edges,
		PageInfo: pageInfo,
	}
	return connection, nil
}

func parentOrgSladeCodeFromBranch(branch *domain.BusinessPartner) (string, error) {
	if !strings.HasPrefix(branch.SladeCode, "BRA-") {
		return "", fmt.Errorf("%s is not a valid branch Slade Code; expected a BRA- prefix", branch.SladeCode)
	}
	trunc := strings.TrimPrefix(branch.SladeCode, "BRA-")
	split := strings.Split(trunc, "-")
	if len(split) != 3 {
		return "", fmt.Errorf("expected the branch Slade Code to split into 3 parts on -; got %s", split)
	}
	orgParts := split[0:2]
	orgSladeCode := strings.Join(orgParts, "-")
	return orgSladeCode, nil
}
