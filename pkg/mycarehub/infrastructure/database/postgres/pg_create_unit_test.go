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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
)

func TestMyCareHubDb_SaveTemporaryUserPin(t *testing.T) {
	ctx := context.Background()
	ID := uuid.New().String()

	tempPin, err := utils.GenerateTempPIN(ctx)
	if err != nil {
		t.Errorf("failed to generate temporary pin: %v", err)
	}
	salt, encryptedTempPin := utils.EncryptPIN(tempPin, nil)

	pinPayload := &domain.UserPIN{
		UserID:    ID,
		HashedPIN: encryptedTempPin,
		ValidFrom: time.Now(),
		ValidTo:   time.Now(),
		IsValid:   true,
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
					Flavour:     feedlib.FlavourPro,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case -  invalid flavour",
			args: args{
				ctx: ctx,
				otpInput: &domain.OTP{
					UserID:      uuid.New().String(),
					Valid:       true,
					ValidUntil:  time.Now().Add(time.Hour * 1),
					GeneratedAt: time.Now(),
					PhoneNumber: gofakeit.Phone(),
					Channel:     "SMS",
					Flavour:     feedlib.Flavour("invalid"),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case -  missing phone",
			args: args{
				ctx: ctx,
				otpInput: &domain.OTP{
					UserID:      uuid.New().String(),
					Valid:       true,
					ValidUntil:  time.Now().Add(time.Hour * 1),
					GeneratedAt: time.Now(),
					Channel:     "SMS",
					Flavour:     feedlib.FlavourPro,
				},
			},
			wantErr: true,
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
					Flavour:     feedlib.FlavourPro,
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

func TestMyCareHubDb_RegisterCaregiver(t *testing.T) {
	dob := time.Now()
	type args struct {
		ctx   context.Context
		input *domain.CaregiverRegistration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: register caregiver",
			args: args{
				ctx: nil,
				input: &domain.CaregiverRegistration{
					User: &domain.User{
						Username:    gofakeit.Username(),
						Name:        gofakeit.Name(),
						Gender:      enumutils.GenderFemale,
						DateOfBirth: &dob,
						Active:      true,
					},
					Contact: &domain.Contact{
						ContactType:  "PHONE",
						ContactValue: gofakeit.Phone(),
						Active:       true,
						OptedIn:      false,
					},
					Caregiver: &domain.Caregiver{
						Active:          true,
						CaregiverNumber: gofakeit.SSN(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to create caregiver",
			args: args{
				ctx: nil,
				input: &domain.CaregiverRegistration{
					User: &domain.User{
						Username:    gofakeit.Username(),
						Name:        gofakeit.Name(),
						Gender:      enumutils.GenderFemale,
						DateOfBirth: &dob,
						Active:      true,
					},
					Contact: &domain.Contact{
						ContactType:  "PHONE",
						ContactValue: gofakeit.Phone(),
						Active:       true,
						OptedIn:      false,
					},
					Caregiver: &domain.Caregiver{
						Active:          true,
						CaregiverNumber: gofakeit.SSN(),
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

			if tt.name == "sad case: fail to create caregiver" {
				fakeGorm.MockRegisterCaregiverFn = func(ctx context.Context, user *gorm.User, contact *gorm.Contact, caregiver *gorm.Caregiver) error {
					return fmt.Errorf("failed to register")
				}
			}

			_, err := d.RegisterCaregiver(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
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
				fakeGorm.MockCreateHealthDiaryEntryFn = func(ctx context.Context, healthDiaryInput *gorm.ClientHealthDiaryEntry) (*gorm.ClientHealthDiaryEntry, error) {
					return nil, fmt.Errorf("failed to create health diary entry")
				}
			}

			if _, err := d.CreateHealthDiaryEntry(tt.args.ctx, tt.args.healthDiaryInput); (err != nil) != tt.wantErr {
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

func TestMyCareHubDb_CreateCommunity(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx              context.Context
		communityPayload *domain.Community
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
				communityPayload: &domain.Community{
					Name:           "test",
					Description:    "test",
					AgeRange:       &domain.AgeRange{LowerBound: 0, UpperBound: 0},
					Gender:         []enumutils.Gender{enumutils.AllGender[0]},
					ClientType:     []enums.ClientType{enums.AllClientType[0]},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				communityPayload: &domain.Community{
					Name:           "test",
					Description:    "test",
					AgeRange:       &domain.AgeRange{LowerBound: 0, UpperBound: 0},
					Gender:         []enumutils.Gender{enumutils.AllGender[0]},
					ClientType:     []enums.ClientType{enums.AllClientType[0]},
					OrganisationID: uuid.NewString(),
					ProgramID:      uuid.NewString(),
					FacilityID:     uuid.NewString(),
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

			got, err := d.CreateCommunity(tt.args.ctx, tt.args.communityPayload)
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
					ID:     id,
					Type:   "PHONE",
					Value:  gofakeit.Phone(),
					Active: false,
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
	programID := uuid.New().String()

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
					Username:         gofakeit.Username(),
					Name:             gofakeit.Name(),
					Gender:           enumutils.GenderMale,
					DateOfBirth:      &date,
					CurrentProgramID: programID,
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
	UUID := gofakeit.UUID()
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
					Active: true,
					UserID: userID,
					DefaultFacility: &domain.Facility{
						ID: &UUID,
					},
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
		},
		DateOfBirth: &pastYear,
	}

	staff := &domain.StaffProfile{
		ID:          &ID,
		User:        userProfile,
		UserID:      uuid.New().String(),
		Active:      false,
		StaffNumber: gofakeit.BeerAlcohol(),
		DefaultFacility: &domain.Facility{
			ID:   &ID,
			Name: gofakeit.Name(),
		},
	}

	contact := &domain.Contact{
		ID:           &ID,
		ContactType:  "PHONE",
		ContactValue: interserviceclient.TestUserPhoneNumber,
		Active:       true,
		OptedIn:      true,
		UserID:       &staff.UserID,
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
				fakeGorm.MockRegisterStaffFn = func(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, staffProfile *gorm.StaffProfile) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("cannot register staff")
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
					},
					ClientIdentifier: domain.Identifier{
						IdentifierType:  "CCC",
						IdentifierValue: "123456789",
					},
					Client: domain.ClientProfile{
						ID:          &UID,
						ClientTypes: []enums.ClientType{"PMTCT"},
						DefaultFacility: &domain.Facility{
							ID: &UID,
						},
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
					},
					ClientIdentifier: domain.Identifier{
						IdentifierType:  "CCC",
						IdentifierValue: "123456789",
					},
					Client: domain.ClientProfile{
						ID:          &UID,
						ClientTypes: []enums.ClientType{"PMTCT"},
						DefaultFacility: &domain.Facility{
							ID: &UID,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad Case: unable register client" {
				fakeGorm.MockRegisterClientFn = func(ctx context.Context, user *gorm.User, contact *gorm.Contact, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error) {
					return nil, fmt.Errorf("cannot register client")
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
					QuestionResponses: []*domain.QuestionnaireScreeningToolQuestionResponse{
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
					QuestionResponses: []*domain.QuestionnaireScreeningToolQuestionResponse{
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

func TestMyCareHubDb_AddCaregiverToClient(t *testing.T) {
	type args struct {
		ctx             context.Context
		clientCaregiver *domain.CaregiverClient
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: add caregiver to client",
			args: args{
				ctx: context.Background(),
				clientCaregiver: &domain.CaregiverClient{
					CaregiverID:      uuid.NewString(),
					ClientID:         uuid.NewString(),
					RelationshipType: enums.CaregiverTypeFather,
					Active:           true,
					AssignedBy:       uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to add caregiver to client",
			args: args{
				ctx: context.Background(),
				clientCaregiver: &domain.CaregiverClient{
					CaregiverID:      uuid.NewString(),
					ClientID:         uuid.NewString(),
					RelationshipType: enums.CaregiverTypeFather,
					Active:           true,
					AssignedBy:       uuid.NewString(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to add caregiver to client" {
				fakeGorm.MockAddCaregiverToClientFn = func(ctx context.Context, clientCaregiver *gorm.CaregiverClient) error {
					return fmt.Errorf("unable to add caregiver to client")
				}
			}

			if err := d.AddCaregiverToClient(tt.args.ctx, tt.args.clientCaregiver); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AddCaregiverToClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateCaregiver(t *testing.T) {

	type args struct {
		ctx       context.Context
		caregiver domain.Caregiver
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create caregiver",
			args: args{
				ctx: context.Background(),
				caregiver: domain.Caregiver{
					UserID:          gofakeit.UUID(),
					CaregiverNumber: gofakeit.SSN(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: create caregiver error",
			args: args{
				ctx: context.Background(),
				caregiver: domain.Caregiver{
					UserID:          gofakeit.UUID(),
					CaregiverNumber: gofakeit.SSN(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: create caregiver error" {
				fakeGorm.MockCreateCaregiverFn = func(ctx context.Context, caregiver *gorm.Caregiver) error {
					return fmt.Errorf("failed to create caregiver")
				}
			}

			got, err := d.CreateCaregiver(tt.args.ctx, tt.args.caregiver)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}

func TestMyCareHubDb_CreateOrganisation(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx          context.Context
		organisation *domain.Organisation
		programs     []*domain.Program
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: create organisation",
			args: args{
				ctx: context.Background(),
				organisation: &domain.Organisation{
					Active:          true,
					Code:            uuid.New().String(),
					Name:            gofakeit.Company(),
					Description:     gofakeit.Sentence(5),
					EmailAddress:    gofakeit.Email(),
					PhoneNumber:     gofakeit.Phone(),
					PostalAddress:   gofakeit.Address().Address,
					PhysicalAddress: gofakeit.Address().Address,
					DefaultCountry:  gofakeit.Country(),
				},
				programs: []*domain.Program{
					{
						Active:      true,
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: create organisation, no program",
			args: args{
				ctx: context.Background(),
				organisation: &domain.Organisation{
					Active:          true,
					Code:            uuid.New().String(),
					Name:            gofakeit.Company(),
					Description:     gofakeit.Sentence(5),
					EmailAddress:    gofakeit.Email(),
					PhoneNumber:     gofakeit.Phone(),
					PostalAddress:   gofakeit.Address().Address,
					PhysicalAddress: gofakeit.Address().Address,
					DefaultCountry:  gofakeit.Country(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to create organisation",
			args: args{
				ctx: context.Background(),
				organisation: &domain.Organisation{
					Active:          true,
					Code:            uuid.New().String(),
					Name:            gofakeit.Company(),
					Description:     gofakeit.Sentence(5),
					EmailAddress:    gofakeit.Email(),
					PhoneNumber:     gofakeit.Phone(),
					PostalAddress:   gofakeit.Address().Address,
					PhysicalAddress: gofakeit.Address().Address,
					DefaultCountry:  gofakeit.Country(),
				},
				programs: []*domain.Program{
					{
						Active:      true,
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to create program",
			args: args{
				ctx: context.Background(),
				organisation: &domain.Organisation{
					Active:          true,
					Code:            uuid.New().String(),
					Name:            gofakeit.Company(),
					Description:     gofakeit.Sentence(5),
					EmailAddress:    gofakeit.Email(),
					PhoneNumber:     gofakeit.Phone(),
					PostalAddress:   gofakeit.Address().Address,
					PhysicalAddress: gofakeit.Address().Address,
					DefaultCountry:  gofakeit.Country(),
				},
				programs: []*domain.Program{
					{
						Active:      true,
						Name:        gofakeit.BS(),
						Description: gofakeit.BS(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: unable to create organisation" {
				fakeGorm.MockCreateOrganisationFn = func(ctx context.Context, organization *gorm.Organisation) (*gorm.Organisation, error) {
					return nil, fmt.Errorf("failed to create organisation")
				}
			}

			if tt.name == "sad case: unable to create program" {
				fakeGorm.MockCreateProgramFn = func(ctx context.Context, program *gorm.Program) (*gorm.Program, error) {
					return nil, fmt.Errorf("failed to create program")
				}
			}

			_, err := d.CreateOrganisation(tt.args.ctx, tt.args.organisation, tt.args.programs)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateOrganisation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateProgram(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *dto.ProgramInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create program",
			args: args{
				ctx: context.Background(),
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerBlg(),
					OrganisationID: uuid.NewString(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create program",
			args: args{
				ctx: context.Background(),
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerBlg(),
					OrganisationID: uuid.NewString(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()

			if tt.name == "Sad case: failed to create program" {
				fakeGorm.MockCreateProgramFn = func(ctx context.Context, program *gorm.Program) (*gorm.Program, error) {
					return nil, fmt.Errorf("failed to create program")
				}
			}
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			_, err := d.CreateProgram(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_RegisterExistingUserAsClient(t *testing.T) {
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
		wantErr bool
	}{
		{
			name: "Happy case: register existing user as client",
			args: args{
				ctx: context.Background(),
				payload: &domain.ClientRegistrationPayload{
					ClientIdentifier: domain.Identifier{
						IdentifierType:  "CCC",
						IdentifierValue: "123456789",
					},
					Client: domain.ClientProfile{
						ID:          &UID,
						ClientTypes: []enums.ClientType{"PMTCT"},
						DefaultFacility: &domain.Facility{
							ID: &UID,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register existing user as client",
			args: args{
				ctx: context.Background(),
				payload: &domain.ClientRegistrationPayload{
					ClientIdentifier: domain.Identifier{
						IdentifierType:  "CCC",
						IdentifierValue: "123456789",
					},
					Client: domain.ClientProfile{
						ID:          &UID,
						ClientTypes: []enums.ClientType{"PMTCT"},
						DefaultFacility: &domain.Facility{
							ID: &UID,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to register existing user as client" {
				fakeGorm.MockRegisterExistingUserAsClientFn = func(ctx context.Context, identifier *gorm.Identifier, client *gorm.Client) (*gorm.Client, error) {
					return nil, fmt.Errorf("failed to register existing user as client")
				}
			}
			_, err := d.RegisterExistingUserAsClient(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterExistingUserAsClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_AddFacilityToProgram(t *testing.T) {
	type args struct {
		ctx        context.Context
		programID  string
		facilityID []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: add facility to program",
			args: args{
				ctx:        context.Background(),
				programID:  uuid.NewString(),
				facilityID: []string{uuid.NewString()},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to add facility to program",
			args: args{
				ctx:        context.Background(),
				programID:  uuid.NewString(),
				facilityID: []string{uuid.NewString()},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get program facilities",
			args: args{
				ctx:        context.Background(),
				programID:  uuid.NewString(),
				facilityID: []string{uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to add facility to program" {
				fakeGorm.MockAddFacilityToProgramFn = func(ctx context.Context, programID string, facilityID []string) error {
					return fmt.Errorf("failed to add facility to program")
				}
			}
			if tt.name == "Sad case: unable to get program facilities" {
				fakeGorm.MockGetProgramFacilitiesFn = func(ctx context.Context, programID string) ([]*gorm.ProgramFacility, error) {
					return nil, fmt.Errorf("unable to get program facilities")
				}
			}

			_, err := d.AddFacilityToProgram(tt.args.ctx, tt.args.programID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AddFacilityToProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_RegisterExistingUserAsStaff(t *testing.T) {
	UID := uuid.New().String()

	staff := &domain.StaffProfile{
		ID:          &UID,
		UserID:      uuid.New().String(),
		Active:      false,
		StaffNumber: gofakeit.BeerAlcohol(),
		DefaultFacility: &domain.Facility{
			ID:   &UID,
			Name: gofakeit.Name(),
		},
	}

	identifierData := &domain.Identifier{
		ID:                  UID,
		IdentifierType:      UID,
		IdentifierValue:     UID,
		IdentifierUse:       UID,
		Description:         "Valid Identifier",
		ValidFrom:           time.Now(),
		ValidTo:             time.Now(),
		IsPrimaryIdentifier: true,
		Active:              true,
	}

	type args struct {
		ctx     context.Context
		payload *domain.StaffRegistrationPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: register existing user as staff",
			args: args{
				ctx: context.Background(),
				payload: &domain.StaffRegistrationPayload{
					StaffIdentifier: *identifierData,
					Staff:           *staff,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register existing user as staff",
			args: args{
				ctx: context.Background(),
				payload: &domain.StaffRegistrationPayload{
					StaffIdentifier: *identifierData,
					Staff:           *staff,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to register existing user as staff" {
				fakeGorm.MockRegisterExistingUserAsStaffFn = func(ctx context.Context, identifier *gorm.Identifier, staff *gorm.StaffProfile) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("unable to register existing user as staff")
				}
			}

			_, err := d.RegisterExistingUserAsStaff(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterExistingUserAsStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_RegisterExistingUserAsCaregiver(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *domain.CaregiverRegistration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: register existing user as caregiver",
			args: args{
				ctx: context.Background(),
				input: &domain.CaregiverRegistration{
					Caregiver: &domain.Caregiver{
						ID:              uuid.NewString(),
						UserID:          uuid.NewString(),
						CaregiverNumber: "123456789",
						Active:          true,
						OrganisationID:  uuid.NewString(),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to register existing user as caregiver",
			args: args{
				ctx: context.Background(),
				input: &domain.CaregiverRegistration{
					Caregiver: &domain.Caregiver{
						ID:              uuid.NewString(),
						UserID:          uuid.NewString(),
						CaregiverNumber: "123456789",
						Active:          true,
						OrganisationID:  uuid.NewString(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to get user profile bu user id",
			args: args{
				ctx: context.Background(),
				input: &domain.CaregiverRegistration{
					Caregiver: &domain.Caregiver{
						ID:              uuid.NewString(),
						UserID:          uuid.NewString(),
						CaregiverNumber: "123456789",
						Active:          true,
						OrganisationID:  uuid.NewString(),
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

			if tt.name == "Sad case: unable to register existing user as caregiver" {
				fakeGorm.MockRegisterExistingUserAsCaregiverFn = func(ctx context.Context, caregiver *gorm.Caregiver) (*gorm.Caregiver, error) {
					return nil, fmt.Errorf("failed to register existing user as caregiver")
				}
			}
			if tt.name == "Sad case: unable to get user profile bu user id" {
				fakeGorm.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID *string) (*gorm.User, error) {
					return nil, fmt.Errorf("failed to get user profile by user id")
				}
			}

			_, err := d.RegisterExistingUserAsCaregiver(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RegisterExistingUserAsCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_CreateFacilities(t *testing.T) {
	type args struct {
		ctx        context.Context
		facilities []*domain.Facility
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create facilities",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						Name:        gofakeit.BS(),
						Phone:       "0777777777",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: domain.FacilityIdentifier{
							Active: true,
							Type:   enums.FacilityIdentifierTypeMFLCode,
							Value:  "392893828",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to create facilities",
			args: args{
				ctx: context.Background(),
				facilities: []*domain.Facility{
					{
						Name:        gofakeit.BS(),
						Phone:       "0999999999",
						Active:      true,
						Country:     "Kenya",
						Description: gofakeit.BS(),
						Identifier: domain.FacilityIdentifier{
							Active: true,
							Type:   enums.FacilityIdentifierTypeMFLCode,
							Value:  "09090908",
						},
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

			if tt.name == "Sad case: failed to create facilities" {
				fakeGorm.MockCreateFacilitiesFn = func(ctx context.Context, facilities []*gorm.Facility) ([]*gorm.Facility, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.CreateFacilities(tt.args.ctx, tt.args.facilities)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a a value to be returned, got: %v", got)
			}
		})
	}
}

func TestMyCareHubDb_CreateSecurityQuestions(t *testing.T) {
	type args struct {
		ctx               context.Context
		securityQuestions []*domain.SecurityQuestion
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: create security questions",
			args: args{
				ctx: context.Background(),
				securityQuestions: []*domain.SecurityQuestion{{
					SecurityQuestionID: gofakeit.UUID(),
					QuestionStem:       gofakeit.Question(),
					Description:        gofakeit.BS(),
					Flavour:            feedlib.FlavourPro,
					Active:             true,
					ResponseType:       enums.SecurityQuestionResponseTypeText,
					Sequence:           1,
				}},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to create security questions",
			args: args{
				ctx: context.Background(),
				securityQuestions: []*domain.SecurityQuestion{{
					SecurityQuestionID: gofakeit.UUID(),
					QuestionStem:       gofakeit.Question(),
					Description:        gofakeit.BS(),
					Flavour:            feedlib.FlavourPro,
					Active:             true,
					ResponseType:       enums.SecurityQuestionResponseTypeText,
					Sequence:           1,
				}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case: failed to create security questions" {
				fakeGorm.MockCreateSecurityQuestionsFn = func(ctx context.Context, securityQuestions []*gorm.SecurityQuestion) ([]*gorm.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := d.CreateSecurityQuestions(tt.args.ctx, tt.args.securityQuestions)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateSecurityQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_CreateTermsOfService(t *testing.T) {
	dummyText := gofakeit.BS()
	type args struct {
		ctx            context.Context
		termsOfService *domain.TermsOfService
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create terms of service",
			args: args{
				ctx: context.Background(),
				termsOfService: &domain.TermsOfService{
					TermsID:   1,
					Text:      &dummyText,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case: failed to create terms of service",
			args: args{
				ctx: context.Background(),
				termsOfService: &domain.TermsOfService{
					TermsID:   1,
					Text:      &dummyText,
					ValidFrom: time.Now(),
					ValidTo:   time.Now(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Sad Case: failed to create terms of service" {
				fakeGorm.MockCreateTermsOfServiceFn = func(ctx context.Context, termsOfService *gorm.TermsOfService) (*gorm.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			_, err := d.CreateTermsOfService(tt.args.ctx, tt.args.termsOfService)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateTermsOfService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
