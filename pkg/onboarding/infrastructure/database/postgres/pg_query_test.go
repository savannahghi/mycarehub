package postgres

import (
	"context"
	"fmt"
	"strconv"
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
)

func TestOnboardingDb_RetrieveFacility_Unittest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        "KN001",
		County:      "Kanairo",
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	id := facility.ID

	invalidID := uuid.New().String()

	type args struct {
		ctx    context.Context
		id     *string
		active bool
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Facility
		wantErr bool
	}{
		{
			name: "happy case - valid ID passed",
			args: args{
				ctx:    ctx,
				id:     id,
				active: true,
			},
			wantErr: false,
		},
		{
			name: "sad case - no ID passed",
			args: args{
				ctx:    ctx,
				active: false,
			},
			wantErr: true,
		},
		{
			name: "sad case - invalid ID",
			args: args{
				ctx:    ctx,
				id:     &invalidID,
				active: false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := d.RetrieveFacility(ctx, tt.args.id, tt.args.active)

			if tt.name == "happy case - valid ID passed" {
				fakeGorm.RetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return &gorm.Facility{
						FacilityID:  facility.ID,
						Name:        facility.Name,
						Code:        facility.Code,
						Active:      strconv.FormatBool(facility.Active),
						County:      facility.County,
						Description: facility.Description,
					}, nil
				}
			}

			if tt.name == "sad case - no ID passed" {
				fakeGorm.RetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - invalid ID" {
				fakeGorm.RetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RetrieveFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestOnboardingDb_GetFacilities(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()
	name := "Kanairo One"
	code := "KN001"
	county := "Kanairo"
	description := "This is just for mocking"

	facility := &domain.Facility{
		ID:          &id,
		Name:        name,
		Code:        code,
		Active:      true,
		County:      county,
		Description: description,
	}

	facilityData := []*domain.Facility{}
	facilityData = append(facilityData, facility)
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Facility
		wantErr bool
	}{
		{
			name:    "happy case - valid payload",
			args:    args{ctx: ctx},
			want:    facilityData,
			wantErr: false,
		},
		{
			name:    "sad case - facility want data not given",
			args:    args{ctx: ctx},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case - facility want data not given" {
				fakeGorm.GetFacilitiesFn = func(ctx context.Context) ([]gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get facilities")
				}
			}
			if tt.name == "happy case - valid payload" {
				fakeGorm.GetFacilitiesFn = func(ctx context.Context) ([]gorm.Facility, error) {
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
				}
			}
			got, err := d.GetFacilities(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_RetrieveByFacilityMFLCode(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	facilityInput := &dto.FacilityInput{
		Name:        "Kanairo One",
		Code:        "KN001",
		Active:      true,
		County:      "Kanairo",
		Description: "This is just for mocking",
	}

	// Setup, create a facility
	facility, err := d.GetOrCreateFacility(ctx, facilityInput)
	if err != nil {
		t.Errorf("failed to create new facility: %v", err)
	}

	mflCode := facility.Code

	invalidMFLCode := ksuid.New().String()

	type args struct {
		ctx      context.Context
		MFLCode  string
		isActive bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:      ctx,
				MFLCode:  mflCode,
				isActive: true,
			},
			wantErr: false,
		},

		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				MFLCode:  invalidMFLCode,
				isActive: true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy case" {
				fakeGorm.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
					return &gorm.Facility{
						FacilityID:  facility.ID,
						Name:        facility.Name,
						Code:        facility.Code,
						Active:      strconv.FormatBool(facility.Active),
						County:      facility.County,
						Description: facility.Description,
					}, nil
				}
			}
			if tt.name == "Sad case" {
				fakeGorm.RetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode string, isActive bool) (*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred while retrieving facility by MFL code")
				}
			}
			got, err := d.RetrieveByFacilityMFLCode(tt.args.ctx, tt.args.MFLCode, tt.args.isActive)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.RetrieveByFacilityMFLCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestOnboardingDb_GetClientProfileByClientID(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewOnboardingDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully fetch client profile",
			args: args{
				ctx:      ctx,
				clientID: "1234",
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get profile",
			args: args{
				ctx:      ctx,
				clientID: "1234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case - Fail to get profile" {
				fakeGorm.GetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*gorm.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by ID")
				}
			}
			got, err := d.GetClientProfileByClientID(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.GetClientProfileByClientID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got :%v", got)
			}
		})
	}
}

func TestOnboardingDb_GetStaffProfileByStaffID(t *testing.T) {
	ctx := context.Background()

	testFacilityID := uuid.New().String()
	testUserID := uuid.New().String()
	testTime := time.Now()
	testID := uuid.New().String()
	staffProfileID := uuid.New().String()

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
							StaffProfileID:    &staffProfileID,
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
						},
					}, nil
				}
				fakeGorm.GetStaffProfileByStaffIDFn = func(ctx context.Context, staffProfileID string) (*gorm.StaffUserProfile, error) {
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
							StaffProfileID:    &staffProfileID,
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

func TestOnboardingDb_GetStaffProfileByStaffNumber(t *testing.T) {
	ctx := context.Background()

	testFacilityID := uuid.New().String()
	testUserID := uuid.New().String()
	testTime := time.Now()
	testID := uuid.New().String()
	staffProfileID := uuid.New().String()
	staffNumber := ksuid.New().String()

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
							StaffProfileID:    &staffProfileID,
							UserID:            &testUserID,
							StaffNumber:       staffNumber,
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
						},
					}, nil
				}
				fakeGorm.GetStaffProfileByStaffIDFn = func(ctx context.Context, staffProfileID string) (*gorm.StaffUserProfile, error) {
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
							StaffProfileID:    &staffProfileID,
							UserID:            &testUserID,
							StaffNumber:       staffNumber,
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
