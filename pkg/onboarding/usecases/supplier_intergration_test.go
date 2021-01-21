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

func TestSubmitProcessOrganizationRiderKycRequest(t *testing.T) {
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

	validInput := domain.OrganizationRider{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		OrganizationCertificate:     "org-cert",
		OrganizationTypeName:        domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "123456789",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddOrganizationRiderKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddOrganizationRiderKyc(authenticatedContext, validInput)
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

func TestSubmitProcessIndividualPractitionerKyc(t *testing.T) {
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

	validInput := domain.IndividualPractitioner{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		RegistrationNumber:          "reg-num",
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
		PracticeServices:            []domain.PractitionerService{domain.PractitionerServiceOutpatientServices},
		Cadre:                       domain.PractitionerCadreDoctor,
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddIndividualPractitionerKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddIndividualPractitionerKyc(authenticatedContext, validInput)
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

func TestSubmitProcessOrganizationPractitionerKyc(t *testing.T) {
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

	validInput := domain.OrganizationPractitioner{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		OrganizationCertificate:     "org-cert",
		OrganizationTypeName:        domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:          "reg-num",
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
		PracticeServices:            []domain.PractitionerService{domain.PractitionerServiceOutpatientServices},
		Cadre:                       domain.PractitionerCadreDoctor,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "123456789",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddOrganizationPractitionerKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddOrganizationPractitionerKyc(authenticatedContext, validInput)
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

func TestSubmitProcessOrganizationProviderKyc(t *testing.T) {
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

	validInput := domain.OrganizationProvider{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		OrganizationCertificate:     "org-cert",
		OrganizationTypeName:        domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:          "reg-num",
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
		PracticeServices:            []domain.PractitionerService{domain.PractitionerServiceOutpatientServices},
		Cadre:                       domain.PractitionerCadreDoctor,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "123456789",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddOrganizationProviderKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddOrganizationProviderKyc(authenticatedContext, validInput)
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

func TestSubmitProcessIndividualPharmaceuticalKyc(t *testing.T) {
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

	validInput := domain.IndividualPharmaceutical{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "123456789",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		RegistrationNumber:          "reg-num",
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddIndividualPharmaceuticalKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddIndividualPharmaceuticalKyc(authenticatedContext, validInput)
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

func TestSubmitProcessOrganizationPharmaceuticalKyc(t *testing.T) {
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

	validInput := domain.OrganizationPharmaceutical{
		KRAPIN:                             "someKRAPIN",
		KRAPINUploadID:                     "KRAPINUploadID",
		SupportingDocumentsUploadID:        []string{"SupportingDocumentsUploadID", "Support"},
		OrganizationCertificate:            "org-cert",
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:                 "reg-num",
		PracticeLicenseID:                  "PracticeLicenseID",
		PracticeLicenseUploadID:            "PracticeLicenseUploadID",
		CertificateOfIncorporation:         "cert-org",
		CertificateOfInCorporationUploadID: "cert-org-upload",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "123456789",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddOrganizationPharmaceuticalKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddOrganizationPharmaceuticalKyc(authenticatedContext, validInput)
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

func TestSubmitProcessIndividualCoachKyc(t *testing.T) {
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

	validInput := domain.IndividualCoach{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           domain.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "123456789",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		PracticeLicenseID:           "PracticeLicenseID",
		PracticeLicenseUploadID:     "PracticeLicenseUploadID",
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddIndividualCoachKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddIndividualCoachKyc(authenticatedContext, validInput)
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

func TestSubmitProcessOrganizationCoachKycRequest(t *testing.T) {
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

	validInput := domain.OrganizationCoach{
		KRAPIN:                      "someKRAPIN",
		KRAPINUploadID:              "KRAPINUploadID",
		SupportingDocumentsUploadID: []string{"SupportingDocumentsUploadID", "Support"},
		OrganizationCertificate:     "org-cert",
		OrganizationTypeName:        domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           domain.IdentificationDocTypeNationalid,
				IdentificationDocNumber:         "123456789",
				IdentificationDocNumberUploadID: "id-upload",
			},
		},
	}

	// submit first kyc. this should pass
	kyc1, err := s.Supplier.AddOrganizationCoachKyc(authenticatedContext, validInput)
	assert.Nil(t, err)
	assert.NotNil(t, kyc1)

	// submit another kyc. this should fail
	kyc2, err := s.Supplier.AddOrganizationCoachKyc(authenticatedContext, validInput)
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
		OrganizationTypeName:        domain.OrganizationTypeLimitedCompany,
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

func TestSupplierSetDefaultLocation(t *testing.T) {
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

	cmParentOrgId := testChargeMasterParentOrgId
	filter := []*resources.BranchFilterInput{
		{
			ParentOrganizationID: &cmParentOrgId,
		},
	}

	br, err := s.ChargeMaster.FindBranch(authenticatedContext, nil, filter, nil)
	assert.Nil(t, err)
	assert.NotNil(t, br)
	assert.NotEqual(t, 0, len(br.Edges))

	// call set supplier default location
	spr, err := s.Supplier.SupplierSetDefaultLocation(authenticatedContext, br.Edges[0].Node.ID)
	assert.Nil(t, err)
	assert.NotNil(t, spr)
	assert.Equal(t, br.Edges[0].Node.ID, spr.Location.ID)
}

func TestFindSupplierByUID(t *testing.T) {

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
			Flavour:     base.FlavourPro,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.NotNil(t, login1.SupplierProfile)
	assert.Equal(t, resp1.SupplierProfile.ID, login1.SupplierProfile.ID)
	assert.Equal(t, resp1.SupplierProfile.ProfileID, login1.SupplierProfile.ProfileID)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the supplier profile with the uid
	spr, err := s.Supplier.FindSupplierByUID(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, spr)
	assert.Equal(t, login1.SupplierProfile.ID, spr.ID)
	assert.Equal(t, login1.SupplierProfile.ProfileID, spr.ProfileID)
	assert.Equal(t, login1.SupplierProfile.Active, spr.Active)
	assert.Equal(t, login1.SupplierProfile.AccountType.String(), spr.AccountType.String())

	// try using the wrong context. shoild should fail
	spr, err = s.Supplier.FindSupplierByUID(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, spr)
}

func TestFindSupplierByID(t *testing.T) {

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
			Flavour:     base.FlavourPro,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.NotNil(t, login1.SupplierProfile)
	assert.Equal(t, resp1.SupplierProfile.ID, login1.SupplierProfile.ID)
	assert.Equal(t, resp1.SupplierProfile.ProfileID, login1.SupplierProfile.ProfileID)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the supplier profile with the id
	spr, err := s.Supplier.FindSupplierByID(authenticatedContext, login1.SupplierProfile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, spr)
	assert.Equal(t, login1.SupplierProfile.ID, spr.ID)
	assert.Equal(t, login1.SupplierProfile.ProfileID, spr.ProfileID)
	assert.Equal(t, login1.SupplierProfile.Active, spr.Active)
	assert.Equal(t, login1.SupplierProfile.AccountType.String(), spr.AccountType.String())

	// try using the wrong context. shoild should not fail
	spr, err = s.Supplier.FindSupplierByID(context.Background(), login1.SupplierProfile.ID)
	assert.Nil(t, err)
	assert.NotNil(t, spr)
}

func TestSupplierEDILogin(t *testing.T) {
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
			Flavour:     base.FlavourPro,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.NotNil(t, login1.SupplierProfile)
	assert.Equal(t, resp1.SupplierProfile.ID, login1.SupplierProfile.ID)
	assert.Equal(t, resp1.SupplierProfile.ProfileID, login1.SupplierProfile.ProfileID)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	name := "Makmende And Sons"
	partnerPractitioner := base.PartnerTypePractitioner
	resp2, err := s.Supplier.AddPartnerType(authenticatedContext, &name, &partnerPractitioner)
	assert.Nil(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)

	resp3, err := s.Supplier.SetUpSupplier(authenticatedContext, base.AccountTypeOrganisation)
	assert.Nil(t, err)
	assert.NotNil(t, resp3)
	assert.Equal(t, false, resp3.Active)
	assert.Nil(t, resp3.EDIUserProfile)

	resp4, err := s.Supplier.SupplierEDILogin(authenticatedContext, testEDIPortalUsername, testEDIPortalPassword, testSladeCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp4)
	assert.Nil(t, resp4.Supplier)
	assert.NotNil(t, resp4.Branches)
}

func TestFetchSupplierAllowedLocations(t *testing.T) {

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
			Flavour:     base.FlavourPro,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.NotNil(t, login1.SupplierProfile)
	assert.Equal(t, resp1.SupplierProfile.ID, login1.SupplierProfile.ID)
	assert.Equal(t, resp1.SupplierProfile.ProfileID, login1.SupplierProfile.ProfileID)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	name := "Makmende And Sons"
	partnerPractitioner := base.PartnerTypePractitioner
	resp2, err := s.Supplier.AddPartnerType(authenticatedContext, &name, &partnerPractitioner)
	assert.Nil(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)

	resp3, err := s.Supplier.SetUpSupplier(authenticatedContext, base.AccountTypeOrganisation)
	assert.Nil(t, err)
	assert.NotNil(t, resp3)
	assert.Equal(t, false, resp3.Active)
	assert.Nil(t, resp3.EDIUserProfile)

	resp4, err := s.Supplier.SupplierEDILogin(authenticatedContext, testEDIPortalUsername, testEDIPortalPassword, testSladeCode)
	assert.Nil(t, err)
	assert.NotNil(t, resp4)
	assert.Nil(t, resp4.Supplier)
	assert.NotNil(t, resp4.Branches)

	// fetch all AllowedLocations for the suppier
	resp5, err := s.Supplier.FetchSupplierAllowedLocations(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, resp5)
	assert.Equal(t, len(resp4.Branches.Edges), len(resp5.Edges))

}

func TestSuspendSupplier(t *testing.T) {

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
			Flavour:     base.FlavourPro,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourPro)
	assert.Nil(t, err)
	assert.NotNil(t, login1)
	assert.NotNil(t, login1.SupplierProfile)
	assert.Equal(t, resp1.SupplierProfile.ID, login1.SupplierProfile.ID)
	assert.Equal(t, resp1.SupplierProfile.ProfileID, login1.SupplierProfile.ProfileID)
	assert.Equal(t, resp1.SupplierProfile.ProfileID, login1.SupplierProfile.ProfileID)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	name := "Makmende And Sons"
	partnerPractitioner := base.PartnerTypePractitioner
	resp2, err := s.Supplier.AddPartnerType(authenticatedContext, &name, &partnerPractitioner)
	assert.Nil(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)

	resp3, err := s.Supplier.SetUpSupplier(authenticatedContext, base.AccountTypeOrganisation)
	assert.Nil(t, err)
	assert.NotNil(t, resp3)
	assert.Equal(t, false, resp3.Active)
	assert.Nil(t, resp3.EDIUserProfile)

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

	// now fetch kyc processing requests
	kycrequests, err := s.Supplier.FetchKYCProcessingRequests(authenticatedContext)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(kycrequests))

	firstKYC := kycrequests[0]
	assert.Equal(t, false, firstKYC.Processed)

	response, err := s.Supplier.ProcessKYCRequest(authenticatedContext, firstKYC.ID, domain.KYCProcessStatusApproved, nil)
	assert.Nil(t, err)
	assert.Equal(t, true, response)

	// fetch the supplier profile. active should be true now
	sup, err := s.Supplier.FindSupplierByUID(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, sup)
	assert.Equal(t, true, sup.Active)

	// now suspend the susplier
	v, err := s.Supplier.SuspendSupplier(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, true, v)

	// fetch the supplier profile. active should be false now
	sup, err = s.Supplier.FindSupplierByUID(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, sup)
	assert.Equal(t, false, sup.Active)

}
