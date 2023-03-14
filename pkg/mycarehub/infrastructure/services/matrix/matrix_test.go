package matrix_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/matrix"
)

func TestServiceImpl_RegisterUser(t *testing.T) {
	type args struct {
		ctx                 context.Context
		auth                *domain.MatrixAuth
		registrationPayload *domain.MatrixUserRegistration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: Successfully register user",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: "test",
					Password: "test",
				},
				registrationPayload: &domain.MatrixUserRegistration{
					Username: "test",
					Password: "test",
					Admin:    true,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to register user",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: "test",
					Password: "test",
				},
				registrationPayload: &domain.MatrixUserRegistration{
					Username: "test",
					Password: "test",
					Admin:    true,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid method",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: "test",
					Password: "test",
				},
				registrationPayload: &domain.MatrixUserRegistration{
					Username: "test",
					Password: "test",
					Admin:    true,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseURL := "https://example.com"

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			m := matrix.NewMatrixImpl(baseURL)

			if tt.name == "happy case: Successfully register user" {
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/login",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"auth": map[string]interface{}{
								"type": "m.login.dummy",
							},
							"username": gofakeit.BeerMalt(),
							"password": gofakeit.BeerName(),
						})
						return resp, err
					},
				)
				url := fmt.Sprintf("/_synapse/admin/v2/users/@%s:prohealth360.org", "test")
				httpmock.RegisterResponder(http.MethodPut, url,
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(201, map[string]interface{}{
							"password": "test",
						})
						return resp, err
					},
				)
			}
			if tt.name == "sad case: unable to register user" {
				url := fmt.Sprintf("/_synapse/admin/v2/users/@%s:prohealth360.org", "test")
				httpmock.RegisterResponder(http.MethodPut, url,
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(400, map[string]interface{}{
							"error": gofakeit.BeerName(),
						})
						return resp, err
					},
				)
			}
			if tt.name == "sad case: invalid method" {
				url := fmt.Sprintf("/_synapse/admin/v2/users/@%s:prohealth360.org", "test")
				httpmock.RegisterResponder(http.MethodGet, url,
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(405, map[string]interface{}{
							"error": gofakeit.BeerName(),
						})
						return resp, err
					},
				)
			}

			_, err := m.RegisterUser(tt.args.ctx, tt.args.auth, tt.args.registrationPayload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.RegisterUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_CreateCommunity(t *testing.T) {
	type args struct {
		ctx  context.Context
		auth *domain.MatrixAuth
		room *dto.CommunityInput
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy case: create community",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.Email(),
					Password: gofakeit.Email(),
				},
				room: &dto.CommunityInput{
					Name:  "test",
					Topic: "test",
					AgeRange: &dto.AgeRangeInput{
						LowerBound: 0,
						UpperBound: 0,
					},
					Gender:         []*enumutils.Gender{},
					Visibility:     "public",
					Preset:         "public_chat",
					ClientType:     []*enums.ClientType{},
					OrganisationID: gofakeit.UUID(),
					ProgramID:      gofakeit.UUID(),
					FacilityID:     gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to create matrix room invalid url",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
				},
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
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
				},
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
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
				},
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
			baseURL := "https://example.com"

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			m := matrix.NewMatrixImpl(baseURL)

			if tt.name == "happy case: create community" {
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/login",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"identifier": map[string]interface{}{
								"type": "m.id.user",
								"user": "test",
							},
							"type":     "m.login.password",
							"password": "test@matrix",
						})

						return resp, err
					},
				)
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

			_, err := m.CreateCommunity(tt.args.ctx, tt.args.auth, tt.args.room)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.CreateCommunity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
