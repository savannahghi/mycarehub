package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/segmentio/ksuid"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	emailSignupSubject     = "Thank you for signing up"
	active                 = true
	country                = "KEN" // Anticipate worldwide expansion
	supplierCollectionName = "suppliers"
	futureHours            = 878400
	savannahSladeCode      = "1"
	savannahOrgName        = "Savannah Informatics"

	// PartnerAccountSetupNudgeTitle is the title defined in the `engagement service`
	// for the `PartnerAccountSetupNudge`
	PartnerAccountSetupNudgeTitle = "Setup your partner account"

	// PublishKYCNudgeTitle is the title for the PublishKYCNudge.
	// It takes a partner type as an argument
	PublishKYCNudgeTitle = "Complete your %s KYC"
)

// SupplierUseCases represent the business logic required for management of suppliers
type SupplierUseCases interface {
	AddPartnerType(ctx context.Context, name *string, partnerType *base.PartnerType) (bool, error)

	FindSupplierByID(ctx context.Context, id string) (*base.Supplier, error)

	FindSupplierByUID(ctx context.Context) (*base.Supplier, error)

	SetUpSupplier(ctx context.Context, accountType base.AccountType) (*base.Supplier, error)

	SuspendSupplier(ctx context.Context) (bool, error)

	EDIUserLogin(username, password *string) (*base.EDIUserProfile, error)

	CoreEDIUserLogin(username, password string) (*base.EDIUserProfile, error)

	FetchSupplierAllowedLocations(ctx context.Context) (*resources.BranchConnection, error)

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

	SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*resources.SupplierLogin, error)

	SupplierSetDefaultLocation(ctx context.Context, locationID string) (*base.Supplier, error)

	SaveKYCResponseAndNotifyAdmins(ctx context.Context, sup *base.Supplier) error

	SendKYCEmail(ctx context.Context, text, emailaddress string) error

	StageKYCProcessingRequest(ctx context.Context, sup *base.Supplier) error

	ProcessKYCRequest(ctx context.Context, id string, status domain.KYCProcessStatus, rejectionReason *string) (bool, error)

	RetireKYCRequest(ctx context.Context) error

	PublishKYCFeedItem(ctx context.Context, uids ...string) error

	CreateCustomerAccount(
		ctx context.Context,
		name string,
		partnerType base.PartnerType,
	) error

	CreateSupplierAccount(
		ctx context.Context,
		name string,
		partnerType base.PartnerType,
	) error
}

// SupplierUseCasesImpl represents usecase implementation object
type SupplierUseCasesImpl struct {
	repo    repository.OnboardingRepository
	profile ProfileUseCase

	erp          erp.ServiceERP
	chargemaster chargemaster.ServiceChargeMaster
	engagement   engagement.ServiceEngagement
	messaging    messaging.ServiceMessaging
	baseExt      extension.BaseExtension
	pubsub       pubsubmessaging.ServicePubSub
}

// NewSupplierUseCases returns a new a onboarding usecase
func NewSupplierUseCases(
	r repository.OnboardingRepository,
	p ProfileUseCase,
	er erp.ServiceERP,
	chrg chargemaster.ServiceChargeMaster,
	eng engagement.ServiceEngagement,
	messaging messaging.ServiceMessaging,
	ext extension.BaseExtension,
	pubsub pubsubmessaging.ServicePubSub,
) SupplierUseCases {

	return &SupplierUseCasesImpl{
		repo:         r,
		profile:      p,
		erp:          er,
		chargemaster: chrg,
		engagement:   eng,
		messaging:    messaging,
		baseExt:      ext,
		pubsub:       pubsub,
	}
}

// AddPartnerType create the initial supplier record
func (s SupplierUseCasesImpl) AddPartnerType(ctx context.Context, name *string, partnerType *base.PartnerType) (bool, error) {

	if name == nil || partnerType == nil {
		return false, fmt.Errorf("expected `name` to be defined and `partnerType` to be valid")
	}

	if !partnerType.IsValid() {
		return false, exceptions.InvalidPartnerTypeError()
	}

	if *partnerType == base.PartnerTypeConsumer {
		return false, exceptions.WrongEnumTypeError(partnerType.String())
	}

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		return false, err
	}

	v, err := s.repo.AddPartnerType(ctx, profile.ID, name, partnerType)
	if !v || err != nil {
		return false, exceptions.AddPartnerTypeError(err)
	}

	return true, nil
}

// CreateCustomerAccount makes an external call to the Slade 360 ERP to create
// a customer business partner account
func (s SupplierUseCasesImpl) CreateCustomerAccount(
	ctx context.Context,
	name string,
	partnerType base.PartnerType,
) error {
	UID, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return err
	}

	if partnerType != base.PartnerTypeConsumer {
		return exceptions.WrongEnumTypeError(partnerType.String())
	}

	currency, err := s.baseExt.FetchDefaultCurrency(s.erp.FetchERPClient())
	if err != nil {
		return exceptions.FetchDefaultCurrencyError(err)
	}

	customerPayload := resources.CustomerPayload{
		Active:       active,
		PartnerName:  name,
		Country:      country,
		Currency:     *currency.ID,
		IsCustomer:   true,
		CustomerType: partnerType,
	}

	customerPubSubPayload := resources.CustomerPubSubMessage{
		CustomerPayload: customerPayload,
		UID:             UID,
	}

	bs, err := json.Marshal(customerPubSubPayload)
	if err != nil {
		return err
	}

	topicName := s.pubsub.AddPubSubNamespace(pubsubmessaging.CreateCustomerTopic)
	return s.pubsub.PublishToPubsub(
		ctx,
		topicName,
		bs,
	)
}

// CreateSupplierAccount makes a call to our own ERP and creates a supplier account based
// on the provided partnerType
func (s SupplierUseCasesImpl) CreateSupplierAccount(
	ctx context.Context,
	name string,
	partnerType base.PartnerType,
) error {
	UID, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return err
	}

	if partnerType == base.PartnerTypeConsumer {
		return exceptions.WrongEnumTypeError(partnerType.String())
	}

	currency, err := s.baseExt.FetchDefaultCurrency(s.erp.FetchERPClient())
	if err != nil {
		return exceptions.FetchDefaultCurrencyError(err)
	}

	supplierPayload := resources.SupplierPayload{
		Active:       active,
		PartnerName:  name,
		Country:      country,
		Currency:     *currency.ID,
		IsSupplier:   true,
		SupplierType: partnerType,
	}

	supplierPubSubPayload := resources.SupplierPubSubMessage{
		SupplierPayload: supplierPayload,
		UID:             UID,
	}

	bs, err := json.Marshal(supplierPubSubPayload)
	if err != nil {
		return err
	}

	topicName := s.pubsub.AddPubSubNamespace(pubsubmessaging.CreateSupplierTopic)
	return s.pubsub.PublishToPubsub(
		context.Background(),
		topicName,
		bs,
	)
}

// FindSupplierByID fetches a supplier by their id
func (s SupplierUseCasesImpl) FindSupplierByID(ctx context.Context, id string) (*base.Supplier, error) {
	return s.repo.GetSupplierProfileByID(ctx, id)
}

// FindSupplierByUID fetches a supplier by logged in user uid
func (s SupplierUseCasesImpl) FindSupplierByUID(ctx context.Context) (*base.Supplier, error) {
	pr, err := s.profile.UserProfile(ctx)
	if err != nil {
		return nil, err
	}
	return s.repo.GetSupplierProfileByProfileID(ctx, pr.ID)
}

// SetUpSupplier performs initial account set up during onboarding
func (s SupplierUseCasesImpl) SetUpSupplier(ctx context.Context, accountType base.AccountType) (*base.Supplier, error) {

	validAccountType := accountType.IsValid()
	if !validAccountType {
		return nil, fmt.Errorf("%v is not an allowed AccountType choice", accountType.String())
	}

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	sup, err := s.repo.AddSupplierAccountType(ctx, profile.ID, accountType)
	if err != nil {
		return nil, err
	}

	go func(u string, pnt base.PartnerType, acnt base.AccountType) {
		op := func() error {
			return s.PublishKYCNudge(ctx, u, &pnt, &acnt)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}(uid, sup.PartnerType, *sup.AccountType)

	go func() {
		pro := func() error {
			return s.engagement.ResolveDefaultNudgeByTitle(
				uid,
				base.FlavourPro,
				PartnerAccountSetupNudgeTitle,
			)
		}
		if err := backoff.Retry(
			pro,
			backoff.NewExponentialBackOff(),
		); err != nil {
			logrus.Error(err)
		}
	}()

	return sup, nil
}

// SuspendSupplier flips the active boolean on the erp partner from true to false
func (s SupplierUseCasesImpl) SuspendSupplier(ctx context.Context) (bool, error) {
	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return false, exceptions.UserNotFoundError(err)
	}
	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return false, err
	}
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return false, err
	}
	sup.Active = false

	if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup); err != nil {
		return false, err
	}

	//TODO(dexter) notify the supplier of the suspension

	return true, nil

}

// EDIUserLogin used to login a user to EDI (Portal Authserver) and return their
// EDI (Portal Authserver) profile
func (s SupplierUseCasesImpl) EDIUserLogin(username, password *string) (*base.EDIUserProfile, error) {
	if username == nil || password == nil {
		return nil, exceptions.InvalidCredentialsError()
	}

	ediClient, err := s.baseExt.LoginClient(*username, *password)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := s.baseExt.FetchUserProfile(ediClient)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
	}

	return userProfile, nil

}

// CoreEDIUserLogin used to login a user to EDI (Core Authserver) and return their EDI
// EDI (Core Authserver) profile
func (s SupplierUseCasesImpl) CoreEDIUserLogin(username, password string) (*base.EDIUserProfile, error) {

	if username == "" || password == "" {
		return nil, exceptions.InvalidCredentialsError()
	}

	ediClient, err := s.baseExt.LoginClient(username, password)
	if err != nil {
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := s.baseExt.FetchUserProfile(ediClient)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
	}

	return userProfile, nil

}

// SupplierEDILogin it used to instantiate as call when setting up a supplier's account's who
// has an affiliation to a provider with the slade ecosystem. The logic is as follows;
// 1 . login to the relevant edi to assert the user has an account
// 2 . fetch the branches of the provider given the slade code which we have
// 3 . update the user's supplier record
// 4. return the list of branches to the frontend so that a default location can be set
func (s SupplierUseCasesImpl) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*resources.SupplierLogin, error) {
	var resp resources.SupplierLogin

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	supplier, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	accType := base.AccountTypeIndividual
	supplier.AccountType = &accType
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
		supplier.EDIUserProfile = ediUserProfile
		supplier.IsOrganizationVerified = true
		supplier.SladeCode = sladeCode
		supplier.Active = true
		supplier.KYCSubmitted = true
		supplier.PartnerSetupComplete = true
		supplier.OrganizationName = savannahOrgName

		if err := s.repo.UpdateSupplierProfile(ctx, *supplier.ProfileID, supplier); err != nil {
			return nil, err
		}

		if err := s.profile.UpdatePermissions(ctx, base.DefaultAdminPermissions); err != nil {
			return nil, err
		}

		resp.Supplier = supplier
		return &resp, nil
	}
	// verify slade code.
	if ediUserProfile.BusinessPartner != orgSladeCode {
		return nil, exceptions.InvalidSladeCodeError()
	}
	supplier.EDIUserProfile = ediUserProfile
	supplier.IsOrganizationVerified = true
	supplier.SladeCode = sladeCode

	filter := []*resources.BusinessPartnerFilterInput{
		{
			SladeCode: &sladeCode,
		},
	}

	partner, err := s.chargemaster.FindProvider(ctx, nil, filter, nil)
	if err != nil {
		return nil, exceptions.FindProviderError(err)
	}
	if len(partner.Edges) != 1 {
		return nil, fmt.Errorf("expected one business partner, found: %v", len(partner.Edges))
	}

	businessPartner := *partner.Edges[0].Node

	go func() {
		op := func() error {
			return s.PublishKYCNudge(ctx, uid, &supplier.PartnerType, supplier.AccountType)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}()

	go func() {
		pro := func() error {
			return s.engagement.ResolveDefaultNudgeByTitle(
				uid,
				base.FlavourPro,
				PartnerAccountSetupNudgeTitle,
			)
		}
		if err := backoff.Retry(
			pro,
			backoff.NewExponentialBackOff(),
		); err != nil {
			logrus.Error(err)
		}
	}()

	if businessPartner.Parent != nil {
		supplier.HasBranches = true
		supplier.ParentOrganizationID = *businessPartner.Parent

		partner, err := s.chargemaster.FetchProviderByID(ctx, *businessPartner.Parent)
		if err != nil {
			return nil, exceptions.FindProviderError(err)
		}

		supplier.OrganizationName = partner.Name

		if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, supplier); err != nil {
			return nil, err
		}

		// fetch all locations of the business partner
		filter := []*resources.BranchFilterInput{
			{
				ParentOrganizationID: &supplier.ParentOrganizationID,
			},
		}

		brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
		if len(brs.Edges) == 0 || err != nil {
			return nil, exceptions.FindProviderError(err)
		}

		if len(brs.Edges) > 1 {
			// set branches in the final response object
			resp.Branches = brs
			resp.Supplier = supplier
			return &resp, nil
		}

		if len(brs.Edges) == 1 {
			spr, err := s.SupplierSetDefaultLocation(ctx, brs.Edges[0].Node.ID)
			if err != nil {
				return nil, err
			}

			resp.Supplier = spr
			return &resp, nil

		}
		return nil, exceptions.InternalServerError(nil)

	}

	// set the main branch as the supplier's location
	supplier.OrganizationName = businessPartner.Name
	loc := base.Location{
		ID:   businessPartner.ID,
		Name: businessPartner.Name,
	}
	supplier.Location = &loc

	if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, supplier); err != nil {
		return nil, err
	}

	// set the supplier profile in the final response object
	resp.Supplier = supplier

	return &resp, nil
}

// SupplierSetDefaultLocation updates the default location ot the supplier by the given location id
func (s SupplierUseCasesImpl) SupplierSetDefaultLocation(ctx context.Context, locationID string) (*base.Supplier, error) {

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	// fetch the branches of the provider filtered by ParentOrganizationID
	filter := []*resources.BranchFilterInput{
		{
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		return nil, exceptions.FindProviderError(err)
	}

	branch := func(brs *resources.BranchConnection, location string) *resources.BranchEdge {
		for _, b := range brs.Edges {
			if b.Node.ID == location {
				return b
			}
		}
		return nil
	}(brs, locationID)

	if branch != nil {
		loc := base.Location{
			ID:              branch.Node.ID,
			Name:            branch.Node.Name,
			BranchSladeCode: &branch.Node.BranchSladeCode,
		}
		sup.Location = &loc

		if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup); err != nil {
			return nil, err
		}

		// refetch the supplier profile and return it
		return s.FindSupplierByUID(ctx)
	}

	return nil, fmt.Errorf("unable to get location of id %v : %v", locationID, err)
}

// FetchSupplierAllowedLocations retrieves all the locations that the user in context can work on.
func (s *SupplierUseCasesImpl) FetchSupplierAllowedLocations(ctx context.Context) (*resources.BranchConnection, error) {

	// fetch the supplier's profile
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	// fetch the branches of the provider filtered by ParentOrganizationID
	filter := []*resources.BranchFilterInput{
		{
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		return nil, exceptions.FindProviderError(err)
	}

	if sup.Location != nil {
		var resp resources.BranchConnection
		resp.PageInfo = brs.PageInfo

		newEdges := []*resources.BranchEdge{}
		for _, edge := range brs.Edges {
			if edge.Node.ID == sup.Location.ID {
				loc := &resources.BranchEdge{
					Cursor: edge.Cursor,
					Node: &domain.Branch{
						ID:                    edge.Node.ID,
						Name:                  edge.Node.Name,
						OrganizationSladeCode: edge.Node.OrganizationSladeCode,
						BranchSladeCode:       edge.Node.BranchSladeCode,
						Default:               true,
					},
				}
				newEdges = append(newEdges, loc)
			} else {
				newEdges = append(newEdges, edge)
			}
		}

		resp.Edges = newEdges
		return &resp, nil
	}
	return brs, nil
}

// PublishKYCNudge pushes a KYC nudge to the user feed
func (s *SupplierUseCasesImpl) PublishKYCNudge(
	ctx context.Context,
	uid string,
	partner *base.PartnerType,
	account *base.AccountType,
) error {
	if partner == nil || account == nil {
		return exceptions.PublishKYCNudgeError(
			fmt.Errorf("expected `partner` to be defined and to be valid"),
		)
	}

	if *partner == base.PartnerTypeConsumer {
		return exceptions.WrongEnumTypeError(partner.String())
	}

	if !account.IsValid() || !partner.IsValid() {
		return exceptions.WrongEnumTypeError(
			fmt.Sprintf(
				"%s, %s",
				account.String(),
				partner.String(),
			),
		)
	}

	title := fmt.Sprintf(
		PublishKYCNudgeTitle,
		strings.ToLower(partner.String()),
	)
	text := "Fill in your Be.Well business KYC in order to start transacting."

	nudge := base.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: int(time.Now().Unix()),
		Visibility:     base.VisibilityShow,
		Status:         base.StatusPending,
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Title:          title,
		Text:           text,
		Links: []base.Link{
			{
				ID:          ksuid.New().String(),
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC",
				Description: fmt.Sprintf("KYC for %v", partner.String()),
				Thumbnail:   base.LogoURL,
			},
		},
		Actions: []base.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: int(time.Now().Unix()),
				Name: strings.ToUpper(fmt.Sprintf(
					"COMPLETE_%v_%v_KYC",
					account.String(),
					partner.String(),
				),
				),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
				Icon: base.Link{
					ID:          ksuid.New().String(),
					URL:         base.LogoURL,
					LinkType:    base.LinkTypePngImage,
					Title:       title,
					Description: text,
					Thumbnail:   base.LogoURL,
				},
			},
		},
		Users:  []string{uid},
		Groups: []string{uid},
		NotificationChannels: []base.Channel{
			base.ChannelEmail,
			base.ChannelFcm,
		},
	}

	// TODO: This call should be asynchronous (Pub/Sub)
	resp, err := s.engagement.PublishKYCNudge(uid, nudge)
	if err != nil {
		return exceptions.PublishKYCNudgeError(err)
	}

	if resp.StatusCode != http.StatusOK {
		if err := s.SaveProfileNudge(ctx, &nudge); err != nil {
			logrus.Errorf("failed to stage nudge : %v", err)
		}

		return exceptions.PublishKYCNudgeError(
			fmt.Errorf("unable to publish kyc nudge. unexpected status code  %v",
				resp.StatusCode,
			),
		)
	}
	return nil
}

// PublishKYCFeedItem notifies admin users of a KYC approval request
func (s SupplierUseCasesImpl) PublishKYCFeedItem(ctx context.Context, uids ...string) error {

	for _, uid := range uids {
		payload := base.Item{
			ID:             ksuid.New().String(),
			SequenceNumber: int(time.Now().Unix()),
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
				ID:          ksuid.New().String(),
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC Review",
				Description: "Review KYC for the partner and either approve or reject",
				Thumbnail:   base.LogoURL,
			},
			Timestamp: time.Now(),
			Actions: []base.Action{
				{
					ID:             ksuid.New().String(),
					SequenceNumber: int(time.Now().Unix()),
					Name:           "Review KYC details",
					Icon: base.Link{
						ID:          ksuid.New().String(),
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
					ID:          ksuid.New().String(),
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
func (s *SupplierUseCasesImpl) SaveProfileNudge(
	ctx context.Context,
	nudge *base.Nudge,
) error {
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

// SaveKYCResponseAndNotifyAdmins saves the kyc information provided by the user
// and sends a notification to all admins for a pending KYC review request
func (s *SupplierUseCasesImpl) SaveKYCResponseAndNotifyAdmins(ctx context.Context, sup *base.Supplier) error {
	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}
	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup); err != nil {
		return err
	}
	if err := s.StageKYCProcessingRequest(ctx, sup); err != nil {
		return err
	}

	go func() {
		op := func() error {
			a, err := s.repo.FetchAdminUsers(ctx)
			if err != nil {
				return err
			}
			var uids []string
			for _, u := range a {
				uids = append(uids, u.ID)
			}

			return s.PublishKYCFeedItem(ctx, uids...)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			logrus.Error(err)
		}
	}()

	return nil
}

// StageKYCProcessingRequest saves kyc processing requests
func (s *SupplierUseCasesImpl) StageKYCProcessingRequest(ctx context.Context, sup *base.Supplier) error {
	r := &domain.KYCRequest{
		ID:             uuid.New().String(),
		ReqPartnerType: sup.PartnerType,
		ReqRaw:         sup.SupplierKYC,
		Processed:      false,
		SupplierRecord: sup,
		Status:         domain.KYCProcessStatusPending,
		FiledTimestamp: time.Now().In(domain.TimeLocation),
	}

	return s.repo.StageKYCProcessingRequest(ctx, r)
}

// AddIndividualRiderKyc adds KYC for an individual rider
func (s *SupplierUseCasesImpl) AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}
	if !sup.KYCSubmitted {
		if !input.IdentificationDoc.IdentificationDocType.IsValid() {
			return nil, exceptions.WrongEnumTypeError(input.IdentificationDoc.IdentificationDocType.String())
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
			SupportingDocuments:            input.SupportingDocuments,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationRiderKyc adds KYC for an organization rider
func (s *SupplierUseCasesImpl) AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error) {

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName)
	}

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
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

			KRAPIN:              input.KRAPIN,
			KRAPINUploadID:      input.KRAPINUploadID,
			SupportingDocuments: input.SupportingDocuments,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualPractitionerKyc adds KYC for an individual practitioner
func (s *SupplierUseCasesImpl) AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		for _, p := range input.PracticeServices {
			if !p.IsValid() {
				return nil, fmt.Errorf("invalid `PracticeService` provided : %v", p.String())
			}
		}

		kyc := domain.IndividualPractitioner{

			IdentificationDoc: func(p domain.Identification) domain.Identification {
				return domain.Identification(p)
			}(input.IdentificationDoc),

			KRAPIN:                  input.KRAPIN,
			KRAPINUploadID:          input.KRAPINUploadID,
			RegistrationNumber:      input.RegistrationNumber,
			PracticeLicenseID:       input.PracticeLicenseID,
			PracticeLicenseUploadID: input.PracticeLicenseUploadID,
			PracticeServices:        input.PracticeServices,
			Cadre:                   input.Cadre,
			SupportingDocuments:     input.SupportingDocuments,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()

}

// AddOrganizationPractitionerKyc adds KYC for an organization practitioner
func (s *SupplierUseCasesImpl) AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.OrganizationTypeName.IsValid() {
			return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName.String())
		}

		kyc := domain.OrganizationPractitioner{
			OrganizationTypeName:               input.OrganizationTypeName,
			KRAPIN:                             input.KRAPIN,
			KRAPINUploadID:                     input.KRAPINUploadID,
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
			SupportingDocuments:     input.SupportingDocuments,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationProviderKyc adds KYC for an organization provider
func (s *SupplierUseCasesImpl) AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.OrganizationTypeName.IsValid() {
			return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName.String())
		}

		for _, practiceService := range input.PracticeServices {
			if !practiceService.IsValid() {
				return nil, fmt.Errorf("invalid `PracticeService` provided : %v", practiceService.String())
			}
		}

		kyc := domain.OrganizationProvider{
			OrganizationTypeName:               input.OrganizationTypeName,
			KRAPIN:                             input.KRAPIN,
			KRAPINUploadID:                     input.KRAPINUploadID,
			RegistrationNumber:                 input.RegistrationNumber,
			PracticeLicenseID:                  input.PracticeLicenseID,
			PracticeLicenseUploadID:            input.PracticeLicenseUploadID,
			PracticeServices:                   input.PracticeServices,
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
			SupportingDocuments:     input.SupportingDocuments,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualPharmaceuticalKyc adds KYC for an individual Pharmaceutical kyc
func (s *SupplierUseCasesImpl) AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error) {

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.IdentificationDoc.IdentificationDocType.IsValid() {
			return nil, fmt.Errorf("invalid `IdentificationDocType` provided : %v", input.IdentificationDoc.IdentificationDocType.String())
		}

		kyc := domain.IndividualPharmaceutical{
			IdentificationDoc: func(p domain.Identification) domain.Identification {
				return domain.Identification(p)
			}(input.IdentificationDoc),
			KRAPIN:                  input.KRAPIN,
			KRAPINUploadID:          input.KRAPINUploadID,
			RegistrationNumber:      input.RegistrationNumber,
			PracticeLicenseID:       input.PracticeLicenseID,
			PracticeLicenseUploadID: input.PracticeLicenseUploadID,
			SupportingDocuments:     input.SupportingDocuments,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationPharmaceuticalKyc adds KYC for a pharmacy organization
func (s *SupplierUseCasesImpl) AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.OrganizationTypeName.IsValid() {
			return nil, fmt.Errorf("invalid `OrganizationTypeName` provided : %v", input.OrganizationTypeName.String())
		}

		kyc := domain.OrganizationPharmaceutical{
			OrganizationTypeName:               input.OrganizationTypeName,
			KRAPIN:                             input.KRAPIN,
			KRAPINUploadID:                     input.KRAPINUploadID,
			SupportingDocuments:                input.SupportingDocuments,
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
			PracticeLicenseID:       input.PracticeLicenseID,
			PracticeLicenseUploadID: input.PracticeLicenseUploadID,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualCoachKyc adds KYC for an individual coach
func (s *SupplierUseCasesImpl) AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.IdentificationDoc.IdentificationDocType.IsValid() {
			return nil, fmt.Errorf("invalid `IdentificationDocType` provided : %v", input.IdentificationDoc.IdentificationDocType.String())
		}

		kyc := domain.IndividualCoach{
			IdentificationDoc: func(p domain.Identification) domain.Identification {
				return domain.Identification(p)
			}(input.IdentificationDoc),
			KRAPIN:                  input.KRAPIN,
			KRAPINUploadID:          input.KRAPINUploadID,
			SupportingDocuments:     input.SupportingDocuments,
			PracticeLicenseID:       input.PracticeLicenseID,
			PracticeLicenseUploadID: input.PracticeLicenseUploadID,
			AccreditationUploadID:   input.AccreditationUploadID,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationCoachKyc adds KYC for an organization coach
func (s *SupplierUseCasesImpl) AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		kyc := domain.OrganizationCoach{
			OrganizationTypeName:               input.OrganizationTypeName,
			KRAPIN:                             input.KRAPIN,
			KRAPINUploadID:                     input.KRAPINUploadID,
			SupportingDocuments:                input.SupportingDocuments,
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

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualNutritionKyc adds KYC for an individual nutritionist
func (s *SupplierUseCasesImpl) AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}
	if !sup.KYCSubmitted {
		kyc := domain.IndividualNutrition{
			IdentificationDoc: func(p domain.Identification) domain.Identification {
				return domain.Identification(p)
			}(input.IdentificationDoc),
			KRAPIN:                  input.KRAPIN,
			KRAPINUploadID:          input.KRAPINUploadID,
			SupportingDocuments:     input.SupportingDocuments,
			PracticeLicenseID:       input.PracticeLicenseID,
			PracticeLicenseUploadID: input.PracticeLicenseUploadID,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationNutritionKyc adds kyc for a nutritionist organisation
func (s *SupplierUseCasesImpl) AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error) {
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		kyc := domain.OrganizationNutrition{
			OrganizationTypeName:               input.OrganizationTypeName,
			KRAPIN:                             input.KRAPIN,
			KRAPINUploadID:                     input.KRAPINUploadID,
			SupportingDocuments:                input.SupportingDocuments,
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
			PracticeLicenseID:       input.PracticeLicenseID,
			PracticeLicenseUploadID: input.PracticeLicenseUploadID,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// FetchKYCProcessingRequests fetches a list of all unprocessed kyc approval requests
func (s *SupplierUseCasesImpl) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	return s.repo.FetchKYCProcessingRequests(ctx)
}

// SendKYCEmail will send a KYC processing request email to the supplier
func (s *SupplierUseCasesImpl) SendKYCEmail(ctx context.Context, text, emailaddress string) error {
	return s.engagement.SendMail(emailaddress, text, emailSignupSubject)
}

// ProcessKYCRequest transitions a kyc request to a given state
func (s *SupplierUseCasesImpl) ProcessKYCRequest(
	ctx context.Context,
	id string,
	status domain.KYCProcessStatus,
	rejectionReason *string,
) (bool, error) {
	reviewerProfile, err := s.profile.UserProfile(ctx)
	if err != nil {
		return false, err
	}

	s.repo.CheckIfAdmin(reviewerProfile)

	KYCRequest, err := s.repo.FetchKYCProcessingRequestByID(ctx, id)
	if err != nil {
		return false, err
	}

	KYCRequest.Status = status
	KYCRequest.Processed = true
	if rejectionReason != nil {
		KYCRequest.RejectionReason = rejectionReason
	}
	KYCRequest.ProcessedTimestamp = time.Now().In(domain.TimeLocation)
	KYCRequest.ProcessedBy = reviewerProfile.ID

	if err := s.repo.UpdateKYCProcessingRequest(
		ctx,
		KYCRequest,
	); err != nil {
		return false, fmt.Errorf("unable to update KYC request record: %v", err)
	}

	supplierProfile, err := s.profile.GetProfileByID(
		ctx,
		KYCRequest.SupplierRecord.ProfileID,
	)
	if err != nil {
		return false, err
	}

	var email string
	var message string

	switch status {
	case domain.KYCProcessStatusApproved:
		go func() {
			if err := s.CreateSupplierAccount(
				ctx,
				KYCRequest.SupplierRecord.SupplierName,
				KYCRequest.ReqPartnerType,
			); err != nil {
				logrus.Error(fmt.Errorf("unable to create erp supplier account: %v", err))
			}
		}()

		email = s.generateProcessKYCApprovalEmailTemplate()
		message = "Your KYC details have been reviewed and approved. We look forward to working with you."

		supplier, err := s.FindSupplierByUID(ctx)
		if err != nil {
			return false, err
		}

		supplier.Active = true

		if err := s.repo.UpdateSupplierProfile(
			ctx,
			supplierProfile.ID,
			supplier,
		); err != nil {
			return false, err
		}

	case domain.KYCProcessStatusRejected:
		email = s.generateProcessKYCRejectionEmailTemplate()
		message = "Your KYC details have been reviewed and not verified. Incase of any queries, please contact us via +254 790 360 360"

	}

	nudgeTitle := fmt.Sprintf(
		PublishKYCNudgeTitle,
		strings.ToLower(string(KYCRequest.ReqPartnerType)),
	)
	supplierVerifiedUIDs := supplierProfile.VerifiedUIDS
	go func() {
		for _, UID := range supplierVerifiedUIDs {
			if err = s.engagement.ResolveDefaultNudgeByTitle(
				UID,
				base.FlavourPro,
				nudgeTitle,
			); err != nil {
				logrus.Print(err)
			}
		}
	}()

	supplierEmails := func(profile *base.UserProfile) []string {
		var emails []string
		if profile.PrimaryEmailAddress != nil {
			emails = append(emails, *profile.PrimaryEmailAddress)
		}
		emails = append(emails, profile.SecondaryEmailAddresses...)
		return emails
	}(supplierProfile)

	for _, supplierEmail := range supplierEmails {
		err = s.SendKYCEmail(ctx, email, supplierEmail)
		if err != nil {
			return false, fmt.Errorf("unable to send KYC processing email: %w", err)
		}
	}

	supplierPhones := func(profile *base.UserProfile) []string {
		var phones []string
		phones = append(phones, *profile.PrimaryPhone)
		phones = append(phones, profile.SecondaryPhoneNumbers...)
		return phones
	}(supplierProfile)

	if err := s.messaging.SendSMS(supplierPhones, message); err != nil {
		return false, fmt.Errorf("unable to send KYC processing message: %w", err)
	}

	return true, nil
}

// RetireKYCRequest retires the KYC process request of a supplier
func (s *SupplierUseCasesImpl) RetireKYCRequest(ctx context.Context) error {
	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		return exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	sup, err := s.repo.GetSupplierProfileByProfileID(ctx, profile.ID)
	if err != nil {
		// this is a wrapped error. No need to wrap it again
		return err
	}
	if err := s.repo.RemoveKYCProcessingRequest(ctx, sup.ID); err != nil {
		// the error is a custom error already. No need to wrap it here
		return err
	}

	return nil

}

func (s *SupplierUseCasesImpl) generateProcessKYCApprovalEmailTemplate() string {
	t := template.Must(template.New("approvalKYCEmail").Parse(utils.ProcessKYCApprovalEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating KYC approval email template: %s", err)
	}
	return buf.String()
}

func (s *SupplierUseCasesImpl) generateProcessKYCRejectionEmailTemplate() string {
	t := template.Must(template.New("rejectionKYCEmail").Parse(utils.ProcessKYCRejectionEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, "")
	if err != nil {
		log.Fatalf("Error while generating KYC rejection email template: %s", err)
	}
	return buf.String()
}
