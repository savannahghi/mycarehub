// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
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
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor() (*Interactor, error) {

	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	profile := usecases.NewProfileUseCase(fr)
	otp := otp.NewOTPService(fr)
	erp := erp.NewERPService(fr)
	chrg := chargemaster.NewChargeMasterUseCasesImpl(fr)
	engage := engagement.NewServiceEngagementImpl(fr)
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes)
	login := usecases.NewLoginUseCases(fr)
	survey := usecases.NewSurveyUseCases(fr)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier)

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
	}, nil
}
