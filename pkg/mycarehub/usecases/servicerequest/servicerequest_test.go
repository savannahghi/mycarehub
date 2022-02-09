package servicerequest_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
)

func TestUseCasesServiceRequestImpl_CreateServiceRequest(t *testing.T) {
	type args struct {
		ctx         context.Context
		clientID    string
		requestType string
		request     string
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

			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB)
			got, err := u.CreateServiceRequest(tt.args.ctx, tt.args.clientID, tt.args.requestType, tt.args.request)
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
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB)

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
				ctx: context.Background(),
				facilityID: &facilityID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get service requests type, invalid type",
			args: args{
				ctx:         context.Background(),
				requestType: &invalidRequestType,
				facilityID: &facilityID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests status, invalid status",
			args: args{
				ctx:           context.Background(),
				requestStatus: &invalidStatus,
				facilityID: &facilityID,
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get service requests",
			args: args{
				ctx: context.Background(),
				facilityID: &facilityID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB)

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
			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB, fakeDB, fakeDB)

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
