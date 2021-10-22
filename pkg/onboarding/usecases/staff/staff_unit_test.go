package staff_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

func TestOnboardingDb_RegisterStaffUser(t *testing.T) {
	ctx := context.Background()

	testFacilityID := uuid.New().String()
	testUserID := uuid.New().String()
	testTime := time.Now()
	testID := uuid.New().String()

	d := testFakeInfrastructureInteractor

	userInput := &dto.UserInput{
		UserName:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
	}

	staffInput := &dto.StaffProfileInput{
		StaffNumber:       "s123",
		DefaultFacilityID: &testFacilityID,
	}

	staffNoFacilityInput := &dto.StaffProfileInput{
		StaffNumber: "s123",
	}

	type args struct {
		ctx   context.Context
		user  *dto.UserInput
		staff *dto.StaffProfileInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.StaffUserProfile
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:   ctx,
				user:  userInput,
				staff: staffInput,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing facility",
			args: args{
				ctx:   ctx,
				user:  userInput,
				staff: staffNoFacilityInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy case" {

				fakeCreate.GetOrCreateFacilityFn = func(ctx context.Context, facility dto.FacilityInput) (*domain.Facility, error) {
					return &domain.Facility{
						ID:          &testFacilityID,
						Name:        "test",
						Code:        "f1234",
						Active:      true,
						County:      "test",
						Description: "test description",
					}, nil
				}
				fakeCreate.RegisterStaffUserFn = func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
					return &domain.StaffUserProfile{
						User: &domain.User{
							ID:                  &testUserID,
							Username:            "test",
							DisplayName:         "test",
							FirstName:           "test",
							MiddleName:          "test",
							LastName:            "test",
							Active:              true,
							LastSuccessfulLogin: &testTime,
							LastFailedLogin:     &testTime,
							NextAllowedLogin:    &testTime,
							FailedLoginCount:    "0",
							TermsAccepted:       true,
							AcceptedTermsID:     testID,
						},
						Staff: &domain.StaffProfile{
							ID:                &testID,
							UserID:            &testUserID,
							StaffNumber:       "s123",
							DefaultFacilityID: &testFacilityID,
						},
					}, nil
				}
			}

			if tt.name == "invalid: missing facility" {
				fakeCreate.RegisterStaffUserFn = func(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
					return nil, fmt.Errorf("test error")
				}
			}

			_, err := d.RegisterStaffUser(tt.args.ctx, tt.args.user, tt.args.staff)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RegisterStaffUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
