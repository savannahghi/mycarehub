package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

// UserUseCaseMock mocks the implementation of usecase methods.
type UserUseCaseMock struct {
	MockLoginFn                         func(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, error)
	MockInviteUserFn                    func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error)
	MockSavePinFn                       func(ctx context.Context, input dto.PINInput) (bool, error)
	MockVerifyLoginPINFn                func(ctx context.Context, userProfile *domain.User, pin string, flavour feedlib.Flavour) (bool, error)
	MockSetNickNameFn                   func(ctx context.Context, userID string, nickname string) (bool, error)
	MockRequestPINResetFn               func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (string, error)
	MockResetPINFn                      func(ctx context.Context, input dto.UserResetPinInput) (bool, error)
	MockRefreshTokenFn                  func(ctx context.Context, userID string) (*domain.AuthCredentials, error)
	MockVerifyPINFn                     func(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error)
	MockGetClientCaregiverFn            func(ctx context.Context, clientID string) (*domain.Caregiver, error)
	MockCreateOrUpdateClientCaregiverFn func(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error)
	MockRegisterClientFn                func(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error)
	MockRefreshGetStreamTokenFn         func(ctx context.Context, userID string) (*domain.GetStreamToken, error)
	MockRegisterStaffFn                 func(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error)
	MockSearchClientsByCCCNumberFn      func(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error)
	MockCompleteOnboardingTourFn        func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error)
	MockRegisterKenyaEMRPatientsFn      func(ctx context.Context, input []*dto.PatientRegistrationPayload) ([]*dto.PatientRegistrationPayload, error)
	MockRegisteredFacilityPatientsFn    func(ctx context.Context, input dto.PatientSyncPayload) (*dto.PatientSyncResponse, error)
	MockSetUserPINFn                    func(ctx context.Context, input dto.PINInput) (bool, error)
	MockSearchStaffByStaffNumberFn      func(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error)
	MockConsentFn                       func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour, active bool) (bool, error)
	MockGetUserProfileFn                func(ctx context.Context, userID string) (*domain.User, error)
	MockAddClientFHIRIDFn               func(ctx context.Context, input dto.ClientFHIRPayload) error
	MockGenerateTemporaryPinFn          func(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error)
	MockRegisterPushTokenFn             func(ctx context.Context, token string) (bool, error)
	MockFetchNotificationsFn            func(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error)
}

// NewUserUseCaseMock creates in initializes create type mocks
func NewUserUseCaseMock() *UserUseCaseMock {
	var UUID = uuid.New().String()
	caregiver := &domain.Caregiver{
		ID:            UUID,
		FirstName:     gofakeit.FirstName(),
		LastName:      gofakeit.LastName(),
		PhoneNumber:   gofakeit.Phone(),
		CaregiverType: enums.CaregiverTypeFather,
	}

	staff := &domain.StaffProfile{
		ID:                &UUID,
		User:              &domain.User{},
		UserID:            uuid.New().String(),
		Active:            true,
		StaffNumber:       "test-staff-101",
		Facilities:        []domain.Facility{},
		DefaultFacilityID: uuid.New().String(),
	}

	return &UserUseCaseMock{

		MockLoginFn: func(ctx context.Context, phoneNumber, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, error) {
			ID := uuid.New().String()
			time := time.Now()
			resp := &domain.Response{
				Client:          &domain.ClientProfile{ID: &ID, User: &domain.User{ID: &ID, Username: gofakeit.Username(), TermsAccepted: true, Active: true, NextAllowedLogin: &time, FailedLoginCount: 1}},
				Staff:           &domain.StaffProfile{},
				AuthCredentials: domain.AuthCredentials{RefreshToken: gofakeit.HipsterSentence(15), IDToken: gofakeit.BeerAlcohol(), ExpiresIn: gofakeit.BeerHop()},
				GetStreamToken:  "",
			}
			return &domain.LoginResponse{
				Response: resp,
				Attempts: 10,
				Message:  "Success",
				Code:     10,
			}, nil
		},
		MockFetchNotificationsFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error) {
			return &domain.NotificationsPage{
				Notifications: []*domain.Notification{
					{
						ID:         UUID,
						Title:      "New Teleconsult",
						Body:       "Teleconsult with Doctor Who at the Tardis",
						Type:       "TELECONSULT",
						IsRead:     false,
						UserID:     &UUID,
						FacilityID: &UUID,
						Flavour:    flavour,
					},
				},
				Pagination: domain.Pagination{},
			}, nil
		},
		MockInviteUserFn: func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
			return true, nil
		},
		MockSavePinFn: func(ctx context.Context, input dto.PINInput) (bool, error) {
			return true, nil
		},
		MockVerifyLoginPINFn: func(ctx context.Context, userProfile *domain.User, pin string, flavour feedlib.Flavour) (bool, error) {
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
		MockRefreshTokenFn: func(ctx context.Context, userID string) (*domain.AuthCredentials, error) {
			return &domain.AuthCredentials{
				RefreshToken: uuid.New().String(),
				ExpiresIn:    "3600",
				IDToken:      uuid.New().String(),
			}, nil
		},
		MockVerifyPINFn: func(ctx context.Context, userID string, flavour feedlib.Flavour, pin string) (bool, error) {
			return true, nil
		},
		MockSearchStaffByStaffNumberFn: func(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
			return []*domain.StaffProfile{staff}, nil
		},
		MockGetClientCaregiverFn: func(ctx context.Context, clientID string) (*domain.Caregiver, error) {
			return caregiver, nil
		},
		MockCreateOrUpdateClientCaregiverFn: func(ctx context.Context, caregiverInput *dto.CaregiverInput) (bool, error) {
			return true, nil
		},
		MockRegisterClientFn: func(ctx context.Context, input *dto.ClientRegistrationInput) (*dto.ClientRegistrationOutput, error) {
			return &dto.ClientRegistrationOutput{
				ID: uuid.New().String(),
			}, nil
		},
		MockRefreshGetStreamTokenFn: func(ctx context.Context, userID string) (*domain.GetStreamToken, error) {
			return &domain.GetStreamToken{
				Token: uuid.New().String(),
			}, nil
		},
		MockRegisterStaffFn: func(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
			return &dto.StaffRegistrationOutput{
				ID: uuid.New().String(),
			}, nil
		},
		MockSearchClientsByCCCNumberFn: func(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
			clientID := uuid.New().String()
			client := &domain.ClientProfile{
				ID:                      &clientID,
				User:                    &domain.User{},
				Active:                  true,
				ClientType:              "PMTCT",
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
					ClientType:  "KenyaEMR",
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
		MockConsentFn: func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour, active bool) (bool, error) {
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
		MockGenerateTemporaryPinFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error) {
			return "1234", nil
		},
		MockRegisterPushTokenFn: func(ctx context.Context, token string) (bool, error) {
			return true, nil
		},
	}
}

// Login mocks the login functionality
func (f *UserUseCaseMock) Login(ctx context.Context, phoneNumber string, pin string, flavour feedlib.Flavour) (*domain.LoginResponse, error) {
	return f.MockLoginFn(ctx, phoneNumber, pin, flavour)
}

// InviteUser mocks the invite functionality
func (f *UserUseCaseMock) InviteUser(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
	return f.MockInviteUserFn(ctx, userID, phoneNumber, flavour)
}

// SavePin mocks the save pin functionality
func (f *UserUseCaseMock) SavePin(ctx context.Context, input dto.PINInput) (bool, error) {
	return f.MockSavePinFn(ctx, input)
}

// VerifyLoginPIN mocks the verify pin functionality
func (f *UserUseCaseMock) VerifyLoginPIN(ctx context.Context, userProfile *domain.User, pin string, flavour feedlib.Flavour) (bool, error) {
	return f.MockVerifyLoginPINFn(ctx, userProfile, pin, flavour)
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
func (f *UserUseCaseMock) RefreshToken(ctx context.Context, userID string) (*domain.AuthCredentials, error) {
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
func (f *UserUseCaseMock) RefreshGetStreamToken(ctx context.Context, userID string) (*domain.GetStreamToken, error) {
	return f.MockRefreshGetStreamTokenFn(ctx, userID)
}

// RegisterStaff mocks the implementation of registering a staff user
func (f *UserUseCaseMock) RegisterStaff(ctx context.Context, input dto.StaffRegistrationInput) (*dto.StaffRegistrationOutput, error) {
	return f.MockRegisterStaffFn(ctx, input)
}

// SearchClientsByCCCNumber mocks the implementation getting the client by CCC number
func (f *UserUseCaseMock) SearchClientsByCCCNumber(ctx context.Context, CCCNumber string) ([]*domain.ClientProfile, error) {
	return f.MockSearchClientsByCCCNumberFn(ctx, CCCNumber)
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

// SearchStaffByStaffNumber mocks the implementation of getting staff profile using their staff number.
func (f *UserUseCaseMock) SearchStaffByStaffNumber(ctx context.Context, staffNumber string) ([]*domain.StaffProfile, error) {
	return f.MockSearchStaffByStaffNumberFn(ctx, staffNumber)
}

// Consent mocks the implementation of a user withdrawing or offering their consent to the app
func (f *UserUseCaseMock) Consent(ctx context.Context, phoneNumber string, flavour feedlib.Flavour, active bool) (bool, error) {
	return f.MockConsentFn(ctx, phoneNumber, flavour, active)
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

// FetchNotifications retrieves a users notifications
func (f *UserUseCaseMock) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error) {
	return f.MockFetchNotificationsFn(ctx, userID, flavour, paginationInput)
}
