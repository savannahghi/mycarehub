package surveys

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	mockSurveys "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys/mock"
	mockNotification "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
)

func TestUsecaseSurveysImpl_ListSurveys(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeSurveys := mockSurveys.NewSurveysMock()
	fakeNotification := mockNotification.NewServiceNotificationMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeNotification)
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
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeNotification)

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
				fakeDB.MockGetUserSurveyFormsFn = func(ctx context.Context, userID string) ([]*domain.UserSurvey, error) {
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
	fakeSurveys := mockSurveys.NewSurveysMock()
	fakeDB := pgMock.NewPostgresMock()
	fakeNotification := mockNotification.NewServiceNotificationMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB, fakeDB, fakeNotification)

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
			name: "sad case: failed to generate survey access link",
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
			want:    false, // only report the error to sentry
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: failed to get filtered clients" {
				fakeDB.MockGetClientsByFilterParamsFn = func(ctx context.Context, facilityID *string, filterParams *dto.ClientFilterParamsInput) ([]*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get filtered clients")
				}

			}
			if tt.name == "sad case: failed to get survey form" {
				fakeSurveys.MockGetSurveyFormFn = func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
					return nil, fmt.Errorf("failed to get survey form")
				}
			}
			if tt.name == "sad case: failed to generate survey access link" {
				fakeSurveys.MockGeneratePublickAccessLinkFn = func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
					return nil, fmt.Errorf("failed to generate survey access link")
				}
			}

			if tt.name == "sad case: failed to create survey" {
				fakeDB.MockCreateUserSurveyFn = func(ctx context.Context, userSurvey []*dto.UserSurveyInput) error {
					return fmt.Errorf("failed to create survey")
				}
			}

			if tt.name == "sad case: failed to notify user" {
				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify user")
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
