package mock

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SurveysMock is a mock of `Surveys` interface.
type SurveysMock struct {
	MockGetUsersWithSurveyServiceRequestFn func(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error)
}

// NewSurveysMock initializes a new instance of `Survey Mock` then mocking the case of success.
func NewSurveysMock() *SurveysMock {
	return &SurveysMock{
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

// GetSurveyServiceRequestUser mocks the case of success.
func (m *SurveysMock) GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error) {
	return m.MockGetUsersWithSurveyServiceRequestFn(ctx, facilityID, projectID, formID, paginationInput)
}
