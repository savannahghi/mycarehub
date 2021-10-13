package facility_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/segmentio/ksuid"
)

func TestUseCaseFacilityImpl_CreateFacility(t *testing.T) {
	f := testInfrastructureInteractor
	ctx := context.Background()
	name := "Kanairo One"
	code := ksuid.New().String()
	county := "Kanairo"
	description := "This is just for mocking"

	type args struct {
		ctx      context.Context
		facility dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: dto.FacilityInput{
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx: ctx,
				facility: dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := f.CreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.CreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantNil && got == nil {
				t.Errorf("UseCaseFacilityImpl.CreateFacility() expected to return a value, got:  %v", got)
			}
		})
	}
	// TODO: add teardown
}

func TestUseCaseFacilityImpl_RetrieveFacility(t *testing.T) {
	f := testInfrastructureInteractor

	ctx := context.Background()

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        ksuid.New().String(),
		County:      "Kanairo",
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := f.CreateFacility(ctx, *facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	id := facility.ID

	invalidID := int64(-100)

	type args struct {
		ctx context.Context
		id  *int64
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid ID passed",
			args: args{
				ctx: ctx,
				id:  &id,
			},
			wantErr: false,
			want:    facility,
		},
		{
			name: "sad case - no ID passed",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "sad case - invalid ID",
			args: args{
				ctx: ctx,
				id:  &invalidID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.RetrieveFacility(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UseCaseFacilityImpl.RetrieveFacility() = %v, want %v", got, tt.want)
			}
		})
	}
	// TODO: add teardown
}
