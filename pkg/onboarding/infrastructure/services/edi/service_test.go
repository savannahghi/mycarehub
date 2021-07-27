package edi_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	extMock "github.com/savannahghi/onboarding/pkg/onboarding/application/extension/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
	engagementMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	mockRepo "github.com/savannahghi/onboarding/pkg/onboarding/repository/mock"
	"gitlab.slade360emr.com/go/apiclient"
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
				phoneNumber: interserviceclient.TestUserPhoneNumber,
				uid:         uuid.New().String(),
				pushToken:   []string{uuid.New().String()},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "Sad Case - Fail to link a cover",
			args: args{
				phoneNumber: interserviceclient.TestUserPhoneNumber,
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
				fakeRepo.SaveCoverAutolinkingEventsFn = func(ctx context.Context, input *dto.CoverLinkingEvent) (*dto.CoverLinkingEvent, error) {
					return &dto.CoverLinkingEvent{ID: uuid.NewString()}, nil
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
					}, fmt.Errorf("an error occurred!")
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
