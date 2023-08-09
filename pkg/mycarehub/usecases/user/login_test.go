package user

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	clinicalMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical/mock"
	matrixMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	smsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/sms/mock"
	twilioMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/twilio/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	"gorm.io/gorm"
)

func TestUseCasesUserImpl_caregiverProfileCheck(t *testing.T) {
	pin := "1234"
	id := gofakeit.UUID()

	type args struct {
		ctx         context.Context
		credentials *dto.LoginInput
		response    dto.ILoginResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "happy case: user type caregiver",
			args: args{
				ctx: context.Background(),
				credentials: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      pin,
					Flavour:  feedlib.FlavourConsumer,
				},
				response: &dto.LoginResponse{
					Response: &dto.Response{
						User: &dto.User{
							ID:       id,
							Username: gofakeit.Username(),
							Name:     gofakeit.Name(),
						},
					},
				},
			},
			want: true,
		},
		{
			name: "sad case: fail to check caregiver profile",
			args: args{
				ctx: context.Background(),
				credentials: &dto.LoginInput{
					Username: gofakeit.Username(),
					PIN:      pin,
					Flavour:  feedlib.FlavourConsumer,
				},
				response: &dto.LoginResponse{
					Response: &dto.Response{
						User: &dto.User{
							ID:       id,
							Username: gofakeit.Username(),
							Name:     gofakeit.Name(),
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "sad case: fail to check caregiver profile" {
				fakeDB.MockCheckCaregiverExistsFn = func(ctx context.Context, userID string) (bool, error) {
					return false, fmt.Errorf("caregiver not found: %w", gorm.ErrRecordNotFound)
				}
			}

			if got := us.caregiverProfileCheck(tt.args.ctx, tt.args.credentials, tt.args.response); got != tt.want {
				t.Errorf("UseCasesUserImpl.caregiverProfileCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesUserImpl_pinResetRequestCheck(t *testing.T) {
	type args struct {
		ctx         context.Context
		credentials *dto.LoginInput
		response    dto.ILoginResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Happy case: pin reset request check",
			args: args{
				ctx: context.Background(),
				credentials: &dto.LoginInput{
					Username: gofakeit.BeerName(),
					PIN:      "1234",
					Flavour:  feedlib.FlavourConsumer,
				},
				response: &dto.LoginResponse{
					Response: &dto.Response{
						User: &dto.User{
							ID: gofakeit.UUID(),
						},
						AuthCredentials: dto.AuthCredentials{
							RefreshToken: gofakeit.UUID(),
							IDToken:      gofakeit.UUID(),
							ExpiresIn:    "",
						},
						GetStreamToken: "",
					},
					IsCaregiver: false,
					IsClient:    true,
				},
			},
			want: false,
		},
		{
			name: "Sad case: unable to get user profile by username",
			args: args{
				ctx: context.Background(),
				credentials: &dto.LoginInput{
					Username: gofakeit.BeerName(),
					PIN:      "1234",
					Flavour:  feedlib.FlavourConsumer,
				},
				response: &dto.LoginResponse{
					Response: &dto.Response{
						User: &dto.User{
							ID: gofakeit.UUID(),
						},
						AuthCredentials: dto.AuthCredentials{
							RefreshToken: gofakeit.UUID(),
							IDToken:      gofakeit.UUID(),
							ExpiresIn:    "",
						},
						GetStreamToken: "",
					},
					IsCaregiver: false,
					IsClient:    true,
				},
			},
			want: false,
		},
		{
			name: "Sad case: unable to get client profile",
			args: args{
				ctx: context.Background(),
				credentials: &dto.LoginInput{
					Username: gofakeit.BeerName(),
					PIN:      "1234",
					Flavour:  feedlib.FlavourConsumer,
				},
				response: &dto.LoginResponse{
					Response: &dto.Response{
						User: &dto.User{
							ID: gofakeit.UUID(),
						},
						AuthCredentials: dto.AuthCredentials{
							RefreshToken: gofakeit.UUID(),
							IDToken:      gofakeit.UUID(),
							ExpiresIn:    "",
						},
						GetStreamToken: "",
					},
					IsCaregiver: false,
					IsClient:    true,
				},
			},
			want: false,
		},
		{
			name: "Sad case: unable to get client service request",
			args: args{
				ctx: context.Background(),
				credentials: &dto.LoginInput{
					Username: gofakeit.BeerName(),
					PIN:      "1234",
					Flavour:  feedlib.FlavourConsumer,
				},
				response: &dto.LoginResponse{
					Response: &dto.Response{
						User: &dto.User{
							ID: gofakeit.UUID(),
						},
						AuthCredentials: dto.AuthCredentials{
							RefreshToken: gofakeit.UUID(),
							IDToken:      gofakeit.UUID(),
							ExpiresIn:    "",
						},
						GetStreamToken: "",
					},
					IsCaregiver: false,
					IsClient:    true,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeOTP := otpMock.NewOTPUseCaseMock()
			fakeAuthority := authorityMock.NewAuthorityUseCaseMock()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeClinical := clinicalMock.NewClinicalServiceMock()
			fakeSMS := smsMock.NewSMSServiceMock()
			fakeTwilio := twilioMock.NewTwilioServiceMock()
			fakeMatrix := matrixMock.NewMatrixMock()

			us := NewUseCasesUserImpl(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension, fakeOTP, fakeAuthority, fakePubsub, fakeClinical, fakeSMS, fakeTwilio, fakeMatrix)

			if tt.name == "Sad case: unable to get user profile by username" {
				fakeDB.MockGetUserProfileByUsernameFn = func(ctx context.Context, username string) (*domain.User, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to get client profile" {
				fakeDB.MockGetClientProfileFn = func(ctx context.Context, userID, programID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: unable to get client service request" {
				fakeDB.MockGetClientServiceRequestsFn = func(ctx context.Context, requestType, status, clientID, facilityID string) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if got := us.pinResetRequestCheck(tt.args.ctx, tt.args.credentials, tt.args.response); got != tt.want {
				t.Errorf("UseCasesUserImpl.pinResetRequestCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}
