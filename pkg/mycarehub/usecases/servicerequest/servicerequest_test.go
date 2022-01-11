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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()

			if tt.name == "Sad Case - Fail to create a service request" {
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *domain.ClientServiceRequest) error {
					return fmt.Errorf("failed to create service request")
				}
			}

			u := servicerequest.NewUseCaseServiceRequestImpl(fakeDB)
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
