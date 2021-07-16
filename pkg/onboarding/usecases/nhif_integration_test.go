package usecases_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestAddNHIFDetails(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	validInput := dto.NHIFDetailsInput{
		MembershipNumber:          "123456",
		Employment:                domain.EmploymentTypeEmployed,
		NHIFCardPhotoID:           uuid.New().String(),
		IDDocType:                 enumutils.IdentificationDocTypeMilitary,
		IdentificationCardPhotoID: uuid.New().String(),
		IDNumber:                  "11111111",
	}
	type args struct {
		ctx   context.Context
		input dto.NHIFDetailsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:) successfully add NHIF Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
		{
			name: "sad:( unsuccessfully add NHIF Details",
			args: args{
				ctx:   context.Background(),
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( unsuccessfully add NHIF details since it exists",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nhif, err := s.NHIF.AddNHIFDetails(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NHIFUseCaseImpl.AddNHIFDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && nhif != nil {
				t.Errorf("the error was not expected")
				return
			}

			if !tt.wantErr && nhif == nil {
				t.Errorf("an error was expected: %v", err)
				return
			}
		})
	}
}

func TestNHIFDetails(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	_, err = s.NHIF.AddNHIFDetails(
		ctx,
		dto.NHIFDetailsInput{
			MembershipNumber:          fmt.Sprintln(time.Now().Unix()),
			Employment:                domain.EmploymentTypeEmployed,
			NHIFCardPhotoID:           uuid.New().String(),
			IDDocType:                 enumutils.IdentificationDocTypeMilitary,
			IdentificationCardPhotoID: uuid.New().String(),
			IDNumber:                  fmt.Sprintln(time.Now().Unix()),
		},
	)
	if err != nil {
		t.Errorf("expected NHIF: %v", err)
		return
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.NHIFDetails
		wantErr bool
	}{
		{
			name: "happy:) successfully get NHIF Details",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "sad:( unsuccessfully get NHIF Details",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nhif, err := s.NHIF.NHIFDetails(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NHIFUseCaseImpl.NHIFDetails() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if tt.wantErr && nhif != nil {
				t.Errorf("the error was not expected")
				return
			}

			if !tt.wantErr && nhif == nil {
				t.Errorf("an error was expected: %v", err)
				return
			}
		})
	}
}
