package engagement_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
)

var fakeISCExt extMock.ISCClientExtension
var engClient extension.ISCClientExtension = &fakeISCExt

func TestServiceEngagementImpl_ResolveDefaultNudgeByTitle(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient)

	type args struct {
		UID        string
		flavour    base.Flavour
		nudgeTitle string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid:_resolve_a_default_nudge",
			args: args{
				UID:        uuid.New().String(),
				flavour:    base.FlavourConsumer,
				nudgeTitle: "Nudge Title",
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid:_nudge_not_found",
			args: args{
				UID:        uuid.New().String(),
				flavour:    base.FlavourConsumer,
				nudgeTitle: "Nudge Title",
			},
			wantErr:    true,
			wantStatus: http.StatusNotFound,
		},
		{
			name: "invalid:_bad_request_sent",
			args: args{
				UID:        uuid.New().String(),
				flavour:    base.FlavourConsumer,
				nudgeTitle: "Nudge Title",
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid:_error_occured_when_sending_the_request",
			args: args{
				UID:        uuid.New().String(),
				flavour:    base.FlavourConsumer,
				nudgeTitle: "Nudge Title",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_resolve_a_default_nudge" {
				fakeISCExt.MakeRequestFn = func(
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       nil,
					}, nil
				}
			}

			if tt.name == "invalid:_nudge_not_found" {
				fakeISCExt.MakeRequestFn = func(
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return &http.Response{
						Status:     "NOT FOUND",
						StatusCode: 404,
						Body:       nil,
					}, fmt.Errorf("nil nudge")
				}
			}

			if tt.name == "invalid:_bad_request_sent" {
				fakeISCExt.MakeRequestFn = func(
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: 400,
						Body:       nil,
					}, fmt.Errorf("error occured")
				}
			}

			if tt.name == "invalid:_error_occured_when_sending_the_request" {
				fakeISCExt.MakeRequestFn = func(
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return nil, fmt.Errorf("error occured")
				}
			}

			resp, err := e.ResolveDefaultNudgeByTitle(
				tt.args.UID,
				tt.args.flavour,
				tt.args.nudgeTitle,
			)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}

				if resp.StatusCode != tt.wantStatus {
					t.Errorf("expected status code 200 but got %v", resp.StatusCode)
					return
				}
			}
		})
	}
}
