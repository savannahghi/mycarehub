package client_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/mock"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

func TestUseCasesClientImpl_RegisterClient(t *testing.T) {
	ctx := context.Background()
	d := testFakeInfrastructureInteractor

	userPayload := &dto.UserInput{
		FirstName:   "FirstName",
		LastName:    "Last Name",
		Username:    "User Name",
		MiddleName:  "Middle Name",
		DisplayName: "Display Name",
		Gender:      enumutils.GenderMale,
	}

	clientPayload := &dto.ClientProfileInput{
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
		// {
		// 	name: "Happy Case",
		// 	args: args{
		// 		ctx:         ctx,
		// 		userInput:   userPayload,
		// 		clientInput: clientPayload,
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "Sad Case: Fail to register user",
			args: args{
				ctx:         ctx,
				userInput:   userPayload,
				clientInput: clientPayload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = mock.NewCreateMock()

			if tt.name == "Sad Case: Fail to register user" {
				fakeCreate.RegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to register a client")
				}
			}

			got, err := d.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClientImpl.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", nil)
			}
		})
	}
}

func TestUseCasesClientImpl_TransferClient_Unittest(t *testing.T) {
	ctx := context.Background()
	f := testFakeInfrastructureInteractor

	ID := uuid.New().String()

	type args struct {
		ctx                   context.Context
		clientID              string
		originFacilityID      string
		destinationFacilityID string
		reason                enums.TransferReason
		notes                 string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                   ctx,
				clientID:              ksuid.New().String(),
				originFacilityID:      ksuid.New().String(),
				destinationFacilityID: ksuid.New().String(),
				reason:                enums.RelocationTransferReason,
				notes:                 ksuid.New().String(),
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:                   ctx,
				clientID:              "",
				originFacilityID:      ksuid.New().String(),
				destinationFacilityID: ksuid.New().String(),
				reason:                enums.OtherTransferReason,
				notes:                 "other-reason",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeCreate.RegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return &domain.ClientUserProfile{
						User: &domain.User{
							ID:                  &ID,
							Username:            "User Name",
							DisplayName:         "Display Name",
							FirstName:           "FirstName",
							MiddleName:          userInput.MiddleName,
							LastName:            "Last Name",
							UserType:            enums.HealthcareWorkerUser,
							Gender:              enumutils.GenderMale,
							Active:              false,
							Contacts:            []*domain.Contact{},
							Languages:           []enumutils.Language{},
							LastSuccessfulLogin: &time.Time{},
							LastFailedLogin:     &time.Time{},
							FailedLoginCount:    "0",
							NextAllowedLogin:    &time.Time{},
							TermsAccepted:       false,
							AcceptedTermsID:     ksuid.New().String(),
							Flavour:             feedlib.FlavourConsumer,
						},
						Client: &domain.ClientProfile{
							ClientType: enums.ClientTypeOvc,
						},
					}, nil
				}

				fakeQuery.GetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {

					healthRecordID := ksuid.New().String()
					return &domain.ClientProfile{
						ID:                      &ID,
						UserID:                  new(string),
						TreatmentEnrollmentDate: &time.Time{},
						ClientType:              enums.ClientTypeOvc,
						Active:                  false,
						HealthRecordID:          &healthRecordID,
						Identifiers:             []*domain.Identifier{},
						Addresses:               []*domain.Addresses{},
						RelatedPersons:          []*domain.RelatedPerson{},
						FacilityID:              ksuid.New().String(),
						TreatmentBuddyUserID:    ksuid.New().String(),
						CHVUserID:               ksuid.New().String(),
						ClientCounselled:        false,
					}, nil
				}

				fakeUpdate.TransferClientFn = func(ctx context.Context, clientID, originFacilityID, destinationFacilityID string, reason enums.TransferReason, notes string) (bool, error) {
					return true, nil
				}
			}

			if tt.name == "Sad case" {
				fakeCreate.RegisterClientFn = func(ctx context.Context, userInput *dto.UserInput, clientInput *dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to register a client")
				}

				fakeUpdate.TransferClientFn = func(ctx context.Context, clientID, originFacilityID, destinationFacilityID string, reason enums.TransferReason, notes string) (bool, error) {
					return false, fmt.Errorf("failed to register a client")
				}
			}

			got, err := f.TransferClient(tt.args.ctx, tt.args.clientID, tt.args.originFacilityID, tt.args.destinationFacilityID, tt.args.reason, tt.args.notes)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClientImpl.TransferClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
