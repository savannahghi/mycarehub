package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	welcomeMessage = "CON Welcome to Be.Well"
	todoMessage    = "What would you like to do?"
)

//UssdUsecase represent the logic involved in receiving a USSD
type UssdUsecase interface {
	GenerateUSSD(payload *dto.SessionDetails) string
	CreateUSSDData(ctx context.Context, input *dto.EndSessionDetails) error
}

//UssdImpl represents usecase implementation object
type UssdImpl struct {
	baseExt              extension.BaseExtension
	onboardingRepository repository.OnboardingRepository
}

//NewUssdUsecases returns a new Ussd usecase
func NewUssdUsecases(repository repository.OnboardingRepository, ext extension.BaseExtension) UssdUsecase {
	return &UssdImpl{
		baseExt:              ext,
		onboardingRepository: repository,
	}
}

//GenerateUSSD generates the USSD response
func (u *UssdImpl) GenerateUSSD(payload *dto.SessionDetails) string {
	var resp string
	if len(payload.Text) == 0 {
		resp = welcomeMessage + ".\r\n"
		resp += todoMessage + "\r\n"
		resp += utils.Menu()
	} else {
		resp = utils.ResponseMenu(payload.Text)
	}
	return resp
}

//CreateUSSDData persists USSD details
func (u *UssdImpl) CreateUSSDData(ctx context.Context, input *dto.EndSessionDetails) error {
	phone, err := base.NormalizeMSISDN(*input.PhoneNumber)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}
	text := utils.GetTextValue(input.Input)
	sessionDetails := &dto.EndSessionDetails{
		PhoneNumber: phone,
		Input:       text,
		SessionID:   input.SessionID,
		Status:      input.Status,
	}
	err = u.onboardingRepository.AddIncomingUSSDData(ctx, sessionDetails)
	if err != nil {
		return err
	}
	return nil
}
