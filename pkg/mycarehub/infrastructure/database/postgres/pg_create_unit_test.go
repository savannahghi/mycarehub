package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestMyCareHubDb_GetOrCreateFacility(t *testing.T) {
	ctx := context.Background()

	name := gofakeit.Name()
	code := gofakeit.Number(300, 400)
	phone := interserviceclient.TestUserPhoneNumber
	county := "Nairobi"
	description := gofakeit.HipsterSentence(15)

	facility := &dto.FacilityInput{
		Name:        name,
		Code:        code,
		Phone:       phone,
		Active:      true,
		County:      county,
		Description: description,
	}

	invalidFacility := &dto.FacilityInput{
		Name:        name,
		Active:      true,
		County:      county,
		Description: description,
	}

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
				ctx:      ctx,
				facility: facility,
			},
			wantErr: false,
		},
		{
			name: "sad case - facility code not defined",
			args: args{
				ctx:      ctx,
				facility: invalidFacility,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get ot create facility",
			args: args{
				ctx:      ctx,
				facility: facility,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case - facility code not defined" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "sad case - nil facility input" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to create facility")
				}
			}

			if tt.name == "Sad Case - Fail to get ot create facility" {
				fakeGorm.MockGetOrCreateFacilityFn = func(ctx context.Context, facility *gorm.Facility) (*gorm.Facility, error) {
					return nil, fmt.Errorf("failed to get or create facility")
				}
			}

			got, err := d.GetOrCreateFacility(tt.args.ctx, tt.args.facility)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetOrCreateFacility() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_SaveTemporaryUserPin(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()
	flavor := feedlib.FlavourConsumer

	newExtension := extension.NewExternalMethodsImpl()

	tempPin, err := newExtension.GenerateTempPIN(ctx)
	if err != nil {
		t.Errorf("failed to generate temporary pin: %v", err)
	}
	salt, encryptedTempPin := newExtension.EncryptPIN(tempPin, nil)

	pinPayload := &domain.UserPIN{
		UserID:    ID,
		HashedPIN: encryptedTempPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
		Flavour:   flavor,
		Salt:      salt,
	}
	type args struct {
		ctx     context.Context
		pinData *domain.UserPIN
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx:     ctx,
				pinData: pinPayload,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "invalid: invalid user id provided",
			args: args{
				ctx: ctx,
				pinData: &domain.UserPIN{
					UserID:    gofakeit.Sentence(200),
					HashedPIN: encryptedTempPin,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					IsValid:   true,
					Flavour:   flavor,
					Salt:      salt,
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad Case - Fail to save temporary pin",
			args: args{
				ctx:     ctx,
				pinData: pinPayload,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "invalid: missing userID" {
				fakeGorm.MockSaveTemporaryUserPinFn = func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
					return false, fmt.Errorf("user id must be provided")
				}
			}

			if tt.name == "Sad Case - Fail to save temporary pin" {
				fakeGorm.MockSaveTemporaryUserPinFn = func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
					return false, fmt.Errorf("fail to save temporary pin")
				}
			}

			got, err := d.SaveTemporaryUserPin(tt.args.ctx, tt.args.pinData)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveTemporaryUserPin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.SaveTemporaryUserPin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnboardingDb_SavePin(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		pinInput *domain.UserPIN
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save pin",
			args: args{
				ctx: ctx,
				pinInput: &domain.UserPIN{
					UserID:    uuid.New().String(),
					HashedPIN: "12345",
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					Flavour:   feedlib.FlavourConsumer,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save pin",
			args: args{
				ctx: ctx,
				pinInput: &domain.UserPIN{
					UserID:    uuid.New().String(),
					HashedPIN: "12345",
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
					Flavour:   feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save pin" {
				fakeGorm.MockSavePinFn = func(ctx context.Context, pinData *gorm.PINData) (bool, error) {
					return false, fmt.Errorf("failed to save pin")
				}
			}

			got, err := d.SavePin(tt.args.ctx, tt.args.pinInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("OnboardingDb.SavePin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("OnboardingDb.SavePin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_SaveOTP(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx      context.Context
		otpInput *domain.OTP
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save otp",
			args: args{
				ctx: ctx,
				otpInput: &domain.OTP{
					UserID:      uuid.New().String(),
					Valid:       true,
					ValidUntil:  time.Now().Add(time.Hour * 1),
					GeneratedAt: time.Now(),
					PhoneNumber: gofakeit.Phone(),
					Channel:     "SMS",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save otp",
			args: args{
				ctx: ctx,
				otpInput: &domain.OTP{
					UserID:      uuid.New().String(),
					Valid:       true,
					ValidUntil:  time.Now().Add(time.Hour * 1),
					GeneratedAt: time.Now(),
					PhoneNumber: gofakeit.Phone(),
					Channel:     "SMS",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save otp" {
				fakeGorm.MockSaveOTPFn = func(ctx context.Context, otpInput *gorm.UserOTP) error {
					return fmt.Errorf("failed to save otp")
				}
			}

			if err := d.SaveOTP(tt.args.ctx, tt.args.otpInput); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveOTP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_SaveSecurityQuestionResponse(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                      context.Context
		securityQuestionResponse []*dto.SecurityQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save question response",
			args: args{
				ctx: ctx,
				securityQuestionResponse: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             uuid.New().String(),
						SecurityQuestionID: uuid.New().String(),
						Response:           "A valid response",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save question response",
			args: args{
				ctx: ctx,
				securityQuestionResponse: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             uuid.New().String(),
						SecurityQuestionID: uuid.New().String(),
						Response:           "A valid response",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save question response" {
				fakeGorm.MockSaveSecurityQuestionResponseFn = func(ctx context.Context, securityQuestionResponse []*gorm.SecurityQuestionResponse) error {
					return fmt.Errorf("failed to save security question response")
				}
			}
			if err := d.SaveSecurityQuestionResponse(tt.args.ctx, tt.args.securityQuestionResponse); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateHealthDiaryEntry(t *testing.T) {
	ctx := context.Background()
	currentTime := time.Now()
	type args struct {
		ctx              context.Context
		healthDiaryInput *domain.ClientHealthDiaryEntry
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create a health diary entry",
			args: args{
				ctx: ctx,
				healthDiaryInput: &domain.ClientHealthDiaryEntry{
					Active:   true,
					Mood:     enums.MoodHappy.String(),
					SharedAt: &currentTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create health diary entry",
			args: args{
				ctx: ctx,
				healthDiaryInput: &domain.ClientHealthDiaryEntry{
					Active:   true,
					Mood:     enums.MoodHappy.String(),
					SharedAt: &currentTime,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to create health diary entry" {
				fakeGorm.MockCreateHealthDiaryEntryFn = func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) error {
					return fmt.Errorf("failed to create health diary entry")
				}
			}

			if err := d.CreateHealthDiaryEntry(tt.args.ctx, tt.args.healthDiaryInput); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateHealthDiaryEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateServiceRequest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx                 context.Context
		serviceRequestInput *dto.ServiceRequestInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Create a service request",
			args: args{
				ctx: ctx,
				serviceRequestInput: &dto.ServiceRequestInput{
					Active: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create a service request",
			args: args{
				ctx: ctx,
				serviceRequestInput: &dto.ServiceRequestInput{
					Active: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to create a service request" {
				fakeGorm.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *gorm.ClientServiceRequest) error {
					return fmt.Errorf("failed to create a service request")
				}
			}

			if err := d.CreateServiceRequest(tt.args.ctx, tt.args.serviceRequestInput); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateClientCaregiver(t *testing.T) {
	type args struct {
		ctx            context.Context
		caregiverInput *dto.CaregiverInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Create a client caregiver",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					FirstName:     "John",
					LastName:      "Doe",
					PhoneNumber:   "+1234567890",
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			wantErr: false,
		},

		{
			name: "Sad Case - Fail to create a client caregiver",
			args: args{
				ctx: context.Background(),
				caregiverInput: &dto.CaregiverInput{
					FirstName:     "John",
					LastName:      "Doe",
					PhoneNumber:   "+1234567890",
					CaregiverType: enums.CaregiverTypeFather,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to create a client caregiver" {
				fakeGorm.MockCreateClientCaregiverFn = func(ctx context.Context, clientID string, clientCaregiver *gorm.Caregiver) error {
					return fmt.Errorf("failed to create a client caregiver")
				}
			}
			if err := d.CreateClientCaregiver(tt.args.ctx, tt.args.caregiverInput); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateChannel(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx            context.Context
		communityInput *dto.CommunityInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				communityInput: &dto.CommunityInput{
					Name:        "test",
					Description: "test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:     []*enumutils.Gender{&enumutils.AllGender[0]},
					ClientType: []*enums.ClientType{&enums.AllClientType[0]},
					InviteOnly: true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				communityInput: &dto.CommunityInput{
					Name:        "test",
					Description: "test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:     []*enumutils.Gender{&enumutils.AllGender[0]},
					ClientType: []*enums.ClientType{&enums.AllClientType[0]},
					InviteOnly: true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockCreateCommunityFn = func(ctx context.Context, community *gorm.Community) (*gorm.Community, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CreateCommunity(tt.args.ctx, tt.args.communityInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_GetOrCreateNextOfKin(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx       context.Context
		person    *dto.NextOfKinPayload
		clientID  string
		contactID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create a next of kin",
			args: args{
				ctx: context.Background(),
				person: &dto.NextOfKinPayload{
					Name:         "John Doe",
					Relationship: "Next of Kin",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		if tt.name == "Happy case: create a next of kin" {
			fakeGorm.MockGetOrCreateNextOfKin = func(ctx context.Context, person *gorm.RelatedPerson, clientID, contactID string) error {
				return nil
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := d.GetOrCreateNextOfKin(tt.args.ctx, tt.args.person, tt.args.clientID, tt.args.contactID); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetOrCreateNextOfKin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_GetOrCreateContact(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx     context.Context
		contact *domain.Contact
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create a new contact",
			args: args{
				ctx: context.Background(),
				contact: &domain.Contact{
					Active:       true,
					ContactType:  "PHONE",
					ContactValue: gofakeit.Phone(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: Failed to create a new contact",
			args: args{
				ctx: context.Background(),
				contact: &domain.Contact{
					Active:       true,
					ContactType:  "PHONE",
					ContactValue: gofakeit.Phone(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		if tt.name == "Happy case: create a new contact" {
			fakeGorm.MockGetOrCreateContact = func(ctx context.Context, contact *gorm.Contact) (*gorm.Contact, error) {
				id := gofakeit.UUID()
				return &gorm.Contact{
					ContactID:    &id,
					ContactType:  "PHONE",
					ContactValue: gofakeit.Phone(),
					Active:       false,
				}, nil
			}
		}
		if tt.name == "Sad case: Failed to create a new contact" {
			fakeGorm.MockGetOrCreateContact = func(ctx context.Context, contact *gorm.Contact) (*gorm.Contact, error) {
				return nil, fmt.Errorf("an error occurred")
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if _, err := d.GetOrCreateContact(tt.args.ctx, tt.args.contact); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.GetOrCreateContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateAppointment(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx         context.Context
		appointment domain.Appointment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create a new appointment",
			args: args{
				ctx: context.Background(),
				appointment: domain.Appointment{
					ID:         gofakeit.UUID(),
					Reason:     "Dental",
					ClientID:   gofakeit.UUID(),
					ExternalID: gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.CreateAppointment(tt.args.ctx, tt.args.appointment); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateAppointment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateUser(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	date := gofakeit.Date()

	type args struct {
		ctx  context.Context
		user domain.User
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.User
		wantErr bool
	}{
		{
			name: "happy case: create user",
			args: args{
				ctx: context.Background(),
				user: domain.User{
					Username:    gofakeit.Username(),
					Name:        gofakeit.Name(),
					Gender:      enumutils.GenderMale,
					DateOfBirth: &date,
					UserType:    enums.ClientUser,
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CreateUser(tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected user to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected user not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_CreateClient(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	enrollment := time.Now()
	userID := gofakeit.UUID()
	var clientTypeList []enums.ClientType
	clientTypeList = append(clientTypeList, enums.AllClientType...)

	type args struct {
		ctx          context.Context
		client       domain.ClientProfile
		contactID    string
		identifierID string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientProfile
		wantErr bool
	}{
		{
			name: "happy case: create client",
			args: args{
				ctx: context.Background(),
				client: domain.ClientProfile{
					Active:                  true,
					UserID:                  userID,
					FacilityID:              gofakeit.UUID(),
					ClientCounselled:        true,
					ClientTypes:             clientTypeList,
					TreatmentEnrollmentDate: &enrollment,
				},
				contactID:    gofakeit.UUID(),
				identifierID: gofakeit.UUID(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CreateClient(tt.args.ctx, tt.args.client, tt.args.contactID, tt.args.identifierID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected client to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected client not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_CreateIdentifier(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		identifier domain.Identifier
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.Identifier
		wantErr bool
	}{
		{
			name: "happy case: create identifier",
			args: args{
				ctx: context.Background(),
				identifier: domain.Identifier{
					IdentifierType:      "CCC",
					IdentifierValue:     "5678901234789",
					IdentifierUse:       "OFFICIAL",
					Description:         "CCC Number, Primary Identifier",
					IsPrimaryIdentifier: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CreateIdentifier(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && got != nil {
				t.Errorf("expected identifier to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected identifier not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestMyCareHubDb_AnswerScreeningToolQuestions(t *testing.T) {
	type args struct {
		ctx                    context.Context
		screeningToolResponses []*dto.ScreeningToolQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: answer screening tool questions",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to answer screening tool questions",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to answer screening tool questions" {
				fakeGorm.MockAnswerScreeningToolQuestionsFn = func(ctx context.Context, screeningToolResponses []*gorm.ScreeningToolsResponse) error {
					return fmt.Errorf("an error occurred")
				}
			}

			if err := d.AnswerScreeningToolQuestions(tt.args.ctx, tt.args.screeningToolResponses); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AnswerScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateStaffServiceRequest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	meta := map[string]interface{}{
		"client_id": uuid.New().String(),
	}

	type args struct {
		ctx                 context.Context
		serviceRequestInput *dto.ServiceRequestInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				serviceRequestInput: &dto.ServiceRequestInput{
					Active:      true,
					RequestType: uuid.New().String(),
					Status:      uuid.New().String(),
					Request:     uuid.New().String(),
					StaffID:     uuid.New().String(),
					FacilityID:  uuid.New().String(),
					StaffName:   uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				serviceRequestInput: &dto.ServiceRequestInput{
					Active:      true,
					RequestType: uuid.New().String(),
					Status:      uuid.New().String(),
					Request:     uuid.New().String(),
					StaffID:     uuid.New().String(),
					FacilityID:  uuid.New().String(),
					StaffName:   uuid.New().String(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid meta",
			args: args{
				ctx: ctx,
				serviceRequestInput: &dto.ServiceRequestInput{
					Active:      true,
					RequestType: uuid.New().String(),
					Status:      uuid.New().String(),
					Request:     uuid.New().String(),
					StaffID:     uuid.New().String(),
					FacilityID:  uuid.New().String(),
					StaffName:   uuid.New().String(),
					Meta:        meta,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockCreateStaffServiceRequestFn = func(ctx context.Context, serviceRequestInput *gorm.StaffServiceRequest) error {
					return fmt.Errorf("failed to create staff service request")
				}
			}
			if tt.name == "Sad case - invalid meta" {
				fakeGorm.MockCreateStaffServiceRequestFn = func(ctx context.Context, serviceRequestInput *gorm.StaffServiceRequest) error {
					return fmt.Errorf("failed to create staff service request")
				}
			}
			if err := d.CreateStaffServiceRequest(tt.args.ctx, tt.args.serviceRequestInput); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateStaffServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_SaveNotification(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()
	type args struct {
		ctx     context.Context
		payload *domain.Notification
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully save a notification",
			args: args{
				ctx: ctx,
				payload: &domain.Notification{
					Title:      "An introduction",
					Body:       "This is a new introduction",
					Type:       "Test Notification",
					IsRead:     false,
					UserID:     &UUID,
					FacilityID: &UUID,
					Flavour:    feedlib.FlavourConsumer,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to save a notification",
			args: args{
				ctx: ctx,
				payload: &domain.Notification{
					Title:      "An introduction",
					Body:       "This is a new introduction",
					Type:       "Test Notification",
					IsRead:     false,
					UserID:     &UUID,
					FacilityID: &UUID,
					Flavour:    feedlib.FlavourConsumer,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to save a notification" {
				fakeGorm.MockCreateNotificationFn = func(ctx context.Context, notification *gorm.Notification) error {
					return fmt.Errorf("failed to save notification")
				}
			}

			if err := d.SaveNotification(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateMetric(t *testing.T) {
	userID := gofakeit.UUID()

	type args struct {
		ctx     context.Context
		payload *domain.Metric
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create a metric",
			args: args{
				ctx: context.Background(),
				payload: &domain.Metric{
					UserID:    &userID,
					Timestamp: time.Now(),
					Type:      enums.MetricTypeContent,
					Event: map[string]interface{}{
						"contentID": 20,
						"duration":  time.Since(time.Now()),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error creating metric",
			args: args{
				ctx: context.Background(),
				payload: &domain.Metric{
					UserID:    &userID,
					Timestamp: time.Now(),
					Type:      enums.MetricTypeContent,
					Event: map[string]interface{}{
						"contentID": 10,
						"duration":  time.Since(time.Now()),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: error creating metric" {
				fakeGorm.MockCreateMetricFn = func(ctx context.Context, metric *gorm.Metric) error {
					return fmt.Errorf("cannot create metric")
				}
			}

			if err := d.CreateMetric(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateMetric() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateUserSurvey(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx         context.Context
		userSurveys []*dto.UserSurveyInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case",
			args: args{
				ctx: context.Background(),
				userSurveys: []*dto.UserSurveyInput{
					{
						UserID:      uuid.New().String(),
						Title:       "Test Survey",
						Description: "This is a test survey",
						Link:        gofakeit.URL(),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := d.CreateUserSurveys(tt.args.ctx, tt.args.userSurveys); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateUserSurveys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_RegisterStaff(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	currentTime := time.Now()

	pastYear := time.Now().AddDate(-3, 0, 0)
	ID := uuid.New().String()

	userProfile := &domain.User{
		ID:                  &ID,
		Username:            gofakeit.Name(),
		Name:                gofakeit.Name(),
		Active:              true,
		TermsAccepted:       true,
		Gender:              enumutils.GenderMale,
		LastSuccessfulLogin: &currentTime,
		NextAllowedLogin:    &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    3,
		Contacts: &domain.Contact{
			ID:           &ID,
			ContactType:  "PHONE",
			ContactValue: "+254711223344",
			Active:       true,
			OptedIn:      true,
			UserID:       &ID,
			Flavour:      "CONSUMER",
		},
		DateOfBirth: &pastYear,
	}

	staff := &domain.StaffProfile{
		ID:                &ID,
		User:              userProfile,
		UserID:            uuid.New().String(),
		Active:            false,
		StaffNumber:       gofakeit.BeerAlcohol(),
		DefaultFacilityID: gofakeit.BeerAlcohol(),
	}

	contact := &domain.Contact{
		ID:           &ID,
		ContactType:  "PHONE",
		ContactValue: interserviceclient.TestUserPhoneNumber,
		Active:       true,
		OptedIn:      true,
		UserID:       &staff.UserID,
		Flavour:      feedlib.FlavourPro,
	}

	identifierData := &domain.Identifier{
		ID:                  ID,
		IdentifierType:      ID,
		IdentifierValue:     ID,
		IdentifierUse:       ID,
		Description:         "Valid Identifier",
		ValidFrom:           time.Now(),
		ValidTo:             time.Now(),
		IsPrimaryIdentifier: true,
		Active:              true,
	}

	payload := &domain.StaffRegistrationPayload{
		UserProfile:     *userProfile,
		Phone:           *contact,
		StaffIdentifier: *identifierData,
		Staff:           *staff,
	}
	type args struct {
		ctx     context.Context
		payload *domain.StaffRegistrationPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy Case: Register staff",
			args: args{
				ctx:     context.Background(),
				payload: payload,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Unable to register staff",
			args: args{
				ctx:     context.Background(),
				payload: payload,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case: Unable to register staff" {
				fakeGorm.MockRegisterStaffFn = func(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) error {
					return fmt.Errorf("cannot register staff")
				}
			}
			_, err := d.RegisterStaff(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_SaveFeedback(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	feedback := &domain.FeedbackResponse{
		UserID:            uuid.New().String(),
		FeedbackType:      "TEST",
		SatisfactionLevel: 1,
		ServiceName:       "TEST",
		Feedback:          "TEST",
		RequiresFollowUp:  true,
		PhoneNumber:       interserviceclient.TestUserPhoneNumber,
	}

	type args struct {
		ctx     context.Context
		payload *domain.FeedbackResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully save feedback",
			args: args{
				ctx:     context.Background(),
				payload: feedback,
			},
			wantErr: false,
		},
		{
			name: "Sade Case: Unable to save feedback",
			args: args{
				ctx:     context.Background(),
				payload: feedback,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sade Case: Unable to save feedback" {
				fakeGorm.MockSaveFeedbackFn = func(ctx context.Context, feedback *gorm.Feedback) error {
					return fmt.Errorf("cannot save feedback")
				}
			}
			if err := d.SaveFeedback(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SaveFeedback() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_RegisterClient(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	UID := uuid.New().String()

	type args struct {
		ctx     context.Context
		payload *domain.ClientRegistrationPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully register client",
			args: args{
				ctx: context.Background(),
				payload: &domain.ClientRegistrationPayload{
					UserProfile: domain.User{
						ID:       &UID,
						Username: "test",
					},
					Phone: domain.Contact{
						ID:           &UID,
						ContactType:  "PHONE",
						ContactValue: interserviceclient.TestUserPhoneNumber,
						Active:       true,
						OptedIn:      true,
						UserID:       &UID,
						Flavour:      feedlib.FlavourConsumer,
					},
					ClientIdentifier: domain.Identifier{
						IdentifierType:  "CCC",
						IdentifierValue: "123456789",
					},
					Client: domain.ClientProfile{
						ID:          &UID,
						ClientTypes: []enums.ClientType{"PMTCT"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: unable register client",
			args: args{
				ctx: context.Background(),
				payload: &domain.ClientRegistrationPayload{
					UserProfile: domain.User{
						ID:       &UID,
						Username: "test",
					},
					Phone: domain.Contact{
						ID:           &UID,
						ContactType:  "PHONE",
						ContactValue: interserviceclient.TestUserPhoneNumber,
						Active:       true,
						OptedIn:      true,
						UserID:       &UID,
						Flavour:      feedlib.FlavourConsumer,
					},
					ClientIdentifier: domain.Identifier{
						IdentifierType:  "CCC",
						IdentifierValue: "123456789",
					},
					Client: domain.ClientProfile{
						ID:          &UID,
						ClientTypes: []enums.ClientType{"PMTCT"},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case: unable register client" {
				fakeGorm.MockRegisterClientFn = func(ctx context.Context, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) error {
					return fmt.Errorf("cannot register client")
				}
			}
			_, err := d.RegisterClient(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_CreateScreeningTool(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	UID := uuid.New().String()

	screeningTool := &domain.ScreeningTool{
		ID:              UID,
		Active:          true,
		QuestionnaireID: uuid.New().String(),
		Threshold:       3,
		ClientTypes:     []enums.ClientType{"PMTCT"},
		Genders:         []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: domain.AgeRange{
			LowerBound: 10,
			UpperBound: 20,
		},
		Questionnaire: domain.Questionnaire{
			ID:          UID,
			Active:      true,
			Name:        gofakeit.BeerName(),
			Description: "Test Description",
			Questions: []domain.Question{
				{
					ID:                UID,
					Active:            true,
					QuestionnaireID:   UID,
					Text:              gofakeit.Sentence(50),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeNumber,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []domain.QuestionInputChoice{
						{
							ID:         UID,
							Active:     true,
							QuestionID: UID,
							Choice:     "YES",
							Value:      "yes",
							Score:      1,
						},
					},
				},
			},
		},
	}

	type args struct {
		ctx   context.Context
		input *domain.ScreeningTool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully create screening tool",
			args: args{
				ctx:   context.Background(),
				input: screeningTool,
			},
			wantErr: false,
		},
		{
			name: "Sad Case: Unable to create questionnaire",
			args: args{
				ctx:   context.Background(),
				input: screeningTool,
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Unable to create screening tool",
			args: args{
				ctx:   context.Background(),
				input: screeningTool,
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Unable to create question",
			args: args{
				ctx:   context.Background(),
				input: screeningTool,
			},
			wantErr: true,
		},
		{
			name: "Sad Case: Unable to create question input choice",
			args: args{
				ctx:   context.Background(),
				input: screeningTool,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case: Unable to create questionnaire" {
				fakeGorm.MockCreateQuestionnaireFn = func(ctx context.Context, questionnaire *gorm.Questionnaire) error {
					return fmt.Errorf("cannot create questionnaire")
				}
			}
			if tt.name == "Sad Case: Unable to create screening tool" {
				fakeGorm.MockCreateQuestionnaireFn = func(ctx context.Context, input *gorm.Questionnaire) error {
					return nil
				}
				fakeGorm.MockCreateScreeningToolFn = func(ctx context.Context, screeningTool *gorm.ScreeningTool) error {
					return fmt.Errorf("cannot create screening tool")
				}
			}
			if tt.name == "Sad Case: Unable to create question" {
				fakeGorm.MockCreateQuestionnaireFn = func(ctx context.Context, input *gorm.Questionnaire) error {
					return nil
				}
				fakeGorm.MockCreateScreeningToolFn = func(ctx context.Context, screeningTool *gorm.ScreeningTool) error {
					return nil
				}
				fakeGorm.MockCreateQuestionFn = func(ctx context.Context, question *gorm.Question) error {
					return fmt.Errorf("cannot create question")
				}
			}
			if tt.name == "Sad Case: Unable to create question input choice" {
				fakeGorm.MockCreateQuestionnaireFn = func(ctx context.Context, input *gorm.Questionnaire) error {
					return nil
				}
				fakeGorm.MockCreateScreeningToolFn = func(ctx context.Context, screeningTool *gorm.ScreeningTool) error {
					return nil
				}
				fakeGorm.MockCreateQuestionFn = func(ctx context.Context, question *gorm.Question) error {
					return nil
				}
				fakeGorm.MockCreateQuestionChoiceFn = func(ctx context.Context, input *gorm.QuestionInputChoice) error {
					return fmt.Errorf("cannot create question input choice")
				}
			}
			if err := d.CreateScreeningTool(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateScreeningTool() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateScreeningToolResponse(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx   context.Context
		input *domain.QuestionnaireScreeningToolResponse
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Successfully create screening tool response",
			args: args{
				ctx: context.Background(),
				input: &domain.QuestionnaireScreeningToolResponse{
					ID:              uuid.NewString(),
					Active:          true,
					ScreeningToolID: uuid.NewString(),
					FacilityID:      uuid.NewString(),
					ClientID:        uuid.NewString(),
					AggregateScore:  3,
					QuestionResponses: []domain.QuestionnaireScreeningToolQuestionResponse{
						{
							ID:                      uuid.NewString(),
							Active:                  true,
							ScreeningToolResponseID: uuid.NewString(),
							QuestionID:              uuid.NewString(),
							Response:                "response",
							Score:                   0,
						},
						{
							ID:                      uuid.NewString(),
							Active:                  true,
							ScreeningToolResponseID: uuid.NewString(),
							QuestionID:              uuid.NewString(),
							Response:                "response",
							Score:                   0,
						},
					},
				},
			},
		},
		{
			name: "Sad Case: failed to create screening tool response",
			args: args{
				ctx: context.Background(),
				input: &domain.QuestionnaireScreeningToolResponse{
					ID:              uuid.NewString(),
					Active:          true,
					ScreeningToolID: uuid.NewString(),
					FacilityID:      uuid.NewString(),
					ClientID:        uuid.NewString(),
					AggregateScore:  3,
					QuestionResponses: []domain.QuestionnaireScreeningToolQuestionResponse{
						{
							ID:                      uuid.NewString(),
							Active:                  true,
							ScreeningToolResponseID: uuid.NewString(),
							QuestionID:              uuid.NewString(),
							Response:                "response",
							Score:                   0,
						},
						{
							ID:                      uuid.NewString(),
							Active:                  true,
							ScreeningToolResponseID: uuid.NewString(),
							QuestionID:              uuid.NewString(),
							Response:                "response",
							Score:                   0,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case: failed to create screening tool response" {
				fakeGorm.MockCreateScreeningToolResponseFn = func(ctx context.Context, screeningToolResponse *gorm.ScreeningToolResponse, screeningToolQuestionResponses []*gorm.ScreeningToolQuestionResponse) (*string, error) {
					return nil, fmt.Errorf("cannot create screening tool response")
				}
			}

			got, err := d.CreateScreeningToolResponse(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateScreeningToolResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
