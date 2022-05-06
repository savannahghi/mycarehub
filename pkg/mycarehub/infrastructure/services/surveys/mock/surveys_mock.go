package mock

import (
	"context"
	"net/http"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SurveysMock mocks the surveys service
type SurveysMock struct {
	MockListSurveyFormsFn func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	MockMakeRequestFn     func(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error)
}

// NewSurveysMock initializes the surveys mock service
func NewSurveysMock() *SurveysMock {
	return &SurveysMock{
		MockListSurveyFormsFn: func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
			return []*domain.SurveyForm{
				{
					ProjectID: 2,
					Name:      gofakeit.Name(),
				},
			}, nil
		},
		MockMakeRequestFn: func(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
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
