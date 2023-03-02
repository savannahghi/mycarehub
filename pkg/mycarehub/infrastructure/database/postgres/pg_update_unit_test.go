package postgres

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	gormMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm/mock"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/segmentio/ksuid"
)

func TestMyCareHubDb_InactivateFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		identifier *dto.FacilityIdentifierInput
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
				ctx: ctx,
				identifier: &dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "8893278372",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - empty mflCode" {
				fakeGorm.MockInactivateFacilityFn = func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}
			if tt.name == "Sad Case - very bad mflCode" {
				fakeGorm.MockInactivateFacilityFn = func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}

			got, err := d.InactivateFacility(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_ReactivateFacility(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		identifier *dto.FacilityIdentifierInput
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
				ctx: ctx,
				identifier: &dto.FacilityIdentifierInput{
					Type:  enums.FacilityIdentifierTypeMFLCode,
					Value: "88932783729",
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case - empty mflCode" {
				fakeGorm.MockReactivateFacilityFn = func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}
			if tt.name == "Sad Case - very bad mflCode" {
				fakeGorm.MockReactivateFacilityFn = func(ctx context.Context, identifier *gorm.FacilityIdentifier) (bool, error) {
					return false, fmt.Errorf("failed to inactivate facility")
				}
			}

			got, err := d.ReactivateFacility(tt.args.ctx, tt.args.identifier)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.InactivateFacility() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.InactivateFacility() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_AcceptTerms(t *testing.T) {
	ctx := context.Background()

	userID := ksuid.New().String()
	termsID := gofakeit.Number(0, 100000)
	negativeTermsID := gofakeit.Number(-100000, -1)

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
				userID:  &userID,
				termsID: &termsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no termsID",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &termsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - no userID and termsID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - negative termsID",
			args: args{
				ctx:     ctx,
				userID:  nil,
				termsID: &negativeTermsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - userID and negative termsID",
			args: args{
				ctx:     ctx,
				userID:  &userID,
				termsID: &negativeTermsID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - no userID and termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - negative termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - userID and negative termsID" {
				fakeGorm.MockAcceptTermsFn = func(ctx context.Context, userID *string, termsID *int) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_CompleteOnboardingTour(t *testing.T) {
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
			name: "Happy Case - Successfully change status",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to update",
			args: args{
				ctx:     ctx,
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Missing user id",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - No user id and flavour",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad Case - Fail to update" {
				fakeGorm.MockCompleteOnboardingTourFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to update status")
				}
			}

			if tt.name == "Sad Case - Missing user id" {
				fakeGorm.MockCompleteOnboardingTourFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to update status")
				}
			}
			if tt.name == "Sad Case - No user id and flavour" {
				fakeGorm.MockCompleteOnboardingTourFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to update status")
				}
			}

			got, err := d.CompleteOnboardingTour(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CompleteOnboardingTour() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CompleteOnboardingTour() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_InvalidatePIN(t *testing.T) {

	ctx := context.Background()
	userID := uuid.New().String()
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
				userID:  userID,
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
		{
			name: "invalid: invalid flavour",
			args: args{
				ctx:     ctx,
				flavour: "invalid-flavour",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.InvalidatePIN(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.InvalidatePIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.InvalidatePIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_UpdateIsCorrectSecurityQuestionResponse(t *testing.T) {

	ctx := context.Background()
	userID := uuid.New().String()

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
			name: "invalid: no user id provided",
			args: args{
				ctx: ctx,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			got, err := d.UpdateIsCorrectSecurityQuestionResponse(tt.args.ctx, tt.args.userID, tt.args.isCorrectSecurityQuestionResponse)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateIsCorrectSecurityQuestionResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.UpdateIsCorrectSecurityQuestionResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_InProgressBy(t *testing.T) {
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
				requestID: uuid.New().String(),
				staffID:   uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:       ctx,
				requestID: uuid.New().String(),
				staffID:   uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty request ID",
			args: args{
				ctx:       ctx,
				requestID: "",
				staffID:   uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - empty staff ID",
			args: args{
				ctx:       ctx,
				requestID: uuid.New().String(),
				staffID:   "",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			_ = pgMock.NewPostgresMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case" {
				fakeGorm.MockInProgressByFn = func(ctx context.Context, requestID, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty request ID" {
				fakeGorm.MockInProgressByFn = func(ctx context.Context, requestID, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty staff ID" {
				fakeGorm.MockInProgressByFn = func(ctx context.Context, requestID, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := d.SetInProgressBy(tt.args.ctx, tt.args.requestID, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.SetInProgressBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.SetInProgressBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_ResolveServiceRequest(t *testing.T) {
	testUUD := uuid.New().String()
	comment := "test comment"
	type args struct {
		ctx              context.Context
		staffID          *string
		serviceRequestID *string
		status           string
		action           []string
		comment          *string
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
				ctx:              context.Background(),
				staffID:          &testUUD,
				serviceRequestID: &testUUD,
				status:           enums.ServiceRequestStatusResolved.String(),
				action:           []string{"resolve"},
			},
			wantErr: false,
		},
		{
			name: "Happy case: no comment passed",
			args: args{
				ctx:              context.Background(),
				staffID:          &testUUD,
				serviceRequestID: &testUUD,
				status:           enums.ServiceRequestStatusResolved.String(),
				comment:          &comment,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get service request by ID",
			args: args{
				ctx:              context.Background(),
				staffID:          &testUUD,
				serviceRequestID: &testUUD,
				status:           enums.ServiceRequestStatusResolved.String(),
				comment:          &comment,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to convert json to map",
			args: args{
				ctx:              context.Background(),
				staffID:          &testUUD,
				serviceRequestID: &testUUD,
				status:           enums.ServiceRequestStatusResolved.String(),
				comment:          &comment,
			},
			wantErr: true,
		},
		{
			name: "Sad case: unable to update client service request",
			args: args{
				ctx:              context.Background(),
				staffID:          &testUUD,
				serviceRequestID: &testUUD,
				status:           enums.ServiceRequestStatusResolved.String(),
				comment:          &comment,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Happy case: no comment passed" {
				fakeGorm.MockGetServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
					return &gorm.ClientServiceRequest{
						Meta: "",
					}, nil
				}
			}
			if tt.name == "Sad case: unable to get service request by ID" {
				fakeGorm.MockGetServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: unable to convert json to map" {
				fakeGorm.MockGetServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
					return &gorm.ClientServiceRequest{
						Meta: `["yes","no"]`,
					}, nil
				}
			}
			if tt.name == "Sad case: unable to update client service request" {
				fakeGorm.MockGetServiceRequestByIDFn = func(ctx context.Context, serviceRequestID string) (*gorm.ClientServiceRequest, error) {
					return &gorm.ClientServiceRequest{
						Meta: ``,
					}, nil
				}
				fakeGorm.MockUpdateClientServiceRequestFn = func(ctx context.Context, clientServiceRequest *gorm.ClientServiceRequest, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}

			err := d.ResolveServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID, tt.args.status, tt.args.action, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ResolveServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_InvalidateScreeningToolResponse(t *testing.T) {
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
				ctx:        context.Background(),
				clientID:   uuid.New().String(),
				questionID: uuid.New().String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fakeGorm = gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if err := d.InvalidateScreeningToolResponse(tt.args.ctx, tt.args.clientID, tt.args.questionID); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.InvalidateScreeningToolResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateAppointment(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
	type args struct {
		ctx         context.Context
		appointment *domain.Appointment
		updateData  map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad case: invalid date format",
			args: args{
				ctx: context.Background(),
				appointment: &domain.Appointment{
					ExternalID: gofakeit.UUID(),
				},
				updateData: map[string]interface{}{
					"type":   "Dental",
					"status": "COMPLETED",
					"reason": "Knocked out",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.UpdateAppointment(tt.args.ctx, tt.args.appointment, tt.args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateAppointment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got == nil {
				t.Errorf("MyCareHubDb.UpdateAppointment() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateServiceRequests(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	currentTime := time.Now()
	inProgressBy := uuid.New().String()

	payload := domain.ServiceRequest{
		ID:           uuid.New().String(),
		RequestType:  gofakeit.BeerName(),
		Status:       "STATUS",
		InProgressAt: &currentTime,
		InProgressBy: &inProgressBy,
		ResolvedAt:   &currentTime,
		ResolvedBy:   &inProgressBy,
	}

	serviceReq := &domain.UpdateServiceRequestsPayload{
		ServiceRequests: []domain.ServiceRequest{
			payload,
		},
	}

	type args struct {
		ctx     context.Context
		payload *domain.UpdateServiceRequestsPayload
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
				payload: serviceReq,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				payload: serviceReq,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockUpdateServiceRequestsFn = func(ctx context.Context, payload []*gorm.ClientServiceRequest) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.UpdateServiceRequests(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.UpdateServiceRequests() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_UpdateUserPinUpdateRequiredStatus(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
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
			name: "Happy Case - Successfully update user pin update required status",
			args: args{
				ctx:     context.Background(),
				userID:  uuid.New().String(),
				flavour: feedlib.FlavourConsumer,
				status:  true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := d.UpdateUserPinUpdateRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateUserPinUpdateRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateClient(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	clientID := gofakeit.UUID()
	type args struct {
		ctx     context.Context
		client  *domain.ClientProfile
		updates map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy case: update client details",
			args: args{
				ctx: context.Background(),
				client: &domain.ClientProfile{
					ID: &clientID,
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: update client details error",
			args: args{
				ctx: context.Background(),
				client: &domain.ClientProfile{
					ID: &clientID,
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {

		if tt.name == "Sad case: update client details error" {
			fakeGorm.MockUpdateClientFn = func(ctx context.Context, client *gorm.Client, updates map[string]interface{}) (*gorm.Client, error) {
				return nil, fmt.Errorf("error cannot update client")
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			got, err := d.UpdateClient(tt.args.ctx, tt.args.client, tt.args.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateClient() error = %v, wantErr %v", err, tt.wantErr)
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

func TestMyCareHubDb_UpdateHealthDiary(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx              context.Context
		healthDairyEntry *domain.ClientHealthDiaryEntry
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
				healthDairyEntry: &domain.ClientHealthDiaryEntry{
					ID:       &UUID,
					ClientID: uuid.New().String(),
				},
				updateData: map[string]interface{}{
					"share_with_health_worker": true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				healthDairyEntry: &domain.ClientHealthDiaryEntry{
					ID:       &UUID,
					ClientID: uuid.New().String(),
				},
				updateData: map[string]interface{}{
					"share_with_health_worker": true,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty ids",
			args: args{
				ctx: ctx,
				healthDairyEntry: &domain.ClientHealthDiaryEntry{
					ClientID: "",
				},
				updateData: map[string]interface{}{
					"share_with_health_worker": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockUpdateHealthDiaryFn = func(ctx context.Context, clientHealthDiaryEntry *gorm.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty ids" {
				fakeGorm.MockUpdateHealthDiaryFn = func(ctx context.Context, clientHealthDiaryEntry *gorm.ClientHealthDiaryEntry, updateData map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			err := d.UpdateHealthDiary(tt.args.ctx, tt.args.healthDairyEntry, tt.args.updateData)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateHealthDiary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_UpdateFailedSecurityQuestionsAnsweringAttempts(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
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
			name: "Happy case",
			args: args{
				ctx:       context.Background(),
				userID:    uuid.New().String(),
				failCount: 5,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.UpdateFailedSecurityQuestionsAnsweringAttempts(tt.args.ctx, tt.args.userID, tt.args.failCount); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateFailedSecurityQuestionsAnsweringAttempts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateUser(t *testing.T) {
	ctx := context.Background()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	UUID := uuid.New().String()

	type args struct {
		ctx        context.Context
		user       *domain.User
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
				ctx:  ctx,
				user: &domain.User{ID: &UUID},
				updateData: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:  ctx,
				user: &domain.User{ID: &UUID},
				updateData: map[string]interface{}{
					"test": "test",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockUpdateUserFn = func(ctx context.Context, user *gorm.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user")
				}
			}
			if err := d.UpdateUser(tt.args.ctx, tt.args.user, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_ResolveStaffServiceRequest(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
				staffID:            &UUID,
				serviceRequestID:   &UUID,
				verificationStatus: "APPROVED",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:                ctx,
				staffID:            &UUID,
				serviceRequestID:   &UUID,
				verificationStatus: "REJECTED",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockResolveStaffServiceRequestFn = func(ctx context.Context, staffID, serviceRequestID *string, verificationStatus string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			got, err := d.ResolveStaffServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID, tt.args.verificationStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.ResolveStaffServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.ResolveStaffServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_UpdateUserPinChangeRequiredStatus(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

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
				userID:  UUID,
				flavour: feedlib.FlavourConsumer,
				status:  true,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				userID:  UUID,
				flavour: feedlib.FlavourConsumer,
				status:  true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockUpdateUserPinChangeRequiredStatusFn = func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
					return fmt.Errorf("failed to update user")
				}
			}

			if err := d.UpdateUserPinChangeRequiredStatus(tt.args.ctx, tt.args.userID, tt.args.flavour, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateUserPinChangeRequiredStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateFacility(t *testing.T) {
	ctx := context.Background()
	UUID := uuid.New().String()
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		facility   *domain.Facility
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
				facility: &domain.Facility{
					ID:                 &UUID,
					FHIROrganisationID: UUID,
				},
				updateData: map[string]interface{}{"name": "new name"},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				facility: &domain.Facility{
					ID:                 &UUID,
					FHIROrganisationID: UUID,
				},
				updateData: map[string]interface{}{"name": "new name"},
			},
			wantErr: true,
		},
		{
			name: "Sad case - no ID",
			args: args{
				ctx: ctx,
				facility: &domain.Facility{
					ID:                 &UUID,
					FHIROrganisationID: UUID,
				},
				updateData: map[string]interface{}{"name": "new name"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, facility *gorm.Facility, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update facility")
				}
			}
			if tt.name == "Sad case - no ID" {
				fakeGorm.MockUpdateFacilityFn = func(ctx context.Context, facility *gorm.Facility, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update facility")
				}
			}
			if err := d.UpdateFacility(tt.args.ctx, tt.args.facility, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateFacility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_CheckAppointmentExistsByExternalID(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		externalID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: check appointment exists",
			args: args{
				ctx:        context.Background(),
				externalID: uuid.New().String(),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := d.CheckAppointmentExistsByExternalID(tt.args.ctx, tt.args.externalID)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.CheckAppointmentExistsByExternalID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MyCareHubDb.CheckAppointmentExistsByExternalID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMyCareHubDb_UpdateNotification(t *testing.T) {
	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx          context.Context
		notification *domain.Notification
		updateData   map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: update a notification",
			args: args{
				ctx: context.Background(),
				notification: &domain.Notification{
					ID: gofakeit.UUID(),
				},
				updateData: map[string]interface{}{
					"is_read": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := d.UpdateNotification(tt.args.ctx, tt.args.notification, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateUserSurveys(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx        context.Context
		survey     *domain.UserSurvey
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
				survey: &domain.UserSurvey{
					ID:     uuid.New().String(),
					UserID: uuid.New().String(),
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
				survey: &domain.UserSurvey{
					ID:     uuid.New().String(),
					UserID: uuid.New().String(),
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
			if tt.name == "Sad case: unable to update a user survey" {
				fakeGorm.MockUpdateUserSurveysFn = func(ctx context.Context, survey *gorm.UserSurvey, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user survey")
				}
			}
			if err := d.UpdateUserSurveys(tt.args.ctx, tt.args.survey, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateUserSurveys() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateClientServiceRequest(t *testing.T) {
	ctx := context.Background()

	var fakeGorm = gormMock.NewGormMock()
	d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

	type args struct {
		ctx                  context.Context
		clientServiceRequest *domain.ServiceRequest
		updateData           map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: update a client service request",
			args: args{
				ctx: ctx,
				clientServiceRequest: &domain.ServiceRequest{
					ID: uuid.New().String(),
				},
				updateData: map[string]interface{}{
					"active": true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := d.UpdateClientServiceRequest(tt.args.ctx, tt.args.clientServiceRequest, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateClientServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateStaff(t *testing.T) {
	staffID := uuid.NewString()

	type args struct {
		ctx     context.Context
		staff   *domain.StaffProfile
		updates map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update staff details",
			args: args{
				ctx: context.Background(),
				staff: &domain.StaffProfile{
					ID: &staffID,
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: Failed to update staff, update error occurred",
			args: args{
				ctx: context.Background(),
				staff: &domain.StaffProfile{
					ID: &staffID,
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: Failed to update staff, update error occurred" {
				fakeGorm.MockUpdateStaffFn = func(ctx context.Context, staff *gorm.StaffProfile, updates map[string]interface{}) (*gorm.StaffProfile, error) {
					return nil, fmt.Errorf("error cannot update staff")
				}
			}

			err := d.UpdateStaff(tt.args.ctx, tt.args.staff, tt.args.updates)
			if (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateStaff() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestMyCareHubDb_AddFacilitiesToStaffProfile(t *testing.T) {
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
			name: "Happy case: add facilities to staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to add facilities to staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to add facilities to staff profile" {
				fakeGorm.MockAddFacilitiesToStaffProfileFn = func(ctx context.Context, staffID string, facilities []string) error {
					return fmt.Errorf("failed to add facilities to staff profile")
				}
			}
			if err := d.AddFacilitiesToStaffProfile(tt.args.ctx, tt.args.staffID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AddFacilitiesToStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_AddFacilitiesToClientProfile(t *testing.T) {
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
			name: "Happy case: add facilities to client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to add facilities to client profile",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: failed to add facilities to client profile" {
				fakeGorm.MockAddFacilitiesToClientProfileFn = func(ctx context.Context, clientID string, facilities []string) error {
					return fmt.Errorf("failed to add facilities to client profile")
				}
			}

			if err := d.AddFacilitiesToClientProfile(tt.args.ctx, tt.args.clientID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.AddFacilitiesToClientProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_RemoveFacilitiesFromStaffProfile(t *testing.T) {
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
			name: "Happy case: remove facilities from  staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to remove facilities from  staff profile",
			args: args{
				ctx:        context.Background(),
				staffID:    uuid.NewString(),
				facilities: []string{uuid.NewString()},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)
			if tt.name == "Sad case: failed to remove facilities from  staff profile" {
				fakeGorm.MockRemoveFacilitiesFromStaffProfileFn = func(ctx context.Context, staffID string, facilities []string) error {
					return fmt.Errorf("failed to remove facilities from staff profile")
				}
			}
			if err := d.RemoveFacilitiesFromStaffProfile(tt.args.ctx, tt.args.staffID, tt.args.facilities); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.RemoveFacilitiesFromStaffProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateCaregiverClient(t *testing.T) {
	type args struct {
		ctx             context.Context
		caregiverClient *domain.CaregiverClient
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
				ctx: context.Background(),
				caregiverClient: &domain.CaregiverClient{
					CaregiverID: uuid.NewString(),
					ClientID:    uuid.NewString(),
				},
				updateData: map[string]interface{}{
					"caregiver_consent_at": time.Now(),
					"caregiver_consent":    true,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update consent",
			args: args{
				ctx: context.Background(),
				caregiverClient: &domain.CaregiverClient{
					CaregiverID: uuid.NewString(),
					ClientID:    uuid.NewString(),
				},
				updateData: map[string]interface{}{
					"caregiver_consent_at": time.Now(),
					"caregiver_consent":    true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to update consent" {
				fakeGorm.MockUpdateCaregiverClientFn = func(ctx context.Context, caregiverClient *gorm.CaregiverClient, updates map[string]interface{}) error {
					return fmt.Errorf("unable to update caregiver consent")
				}
			}
			if err := d.UpdateCaregiverClient(tt.args.ctx, tt.args.caregiverClient, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateCaregiverClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateCaregiver(t *testing.T) {
	type args struct {
		ctx       context.Context
		caregiver *domain.CaregiverProfile
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
				ctx: context.Background(),
				caregiver: &domain.CaregiverProfile{
					ID: uuid.NewString(),
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to update caregiver",
			args: args{
				ctx: context.Background(),
				caregiver: &domain.CaregiverProfile{
					ID: uuid.NewString(),
				},
				updates: map[string]interface{}{
					"active": true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "sad case: failed to update caregiver" {
				fakeGorm.MockUpdateCaregiverFn = func(ctx context.Context, caregiver *gorm.Caregiver, updates map[string]interface{}) error {
					return fmt.Errorf("an error occurred")
				}
			}
			if err := d.UpdateCaregiver(tt.args.ctx, tt.args.caregiver, tt.args.updates); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateCaregiver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateUserContact(t *testing.T) {
	UID := uuid.NewString()
	type args struct {
		ctx        context.Context
		contact    *domain.Contact
		updateData map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: update user contact",
			args: args{
				ctx: context.Background(),
				contact: &domain.Contact{
					UserID: &UID,
				},
				updateData: map[string]interface{}{
					"contact_value": interserviceclient.TestUserPhoneNumber,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update user contact",
			args: args{
				ctx: context.Background(),
				contact: &domain.Contact{
					UserID: &UID,
				},
				updateData: map[string]interface{}{
					"contact_value": interserviceclient.TestUserPhoneNumber,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: unable to update user contact" {
				fakeGorm.MockUpdateUserContactFn = func(ctx context.Context, contact *gorm.Contact, updateData map[string]interface{}) error {
					return fmt.Errorf("unable to update user contact")
				}
			}

			if err := d.UpdateUserContact(tt.args.ctx, tt.args.contact, tt.args.updateData); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateUserContact() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMyCareHubDb_UpdateClientIdentifier(t *testing.T) {
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
				ctx:             context.Background(),
				clientID:        uuid.NewString(),
				identifierType:  "phone_number",
				identifierValue: "1234567890",
				programID:       uuid.NewString(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: update client identifier",
			args: args{
				ctx:             context.Background(),
				clientID:        uuid.NewString(),
				identifierType:  "phone_number",
				identifierValue: "1234567890",
				programID:       uuid.NewString(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGorm := gormMock.NewGormMock()
			d := NewMyCareHubDb(fakeGorm, fakeGorm, fakeGorm, fakeGorm)

			if tt.name == "Sad case: update client identifier" {
				fakeGorm.MockUpdateClientIdentifierFn = func(ctx context.Context, clientID string, identifierType string, identifierValue string, programID string) error {
					return fmt.Errorf("unable to update client identifier")
				}
			}

			if err := d.UpdateClientIdentifier(tt.args.ctx, tt.args.clientID, tt.args.identifierType, tt.args.identifierValue, tt.args.programID); (err != nil) != tt.wantErr {
				t.Errorf("MyCareHubDb.UpdateClientIdentifier() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
