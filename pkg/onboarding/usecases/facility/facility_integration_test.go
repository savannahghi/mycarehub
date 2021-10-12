package facility_test

import (
	"context"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
)

func TestUseCaseFacilityImpl_DeleteFacility_Integration_Test(t *testing.T) {
	ctx := context.Background()

	i := testInfrastructureInteractor

	//Create facility
	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        "KN001",
		County:      "Kanairo",
		Description: "This is just for integration testing",
	}

	// create a facility
	facility, err := i.CreateFacility(ctx, *facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflcode := facility.Code

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				id:  mflcode,
			},
			wantErr: false,
		},

		{
			name: "Sad case: Empty id",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},

		{
			name: "Sad case: invalid id",
			args: args{
				ctx: ctx,
				id:  "mflcode",
			},
			//should not fail: assume the id has already been deleted
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := i.DeleteFacility(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.DeleteFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
