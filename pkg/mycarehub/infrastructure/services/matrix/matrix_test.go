package matrix_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix"
)

func TestMatrix_CreateCommunity(t *testing.T) {
	type args struct {
		ctx  context.Context
		room *dto.CommunityInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: Create room successfully",
			args: args{
				ctx: context.Background(),
				room: &dto.CommunityInput{
					Name:       gofakeit.BeerName(),
					Topic:      gofakeit.BeerMalt(),
					AgeRange:   &dto.AgeRangeInput{},
					Visibility: enums.PrivateVisibility,
					Preset:     enums.PresetPrivateChat,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to create matrix room invalid url",
			args: args{
				ctx: context.Background(),
				room: &dto.CommunityInput{
					Name:       gofakeit.BeerName(),
					Topic:      gofakeit.BeerMalt(),
					AgeRange:   &dto.AgeRangeInput{},
					Visibility: enums.PrivateVisibility,
					Preset:     enums.PresetPrivateChat,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to create matrix room",
			args: args{
				ctx: context.Background(),
				room: &dto.CommunityInput{
					Name:       gofakeit.BeerName(),
					Topic:      gofakeit.BeerMalt(),
					AgeRange:   &dto.AgeRangeInput{},
					Visibility: enums.PrivateVisibility,
					Preset:     enums.PresetPrivateChat,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: unable to create matrix, an error occurred",
			args: args{
				ctx: context.Background(),
				room: &dto.CommunityInput{
					Name:       gofakeit.BeerName(),
					Topic:      gofakeit.BeerMalt(),
					AgeRange:   &dto.AgeRangeInput{},
					Visibility: enums.PrivateVisibility,
					Preset:     enums.PresetPrivateChat,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			matrixClient := matrix.ServiceImpl{
				BaseURL: "https://example.com",
			}

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			m := matrix.NewMatrixImpl(matrixClient.BaseURL)

			if tt.name == "happy case: Create room successfully" {
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/createRoom",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"name":       gofakeit.BeerName(),
							"topic":      gofakeit.BeerMalt(),
							"visibility": enums.PrivateVisibility,
							"preset":     enums.PresetPublicChat,
							"room_id":    gofakeit.BeerName(),
						})
						return resp, err
					},
				)
			}
			if tt.name == "sad case: unable to create matrix room invalid url" {
				httpmock.RegisterResponder(http.MethodPost, "/_matix/client/v7/createRoom",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"name":       gofakeit.BeerName(),
							"topic":      gofakeit.BeerMalt(),
							"visibility": enums.PrivateVisibility,
							"preset":     enums.PresetPublicChat,
							"room_id":    gofakeit.BeerName(),
						})
						return resp, err
					},
				)
			}

			if tt.name == "sad case: unable to create matrix room" {
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/createRoom",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(400, map[string]interface{}{
							"error": gofakeit.BeerName(),
						})
						return resp, err
					},
				)
			}

			if tt.name == "sad case: unable to create matrix, an error occurred" {
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/createRoom",
					func(req *http.Request) (*http.Response, error) {
						return nil, fmt.Errorf("unable to create matrix room")
					},
				)
			}

			got, err := m.CreateCommunity(tt.args.ctx, tt.args.room)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.CreateCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == "" {
				t.Errorf("ServiceImpl.CreateCommunity() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
