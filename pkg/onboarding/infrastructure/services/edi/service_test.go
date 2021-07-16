package edi_test

// import (
// 	"context"
// 	"net/http"
// 	"testing"

// 	"github.com/google/uuid"
// 	"gitlab.slade360emr.com/go/base"
// 	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
// 	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
// 	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
// 	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/edi"
// 	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
// 	mockRepo "gitlab.slade360emr.com/go/profile/pkg/onboarding/repository/mock"
// )

// var fakeISCExt extMock.ISCClientExtension
// var ediClient extension.ISCClientExtension = &fakeISCExt
// var fakeRepo mockRepo.FakeOnboardingRepository
// var r repository.OnboardingRepository = &fakeRepo

// func TestServiceEDIImpl_LinkCover(t *testing.T) {
// 	e := edi.NewEdiService(ediClient, r)

// 	type args struct {
// 		ctx         context.Context
// 		phoneNumber string
// 		uid         string
// 		pushToken   []string
// 	}
// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantErr    bool
// 		wantStatus int
// 	}{
// 		{
// 			name: "Happy Case - Successfully link a cover",
// 			args: args{
// 				phoneNumber: interserviceclient.TestUserPhoneNumber,
// 				uid:         uuid.New().String(),
// 				pushToken:   []string{uuid.New().String()},
// 			},
// 			wantErr:    true, // TODO: Fix and make false
// 			wantStatus: http.StatusOK,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			if tt.name == "Happy Case - Successfully link a cover" {
// 				fakeRepo.GetUserMarketingDataFn = func(ctx context.Context, phoneNumber string) (*dto.Segment, error) {
// 					return &dto.Segment{
// 						MemberNumber: uuid.New().String(),
// 					}, nil
// 				}

// 				fakeISCExt.MakeRequestFn = func(
// 					ctx context.Context,
// 					method string,
// 					path string,
// 					body interface{},
// 				) (*http.Response, error) {
// 					return &http.Response{
// 						Status:     "OK",
// 						StatusCode: 200,
// 						Body:       nil,
// 					}, nil
// 				}
// 			}

// 			got, err := e.LinkCover(tt.args.ctx, tt.args.phoneNumber, tt.args.uid, tt.args.pushToken)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ServiceEDIImpl.LinkCover() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if got.StatusCode != http.StatusOK {
// 				t.Errorf("invalid status code returned %v", got.StatusCode)
// 				return
// 			}
// 		})
// 	}
// }
