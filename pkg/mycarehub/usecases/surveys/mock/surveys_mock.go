package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SurveysMock is a mock of `Surveys` interface.
type SurveysMock struct {
	MockGetSurveysWithServiceRequestsFn    func(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error)
	MockGetUsersWithSurveyServiceRequestFn func(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error)
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
		MockGetUsersWithSurveyServiceRequestFn: func(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error) {
			return &domain.SurveyServiceRequestUserPage{
				Users: []*domain.SurveyServiceRequestUser{
					{
						Name:        "test",
						FormID:      "test",
						ProjectID:   1,
						SubmitterID: 1,
						SurveyName:  "test",
					},
				},
				Pagination: domain.Pagination{
					Limit:        10,
					CurrentPage:  1,
					Count:        20,
					TotalPages:   20,
					NextPage:     new(int),
					PreviousPage: new(int),
				},
			}, nil
		},
	}
}

// GetSurveysWithServiceRequests mocks the case of success of getting surveys with service requests.
func (m *SurveysMock) GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
	return m.MockGetSurveysWithServiceRequestsFn(ctx, facilityID)
}

// GetSurveyServiceRequestUser mocks the case of success.
func (m *SurveysMock) GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error) {
	return m.MockGetUsersWithSurveyServiceRequestFn(ctx, facilityID, projectID, formID, paginationInput)
}
