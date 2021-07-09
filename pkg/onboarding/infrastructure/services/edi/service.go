package edi

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

const (
	LinkCover = "internal/link_cover"
)

// ServiceEdi defines the business logic required to interact with EDI
type ServiceEdi interface {
	LinkCover(
		ctx context.Context,
		phoneNumber string,
		uid string,
	) (*http.Response, error)
}

// ServiceEDIImpl represents EDI usecases
type ServiceEDIImpl struct {
	EdiExt               extension.ISCClientExtension
	onboardingRepository repository.OnboardingRepository
}

// NewEdiService returns a new instance of edi implementations
func NewEdiService(edi extension.ISCClientExtension, r repository.OnboardingRepository) *ServiceEDIImpl {
	return &ServiceEDIImpl{
		EdiExt:               edi,
		onboardingRepository: r,
	}
}

// LinkCover calls the `EDI` service to link a cover to a converted, verified slade member user profile.
func (e *ServiceEDIImpl) LinkCover(
	ctx context.Context,
	phoneNumber string,
	uid string,
) (*http.Response, error) {
	// Get the user marketing data
	userMarketingData, err := e.onboardingRepository.GetUserMarketingData(ctx, phoneNumber)
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
