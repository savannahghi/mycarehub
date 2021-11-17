package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
)

func TestMyCareHubDb_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	name := gofakeit.Name()
	code := gofakeit.Number(300, 400)
	county := "Nairobi"
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
		{
			name: "Sad Case - Fail to get ot create facility",
			args: args{
				ctx:      ctx,
				facility: facility,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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

			if tt.name == "Sad Case - Fail to get ot create facility" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get or create facility")
				}
			}

			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
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

func TestMyCareHubDb_SaveTemporaryUserPin(t *testing.T) {
	ctx := context.Background()
	ID := ksuid.New().String()
	flavor := feedlib.FlavourConsumer

	newExtension := extension.NewExternalMethodsImpl()

	tempPin, err := newExtension.GenerateTempPIN(ctx)
	if err != nil {
		t.Errorf("failed to generate temporary pin: %v", err)
	}
	salt, encryptedTempPin := newExtension.EncryptPIN(tempPin, nil)

	pinPayload := &domain.UserPIN{
		UserID:    ID,
		HashedPIN: encryptedTempPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   flavor,
		Salt:      salt,
	}
	type args struct {
		ctx     context.Context
		pinData *domain.UserPIN
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     ctx,
				pinData: pinPayload,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: invalid user id provided",
			args: args{
				ctx: ctx,
				pinData: &domain.UserPIN{
					UserID:    gofakeit.Sentence(200),
					HashedPIN: encryptedTempPin,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					IsValid:   true,
					Flavour:   flavor,
					Salt:      salt,
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad Case - Fail to save temporary pin",
			args: args{
				ctx:     ctx,
				pinData: pinPayload,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "invalid: missing userID" {
				fakeGorm.MockSaveTemporaryUserPinFn = func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
					return false, fmt.Errorf("user id must be provided")
				}
			}

			if tt.name == "Sad Case - Fail to save temporary pin" {
				fakeGorm.MockSaveTemporaryUserPinFn = func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
					return false, fmt.Errorf("fail to save temporary pin")
				}
			}

			got, err := d.SaveTemporaryUserPin(tt.args.ctx, tt.args.pinData)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveTemporaryUserPin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.SaveTemporaryUserPin() = %v, want %v", got, tt.want)
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
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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

func TestMyCareHubDb_SaveOTP(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		otpInput *domain.OTP
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save otp",
			args: args{
				ctx: ctx,
				otpInput: &domain.OTP{
					UserID:      "12345",
					Valid:       true,
					ValidUntil:  time.Now().Add(time.Hour * 1),
					GeneratedAt: time.Now(),
					PhoneNumber: gofakeit.Phone(),
					Channel:     "SMS",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save otp",
			args: args{
				ctx: ctx,
				otpInput: &domain.OTP{
					UserID:      "12345",
					Valid:       true,
					ValidUntil:  time.Now().Add(time.Hour * 1),
					GeneratedAt: time.Now(),
					PhoneNumber: gofakeit.Phone(),
					Channel:     "SMS",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save otp" {
				fakeGorm.MockSaveOTPFn = func(ctx context.Context, otpInput *gorm.UserOTP) error {
					return fmt.Errorf("failed to save otp")
				}
			}

			if err := d.SaveOTP(tt.args.ctx, tt.args.otpInput); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_SaveSecurityQuestionResponse(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                      context.Context
		securityQuestionResponse *dto.SecurityQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save question response",
			args: args{
				ctx: ctx,
				securityQuestionResponse: &dto.SecurityQuestionResponseInput{
					UserID:             uuid.New().String(),
					SecurityQuestionID: uuid.New().String(),
					Response:           "A valid response",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save question response",
			args: args{
				ctx: ctx,
				securityQuestionResponse: &dto.SecurityQuestionResponseInput{
					UserID:             uuid.New().String(),
					SecurityQuestionID: uuid.New().String(),
					Response:           "A valid response",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save question response" {
				fakeGorm.MockSaveSecurityQuestionResponseFn = func(ctx context.Context, securityQuestionResponse *gorm.SecurityQuestionResponse) error {
					return fmt.Errorf("failed to save security question response")
				}
			}
			if err := d.SaveSecurityQuestionResponse(tt.args.ctx, tt.args.securityQuestionResponse); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
