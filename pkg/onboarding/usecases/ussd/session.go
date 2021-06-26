package ussd

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// UpdateSessionLevel updates user current level of interaction with USSD
func (u *Impl) UpdateSessionLevel(ctx context.Context, level int, sessionID string) error {
	_, err := u.onboardingRepository.UpdateSessionLevel(ctx, sessionID, level)
	if err != nil {
		return err
	}
	return nil

}

//AddAITSessionDetails persists USSD details
func (u *Impl) AddAITSessionDetails(ctx context.Context, input *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
	phone, err := base.NormalizeMSISDN(*input.PhoneNumber)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}
	sessionDetails := &dto.SessionDetails{
		PhoneNumber: phone,
		SessionID:   input.SessionID,
		Level:       input.Level,
	}
	result, err := u.onboardingRepository.AddAITSessionDetails(ctx, sessionDetails)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetOrCreateSessionState is used to set or return a user session
func (u *Impl) GetOrCreateSessionState(ctx context.Context, payload *dto.SessionDetails) (*domain.USSDLeadDetails, error) {
	sessionDetails, err := u.onboardingRepository.GetAITSessionDetails(ctx, payload.SessionID)
	if err != nil {
		return nil, err
	}
	if sessionDetails == nil {
		payload.Level = 0
		sessionDetails, err = u.AddAITSessionDetails(ctx, payload)
		if err != nil {
			return nil, err
		}
	}
	return sessionDetails, nil
}
