package surveys

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	mockSurveys "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys/mock"
)

func TestUsecaseSurveysImpl_ListSurveys(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeSurveys := mockSurveys.NewSurveysMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB)
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
	fakeDB := pgMock.NewPostgresMock()
	fakeSurveys := mockSurveys.NewSurveysMock()
	u := NewUsecaseSurveys(fakeSurveys, fakeDB)

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
				fakeDB.MockGetUserSurveyFormsFn = func(ctx context.Context, userID string) ([]*domain.UserSurveys, error) {
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
