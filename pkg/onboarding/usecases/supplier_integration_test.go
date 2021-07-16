package usecases_test

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
	"github.com/savannahghi/serverutils"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"
)

const (
	// TestSladeCode is a test slade code for `test` EDI Login
	TestSladeCode = "BRA-PRO-3873-4"

	// TestEDIPortalUsername is a test username for `test` EDI Login
	TestEDIPortalUsername = "malibu.pharmacy-3873@healthcloud.co.ke"

	// TestEDIPortalPassword is a test password for `test` EDI Login
	TestEDIPortalPassword = "test provider one"

	testChargeMasterParentOrgId = "83d3479d-e902-4aab-a27d-6d5067454daf"
	testChargeMasterBranchID    = "94294577-6b27-4091-9802-1ce0f2ce4153"
	primaryEmail                = "test@bewell.co.ke"
)

func cleanUpFirebase(ctx context.Context, t *testing.T) {
	r := fb.Repository{}
	fsc, _ := InitializeTestFirebaseClient(ctx)
	ref := fsc.Collection(r.GetKCYProcessCollectionName())
	firebasetools.DeleteCollection(ctx, fsc, ref, 10)
}

func TestSubmitProcessAddIndividualRiderKycRequest(t *testing.T) {
	// clean kyc processing requests collection because other tests have written to it
	ctx1 := context.Background()
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))
	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual, *spr2.AccountType)
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualRider{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "123456789",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:                         "someKRAPIN",
		KRAPINUploadID:                 "KRAPINUploadID",
		DrivingLicenseID:               "license",
		CertificateGoodConductUploadID: "upload1",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationRider{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualPractitioner{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		RegistrationNumber:      "reg-num",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
		PracticeServices:        []domain.PractitionerService{domain.PractitionerServiceOutpatientServices},
		Cadre:                   domain.PractitionerCadreDoctor,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationPractitioner{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		OrganizationTypeName:    domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:      "reg-num",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
		PracticeServices:        []domain.PractitionerService{domain.PractitionerServiceOutpatientServices},
		Cadre:                   domain.PractitionerCadreDoctor,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationProvider{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		OrganizationTypeName:    domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:      "reg-num",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
		PracticeServices:        []domain.PractitionerService{domain.PractitionerServiceOutpatientServices},
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualPharmaceutical{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "123456789",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		RegistrationNumber:      "reg-num",
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationPharmaceutical{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		OrganizationTypeName:               domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:                 "reg-num",
		PracticeLicenseID:                  "PracticeLicenseID",
		PracticeLicenseUploadID:            "PracticeLicenseUploadID",
		CertificateOfIncorporation:         "cert-org",
		CertificateOfInCorporationUploadID: "cert-org-upload",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeRider

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualCoach{
		IdentificationDoc: domain.Identification{
			IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
			IdentificationDocNumber:         "123456789",
			IdentificationDocNumberUploadID: "id-upload",
		},
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
		AccreditationID:         "ACR-12344568",
		AccreditationUploadID:   "ACR-UPLOAD-12344568",
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "rider"
	partnerType := profileutils.PartnerTypeCoach

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationCoach{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		OrganizationTypeName: domain.OrganizationTypeLimitedCompany,
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)

	// add a partner type for the logged in user
	partnerName := "nutrition"
	partnerType := profileutils.PartnerTypeNutrition

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.IndividualNutrition{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		PracticeLicenseID:       "PracticeLicenseID",
		PracticeLicenseUploadID: "PracticeLicenseUploadID",
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// fetch the profile and assert  the permissions slice is empty
	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 0, len(pr.Permissions))

	// now update the permissions
	perms := []profileutils.PermissionType{profileutils.PermissionTypeAdmin}
	err = s.Onboarding.UpdatePermissions(authenticatedContext, perms)
	assert.Nil(t, err)

	// fetch the profile and assert  the permissions slice is not empty
	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Permissions))

	err = s.Onboarding.UpdatePrimaryEmailAddress(authenticatedContext, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(authenticatedContext, completeUserDetails)
	assert.Nil(t, err)
	// add a partner type for the logged in user
	partnerName := "nutrition"
	partnerType := profileutils.PartnerTypeNutrition

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

	spr2, err := s.Supplier.SetUpSupplier(authenticatedContext, profileutils.AccountTypeIndividual)
	assert.Nil(t, err)
	assert.NotNil(t, spr2)
	assert.Equal(t, profileutils.AccountTypeIndividual.String(), spr2.AccountType.String())
	assert.Equal(t, false, spr2.UnderOrganization)
	assert.Equal(t, false, spr2.IsOrganizationVerified)
	assert.Equal(t, false, spr2.HasBranches)
	assert.Equal(t, false, spr2.Active)

	validInput := domain.OrganizationNutrition{
		KRAPIN:         "someKRAPIN",
		KRAPINUploadID: "KRAPINUploadID",
		SupportingDocuments: []domain.SupportingDocument{
			{
				SupportingDocumentTitle:       "support-title",
				SupportingDocumentDescription: "support-description",
				SupportingDocumentUpload:      "support-upload-id",
			},
		},
		OrganizationTypeName:    domain.OrganizationTypeLimitedCompany,
		RegistrationNumber:      "org-reg-number",
		PracticeLicenseID:       "org-practice-license",
		PracticeLicenseUploadID: "org-practice-license-upload",
		DirectorIdentifications: []domain.Identification{
			{
				IdentificationDocType:           enumutils.IdentificationDocTypeNationalid,
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
	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		cleanUpFirebase(ctx1, t)
	}

	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}

	primaryPhone := interserviceclient.TestUserPhoneNumber

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
		&dto.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     feedlib.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp1)
	assert.NotNil(t, resp1.Profile)
	assert.NotNil(t, resp1.CustomerProfile)
	assert.NotNil(t, resp1.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, feedlib.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		firebasetools.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	cmParentOrgId := testChargeMasterParentOrgId
	filter := []*dto.BranchFilterInput{
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
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.Supplier
		wantErr bool
	}{
		{
			name: "happy :) find supplier by UID",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "sad :( fail to find supplier by UID",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supplier, err := s.Supplier.FindSupplierByUID(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FindSupplierByUID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if supplier != nil {
				if supplier.ID == "" {
					t.Errorf("expected a supplier.")
					return
				}
			}
		})
	}
}

func TestFindSupplierByID(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	name := "Makmende And Sons"
	partnerPractitioner := profileutils.PartnerTypePractitioner
	partnerType, err := s.Supplier.AddPartnerType(ctx, &name, &partnerPractitioner)
	assert.Nil(t, err)
	assert.NotNil(t, partnerType)
	assert.Equal(t, true, partnerType)

	supplier, err := s.Supplier.SetUpSupplier(ctx, profileutils.AccountTypeOrganisation)
	assert.Nil(t, err)

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.Supplier
		wantErr bool
	}{
		{
			name: "happy :) find supplier by ID",
			args: args{
				ctx: ctx,
				id:  supplier.ID,
			},
			wantErr: false,
		},
		{
			name: "happy :) find supplier by ID using wrong context, should not fail",
			args: args{
				ctx: context.Background(),
				id:  supplier.ID,
			},
			wantErr: false,
		},
		{
			name: "sad :( fail to find supplier by ID - wrong supplier ID",
			args: args{
				ctx: context.Background(),
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supplier, err := s.Supplier.FindSupplierByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FindSupplierByID() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
			if supplier != nil {
				if supplier.ID == "" {
					t.Errorf("expected a supplier.")
					return
				}
			}
		})
	}
}

func TestSupplierEDIUserLogin(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	name := "Makmende And Sons"
	partnerPractitioner := profileutils.PartnerTypePractitioner

	// TestEDIPortalUsername is a test username for `test` EDI Login
	TestEDIPortalUsername := "malibu.pharmacy-3873@healthcloud.co.ke"

	// TestEDIPortalPassword is a test password for `test` EDI Login
	TestEDIPortalPassword := "test provider one"

	WrongTestEDIPortalUsername := "username"
	WrongTestEDIPortalPassword := "password"
	EmptyWrongTestEDIPortalUsername := ""
	EmptyTestEDIPortalPassword := ""
	resp2, err := s.Supplier.AddPartnerType(ctx, &name, &partnerPractitioner)
	assert.Nil(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)

	resp3, err := s.Supplier.SetUpSupplier(ctx, profileutils.AccountTypeOrganisation)
	assert.Nil(t, err)
	assert.NotNil(t, resp3)
	assert.Equal(t, false, resp3.Active)
	assert.Nil(t, resp3.EDIUserProfile)

	type args struct {
		ctx      context.Context
		username *string
		password *string
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.EDIUserProfile
		wantErr bool
	}{
		{
			name: "Happy Case: valid credentials",
			args: args{
				ctx:      context.Background(),
				username: &TestEDIPortalUsername,
				password: &TestEDIPortalPassword,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Wrong username and password",
			args: args{
				ctx:      context.Background(),
				username: &WrongTestEDIPortalUsername,
				password: &WrongTestEDIPortalPassword,
			},
			wantErr: true,
		},
		{
			name: "sad case: empty username and password",
			args: args{
				ctx:      context.Background(),
				username: &EmptyWrongTestEDIPortalUsername,
				password: &EmptyTestEDIPortalPassword,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.Supplier.EDIUserLogin(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.EDIUserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestFetchSupplierAllowedLocations(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case :)",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "sad case :( unable to get supplier",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branchConnection, err := s.Supplier.FetchSupplierAllowedLocations(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.FetchSupplierAllowedLocations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !tt.wantErr {
				if branchConnection == nil {
					t.Errorf("expected branch connection")
					return
				}
			}
		})
	}
}

func TestSuspendSupplier(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	suspensionReason := `
	"This email is to inform you that as a result of your actions on April 12th, 2021, you have been issued a suspension for 1 week (7 days)"
	`
	err = s.Onboarding.UpdatePrimaryEmailAddress(ctx, primaryEmail)
	assert.Nil(t, err)

	dateOfBirth2 := scalarutils.Date{
		Day:   12,
		Year:  1995,
		Month: 10,
	}
	firstName2 := "makmende"
	lastName2 := "juha"

	completeUserDetails := profileutils.BioData{
		DateOfBirth: &dateOfBirth2,
		FirstName:   &firstName2,
		LastName:    &lastName2,
	}

	// update biodata
	err = s.Onboarding.UpdateBioData(ctx, completeUserDetails)
	assert.Nil(t, err)

	name := "Makmende And Sons"
	partnerPractitioner := profileutils.PartnerTypePractitioner

	// Add PartnerType
	resp2, err := s.Supplier.AddPartnerType(ctx, &name, &partnerPractitioner)
	assert.Nil(t, err)
	assert.NotNil(t, resp2)
	assert.Equal(t, true, resp2)
	type args struct {
		ctx              context.Context
		suspensionReason *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "sad case: suspend a nonexisting supplier",
			args: args{
				ctx:              context.Background(),
				suspensionReason: &suspensionReason,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Happy case: suspend an existing supplier",
			args: args{
				ctx:              ctx,
				suspensionReason: &suspensionReason,
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Supplier.SuspendSupplier(tt.args.ctx, tt.args.suspensionReason)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.SuspendSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SupplierUseCasesImpl.SuspendSupplier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSupplierUseCasesImpl_AddPartnerType(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	testRiderName := "Test Rider"
	rider := profileutils.PartnerTypeRider
	testPractitionerName := "Test Practitioner"
	practitioner := profileutils.PartnerTypePractitioner
	testProviderName := "Test Provider"
	provider := profileutils.PartnerTypeProvider
	testPharmaceuticalName := "Test Pharmaceutical"
	pharmaceutical := profileutils.PartnerTypePharmaceutical
	testCoachName := "Test Coach"
	coach := profileutils.PartnerTypeCoach
	testNutritionName := "Test Nutrition"
	nutrition := profileutils.PartnerTypeNutrition
	testConsumerName := "Test Consumer"
	consumer := profileutils.PartnerTypeConsumer

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		ctx         context.Context
		name        *string
		partnerType *profileutils.PartnerType
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		wantErr     bool
		expectedErr string
	}{
		{
			name: "valid: add PartnerTypeRider ",
			args: args{
				ctx:         ctx,
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypePractitioner ",
			args: args{
				ctx:         ctx,
				name:        &testPractitionerName,
				partnerType: &practitioner,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeProvider ",
			args: args{
				ctx:         ctx,
				name:        &testProviderName,
				partnerType: &provider,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypePharmaceutical",
			args: args{
				ctx:         ctx,
				name:        &testPharmaceuticalName,
				partnerType: &pharmaceutical,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeCoach",
			args: args{
				ctx:         ctx,
				name:        &testCoachName,
				partnerType: &coach,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "valid: add PartnerTypeNutrition",
			args: args{
				ctx:         ctx,
				name:        &testNutritionName,
				partnerType: &nutrition,
			},
			want:    true,
			wantErr: false,
		},

		{
			name: "invalid: add PartnerTypeConsumer",
			args: args{
				ctx:         ctx,
				name:        &testConsumerName,
				partnerType: &consumer,
			},
			want:        false,
			wantErr:     true,
			expectedErr: "invalid `partnerType`. cannot use CONSUMER in this context",
		},

		{
			name: "invalid : invalid context",
			args: args{
				ctx:         context.Background(),
				name:        &testRiderName,
				partnerType: &rider,
			},
			want:        false,
			wantErr:     true,
			expectedErr: `unable to get the logged in user: auth token not found in context: unable to get auth token from context with key "UID" `,
		},
		{
			name: "invalid : missing name arg",
			args: args{
				ctx: ctx,
			},
			want:        false,
			wantErr:     true,
			expectedErr: "expected `name` to be defined and `partnerType` to be valid",
		},
		{
			name: "invalid: missing partnerType",
			args: args{
				ctx:  ctx,
				name: &testPharmaceuticalName,
			},
			want:        false,
			wantErr:     true,
			expectedErr: "expected `partnerType` should be valid",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supplier := s
			got, err := supplier.Supplier.AddPartnerType(tt.args.ctx, tt.args.name, tt.args.partnerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.AddPartnerType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SupplierUseCasesImpl.AddPartnerType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetUpSupplier(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}

	individualPartner := profileutils.AccountTypeIndividual
	organizationPartner := profileutils.AccountTypeOrganisation

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx         context.Context
		accountType profileutils.AccountType
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Successful individual supplier account setup",
			args: args{
				ctx:         ctx,
				accountType: individualPartner,
			},
			wantErr: false,
		},
		{
			name: "Successful organization supplier account setup",
			args: args{
				ctx:         ctx,
				accountType: organizationPartner,
			},
			wantErr: false,
		},
		{
			name: "SadCase - Invalid supplier setup",
			args: args{
				ctx:         ctx,
				accountType: "non existent type",
			},
			wantErr: true,
		},
		{
			name: "SadCase - unauthenticated context",
			args: args{
				ctx:         context.Background(),
				accountType: organizationPartner,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			supplier, err := s.Supplier.SetUpSupplier(tt.args.ctx, tt.args.accountType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetUpSupplier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if supplier == nil && !tt.wantErr {
				t.Errorf("expected a supplier and nil error but got: %v", err)
				return
			}

			if supplier != nil && tt.wantErr {
				t.Errorf("expected an error but instead got a nil")
				return
			}
		})
	}

}

func TestSupplierUseCasesImpl_EDIUserLogin(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	validUsername := TestEDIPortalUsername
	validPassword := TestEDIPortalPassword

	invalidUsername := "username"
	invalidPassword := "password"

	emptyUsername := ""
	emptyPassword := ""
	type args struct {
		ctx      context.Context
		username *string
		password *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: valid credentials",
			args: args{
				ctx:      context.Background(),
				username: &validUsername,
				password: &validPassword,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Wrong username and password",
			args: args{
				ctx:      context.Background(),
				username: &invalidUsername,
				password: &invalidPassword,
			},
			wantErr: true,
		},
		{
			name: "sad case: empty username and password",
			args: args{
				ctx:      context.Background(),
				username: &emptyUsername,
				password: &emptyPassword,
			},
			wantErr: true,
		},
		{
			name: "sad case:valid username invalid password",
			args: args{
				ctx:      context.Background(),
				username: &validUsername,
				password: &invalidPassword,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ediLogin := s
			_, err := ediLogin.Supplier.EDIUserLogin(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.EDIUserLogin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// func TestSupplierUseCasesImpl_CoreEDIUserLogin(t *testing.T) {
// 	ctx, _, err := GetTestAuthenticatedContext(t)
// 	if err != nil {
// 		t.Errorf("failed to get test authenticated context: %v", err)
// 		return
// 	}
// 	s, err := InitializeTestService(ctx)
// 	if err != nil {
// 		t.Errorf("unable to initialize test service")
// 		return
// 	}
// 	type args struct {
// 		ctx      context.Context
// 		username string
// 		password string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Happy Case: valid credentials",
// 			args: args{
// 				ctx:      context.Background(),
// 				username: "bewell@slade360.co.ke",
// 				password: "please change me",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Sad Case: Wrong userame and password",
// 			args: args{
// 				ctx:      context.Background(),
// 				username: "invalid Username",
// 				password: "invalid Password",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "sad case: empty username and password",
// 			args: args{
// 				ctx:      context.Background(),
// 				username: "",
// 				password: "",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "sad case: valid username and wrong password",
// 			args: args{
// 				username: "bewell@slade360.co.ke",
// 				password: "invalid Password",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			coreEdiLogin := s
// 			_, err := coreEdiLogin.Supplier.CoreEDIUserLogin(tt.args.ctx, tt.args.username, tt.args.password)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("SupplierUseCasesImpl.CoreEDIUserLogin() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

func clean(newCtx context.Context, testPhoneNumber string, t *testing.T, service *interactor.Interactor) {
	err := service.Signup.RemoveUserByPhoneNumber(newCtx, testPhoneNumber)
	if err != nil {
		t.Errorf("failed to clean data after test error: %v", err)
		return
	}
}

func TestCreateCustomerAccount(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		ctx         context.Context
		name        string
		partnerType profileutils.PartnerType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy:) create customer account",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: profileutils.PartnerTypeConsumer,
			},
			wantErr: false,
		},
		{
			name: "sad:( wrong partner type",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: profileutils.PartnerTypeCoach,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Supplier.CreateCustomerAccount(
				tt.args.ctx,
				tt.args.name,
				tt.args.partnerType,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.CreateCustomerAccount() error = %v, wantErr %v",
					err,
					tt.wantErr,
				)
				return
			}
		})
	}
}

func TestCreateSupplierAccount(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		ctx         context.Context
		name        string
		partnerType profileutils.PartnerType
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.Supplier
		wantErr bool
	}{
		{
			name: "happy:) create supplier account",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: profileutils.PartnerTypeRider,
			},
			wantErr: false,
		},
		{
			name: "sad:( wrong partner type",
			args: args{
				ctx:         ctx,
				name:        *utils.GetRandomName(),
				partnerType: profileutils.PartnerTypeConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.Supplier.CreateSupplierAccount(tt.args.ctx, tt.args.name, tt.args.partnerType)
			if (err != nil) != tt.wantErr {
				t.Errorf("SupplierUseCasesImpl.CreateSupplierAccount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// func TestSupplierUseCasesImpl_CheckSupplierKYCSubmitted(t *testing.T) {
// 	ctx, _, err := GetTestAuthenticatedContext(t)
// 	if err != nil {
// 		t.Errorf("failed to get test authenticated context: %v", err)
// 		return
// 	}
// 	s, err := InitializeTestService(ctx)
// 	if err != nil {
// 		t.Errorf("unable to initialize test service")
// 		return
// 	}
//
// 	err = s.Onboarding.UpdatePrimaryEmailAddress(ctx, primaryEmail)
// 	assert.Nil(t, err)

// 	dateOfBirth2 := scalarutils.Date{
// 		Day:   12,
// 		Year:  1995,
// 		Month: 10,
// 	}
// 	firstName2 := "makmende"
// 	lastName2 := "juha"

// 	completeUserDetails := profileutils.BioData{
// 		DateOfBirth: &dateOfBirth2,
// 		FirstName:   &firstName2,
// 		LastName:    &lastName2,
// 	}

// 	// update biodata
// 	err = s.Onboarding.UpdateBioData(ctx, completeUserDetails)
// 	assert.Nil(t, err)

// 	// add a partner type for the logged in user
// 	partnerName := "nutrition"
// 	partnerType := profileutils.PartnerTypeNutrition

// 	resp2, err := s.Supplier.AddPartnerType(ctx, &partnerName, &partnerType)
// 	assert.Nil(t, err)
// 	assert.Equal(t, true, resp2)
// 	_, err = s.Supplier.SetUpSupplier(ctx, profileutils.AccountTypeIndividual)
// 	if err != nil {
// 		t.Errorf("unable to setup supplier")
// 		return
// 	}
// 	validInput := domain.IndividualNutrition{
// 		KRAPIN:         "someKRAPIN",
// 		KRAPINUploadID: "KRAPINUploadID",
// 		SupportingDocuments: []domain.SupportingDocument{
// 			{
// 				SupportingDocumentTitle:       "support-title",
// 				SupportingDocumentDescription: "support-description",
// 				SupportingDocumentUpload:      "support-upload-id",
// 			},
// 		},
// 		PracticeLicenseID:       "PracticeLicenseID",
// 		PracticeLicenseUploadID: "PracticeLicenseUploadID",
// 	}
// 	_, err = s.Supplier.AddIndividualNutritionKyc(ctx, validInput)
// 	if err != nil {
// 		t.Errorf("unable to add individual nutrition KYC")
// 		return
// 	}
// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "happy:) check supplier KYC submitted",
// 			args: args{
// 				ctx: ctx,
// 			},
// 			want:    true,
// 			wantErr: false,
// 		},
// 		{
// 			name: "sad:( check supplier KYC submitted",
// 			args: args{
// 				ctx: context.Background(),
// 			},
// 			want:    false,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := s.Supplier.CheckSupplierKYCSubmitted(tt.args.ctx)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("SupplierUseCasesImpl.CheckSupplierKYCSubmitted() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got != tt.want {
// 				t.Errorf("SupplierUseCasesImpl.CheckSupplierKYCSubmitted() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
