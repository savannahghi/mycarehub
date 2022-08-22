package surveys

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	mockSurveys "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys/mock"
	mockNotification "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
)

func TestUsecaseSurveysImpl_GetSurveyResponse(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.SurveyResponseInput
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.SurveyResponse
		wantErr bool
	}{
		{
			name: "happy case: get survey response",
			args: args{
				ctx: context.Background(),
				input: dto.SurveyResponseInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 1096,
				},
			},
			want:    []*domain.SurveyResponse{},
			wantErr: false,
		},
		{
			name: "sad case: fail to get submissions",
			args: args{
				ctx: context.Background(),
				input: dto.SurveyResponseInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 1096,
				},
			},
			want:    []*domain.SurveyResponse{},
			wantErr: true,
		},
		{
			name: "sad case: fail to get submission instanceID",
			args: args{
				ctx: context.Background(),
				input: dto.SurveyResponseInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 1096,
				},
			},
			want:    []*domain.SurveyResponse{},
			wantErr: true,
		},
		{
			name: "sad case: fail to get submission xml",
			args: args{
				ctx: context.Background(),
				input: dto.SurveyResponseInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 1096,
				},
			},
			want:    []*domain.SurveyResponse{},
			wantErr: true,
		},
		{
			name: "sad case: fail to get form xml",
			args: args{
				ctx: context.Background(),
				input: dto.SurveyResponseInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 1096,
				},
			},
			want:    []*domain.SurveyResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeSurveys := mockSurveys.NewSurveysMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeNotification := mockNotification.NewServiceNotificationMock()
			u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeDB, fakeNotification)

			if tt.name == "sad case: fail to get submissions" {
				fakeSurveys.MockGetSubmissionsFn = func(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
					return []domain.Submission{}, fmt.Errorf("failed to get submissions")
				}
			}

			if tt.name == "sad case: fail to get submission instanceID" {
				fakeSurveys.MockGetSubmissionsFn = func(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
					return []domain.Submission{
						{
							InstanceID:  "test",
							SubmitterID: 1,
							DeviceID:    "artghjkl",
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Now(),
							ReviewState: "good",
							Submitter: domain.Submitter{
								ID:          1,
								Type:        "test",
								DisplayName: gofakeit.Name(),
								CreatedAt:   time.Now(),
								UpdatedAt:   time.Now(),
								DeletedAt:   time.Now(),
							},
						},
					}, nil
				}
			}

			if tt.name == "sad case: fail to get submission xml" {
				fakeSurveys.MockGetSubmissionXMLFn = func(ctx context.Context, projectID int, formID, instanceID string) (map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to get submission xml")
				}
			}

			if tt.name == "sad case: fail to get form xml" {
				fakeSurveys.MockGetFormXMLFn = func(ctx context.Context, projectID int, formID, version string) (map[string]interface{}, error) {
					return nil, fmt.Errorf("failed to get form xml")
				}
			}

			_, err := u.GetSurveyResponse(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseSurveysImpl.GetSurveyResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseSurveysImpl_ListSurveys(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeSurveys := mockSurveys.NewSurveysMock()
	fakeNotification := mockNotification.NewServiceNotificationMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeDB, fakeNotification)
	projectID := 2

	type args struct {
		ctx       context.Context
		projectID *int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: list surveys",
			args: args{
				ctx:       context.Background(),
				projectID: &projectID,
			},
		},
		{
			name: "sad case: failed to list surveys",
			args: args{
				ctx:       context.Background(),
				projectID: &projectID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: failed to list surveys" {
				fakeSurveys.MockListSurveyFormsFn = func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
					return nil, fmt.Errorf("failed to list surveys")
				}
			}

			got, err := u.ListSurveys(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseSurveysImpl.ListSurveys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("UsecaseSurveysImpl.ListSurveys() = %v, want %v", got, tt.wantErr)
			}

		})
	}
}

func TestUsecaseSurveysImpl_GetUserSurveyForms(t *testing.T) {
	fakeSurveys := mockSurveys.NewSurveysMock()
	fakeDB := pgMock.NewPostgresMock()
	fakeNotification := mockNotification.NewServiceNotificationMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeDB, fakeNotification)

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get user survey forms",
			args: args{
				ctx:    context.Background(),
				userID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get user survey forms",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get user survey forms" {
				fakeDB.MockGetUserSurveyFormsFn = func(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error) {
					return nil, fmt.Errorf("failed to get user survey forms")
				}
			}
			got, err := u.GetUserSurveyForms(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseSurveysImpl.GetUserSurveyForms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UsecaseSurveysImpl.ListSurveys() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestUsecaseSurveysImpl_SendClientSurveyLinks(t *testing.T) {
	facilityID := uuid.New().String()
	formID := uuid.New().String()
	projectID := 2

	type args struct {
		ctx          context.Context
		facilityID   *string
		formID       *string
		projectID    *int
		filterParams *dto.ClientFilterParamsInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: send client survey links",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: skip client duplicate link",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "happy case: no clients to send surveys to",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: failed to get filtered clients",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to get survey form",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: failed to create survey",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user survey forms",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    true, // only report the error to sentry
			wantErr: false,
		},
		{
			name: "sad case: failed to notify user",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
				formID:     &formID,
				projectID:  &projectID,
				filterParams: &dto.ClientFilterParamsInput{
					ClientTypes: []enums.ClientType{enums.ClientTypePmtct},
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 20,
						UpperBound: 25,
					},
					Gender: []enumutils.Gender{enumutils.GenderMale},
				},
			},
			want:    true, // only report the error to sentry
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeSurveys := mockSurveys.NewSurveysMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeNotification := mockNotification.NewServiceNotificationMock()
			u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeDB, fakeNotification)
			if tt.name == "happy case: no clients to send surveys to" {
				fakeDB.MockGetClientsByFilterParamsFn = func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
					return []*domain.ClientProfile{}, nil
				}
			}

			if tt.name == "Sad case: unable to get user survey forms" {
				ID := uuid.New().String()
				fakeDB.MockGetClientsByFilterParamsFn = func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
					return []*domain.ClientProfile{
						{
							ID:   &ID,
							User: &domain.User{ID: &ID},
						},
					}, nil
				}
				fakeDB.MockGetUserSurveyFormsFn = func(ctx context.Context, params map[string]interface{}) ([]*domain.UserSurvey, error) {
					return nil, fmt.Errorf("failed to get user survey forms")
				}
			}

			if tt.name == "sad case: failed to get filtered clients" {
				fakeDB.MockGetClientsByFilterParamsFn = func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get filtered clients")
				}

			}
			if tt.name == "sad case: failed to get survey form" {
				fakeDB.MockGetClientsByFilterParamsFn = func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
					ID := uuid.NewString()
					return []*domain.ClientProfile{
						{
							ID: &ID,
							User: &domain.User{
								ID: &ID,
							},
							Active: false,
						},
					}, nil
				}
				fakeSurveys.MockGetSurveyFormFn = func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
					return nil, fmt.Errorf("failed to get survey form")
				}
			}

			if tt.name == "sad case: failed to create survey" {
				fakeDB.MockCreateUserSurveyFn = func(ctx context.Context, userSurvey []*dto.UserSurveyInput) error {
					return fmt.Errorf("failed to create survey")
				}
			}

			if tt.name == "sad case: failed to notify user" {
				fakeSurveys.MockGeneratePublicAccessLinkFn = func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
					return &dto.SurveyPublicLink{
						Once:        true,
						ID:          projectID,
						DisplayName: "Test",
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
						DeletedAt:   &time.Time{},
						Token:       "",
					}, nil
				}
				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify user")
				}
			}

			if tt.name == "sad case: failed to list public access link" {
				fakeSurveys.MockListPublicAccessLinksFn = func(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error) {
					return nil, fmt.Errorf("error listing public access link")
				}
			}

			got, err := u.SendClientSurveyLinks(tt.args.ctx, tt.args.facilityID, tt.args.formID, tt.args.projectID, tt.args.filterParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseSurveysImpl.SendClientSurveyLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseSurveysImpl.SendClientSurveyLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseSurveysImpl_VerifySurveySubmission(t *testing.T) {
	ctx := context.Background()

	fakeSurveys := mockSurveys.NewSurveysMock()
	fakeDB := pgMock.NewPostgresMock()
	fakeNotification := mockNotification.NewServiceNotificationMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeDB, fakeNotification)

	type args struct {
		ctx   context.Context
		input dto.VerifySurveySubmissionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case - successfully-verify-survey-submission",
			args: args{
				ctx: ctx,
				input: dto.VerifySurveySubmissionInput{
					ProjectID:   10000000000000,
					FormID:      uuid.New().String(),
					SubmitterID: 10,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - fail to verify survey submission",
			args: args{
				ctx: ctx,
				input: dto.VerifySurveySubmissionInput{
					ProjectID:   10000000000000,
					FormID:      uuid.New().String(),
					SubmitterID: 10,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - unable to update survey",
			args: args{
				ctx: ctx,
				input: dto.VerifySurveySubmissionInput{
					ProjectID:   10000000000000,
					FormID:      uuid.New().String(),
					SubmitterID: 100000000000,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - fail to verify survey submission" {
				fakeSurveys.MockListSubmittersFn = func(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
					return nil, fmt.Errorf("failed to get submitters")
				}
			}
			if tt.name == "Sad case - unable to update survey" {
				fakeSurveys.MockListSubmittersFn = func(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
					return []domain.Submitter{
						{
							ID:          100000000000,
							Type:        gofakeit.BeerName(),
							DisplayName: gofakeit.BeerName(),
							CreatedAt:   time.Now(),
							UpdatedAt:   time.Time{},
							DeletedAt:   time.Time{},
						},
					}, nil
				}
				fakeDB.MockUpdateUserSurveysFn = func(ctx context.Context, survey *domain.UserSurvey, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update survey")
				}
			}
			got, err := u.VerifySurveySubmission(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseSurveysImpl.VerifySurveySubmission() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseSurveysImpl.VerifySurveySubmission() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseSurveysImpl_ListSurveyRespondents(t *testing.T) {
	ctx := context.Background()

	fakeSurveys := mockSurveys.NewSurveysMock()
	fakeDB := pgMock.NewPostgresMock()
	fakeNotification := mockNotification.NewServiceNotificationMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeDB, fakeNotification)
	type args struct {
		ctx             context.Context
		projectID       int
		paginationInput dto.PaginationsInput
		formID          string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.SurveyRespondentPage
		wantErr bool
	}{
		{
			name: "happy case - list survey respondents",
			args: args{
				ctx:       ctx,
				projectID: 10000000000000,
				formID:    uuid.New().String(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case - unable to list survey respondents",
			args: args{
				ctx:       ctx,
				projectID: 10000000000000,
				formID:    uuid.New().String(),
				paginationInput: dto.PaginationsInput{
					Limit:       10,
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - unable to list survey respondents" {
				fakeDB.MockListSurveyRespondentsFn = func(ctx context.Context, projectID int, formID string, pagination *domain.Pagination) ([]*domain.SurveyRespondent, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("failed to list survey respondents")
				}
			}
			_, err := u.ListSurveyRespondents(tt.args.ctx, tt.args.projectID, tt.args.formID, tt.args.paginationInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseSurveysImpl.ListSurveyRespondents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
