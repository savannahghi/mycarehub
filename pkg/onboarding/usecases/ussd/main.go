package ussd

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	// LoginUserState handles workflow required to authenticate/login a user
	LoginUserState = 0
	//HomeMenuState represents inner submenu once user is logged in
	HomeMenuState = 5
	// UserPINState represents workflows required to set a user PIN
	UserPINState = 50
	// UserPINResetState represents workflows required to reset a forgotten user PIN
	UserPINResetState = 10
	// EmptyInput is used to load a default menu when user has not supplied any input
	EmptyInput = ""
	// GoBackHomeInput represents the user intention to go back to the main menu
	GoBackHomeInput = "0"
)

//Usecase represent the logic involved in processing USSD requests
type Usecase interface {
	HandleResponseFromUSSDGateway(context context.Context, input *dto.SessionDetails) string
	HandleUserRegistration(ctx context.Context, sessionDetails *domain.USSDLeadDetails, userResponse string) string
	HandleHomeMenu(ctx context.Context, level int, session *domain.USSDLeadDetails, userResponse string) string
	CreateUsddUserProfile(ctx context.Context, phoneNumber string, PIN string, userProfile *dto.UserProfileInput) error
	// session usecases
	GetOrCreateSessionState(ctx context.Context, payload *dto.SessionDetails) (*domain.USSDLeadDetails, error)
	AddAITSessionDetails(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error)
	UpdateSessionLevel(ctx context.Context, level int, sessionID string) error
	// USSD PIN usecases
	HandleChangePIN(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string
	HandlePINReset(ctx context.Context, session *domain.USSDLeadDetails, userResponse string) string
	SetUSSDUserPin(ctx context.Context, phoneNumber string, PIN string) error
	ChangeUSSDUserPIN(ctx context.Context, phone string, pin string) (bool, error)
}

//Impl represents usecase implementation
type Impl struct {
	baseExt              extension.BaseExtension
	onboardingRepository repository.OnboardingRepository
	profile              usecases.ProfileUseCase
	pinUsecase           usecases.UserPINUseCases
	signUp               usecases.SignUpUseCases
	pinExt               extension.PINExtension
	pubsub               pubsubmessaging.ServicePubSub
}

//NewUssdUsecases returns a new USSD usecase
func NewUssdUsecases(
	repository repository.OnboardingRepository,
	ext extension.BaseExtension,
	profileUsecase usecases.ProfileUseCase,
	pinUsecase usecases.UserPINUseCases,
	signUp usecases.SignUpUseCases,
	pinExt extension.PINExtension,
	pubsub pubsubmessaging.ServicePubSub,
) Usecase {
	return &Impl{
		baseExt:              ext,
		onboardingRepository: repository,
		profile:              profileUsecase,
		pinUsecase:           pinUsecase,
		signUp:               signUp,
		pinExt:               pinExt,
		pubsub:               pubsub,
	}
}

//HandleResponseFromUSSDGateway receives and processes the USSD response from the USSD gateway
func (u *Impl) HandleResponseFromUSSDGateway(ctx context.Context, payload *dto.SessionDetails) string {
	sessionDetails, err := u.GetOrCreateSessionState(ctx, payload)
	if err != nil {
		return "END something wrong happened"
	}
	userResponse := utils.GetUserResponse(payload.Text)
	exists, err := u.profile.CheckPhoneExists(ctx, *payload.PhoneNumber)
	if err != nil {
		return "END something wrong happened"
	}
	if !exists {
		return u.HandleUserRegistration(ctx, sessionDetails, userResponse)
	}
	switch {
	case sessionDetails.Level == LoginUserState:
		return u.HandleLogin(ctx, sessionDetails, userResponse)
	case sessionDetails.Level == HomeMenuState:
		return u.HandleHomeMenu(ctx, HomeMenuState, sessionDetails, userResponse)
	case sessionDetails.Level >= UserPINState:
		return u.HandleChangePIN(ctx, sessionDetails, userResponse)
	case sessionDetails.Level >= UserPINResetState:
		return u.HandlePINReset(ctx, sessionDetails, userResponse)
	default:
		return "END something wrong happened"
	}

}
