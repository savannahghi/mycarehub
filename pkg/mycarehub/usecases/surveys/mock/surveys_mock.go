package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// SurveysMock is a mock of `Surveys` interface.
type SurveysMock struct {
	MockListSurveysFn                   func(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error)
	MockGetUserSurveyFormsFn            func(ctx context.Context, userID string) ([]*domain.UserSurvey, error)
	MockSendClientSurveyLinksFn         func(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error)
	MockListSurveyRespondentsFn         func(ctx context.Context, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyRespondentPage, error)
	MockGetSurveyResponseFn             func(ctx context.Context, input dto.SurveyResponseInput) ([]*domain.SurveyResponse, error)
	MockGetSurveysWithServiceRequestsFn func(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error)
	MockVerifySurveySubmissionFn        func(ctx context.Context, input dto.VerifySurveySubmissionInput) (bool, error)
	MockGetSurveyServiceRequestUserFn   func(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error)
}

// NewSurveysMock initializes a new instance of `Survey Mock` then mocking the case of success.
func NewSurveysMock() *SurveysMock {
	UUID := gofakeit.UUID()
	bs := gofakeit.BS()
	now := time.Now()
	return &SurveysMock{
		MockListSurveysFn: func(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error) {
			return []*domain.SurveyForm{
				{
					ProjectID: 1,
					XMLFormID: UUID,
					Name:      bs,
					EnketoID:  UUID,
				},
			}, nil
		},
		MockGetUserSurveyFormsFn: func(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
			return []*domain.UserSurvey{
				{
					ID:           UUID,
					Active:       false,
					Created:      now,
					Link:         bs,
					Title:        bs,
					Description:  bs,
					HasSubmitted: true,
					UserID:       UUID,
					Token:        bs,
					ProjectID:    1,
					FormID:       UUID,
					LinkID:       1,
					SubmittedAt:  now,
					ProgramID:    UUID,
				},
			}, nil
		},
		MockSendClientSurveyLinksFn: func(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error) {
			return true, nil
		},
		MockListSurveyRespondentsFn: func(ctx context.Context, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyRespondentPage, error) {
			return &domain.SurveyRespondentPage{
				SurveyRespondents: []*domain.SurveyRespondent{
					{
						ID:          UUID,
						Name:        bs,
						SubmittedAt: now,
						ProjectID:   1,
						SubmitterID: 1,
						FormID:      UUID,
					},
				},
				Pagination: domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			}, nil
		},
		MockGetSurveyResponseFn: func(ctx context.Context, input dto.SurveyResponseInput) ([]*domain.SurveyResponse, error) {
			return []*domain.SurveyResponse{
				{
					Question:     bs,
					QuestionType: bs,
					Answer:       []string{"yes"},
				},
			}, nil
		},
		MockGetSurveysWithServiceRequestsFn: func(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
			return []*dto.SurveysWithServiceRequest{{Title: "test", ProjectID: 1, LinkID: 1, FormID: "test"}}, nil
		},
		MockVerifySurveySubmissionFn: func(ctx context.Context, input dto.VerifySurveySubmissionInput) (bool, error) {
			return true, nil
		},
		MockGetSurveyServiceRequestUserFn: func(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error) {
			return &domain.SurveyServiceRequestUserPage{
				Users: []*domain.SurveyServiceRequestUser{
					{
						Name:             bs,
						FormID:           UUID,
						ProjectID:        1,
						SubmitterID:      1,
						SurveyName:       bs,
						ServiceRequestID: UUID,
						PhoneNumber:      "09999999",
					},
				},
				Pagination: domain.Pagination{
					Limit:       1,
					CurrentPage: 1,
				},
			}, nil
		},
	}
}

// ListSurveys mock the implementation of the ListSurveys method
func (m *SurveysMock) ListSurveys(ctx context.Context, projectID *int) ([]*domain.SurveyForm, error) {
	return m.MockListSurveysFn(ctx, projectID)
}

// GetUserSurveyForms mock the implementation of the GetUserSurveyForms method
func (m *SurveysMock) GetUserSurveyForms(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
	return m.MockGetUserSurveyFormsFn(ctx, userID)
}

// SendClientSurveyLinks mock the implementation of the SendClientSurveyLinks method
func (m *SurveysMock) SendClientSurveyLinks(ctx context.Context, facilityID *string, formID *string, projectID *int, filterParams *dto.ClientFilterParamsInput) (bool, error) {
	return m.MockSendClientSurveyLinksFn(ctx, facilityID, formID, projectID, filterParams)
}

// ListSurveyRespondents mock the implementation of the ListSurveyRespondents method
func (m *SurveysMock) ListSurveyRespondents(ctx context.Context, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyRespondentPage, error) {
	return m.MockListSurveyRespondentsFn(ctx, projectID, formID, paginationInput)
}

// GetSurveyResponse mock the implementation of the GetSurveyResponse method
func (m *SurveysMock) GetSurveyResponse(ctx context.Context, input dto.SurveyResponseInput) ([]*domain.SurveyResponse, error) {
	return m.MockGetSurveyResponseFn(ctx, input)
}

// GetSurveysWithServiceRequests mock the implementation of the GetSurveysWithServiceRequests method
func (m *SurveysMock) GetSurveysWithServiceRequests(ctx context.Context, facilityID string) ([]*dto.SurveysWithServiceRequest, error) {
	return m.MockGetSurveysWithServiceRequestsFn(ctx, facilityID)
}

// VerifySurveySubmission mock the implementation of the VerifySurveySubmission method
func (m *SurveysMock) VerifySurveySubmission(ctx context.Context, input dto.VerifySurveySubmissionInput) (bool, error) {
	return m.MockVerifySurveySubmissionFn(ctx, input)
}

// GetSurveyServiceRequestUser mock the implementation of the GetSurveyServiceRequestUser method
func (m *SurveysMock) GetSurveyServiceRequestUser(ctx context.Context, facilityID string, projectID int, formID string, paginationInput dto.PaginationsInput) (*domain.SurveyServiceRequestUserPage, error) {
	return m.MockGetSurveyServiceRequestUserFn(ctx, facilityID, projectID, formID, paginationInput)
}
