package content_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content/mock"
)

func TestUsecaseContentImpl_ListContentCategories(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ContentItemCategory
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
			name: "Sad case",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewContentUsecaseMock()
			c := content.NewUsecaseContent(fakeDB)

			if tt.name == "Sad case" {
				fakeDB.MockListContentCategoriesFn = func(ctx context.Context) ([]*domain.ContentItemCategory, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := c.ListContentCategories(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseContentImpl.ListContentCategories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected content to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected content not to be nil for %v", tt.name)
				return
			}
		})
	}
}
