package servicerequest_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
	userMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func TestUseCasesServiceRequestImpl_CreateServiceRequest(t *testing.T) {
	type args struct {
		ctx                 context.Context
		serviceRequestInput *dto.ServiceRequestInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create a service request",
			args: args{
				ctx: context.Background(),
				serviceRequestInput: &dto.ServiceRequestInput{
					ClientID:    uuid.New().String(),
					RequestType: "HEALTH_DIARY_ENTRY",
					Request:     "A random request",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create a service request",
			args: args{
				ctx: context.Background(),
				serviceRequestInput: &dto.ServiceRequestInput{
					ClientID:    uuid.New().String(),
					RequestType: "HEALTH_DIARY_ENTRY",
					Request:     "A random request",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get client profile",
			args: args{
				ctx: context.Background(),
				serviceRequestInput: &dto.ServiceRequestInput{
					ClientID:    uuid.New().String(),
					RequestType: "HEALTH_DIARY_ENTRY",
					Request:     "A random request",
					Flavour:     feedlib.FlavourConsumer,
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff profile by staff ID",
			args: args{
				ctx: context.Background(),
				serviceRequestInput: &dto.ServiceRequestInput{
					ClientID:    uuid.New().String(),
					RequestType: "PIN_RESET",
					Request:     "A random request",
					Flavour:     feedlib.FlavourPro,
					StaffID:     uuid.New().String(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - No staff ID",
			args: args{
				ctx: context.Background(),
				serviceRequestInput: &dto.ServiceRequestInput{
					ClientID:    uuid.New().String(),
					RequestType: "PIN_RESET",
					Request:     "A random request",
					Flavour:     feedlib.FlavourPro,
					StaffID:     "",
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - No client ID",
			args: args{
				ctx: context.Background(),
				serviceRequestInput: &dto.ServiceRequestInput{
					ClientID:    "",
					RequestType: "PIN_RESET",
					Request:     "A random request",
					Flavour:     feedlib.FlavourConsumer,
					StaffID:     uuid.New().String(),
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()

			if tt.name == "Sad Case - Fail to create a service request" {
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			if tt.name == "Sad Case - Unable to get client profile" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}
			if tt.name == "Sad Case - Fail to get staff profile by staff ID" {
				fakeDB.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred while getting staff profile")
				}
			}
			if tt.name == "Sad Case - No staff ID" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}
			if tt.name == "Sad Case - No client IDD" {
				fakeDB.MockGetStaffProfileByStaffIDFn = func(ctx context.Context, staffID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("an error occurred while getting staff profile")
				}
			}

			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)
			got, err := u.CreateServiceRequest(tt.args.ctx, tt.args.serviceRequestInput)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.CreateServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.CreateServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_InProgressBy(t *testing.T) {
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
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad case" {
				fakeDB.MockInProgressByFn = func(ctx context.Context, requestID, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty request ID" {
				fakeDB.MockInProgressByFn = func(ctx context.Context, requestID, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty staff ID" {
				fakeDB.MockInProgressByFn = func(ctx context.Context, requestID, staffID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			got, err := u.SetInProgressBy(tt.args.ctx, tt.args.requestID, tt.args.staffID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.SetInProgressBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.SetInProgressBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_GetServiceRequests(t *testing.T) {
	invalidRequestType := "invalid"
	invalidStatus := "invalid"
	facilityID := uuid.New().String()
	type args struct {
		ctx           context.Context
		requestType   *string
		requestStatus *string
		facilityID    string
		flavour       feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get service requests type, invalid type",
			args: args{
				ctx:         context.Background(),
				requestType: &invalidRequestType,
				facilityID:  facilityID,
				flavour:     feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests status, invalid status",
			args: args{
				ctx:           context.Background(),
				requestStatus: &invalidStatus,
				facilityID:    facilityID,
				flavour:       feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad Case - Fail to get service requests" {
				fakeDB.MockGetServiceRequestsFn = func(ctx context.Context, requestType, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service requests")
				}
			}
			_, err := u.GetServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.requestStatus, tt.args.facilityID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.GetServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if got != tt.want {
			// 	t.Errorf("UseCasesServiceRequestImpl.CreateServiceRequest() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestUseCasesServiceRequestImpl_ResolveServiceRequest(t *testing.T) {
	testID := uuid.New().String()
	type args struct {
		ctx              context.Context
		staffID          *string
		serviceRequestID *string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully resolve service request",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Failed to get service request by id",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to get client profile by client id",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Failed to update security question answering attempts",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - no staff ID present",
			args: args{
				ctx:              context.Background(),
				staffID:          nil,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - no service request ID present",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: nil,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to resolve service request",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to resolve service request, return false",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user",
			args: args{
				ctx:              context.Background(),
				staffID:          &testID,
				serviceRequestID: &testID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad Case - Fail to resolve service request" {
				fakeDB.MockResolveServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string, status string) (bool, error) {
					return false, fmt.Errorf("failed to resolve service request")
				}
			}
			if tt.name == "Sad Case - Fail to resolve service request, return false" {
				fakeDB.MockResolveServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string, status string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad Case - Failed to get service request by id" {
				fakeDB.MockGetServiceRequestByIDFn = func(ctx context.Context, id string) (*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service request by id")
				}
			}

			if tt.name == "Sad Case - Failed to get client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, id string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by client id")
				}
			}
			if tt.name == "Sad Case - Fail to update user" {
				fakeDB.MockGetServiceRequestByIDFn = func(ctx context.Context, id string) (*domain.ServiceRequest, error) {
					return &domain.ServiceRequest{
						ID:             testID,
						RequestType:    enums.ServiceRequestTypePinReset.String(),
						Request:        gofakeit.Sentence(5),
						Status:         enums.ServiceRequestStatusPending.String(),
						Active:         false,
						ClientID:       testID,
						CreatedAt:      time.Time{},
						InProgressAt:   &time.Time{},
						InProgressBy:   new(string),
						ResolvedAt:     &time.Time{},
						ResolvedBy:     new(string),
						ResolvedByName: new(string),
						FacilityID:     testID,
						ClientName:     new(string),
						ClientContact:  new(string),
						Meta:           map[string]interface{}{},
					}, nil
				}

				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user")
				}
			}

			if tt.name == "Sad Case - Failed to update security question answering attempts" {
				// pin request service request
				fakeDB.MockGetServiceRequestByIDFn = func(ctx context.Context, id string) (*domain.ServiceRequest, error) {
					return &domain.ServiceRequest{
						ID:             testID,
						RequestType:    enums.ServiceRequestTypePinReset.String(),
						Request:        gofakeit.Sentence(5),
						Status:         enums.ServiceRequestStatusPending.String(),
						Active:         false,
						ClientID:       testID,
						CreatedAt:      time.Time{},
						InProgressAt:   &time.Time{},
						InProgressBy:   new(string),
						ResolvedAt:     &time.Time{},
						ResolvedBy:     new(string),
						ResolvedByName: new(string),
						FacilityID:     testID,
						ClientName:     new(string),
						ClientContact:  new(string),
						Meta:           map[string]interface{}{},
					}, nil
				}
				fakeDB.MockUpdateFailedSecurityQuestionsAnsweringAttemptsFn = func(ctx context.Context, userID string, failCount int) error {
					return fmt.Errorf("failed to update security question answering attempts")

				}
			}

			got, err := u.ResolveServiceRequest(tt.args.ctx, tt.args.staffID, tt.args.serviceRequestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.ResolveServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.ResolveServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_GetPendingServiceRequestsCount(t *testing.T) {
	ctx := context.Background()
	facilityID := uuid.New().String()
	type args struct {
		ctx        context.Context
		facilityID *string
		flavour    feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *int64
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:        ctx,
				facilityID: &facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        ctx,
				facilityID: &facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty facility id",
			args: args{
				ctx:        ctx,
				facilityID: &facilityID,
				flavour:    feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:        ctx,
				facilityID: &facilityID,
				flavour:    "invalid-flavour",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad case" {
				fakeDB.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty facility id" {
				fakeDB.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - invalid flavour" {
				fakeDB.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCountResponse, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := u.GetPendingServiceRequestsCount(tt.args.ctx, *tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.GetPendingServiceRequestsCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("PGInstance.GetPendingServiceRequestsCount() = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("PGInstance.GetPendingServiceRequestsCount() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_GetServiceRequestsForKenyaEMR(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeUser := userMock.NewUserUseCaseMock()
	u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

	currentTime := time.Now()

	type args struct {
		ctx     context.Context
		payload *dto.ServiceRequestPayload
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ServiceRequest
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      1234,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode:      1234,
					LastSyncTime: &currentTime,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case - Bad Input",
			args: args{
				ctx: ctx,
				payload: &dto.ServiceRequestPayload{
					MFLCode: 0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockGetServiceRequestsForKenyaEMRFn = func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - Bad Input" {
				fakeDB.MockGetServiceRequestsForKenyaEMRFn = func(ctx context.Context, payload *dto.ServiceRequestPayload) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			got, err := u.GetServiceRequestsForKenyaEMR(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.GetServiceRequestsForKenyaEMR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("UseCasesServiceRequestImpl.GetServiceRequestsForKenyaEMR = %v, want %v", got, tt.want)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesServiceRequestImpl.GetServiceRequestsForKenyaEMR = %v, want %v", got, tt.want)
				return
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_UpdateServiceRequestsFromKenyaEMR(t *testing.T) {
	ctx := context.Background()
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeUser := userMock.NewUserUseCaseMock()
	u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

	payload := dto.UpdateServiceRequestPayload{
		ID:           uuid.New().String(),
		RequestType:  gofakeit.BeerName(),
		Status:       "STATUS",
		InProgressAt: time.Time{},
		InProgressBy: uuid.New().String(),
		ResolvedAt:   time.Time{},
		ResolvedBy:   uuid.New().String(),
	}

	serviceReq := &dto.UpdateServiceRequestsPayload{
		ServiceRequests: []dto.UpdateServiceRequestPayload{
			payload,
		},
	}

	type args struct {
		ctx     context.Context
		payload *dto.UpdateServiceRequestsPayload
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
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := u.UpdateServiceRequestsFromKenyaEMR(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.UpdateServiceRequestsFromKenyaEMR() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.UpdateServiceRequestsFromKenyaEMR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_CreatePinResetServiceRequest(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		phoneNumber string
		cccNumber   string
		flavour     feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully create service request",
			args: args{
				ctx:         ctx,
				phoneNumber: "12345",
				cccNumber:   "12345",
				flavour:     feedlib.FlavourPro,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create service request",
			args: args{
				ctx:         ctx,
				phoneNumber: "12345",
				cccNumber:   "12345",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by phonenumber",
			args: args{
				ctx:         ctx,
				phoneNumber: "12345",
				cccNumber:   "12345",
				flavour:     feedlib.FlavourPro,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile by user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: "12345",
				cccNumber:   "12345",
				flavour:     feedlib.FlavourConsumer,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad Case - Fail to create service request" {
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			if tt.name == "Sad Case - Fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phonenumber")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile by user ID" {
				fakeDB.MockGetClientProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by user id")
				}
			}

			got, err := u.CreatePinResetServiceRequest(tt.args.ctx, tt.args.phoneNumber, tt.args.cccNumber, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.CreatePinResetServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.CreatePinResetServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_VerifyPinResetServiceRequest(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                      context.Context
		clientID                 string
		serviceRequestID         string
		cccNumber                string
		phoneNumber              string
		physicalIdentityVerified bool
		state                    string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully approve pin reset service request",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
				state:                    enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully reject pin reset service request",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
				state:                    enums.VerifyServiceRequestStateRejected.String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get logged in user",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff profile by user ID",
			args: args{
				ctx:      ctx,
				clientID: uuid.New().String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Patient not verified by healthcare worker",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: false,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get ccc number",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to mark service request as in progress",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
				state:                    enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get useer profile by phonenumber",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to generate temporary pin",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
				state:                    enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user pin changed required status",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
				state:                    enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to resolve service request",
			args: args{
				ctx:                      ctx,
				clientID:                 "26b20a42-cbb8-4553-aedb-c539602d04fc",
				serviceRequestID:         uuid.New().String(),
				cccNumber:                "123456",
				phoneNumber:              "+254711111111",
				physicalIdentityVerified: true,
				state:                    enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			fakeServiceRequest := mock.NewServiceRequestUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad Case - Fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user UID")
				}
			}

			if tt.name == "Sad Case - Fail to get staff profile by user ID" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}

			if tt.name == "Sad Case - Patient not verified by healthcare worker" {
				fakeServiceRequest.MockVerifyClientPinResetServiceRequestFn = func(
					ctx context.Context,
					clientID string,
					serviceRequestID string,
					cccNumber string,
					phoneNumber string,
					physicalIdentityVerified bool,
					state string,
				) (bool, error) {
					return false, fmt.Errorf("patient not verified")
				}
			}

			if tt.name == "Sad Case - Fail to get ccc number" {
				fakeDB.MockGetClientCCCIdentifier = func(ctx context.Context, clientID string) (*domain.Identifier, error) {
					return nil, fmt.Errorf("fail to get client ccc number")
				}
			}

			if tt.name == "Sad Case - Fail to mark service request as in progress" {
				fakeDB.MockInProgressByFn = func(ctx context.Context, requestID string, staffID string) (bool, error) {
					return false, fmt.Errorf("failed to mark service request as in progress")
				}
			}

			if tt.name == "Sad Case - Fail to get useer profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone")
				}
			}

			if tt.name == "Sad Case - Fail to generate temporary pin" {
				fakeUser.MockGenerateTemporaryPinFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("failed to generate temporary pin")
				}
			}

			if tt.name == "Sad Case - Fail to update user pin changed required status" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user pin changed required status")
				}
			}

			if tt.name == "Sad Case - Fail to resolve service request" {
				fakeDB.MockResolveServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string, status string) (bool, error) {
					return false, fmt.Errorf("failed to resolve service request")
				}
			}

			got, err := u.VerifyClientPinResetServiceRequest(tt.args.ctx, tt.args.clientID, tt.args.serviceRequestID, tt.args.cccNumber, tt.args.phoneNumber, tt.args.physicalIdentityVerified, tt.args.state)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.VerifyPinResetServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.VerifyPinResetServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesServiceRequestImpl_VerifyStaffPinResetServiceRequest(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakeUser := userMock.NewUserUseCaseMock()
	_ = mock.NewServiceRequestUseCaseMock()
	u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

	type args struct {
		ctx                context.Context
		phoneNumber        string
		serviceRequestID   string
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
				phoneNumber:        uuid.New().String(),
				serviceRequestID:   uuid.New().String(),
				verificationStatus: "APPROVED",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Happy Case - Successfully reject pin reset service request",
			args: args{
				ctx:                ctx,
				serviceRequestID:   uuid.New().String(),
				phoneNumber:        "+254711111111",
				verificationStatus: enums.VerifyServiceRequestStateRejected.String(),
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:                ctx,
				phoneNumber:        uuid.New().String(),
				serviceRequestID:   uuid.New().String(),
				verificationStatus: "APPROVED",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get logged in user",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254711111111",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get staff profile by user ID",
			args: args{
				ctx:         ctx,
				phoneNumber: "+254711111111",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to mark service request as in progress",
			args: args{
				ctx:                ctx,
				serviceRequestID:   uuid.New().String(),
				phoneNumber:        "+254711111111",
				verificationStatus: enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get user profile by phonenumber",
			args: args{
				ctx:              ctx,
				serviceRequestID: uuid.New().String(),
				phoneNumber:      "+254711111111",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to generate temporary pin",
			args: args{
				ctx:                ctx,
				serviceRequestID:   uuid.New().String(),
				phoneNumber:        "+254711111111",
				verificationStatus: enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to update user pin changed required status",
			args: args{
				ctx:                ctx,
				serviceRequestID:   uuid.New().String(),
				phoneNumber:        "+254711111111",
				verificationStatus: enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to resolve service request",
			args: args{
				ctx:                ctx,
				serviceRequestID:   uuid.New().String(),
				phoneNumber:        "+254711111111",
				verificationStatus: enums.VerifyServiceRequestStateApproved.String(),
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case" {
				fakeDB.MockResolveStaffServiceRequestFn = func(ctx context.Context, staffID, serviceRequestID *string, verificationStatus string) (bool, error) {
					return false, fmt.Errorf("failed to resolve service request")
				}
			}
			if tt.name == "Sad Case - Fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user UID")
				}
			}
			if tt.name == "Sad Case - Fail to get staff profile by user ID" {
				fakeDB.MockGetStaffProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.StaffProfile, error) {
					return nil, fmt.Errorf("failed to get staff profile by user ID")
				}
			}
			if tt.name == "Sad Case - Fail to mark service request as in progress" {
				fakeDB.MockInProgressByFn = func(ctx context.Context, requestID string, staffID string) (bool, error) {
					return false, fmt.Errorf("failed to mark service request as in progress")
				}
			}
			if tt.name == "Sad Case - Fail to get user profile by phonenumber" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phonenumber")
				}
			}
			if tt.name == "Sad Case - Fail to generate temporary pin" {
				fakeUser.MockGenerateTemporaryPinFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (string, error) {
					return "", fmt.Errorf("failed to generate temporary pin")
				}
			}
			if tt.name == "Sad Case - Fail to update user pin changed required status" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user pin changed required status")
				}
			}
			if tt.name == "Sad Case - Fail to resolve service request" {
				fakeDB.MockResolveStaffServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string, status string) (bool, error) {
					return false, fmt.Errorf("failed to resolve service request")
				}
			}
			got, err := u.VerifyStaffPinResetServiceRequest(tt.args.ctx, tt.args.phoneNumber, tt.args.serviceRequestID, tt.args.verificationStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.VerifyStaffPinResetServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.VerifyStaffPinResetServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
