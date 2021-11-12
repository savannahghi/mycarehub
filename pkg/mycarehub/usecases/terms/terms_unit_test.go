package terms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms/mock"
)

func TestTermsOfServiceImpl_GetCurrentTerms_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
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
			},
			wantErr: false,
		},
		{
			name: "Sad case - empty flavour",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "Sad case - bad context",
			args: args{
				ctx: context.TODO(),
			},
			wantErr: true,
		},
		{
			name: "Sad case - nil context",
			args: args{
				ctx: nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewTermsUseCaseMock()

			j := terms.NewUseCasesTermsOfService(fakeDB)

			if tt.name == "Sad case - empty flavour" {
				fakeDB.MockGetCurrentTermsFn = func(ctx context.Context) (*domain.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - bad context" {
				fakeDB.MockGetCurrentTermsFn = func(ctx context.Context) (*domain.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - nil context" {
				fakeDB.MockGetCurrentTermsFn = func(ctx context.Context) (*domain.TermsOfService, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			_, err := j.GetCurrentTerms(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("TermsOfServiceImpl.GetCurrentTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
