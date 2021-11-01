package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestMyCareHubDb_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	name := gofakeit.Name()
	code := "KN001"
	county := enums.CountyTypeNairobi
	description := gofakeit.HipsterSentence(15)

	facility := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	invalidFacility := &dto.FacilityInput{
		Name:        name,
		Active:      true,
		County:      county,
		Description: description,
	}

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
				ctx:      ctx,
				facility: facility,
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx:      ctx,
				facility: invalidFacility,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm)
			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if tt.name == "sad case - facility code not defined" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - nil facility input" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_RegisterClient(t *testing.T) {
	ctx := context.Background()
	userInput := &dto.UserInput{
		FirstName:   gofakeit.FirstName(),
		LastName:    gofakeit.LastName(),
		Username:    gofakeit.Username(),
		MiddleName:  gofakeit.Name(),
		DisplayName: gofakeit.BeerAlcohol(),
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
			name: "Sad Case: Fail to create user with nil user input",
			args: args{
				ctx:         ctx,
				userInput:   nil,
				clientInput: &clientInput,
			},
			wantErr: true,
		},

		{
			name: "Sad Case: Fail to create user with nil client input",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: nil,
			},
			wantErr: true,
		},

		{
			name: "Sad Case: Fail to register client",
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
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case: Fail to create user with nil user input" {
				fakeGorm.MockRegisterClientFn = func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to create a client user")
				}
			}

			if tt.name == "Sad Case: Fail to create user with nil client input" {
				fakeGorm.MockRegisterClientFn = func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to create a client user")
				}
			}

			if tt.name == "Sad Case: Fail to register client" {
				fakeGorm.MockRegisterClientFn = func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to create a client user")
				}
			}

			got, err := d.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got :%v", got)
			}
		})
	}
}

func TestOnboardingDb_SavePin(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		pinInput *domain.UserPIN
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save pin",
			args: args{
				ctx: ctx,
				pinInput: &domain.UserPIN{
					UserID:    "123456",
					HashedPIN: "12345",
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					Flavour:   feedlib.FlavourConsumer,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save pin",
			args: args{
				ctx: ctx,
				pinInput: &domain.UserPIN{
					UserID:    "123456",
					HashedPIN: "12345",
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					Flavour:   feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save pin" {
				fakeGorm.MockSavePinFn = func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
					return false, fmt.Errorf("failed to save pin")
				}
			}

			got, err := d.SavePin(tt.args.ctx, tt.args.pinInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.SavePin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OnboardingDb.SavePin() = %v, want %v", got, tt.want)
			}
		})
	}
}
