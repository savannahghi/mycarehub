package usecases

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/profileutils"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization/permission"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/chargemaster"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"

	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
	erp "gitlab.slade360emr.com/go/commontools/accounting/pkg/usecases"
)

// Supplier constants
const (
	emailKYCSubject        = "KYC Request"
	active                 = true
	country                = "KEN" // Anticipate worldwide expansion
	supplierCollectionName = "suppliers"
	futureHours            = 878400
	SavannahSladeCode      = "1"
	SavannahOrgName        = "Savannah Informatics"
	adminEmailBody         = `
	The below supplier KYC request has been made.
	To view and process the request, please log in to Be.Well Professional.
	`
	// KYC Acknowledgement Email
	supplierEmailSubjectTitle = "KYC Acknowledgement Email"
	supplierEmailBody         = `
		we acknowledge receipt of your.
		`
	// Supplier Suspension EmailSubject Title
	supplierSuspensionEmailSubjectTitle = "Suspension from Be.Well"
	// PartnerAccountSetupNudgeTitle is the title defined in the `engagement service`
	// for the `PartnerAccountSetupNudge`
	PartnerAccountSetupNudgeTitle = "Setup your partner account"

	// PublishKYCNudgeTitle is the title for the PublishKYCNudge.
	// It takes a partner type as an argument
	PublishKYCNudgeTitle = "Complete your %s KYC"
)

// SupplierUseCases represent the business logic required for management of suppliers
type SupplierUseCases interface {
	AddPartnerType(ctx context.Context, name *string, partnerType *profileutils.PartnerType) (bool, error)

	FindSupplierByID(ctx context.Context, id string) (*profileutils.Supplier, error)

	FindSupplierByUID(ctx context.Context) (*profileutils.Supplier, error)

	SetUpSupplier(ctx context.Context, accountType profileutils.AccountType) (*profileutils.Supplier, error)

	SuspendSupplier(ctx context.Context, suspensionReason *string) (bool, error)

	EDIUserLogin(ctx context.Context, username, password *string) (*profileutils.EDIUserProfile, error)

	CoreEDIUserLogin(ctx context.Context, username, password string) (*profileutils.EDIUserProfile, error)

	FetchSupplierAllowedLocations(ctx context.Context) (*dto.BranchConnection, error)
	CheckSupplierKYCSubmitted(ctx context.Context) (bool, error)

	AddIndividualRiderKyc(
		ctx context.Context,
		input domain.IndividualRider,
	) (*domain.IndividualRider, error)

	AddOrganizationRiderKyc(
		ctx context.Context,
		input domain.OrganizationRider,
	) (*domain.OrganizationRider, error)

	AddIndividualPractitionerKyc(
		ctx context.Context,
		input domain.IndividualPractitioner,
	) (*domain.IndividualPractitioner, error)

	AddOrganizationPractitionerKyc(
		ctx context.Context,
		input domain.OrganizationPractitioner,
	) (*domain.OrganizationPractitioner, error)

	AddOrganizationProviderKyc(
		ctx context.Context,
		input domain.OrganizationProvider,
	) (*domain.OrganizationProvider, error)

	AddIndividualPharmaceuticalKyc(
		ctx context.Context,
		input domain.IndividualPharmaceutical,
	) (*domain.IndividualPharmaceutical, error)

	AddOrganizationPharmaceuticalKyc(
		ctx context.Context,
		input domain.OrganizationPharmaceutical,
	) (*domain.OrganizationPharmaceutical, error)

	AddIndividualCoachKyc(
		ctx context.Context,
		input domain.IndividualCoach,
	) (*domain.IndividualCoach, error)

	AddOrganizationCoachKyc(
		ctx context.Context,
		input domain.OrganizationCoach,
	) (*domain.OrganizationCoach, error)

	AddIndividualNutritionKyc(
		ctx context.Context,
		input domain.IndividualNutrition,
	) (*domain.IndividualNutrition, error)

	AddOrganizationNutritionKyc(
		ctx context.Context,
		input domain.OrganizationNutrition,
	) (*domain.OrganizationNutrition, error)

	FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error)

	SupplierEDILogin(
		ctx context.Context,
		username string,
		password string,
		sladeCode string,
	) (*dto.SupplierLogin, error)

	SupplierSetDefaultLocation(ctx context.Context, locationID string) (*profileutils.Supplier, error)

	SaveKYCResponseAndNotifyAdmins(ctx context.Context, sup *profileutils.Supplier) error

	SendKYCEmail(ctx context.Context, text, emailaddress string) error

	StageKYCProcessingRequest(ctx context.Context, sup *profileutils.Supplier) error

	ProcessKYCRequest(
		ctx context.Context,
		id string,
		status domain.KYCProcessStatus,
		rejectionReason *string,
	) (bool, error)

	RetireKYCRequest(ctx context.Context) error

	PublishKYCFeedItem(ctx context.Context, uids ...string) error

	CreateCustomerAccount(
		ctx context.Context,
		name string,
		partnerType profileutils.PartnerType,
	) error

	CreateSupplierAccount(
		ctx context.Context,
		name string,
		partnerType profileutils.PartnerType,
	) error
}

// SupplierUseCasesImpl represents usecase implementation object
type SupplierUseCasesImpl struct {
	repo         repository.OnboardingRepository
	profile      ProfileUseCase
	erp          erp.AccountingUsecase
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
	er erp.AccountingUsecase,
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
func (s SupplierUseCasesImpl) AddPartnerType(
	ctx context.Context,
	name *string,
	partnerType *profileutils.PartnerType,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "AddPartnerType")
	defer span.End()

	if name == nil || partnerType == nil {
		return false, fmt.Errorf("expected `name` to be defined and `partnerType` to be valid")
	}

	if !partnerType.IsValid() {
		return false, exceptions.InvalidPartnerTypeError()
	}

	if *partnerType == profileutils.PartnerTypeConsumer {
		return false, exceptions.WrongEnumTypeError(partnerType.String())
	}

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.PartnerTypeCreate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	if !isAuthorized {
		return false, fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	v, err := s.repo.AddPartnerType(ctx, profile.ID, name, partnerType)
	if !v || err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.AddPartnerTypeError(err)
	}

	return true, nil
}

// CreateCustomerAccount makes an external call to the Slade 360 ERP to create
// a customer business partner account
func (s SupplierUseCasesImpl) CreateCustomerAccount(
	ctx context.Context,
	name string,
	partnerType profileutils.PartnerType,
) error {
	ctx, span := tracer.Start(ctx, "CreateCustomerAccount")
	defer span.End()

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.CustomerAccountCreate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	if partnerType != profileutils.PartnerTypeConsumer {
		return exceptions.WrongEnumTypeError(partnerType.String())
	}

	currency, err := s.baseExt.FetchDefaultCurrency(s.erp.FetchERPClient())
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.FetchDefaultCurrencyError(err)
	}

	customerPayload := dto.CustomerPayload{
		Active:       active,
		PartnerName:  name,
		Country:      country,
		Currency:     *currency.ID,
		IsCustomer:   true,
		CustomerType: partnerType,
	}

	customerPubSubPayload := dto.CustomerPubSubMessage{
		CustomerPayload: customerPayload,
		UID:             user.UID,
	}

	if err = s.pubsub.NotifyCreateCustomer(ctx, customerPubSubPayload); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to customers.create topic: %v", err)
	}

	return nil
}

// CreateSupplierAccount makes a call to our own ERP and creates a supplier account based
// on the provided partnerType
func (s SupplierUseCasesImpl) CreateSupplierAccount(
	ctx context.Context,
	name string,
	partnerType profileutils.PartnerType,
) error {
	ctx, span := tracer.Start(ctx, "CreateSupplierAccount")
	defer span.End()

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.SupplierAccountCreate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}
	if !isAuthorized {
		return fmt.Errorf("user not authorized to access this resource")
	}

	if partnerType == profileutils.PartnerTypeConsumer {
		return exceptions.WrongEnumTypeError(partnerType.String())
	}

	currency, err := s.baseExt.FetchDefaultCurrency(s.erp.FetchERPClient())
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.FetchDefaultCurrencyError(err)
	}

	supplierPayload := dto.SupplierPayload{
		Active:       active,
		PartnerName:  name,
		Country:      country,
		Currency:     *currency.ID,
		IsSupplier:   true,
		SupplierType: partnerType,
	}

	supplierPubSubPayload := dto.SupplierPubSubMessage{
		SupplierPayload: supplierPayload,
		UID:             user.UID,
	}

	if err = s.pubsub.NotifyCreateSupplier(ctx, supplierPubSubPayload); err != nil {
		utils.RecordSpanError(span, err)
		log.Printf("failed to publish to suppliers.create topic: %v", err)
	}

	return nil
}

// FindSupplierByID fetches a supplier by their id
func (s SupplierUseCasesImpl) FindSupplierByID(
	ctx context.Context,
	id string,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "FindSupplierByID")
	defer span.End()

	return s.repo.GetSupplierProfileByID(ctx, id)
}

// FindSupplierByUID fetches a supplier by logged in user uid
func (s SupplierUseCasesImpl) FindSupplierByUID(ctx context.Context) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "FindSupplierByUID")
	defer span.End()

	pr, err := s.profile.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	return s.repo.GetSupplierProfileByProfileID(ctx, pr.ID)
}

// CheckSupplierKYCSubmitted checks if a supplier has submitted KYC already.
func (s SupplierUseCasesImpl) CheckSupplierKYCSubmitted(ctx context.Context) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckSupplierKYCSubmitted")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return false, err
	}
	if !sup.KYCSubmitted {
		return false, fmt.Errorf("the supplier has no KYC submitted")
	}
	exists := false
	if sup.KYCSubmitted {
		exists = true
	}
	return exists, nil

}

// SetUpSupplier performs initial account set up during onboarding
func (s SupplierUseCasesImpl) SetUpSupplier(
	ctx context.Context,
	accountType profileutils.AccountType,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "SetUpSupplier")
	defer span.End()

	validAccountType := accountType.IsValid()
	if !validAccountType {
		return nil, fmt.Errorf("%v is not an allowed AccountType choice", accountType.String())
	}

	user, err := s.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("can't get user: %w", err)
	}
	isAuthorized, err := authorization.IsAuthorized(user, permission.SupplierCreate)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !isAuthorized {
		return nil, fmt.Errorf("user not authorized to access this resource")
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	sup, err := s.repo.AddSupplierAccountType(ctx, profile.ID, accountType)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if *sup.AccountType == profileutils.AccountTypeOrganisation ||
		*sup.AccountType == profileutils.AccountTypeIndividual {
		sup.OrganizationName = sup.SupplierName
		err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}
	}

	go func(u string, pnt profileutils.PartnerType, acnt profileutils.AccountType) {
		op := func() error {
			return s.PublishKYCNudge(ctx, u, &pnt, &acnt)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}(user.UID, sup.PartnerType, *sup.AccountType)

	go func() {
		pro := func() error {
			return s.engagement.ResolveDefaultNudgeByTitle(
				ctx,
				user.UID,
				feedlib.FlavourPro,
				PartnerAccountSetupNudgeTitle,
			)
		}
		if err := backoff.Retry(
			pro,
			backoff.NewExponentialBackOff(),
		); err != nil {
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}()

	return sup, nil
}

// SuspendSupplier flips the active boolean on the erp partner from true to false
func (s SupplierUseCasesImpl) SuspendSupplier(ctx context.Context, suspensionReason *string) (bool, error) {
	ctx, span := tracer.Start(ctx, "SuspendSupplier")
	defer span.End()

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, exceptions.UserNotFoundError(err)
	}
	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return false, err
	}
	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return false, err
	}
	sup.Active = false

	if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup); err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	supplierEmailPayload := dto.EmailNotificationPayload{
		SupplierName: *profile.UserBioData.FirstName,
		SubjectTitle: supplierSuspensionEmailSubjectTitle,
		EmailBody:    *suspensionReason,
		EmailAddress: *profile.PrimaryEmailAddress,
		PrimaryPhone: *profile.PrimaryPhone,
	}
	err = s.engagement.NotifySupplierOnSuspension(ctx, supplierEmailPayload)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	return true, nil

}

// EDIUserLogin used to login a user to EDI (Portal Authserver) and return their
// EDI (Portal Authserver) profile
func (s SupplierUseCasesImpl) EDIUserLogin(
	ctx context.Context,
	username, password *string,
) (*profileutils.EDIUserProfile, error) {
	_, span := tracer.Start(ctx, "EDIUserLogin")
	defer span.End()

	if username == nil || password == nil {
		return nil, exceptions.InvalidCredentialsError()
	}

	ediClient, err := s.baseExt.LoginClient(*username, *password)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := s.baseExt.FetchUserProfile(ediClient)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("cannot retrieve EDI user profile: %w", err)
	}

	return userProfile, nil

}

// CoreEDIUserLogin used to login a user to EDI (Core Authserver) and return their EDI
// EDI (Core Authserver) profile
func (s SupplierUseCasesImpl) CoreEDIUserLogin(
	ctx context.Context,
	username, password string,
) (*profileutils.EDIUserProfile, error) {
	_, span := tracer.Start(ctx, "CoreEDIUserLogin")
	defer span.End()

	if username == "" || password == "" {
		return nil, exceptions.InvalidCredentialsError()
	}

	ediClient, err := s.baseExt.LoginClient(username, password)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, fmt.Errorf("cannot initialize edi client with supplied credentials: %w", err)
	}

	userProfile, err := s.baseExt.FetchUserProfile(ediClient)
	if err != nil {
		utils.RecordSpanError(span, err)
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
func (s SupplierUseCasesImpl) SupplierEDILogin(
	ctx context.Context,
	username string,
	password string,
	sladeCode string,
) (*dto.SupplierLogin, error) {
	ctx, span := tracer.Start(ctx, "SupplierEDILogin")
	defer span.End()

	var resp dto.SupplierLogin

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	supplier, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	accType := profileutils.AccountTypeIndividual
	supplier.AccountType = &accType
	supplier.UnderOrganization = true

	ediUserProfile, err := func(sladeCode string) (*profileutils.EDIUserProfile, error) {
		var ediUserProfile *profileutils.EDIUserProfile
		var err error

		switch sladeCode {
		case SavannahSladeCode:
			// login to core
			ediUserProfile, err = s.CoreEDIUserLogin(ctx, username, password)
			if err != nil {
				utils.RecordSpanError(span, err)
				supplier.IsOrganizationVerified = false
				return nil, fmt.Errorf("cannot get Core  user profile: %w", err)
			}

			if ediUserProfile == nil {
				return nil, fmt.Errorf("edi user profile not found")
			}

		default:
			//Login to portal edi
			ediUserProfile, err = s.EDIUserLogin(ctx, &username, &password)
			if err != nil {
				utils.RecordSpanError(span, err)
				supplier.IsOrganizationVerified = false
				return nil, fmt.Errorf("cannot get EDI user profile: %w", err)
			}

			if ediUserProfile == nil {
				return nil, fmt.Errorf("edi user profile not found")
			}

		}
		return ediUserProfile, nil
	}(sladeCode)
	if err != nil {
		utils.RecordSpanError(span, err)
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

	if orgSladeCode == SavannahSladeCode {
		supplier.EDIUserProfile = ediUserProfile
		supplier.IsOrganizationVerified = true
		supplier.SladeCode = sladeCode
		supplier.Active = true
		supplier.KYCSubmitted = true
		supplier.PartnerSetupComplete = true
		supplier.OrganizationName = SavannahOrgName

		if err := s.repo.UpdateSupplierProfile(ctx, *supplier.ProfileID, supplier); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		if err := s.profile.UpdatePermissions(ctx, profileutils.DefaultAdminPermissions); err != nil {
			utils.RecordSpanError(span, err)
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

	filter := []*dto.BusinessPartnerFilterInput{
		{
			SladeCode: &sladeCode,
		},
	}

	partner, err := s.chargemaster.FindProvider(ctx, nil, filter, nil)
	if err != nil {
		utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}()

	go func() {
		pro := func() error {
			return s.engagement.ResolveDefaultNudgeByTitle(
				ctx,
				uid,
				feedlib.FlavourPro,
				PartnerAccountSetupNudgeTitle,
			)
		}
		if err := backoff.Retry(
			pro,
			backoff.NewExponentialBackOff(),
		); err != nil {
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}()

	if businessPartner.Parent != nil {
		supplier.HasBranches = true
		supplier.ParentOrganizationID = *businessPartner.Parent

		partner, err := s.chargemaster.FetchProviderByID(ctx, *businessPartner.Parent)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, exceptions.FindProviderError(err)
		}

		supplier.OrganizationName = partner.Name

		if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, supplier); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		// fetch all locations of the business partner
		filter := []*dto.BranchFilterInput{
			{
				ParentOrganizationID: &supplier.ParentOrganizationID,
			},
		}

		brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
		if len(brs.Edges) == 0 || err != nil {
			utils.RecordSpanError(span, err)
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
				utils.RecordSpanError(span, err)
				return nil, err
			}

			resp.Supplier = spr
			return &resp, nil

		}
		return nil, exceptions.InternalServerError(nil)

	}

	// set the main branch as the supplier's location
	supplier.OrganizationName = businessPartner.Name
	loc := profileutils.Location{
		ID:   businessPartner.ID,
		Name: businessPartner.Name,
	}
	supplier.Location = &loc

	if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, supplier); err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// set the supplier profile in the final response object
	resp.Supplier = supplier

	return &resp, nil
}

// SupplierSetDefaultLocation updates the default location ot the supplier by the given location id
func (s SupplierUseCasesImpl) SupplierSetDefaultLocation(
	ctx context.Context,
	locationID string,
) (*profileutils.Supplier, error) {
	ctx, span := tracer.Start(ctx, "SupplierSetDefaultLocation")
	defer span.End()

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	// fetch the branches of the provider filtered by ParentOrganizationID
	filter := []*dto.BranchFilterInput{
		{
			ParentOrganizationID: &sup.ParentOrganizationID,
		},
	}

	brs, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.FindProviderError(err)
	}

	branch := func(brs *dto.BranchConnection, location string) *dto.BranchEdge {
		for _, b := range brs.Edges {
			if b.Node.ID == location {
				return b
			}
		}
		return nil
	}(brs, locationID)

	if branch != nil {
		loc := profileutils.Location{
			ID:              branch.Node.ID,
			Name:            branch.Node.Name,
			BranchSladeCode: &branch.Node.BranchSladeCode,
		}
		sup.Location = &loc

		if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		// refetch the supplier profile and return it
		return s.FindSupplierByUID(ctx)
	}

	return nil, fmt.Errorf("unable to get location of id %v : %v", locationID, err)
}

// FetchSupplierAllowedLocations retrieves all the locations that the user in context can work on.
func (s *SupplierUseCasesImpl) FetchSupplierAllowedLocations(
	ctx context.Context,
) (*dto.BranchConnection, error) {
	ctx, span := tracer.Start(ctx, "FetchSupplierAllowedLocations")
	defer span.End()

	supplier, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// fetch the branches of the provider filtered by ParentOrganizationID
	filter := []*dto.BranchFilterInput{
		{
			ParentOrganizationID: &supplier.ParentOrganizationID,
		},
	}

	branchConnection, err := s.chargemaster.FindBranch(ctx, nil, filter, nil)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, exceptions.FindProviderError(err)
	}

	if supplier.Location != nil {
		var resp dto.BranchConnection
		resp.PageInfo = branchConnection.PageInfo

		newEdges := []*dto.BranchEdge{}
		for _, edge := range branchConnection.Edges {
			if edge.Node.ID == supplier.Location.ID {
				loc := &dto.BranchEdge{
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
	return branchConnection, nil
}

// PublishKYCNudge pushes a KYC nudge to the user feed
func (s *SupplierUseCasesImpl) PublishKYCNudge(
	ctx context.Context,
	uid string,
	partner *profileutils.PartnerType,
	account *profileutils.AccountType,
) error {
	ctx, span := tracer.Start(ctx, "PublishKYCNudge")
	defer span.End()

	if partner == nil || account == nil {
		return exceptions.PublishKYCNudgeError(
			fmt.Errorf("expected `partner` to be defined and to be valid"),
		)
	}

	if *partner == profileutils.PartnerTypeConsumer {
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

	nudge := feedlib.Nudge{
		ID:             ksuid.New().String(),
		SequenceNumber: int(time.Now().Unix()),
		Visibility:     feedlib.VisibilityShow,
		Status:         feedlib.StatusPending,
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Title:          title,
		Text:           text,
		Links: []feedlib.Link{
			{
				ID:          ksuid.New().String(),
				URL:         feedlib.LogoURL,
				LinkType:    feedlib.LinkTypePngImage,
				Title:       "KYC",
				Description: fmt.Sprintf("KYC for %v", partner.String()),
				Thumbnail:   feedlib.LogoURL,
			},
		},
		Actions: []feedlib.Action{
			{
				ID:             ksuid.New().String(),
				SequenceNumber: int(time.Now().Unix()),
				Name: strings.ToUpper(fmt.Sprintf(
					"COMPLETE_%v_%v_KYC",
					account.String(),
					partner.String(),
				)),
				ActionType:     feedlib.ActionTypePrimary,
				Handling:       feedlib.HandlingFullPage,
				AllowAnonymous: false,
				Icon: feedlib.Link{
					ID:          ksuid.New().String(),
					URL:         feedlib.LogoURL,
					LinkType:    feedlib.LinkTypePngImage,
					Title:       title,
					Description: text,
					Thumbnail:   feedlib.LogoURL,
				},
			},
		},
		Users:  []string{uid},
		Groups: []string{uid},
		NotificationChannels: []feedlib.Channel{
			feedlib.ChannelEmail,
			feedlib.ChannelFcm,
		},
		NotificationBody: feedlib.NotificationBody{
			PublishMessage: "Kindly complete your KYC details and await approval.",
			ResolveMessage: "Thank you for adding your KYC details.",
		},
	}

	resp, err := s.engagement.PublishKYCNudge(ctx, uid, nudge)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.PublishKYCNudgeError(err)
	}

	// Status conflict means a similar KYC nudge already exists
	if resp.StatusCode == http.StatusConflict {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		if err := s.SaveProfileNudge(ctx, &nudge); err != nil {
			utils.RecordSpanError(span, err)
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
	ctx, span := tracer.Start(ctx, "PublishKYCFeedItem")
	defer span.End()

	for _, uid := range uids {
		payload := feedlib.Item{
			ID:             ksuid.New().String(),
			SequenceNumber: int(time.Now().Unix()),
			Expiry:         time.Now().Add(time.Hour * futureHours),
			Persistent:     true,
			Status:         feedlib.StatusPending,
			Visibility:     feedlib.VisibilityShow,
			Author:         "Be.Well Team",
			Label:          "KYC",
			Tagline:        "Process incoming KYC",
			Text:           "Review KYC for the partner and either approve or reject",
			TextType:       feedlib.TextTypeMarkdown,
			Icon: feedlib.Link{
				ID:          ksuid.New().String(),
				URL:         feedlib.LogoURL,
				LinkType:    feedlib.LinkTypePngImage,
				Title:       "KYC Review",
				Description: "Review KYC for the partner and either approve or reject",
				Thumbnail:   feedlib.LogoURL,
			},
			Timestamp: time.Now(),
			Actions: []feedlib.Action{
				{
					ID:             ksuid.New().String(),
					SequenceNumber: int(time.Now().Unix()),
					Name:           "Review KYC details",
					Icon: feedlib.Link{
						ID:          ksuid.New().String(),
						URL:         feedlib.LogoURL,
						LinkType:    feedlib.LinkTypePngImage,
						Title:       "Review KYC details",
						Description: "Review and approve or reject KYC details for the supplier",
						Thumbnail:   feedlib.LogoURL,
					},
					ActionType:     feedlib.ActionTypePrimary,
					Handling:       feedlib.HandlingFullPage,
					AllowAnonymous: false,
				},
			},
			Links: []feedlib.Link{
				{
					ID:          ksuid.New().String(),
					URL:         feedlib.LogoURL,
					LinkType:    feedlib.LinkTypePngImage,
					Title:       "KYC process request",
					Description: "Process KYC request",
					Thumbnail:   feedlib.LogoURL,
				},
			},

			Summary: "Process incoming KYC",
			Users:   uids,
			NotificationChannels: []feedlib.Channel{
				feedlib.ChannelFcm,
				feedlib.ChannelEmail,
				feedlib.ChannelSms,
			},
		}
		resp, err := s.engagement.PublishKYCFeedItem(ctx, uid, payload)
		if err != nil {
			utils.RecordSpanError(span, err)
			return fmt.Errorf("unable to publish kyc admin notification feed item : %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf(
				"unable to publish kyc admin notification feed item. unexpected status code  %v",
				resp.StatusCode,
			)
		}
	}

	return nil
}

// SaveProfileNudge stages nudges published from this service. These nudges will be
// referenced later to support some specialized use-case. A nudge will be uniquely
// identified by its id and sequenceNumber
func (s *SupplierUseCasesImpl) SaveProfileNudge(
	ctx context.Context,
	nudge *feedlib.Nudge,
) error {
	ctx, span := tracer.Start(ctx, "SaveProfileNudge")
	defer span.End()

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
func (s *SupplierUseCasesImpl) SaveKYCResponseAndNotifyAdmins(
	ctx context.Context,
	sup *profileutils.Supplier,
) error {
	ctx, span := tracer.Start(ctx, "SaveKYCResponseAndNotifyAdmins")
	defer span.End()

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.UserNotFoundError(err)
	}
	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}

	if profile.PrimaryEmailAddress == nil {
		return fmt.Errorf("supplier does not have a primary email address")
	}

	if err := s.repo.UpdateSupplierProfile(ctx, profile.ID, sup); err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	if err := s.StageKYCProcessingRequest(ctx, sup); err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	supplierEmailPayload := dto.EmailNotificationPayload{
		SupplierName: *profile.UserBioData.FirstName,
		PartnerType:  string(sup.PartnerType),
		AccountType:  string(*sup.AccountType),
		SubjectTitle: supplierEmailSubjectTitle,
		EmailBody:    supplierEmailBody,
		EmailAddress: *profile.PrimaryEmailAddress,
		PrimaryPhone: *profile.PrimaryPhone,
	}
	err = s.engagement.SendAlertToSupplier(ctx, supplierEmailPayload)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	adminEmailPayload := dto.EmailNotificationPayload{
		SupplierName: *profile.UserBioData.FirstName,
		PartnerType:  string(sup.PartnerType),
		AccountType:  string(*sup.AccountType),
		SubjectTitle: emailKYCSubject,
		EmailBody:    adminEmailBody,
		EmailAddress: *profile.PrimaryEmailAddress,
		PrimaryPhone: *profile.PrimaryPhone,
	}
	err = s.engagement.NotifyAdmins(ctx, adminEmailPayload)
	if err != nil {
		utils.RecordSpanError(span, err)
		return err
	}

	go func() {
		op := func() error {
			a, err := s.repo.FetchAdminUsers(ctx)
			if err != nil {
				utils.RecordSpanError(span, err)
				return err
			}
			var uids []string
			for _, u := range a {
				uids = append(uids, u.ID)
			}

			return s.PublishKYCFeedItem(ctx, uids...)
		}

		if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
			utils.RecordSpanError(span, err)
			logrus.Error(err)
		}
	}()

	return nil
}

// StageKYCProcessingRequest saves kyc processing requests
func (s *SupplierUseCasesImpl) StageKYCProcessingRequest(
	ctx context.Context,
	sup *profileutils.Supplier,
) error {
	ctx, span := tracer.Start(ctx, "StageKYCProcessingRequest")
	defer span.End()

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
func (s *SupplierUseCasesImpl) AddIndividualRiderKyc(
	ctx context.Context,
	input domain.IndividualRider,
) (*domain.IndividualRider, error) {
	ctx, span := tracer.Start(ctx, "AddIndividualRiderKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}
	if !sup.KYCSubmitted {
		if !input.IdentificationDoc.IdentificationDocType.IsValid() {
			return nil, exceptions.WrongEnumTypeError(
				input.IdentificationDoc.IdentificationDocType.String(),
			)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		// Notify the supplier and admins
		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationRiderKyc adds KYC for an organization rider
func (s *SupplierUseCasesImpl) AddOrganizationRiderKyc(
	ctx context.Context,
	input domain.OrganizationRider,
) (*domain.OrganizationRider, error) {
	ctx, span := tracer.Start(ctx, "AddOrganizationRiderKyc")
	defer span.End()

	if !input.OrganizationTypeName.IsValid() {
		return nil, fmt.Errorf(
			"invalid `OrganizationTypeName` provided : %v",
			input.OrganizationTypeName,
		)
	}

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualPractitionerKyc adds KYC for an individual practitioner
func (s *SupplierUseCasesImpl) AddIndividualPractitionerKyc(
	ctx context.Context,
	input domain.IndividualPractitioner,
) (*domain.IndividualPractitioner, error) {
	ctx, span := tracer.Start(ctx, "AddIndividualPractitionerKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()

}

// AddOrganizationPractitionerKyc adds KYC for an organization practitioner
func (s *SupplierUseCasesImpl) AddOrganizationPractitionerKyc(
	ctx context.Context,
	input domain.OrganizationPractitioner,
) (*domain.OrganizationPractitioner, error) {
	ctx, span := tracer.Start(ctx, "AddOrganizationPractitionerKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.OrganizationTypeName.IsValid() {
			return nil, fmt.Errorf(
				"invalid `OrganizationTypeName` provided : %v",
				input.OrganizationTypeName.String(),
			)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationProviderKyc adds KYC for an organization provider
func (s *SupplierUseCasesImpl) AddOrganizationProviderKyc(
	ctx context.Context,
	input domain.OrganizationProvider,
) (*domain.OrganizationProvider, error) {
	ctx, span := tracer.Start(ctx, "AddOrganizationProviderKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.OrganizationTypeName.IsValid() {
			return nil, fmt.Errorf(
				"invalid `OrganizationTypeName` provided : %v",
				input.OrganizationTypeName.String(),
			)
		}

		for _, practiceService := range input.PracticeServices {
			if !practiceService.IsValid() {
				return nil, fmt.Errorf(
					"invalid `PracticeService` provided : %v",
					practiceService.String(),
				)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualPharmaceuticalKyc adds KYC for an individual Pharmaceutical kyc
func (s *SupplierUseCasesImpl) AddIndividualPharmaceuticalKyc(
	ctx context.Context,
	input domain.IndividualPharmaceutical,
) (*domain.IndividualPharmaceutical, error) {
	ctx, span := tracer.Start(ctx, "AddIndividualPharmaceuticalKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.IdentificationDoc.IdentificationDocType.IsValid() {
			return nil, fmt.Errorf(
				"invalid `IdentificationDocType` provided : %v",
				input.IdentificationDoc.IdentificationDocType.String(),
			)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationPharmaceuticalKyc adds KYC for a pharmacy organization
func (s *SupplierUseCasesImpl) AddOrganizationPharmaceuticalKyc(
	ctx context.Context,
	input domain.OrganizationPharmaceutical,
) (*domain.OrganizationPharmaceutical, error) {
	ctx, span := tracer.Start(ctx, "AddOrganizationPharmaceuticalKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.OrganizationTypeName.IsValid() {
			return nil, fmt.Errorf(
				"invalid `OrganizationTypeName` provided : %v",
				input.OrganizationTypeName.String(),
			)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualCoachKyc adds KYC for an individual coach
func (s *SupplierUseCasesImpl) AddIndividualCoachKyc(
	ctx context.Context,
	input domain.IndividualCoach,
) (*domain.IndividualCoach, error) {
	ctx, span := tracer.Start(ctx, "AddIndividualCoachKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return nil, err
	}

	if !sup.KYCSubmitted {
		if !input.IdentificationDoc.IdentificationDocType.IsValid() {
			return nil, fmt.Errorf(
				"invalid `IdentificationDocType` provided : %v",
				input.IdentificationDoc.IdentificationDocType.String(),
			)
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
			AccreditationID:         input.AccreditationID,
			AccreditationUploadID:   input.AccreditationUploadID,
		}

		kycAsMap, err := s.parseKYCAsMap(kyc)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationCoachKyc adds KYC for an organization coach
func (s *SupplierUseCasesImpl) AddOrganizationCoachKyc(
	ctx context.Context,
	input domain.OrganizationCoach,
) (*domain.OrganizationCoach, error) {
	ctx, span := tracer.Start(ctx, "AddOrganizationCoachKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddIndividualNutritionKyc adds KYC for an individual nutritionist
func (s *SupplierUseCasesImpl) AddIndividualNutritionKyc(
	ctx context.Context,
	input domain.IndividualNutrition,
) (*domain.IndividualNutrition, error) {
	ctx, span := tracer.Start(ctx, "AddIndividualNutritionKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// AddOrganizationNutritionKyc adds kyc for a nutritionist organisation
func (s *SupplierUseCasesImpl) AddOrganizationNutritionKyc(
	ctx context.Context,
	input domain.OrganizationNutrition,
) (*domain.OrganizationNutrition, error) {
	ctx, span := tracer.Start(ctx, "AddOrganizationNutritionKyc")
	defer span.End()

	sup, err := s.FindSupplierByUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
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
			utils.RecordSpanError(span, err)
			return nil, fmt.Errorf("cannot marshal kyc to json")
		}

		sup.SupplierKYC = kycAsMap
		sup.KYCSubmitted = true

		if err := s.SaveKYCResponseAndNotifyAdmins(ctx, sup); err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		return &kyc, nil
	}

	return nil, exceptions.SupplierKYCAlreadySubmittedNotFoundError()
}

// FetchKYCProcessingRequests fetches a list of all unprocessed kyc approval requests
func (s *SupplierUseCasesImpl) FetchKYCProcessingRequests(
	ctx context.Context,
) ([]*domain.KYCRequest, error) {
	ctx, span := tracer.Start(ctx, "FetchKYCProcessingRequests")
	defer span.End()

	return s.repo.FetchKYCProcessingRequests(ctx)
}

// SendKYCEmail will send a KYC processing request email to the supplier
func (s *SupplierUseCasesImpl) SendKYCEmail(ctx context.Context, text, emailaddress string) error {
	ctx, span := tracer.Start(ctx, "SendKYCEmail")
	defer span.End()

	return s.engagement.SendMail(ctx, emailaddress, text, emailKYCSubject)
}

// ProcessKYCRequest transitions a kyc request to a given state
func (s *SupplierUseCasesImpl) ProcessKYCRequest(
	ctx context.Context,
	id string,
	status domain.KYCProcessStatus,
	rejectionReason *string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "ProcessKYCRequest")
	defer span.End()

	reviewerProfile, err := s.profile.UserProfile(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	s.repo.CheckIfAdmin(reviewerProfile)

	KYCRequest, err := s.repo.FetchKYCProcessingRequestByID(ctx, id)
	if err != nil {
		utils.RecordSpanError(span, err)
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
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to update KYC request record: %v", err)
	}

	supplierProfile, err := s.profile.GetProfileByID(
		ctx,
		KYCRequest.SupplierRecord.ProfileID,
	)
	if err != nil {
		utils.RecordSpanError(span, err)
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
				utils.RecordSpanError(span, err)
				logrus.Error(fmt.Errorf("unable to create erp supplier account: %v", err))
			}
		}()

		email = s.generateProcessKYCApprovalEmailTemplate()
		message = "Your KYC details have been reviewed and approved. We look forward to working with you. For enquiries call us on 0790360360"

		supplier, err := s.FindSupplierByUID(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return false, err
		}

		supplier.Active = true

		if err := s.repo.UpdateSupplierProfile(
			ctx,
			supplierProfile.ID,
			supplier,
		); err != nil {
			utils.RecordSpanError(span, err)
			return false, err
		}

	case domain.KYCProcessStatusRejected:
		email = s.generateProcessKYCRejectionEmailTemplate(*rejectionReason)
		message = "Your KYC details have been reviewed and have not been approved. Please check your email for detailed information. For enquiries call us on 0790360360"

	}

	nudgeTitle := fmt.Sprintf(
		PublishKYCNudgeTitle,
		strings.ToLower(string(KYCRequest.ReqPartnerType)),
	)
	supplierVerifiedUIDs := supplierProfile.VerifiedUIDS
	go func() {
		for _, UID := range supplierVerifiedUIDs {
			if err = s.engagement.ResolveDefaultNudgeByTitle(
				ctx,
				UID,
				feedlib.FlavourPro,
				nudgeTitle,
			); err != nil {
				utils.RecordSpanError(span, err)
				logrus.Print(err)
			}
		}
	}()

	supplierEmails := func(profile *profileutils.UserProfile) []string {
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
			utils.RecordSpanError(span, err)
			return false, fmt.Errorf("unable to send KYC processing email: %w", err)
		}
	}

	supplierPhones := func(profile *profileutils.UserProfile) []string {
		var phones []string
		phones = append(phones, *profile.PrimaryPhone)
		phones = append(phones, profile.SecondaryPhoneNumbers...)
		return phones
	}(supplierProfile)

	if err := s.engagement.SendSMS(ctx, supplierPhones, message); err != nil {
		utils.RecordSpanError(span, err)
		return false, fmt.Errorf("unable to send KYC processing message: %w", err)
	}

	return true, nil
}

// RetireKYCRequest retires the KYC process request of a supplier
func (s *SupplierUseCasesImpl) RetireKYCRequest(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "RetireKYCRequest")
	defer span.End()

	uid, err := s.baseExt.GetLoggedInUserUID(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return exceptions.UserNotFoundError(err)
	}

	profile, err := s.repo.GetUserProfileByUID(ctx, uid, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	sup, err := s.repo.GetSupplierProfileByProfileID(ctx, profile.ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		// this is a wrapped error. No need to wrap it again
		return err
	}
	if err := s.repo.RemoveKYCProcessingRequest(ctx, sup.ID); err != nil {
		utils.RecordSpanError(span, err)
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

func (s *SupplierUseCasesImpl) generateProcessKYCRejectionEmailTemplate(reason string) string {
	type rejectionData struct {
		Reason string
	}
	t := template.Must(template.New("rejectionKYCEmail").Parse(utils.ProcessKYCRejectionEmail))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, rejectionData{Reason: reason})
	if err != nil {
		log.Fatalf("Error while generating KYC rejection email template: %s", err)
	}
	return buf.String()
}
