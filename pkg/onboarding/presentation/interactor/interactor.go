// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	Onboarding   usecases.ProfileUseCase
	Signup       usecases.SignUpUseCases
	Supplier     usecases.SupplierUseCases
	Login        usecases.LoginUseCases
	Survey       usecases.SurveyUseCases
	UserPIN      usecases.UserPINUseCases
	ERP          erp.ServiceERP
	ChargeMaster chargemaster.ServiceChargeMaster
	Engagement   engagement.ServiceEngagement
	Messaging    messaging.ServiceMessaging
	NHIF         usecases.NHIFUseCases
	PubSub       pubsubmessaging.ServicePubSub
	SMS          usecases.SMSUsecase
	AITUSSD      usecases.UssdUsecase
	CRM          hubspot.ServiceHubSpotInterface
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	fr repository.OnboardingRepository,
	profile usecases.ProfileUseCase,
	su usecases.SignUpUseCases,
	supplier usecases.SupplierUseCases,
	login usecases.LoginUseCases,
	survey usecases.SurveyUseCases,
	userpin usecases.UserPINUseCases,
	erp erp.ServiceERP,
	chrg chargemaster.ServiceChargeMaster,
	engage engagement.ServiceEngagement,
	mes messaging.ServiceMessaging,
	nhif usecases.NHIFUseCases,
	pubsub pubsubmessaging.ServicePubSub,
	sms usecases.SMSUsecase,
	aitussd usecases.UssdUsecase,
	crm hubspot.ServiceHubSpotInterface,
) (*Interactor, error) {

	return &Interactor{
		Onboarding:   profile,
		Signup:       su,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engagement:   engage,
		Messaging:    mes,
		NHIF:         nhif,
		PubSub:       pubsub,
		SMS:          sms,
		AITUSSD:      aitussd,
		CRM:          crm,
	}, nil
}
