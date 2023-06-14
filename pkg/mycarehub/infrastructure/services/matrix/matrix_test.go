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

func TestServiceImpl_CheckIfUserIsAdmin(t *testing.T) {
	type args struct {
		ctx    context.Context
		auth   *domain.MatrixAuth
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: check if a user is admin",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
				},
				userID: "test",
			},
			wantErr: false,
		},
		{
			name: "sad case: unable to check if a user is admin",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
				},
				userID: "test",
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

			if tt.name == "happy case: check if a user is admin" {
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

				httpmock.RegisterResponder(http.MethodGet, "/_synapse/admin/v1/users/@test:prohealth360.org/admin",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, nil)

						return resp, err
					},
				)
			}
			if tt.name == "sad case: unable to check if a user is admin" {
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

				httpmock.RegisterResponder(http.MethodGet, "/_synapse/admin/v1/users/@test:prohealth360.org/admin",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(500, nil)

						return resp, err
					},
				)
			}

			_, err := m.CheckIfUserIsAdmin(tt.args.ctx, tt.args.auth, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.CheckIfUserIsAdmin() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_SearchUsers(t *testing.T) {
	type args struct {
		ctx        context.Context
		limit      int
		searchTerm string
		auth       *domain.MatrixAuth
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: successfully search users",
			args: args{
				ctx:        context.Background(),
				limit:      10,
				searchTerm: "test",
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to search users",
			args: args{
				ctx:        context.Background(),
				limit:      10,
				searchTerm: "test",
				auth: &domain.MatrixAuth{
					Username: gofakeit.Name(),
					Password: gofakeit.BeerName(),
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

			if tt.name == "Happy case: successfully search users" {
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
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/user_directory/search",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"limit":       10,
							"search_term": gofakeit.BeerMalt(),
						})

						return resp, err
					},
				)
			}
			if tt.name == "Sad case: unable to search users" {
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
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/user_directory/search",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(400, map[string]interface{}{
							"limit":       10,
							"search_term": gofakeit.BeerMalt(),
						})

						return resp, err
					},
				)
			}

			_, err := m.SearchUsers(tt.args.ctx, tt.args.limit, tt.args.searchTerm, tt.args.auth)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.SearchUsers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestServiceImpl_DeactivateUser(t *testing.T) {
	type args struct {
		ctx    context.Context
		userID string
		auth   *domain.MatrixAuth
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: deactivate user",
			args: args{
				ctx:    context.Background(),
				userID: "@test:prohealth360.org",
				auth: &domain.MatrixAuth{
					Username: gofakeit.BeerName(),
					Password: gofakeit.UUID(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to deactivate user",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.BeerName(),
					Password: gofakeit.UUID(),
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

			if tt.name == "Happy case: deactivate user" {
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
				httpmock.RegisterResponder(http.MethodPost, "/_synapse/admin/v1/deactivate/@test:prohealth360.org",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"erase": true,
						})

						return resp, err
					},
				)
			}
			if tt.name == "Sad case: unable to deactivate user" {
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
				httpmock.RegisterResponder(http.MethodPost, "/_synapse/admin/v1/deactivate/%s",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(400, map[string]interface{}{
							"erase": true,
						})

						return resp, err
					},
				)
			}

			if err := m.DeactivateUser(tt.args.ctx, tt.args.userID, tt.args.auth); (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.DeactivateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceImpl_SetPusher(t *testing.T) {
	kind := "http"
	type args struct {
		ctx     context.Context
		auth    *domain.MatrixAuth
		payload *domain.PusherPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Set pusher",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.BeerName(),
					Password: gofakeit.UUID(),
				},
				payload: &domain.PusherPayload{
					AppDisplayName: "MCH",
					AppID:          "com.example.app.ios",
					Append:         false,
					PusherData: domain.PusherData{
						Format: "event_id_only",
						URL:    "https://push-gateway.location.here/_matrix/push/v1/notify",
					},
					DeviceDisplayName: "Samsung Galaxy",
					Kind:              &kind,
					Lang:              "en-US",
					Pushkey:           gofakeit.HipsterSentence(50),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to set pusher",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.BeerName(),
					Password: gofakeit.UUID(),
				},
				payload: &domain.PusherPayload{
					AppDisplayName: "MCH",
					AppID:          "com.example.app.ios",
					Append:         false,
					PusherData: domain.PusherData{
						Format: "event_id_only",
						URL:    "https://push-gateway.location.here/_matrix/push/v1/notify",
					},
					DeviceDisplayName: "Samsung Galaxy",
					Kind:              &kind,
					Lang:              "en-US",
					Pushkey:           gofakeit.HipsterSentence(50),
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

			if tt.name == "Happy case: Set pusher" {
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
				httpmock.RegisterResponder(http.MethodPost, "/_matrix/client/v3/pushers/set",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{
							"app_display_name": "Mat Rix",
							"app_id":           "com.example.app.ios",
							"append":           false,
							"data": map[string]interface{}{
								"format": "event_id_only",
								"url":    "https://push-gateway.location.here/_matrix/push/v1/notify",
							},
							"device_display_name": "iPhone 9",
							"kind":                "http",
							"lang":                "en",
							"profile_tag":         "xxyyzz",
							"pushkey":             "APA91bHPRgkF3JUikC4ENAHEeMrd41Zxv3hVZjC9KtT8OvPVGJ-hQMRKRrZuJAEcl7B338qju59zJMjw2DELjzEvxwYv7hH5Ynpc1ODQ0aT4U4OFEeco8ohsN5PjL1iC2dNtk2BAokeMCg2ZXKqpc8FXKmhX94kIxQ",
						})

						return resp, err
					},
				)
			}
			if tt.name == "Sad case: unable to set pusher" {
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
				httpmock.RegisterResponder(http.MethodPost, "/_synapse/admin/v1/deactivate/%s",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(400, map[string]interface{}{
							"app_display_name": "Mat Rix",
							"app_id":           "com.example.app.ios",
						})

						return resp, err
					},
				)
			}

			if err := m.SetPusher(tt.args.ctx, tt.args.auth, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.SetPusher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServiceImpl_SetPushRule(t *testing.T) {
	type args struct {
		ctx             context.Context
		auth            *domain.MatrixAuth
		queryPathValues *domain.QueryPathValues
		payload         *domain.PushRulePayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: should return an empty response when the request is successful",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.BeerName(),
					Password: gofakeit.UUID(),
				},
				queryPathValues: &domain.QueryPathValues{
					Scope:  "global",
					RuleID: "m.room.message",
					Kind:   "room",
				},
				payload: &domain.PushRulePayload{
					Conditions: []domain.Conditions{
						{
							Kind:    "event_match",
							Key:     "type",
							Pattern: "m.room.message",
						},
					},
					Actions: []any{
						"notify",
						map[string]interface{}{
							"set_tweak": "highlight",
						},
						map[string]interface{}{
							"set_tweak": "sound",
							"value":     "default",
						},
					},
					Kind: "room",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: fail to set push rule",
			args: args{
				ctx: context.Background(),
				auth: &domain.MatrixAuth{
					Username: gofakeit.BeerName(),
					Password: gofakeit.UUID(),
				},
				queryPathValues: &domain.QueryPathValues{
					Scope:  "global",
					RuleID: "m.room.message",
					Kind:   "room",
				},
				payload: &domain.PushRulePayload{
					Conditions: []domain.Conditions{
						{
							Kind:    "event_match",
							Key:     "type",
							Pattern: "m.room.message",
						},
					},
					Actions: []any{
						"notify",
						map[string]interface{}{
							"set_tweak": "highlight",
						},
						map[string]interface{}{
							"set_tweak": "sound",
							"value":     "default",
						},
					},
					Kind: "room",
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

			if tt.name == "Happy case: should return an empty response when the request is successful" {
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

				httpmock.RegisterResponder(http.MethodPut, "/_matrix/client/v3/pushrules/global/room/m.room.message",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(200, map[string]interface{}{})

						return resp, err
					},
				)
			}
			if tt.name == "Sad case: fail to set push rule" {
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

				httpmock.RegisterResponder(http.MethodPut, "/_matrix/client/v3/pushrules/global/room/m.room.message",
					func(req *http.Request) (*http.Response, error) {
						resp, err := httpmock.NewJsonResponse(400, map[string]interface{}{})

						return resp, err
					},
				)
			}

			if err := m.SetPushRule(tt.args.ctx, tt.args.auth, tt.args.queryPathValues, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServiceImpl.SetPushRule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
