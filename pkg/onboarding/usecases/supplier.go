package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"

	"github.com/cenkalti/backoff"
	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/config/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	supplierAPIPath        = "/api/business_partners/suppliers/"
	customerAPIPath        = "/api/business_partners/customers/"
	active                 = true
	country                = "KEN" // Anticipate worldwide expansion
	supplierCollectionName = "suppliers"
	futureHours            = 878400
	savannahSladeCode      = "1"
)

// SupplierUseCases represent the business logic required for management of suppliers
type SupplierUseCases interface {
	AddPartnerType(ctx context.Context, name *string, partnerType *domain.PartnerType) (bool, error)

	AddCustomerSupplierERPAccount(ctx context.Context, name string, partnerType domain.PartnerType) (*domain.Supplier, error)

	FindSupplierByID(ctx context.Context, id string) (*domain.Supplier, error)

	FindSupplierByUID(ctx context.Context) (*domain.Supplier, error)

	SetUpSupplier(ctx context.Context, accountType domain.AccountType) (*domain.Supplier, error)

	SuspendSupplier(ctx context.Context) (bool, error)

	EDIUserLogin(username, password *string) (*base.EDIUserProfile, error)

	CoreEDIUserLogin(username, password string) (*base.EDIUserProfile, error)

	FetchSupplierAllowedLocations(ctx context.Context) (*domain.BranchConnection, error)

	AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error)

	AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error)

	AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error)

	AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error)

	AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error)

	AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error)

	AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error)

	AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error)

	AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error)

	AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error)

	AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error)

	FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error)

	SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*domain.BranchConnection, error)

	SupplierSetDefaultLocation(ctx context.Context, locationID string) (bool, error)

	// FetchSupplierAllowedLocations(ctx context.Context) (*domain.BranchConnection, error)
	// SaveProfileNudge(nudge map[string]interface{}) error
	// PublishKYCNudge(uid string, partner *domain.PartnerType, account *domain.AccountType) error
	// PublishKYCFeedItem(ctx context.Context, uids ...string) error
	// StageKYCProcessingRequest(sup *domain.Supplier) error

	// SaveKYCResponse(ctx context.Context, kycJSON []byte, supplier *domain.Supplier, dsnap *firestore.DocumentSnapshot) error

	// ProcessKYCRequest(ctx context.Context, id string, status domain.KYCProcessStatus, rejectionReason *string) (bool, error)
	// SendKYCEmail(ctx context.Context, text, emailaddress string) error
}

// SupplierUseCasesImpl represents usecase implementation object
type SupplierUseCasesImpl struct {
	repo         repository.OnboardingRepository
	profile      ProfileUseCase
	erp          erp.Service
	chargemaster ChargeMasterUseCases
	engagement   engagement.Service
}

// NewSupplierUseCases returns a new a onboarding usecase
func NewSupplierUseCases(r repository.OnboardingRepository, p ProfileUseCase,
	er erp.Service, chrg ChargeMasterUseCases,
	eng engagement.Service) SupplierUseCases {
	return &SupplierUseCasesImpl{repo: r, profile: p,
		erp: er, chargemaster: chrg, engagement: eng}
}

// AddPartnerType create the initial supplier record
func (s SupplierUseCasesImpl) AddPartnerType(ctx context.Context, name *string, partnerType *domain.PartnerType) (bool, error) {

	if name == nil || partnerType == nil || !partnerType.IsValid() {
		return false, fmt.Errorf("expected `name` to be defined and `partnerType` to be valid")
	}

	if !partnerType.IsValid() {
		return false, fmt.Errorf("invalid `partnerType` provided")
	}

	if *partnerType == domain.PartnerTypeConsumer {
		return false, fmt.Errorf("invalid `partnerType`. cannot use CONSUMER in this context")
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid)
	if err != nil {
		return false, fmt.Errorf("unable to read user profile: %w", err)
	}

	v, err := s.repo.AddPartnerType(ctx, profile.ID, name, partnerType)
	if !v || err != nil {
		return false, fmt.Errorf("error occured while adding partner type: %w", err)
	}

	return true, nil
}

// AddCustomerSupplierERPAccount makes a call to our own ERP and creates a  customer account or supplier account  based
// on the provided partnerType
func (s SupplierUseCasesImpl) AddCustomerSupplierERPAccount(ctx context.Context, name string, partnerType domain.PartnerType) (*domain.Supplier, error) {

	userUID, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %v", err)
	}

	profile, err := s.profile.GetProfileByID(ctx, userUID)
	if err != nil {
		return nil, fmt.Errorf("unable to read user profile: %w", err)
	}

	currency, err := base.FetchDefaultCurrency(s.erp.FetchERPClient())
	if err != nil {
		return nil, fmt.Errorf("unable to fetch orgs default currency: %v", err)
	}

	validPartnerType := partnerType.IsValid()
	if !validPartnerType {
		return nil, fmt.Errorf("%v is not an valid partner type choice", partnerType.String())
	}

	var payload map[string]interface{}
	var endpoint string

	if partnerType == domain.PartnerTypeConsumer {
		endpoint = customerAPIPath
		payload = map[string]interface{}{
			"active":        active,
			"partner_name":  name,
			"country":       country,
			"currency":      *currency.ID,
			"is_customer":   true,
			"customer type": partnerType,
		}
	} else {
		endpoint = supplierAPIPath
		payload = map[string]interface{}{
			"active":        active,
			"partner_name":  name,
			"country":       country,
			"currency":      *currency.ID,
			"is_supplier":   true,
			"supplier_type": partnerType,
		}
	}

	if err := s.erp.CreateERPSupplier(string(http.MethodPost), endpoint, payload, partnerType); err != nil {
		return nil, err
	}

	// for customers, we don't return anything. So long as there is not error, we are good
	if partnerType == domain.PartnerTypeConsumer {
		return nil, nil
	}

	return s.repo.ActivateSupplierProfile(ctx, profile.ID)
}

// FindSupplierByID fetches a supplier by their id
func (s SupplierUseCasesImpl) FindSupplierByID(ctx context.Context, id string) (*domain.Supplier, error) {
	return s.repo.GetSupplierProfileByID(ctx, id)
}

// FindSupplierByUID fetches a supplier by logged in user uid
func (s SupplierUseCasesImpl) FindSupplierByUID(ctx context.Context) (*domain.Supplier, error) {
	pr, err := s.profile.UserProfile(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSupplierProfileByProfileID(ctx, pr.ID)

}

// SetUpSupplier performs initial account set up during onboarding
func (s SupplierUseCasesImpl) SetUpSupplier(ctx context.Context, accountType domain.AccountType) (*domain.Supplier, error) {

	validAccountType := accountType.IsValid()
	if !validAccountType {
		return nil, fmt.Errorf("%v is not an allowed AccountType choice", accountType.String())
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
	}

	supplier, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	supplier.AccountType = accountType
	supplier.UnderOrganization = false
	supplier.IsOrganizationVerified = false
	supplier.HasBranches = false
	supplier.Active = false

	sup, err := s.repo.UpdateSupplierProfile(ctx, supplier)
	if err != nil {
		return nil, err
	}

	go func() {
		op := func() error {
			return s.PublishKYCNudge(ctx, uid, &supplier.PartnerType, &supplier.AccountType)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}()

	return sup, nil
}

// SuspendSupplier flips the active boolean on the erp partner from true to false
func (s SupplierUseCasesImpl) SuspendSupplier(ctx context.Context) (bool, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	sup.Active = false

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return false, err
	}

	//todo(dexter) notify the supplier of the suspension

	return true, nil

}

// EDIUserLogin used to login a user to EDI (Portal Authserver) and return their
// EDI (Portal Authserver) profile
func (s SupplierUseCasesImpl) EDIUserLogin(username, password *string) (*base.EDIUserProfile, error) {

	if username == nil || password == nil {
		return nil, fmt.Errorf("invalid credentials, expected a username AND password")
	}

	ediClient, err := base.LoginClient(*username, *password)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := base.FetchUserProfile(ediClient)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
	}

	return userProfile, nil

}

// CoreEDIUserLogin used to login a user to EDI (Core Authserver) and return their EDI
// EDI (Core Authserver) profile
func (s SupplierUseCasesImpl) CoreEDIUserLogin(username, password string) (*base.EDIUserProfile, error) {

	if username == "" || password == "" {
		return nil, fmt.Errorf("invalid credentials, expected a username AND password")
	}

	ediClient, err := utils.LoginClient(username, password)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := base.FetchUserProfile(ediClient)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
	}

	return userProfile, nil

}

// SupplierEDILogin it used to instantiate as call when setting up a supplier's account's who
// has an affliation to a provider with the slade ecosystem. The logic is as follows;
// 1 . login to the relevant edi to assert the user has an account
// 2 . fetch the branches of the provider given the slade code which we have
// 3 . update the user's supplier record
// 4. return the list of branches to the frontend so that a default location can be set
func (s SupplierUseCasesImpl) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*domain.BranchConnection, error) {

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user: %w", err)
	}

	supplier, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	supplier.AccountType = domain.AccountTypeIndividual
	supplier.UnderOrganization = true

	ediUserProfile, err := func(sladeCode string) (*base.EDIUserProfile, error) {
		var ediUserProfile *base.EDIUserProfile
		var err error

		switch sladeCode {
		case savannahSladeCode:
			// login to core
			ediUserProfile, err = s.CoreEDIUserLogin(username, password)
			if err != nil {
				supplier.IsOrganizationVerified = false
				return nil, fmt.Errorf("cannot get edi user profile: %w", err)
			}

			if ediUserProfile == nil {
				return nil, fmt.Errorf("edi user profile not found")
			}

		default:
			//Login to portal edi
			ediUserProfile, err = s.EDIUserLogin(&username, &password)
			if err != nil {
				supplier.IsOrganizationVerified = false
				return nil, fmt.Errorf("cannot get edi user profile: %w", err)
			}

			if ediUserProfile == nil {
				return nil, fmt.Errorf("edi user profile not found")
			}

		}
		return ediUserProfile, nil
	}(sladeCode)

	if err != nil {
		return nil, err
	}

	pageInfo := &base.PageInfo{
		HasNextPage:     false,
		HasPreviousPage: false,
		StartCursor:     nil,
		EndCursor:       nil,
	}

	// The slade code comes in the form 'PRO-1234' or 'BRA-PRO-1234-1'
	// or a single code '1234'
	// we split it to get the interger part of the slade code.
	var orgSladeCode string
	if strings.HasPrefix(sladeCode, "BRA-") {
		orgSladeCode = strings.Split(sladeCode, "-")[2]
	} else if strings.HasPrefix(sladeCode, "PRO-") {
		orgSladeCode = strings.Split(sladeCode, "-")[1]
	} else {
		orgSladeCode = sladeCode
	}

	if orgSladeCode == savannahSladeCode {
		// profile.Permissions = base.DefaultAdminPermissions
		// todo(dexter) add call to update profile permissions in base

		supplier.EDIUserProfile = ediUserProfile
		supplier.IsOrganizationVerified = true
		supplier.SladeCode = sladeCode
		supplier.Active = true
		supplier.KYCSubmitted = true
		supplier.PartnerSetupComplete = true

		_, err := s.repo.UpdateSupplierProfile(ctx, supplier)
		if err != nil {
			return nil, err
		}

		return &domain.BranchConnection{PageInfo: pageInfo}, nil
	}

	// verify slade code.
	if ediUserProfile.BusinessPartner != orgSladeCode {
		return nil, fmt.Errorf("invalid slade code for selected provider: %v, got: %v", sladeCode, ediUserProfile.BusinessPartner)
	}

	supplier.EDIUserProfile = ediUserProfile
	supplier.IsOrganizationVerified = true
	supplier.SladeCode = sladeCode

	filter := []*domain.BusinessPartnerFilterInput{
		{
			SladeCode: &sladeCode,
		},
	}

	partner, err := s.chargemaster.FindProvider(ctx, nil, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch organization branches location: %v", err)
	}

	var businessPartner domain.BusinessPartner

	if len(partner.Edges) != 1 {
		return nil, fmt.Errorf("expected one business partner, found: %v", len(partner.Edges))
	}

	businessPartner = *partner.Edges[0].Node
	var brFilter []*domain.BranchFilterInput

	go func() {
		op := func() error {
			return s.PublishKYCNudge(ctx, uid, &supplier.PartnerType, &supplier.AccountType)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}()

	if businessPartner.Parent != nil {
		supplier.HasBranches = true
		supplier.ParentOrganizationID = *businessPartner.Parent
		filter := &domain.BranchFilterInput{
			ParentOrganizationID: businessPartner.Parent,
		}

		brFilter = append(brFilter, filter)

		_, err := s.repo.UpdateSupplierProfile(ctx, supplier)
		if err != nil {
			return nil, err
		}

		return s.chargemaster.FindBranch(ctx, nil, brFilter, nil)
	}
	loc := domain.Location{
		ID:   businessPartner.ID,
		Name: businessPartner.Name,
	}
	supplier.Location = &loc

	_, err = s.repo.UpdateSupplierProfile(ctx, supplier)
	if err != nil {
		return nil, err
	}

	return &domain.BranchConnection{PageInfo: pageInfo}, nil
}

// SupplierSetDefaultLocation updates the default location ot the supplier by the given location id
func (s SupplierUseCasesImpl) SupplierSetDefaultLocation(ctx context.Context, locationID string) (bool, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return false, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	// fetch the branches of the provider filtered by sladecode and ParentOrganizationID
	filter := []*domain.BranchFilterInput{
		{
			SladeCode:            &sup.SladeCode,
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		return false, fmt.Errorf("unable to fetch organization branches location: %v", err)
	}

	branch := func(brs *domain.BranchConnection, location string) *domain.BranchEdge {
		for _, b := range brs.Edges {
			if b.Node.ID == location {
				return b
			}
		}
		return nil
	}(brs, locationID)

	if branch != nil {
		loc := domain.Location{
			ID:              branch.Node.ID,
			Name:            branch.Node.Name,
			BranchSladeCode: &branch.Node.BranchSladeCode,
		}
		sup.Location = &loc

		_, err = s.repo.UpdateSupplierProfile(ctx, sup)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, fmt.Errorf("unable to get location of id %v : %v", locationID, err)
}

// FetchSupplierAllowedLocations retrieves all the locations that the user in context can work on.
func (s *SupplierUseCasesImpl) FetchSupplierAllowedLocations(ctx context.Context) (*domain.BranchConnection, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	// fetch the branches of the provider filtered by sladecode and ParentOrganizationID
	filter := []*domain.BranchFilterInput{
		{
			SladeCode:            &sup.SladeCode,
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch organization branches location: %v", err)
	}

	return brs, nil
}

// PublishKYCNudge pushes a kyc nudge to the user feed
func (s *SupplierUseCasesImpl) PublishKYCNudge(ctx context.Context, uid string, partner *domain.PartnerType, account *domain.AccountType) error {

	if partner == nil || !partner.IsValid() {
		return fmt.Errorf("expected `partner` to be defined and to be valid")
	}

	if *partner == domain.PartnerTypeConsumer {
		return fmt.Errorf("invalid `partner`. cannot use CONSUMER in this context")
	}

	if !account.IsValid() {
		return fmt.Errorf("provided `account` is not valid")
	}

	payload := base.Nudge{
		ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
		SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
		Visibility:     "SHOW",
		Status:         "PENDING",
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Title:          fmt.Sprintf("Complete your %v KYC", strings.ToLower(partner.String())),
		Text:           "Fill in your Be.Well business KYC in order to start transacting",
		Links: []base.Link{
			{
				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC",
				Description: fmt.Sprintf("KYC for %v", partner.String()),
				Thumbnail:   base.LogoURL,
			},
		},
		Actions: []base.Action{
			{
				ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
				SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
				Name:           strings.ToUpper(fmt.Sprintf("COMPLETE_%v_%v_KYC", account.String(), partner.String())),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
				Icon: base.Link{
					ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
					URL:         base.LogoURL,
					LinkType:    base.LinkTypePngImage,
					Title:       fmt.Sprintf("Complete your %v KYC", strings.ToLower(partner.String())),
					Description: "Fill in your Be.Well business KYC in order to start transacting",
					Thumbnail:   base.LogoURL,
				},
			},
		},
		Users:                []string{uid},
		Groups:               []string{uid},
		NotificationChannels: []base.Channel{base.ChannelEmail, base.ChannelFcm},
	}

	resp, err := s.engagement.PublishKYCNudge(uid, payload)
	if err != nil {
		return fmt.Errorf("unable to publish kyc nudge : %v", err)
	}

	//TODO(dexter) to be removed. Just here for debug
	res, _ := httputil.DumpResponse(resp, true)
	log.Println(string(res))

	if resp.StatusCode != http.StatusOK {
		// stage the nudge
		stage := func(pl base.Nudge) error {
			k, err := json.Marshal(payload)
			if err != nil {
				return fmt.Errorf("cannot marshal payload to json")
			}

			var kMap map[string]interface{}
			err = json.Unmarshal(k, &kMap)
			if err != nil {
				return fmt.Errorf("cannot unmarshal payload from json")
			}

			if err := s.SaveProfileNudge(ctx, kMap); err != nil {
				logrus.Errorf("failed to stage nudge : %v", err)
			}
			return nil

		}(payload)

		if err := stage; err != nil {
			logrus.Errorf("failed to stage nudge : %v", err)
		}
		return fmt.Errorf("unable to publish kyc nudge. unexpected status code  %v", resp.StatusCode)
	}

	return nil

}

// PublishKYCFeedItem notifies admin users of a KYC approval request
func (s SupplierUseCasesImpl) PublishKYCFeedItem(ctx context.Context, uids ...string) error {

	for _, uid := range uids {
		payload := base.Item{
			ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
			SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
			Expiry:         time.Now().Add(time.Hour * futureHours),
			Persistent:     true,
			Status:         base.StatusPending,
			Visibility:     base.VisibilityShow,
			Author:         "Be.Well Team",
			Label:          "KYC",
			Tagline:        "Process incoming KYC",
			Text:           "Review KYC for the partner and either approve or reject",
			TextType:       base.TextTypeMarkdown,
			Icon: base.Link{
				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC Review",
				Description: "Review KYC for the partner and either approve or reject",
				Thumbnail:   base.LogoURL,
			},
			Timestamp: time.Now(),
			Actions: []base.Action{
				{
					ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
					SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
					Name:           "Review KYC details",
					Icon: base.Link{
						ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
						URL:         base.LogoURL,
						LinkType:    base.LinkTypePngImage,
						Title:       "Review KYC details",
						Description: "Review and approve or reject KYC details for the supplier",
						Thumbnail:   base.LogoURL,
					},
					ActionType:     base.ActionTypePrimary,
					Handling:       base.HandlingFullPage,
					AllowAnonymous: false,
				},
			},
			Links: []base.Link{
				{
					ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
					URL:         base.LogoURL,
					LinkType:    base.LinkTypePngImage,
					Title:       "KYC process request",
					Description: "Process KYC request",
					Thumbnail:   base.LogoURL,
				},
			},

			Summary: "Process incoming KYC",
			Users:   uids,
			NotificationChannels: []base.Channel{
				base.ChannelFcm,
				base.ChannelEmail,
				base.ChannelSms,
			},
		}

		resp, err := s.engagement.PublishKYCFeedItem(uid, payload)
		if err != nil {
			return fmt.Errorf("unable to publish kyc admin notification feed item : %v", err)
		}

		//TODO(dexter) to be removed. Just here for debug
		res, _ := httputil.DumpResponse(resp, true)
		log.Println(string(res))

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("unable to publish kyc admin notification feed item. unexpected status code  %v", resp.StatusCode)
		}
	}

	return nil
}

// SaveProfileNudge stages nudges published from this service. These nudges will be
// referenced later to support some specialized use-case. A nudge will be uniquely
// identified by its id and sequenceNumber
func (s *SupplierUseCasesImpl) SaveProfileNudge(ctx context.Context, nudge map[string]interface{}) error {
	return s.repo.StageProfileNudge(ctx, nudge)
}

func (s *SupplierUseCasesImpl) parseKYCAsMap(data interface{}) (map[string]interface{}, error) {
	kycJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	var kycAsMap map[string]interface{}

	if err := json.Unmarshal(kycJSON, &kycAsMap); err != nil {
		return nil, fmt.Errorf("cannot unmarshal kyc from json")
	}

	return kycAsMap, nil
}

// AddIndividualRiderKyc adds KYC for an individual rider
func (s *SupplierUseCasesImpl) AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.IndividualRider{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           input.IdentificationDoc.IdentificationDocType,
			IdentificationDocNumber:         input.IdentificationDoc.IdentificationDocNumber,
			IdentificationDocNumberUploadID: input.IdentificationDoc.IdentificationDocNumberUploadID,
		},
		KRAPIN:                         input.KRAPIN,
		KRAPINUploadID:                 input.KRAPINUploadID,
		DrivingLicenseID:               input.DrivingLicenseID,
		DrivingLicenseUploadID:         input.DrivingLicenseUploadID,
		CertificateGoodConductUploadID: input.CertificateGoodConductUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddOrganizationRiderKyc adds KYC for an organization rider
func (s *SupplierUseCasesImpl) AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error) {

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.OrganizationRider{
		OrganizationTypeName:               input.OrganizationTypeName,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []domain.Identification) []domain.Identification {
			pl := []domain.Identification{}
			for _, i := range p {
				pl = append(pl, domain.Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,

		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddIndividualPractitionerKyc adds KYC for an individual pratitioner
func (s *SupplierUseCasesImpl) AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error) {

	for _, p := range input.PracticeServices {
		if !p.IsValid() {
			return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
		}
	}

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.IndividualPractitioner{

		IdentificationDoc: func(p domain.Identification) domain.Identification {
			return domain.Identification(p)
		}(input.IdentificationDoc),

		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		RegistrationNumber:          input.RegistrationNumber,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
		PracticeServices:            input.PracticeServices,
		Cadre:                       input.Cadre,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil

}

// AddOrganizationPractitionerKyc adds KYC for an organization pratitioner
func (s *SupplierUseCasesImpl) AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.OrganizationPractitioner{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		RegistrationNumber:                 input.RegistrationNumber,
		PracticeLicenseID:                  input.PracticeLicenseID,
		PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
		PracticeServices:                   input.PracticeServices,
		Cadre:                              input.Cadre,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []domain.Identification) []domain.Identification {
			pl := []domain.Identification{}
			for _, i := range p {
				pl = append(pl, domain.Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddOrganizationProviderKyc adds KYC for an organization provider
func (s *SupplierUseCasesImpl) AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.OrganizationProvider{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		RegistrationNumber:                 input.RegistrationNumber,
		PracticeLicenseID:                  input.PracticeLicenseID,
		PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
		PracticeServices:                   input.PracticeServices,
		Cadre:                              input.Cadre,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []domain.Identification) []domain.Identification {
			pl := []domain.Identification{}
			for _, i := range p {
				pl = append(pl, domain.Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddIndividualPharmaceuticalKyc adds KYC for an individual Pharmaceutical kyc
func (s *SupplierUseCasesImpl) AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.IndividualPharmaceutical{
		IdentificationDoc: func(p domain.Identification) domain.Identification {
			return domain.Identification(p)
		}(input.IdentificationDoc),
		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		RegistrationNumber:          input.RegistrationNumber,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddOrganizationPharmaceuticalKyc adds KYC for a pharmacy organization
func (s *SupplierUseCasesImpl) AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.OrganizationPharmaceutical{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []domain.Identification) []domain.Identification {
			pl := []domain.Identification{}
			for _, i := range p {
				pl = append(pl, domain.Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
		RegistrationNumber:      input.RegistrationNumber,
		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddIndividualCoachKyc adds KYC for an individual coach
func (s *SupplierUseCasesImpl) AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.IndividualCoach{
		IdentificationDoc: func(p domain.Identification) domain.Identification {
			return domain.Identification(p)
		}(input.IdentificationDoc),
		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddOrganizationCoachKyc adds KYC for an organization coach
func (s *SupplierUseCasesImpl) AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.OrganizationCoach{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []domain.Identification) []domain.Identification {
			pl := []domain.Identification{}
			for _, i := range p {
				pl = append(pl, domain.Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
		RegistrationNumber:      input.RegistrationNumber,
		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddIndividualNutritionKyc adds KYC for an individual nutritionist
func (s *SupplierUseCasesImpl) AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.IndividualNutrition{
		IdentificationDoc: func(p domain.Identification) domain.Identification {
			return domain.Identification(p)
		}(input.IdentificationDoc),
		KRAPIN:                      input.KRAPIN,
		KRAPINUploadID:              input.KRAPINUploadID,
		SupportingDocumentsUploadID: input.SupportingDocumentsUploadID,
		PracticeLicenseID:           input.PracticeLicenseID,
		PracticeLicenseUploadID:     input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// AddOrganizationNutritionKyc adds kyc for a nutritionist organisation
func (s *SupplierUseCasesImpl) AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get the logged in user supplier profile: %w", err)
	}

	kyc := domain.OrganizationNutrition{
		OrganizationTypeName:               input.OrganizationTypeName,
		KRAPIN:                             input.KRAPIN,
		KRAPINUploadID:                     input.KRAPINUploadID,
		SupportingDocumentsUploadID:        input.SupportingDocumentsUploadID,
		CertificateOfIncorporation:         input.CertificateOfIncorporation,
		CertificateOfInCorporationUploadID: input.CertificateOfInCorporationUploadID,
		DirectorIdentifications: func(p []domain.Identification) []domain.Identification {
			pl := []domain.Identification{}
			for _, i := range p {
				pl = append(pl, domain.Identification(i))
			}
			return pl
		}(input.DirectorIdentifications),
		OrganizationCertificate: input.OrganizationCertificate,
		RegistrationNumber:      input.RegistrationNumber,
		PracticeLicenseUploadID: input.PracticeLicenseUploadID,
	}

	if len(input.SupportingDocumentsUploadID) != 0 {
		ids := []string{}
		ids = append(ids, input.SupportingDocumentsUploadID...)

		kyc.SupportingDocumentsUploadID = ids
	}

	kycAsMap, err := s.parseKYCAsMap(kyc)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal kyc to json")
	}

	sup.SupplierKYC = kycAsMap
	sup.KYCSubmitted = true

	_, err = s.repo.UpdateSupplierProfile(ctx, sup)
	if err != nil {
		return nil, err
	}

	return &kyc, nil
}

// FetchKYCProcessingRequests fetches a list of all unprocessed kyc approval requests
func (s *SupplierUseCasesImpl) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	return s.repo.FetchKYCProcessingRequests(ctx)
}

// // StageKYCProcessingRequest saves kyc processing requests
// func (s *SupplierUseCasesImpl) StageKYCProcessingRequest(sup *domain.Supplier) error {
// 	r := domain.KYCRequest{
// 		ID:                  uuid.New().String(),
// 		ReqPartnerType:      sup.PartnerType,
// 		ReqOrganizationType: domain.OrganizationType(sup.AccountType),
// 		ReqRaw:              sup.SupplierKYC,
// 		Proceseed:           false,
// 		SupplierRecord:      sup,
// 		Status:              domain.KYCProcessStatusPending,
// 	}

// 	_, err := base.SaveDataToFirestore(s.firestoreClient, s.GetKCYProcessCollectionName(), r)
// 	if err != nil {
// 		return fmt.Errorf("unable to save kyc processing request: %w", err)
// 	}
// 	return nil
// }

// // ProcessKYCRequest transitions a kyc request to a given state
// func (s *SupplierUseCasesImpl) ProcessKYCRequest(ctx context.Context, id string, status KYCProcessStatus, rejectionReason *string) (bool, error) {
// 	collection := s.firestoreClient.Collection(s.GetKCYProcessCollectionName())
// 	query := collection.Where("id", "==", id)
// 	docs, err := query.Documents(ctx).GetAll()
// 	if err != nil {
// 		return false, fmt.Errorf("unable to fetch kyc request documents: %v", err)
// 	}

// 	doc := docs[0]
// 	req := &KYCRequest{}
// 	err = doc.DataTo(req)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to read supplier: %w", err)
// 	}

// 	req.Status = status
// 	req.Proceseed = true
// 	req.RejectionReason = rejectionReason

// 	err = base.UpdateRecordOnFirestore(s.firestoreClient, s.GetKCYProcessCollectionName(), doc.Ref.ID, req)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to update KYC request record: %v", err)
// 	}

// 	var email string
// 	var message string

// 	switch status {
// 	case KYCProcessStatusApproved:
// 		// create supplier erp account
// 		if _, err := s.AddSupplier(ctx, *req.SupplierRecord.UserProfile.Name, req.ReqPartnerType); err != nil {
// 			return false, fmt.Errorf("unable to create erp supplier account: %v", err)
// 		}

// 		email = generateProcessKYCApprovalEmailTemplate()
// 		message = "Your KYC details have been reviewed and approved. We look forward to working with you."

// 	case KYCProcessStatusRejected:
// 		email = generateProcessKYCRejectionEmailTemplate()
// 		message = "Your KYC details have been reviewed and not verified. Incase of any queries, please contact us via +254 790 360 360"

// 	}

// 	for _, supplierEmail := range req.SupplierRecord.UserProfile.Emails {
// 		err = s.SendKYCEmail(ctx, email, supplierEmail)
// 		if err != nil {
// 			return false, fmt.Errorf("unable to send KYC processing email: %w", err)
// 		}
// 	}

// 	smsISC := base.SmsISC{
// 		Isc:      s.sms,
// 		EndPoint: sendSMS,
// 	}

// 	twilioISC := base.SmsISC{
// 		Isc:      s.twilio,
// 		EndPoint: sendTwilioSMS,
// 	}

// 	err = base.SendSMS(req.SupplierRecord.UserProfile.Msisdns, message, smsISC, twilioISC)
// 	if err != nil {
// 		return false, fmt.Errorf("unable to send KYC processing message: %w", err)
// 	}

// 	return true, nil

// }

// // SendKYCEmail will send a KYC processing request email to the supplier
// func (s *SupplierUseCasesImpl) SendKYCEmail(ctx context.Context, text, emailaddress string) error {
// 	if !govalidator.IsEmail(emailaddress) {
// 		return nil
// 	}

// 	body := map[string]interface{}{
// 		"to":      []string{emailaddress},
// 		"text":    text,
// 		"subject": emailSignupSubject,
// 	}

// 	resp, err := s.Mailgun.MakeRequest(http.MethodPost, SendEmail, body)
// 	if err != nil {
// 		return fmt.Errorf("unable to send KYC email: %w", err)
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("unable to send KYC email : %w, with status code %v", err, resp.StatusCode)
// 	}

// 	return nil
// }
