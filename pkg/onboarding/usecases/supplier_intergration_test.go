package usecases_test

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
)

func TestSubmitProcessAddIndividualRiderKycRequest(t *testing.T) {
	// clean kyc processing requests collection because other tests have written to it
	ctx1 := context.Background()
	r := database.Repository{} // They are nil
	fsc, _ := InitializeTestFirebaseClient(ctx1)
	ref := fsc.Collection(r.GetKCYProcessCollectionName())
	base.DeleteCollection(ctx1, fsc, ref, 10)

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := base.TestUserPhoneNumber

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp1, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := base.PartnerTypeRider

	resp2, err := s.Supplier.AddPartnerType(authenticatedContext, &partnerName, &partnerType)
	assert.Nil(t, err)
	assert.Equal(t, true, resp2)

	// fetch the supplier profile and assert that the partner type and name is as was added above

	spr1, err := s.Supplier.FindSupplierByUID(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, spr1)
	assert.NotNil(t, spr1.PartnerType)
	assert.NotNil(t, spr1.SupplierName)
	assert.NotNil(t, spr1.PartnerSetupComplete)
	assert.Equal(t, partnerType.String(), spr1.PartnerType.String())
	assert.Equal(t, partnerName, spr1.SupplierName)
	assert.Equal(t, true, spr1.PartnerSetupComplete)

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, base.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, base.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualRider{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "123456789",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:                         "someKRAPIN",
		KRAPINUploadID:                 "KRAPINUploadID",
		DrivingLicenseID:               "license",
		CertificateGoodConductUploadID: "upload1",
		SupportingDocumentsUploadID:    []string{"SupportingDocumentsUploadID", "Support"},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddIndividualRiderKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddIndividualRiderKyc(authenticatedContext, validInput)
	assert.NotNil(t, err)
	assert.Nil(t, kyc2)

	// now fetch kyc processing requests
	kycrequests, err := s.Supplier.FetchKYCProcessingRequests(authenticatedContext)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(kycrequests))

	firstKYC := kycrequests[0]
	assert.Equal(t, false, firstKYC.Processed)

	response, err := s.Supplier.ProcessKYCRequest(authenticatedContext, firstKYC.ID, domain.KYCProcessStatusApproved, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, response)

	clean(authenticatedContext, primaryPhone, t, s)
}

//todo(dexter): TestSubmitProcessOrganizationRiderKyc

// todo(dexter): TestSubmitProcessIndividualPractitionerKyc

// todo(dexter) : TestSubmitProcessOrganizationPractitionerKyc

// todo(dexter): TestSubmitProcessOrganizationProviderKyc

// todo(dexter) : TestSubmitProcessIndividualPharmaceuticalKyc

// todo(dexter): TestSubmitProcessOrganizationPharmaceuticalKyc

// todo(dexter) : TestSubmitProcessIndividualCoachKyc

// todo(dexter) : TestSubmitProcessOrganizationCoachKyc

func TestSubmitProcessIndividualNutritionKycRequest(t *testing.T) {
	// clean kyc processing requests collection because other tests have written to it
	ctx1 := context.Background()
	r := database.Repository{} // They are nil
	fsc, _ := InitializeTestFirebaseClient(ctx1)
	ref := fsc.Collection(r.GetKCYProcessCollectionName())
	base.DeleteCollection(ctx1, fsc, ref, 10)

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := base.TestUserPhoneNumber

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp1, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// add a partner type for the logged in user
	partnerName := "nutrition"
	partnerType := base.PartnerTypeNutrition

	resp2, err := s.Supplier.AddPartnerType(authenticatedContext, &partnerName, &partnerType)
	assert.Nil(t, err)
	assert.Equal(t, true, resp2)

	// fetch the supplier profile and assert that the partner type and name is as was added above

	spr1, err := s.Supplier.FindSupplierByUID(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, spr1)
	assert.NotNil(t, spr1.PartnerType)
	assert.NotNil(t, spr1.SupplierName)
	assert.NotNil(t, spr1.PartnerSetupComplete)
	assert.Equal(t, partnerType.String(), spr1.PartnerType.String())
	assert.Equal(t, partnerName, spr1.SupplierName)
	assert.Equal(t, true, spr1.PartnerSetupComplete)

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, base.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, base.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualNutrition{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddIndividualNutritionKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddIndividualNutritionKyc(authenticatedContext, validInput)
	assert.NotNil(t, err)
	assert.Nil(t, kyc2)

	// now fetch kyc processing requests
	kycrequests, err := s.Supplier.FetchKYCProcessingRequests(authenticatedContext)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(kycrequests))

	firstKYC := kycrequests[0]
	assert.Equal(t, false, firstKYC.Processed)

	response, err := s.Supplier.ProcessKYCRequest(authenticatedContext, firstKYC.ID, domain.KYCProcessStatusApproved, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, response)

}

func TestSubmitProcessOrganizationNutritionKycRequest(t *testing.T) {
	// clean kyc processing requests collection because other tests have written to it
	ctx1 := context.Background()
	r := database.Repository{} // They are nil
	fsc, _ := InitializeTestFirebaseClient(ctx1)
	ref := fsc.Collection(r.GetKCYProcessCollectionName())
	base.DeleteCollection(ctx1, fsc, ref, 10)

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := base.TestUserPhoneNumber

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp1, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// add a partner type for the logged in user
	partnerName := "nutrition"
	partnerType := base.PartnerTypeNutrition

	resp2, err := s.Supplier.AddPartnerType(authenticatedContext, &partnerName, &partnerType)
	assert.Nil(t, err)
	assert.Equal(t, true, resp2)

	// fetch the supplier profile and assert that the partner type and name is as was added above

	spr1, err := s.Supplier.FindSupplierByUID(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, spr1)
	assert.NotNil(t, spr1.PartnerType)
	assert.NotNil(t, spr1.SupplierName)
	assert.NotNil(t, spr1.PartnerSetupComplete)
	assert.Equal(t, partnerType.String(), spr1.PartnerType.String())
	assert.Equal(t, partnerName, spr1.SupplierName)
	assert.Equal(t, true, spr1.PartnerSetupComplete)

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, base.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, base.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationNutrition{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		OrganizationCertificate:     "org-cert",
		RegistrationNumber:          "org-reg-number",
		PracticeLicenseID:           "org-practice-license",
		PracticeLicenseUploadID:     "org-practice-license-upload",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "123456789",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddOrganizationNutritionKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddOrganizationNutritionKyc(authenticatedContext, validInput)
	assert.NotNil(t, err)
	assert.Nil(t, kyc2)

	// now fetch kyc processing requests
	kycrequests, err := s.Supplier.FetchKYCProcessingRequests(authenticatedContext)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(kycrequests))

	firstKYC := kycrequests[0]
	assert.Equal(t, false, firstKYC.Processed)

	response, err := s.Supplier.ProcessKYCRequest(authenticatedContext, firstKYC.ID, domain.KYCProcessStatusApproved, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, response)
}
