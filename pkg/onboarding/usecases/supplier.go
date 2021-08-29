package usecases

import (
	"context"
	"fmt"

	"github.com/cenkalti/backoff"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/authorization/permission"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
	"github.com/sirupsen/logrus"

	pubsubmessaging "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/pubsub"
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

	// PartnerAccountSetupNudgeTitle is the title defined in the `engagement service`
	// for the `PartnerAccountSetupNudge`
	PartnerAccountSetupNudgeTitle = "Setup your partner account"

	// PublishKYCNudgeTitle is the title for the PublishKYCNudge.
	// It takes a partner type as an argument
	PublishKYCNudgeTitle = "Complete your %s KYC"
)

// SupplierUseCases represent the business logic required for management of suppliers
type SupplierUseCases interface {
	SetUpSupplier(ctx context.Context, accountType profileutils.AccountType) (*profileutils.Supplier, error)
}

// SupplierUseCasesImpl represents usecase implementation object
type SupplierUseCasesImpl struct {
	repo       repository.OnboardingRepository
	profile    ProfileUseCase
	engagement engagement.ServiceEngagement
	baseExt    extension.BaseExtension
	pubsub     pubsubmessaging.ServicePubSub
}

// NewSupplierUseCases returns a new a onboarding usecase
func NewSupplierUseCases(
	r repository.OnboardingRepository,
	p ProfileUseCase,
	eng engagement.ServiceEngagement,
	ext extension.BaseExtension,
	pubsub pubsubmessaging.ServicePubSub,
) SupplierUseCases {

	return &SupplierUseCasesImpl{
		repo:       r,
		profile:    p,
		engagement: eng,
		baseExt:    ext,
		pubsub:     pubsub,
	}
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

	// TODO: restore
	// go func(u string, pnt profileutils.PartnerType, acnt profileutils.AccountType) {
	// 	op := func() error {
	// 		return s.PublishKYCNudge(ctx, u, &pnt, &acnt)
	// 	}

	// 	if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
	// 		utils.RecordSpanError(span, err)
	// 		logrus.Error(err)
	// 	}
	// }(user.UID, sup.PartnerType, *sup.AccountType)

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
