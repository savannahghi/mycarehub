package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm/mock"
	"github.com/segmentio/ksuid"
	"gorm.io/datatypes"
)

func TestOnboardingDb_CreateFacility(t *testing.T) {
	ctx := context.Background()
	name := "Kanairo One"
	code := "KN001"
	county := "Kanairo"
	description := "This is just for mocking"
	type args struct {
		ctx      context.Context
		facility *dto.FacilityInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid payload",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
					Name:        name,
					Code:        code,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx: ctx,
				facility: &dto.FacilityInput{
					Name:        name,
					Active:      true,
					County:      county,
					Description: description,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if tt.name == "sad case - facility code not defined" {
				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return facility, nil
				}
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facility to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facility not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_CollectMetrics_Unittest(t *testing.T) {
	ctx := context.Background()

	metric := &dto.MetricInput{
		Type:      enums.EngagementMetrics,
		Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
		Timestamp: time.Now(),
		UID:       ksuid.New().String(),
	}

	invalidMetric := &dto.MetricInput{
		Type:      "",
		Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
		Timestamp: time.Now(),
		UID:       ksuid.New().String(),
	}

	type args struct {
		ctx    context.Context
		metric *dto.MetricInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    ctx,
				metric: metric,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:    ctx,
				metric: invalidMetric,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Happy case" {
				fakeGorm.CollectMetricsFn = func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
					now := time.Now()
					metricID := uuid.New().String()
					return &gorm.Metric{
						MetricID:  &metricID,
						Type:      enums.EngagementMetrics,
						Payload:   datatypes.JSON([]byte(`{"who": "test user", "keyword": "suicidal"}`)),
						Timestamp: now,
						UID:       ksuid.New().String(),
					}, nil
				}
			}

			if tt.name == "Sad case" {
				fakeGorm.CollectMetricsFn = func(ctx context.Context, metrics *gorm.Metric) (*gorm.Metric, error) {
					return nil, fmt.Errorf("an error occurred while collecting metrics")
				}
			}

			_, err := d.CollectMetrics(tt.args.ctx, tt.args.metric)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.CollectMetrics() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestOnboardingDb_SetUserPIN(t *testing.T) {
	ctx := context.Background()

	validPINDataINput := &domain.UserPIN{
		UserID:    ksuid.New().String(),
		HashedPIN: "test-Pin",
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Flavour:   "CONSUMER",
		IsValid:   false,
		Salt:      "salt",
	}

	invalidPINDataINput := &domain.UserPIN{
		UserID:    "",
		HashedPIN: "test-Pin",
		ValidFrom: time.Time{},
		ValidTo:   time.Time{},
		Flavour:   "CONSUMER",
		IsValid:   false,
		Salt:      "salt",
	}
	type args struct {
		ctx     context.Context
		pinData *domain.UserPIN
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				pinData: validPINDataINput,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				pinData: invalidPINDataINput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			_, err := d.SavePin(tt.args.ctx, tt.args.pinData)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.SetUserPIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestOnboardingDb_GetOrCreateStaffUser(t *testing.T) {
	ctx := context.Background()

	testFacilityID := uuid.New().String()
	testUserID := uuid.New().String()
	testTime := time.Now()
	testID := uuid.New().String()
	rolesInput := []enums.RolesType{enums.RolesTypeCanInviteClient}

	type args struct {
		ctx   context.Context
		user  *dto.UserInput
		staff *dto.StaffProfileInput
	}

	contactInput := &dto.ContactInput{
		Type:    enums.PhoneContact,
		Contact: "+254700000000",
		Active:  true,
		OptedIn: true,
	}

	userInput := &dto.UserInput{
		Username:    "test",
		DisplayName: "test",
		FirstName:   "test",
		MiddleName:  "test",
		LastName:    "test",
		Gender:      enumutils.GenderMale,
		UserType:    enums.HealthcareWorkerUser,
		Contacts:    []*dto.ContactInput{contactInput},
		Languages:   []enumutils.Language{enumutils.LanguageEn},
		Flavour:     feedlib.FlavourPro,
	}

	staffInput := &dto.StaffProfileInput{
		StaffNumber:       "s123",
		DefaultFacilityID: &testFacilityID,
		Addresses: []*dto.AddressesInput{
			{
				Type:       enums.AddressesTypePhysical,
				Text:       "test",
				Country:    enums.CountryTypeKenya,
				PostalCode: "test code",
				County:     enums.CountyTypeBaringo,
				Active:     true,
			},
		},
		Roles: rolesInput,
	}
	staffNoFacilityIDInput := &dto.StaffProfileInput{
		StaffNumber: "s123",
	}

	tests := []struct {
		name    string
		args    args
		want    *domain.StaffUserProfile
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:   ctx,
				user:  userInput,
				staff: staffInput,
			},
			wantErr: false,
		},
		{
			name: "invalid: missing facility",
			args: args{
				ctx:   ctx,
				user:  userInput,
				staff: staffNoFacilityIDInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		var fakeGorm = gormMock.NewGormMock()
		d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy case" {

				fakeGorm.GetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return &gorm.Facility{
						FacilityID:  &testFacilityID,
						Name:        "test",
						Code:        "f1234",
						Active:      "true",
						County:      "test",
						Description: "test description",
					}, nil
				}
				fakeGorm.GetOrCreateStaffUserserFn = func(ctx context.Context, user *gorm.User, staff *gorm.StaffProfile) (*gorm.StaffUserProfile, error) {
					contact := gorm.Contact{
						ContactID: &testID,
						Type:      enums.PhoneContact,
						Contact:   "+254700000000",
						Active:    true,
						OptedIn:   true,
					}
					return &gorm.StaffUserProfile{
						User: &gorm.User{
							UserID:              &testUserID,
							Username:            "test",
							DisplayName:         "test",
							FirstName:           "test",
							MiddleName:          "test",
							LastName:            "test",
							Gender:              enumutils.GenderMale,
							Active:              true,
							Contacts:            []gorm.Contact{contact},
							UserType:            enums.HealthcareWorkerUser,
							Languages:           pq.StringArray{"EN", "SW"},
							LastSuccessfulLogin: &testTime,
							LastFailedLogin:     &testTime,
							NextAllowedLogin:    &testTime,
							FailedLoginCount:    "0",
							TermsAccepted:       true,
							AcceptedTermsID:     testID,
							Flavour:             feedlib.FlavourPro,
						},
						Staff: &gorm.StaffProfile{
							StaffProfileID:    &testID,
							UserID:            &testUserID,
							StaffNumber:       "s123",
							DefaultFacilityID: &testFacilityID,
							Addresses: []*gorm.Addresses{
								{
									AddressesID: &testID,
									Type:        enums.AddressesTypePhysical,
									Text:        "test",
									Country:     enums.CountryTypeKenya,
									PostalCode:  "test code",
									County:      enums.CountyTypeBaringo,
									Active:      true,
								},
							},
							Roles: []string{enums.RolesTypeCanInviteClient.String()},
						},
					}, nil
				}
			}

			if tt.name == "invalid: missing facility" {
				fakeGorm.GetOrCreateStaffUserserFn = func(ctx context.Context, user *gorm.User, staff *gorm.StaffProfile) (*gorm.StaffUserProfile, error) {
					return nil, fmt.Errorf("test error")
				}
			}

			_, err := d.GetOrCreateStaffUser(tt.args.ctx, tt.args.user, tt.args.staff)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetOrCreateStaffUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestOnboardingDb_RegisterClient(t *testing.T) {
	ctx := context.Background()
	userInput := &dto.UserInput{
		FirstName:   "John",
		LastName:    "Joe",
		Username:    "Jontez",
		MiddleName:  "Johnny",
		DisplayName: "jo",
		Gender:      enumutils.GenderMale,
	}

	clientInput := dto.ClientProfileInput{
		ClientType: enums.ClientTypeOvc,
	}
	type args struct {
		ctx         context.Context
		userInput   *dto.UserInput
		clientInput *dto.ClientProfileInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: &clientInput,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Fail to create user",
			args: args{
				ctx:         ctx,
				userInput:   userInput,
				clientInput: &clientInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case: Fail to create user" {
				fakeGorm.RegisterClientFn = func(ctx context.Context, userInput *gorm.User, clientInput *gorm.ClientProfile) (*gorm.ClientUserProfile, error) {
					return nil, fmt.Errorf("failed to create a client user")
				}
			}

			got, err := d.RegisterClient(tt.args.ctx, tt.args.userInput, tt.args.clientInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got :%v", got)
			}
		})
	}
}

func TestOnboardingDb_AddIdentifier(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		clientID  string
		idType    enums.IdentifierType
		idValue   string
		isPrimary bool
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Identifier
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully add identifier",
			args: args{
				ctx:       ctx,
				clientID:  "12345",
				idType:    enums.IdentifierTypeCCC,
				idValue:   "1224",
				isPrimary: true,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to add identifier",
			args: args{
				ctx:       ctx,
				clientID:  "12345",
				idType:    enums.IdentifierTypeCCC,
				idValue:   "1224",
				isPrimary: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to add identifier" {
				fakeGorm.AddIdentifierFn = func(ctx context.Context, identifier *gorm.Identifier) (*gorm.Identifier, error) {
					return nil, fmt.Errorf("failed to add identifier")
				}
			}
			got, err := d.AddIdentifier(tt.args.ctx, tt.args.clientID, tt.args.idType, tt.args.idValue, tt.args.isPrimary)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.AddIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got :%v", got)
			}
		})
	}
}
