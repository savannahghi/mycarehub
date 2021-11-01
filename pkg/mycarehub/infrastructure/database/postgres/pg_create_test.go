package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestOnboardingDb_CreateFacility(t *testing.T) {
	ctx := context.Background()
	name := "Kanairo One"
	code := "KN001"
	county := "Kanairo"
	description := "This is just for mocking"
	type args struct {
		ctx      context.Context
		facility *dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
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
				facility: &dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)
			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if tt.name == "sad case - facility code not defined" {
				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return facility, nil
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facility to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facility not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_RegisterClient(t *testing.T) {
	ctx := context.Background()
	userInput := &dto.UserInput{
		FirstName:   "John",
		LastName:    "Joe",
		Username:    "Jontez",
		MiddleName:  "Johnny",
		DisplayName: "jo",
		Gender:      enumutils.GenderMale,
	}

	clientInput := dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}
	type args struct {
		ctx         context.Context
		userInput   *dto.UserInput
		clientInput *dto.ClientProfileInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: &clientInput,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to create user",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: &clientInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case: Fail to create user" {
				fakeGorm.RegisterClientFn = func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to create a client user")
				}
			}

			got, err := d.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got :%v", got)
			}
		})
	}
}
