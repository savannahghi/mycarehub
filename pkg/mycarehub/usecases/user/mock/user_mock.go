package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockLoginFn                             func(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, bool)
	MockInviteUserFn                        func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error)
	MockSavePinFn                           func(ctx context.Context, input dto.PINInput) (bool, error)
	MockSetNickNameFn                       func(ctx context.Context, userID string, nickname string) (bool, error)
	MockRequestPINResetFn                   func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error)
	MockResetPINFn                          func(ctx context.Context, input dto.UserResetPinInput) (bool, error)
	MockRefreshTokenFn                      func(ctx context.Context, userID string) (*dto.AuthCredentials, error)
	MockVerifyPINFn                         func(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error)
	MockGetClientCaregiverFn                func(ctx context.Context, clientID string) (*domain.Caregiver, error)
	MockCreateOrUpdateClientCaregiverFn     func(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error)
	MockRegisterClientFn                    func(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error)
	MockRefreshGetStreamTokenFn             func(ctx context.Context, userID string) (*dto.GetStreamToken, error)
	MockSearchClientUserFn                  func(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error)
	MockFetchContactOrganisationsFn         func(ctx context.Context, phoneNumber string) ([]*domain.Organisation, error)
	MockCompleteOnboardingTourFn            func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockRegisterKenyaEMRPatientsFn          func(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.PatientRegistrationPayload, error)
	MockRegisteredFacilityPatientsFn        func(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error)
	MockSetUserPINFn                        func(ctx context.Context, input dto.PINInput) (bool, error)
	MockSearchStaffUserFn                   func(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error)
	MockConsentFn                           func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (bool, error)
	MockGetUserProfileFn                    func(ctx context.Context, userID string) (*domain.User, error)
	MockAddClientFHIRIDFn                   func(ctx context.Context, input dto.ClientFHIRPayload) error
	MockGenerateTemporaryPinFn              func(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error)
	MockRegisterPushTokenFn                 func(ctx context.Context, token string) (bool, error)
	MockGetClientProfileByCCCNumberFn       func(ctx context.Context, cccNumber string) (*domain.ClientProfile, error)
	MockRegisterStaffFn                     func(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error)
	MockDeleteUserFn                        func(ctx context.Context, payload *dto.PhoneInput) (bool, error)
	MockTransferClientToFacilityFn          func(ctx context.Context, clientID *string, facilityID *string) (bool, error)
	MockSetStaffDefaultFacilityFn           func(ctx context.Context, staffID string, facilityID string) (*domain.Facility, error)
	MockSetClientDefaultFacilityFn          func(ctx context.Context, clientID string, facilityID string) (*domain.Facility, error)
	MockRemoveFacilitiesFromClientProfileFn func(ctx context.Context, clientID string, facilities []string) (bool, error)
	MockAddFacilitiesToStaffProfileFn       func(ctx context.Context, staffID string, facilities []string) (bool, error)
	MockGetUserLinkedFacilitiesFn           func(ctx context.Context, userID string, paginationInput dto.PaginationsInput) (*dto.FacilityOutputPage, error)
	MockAddFacilitiesToClientProfileFn      func(ctx context.Context, clientID string, facilities []string) (bool, error)
	MockRegisterCaregiver                   func(ctx context.Context, input dto.CaregiverInput) (*domain.CaregiverProfile, error)
	MockSearchCaregiverUserFn               func(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error)
	MockAssignCaregiverFn                   func(ctx context.Context, input dto.ClientCaregiverInput) (bool, error)
	MockRemoveFacilitiesFromStaffProfileFn  func(ctx context.Context, staffID string, facilities []string) (bool, error)
	MockRegisterClientAsCaregiverFn         func(ctx context.Context, clientID string, caregiverNumber string) (*domain.CaregiverProfile, error)
	MockGetCaregiverManagedClientsFn        func(ctx context.Context, caregiverID string, input dto.PaginationsInput) (*dto.ManagedClientOutputPage, error)
	MockListClientsCaregiversFn             func(ctx context.Context, clientID string, pagination *dto.PaginationsInput) (*dto.CaregiverProfileOutputPage, error)
	MockConsentToAClientCaregiverFn         func(ctx context.Context, clientID string, caregiverID string, consent bool) (bool, error)
	MockConsentToManagingClientFn           func(ctx context.Context, caregiverID string, clientID string, consent bool) (bool, error)
}

// NewUserUseCaseMock creates in initializes create type mocks
func NewUserUseCaseMock() *UserUseCaseMock {
	var UUID = uuid.New().String()
	name := gofakeit.Name()
	facilityInput := &domain.Facility{
		ID:          &UUID,
		Name:        name,
		Phone:       gofakeit.Phone(),
		Active:      true,
		County:      gofakeit.Name(),
		Description: gofakeit.Sentence(5),
	}

	staff := &domain.StaffProfile{
		ID:              &UUID,
		User:            &domain.User{},
		UserID:          uuid.New().String(),
		Active:          true,
		StaffNumber:     "test-staff-101",
		Facilities:      []*domain.Facility{},
		DefaultFacility: facilityInput,
	}

	user := &domain.User{
		ID:       &UUID,
		Username: "test",
		Name:     "test",
		Gender:   enumutils.GenderMale,
		Active:   true,
		Contacts: &domain.Contact{
			ID:           &UUID,
			ContactType:  "phone",
			ContactValue: interserviceclient.TestUserPhoneNumber,
			Active:       false,
			OptedIn:      false,
			UserID:       &UUID,
		},
		PushTokens:             []string{},
		LastSuccessfulLogin:    &time.Time{},
		LastFailedLogin:        &time.Time{},
		FailedLoginCount:       0,
		NextAllowedLogin:       &time.Time{},
		PinChangeRequired:      false,
		HasSetPin:              false,
		HasSetSecurityQuestion: false,
		IsPhoneVerified:        false,
		TermsAccepted:          false,
		AcceptedTermsID:        0,
		Suspended:              false,
		Avatar:                 "",
		Roles:                  []*domain.AuthorityRole{},
		Permissions:            []*domain.AuthorityPermission{},
		DateOfBirth:            &time.Time{},
		FailedSecurityCount:    0,
		PinUpdateRequired:      false,
		HasSetNickname:         false,
	}
	clientProfile := &domain.ClientProfile{
		ID:                      &UUID,
		User:                    user,
		Active:                  false,
		ClientTypes:             []enums.ClientType{},
		UserID:                  UUID,
		TreatmentEnrollmentDate: &time.Time{},
		FHIRPatientID:           &UUID,
		HealthRecordID:          &UUID,
		TreatmentBuddy:          "",
		ClientCounselled:        true,
		OrganisationID:          UUID,
		DefaultFacility:         facilityInput,
		CHVUserID:               &UUID,
		CHVUserName:             name,
		CaregiverID:             &UUID,
		CCCNumber:               "123456789",
		Facilities:              []*domain.Facility{facilityInput},
	}

	paginationOutput := &domain.Pagination{
		Limit:        10,
		CurrentPage:  1,
		Count:        1,
		TotalPages:   1,
		NextPage:     nil,
		PreviousPage: nil,
		Sort: &domain.SortParam{
			Field:     "id",
			Direction: enums.SortDataTypeDesc,
		},
	}

	return &UserUseCaseMock{

		MockLoginFn: func(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, bool) {
			resp := &dto.Response{
				User: &dto.User{
					ID:               *user.ID,
					Name:             user.Name,
					Username:         user.Username,
					Active:           user.Active,
					NextAllowedLogin: *user.NextAllowedLogin,
					FailedLoginCount: user.FailedLoginCount,
				},
				AuthCredentials: dto.AuthCredentials{RefreshToken: gofakeit.HipsterSentence(15), IDToken: gofakeit.BeerAlcohol(), ExpiresIn: gofakeit.BeerHop()},
				GetStreamToken:  "",
			}
			return &dto.LoginResponse{
				Response: resp,
				Attempts: 10,
				Message:  "Success",
				Code:     10,
			}, true
		},
		MockRegisterCaregiver: func(ctx context.Context, input dto.CaregiverInput) (*domain.CaregiverProfile, error) {
			return &domain.CaregiverProfile{
				ID: UUID,
				User: domain.User{
					ID: &UUID,
				},
				CaregiverNumber: gofakeit.SSN(),
			}, nil
		},

		MockFetchContactOrganisationsFn: func(ctx context.Context, phoneNumber string) ([]*domain.Organisation, error) {
			return []*domain.Organisation{
				{
					ID:               gofakeit.UUID(),
					Active:           true,
					OrganisationCode: gofakeit.SSN(),
					Name:             gofakeit.Company(),
					Description:      "some description",
					EmailAddress:     gofakeit.Email(),
					PhoneNumber:      gofakeit.Phone(),
					DefaultCountry:   gofakeit.Country(),
				},
			}, nil
		},
		MockRegisterClientAsCaregiverFn: func(ctx context.Context, clientID, caregiverNumber string) (*domain.CaregiverProfile, error) {
			return &domain.CaregiverProfile{
				ID: UUID,
				User: domain.User{
					ID: &UUID,
				},
				CaregiverNumber: gofakeit.SSN(),
			}, nil
		},
		MockInviteUserFn: func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
			return true, nil
		},
		MockSavePinFn: func(ctx context.Context, input dto.PINInput) (bool, error) {
			return true, nil
		},
		MockSetNickNameFn: func(ctx context.Context, userID, nickname string) (bool, error) {
			return true, nil
		},
		MockRequestPINResetFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
			return "111222", nil
		},
		MockResetPINFn: func(ctx context.Context, input dto.UserResetPinInput) (bool, error) {
			return true, nil
		},
		MockRefreshTokenFn: func(ctx context.Context, userID string) (*dto.AuthCredentials, error) {
			return &dto.AuthCredentials{
				RefreshToken: uuid.New().String(),
				ExpiresIn:    "3600",
				IDToken:      uuid.New().String(),
			}, nil
		},
		MockVerifyPINFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
			return true, nil
		},
		MockDeleteUserFn: func(ctx context.Context, payload *dto.PhoneInput) (bool, error) {
			return true, nil
		},
		MockSearchStaffUserFn: func(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
			return []*domain.StaffProfile{staff}, nil
		},

		MockCreateOrUpdateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
			return true, nil
		},
		MockRegisterClientFn: func(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error) {
			return &dto.ClientRegistrationOutput{
				ID: uuid.New().String(),
			}, nil
		},
		MockAssignCaregiverFn: func(ctx context.Context, input dto.ClientCaregiverInput) (bool, error) {
			return true, nil
		},
		MockRefreshGetStreamTokenFn: func(ctx context.Context, userID string) (*dto.GetStreamToken, error) {
			return &dto.GetStreamToken{
				Token: uuid.New().String(),
			}, nil
		},
		MockListClientsCaregiversFn: func(ctx context.Context, clientID string, pagination *dto.PaginationsInput) (*dto.CaregiverProfileOutputPage, error) {
			return &dto.CaregiverProfileOutputPage{
				Pagination: &domain.Pagination{Limit: 10, CurrentPage: 1, TotalPages: 100},
				Caregivers: []*domain.CaregiverProfile{
					{
						ID:              UUID,
						User:            *user,
						CaregiverNumber: "CG001",
						Consent: domain.ConsentStatus{
							ConsentStatus: enums.ConsentStateAccepted,
						},
					},
				},
			}, nil
		},
		MockRegisterStaffFn: func(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
			return &dto.StaffRegistrationOutput{
				ID:              uuid.New().String(),
				Active:          true,
				StaffNumber:     staff.StaffNumber,
				UserID:          staff.UserID,
				DefaultFacility: *staff.DefaultFacility.ID,
			}, nil
		},
		MockSearchClientUserFn: func(ctx context.Context, searchParameter string) ([]*domain.ClientProfile, error) {
			clientID := uuid.New().String()
			client := &domain.ClientProfile{
				ID:                      &clientID,
				User:                    &domain.User{},
				Active:                  true,
				ClientTypes:             []enums.ClientType{enums.ClientTypePmtct},
				UserID:                  uuid.New().String(),
				TreatmentEnrollmentDate: &time.Time{},
				HealthRecordID:          &clientID,
				ClientCounselled:        false,
				CaregiverID:             &clientID,
			}
			return []*domain.ClientProfile{client}, nil
		},
		MockCompleteOnboardingTourFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockConsentToAClientCaregiverFn: func(ctx context.Context, clientID string, caregiverID string, consent bool) (bool, error) {
			return true, nil
		},
		MockRegisterKenyaEMRPatientsFn: func(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.PatientRegistrationPayload, error) {
			return []*dto.PatientRegistrationPayload{
				{
					MFLCode:   "12345",
					CCCNumber: "12345",
					Name:      gofakeit.Name(),
					DateOfBirth: scalarutils.Date{
						Year:  2000,
						Month: 12,
						Day:   12,
					},
					ClientType:  enums.ClientTypeKenyaEMR,
					PhoneNumber: gofakeit.Phone(),
					EnrollmentDate: scalarutils.Date{
						Year:  2000,
						Month: 12,
						Day:   12,
					},
					BirthDateEstimated: false,
					Gender:             "MALE",
					Counselled:         false,
					NextOfKin: dto.NextOfKinPayload{
						Name:         gofakeit.Name(),
						Contact:      gofakeit.Phone(),
						Relationship: "Brother",
					},
				},
			}, nil
		},
		MockRegisteredFacilityPatientsFn: func(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error) {
			return &dto.PatientSyncResponse{
				MFLCode:  1234,
				Patients: []string{"12345"},
			}, nil
		},
		MockSetUserPINFn: func(ctx context.Context, input dto.PINInput) (bool, error) {
			return true, nil
		},
		MockConsentFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockGetUserProfileFn: func(ctx context.Context, userID string) (*domain.User, error) {
			id := gofakeit.UUID()
			return &domain.User{
				ID:       &id,
				Username: gofakeit.Username(),
				Name:     gofakeit.Gender(),
				Gender:   enumutils.GenderOther,
				Active:   true,
			}, nil
		},
		MockAddClientFHIRIDFn: func(ctx context.Context, input dto.ClientFHIRPayload) error {
			return nil
		},
		MockSearchCaregiverUserFn: func(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
			return []*domain.CaregiverProfile{
				{
					ID:              UUID,
					User:            *user,
					CaregiverNumber: "CG001",
				},
			}, nil
		},
		MockGenerateTemporaryPinFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error) {
			return "1234", nil
		},
		MockRegisterPushTokenFn: func(ctx context.Context, token string) (bool, error) {
			return true, nil
		},
		MockGetClientProfileByCCCNumberFn: func(ctx context.Context, cccNumber string) (*domain.ClientProfile, error) {
			id := gofakeit.UUID()
			return &domain.ClientProfile{
				ID:                      &id,
				User:                    &domain.User{},
				Active:                  true,
				ClientTypes:             []enums.ClientType{enums.ClientTypePmtct},
				UserID:                  UUID,
				TreatmentEnrollmentDate: &time.Time{},
				HealthRecordID:          &id,
				ClientCounselled:        false,
				CaregiverID:             &id,
			}, nil
		},
		MockTransferClientToFacilityFn: func(ctx context.Context, clientID *string, facilityID *string) (bool, error) {
			return true, nil
		},
		MockSetStaffDefaultFacilityFn: func(ctx context.Context, staffID string, facilityID string) (*domain.Facility, error) {
			return &domain.Facility{
				ID:                 &UUID,
				Name:               name,
				Phone:              "1234567890",
				Active:             true,
				County:             gofakeit.BS(),
				Description:        gofakeit.BS(),
				FHIROrganisationID: gofakeit.UUID(),
				Identifier: domain.FacilityIdentifier{
					ID:     UUID,
					Active: true,
					Type:   enums.FacilityIdentifierTypeMFLCode,
					Value:  "1234",
				},
				WorkStationDetails: domain.WorkStationDetails{
					Notifications:   0,
					Surveys:         0,
					Articles:        0,
					Messages:        0,
					ServiceRequests: 0,
				},
			}, nil
		},
		MockSetClientDefaultFacilityFn: func(ctx context.Context, clientID string, facilityID string) (*domain.Facility, error) {
			return &domain.Facility{
				ID:                 &UUID,
				Name:               name,
				Phone:              "1234567890",
				Active:             true,
				County:             gofakeit.BS(),
				Description:        gofakeit.BS(),
				FHIROrganisationID: gofakeit.UUID(),
				Identifier: domain.FacilityIdentifier{
					ID:     UUID,
					Active: true,
					Type:   enums.FacilityIdentifierTypeMFLCode,
					Value:  "1234",
				},
				WorkStationDetails: domain.WorkStationDetails{
					Notifications:   0,
					Surveys:         0,
					Articles:        0,
					Messages:        0,
					ServiceRequests: 0,
				},
			}, nil
		},
		MockAddFacilitiesToStaffProfileFn: func(ctx context.Context, staffID string, facilities []string) (bool, error) {
			return true, nil
		},
		MockConsentToManagingClientFn: func(ctx context.Context, caregiverID string, clientID string, consent bool) (bool, error) {
			return true, nil
		},
		MockGetUserLinkedFacilitiesFn: func(ctx context.Context, userID string, paginationInput dto.PaginationsInput) (*dto.FacilityOutputPage, error) {
			id := gofakeit.UUID()
			return &dto.FacilityOutputPage{
				Pagination: &domain.Pagination{
					Limit:       10,
					CurrentPage: 1,
				},
				Facilities: []*domain.Facility{
					{
						ID:                 &id,
						Name:               "Test Facility",
						Phone:              "",
						Active:             false,
						County:             "",
						Description:        "",
						FHIROrganisationID: "",
					},
				},
			}, nil
		},
		MockAddFacilitiesToClientProfileFn: func(ctx context.Context, clientID string, facilities []string) (bool, error) {
			return true, nil
		},
		MockRemoveFacilitiesFromClientProfileFn: func(ctx context.Context, clientID string, facilities []string) (bool, error) {
			return true, nil
		},
		MockRemoveFacilitiesFromStaffProfileFn: func(ctx context.Context, staffID string, facilities []string) (bool, error) {
			return true, nil
		},
		MockGetCaregiverManagedClientsFn: func(ctx context.Context, caregiverID string, input dto.PaginationsInput) (*dto.ManagedClientOutputPage, error) {
			return &dto.ManagedClientOutputPage{
				Pagination: paginationOutput,
				ManagedClients: []*domain.ManagedClient{
					{
						ClientProfile:    clientProfile,
						CaregiverConsent: enums.ConsentStateAccepted,
						ClientConsent:    enums.ConsentStateAccepted,
					},
				},
			}, nil
		},
	}
}

// Login mocks the login functionality
func (f *UserUseCaseMock) Login(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, bool) {
	return f.MockLoginFn(ctx, input)
}

// FetchContactOrganisations fetches organisations associated with a provided phone number
// Provides the organisation options used during login
func (f *UserUseCaseMock) FetchContactOrganisations(ctx context.Context, phoneNumber string) ([]*domain.Organisation, error) {
	return f.MockFetchContactOrganisationsFn(ctx, phoneNumber)
}

// InviteUser mocks the invite functionality
func (f *UserUseCaseMock) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour, reinvite bool) (bool, error) {
	return f.MockInviteUserFn(ctx, userID, phoneNumber, flavour, reinvite)
}

// SavePin mocks the save pin functionality
func (f *UserUseCaseMock) SavePin(ctx context.Context, input dto.PINInput) (bool, error) {
	return f.MockSavePinFn(ctx, input)
}

// SetNickName is used to mock the implementation to offset or change the user's nickname
func (f *UserUseCaseMock) SetNickName(ctx context.Context, userID string, nickname string) (bool, error) {
	return f.MockSetNickNameFn(ctx, userID, nickname)
}

// RequestPINReset mocks the implementation of requesting pin reset
func (f *UserUseCaseMock) RequestPINReset(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error) {
	return f.MockRequestPINResetFn(ctx, phoneNumber, flavour)
}

// ResetPIN mocks the reset pin functionality
func (f *UserUseCaseMock) ResetPIN(ctx context.Context, input dto.UserResetPinInput) (bool, error) {
	return f.MockResetPINFn(ctx, input)
}

// RefreshToken mocks the implementation for refreshing a token
func (f *UserUseCaseMock) RefreshToken(ctx context.Context, userID string) (*dto.AuthCredentials, error) {
	return f.MockRefreshTokenFn(ctx, userID)
}

// VerifyPIN mocks the implementation for verifying a pin
func (f *UserUseCaseMock) VerifyPIN(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
	return f.MockVerifyPINFn(ctx, userID, flavour, pin)
}

// GetClientCaregiver mocks the implementation for getting the caregiver of a client
func (f *UserUseCaseMock) GetClientCaregiver(ctx context.Context, clientID string) (*domain.Caregiver, error) {
	return f.MockGetClientCaregiverFn(ctx, clientID)
}

// CreateOrUpdateClientCaregiver mocks the implementation for creating or updating the caregiver of a client
func (f *UserUseCaseMock) CreateOrUpdateClientCaregiver(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
	return f.MockCreateOrUpdateClientCaregiverFn(ctx, caregiverInput)
}

// RegisterClient mocks the implementation for registering a client
func (f *UserUseCaseMock) RegisterClient(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error) {
	return f.MockRegisterClientFn(ctx, input)
}

// RefreshGetStreamToken mocks the implementation for generating a new getstream token
func (f *UserUseCaseMock) RefreshGetStreamToken(ctx context.Context, userID string) (*dto.GetStreamToken, error) {
	return f.MockRefreshGetStreamTokenFn(ctx, userID)
}

// RegisterStaff mocks the implementation of registering a staff user
func (f *UserUseCaseMock) RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
	return f.MockRegisterStaffFn(ctx, input)
}

// SearchClientUser mocks the implementation getting the client by CCC number, username or phonenumber
func (f *UserUseCaseMock) SearchClientUser(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
	return f.MockSearchClientUserFn(ctx, CCCNumber)
}

// CompleteOnboardingTour mocks the implementation of completing an onboarding tour
func (f *UserUseCaseMock) CompleteOnboardingTour(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
	return f.MockCompleteOnboardingTourFn(ctx, userID, flavour)
}

// RegisterKenyaEMRPatients mocks the implementation of registering kenyaEMR patients
func (f *UserUseCaseMock) RegisterKenyaEMRPatients(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.PatientRegistrationPayload, error) {
	return f.MockRegisterKenyaEMRPatientsFn(ctx, input)
}

// RegisteredFacilityPatients mocks the implementation of syncing the registered patients
func (f *UserUseCaseMock) RegisteredFacilityPatients(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error) {
	return f.MockRegisteredFacilityPatientsFn(ctx, input)
}

// SetUserPIN mocks the implementation of setting a user pin
func (f *UserUseCaseMock) SetUserPIN(ctx context.Context, input dto.PINInput) (bool, error) {
	return f.MockSetUserPINFn(ctx, input)
}

// SearchStaffUser mocks the implementation of getting staff profile using their staff number.
func (f *UserUseCaseMock) SearchStaffUser(ctx context.Context, searchParameter string) ([]*domain.StaffProfile, error) {
	return f.MockSearchStaffUserFn(ctx, searchParameter)
}

// Consent mocks the implementation of a user withdrawing or offering their consent to the app
func (f *UserUseCaseMock) Consent(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	return f.MockConsentFn(ctx, phoneNumber, flavour)
}

// GetUserProfile returns a user profile given the user ID
func (f *UserUseCaseMock) GetUserProfile(ctx context.Context, userID string) (*domain.User, error) {
	return f.MockGetUserProfileFn(ctx, userID)
}

// AddClientFHIRID updates the client profile with the patient fhir ID from clinical
func (f *UserUseCaseMock) AddClientFHIRID(ctx context.Context, input dto.ClientFHIRPayload) error {
	return f.MockAddClientFHIRIDFn(ctx, input)
}

// GenerateTemporaryPin mocks the implementation of generating temporary pin
func (f *UserUseCaseMock) GenerateTemporaryPin(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error) {
	return f.MockGenerateTemporaryPinFn(ctx, userID, flavour)
}

// RegisterPushToken mocks the implementation for adding a push token to a user's profile
func (f *UserUseCaseMock) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	return f.MockRegisterPushTokenFn(ctx, token)
}

// GetClientProfileByCCCNumber mocks the implementation for getting a client profile by CCC number
func (f *UserUseCaseMock) GetClientProfileByCCCNumber(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
	return f.MockGetClientProfileByCCCNumberFn(ctx, CCCNumber)
}

// DeleteUser mocks the implementation of deleting a user
func (f *UserUseCaseMock) DeleteUser(ctx context.Context, payload *dto.PhoneInput) (bool, error) {
	return f.MockDeleteUserFn(ctx, payload)
}

// TransferClientToFacility mocks the implementation of transferring a client to a facility
func (f *UserUseCaseMock) TransferClientToFacility(ctx context.Context, clientID *string, facilityID *string) (bool, error) {
	return f.MockTransferClientToFacilityFn(ctx, clientID, facilityID)
}

// SetStaffDefaultFacility mocks the implementation of setting a default facility for a staff
func (f *UserUseCaseMock) SetStaffDefaultFacility(ctx context.Context, staffID string, facilityID string) (*domain.Facility, error) {
	return f.MockSetStaffDefaultFacilityFn(ctx, staffID, facilityID)
}

// SetClientDefaultFacility mocks the implementation of setting a default facility for a client
func (f *UserUseCaseMock) SetClientDefaultFacility(ctx context.Context, clientID string, facilityID string) (*domain.Facility, error) {
	return f.MockSetClientDefaultFacilityFn(ctx, clientID, facilityID)
}

// AddFacilitiesToStaffProfile mocks the implementation of adding facilities to a staff profile
func (f *UserUseCaseMock) AddFacilitiesToStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error) {
	return f.MockAddFacilitiesToStaffProfileFn(ctx, staffID, facilities)
}

// GetUserLinkedFacilities mocks the implementation of getting a user's linked facilities
func (f *UserUseCaseMock) GetUserLinkedFacilities(ctx context.Context, userID string, paginationInput dto.PaginationsInput) (*dto.FacilityOutputPage, error) {
	return f.MockGetUserLinkedFacilitiesFn(ctx, userID, paginationInput)
}

// AddFacilitiesToClientProfile mocks the implementation of adding facilities to a client profile
func (f *UserUseCaseMock) AddFacilitiesToClientProfile(ctx context.Context, clientID string, facilities []string) (bool, error) {
	return f.MockAddFacilitiesToClientProfileFn(ctx, clientID, facilities)
}

// RegisterClientAsCaregiver adds a caregiver profile to a client
func (f *UserUseCaseMock) RegisterClientAsCaregiver(ctx context.Context, clientID string, caregiverNumber string) (*domain.CaregiverProfile, error) {
	return f.MockRegisterClientAsCaregiverFn(ctx, clientID, caregiverNumber)
}

// RegisterCaregiver is used to register a caregiver
func (f *UserUseCaseMock) RegisterCaregiver(ctx context.Context, input dto.CaregiverInput) (*domain.CaregiverProfile, error) {
	return f.MockRegisterCaregiver(ctx, input)
}

// SearchCaregiverUser mocks the implementation of searching caregiver profile using their caregiver number.
func (f *UserUseCaseMock) SearchCaregiverUser(ctx context.Context, searchParameter string) ([]*domain.CaregiverProfile, error) {
	return f.MockSearchCaregiverUserFn(ctx, searchParameter)
}

// RemoveFacilitiesFromClientProfile mocks the implementation of removing facilities from a client profile
func (f *UserUseCaseMock) RemoveFacilitiesFromClientProfile(ctx context.Context, clientID string, facilities []string) (bool, error) {
	return f.MockRemoveFacilitiesFromClientProfileFn(ctx, clientID, facilities)
}

// AssignCaregiver mocks the implementation of adding a caregiver to a client
func (f *UserUseCaseMock) AssignCaregiver(ctx context.Context, input dto.ClientCaregiverInput) (bool, error) {
	return f.MockAssignCaregiverFn(ctx, input)
}

// RemoveFacilitiesFromStaffProfile mocks the implementation of removing facilities from a staff profile
func (f *UserUseCaseMock) RemoveFacilitiesFromStaffProfile(ctx context.Context, staffID string, facilities []string) (bool, error) {
	return f.MockRemoveFacilitiesFromStaffProfileFn(ctx, staffID, facilities)
}

// GetCaregiverManagedClients mocks the implementation of getting caregiver's managed clients
func (f *UserUseCaseMock) GetCaregiverManagedClients(ctx context.Context, caregiverID string, input dto.PaginationsInput) (*dto.ManagedClientOutputPage, error) {
	return f.MockGetCaregiverManagedClientsFn(ctx, caregiverID, input)
}

// ListClientsCaregivers mocks the implementation of listing a client's caregivers
func (f *UserUseCaseMock) ListClientsCaregivers(ctx context.Context, clientID string, pagination *dto.PaginationsInput) (*dto.CaregiverProfileOutputPage, error) {
	return f.MockListClientsCaregiversFn(ctx, clientID, pagination)
}

// ConsentToAClientCaregiver mocks the implementation of a client acknowledging to having a certain caregiver assigned to them.
func (f *UserUseCaseMock) ConsentToAClientCaregiver(ctx context.Context, clientID string, caregiverID string, consent bool) (bool, error) {
	return f.MockConsentToAClientCaregiverFn(ctx, clientID, caregiverID, consent)
}

// ConsentToManagingClient mock the implementation of a caregiver acknowledging or offering their consent to act on behalf of the client.
func (f *UserUseCaseMock) ConsentToManagingClient(ctx context.Context, caregiverID string, clientID string, consent bool) (bool, error) {
	return f.MockConsentToManagingClientFn(ctx, caregiverID, clientID, consent)
}
