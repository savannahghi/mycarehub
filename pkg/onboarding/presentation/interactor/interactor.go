// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	Onboarding   usecases.ProfileUseCase
	Signup       usecases.SignUpUseCases
	Otp          otp.ServiceOTP
	Supplier     usecases.SupplierUseCases
	Login        usecases.LoginUseCases
	Survey       usecases.SurveyUseCases
	UserPIN      usecases.UserPINUseCases
	ERP          erp.ServiceERP
	ChargeMaster chargemaster.ServiceChargeMaster
	Engagement   engagement.ServiceEngagement
	Mailgun      mailgun.ServiceMailgun
	Messaging    messaging.ServiceMessaging
	NHIF         usecases.NHIFUseCases
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	fr repository.OnboardingRepository,
	profile usecases.ProfileUseCase,
	su usecases.SignUpUseCases,
	otp otp.ServiceOTP,
	supplier usecases.SupplierUseCases,
	login usecases.LoginUseCases,
	survey usecases.SurveyUseCases,
	userpin usecases.UserPINUseCases,
	erp erp.ServiceERP,
	chrg chargemaster.ServiceChargeMaster,
	engage engagement.ServiceEngagement,
	mg mailgun.ServiceMailgun,
	mes messaging.ServiceMessaging,
	nhif usecases.NHIFUseCases,
) (*Interactor, error) {

	return &Interactor{
		Onboarding:   profile,
		Signup:       su,
		Otp:          otp,
		Supplier:     supplier,
		Login:        login,
		Survey:       survey,
		UserPIN:      userpin,
		ERP:          erp,
		ChargeMaster: chrg,
		Engagement:   engage,
		Mailgun:      mg,
		Messaging:    mes,
		NHIF:         nhif,
	}, nil
}
