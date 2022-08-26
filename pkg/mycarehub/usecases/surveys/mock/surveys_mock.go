package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
)

// SurveysMock is a mock of `Surveys` interface.
type SurveysMock struct {
	MockGetSurveysWithServiceRequestsFn func(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error)
}

// NewSurveysMock initializes a new instance of `Survey Mock` then mocking the case of success.
func NewSurveysMock() *SurveysMock {
	return &SurveysMock{
		MockGetSurveysWithServiceRequestsFn: func(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
			return []*dto.SurveysWithServiceRequest{
				{
					Title:     "test",
					ProjectID: 1,
					LinkID:    1,
					FormID:    "test",
				},
			}, nil
		},
	}
}

// GetSurveysWithServiceRequests mocks the case of success of getting surveys with service requests.
func (m *SurveysMock) GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
	return m.MockGetSurveysWithServiceRequestsFn(ctx, facilityID)
}
