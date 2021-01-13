package database_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
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
