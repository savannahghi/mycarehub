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
					SharedAt: currentTime,
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
					SharedAt: currentTime,
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

func TestMyCareHubDb_CreateNextOfKin(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx    context.Context
		person *dto.NextOfKinPayload
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
			fakeGorm.MockCreateRelatedPerson = func(ctx context.Context, person *gorm.RelatedPerson) error {
				return nil
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := d.CreateNextOfKin(tt.args.ctx, tt.args.person); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateNextOfKin() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateContact(t *testing.T) {
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
	}
	for _, tt := range tests {

		if tt.name == "Happy case: create a new contact" {
			fakeGorm.MockCreateContact = func(ctx context.Context, contact *gorm.Contact) error {
				return nil
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			if err := d.CreateContact(tt.args.ctx, tt.args.contact); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CreateAppointment(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx             context.Context
		appointment     domain.Appointment
		appointmentUUID string
		clientID        string
		staffID         string
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
					ID:       gofakeit.UUID(),
					Type:     "Dental",
					Status:   "COMPLETED",
					Reason:   "Knocked out",
					ClientID: gofakeit.UUID(),
				},
				appointmentUUID: gofakeit.UUID(),
				clientID:        gofakeit.UUID(),
				staffID:         gofakeit.UUID(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.CreateAppointment(tt.args.ctx, tt.args.appointment, tt.args.appointmentUUID, tt.args.clientID); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CreateAppointment() error = %v, wantErr %v", err, tt.wantErr)
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
