package communities_test

import (
	"context"
	"fmt"
	"testing"

	stream "github.com/GetStream/stream-chat-go/v5"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
)

func TestUseCasesCommunitiesImpl_ListGetStreamUsers(t *testing.T) {
	type args struct {
		ctx   context.Context
		input *domain.QueryOption
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully list stream users",
			args: args{
				ctx: context.Background(),
				input: &domain.QueryOption{
					Filter: map[string]interface{}{
						"role": "user",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to list stream users",
			args: args{
				ctx: context.Background(),
				input: &domain.QueryOption{
					Filter: map[string]interface{}{
						"role": "user",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			communities := communities.NewUseCaseCommunitiesImpl(fakeGetStream)

			if tt.name == "Sad Case - Fail to list stream users" {
				fakeGetStream.MockListGetStreamUsersFn = func(ctx context.Context, queryOptions *stream.QueryOption) (*stream.QueryUsersResponse, error) {
					return nil, fmt.Errorf("failed to get users")
				}
			}

			got, err := communities.ListGetStreamUsers(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesCommunitiesImpl.ListGetStreamUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected a response but got: %v", got)
				return
			}
		})
	}
}
