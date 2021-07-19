package edi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/apiclient"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/edi"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	engagementMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
	mockRepo "gitlab.slade360emr.com/go/profile/pkg/onboarding/repository/mock"
)

var fakeISCExt extMock.ISCClientExtension
var ediClient extension.ISCClientExtension = &fakeISCExt
var fakeRepo mockRepo.FakeOnboardingRepository
var r repository.OnboardingRepository = &fakeRepo
var fakeEngagementSvs engagementMock.FakeServiceEngagement
var engagementSvc engagement.ServiceEngagement = &fakeEngagementSvs

func TestServiceEDIImpl_LinkCover(t *testing.T) {
	e := edi.NewEdiService(ediClient, r, engagementSvc)

	type args struct {
		ctx         context.Context
		phoneNumber string
		uid         string
		pushToken   []string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "Happy Case - Successfully link a cover",
			args: args{
				phoneNumber: base.TestUserPhoneNumber,
				uid:         uuid.New().String(),
				pushToken:   []string{uuid.New().String()},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "Sad Case - Fail to link a cover",
			args: args{
				phoneNumber: base.TestUserPhoneNumber,
				uid:         uuid.New().String(),
				pushToken:   []string{uuid.New().String()},
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Happy Case - Successfully link a cover" {
				fakeEngagementSvs.GetSladerDataFn = func(ctx context.Context, phoneNumber string) (*apiclient.Segment, error) {
					return &apiclient.Segment{
						MemberNumber:   uuid.New().String(),
						PayerSladeCode: "457",
					}, nil
				}

				fakeISCExt.MakeRequestFn = func(
					ctx context.Context,
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

			if tt.name == "Sad Case - Fail to link a cover" {
				fakeEngagementSvs.GetSladerDataFn = func(ctx context.Context, phoneNumber string) (*apiclient.Segment, error) {
					return &apiclient.Segment{
						MemberNumber:   uuid.New().String(),
						PayerSladeCode: "457",
					}, nil
				}

				fakeISCExt.MakeRequestFn = func(
					ctx context.Context,
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: http.StatusBadRequest,
						Body:       nil,
					}, fmt.Errorf("an error occured!")
				}
			}

			got, err := e.LinkCover(tt.args.ctx, tt.args.phoneNumber, tt.args.uid, tt.args.pushToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEDIImpl.LinkCover() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.StatusCode != http.StatusOK {
					t.Errorf("invalid status code returned %v", got.StatusCode)
					return
				}
			}
		})
	}
}
