package client_test

import (
	"context"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/segmentio/ksuid"
	"github.com/tj/assert"
)

func TestUseCasesClientImplIntegration_RegisterClient(t *testing.T) {
	ctx := context.Background()
	f := testInfrastructureInteractor

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

	// Create an empty client, it should returns an error
	client, err := f.RegisterClient(ctx, nil, nil)
	assert.NotNil(t, err)
	assert.Nil(t, client)

	// Create a client with valid payload
	client1, err := f.RegisterClient(ctx, userPayload, clientPayload)
	assert.Nil(t, err)
	assert.NotNil(t, client1)

	// TODO: Try creating the same user twice, should throw an error after we check for uniqueness
}

func TestUseCasesClientImpl_AddIdentifier(t *testing.T) {
	ctx := context.Background()
	f := testInfrastructureInteractor

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
	client, err := f.RegisterClient(ctx, userPayload, clientPayload)
	if err != nil {
		t.Errorf("failed to create client: %v", err)
		return
	}

	type args struct {
		ctx       context.Context
		clientID  string
		idType    enums.IdentifierType
		idValue   string
		isPrimary bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully add identifier",
			args: args{
				ctx:       ctx,
				clientID:  *client.Client.ID,
				idType:    enums.IdentifierTypeCCC,
				idValue:   "12345",
				isPrimary: true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to add identifier",
			args: args{
				ctx:      ctx,
				clientID: "non-existent",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := f.AddIdentifier(tt.args.ctx, tt.args.clientID, tt.args.idType, tt.args.idValue, tt.args.isPrimary)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesClientImpl.AddIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestUseCasesClientImpl_TransferClient_Integration(t *testing.T) {
	ctx := context.Background()

	m := testInfrastructureInteractor
	i := testInteractor

	originalFacilityID := ksuid.New().String()
	destinationFacilityID := ksuid.New().String()
	reason := enums.RelocationTransferReason
	notes := "Other"

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
	client, err := m.RegisterClient(ctx, userPayload, clientPayload)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	clientProfile, err1 := m.GetClientProfileByClientID(ctx, *client.Client.ID)
	assert.Nil(t, err1)
	assert.NotNil(t, clientProfile)

	bool, err2 := i.ClientUseCase.TransferClient(ctx, *clientProfile.ID, originalFacilityID, destinationFacilityID, reason, notes)
	assert.Nil(t, err2)
	assert.Equal(t, true, bool)

	bool1, err3 := i.ClientUseCase.TransferClient(ctx, "", originalFacilityID, destinationFacilityID, reason, notes)
	assert.NotNil(t, err3)
	assert.Equal(t, false, bool1)

}
