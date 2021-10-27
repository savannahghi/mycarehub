package mock

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

// GormMock struct implements mocks of `gorm's`internal methods.
//
// This mock struct should be separate from our own internal methods.
type GormMock struct {
	GetOrCreateFacilityFn          func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error)
	RetrieveFacilityFn             func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error)
	RetrieveFacilityByMFLCodeFn    func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error)
	GetFacilitiesFn                func(ctx context.Context) ([]gorm.Facility, error)
	DeleteFacilityFn               func(ctx context.Context, mfl_code string) (bool, error)
	CollectMetricsFn               func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error)
	SavePinFn                      func(ctx context.Context, pinData *gorm.PINData) (bool, error)
	GetUserPINByUserIDFn           func(ctx context.Context, userID string) (*gorm.PINData, error)
	GetUserProfileByUserIDFn       func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.User, error)
	GetOrCreateStaffUserserFn      func(ctx context.Context, user *gorm.User, staff *gorm.StaffProfile) (*gorm.StaffUserProfile, error)
	RegisterClientFn               func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error)
	AddIdentifierFn                func(ctx context.Context, identifier *gorm.Identifier) (*gorm.Identifier, error)
	GetClientProfileByClientIDFn   func(ctx context.Context, clientID string) (*gorm.ClientProfile, error)
	GetStaffProfileFn              func(ctx context.Context, staffNumber string) (*gorm.StaffProfile, error)
	GetStaffProfileByStaffIDFn     func(ctx context.Context, staffProfileID string) (*gorm.StaffUserProfile, error)
	GetStaffProfileByStaffNumberFn func(ctx context.Context, staffNumber string) (*gorm.StaffUserProfile, error)

	//Updates
	UpdateUserLastSuccessfulLoginFn func(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateUserLastFailedLoginFn     func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateUserFailedLoginCountFn    func(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error
	UpdateUserNextAllowedLoginFn    func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateStaffUserFn               func(ctx context.Context, userID string, user *gorm.User, staff *gorm.StaffProfile) (bool, error)
	TransferClientFn                func(
		ctx context.Context,
		clientID string,
		originFacilityID string,
		destinationFacilityID string,
		reason enums.TransferReason,
		notes string,
	) (bool, error)
}

// NewGormMock initializes a new instance of `GormMock` then mocking the case of success.
func NewGormMock() *GormMock {
	return &GormMock{
		RegisterClientFn: func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
			return &gorm.ClientUserProfile{
				User: &gorm.User{
					FirstName:   "FirstName",
					LastName:    "Last Name",
					Username:    "User Name",
					MiddleName:  userInput.MiddleName,
					DisplayName: "Display Name",
					Gender:      enumutils.GenderMale,
				},
				Client: &gorm.ClientProfile{
					ClientType: enums.ClientTypeOvc,
				},
			}, nil
		},

		AddIdentifierFn: func(ctx context.Context, identifier *gorm.Identifier) (*gorm.Identifier, error) {
			return &gorm.Identifier{
				ClientID:        identifier.ClientID,
				IdentifierType:  enums.IdentifierTypeCCC,
				IdentifierUse:   enums.IdentifierUseOfficial,
				IdentifierValue: "Just a random value",
				Description:     "Random description",
			}, nil
		},

		GetClientProfileByClientIDFn: func(ctx context.Context, clientID string) (*gorm.ClientProfile, error) {
			ID := uuid.New().String()
			return &gorm.ClientProfile{
				ID:     &clientID,
				UserID: &ID,
			}, nil
		},

		GetOrCreateFacilityFn: func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
			id := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &id,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},

		RetrieveFacilityFn: func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},
		GetFacilitiesFn: func(ctx context.Context) ([]gorm.Facility, error) {
			var facilities []gorm.Facility
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			facilities = append(facilities, gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			})
			return facilities, nil
		},

		DeleteFacilityFn: func(ctx context.Context, mfl_code string) (bool, error) {
			return true, nil
		},

		CollectMetricsFn: func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
			now := time.Now()
			metricID := uuid.New().String()
			return &gorm.Metric{
				MetricID:  &metricID,
				Type:      enums.EngagementMetrics,
				Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
				Timestamp: now,
				UID:       ksuid.New().String(),
			}, nil
		},

		RetrieveFacilityByMFLCodeFn: func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
			facilityID := uuid.New().String()
			name := "Kanairo One"
			code := "KN001"
			county := "Kanairo"
			description := "This is just for mocking"
			return &gorm.Facility{
				FacilityID:  &facilityID,
				Name:        name,
				Code:        code,
				Active:      strconv.FormatBool(true),
				County:      county,
				Description: description,
			}, nil
		},

		SavePinFn: func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
			return true, nil
		},

		GetUserProfileByUserIDFn: func(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.User, error) {
			id := uuid.New().String()
			usercontact := &gorm.Contact{
				ContactID: &id,
				Type:      "test",
				Contact:   "test",
				Active:    true,
				OptedIn:   true,
			}
			time := time.Now()
			return &gorm.User{
				Base:                gorm.Base{},
				UserID:              &id,
				Username:            "test",
				DisplayName:         "test",
				FirstName:           "test",
				MiddleName:          "test",
				LastName:            "test",
				Flavour:             "test",
				UserType:            "test",
				Gender:              "test",
				Active:              false,
				Contacts:            []gorm.Contact{*usercontact},
				Languages:           []string{"en"},
				PushTokens:          []string{"push-token"},
				LastSuccessfulLogin: &time,
				LastFailedLogin:     &time,
				FailedLoginCount:    "test",
				NextAllowedLogin:    &time,
				TermsAccepted:       false,
				AcceptedTermsID:     "test",
			}, nil
		},

		GetUserPINByUserIDFn: func(ctx context.Context, userID string) (*gorm.PINData, error) {
			return &gorm.PINData{
				UserID:    userID,
				HashedPIN: "mbzcbvhbxchjbvhdbvhhjdfskgbfhas832y38hjsdnfkjbh73y73y72",
				ValidFrom: time.Now(),
				ValidTo:   time.Now(),
				Flavour:   "CONSUMER",
				IsValid:   true,
				Salt:      "test-salt",
			}, nil
		},

		GetStaffProfileFn: func(ctx context.Context, staffNumber string) (*gorm.StaffProfile, error) {
			testUID := ksuid.New().String()
			address := &gorm.Addresses{
				Type:           "test",
				Text:           "test",
				Country:        "test",
				PostalCode:     "test",
				County:         "test",
				Active:         false,
				StaffProfileID: new(string),
			}
			return &gorm.StaffProfile{
				StaffProfileID:    &testUID,
				UserID:            &testUID,
				User:              gorm.User{},
				StaffNumber:       "s100",
				DefaultFacilityID: &testUID,
				Addresses:         []*gorm.Addresses{address},
			}, nil
		},

		UpdateUserLastSuccessfulLoginFn: func(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateUserLastFailedLoginFn: func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateUserFailedLoginCountFn: func(ctx context.Context, userID, failedLoginCount string, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateUserNextAllowedLoginFn: func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateStaffUserFn: func(ctx context.Context, userID string, user *gorm.User, staff *gorm.StaffProfile) (bool, error) {
			return true, nil
		},

		GetOrCreateStaffUserserFn: func(ctx context.Context, user *gorm.User, staff *gorm.StaffProfile) (*gorm.StaffUserProfile, error) {
			ID := uuid.New().String()
			testTime := time.Now()
			return &gorm.StaffUserProfile{
				User: &gorm.User{
					UserID:              &ID,
					Username:            "test",
					DisplayName:         "test",
					FirstName:           "test",
					MiddleName:          "test",
					LastName:            "test",
					Active:              true,
					LastSuccessfulLogin: &testTime,
					LastFailedLogin:     &testTime,
					NextAllowedLogin:    &testTime,
					FailedLoginCount:    "0",
					TermsAccepted:       true,
					AcceptedTermsID:     ID,
				},
				Staff: &gorm.StaffProfile{
					StaffProfileID:    &ID,
					UserID:            &ID,
					StaffNumber:       "s123",
					DefaultFacilityID: &ID,
					Addresses: []*gorm.Addresses{
						{
							AddressesID: &ID,
							Type:        enums.AddressesTypePhysical,
							Text:        "test",
							Country:     enums.CountryTypeKenya,
							PostalCode:  "test code",
							County:      "test",
							Active:      true,
						},
					},
				},
			}, nil
		},
		GetStaffProfileByStaffIDFn: func(ctx context.Context, staffProfileID string) (*gorm.StaffUserProfile, error) {
			ID := uuid.New().String()
			testTime := time.Now()
			return &gorm.StaffUserProfile{
				User: &gorm.User{
					UserID:              &ID,
					Username:            "test",
					DisplayName:         "test",
					FirstName:           "test",
					MiddleName:          "test",
					LastName:            "test",
					Active:              true,
					LastSuccessfulLogin: &testTime,
					LastFailedLogin:     &testTime,
					NextAllowedLogin:    &testTime,
					FailedLoginCount:    "0",
					TermsAccepted:       true,
					AcceptedTermsID:     ID,
				},
				Staff: &gorm.StaffProfile{
					StaffProfileID:    &ID,
					UserID:            &ID,
					StaffNumber:       "s123",
					DefaultFacilityID: &ID,
					Addresses: []*gorm.Addresses{
						{
							AddressesID: &ID,
							Type:        enums.AddressesTypePhysical,
							Text:        "test",
							Country:     enums.CountryTypeKenya,
							PostalCode:  "test code",
							County:      "test",
							Active:      true,
						},
					},
				},
			}, nil
		},
		GetStaffProfileByStaffNumberFn: func(ctx context.Context, staffNumber string) (*gorm.StaffUserProfile, error) {
			ID := uuid.New().String()
			testTime := time.Now()
			roles := []string{enums.RolesTypeCanInviteClient.String()}
			languages := []string{string(enumutils.LanguageEn)}
			return &gorm.StaffUserProfile{
				User: &gorm.User{
					UserID:              &ID,
					Username:            "test",
					DisplayName:         "test",
					FirstName:           "test",
					MiddleName:          "test",
					LastName:            "test",
					Active:              true,
					LastSuccessfulLogin: &testTime,
					LastFailedLogin:     &testTime,
					NextAllowedLogin:    &testTime,
					FailedLoginCount:    "0",
					TermsAccepted:       true,
					AcceptedTermsID:     ID,
					Languages:           languages,
				},
				Staff: &gorm.StaffProfile{
					StaffProfileID:    &ID,
					UserID:            &ID,
					StaffNumber:       "s123",
					DefaultFacilityID: &ID,
					Addresses: []*gorm.Addresses{
						{
							AddressesID: &ID,
							Type:        enums.AddressesTypePhysical,
							Text:        "test",
							Country:     enums.CountryTypeKenya,
							PostalCode:  "test code",
							County:      "test",
							Active:      true,
						},
					},
					Roles: roles,
				},
			}, nil
		},

		TransferClientFn: func(ctx context.Context, clientID, originFacilityID, destinationFacilityID string, reason enums.TransferReason, notes string) (bool, error) {
			return true, nil
		},
	}
}

// GetOrCreateFacility mocks the implementation of `gorm's` GetOrCreateFacility method.
func (gm *GormMock) GetOrCreateFacility(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
	return gm.GetOrCreateFacilityFn(ctx, facility)
}

// RetrieveFacility mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
	return gm.RetrieveFacilityFn(ctx, id, isActive)
}

// RetrieveFacilityByMFLCode mocks the implementation of `gorm's` RetrieveFacility method.
func (gm *GormMock) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
	return gm.RetrieveFacilityByMFLCodeFn(ctx, MFLCode, isActive)
}

// GetFacilities mocks the implementation of `gorm's` GetFacilities method.
func (gm *GormMock) GetFacilities(ctx context.Context) ([]gorm.Facility, error) {
	return gm.GetFacilitiesFn(ctx)
}

// DeleteFacility mocks the implementation of  DeleteFacility method.
func (gm *GormMock) DeleteFacility(ctx context.Context, mflcode string) (bool, error) {
	return gm.DeleteFacilityFn(ctx, mflcode)
}

// CollectMetrics mocks the implementation of  CollectMetrics method.
func (gm *GormMock) CollectMetrics(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
	return gm.CollectMetricsFn(ctx, metrics)
}

//SavePin mocks the implementation of SetUserPIN method
func (gm *GormMock) SavePin(ctx context.Context, pinData *gorm.PINData) (bool, error) {
	return gm.SavePinFn(ctx, pinData)
}

// GetUserPINByUserID ...
func (gm *GormMock) GetUserPINByUserID(ctx context.Context, userID string) (*gorm.PINData, error) {
	return gm.GetUserPINByUserIDFn(ctx, userID)
}

// GetUserProfileByUserID gets user profile by user ID
func (gm *GormMock) GetUserProfileByUserID(ctx context.Context, userID string, flavour feedlib.Flavour) (*gorm.User, error) {
	return gm.GetUserProfileByUserIDFn(ctx, userID, flavour)
}

//UpdateUserLastSuccessfulLogin updates the user's last successful login time
func (gm *GormMock) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error {
	return gm.UpdateUserLastSuccessfulLoginFn(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin updates the user's last failed login
func (gm *GormMock) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error {
	return gm.UpdateUserLastFailedLoginFn(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount updates the users failed login count
func (gm *GormMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error {
	return gm.UpdateUserFailedLoginCountFn(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin updates the user's next allowed login time
func (gm *GormMock) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error {
	return gm.UpdateUserNextAllowedLoginFn(ctx, userID, nextAllowedLoginTime, flavour)
}

// GetOrCreateStaffUser mocks the implementation of  GetOrCreateStaffUser method.
func (gm *GormMock) GetOrCreateStaffUser(ctx context.Context, user *gorm.User, staff *gorm.StaffProfile) (*gorm.StaffUserProfile, error) {
	return gm.GetOrCreateStaffUserserFn(ctx, user, staff)
}

// UpdateStaffUserProfile mocks the implementation of  UpdateStaffUserProfile method.
func (gm *GormMock) UpdateStaffUserProfile(ctx context.Context, userID string, user *gorm.User, staff *gorm.StaffProfile) (bool, error) {
	return gm.UpdateStaffUserFn(ctx, userID, user, staff)
}

// RegisterClient mocks the implementation of RegisterClient method
func (gm *GormMock) RegisterClient(
	ctx context.Context,
	userInput *gorm.User,
	clientInput *gorm.ClientProfile,
) (*gorm.ClientUserProfile, error) {
	return gm.RegisterClientFn(ctx, userInput, clientInput)
}

// AddIdentifier mocks the `AddIdentifier` implementation
func (gm *GormMock) AddIdentifier(
	ctx context.Context,
	identifier *gorm.Identifier,
) (*gorm.Identifier, error) {
	return gm.AddIdentifierFn(ctx, identifier)
}

// GetClientProfileByClientID mocks the method that fetches a client profile by the ID
func (gm *GormMock) GetClientProfileByClientID(ctx context.Context, clientID string) (*gorm.ClientProfile, error) {
	return gm.GetClientProfileByClientIDFn(ctx, clientID)
}

// TransferClient mocks the implementation of  TransferClient method
func (gm *GormMock) TransferClient(
	ctx context.Context,
	clientID string,
	originFacilityID string,
	destinationFacilityID string,
	reason enums.TransferReason,
	notes string,
) (bool, error) {
	return gm.TransferClientFn(ctx, clientID, originFacilityID, destinationFacilityID, reason, notes)
}

// GetStaffProfileByStaffID mocks the  GetStaffProfileByStaffID method.
func (gm *GormMock) GetStaffProfileByStaffID(ctx context.Context, staffProfileID string) (*gorm.StaffUserProfile, error) {
	return gm.GetStaffProfileByStaffIDFn(ctx, staffProfileID)
}

// GetStaffProfileByStaffNumber mocks the  GetStaffProfileByStaffNumber method.
func (gm *GormMock) GetStaffProfileByStaffNumber(ctx context.Context, staffNumber string) (*gorm.StaffUserProfile, error) {
	return gm.GetStaffProfileByStaffNumberFn(ctx, staffNumber)
}
