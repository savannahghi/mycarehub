package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

func setup() (usecases.SignUpUseCases, error) {
	fr, err := database.NewFirebaseRepository(context.Background())
	if err != nil {
		return nil, fmt.Errorf("can't instantiate firebase repository in resolver: %w", err)
	}

	profile := usecases.NewProfileUseCase(fr)
	otp := otp.NewOTPService(fr)
	erp := erp.NewERPService(fr)
	chrg := chargemaster.NewChargeMasterUseCasesImpl(fr)
	engage := engagement.NewServiceEngagementImpl(fr)
	mg := mailgun.NewServiceMailgunImpl()
	mes := messaging.NewServiceMessagingImpl()
	supplier := usecases.NewSupplierUseCases(fr, profile, erp, chrg, engage, mg, mes)
	userpin := usecases.NewUserPinUseCase(fr, otp, profile)
	su := usecases.NewSignUpUseCases(fr, profile, userpin, supplier)

	return su, nil
}

func TestCheckPhoneExists(t *testing.T) {
	signup, err := setup()
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	type testArgs struct {
		name        string
		phone       string
		wantErr     bool
		expectedErr string
	}

	testbed := []testArgs{
		{
			name:    "valid : phone number does not exist",
			phone:   base.TestUserPhoneNumber,
			wantErr: false,
		},

		{
			name:    "valid : phone number already exists",
			phone:   "+254718123098", // use a different number since tear down has not happened yet
			wantErr: true,
		},

		{
			name:        "invalid : wrong phone number format",
			phone:       "71812308",
			wantErr:     true,
			expectedErr: "failed to create firebase user: phone number must be a valid",
		},
	}

	for _, tt := range testbed {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				// signup user with the phone number then run phone number check
				resp, err := signup.CreateUserByPhone(context.Background(), tt.phone, "1234", base.FlavourConsumer)
				if tt.expectedErr == "" {
					assert.Nil(t, err)
					assert.NotNil(t, resp)
				}

				if tt.expectedErr != "" {
					assert.NotNil(t, err)
					assert.Contains(t, err.Error(), tt.expectedErr)

					resp2, err2 := signup.CheckPhoneExists(context.Background(), tt.phone)
					assert.Nil(t, err2)
					assert.NotNil(t, resp2)
					assert.Equal(t, false, resp2)
				}

			}

			//todo:(dexter) restore this
			// if !tt.wantErr {

			// 	resp, err := signup.CheckPhoneExists(context.Background(), tt.phone)
			// 	assert.Nil(t, err)
			// 	assert.Equal(t, true, resp)

			// 	// signup user with the phone number then run another phone number check
			// 	resp2, err := signup.CreateUserByPhone(context.Background(), tt.phone, "1234", base.FlavourConsumer)
			// 	assert.Nil(t, err)
			// 	assert.NotNil(t, resp2)

			// 	// now check the phone number that has been used above exists.
			// 	resp3, err := signup.CheckPhoneExists(context.Background(), tt.phone)
			// 	assert.Nil(t, err)
			// 	assert.NotNil(t, resp3)
			// 	assert.Equal(t, true, resp3)
			// }

		})
	}

}

func TestCreateUserWithPhoneNumber(t *testing.T) {
	signup, err := setup()
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	type testArgs struct {
		name        string
		phone       string
		pin         string
		flavour     base.Flavour
		wantErr     bool
		expectedErr string
	}

	testbed := []testArgs{
		{
			name:        "valid : should consumer create user",
			phone:       base.TestUserPhoneNumber,
			pin:         "1234",
			flavour:     base.FlavourConsumer,
			wantErr:     false,
			expectedErr: "",
		},
	}

	for _, tt := range testbed {

		t.Run(tt.name, func(t *testing.T) {

			resp, err := signup.CreateUserByPhone(context.Background(), tt.phone, tt.pin, tt.flavour)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
			//todo:(dexter) restore this
			// if !tt.wantErr {
			// 	assert.Nil(t, err)
			// 	assert.NotNil(t, resp)
			// 	assert.NotNil(t, resp.Profile)
			// 	assert.NotNil(t, resp.CustomerProfile)
			// 	assert.NotNil(t, resp.SupplierProfile)
			// }
		})

	}
}
