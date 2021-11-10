package terms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms/mock"
)

func TestTermsOfServiceImpl_GetCurrentTerms_Unittest(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		flavour enums.Flavour
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
				flavour: enums.PRO,
			},
			wantErr: false,
		},
		{
			name: "Sad case - empty flavour",
			args: args{
				ctx:     ctx,
				flavour: "",
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
				fakeDB.MockGetCurrentTermsFn = func(ctx context.Context, flavour enums.Flavour) (string, error) {
					return "", fmt.Errorf("an error occurred")
				}
			}

			_, err := j.GetCurrentTerms(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("TermsOfServiceImpl.GetCurrentTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
