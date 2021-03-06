package gorm_test

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

func TestPGInstance_InactivateFacility(t *testing.T) {

	ctx := context.Background()

	type args struct {
		ctx     context.Context
		mflCode *int
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
				mflCode: &mflCodeToInactivate,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.InactivateFacility(tt.args.ctx, tt.args.mflCode)
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
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		mflCode *int
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
				mflCode: &inactiveMflCode,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - empty mflCode",
			args: args{
				ctx:     ctx,
				mflCode: nil,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.ReactivateFacility(tt.args.ctx, tt.args.mflCode)
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

func TestPGInstance_SetNickname(t *testing.T) {
	ctx := context.Background()

	invalidUserID := ksuid.New().String()
	invalidNickname := gofakeit.HipsterSentence(50)

	type args struct {
		ctx      context.Context
		userID   *string
		nickname *string
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
				ctx:      ctx,
				userID:   &userID,
				nickname: &userNickname,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:      ctx,
				userID:   &invalidUserID,
				nickname: &userNickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:      ctx,
				userID:   nil,
				nickname: &userNickname,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no nickname",
			args: args{
				ctx:      ctx,
				userID:   &userID,
				nickname: &invalidNickname,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.SetNickName(tt.args.ctx, tt.args.userID, tt.args.nickname)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.SetNickName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.SetNickName() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestPGInstance_InvalidatePIN(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
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
				ctx:     ctx,
				userID:  userIDToInvalidate,
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: no user id provided",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := testingDB.InvalidatePIN(tt.args.ctx, tt.args.userID, tt.args.flavour)
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
	ctx := context.Background()

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
				ctx:                               ctx,
				userID:                            userID,
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: invalid user id",
			args: args{
				ctx:                               ctx,
				userID:                            uuid.New().String(),
				isCorrectSecurityQuestionResponse: true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid: empty user id",
			args: args{
				ctx:                               ctx,
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
	ctx := context.Background()

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
				ctx:     ctx,
				userID:  &userIDToAcceptTerms,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "sad case: missing args",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: no userID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: no terms",
			args: args{
				ctx:     ctx,
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
	ctx := context.Background()
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
				ctx:     ctx,
				userID:  userIDUpdatePinRequireChangeStatus,
				flavour: flavour,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - Empty userID and flavour",
			args: args{
				ctx:     ctx,
				userID:  "",
				flavour: "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - Invalid flavour",
			args: args{
				ctx:     ctx,
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

func TestPGInstance_UpdateClientCaregiver(t *testing.T) {
	ctx := context.Background()

	caregiverInput := dto.CaregiverInput{
		ClientID:      clientID,
		FirstName:     "Updated",
		LastName:      "Updated",
		PhoneNumber:   "+1234567890",
		CaregiverType: enums.CaregiverTypeMother,
	}

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
			name: "happy case",
			args: args{
				ctx:            ctx,
				caregiverInput: &caregiverInput,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testingDB.UpdateClientCaregiver(tt.args.ctx, tt.args.caregiverInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.UpdateClientCaregiver() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPGInstance_InProgressBy(t *testing.T) {
	ctx := context.Background()

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
				ctx:       ctx,
				requestID: clientsServiceRequestID,
				staffID:   staffID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - invalid request ID",
			args: args{
				ctx:       ctx,
				requestID: "clientsServiceRequestID",
				staffID:   staffID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid staff ID",
			args: args{
				ctx:       ctx,
				requestID: clientsServiceRequestID,
				staffID:   "staffID",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid request and  staff ID",
			args: args{
				ctx:       ctx,
				requestID: "clientsServiceRequestID",
				staffID:   "staffID",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid uuid for  staff ID",
			args: args{
				ctx:       ctx,
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
	ctx := context.Background()
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
				ctx:                ctx,
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
				ctx:                ctx,
				staffID:            &fakeString,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non-existent staff",
			args: args{
				ctx:                ctx,
				staffID:            &badUID,
				serviceRequestID:   &serviceRequestID,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid service request id",
			args: args{
				ctx:                ctx,
				staffID:            &staffID,
				serviceRequestID:   &fakeString,
				verificationStatus: enums.ServiceRequestStatusResolved.String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: non existent service request",
			args: args{
				ctx:                ctx,
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

func TestPGInstance_AssignRoles(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
		roles  []enums.UserRoleType
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
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid: invalid user ID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid: invalid role",
			args: args{
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleType("invalid"), enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.AssignRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.AssignRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.AssignRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_RevokeRoles(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		userID string
		roles  []enums.UserRoleType
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
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Invalid: invalid user ID",
			args: args{
				ctx:    ctx,
				userID: uuid.New().String(),
				roles:  []enums.UserRoleType{enums.UserRoleTypeSystemAdministrator, enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Invalid: invalid role",
			args: args{
				ctx:    ctx,
				userID: userID,
				roles:  []enums.UserRoleType{enums.UserRoleType("invalid"), enums.UserRoleTypeContentManagement},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := testingDB.RevokeRoles(tt.args.ctx, tt.args.userID, tt.args.roles)
			if (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.RevokeRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PGInstance.RevokeRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPGInstance_InvalidateScreeningToolResponse(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		clientID   string
		questionID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				clientID:   clientID,
				questionID: screeningToolsQuestionID,
			},
			wantErr: false,
		},
		{
			name: "Invalid: invalid client ID",
			args: args{
				ctx:        ctx,
				clientID:   uuid.New().String(),
				questionID: screeningToolsQuestionID,
			},
			wantErr: true,
		},
		{
			name: "Invalid: invalid question ID",
			args: args{
				ctx:        ctx,
				clientID:   clientID,
				questionID: "12345",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := testingDB.InvalidateScreeningToolResponse(tt.args.ctx, tt.args.clientID, tt.args.questionID); (err != nil) != tt.wantErr {
				t.Errorf("PGInstance.InvalidateScreeningToolResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPGInstance_UpdateServiceRequestsFromKenyaEMR(t *testing.T) {
	ctx := context.Background()
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("failed to initialize new PG instance: %v", err)
		return
	}

	serviceReq := &gorm.ClientServiceRequest{
		Base:           gorm.Base{},
		ID:             &serviceRequestID,
		Active:         true,
		RequestType:    "RED_FLAG",
		Request:        "VERY SAD",
		Status:         "IN PROGRESS",
		InProgressAt:   &time.Time{},
		ResolvedAt:     &time.Time{},
		ClientID:       clientID,
		InProgressByID: &staffID,
		OrganisationID: uuid.New().String(),
		ResolvedByID:   &staffID,
		FacilityID:     facilityID,
		Meta:           `{}`,
	}

	badServiceRequestID := "badServiceRequestID"
	invalidServiceReq := &gorm.ClientServiceRequest{
		ID: &badServiceRequestID,
	}

	err = pg.DB.Create(serviceReq).Error
	if err != nil {
		t.Errorf("Create securityQuestionResponse failed: %v", err)
		return
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
				ctx:     ctx,
				payload: []*gorm.ClientServiceRequest{serviceReq},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
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
				ctx: context.Background(),
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
				ctx: context.Background(),
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
				ctx:     context.Background(),
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
				ctx:     context.Background(),
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
				ctx: context.Background(),
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
				ctx:    context.Background(),
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
				ctx:    context.Background(),
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
	ctx := context.Background()

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
				ctx: ctx,
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
				ctx: ctx,
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
				ctx: ctx,
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
	ctx := context.Background()

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
				ctx:     ctx,
				userID:  userID2,
				flavour: "CONSUMER",
				status:  true,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
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
	ctx := context.Background()
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
				ctx:       ctx,
				userID:    userFailedSecurityCountID,
				failCount: 0,
			},
			wantErr: false,
		},
		{
			name: "Sad case: user not found",
			args: args{
				ctx:    ctx,
				userID: gofakeit.UUID(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid user ID",
			args: args{
				ctx:    ctx,
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
	ctx := context.Background()

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
				ctx: ctx,
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
				ctx: ctx,
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
	ctx := context.Background()

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
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				facility: &gorm.Facility{
					FacilityID: &facilityID,
				},
				updateData: map[string]interface{}{
					"fhir_organization_id": uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
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
	ctx := context.Background()

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
				ctx: ctx,
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
				ctx: ctx,
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
	ctx := context.Background()

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
				ctx: ctx,
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
				ctx: ctx,
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
	ctx := context.Background()
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
				ctx: ctx,
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
				ctx: ctx,
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
