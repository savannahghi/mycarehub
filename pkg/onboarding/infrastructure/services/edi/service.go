package edi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// internal apis definitions
const (
	// LinkCover ISC endpoint to link user cover
	LinkCover = "internal/link_cover"
)

// ServiceEdi defines the business logic required to interact with EDI
type ServiceEdi interface {
	LinkCover(
		ctx context.Context,
		phoneNumber string,
		uid string,
		pushToken []string,
	) (*http.Response, error)
}

// ServiceEDIImpl represents EDI usecases
type ServiceEDIImpl struct {
	EdiExt               extension.ISCClientExtension
	onboardingRepository repository.OnboardingRepository
	engagement           engagement.ServiceEngagement
}

// NewEdiService returns a new instance of edi implementations
func NewEdiService(
	edi extension.ISCClientExtension,
	r repository.OnboardingRepository,
	engagement engagement.ServiceEngagement,
) ServiceEdi {
	return &ServiceEDIImpl{
		EdiExt:               edi,
		onboardingRepository: r,
		engagement:           engagement,
	}
}

// LinkCover calls the `EDI` service to link a cover to a converted, verified slade member user profile.
func (e *ServiceEDIImpl) LinkCover(
	ctx context.Context,
	phoneNumber string,
	uid string,
	pushToken []string,
) (*http.Response, error) {
	userMarketingData, err := e.engagement.GetSladerData(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to query the user's marketing details :%w", err)
	}

	if userMarketingData != nil {
		sladeCode, err := strconv.Atoi(userMarketingData.PayerSladeCode)
		if err != nil {
			return nil, fmt.Errorf("failed to convert slade code to an int: %w", err)
		}
		payload := dto.CoverInput{
			PayerSladeCode: sladeCode,
			MemberNumber:   userMarketingData.MemberNumber,
			UID:            uid,
			PushToken:      pushToken,
		}

		return e.EdiExt.MakeRequest(
			ctx,
			http.MethodPost,
			LinkCover,
			payload,
		)
	}
	return nil, nil
}
