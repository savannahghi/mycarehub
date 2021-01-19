package database_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
)

func TestRemoveKYCProcessingRequest(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	assert.Nil(t, err)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), base.TestUserPhoneNumber)

	ctx, auth, err := GetTestAuthenticatedContext(t)
	assert.Nil(t, err)
	assert.NotNil(t, auth)

	fr, err := database.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	// assemble kyc profile input
	input1 := domain.OrganizationNutrition{
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		KRAPIN:                             "someKRAPIN",
		KRAPINUploadID:                     "KRAPINUploadID",
		SupportingDocumentsUploadID:        []string{"SupportingDocumentsUploadID", "Support"},
		CertificateOfIncorporation:         "CertificateOfIncorporation",
		CertificateOfInCorporationUploadID: "CertificateOfInCorporationUploadID",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeMilitary,
				IdentificationDocNumber:         "IdentificationDocNumber",
				IdentificationDocNumberUploadID: "IdentificationDocNumberUploadID",
			},
		},
		OrganizationCertificate: "OrganizationCertificate",
		RegistrationNumber:      "RegistrationNumber",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
	}

	kycJSON, err := json.Marshal(input1)
	assert.Nil(t, err)

	var kycAsMap map[string]interface{}
	err = json.Unmarshal(kycJSON, &kycAsMap)
	assert.Nil(t, err)

	// get the user profile
	profile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	assert.Nil(t, err)
	assert.NotNil(t, profile)

	// fetch the supplier profile
	sup, err := fr.GetSupplierProfileByProfileID(ctx, profile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, sup)

	//call remove kyc process request. this should fail since the user has not added a kyc yet
	err = fr.RemoveKYCProcessingRequest(ctx, sup.ID)
	assert.NotNil(t, err)

	sup.SupplierKYC = kycAsMap

	// now add the kyc processing request
	req1 := &domain.KYCRequest{
		ID:                  uuid.New().String(),
		ReqPartnerType:      sup.PartnerType,
		ReqOrganizationType: domain.OrganizationType(sup.AccountType),
		ReqRaw:              sup.SupplierKYC,
		Processed:           false,
		SupplierRecord:      sup,
		Status:              domain.KYCProcessStatusPending,
	}
	err = fr.StageKYCProcessingRequest(ctx, req1)
	assert.Nil(t, err)

	// call remove kypc processing request again. this should pass now since there is and existing processing request added
	err = fr.RemoveKYCProcessingRequest(ctx, sup.ID)
	assert.Nil(t, err)
}

func TestPurgeUserByPhoneNumber(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	assert.Nil(t, err)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), base.TestUserPhoneNumber)

	ctx, auth, err := GetTestAuthenticatedContext(t)
	assert.Nil(t, err)
	assert.NotNil(t, auth)

	fr, err := database.NewFirebaseRepository(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, fr)

	// get the user profile and assert the primary phone number
	profile, err := fr.GetUserProfileByUID(ctx, auth.UID)
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, base.TestUserPhoneNumber, *profile.PrimaryPhone)

	// fetch the same profile but now using the primary phone number
	profile, err = fr.GetUserProfileByPrimaryPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, base.TestUserPhoneNumber, *profile.PrimaryPhone)

	// purge the record. this should not fail
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.Nil(t, err)

	// try purging the record again. this should fail since not user profile will be found with the phone number
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.NotNil(t, err)

	// create an invalid user profile
	fakeUID := uuid.New().String()
	invalidpr1, err := fr.CreateUserProfile(context.Background(), base.TestUserPhoneNumber, fakeUID)
	assert.Nil(t, err)
	assert.NotNil(t, invalidpr1)

	// fetch the pins related to invalidpr1. this should fail since no pin has been associated with invalidpr1
	pin, err := fr.GetPINByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, pin)

	// fetch the customer profile related to invalidpr1. this should fail since no customer profile has been associated with invalidpr
	cpr, err := fr.GetCustomerProfileByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, cpr)

	// fetch the supplier profile related to invalidpr1. this should fail since no supplier profile has been associated with invalidpr
	spr, err := fr.GetSupplierProfileByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, spr)

	// call PurgeUserByPhoneNumber using the phone number associated with invalidpr1. this should fail since it does not have
	// an associated pin
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "server error! unable to perform operation")

	// now set a  pin. this should not fail
	userpin := "1234"
	pset, err := s.UserPIN.SetUserPIN(ctx, userpin, base.TestUserPhoneNumber)
	assert.Nil(t, err)
	assert.NotNil(t, pset)
	assert.Equal(t, true, pset)

	// retrieve the pin and assert it matches the one set
	pin, err = fr.GetPINByProfileID(ctx, invalidpr1.ID)
	assert.Nil(t, err)
	assert.NotNil(t, pin)
	matched := utils.ComparePIN(userpin, pin.Salt, pin.PINNumber, nil)
	assert.Equal(t, true, matched)

	// now remove. this should pass even though customer/supplier profile don't exist. What must be removed is the pins
	err = fr.PurgeUserByPhoneNumber(ctx, base.TestUserPhoneNumber)
	assert.Nil(t, err)

	// assert the pin has been removed
	pin, err = fr.GetPINByProfileID(ctx, invalidpr1.ID)
	assert.NotNil(t, err)
	assert.Nil(t, pin)

}
