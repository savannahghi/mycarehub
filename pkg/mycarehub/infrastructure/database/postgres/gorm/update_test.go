package gorm_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

func TestPGInstance_InactivateFacility(t *testing.T) {

	type args struct {
		ctx        context.Context
		identifier *gorm.FacilityIdentifier
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
				ctx: addRequiredContext(context.Background(), t),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: facilityIdentifierToInactivate,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: "53453434",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.InactivateFacility(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ReactivateFacility(t *testing.T) {

	type args struct {
		ctx        context.Context
		identifier *gorm.FacilityIdentifier
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
				ctx: addRequiredContext(context.Background(), t),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: inactiveFacilityIdentifier,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				identifier: &gorm.FacilityIdentifier{
					Type:  enums.FacilityIdentifierTypeMFLCode.String(),
					Value: "434343434",
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ReactivateFacility(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ReactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ReactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_InvalidatePIN(t *testing.T) {

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:    addRequiredContext(context.Background(), t),
				userID: userIDToInvalidate,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: no user id provided",
			args: args{
				ctx:    addRequiredContext(context.Background(), t),
				userID: "userID",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.InvalidatePIN(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InvalidatePIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.InvalidatePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateIsCorrectSecurityQuestionResponse(t *testing.T) {

	type args struct {
		ctx                               context.Context
		userID                            string
		isCorrectSecurityQuestionResponse bool
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                               addRequiredContext(context.Background(), t),
				userID:                            userID,
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid user id",
			args: args{
				ctx:                               addRequiredContext(context.Background(), t),
				userID:                            uuid.New().String(),
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: empty user id",
			args: args{
				ctx:                               addRequiredContext(context.Background(), t),
				userID:                            "",
				isCorrectSecurityQuestionResponse: true,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateIsCorrectSecurityQuestionResponse(tt.args.ctx, tt.args.userID, tt.args.isCorrectSecurityQuestionResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateIsCorrectSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateIsCorrectSecurityQuestionResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_AcceptTerms(t *testing.T) {

	type args struct {
		ctx     context.Context
		userID  *string
		termsID *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  &userIDToAcceptTerms,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing args",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: no userID",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  nil,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: no terms",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  &userIDToAcceptTerms,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_CompleteOnboardingTour(t *testing.T) {

	flavour := feedlib.FlavourConsumer

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name string

		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  userIDUpdatePinRequireChangeStatus,
				flavour: flavour,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - Empty userID and flavour",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  "",
				flavour: "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - Invalid flavour",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  userIDUpdatePinRequireChangeStatus,
				flavour: "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.CompleteOnboardingTour(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.CompleteOnboardingTour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.CompleteOnboardingTour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_InProgressBy(t *testing.T) {

	type args struct {
		ctx       context.Context
		requestID string
		staffID   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:       addRequiredContext(context.Background(), t),
				requestID: clientsServiceRequestID,
				staffID:   staffID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - invalid request ID",
			args: args{
				ctx:       addRequiredContext(context.Background(), t),
				requestID: "clientsServiceRequestID",
				staffID:   staffID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid staff ID",
			args: args{
				ctx:       addRequiredContext(context.Background(), t),
				requestID: clientsServiceRequestID,
				staffID:   "staffID",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid request and  staff ID",
			args: args{
				ctx:       addRequiredContext(context.Background(), t),
				requestID: "clientsServiceRequestID",
				staffID:   "staffID",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid uuid for  staff ID",
			args: args{
				ctx:       addRequiredContext(context.Background(), t),
				requestID: clientsServiceRequestID,
				staffID:   uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SetInProgressBy(tt.args.ctx, tt.args.requestID, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SetInProgressBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SetInProgressBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_ResolveStaffServiceRequest(t *testing.T) {

	fakeString := gofakeit.HipsterSentence(10)
	serviceRequestID := "26b20a42-cbb8-4553-aedb-c539602d04fc"
	badUID := "BadUID"

	type args struct {
		ctx                context.Context
		staffID            *string
		serviceRequestID   *string
		verificationStatus string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:                addRequiredContext(context.Background(), t),
				staffID:            &staffID,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case: invalid staff id",
			args: args{
				ctx:                addRequiredContext(context.Background(), t),
				staffID:            &fakeString,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non-existent staff",
			args: args{
				ctx:                addRequiredContext(context.Background(), t),
				staffID:            &badUID,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid service request id",
			args: args{
				ctx:                addRequiredContext(context.Background(), t),
				staffID:            &staffID,
				serviceRequestID:   &fakeString,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non existent service request",
			args: args{
				ctx:                addRequiredContext(context.Background(), t),
				staffID:            &staffID,
				serviceRequestID:   &badUID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ResolveStaffServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID, tt.args.verificationStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.ResolveStaffServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.ResolveStaffServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateServiceRequests(t *testing.T) {

	serviceReq := &gorm.ClientServiceRequest{
		ID: &clientServiceRequestIDToUpdate,
	}

	badServiceRequestID := "badServiceRequestID"
	invalidServiceReq := &gorm.ClientServiceRequest{
		ID: &badServiceRequestID,
	}

	type args struct {
		ctx     context.Context
		payload []*gorm.ClientServiceRequest
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				payload: []*gorm.ClientServiceRequest{serviceReq},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				payload: []*gorm.ClientServiceRequest{invalidServiceReq},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.UpdateServiceRequests(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.UpdateServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_UpdateAppointment(t *testing.T) {

	type args struct {
		ctx        context.Context
		payload    *gorm.Appointment
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update an appointment using id",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				payload: &gorm.Appointment{
					ID: appointmentID,
				},
				updateData: map[string]interface{}{
					"client_id": clientID,
					"reason":    "Knocked up",
					"date":      time.Now().Add(time.Duration(100)),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update non-existent appointment",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				payload: &gorm.Appointment{
					ID: gofakeit.UUID(),
				},
				updateData: map[string]interface{}{
					"client_id": clientID,
					"reason":    "Knocked up",
					"date":      time.Now().Add(time.Duration(100)),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: update appointment missing ids",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				payload: &gorm.Appointment{},
				updateData: map[string]interface{}{
					"client_id": clientID,
					"reason":    "Knocked up",
					"date":      time.Now().Add(time.Duration(100)),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.UpdateAppointment(tt.args.ctx, tt.args.payload, tt.args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateAppointment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.UpdateAppointment() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserPinUpdateRequiredStatus(t *testing.T) {
	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
		status  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update user pin update required status",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  userID2,
				flavour: feedlib.FlavourConsumer,
				status:  true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserPinUpdateRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserPinUpdateRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateClient(t *testing.T) {
	type args struct {
		ctx     context.Context
		client  *gorm.Client
		updates map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.Client
		wantErr bool
	}{
		{
			name: "Happy case: update client profile",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				client: &gorm.Client{
					ID: &clientID,
				},
				updates: map[string]interface{}{
					"fhir_patient_id": gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update client missing ID",
			args: args{
				ctx:    addRequiredContext(context.Background(), t),
				client: &gorm.Client{},
				updates: map[string]interface{}{
					"fhir_patient_id": gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: update client invalid field",
			args: args{
				ctx:    addRequiredContext(context.Background(), t),
				client: &gorm.Client{},
				updates: map[string]interface{}{
					"invalid_field_id": gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.UpdateClient(tt.args.ctx, tt.args.client, tt.args.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected client to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil && got.FHIRPatientID == nil {
				t.Errorf("expected client not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestPGInstance_UpdateHealthDiary(t *testing.T) {

	type args struct {
		ctx              context.Context
		healthDairyEntry *gorm.ClientHealthDiaryEntry
		updateData       map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				healthDairyEntry: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientsHealthDiaryEntryID,
					ClientID:                 clientID,
				},
				updateData: map[string]interface{}{
					"share_with_health_worker": true,
					"shared_at":                time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case - nil health diary entry ID",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				healthDairyEntry: &gorm.ClientHealthDiaryEntry{
					ClientID: clientID,
				},
				updateData: map[string]interface{}{
					"share_with_health_worker": true,
					"shared_at":                time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				healthDairyEntry: &gorm.ClientHealthDiaryEntry{
					ClientHealthDiaryEntryID: &clientsHealthDiaryEntryID,
					ClientID:                 "clientID",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testingDB.UpdateHealthDiary(tt.args.ctx, tt.args.healthDairyEntry, tt.args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateHealthDiary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_UpdateUserPinChangeRequiredStatus(t *testing.T) {

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
		status  bool
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  userID2,
				flavour: "CONSUMER",
				status:  true,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				userID:  "userID2",
				flavour: "CONSUMER",
				status:  true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserPinChangeRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserPinChangeRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateFailedSecurityQuestionsAnsweringAttempts(t *testing.T) {

	type args struct {
		ctx       context.Context
		userID    string
		failCount int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: reset failed security attempts",
			args: args{
				ctx:       addRequiredContext(context.Background(), t),
				userID:    userFailedSecurityCountID,
				failCount: 0,
			},
			wantErr: false,
		},
		{
			name: "Sad case: user not found",
			args: args{
				ctx:    addRequiredContext(context.Background(), t),
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid user ID",
			args: args{
				ctx:    addRequiredContext(context.Background(), t),
				userID: "32354",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateFailedSecurityQuestionsAnsweringAttempts(tt.args.ctx, tt.args.userID, tt.args.failCount); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateFailedSecurityQuestionsAnsweringAttempts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUser(t *testing.T) {

	invalidUserID := "invalid user"

	type args struct {
		ctx        context.Context
		user       *gorm.User
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				user: &gorm.User{
					UserID: &userID,
				},
				updateData: map[string]interface{}{
					"next_allowed_login": time.Now(),
					"failed_login_count": 0,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				user: &gorm.User{
					UserID: &invalidUserID,
				},
				updateData: map[string]interface{}{
					"next_allowed_login": time.Now(),
					"failed_login_count": 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUser(tt.args.ctx, tt.args.user, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateFacility(t *testing.T) {

	invalidFacilityID := "invalid facility"

	type args struct {
		ctx        context.Context
		facility   *gorm.Facility
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// {
		// 	name: "Happy case",
		// 	args: args{
		// 		ctx:  addRequiredContext(context.Background(), t),
		// 		facility: &gorm.Facility{
		// 			FacilityID: &facilityID,
		// 		},
		// 		updateData: map[string]interface{}{
		// 			"fhir_organization_id": uuid.New().String(),
		// 		},
		// 	},
		// 	wantErr: false,
		// },
		{
			name: "Sad case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				facility: &gorm.Facility{
					FacilityID: &invalidFacilityID,
				},
				updateData: map[string]interface{}{
					"fhir_organization_id": uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateFacility(tt.args.ctx, tt.args.facility, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateNotification(t *testing.T) {

	type args struct {
		ctx          context.Context
		notification *gorm.Notification
		updateData   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				notification: &gorm.Notification{
					ID: notificationID,
				},
				updateData: map[string]interface{}{
					"is_read": true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				notification: &gorm.Notification{
					ID: "invalid notification id",
				},
				updateData: map[string]interface{}{
					"is_read": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateNotification(tt.args.ctx, tt.args.notification, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserSurveys(t *testing.T) {

	type args struct {
		ctx        context.Context
		survey     *gorm.UserSurvey
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: update a user survey",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				survey: &gorm.UserSurvey{
					FormID:    "fe3b9c8e-f8e3-4f0f-b8b1-f8b8b8b8b8b8",
					ProjectID: 2,
					LinkID:    1,
				},
				updateData: map[string]interface{}{
					"has_submitted": true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update a user survey",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				survey: &gorm.UserSurvey{
					UserID: gofakeit.BeerName(),
				},
				updateData: map[string]interface{}{
					"has_submitted": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserSurveys(tt.args.ctx, tt.args.survey, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserSurveys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateClientServiceRequest(t *testing.T) {

	invalidClientServiceRequestID := "invalid client service request"
	type args struct {
		ctx                  context.Context
		clientServiceRequest *gorm.ClientServiceRequest
		updateData           map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				clientServiceRequest: &gorm.ClientServiceRequest{
					ID: &clientsServiceRequestID,
				},
				updateData: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: false,
		},
		{

			name: "Sad case",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				clientServiceRequest: &gorm.ClientServiceRequest{
					ID: &invalidClientServiceRequestID,
				},
				updateData: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := testingDB.UpdateClientServiceRequest(tt.args.ctx, tt.args.clientServiceRequest, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateClientServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateStaff(t *testing.T) {
	invalidID := "invalid"
	type args struct {
		ctx     context.Context
		staff   *gorm.StaffProfile
		updates map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *gorm.StaffProfile
		wantErr bool
	}{
		{
			name: "Happy case: update staff profile",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				staff:   &gorm.StaffProfile{ID: &staffID},
				updates: map[string]interface{}{"active": true},
			},
			want:    &gorm.StaffProfile{},
			wantErr: false,
		},
		{
			name: "Sad case: update staff profile",
			args: args{
				ctx:     addRequiredContext(context.Background(), t),
				staff:   &gorm.StaffProfile{ID: &invalidID},
				updates: map[string]interface{}{"active": true},
			},
			want:    &gorm.StaffProfile{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := testingDB.UpdateStaff(tt.args.ctx, tt.args.staff, tt.args.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_AddFacilitiesToStaffProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		staffID    string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: add new facility",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				staffID:    staffID,
				facilities: []string{facilityToAddToUserProfile},
			},
			wantErr: false,
		},
		{
			name: "Happy case: should not error when user has existing facility",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				staffID:    staffID,
				facilities: []string{facilityToAddToUserProfile, facilityID},
			},
			wantErr: false,
		},
		{
			name: "Sad case: Invalid Client ID",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				staffID:    "invalid",
				facilities: []string{facilityToAddToUserProfile},
			},
			wantErr: true,
		},
		{
			name: "Sad case: Invalid facility ID",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				staffID:    staffID,
				facilities: []string{facilityToAddToUserProfile, "Invalid"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.AddFacilitiesToStaffProfile(tt.args.ctx, tt.args.staffID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AddFacilitiesToStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_AddFacilitiesToClientProfile(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientID   string
		facilities []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: add new facility",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				clientID:   clientID,
				facilities: []string{facilityToAddToUserProfile},
			},
			wantErr: false,
		},
		{
			name: "Happy case: should not error when user has existing facility",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				clientID:   clientID,
				facilities: []string{facilityToAddToUserProfile, facilityID},
			},
			wantErr: false,
		},
		{
			name: "Sad case: Invalid Client ID",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				clientID:   "invalid",
				facilities: []string{facilityToAddToUserProfile},
			},
			wantErr: true,
		},
		{
			name: "Sad case: Invalid facility ID",
			args: args{
				ctx:        addRequiredContext(context.Background(), t),
				clientID:   clientID,
				facilities: []string{facilityToAddToUserProfile, "Invalid"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.AddFacilitiesToClientProfile(tt.args.ctx, tt.args.clientID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AddFacilitiesToClientProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateCaregiverClient(t *testing.T) {
	type args struct {
		ctx             context.Context
		caregiverClient *gorm.CaregiverClient
		updateData      map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update consent",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				caregiverClient: &gorm.CaregiverClient{
					ClientID: clientID,
				},
				updateData: map[string]interface{}{
					"client_consent":    enums.ConsentStateAccepted,
					"client_consent_at": time.Now(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update consent",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				caregiverClient: &gorm.CaregiverClient{
					ClientID: "clientID",
				},
				updateData: map[string]interface{}{
					"client_consent":    enums.ConsentStateAccepted,
					"client_consent_at": time.Now(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateCaregiverClient(tt.args.ctx, tt.args.caregiverClient, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateCaregiverClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateCaregiver(t *testing.T) {
	type args struct {
		ctx       context.Context
		caregiver *gorm.Caregiver
		updates   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update caregiver",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				caregiver: &gorm.Caregiver{
					ID: testCaregiverID,
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update field does not exist",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				caregiver: &gorm.Caregiver{
					ID: testCaregiverID,
				},
				updates: map[string]interface{}{
					"invalid": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateCaregiver(tt.args.ctx, tt.args.caregiver, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateCaregiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateUserContact(t *testing.T) {
	invalidID := "invalid"
	type args struct {
		ctx         context.Context
		userContact *gorm.Contact
		updates     map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update user contact",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				userContact: &gorm.Contact{
					UserID: &userID,
				},
				updates: map[string]interface{}{
					"contact_value": interserviceclient.TestUserPhoneNumber,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update user contact",
			args: args{
				ctx: addRequiredContext(context.Background(), t),
				userContact: &gorm.Contact{
					UserID: &invalidID,
				},
				updates: map[string]interface{}{
					"contact_value": interserviceclient.TestUserPhoneNumber,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateUserContact(tt.args.ctx, tt.args.userContact, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateUserContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateClientIdentifier(t *testing.T) {
	type args struct {
		ctx             context.Context
		clientID        string
		identifierType  string
		identifierValue string
		programID       string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update client identifier",
			args: args{
				ctx:             addRequiredContext(context.Background(), t),
				clientID:        clientID,
				identifierType:  "PHONE",
				identifierValue: interserviceclient.TestUserPhoneNumber,
				programID:       programID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update client identifier",
			args: args{
				ctx:             addRequiredContext(context.Background(), t),
				clientID:        "clientID",
				identifierType:  "PHONE",
				identifierValue: interserviceclient.TestUserPhoneNumber,
				programID:       "programID",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateClientIdentifier(tt.args.ctx, tt.args.clientID, tt.args.identifierType, tt.args.identifierValue, tt.args.programID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateClientIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateProgram(t *testing.T) {
	type args struct {
		ctx        context.Context
		program    *gorm.Program
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update program",
			args: args{
				ctx: context.Background(),
				program: &gorm.Program{
					ID: "4181df12-ca96-4f28-b78b-8e8ad88b25df",
				},
				updateData: map[string]interface{}{
					"fhir_organisation_id": gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update program",
			args: args{
				ctx: context.Background(),
				program: &gorm.Program{
					ID: "4181df12-ca96-4f28-b78b-8e8ad88b25df",
				},
				updateData: map[string]interface{}{
					"fhir_organissation_id": gofakeit.UUID(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateProgram(tt.args.ctx, tt.args.program, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateAuthorizationCode(t *testing.T) {
	type args struct {
		ctx        context.Context
		code       *gorm.AuthorizationCode
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: update an authorisation code",
			args: args{
				ctx: context.Background(),
				code: &gorm.AuthorizationCode{
					ID: oauthAuthorizationCode,
				},
				updateData: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid update field",
			args: args{
				ctx: context.Background(),
				code: &gorm.AuthorizationCode{
					ID: oauthAuthorizationCode,
				},
				updateData: map[string]interface{}{
					"activity": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateAuthorizationCode(tt.args.ctx, tt.args.code, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateAuthorizationCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateAccessToken(t *testing.T) {
	type args struct {
		ctx        context.Context
		code       *gorm.AccessToken
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: retrieve token",
			args: args{
				ctx: context.Background(),
				code: &gorm.AccessToken{
					ID: oauthAccessTokenOne,
				},
				updateData: map[string]interface{}{
					"active": false,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid id",
			args: args{
				ctx: context.Background(),
				code: &gorm.AccessToken{
					ID: "invalid",
				},
				updateData: map[string]interface{}{
					"active": false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateAccessToken(tt.args.ctx, tt.args.code, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateRefreshToken(t *testing.T) {
	type args struct {
		ctx        context.Context
		code       *gorm.RefreshToken
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: retrieve token",
			args: args{
				ctx: context.Background(),
				code: &gorm.RefreshToken{
					ID: oauthRefreshTokenOne,
				},
				updateData: map[string]interface{}{
					"active": false,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: invalid id",
			args: args{
				ctx: context.Background(),
				code: &gorm.RefreshToken{
					ID: "invalid",
				},
				updateData: map[string]interface{}{
					"active": false,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.UpdateRefreshToken(tt.args.ctx, tt.args.code, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
