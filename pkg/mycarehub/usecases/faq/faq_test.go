package faq_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/faq"
)

func TestUsecaseFAQImpl_GetFAQContent(t *testing.T) {
	ctx := context.Background()
	flavour := feedlib.FlavourConsumer
	limit := 10

	type args struct {
		ctx     context.Context
		limit   *int
		flavour feedlib.Flavour
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
				limit:   &limit,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "Happy case with limit 0",
			args: args{
				ctx:     ctx,
				flavour: flavour,
			},
			wantErr: false,
		},
		{
			name: "Invalid: invalid flavour",
			args: args{
				ctx:     ctx,
				limit:   &limit,
				flavour: feedlib.Flavour("invalidFlavour"),
			},
			wantErr: true,
		},
		{
			name: "Invalid: failed to get FAQ content",
			args: args{
				ctx:     ctx,
				limit:   &limit,
				flavour: flavour,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()

			f := faq.NewUsecaseFAQ(fakeDB)

			if tt.name == "Invalid: failed to get FAQ content" {
				fakeDB.MockGetFAQContentFn = func(ctx context.Context, flavour feedlib.Flavour, limit *int) ([]*domain.FAQ, error) {
					return nil, fmt.Errorf("failed to get FAQ content")
				}
			}

			got, err := f.GetFAQContent(tt.args.ctx, tt.args.flavour, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseFAQImpl.GetFAQContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
