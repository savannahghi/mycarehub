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
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
	userMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func TestUseCasesServiceRequestImpl_CreateServiceRequest(t *testing.T) {
	type args struct {
		ctx         context.Context
		clientID    string
		requestType string
		request     string
		cccNumber   string
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
				ctx:         context.Background(),
				clientID:    uuid.New().String(),
				requestType: "HEALTH_DIARY_ENTRY",
				request:     "A random request",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create a service request",
			args: args{
				ctx:         context.Background(),
				clientID:    uuid.New().String(),
				requestType: "HEALTH_DIARY_ENTRY",
				request:     "A random request",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Unable to get client profile",
			args: args{
				ctx:         context.Background(),
				clientID:    uuid.New().String(),
				requestType: "HEALTH_DIARY_ENTRY",
				request:     "A random request",
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
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			if tt.name == "Sad Case - Unable to get client profile" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile")
				}
			}

			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)
			got, err := u.CreateServiceRequest(tt.args.ctx, tt.args.clientID, tt.args.requestType, tt.args.request, tt.args.cccNumber)
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
		facilityID    *string
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
				facilityID: &facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get service requests type, invalid type",
			args: args{
				ctx:         context.Background(),
				requestType: &invalidRequestType,
				facilityID:  &facilityID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests status, invalid status",
			args: args{
				ctx:           context.Background(),
				requestStatus: &invalidStatus,
				facilityID:    &facilityID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: &facilityID,
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
				fakeDB.MockGetServiceRequestsFn = func(ctx context.Context, requestType, requestStatus, facilityID *string) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service requests")
				}
			}
			_, err := u.GetServiceRequests(tt.args.ctx, tt.args.requestType, tt.args.requestStatus, tt.args.facilityID)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakeUser := userMock.NewUserUseCaseMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB, fakeExtension, fakeUser)

			if tt.name == "Sad Case - Fail to resolve service request" {
				fakeDB.MockResolveServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
					return false, fmt.Errorf("failed to resolve service request")
				}
			}
			if tt.name == "Sad Case - Fail to resolve service request, return false" {
				fakeDB.MockResolveServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
					return false, nil
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
		facilityID string
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
				facilityID: facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
			},
			wantErr: true,
		},
		{
			name: "Sad case - empty facility id",
			args: args{
				ctx:        ctx,
				facilityID: facilityID,
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
				fakeDB.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - empty facility id" {
				fakeDB.MockGetPendingServiceRequestsCountFn = func(ctx context.Context, facilityID string) (*domain.ServiceRequestsCount, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := u.GetPendingServiceRequestsCount(tt.args.ctx, tt.args.facilityID)
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
		ctx       context.Context
		cccNumber string
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
				ctx:       ctx,
				cccNumber: "12345",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to create service request",
			args: args{
				ctx:       ctx,
				cccNumber: "12345",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get client profile by ccc number",
			args: args{
				ctx:       ctx,
				cccNumber: "12345",
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
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			if tt.name == "Sad Case - Fail to get client profile by ccc number" {
				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client profile by ccc number")
				}
			}

			got, err := u.CreatePinResetServiceRequest(tt.args.ctx, tt.args.cccNumber)
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

func TestUseCasesServiceRequestImpl_ApprovePinResetServiceRequest(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx                      context.Context
		clientID                 string
		serviceRequestID         string
		cccNumber                string
		phoneNumber              string
		physicalIdentityVerified bool
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
			name: "Sad Case - Fail to invite user",
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
			name: "Sad Case - Fail to update user pin changed required status",
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
			name: "Sad Case - Fail to resolve service request",
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
				fakeServiceRequest.MockApprovePinResetServiceRequestFn = func(
					ctx context.Context,
					clientID string,
					serviceRequestID string,
					cccNumber string,
					phoneNumber string,
					physicalIdentityVerified bool,
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

			if tt.name == "Sad Case - Fail to invite user" {
				fakeUser.MockInviteUserFn = func(ctx context.Context, userID string, phoneNumber string, flavour feedlib.Flavour) (bool, error) {
					return false, fmt.Errorf("failed to invite user")
				}
			}

			if tt.name == "Sad Case - Fail to update user pin changed required status" {
				fakeDB.MockUpdateUserPinChangeRequiredStatusFn = func(ctx context.Context, userID string, flavour feedlib.Flavour, status bool) error {
					return fmt.Errorf("failed to update user pin changed required status")
				}
			}

			if tt.name == "Sad Case - Fail to resolve service request" {
				fakeDB.MockResolveServiceRequestFn = func(ctx context.Context, staffID *string, serviceRequestID *string) (bool, error) {
					return false, fmt.Errorf("failed to resolve service request")
				}
			}

			got, err := u.ApprovePinResetServiceRequest(tt.args.ctx, tt.args.clientID, tt.args.serviceRequestID, tt.args.cccNumber, tt.args.phoneNumber, tt.args.physicalIdentityVerified)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesServiceRequestImpl.ApprovePinResetServiceRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesServiceRequestImpl.ApprovePinResetServiceRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
