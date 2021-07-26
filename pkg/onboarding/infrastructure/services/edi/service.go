package edi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
)

// internal apis definitions
const (
	// LinkCover ISC endpoint to link user cover
	LinkCover = "internal/link_cover"

	// CoverLinkingStatusCompleted ...
	CoverLinkingStatusCompleted = "coverlinking completed"
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

		resp, err := e.EdiExt.MakeRequest(
			ctx,
			http.MethodPost,
			LinkCover,
			payload,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to make an edi request for coverlinking: %w", err)
		}

		currentTime := time.Now()
		coverLinkingEvent := &dto.CoverLinkingEvent{
			ID:                    uuid.NewString(),
			CoverLinkingEventTime: &currentTime,
			CoverStatus:           CoverLinkingStatusCompleted,
			MemberNumber:          userMarketingData.MemberNumber,
			PhoneNumber:           userMarketingData.Properties.Phone,
		}

		if _, err := e.onboardingRepository.SaveCoverAutolinkingEvents(ctx, coverLinkingEvent); err != nil {
			log.Printf("failed to save coverlinking `completed` event: %v", err)
		}

		return resp, nil

	}
	return nil, nil
}
