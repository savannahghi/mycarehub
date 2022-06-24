package mock

import (
	"context"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
)

// SurveysMock mocks the surveys service
type SurveysMock struct {
	MockMakeRequestFn              func(ctx context.Context, payload surveys.RequestHelperPayload) (*http.Response, error)
	MockListSurveyFormsFn          func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	MockGetSurveyFormFn            func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error)
	MockGeneratePublicAccessLinkFn func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error)
	MockGetSubmissionsFn           func(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error)
	MockDeletePublicAccessLinkFn   func(ctx context.Context, input dto.VerifySurveySubmissionInput) error
	MockListSubmittersFn           func(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error)
	MockListPublicAccessLinksFn    func(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error)
}

// NewSurveysMock initializes the surveys mock service
func NewSurveysMock() *SurveysMock {
	return &SurveysMock{

		MockMakeRequestFn: func(ctx context.Context, payload surveys.RequestHelperPayload) (*http.Response, error) {
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
		MockGeneratePublicAccessLinkFn: func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
			return &dto.SurveyPublicLink{
				Once:        true,
				ID:          2,
				DisplayName: gofakeit.Name(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   nil,
				Token:       gofakeit.UUID(),
			}, nil
		},
		MockListPublicAccessLinksFn: func(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error) {
			return []*dto.SurveyPublicLink{
				{
					Once:        true,
					ID:          2,
					DisplayName: gofakeit.Name(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					DeletedAt:   nil,
					Token:       gofakeit.UUID(),
				},
			}, nil
		},
		MockGetSubmissionsFn: func(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
			return []domain.Submission{
				{
					InstanceID:  gofakeit.UUID(),
					SubmitterID: 10,
					DeviceID:    gofakeit.UUID(),
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
					ReviewState: gofakeit.UUID(),
					Submitter: domain.Submitter{
						ID:          10,
						Type:        gofakeit.BeerAlcohol(),
						DisplayName: gofakeit.BeerBlg(),
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
						DeletedAt:   time.Time{},
					},
				},
			}, nil
		},
		MockDeletePublicAccessLinkFn: func(ctx context.Context, input dto.VerifySurveySubmissionInput) error {
			return nil
		},
		MockListSubmittersFn: func(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
			return []domain.Submitter{
				{
					ID:          10,
					Type:        "test",
					DisplayName: "test",
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
					DeletedAt:   time.Time{},
				},
			}, nil
		},
	}
}

// ListSurveyForms lists the survey forms for the given project
func (s *SurveysMock) ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
	return s.MockListSurveyFormsFn(ctx, projectID)
}

// MakeRequest makes a request to the surveys service
func (s *SurveysMock) MakeRequest(ctx context.Context, payload surveys.RequestHelperPayload) (*http.Response, error) {
	return s.MockMakeRequestFn(ctx, payload)
}

// GetSurveyForm gets the survey form for the given project and form ID
func (s *SurveysMock) GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
	return s.MockGetSurveyFormFn(ctx, projectID, formID)
}

// GeneratePublicAccessLink generates a public access link for the given survey
func (s *SurveysMock) GeneratePublicAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
	return s.MockGeneratePublicAccessLinkFn(ctx, input)
}

// GetSubmissions mocks the action of getting the submissions for the given survey
func (s *SurveysMock) GetSubmissions(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
	return s.MockGetSubmissionsFn(ctx, input)
}

// DeletePublicAccessLink mocks the implementation of deleting the public access link for the given survey
func (s *SurveysMock) DeletePublicAccessLink(ctx context.Context, input dto.VerifySurveySubmissionInput) error {
	return s.MockDeletePublicAccessLinkFn(ctx, input)
}

// ListSubmitters mocks the action of listing all the submitters of a given survey
func (s *SurveysMock) ListSubmitters(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
	return s.MockListSubmittersFn(ctx, projectID, formID)
}

// ListPublicAccessLinks returns a list of all public access links created for a particular form
func (s *SurveysMock) ListPublicAccessLinks(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error) {
	return s.MockListPublicAccessLinksFn(ctx, projectID, formID)
}
