package staff_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
)

func TestUseCaseStaffProfileImpl_RegisterStaffUser(t *testing.T) {
	f := testInfrastructureInteractor
	ctx := context.Background()

	testFacilityID := uuid.New().String()

	facilityInput := dto.FacilityInput{
		Name:        "test",
		Code:        "c123",
		Active:      true,
		County:      "test",
		Description: "test description",
	}

	// Create a facility
	facility, err := f.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("Failed to create facility: %v", err)
	}

	userInput := &dto.UserInput{
		UserName:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
	}

	staffInput := &dto.StaffProfileInput{
		StaffNumber:       "s123",
		DefaultFacilityID: facility.ID,
	}

	staffInputNoFacility := &dto.StaffProfileInput{
		StaffNumber:       "s123",
		DefaultFacilityID: &testFacilityID,
	}

	type args struct {
		ctx   context.Context
		user  *dto.UserInput
		staff *dto.StaffProfileInput
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
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
			wantNil: false,
		},
		{
			name: "invalid: missing facility",
			args: args{
				ctx:   ctx,
				user:  userInput,
				staff: staffInputNoFacility,
			},
			wantErr: true,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := f.RegisterStaffUser(tt.args.ctx, tt.args.user, tt.args.staff)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseStaffUserImpl.RegisterStaffUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected staff profile to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected staff profile not to be nil for %v", tt.name)
				return
			}
		})
	}

}
