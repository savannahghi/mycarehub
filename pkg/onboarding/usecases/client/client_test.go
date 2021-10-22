package client_test

import (
	"context"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
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
