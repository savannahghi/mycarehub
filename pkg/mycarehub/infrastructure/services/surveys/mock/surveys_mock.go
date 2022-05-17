package mock

import (
	"context"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SurveysMock mocks the surveys service
type SurveysMock struct {
	MockMakeRequestFn               func(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error)
	MockListSurveyFormsFn           func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	MockGetSurveyFormFn             func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error)
	MockGeneratePublickAccessLinkFn func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error)
}

// NewSurveysMock initializes the surveys mock service
func NewSurveysMock() *SurveysMock {
	return &SurveysMock{

		MockMakeRequestFn: func(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
			}, nil
		},
		MockListSurveyFormsFn: func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
			return []*domain.SurveyForm{
				{
					ProjectID: 2,
					Name:      gofakeit.Name(),
					EnketoID:  gofakeit.UUID(),
				},
			}, nil
		},
		MockGetSurveyFormFn: func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
			return &domain.SurveyForm{
				ProjectID: 2,
				Name:      gofakeit.Name(),
				EnketoID:  gofakeit.UUID(),
			}, nil
		},
		MockGeneratePublickAccessLinkFn: func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
			return &dto.SurveyPublicLink{
				Once:        true,
				ID:          2,
				DisplayName: gofakeit.Name(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   nil,
				Token:       gofakeit.UUID(),
				CSRF:        gofakeit.UUID(),
				ExpiresAt:   time.Now().Add(time.Hour * 24),
			}, nil
		},
	}
}

// ListSurveyForms lists the survey forms for the given project
func (s *SurveysMock) ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
	return s.MockListSurveyFormsFn(ctx, projectID)
}

// MakeRequest makes a request to the surveys service
func (s *SurveysMock) MakeRequest(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error) {
	return s.MockMakeRequestFn(ctx, payload)
}

// GetSurveyForm gets the survey form for the given project and form ID
func (s *SurveysMock) GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
	return s.MockGetSurveyFormFn(ctx, projectID, formID)
}

// GeneratePublickAccessLink generates a public access link for the given survey
func (s *SurveysMock) GeneratePublickAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
	return s.MockGeneratePublickAccessLinkFn(ctx, input)
}
