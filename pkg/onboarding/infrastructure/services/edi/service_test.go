package edi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	extMock "github.com/savannahghi/onboarding/pkg/onboarding/application/extension/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi"
	ediMock "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/edi/mock"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	mockRepo "github.com/savannahghi/onboarding/pkg/onboarding/repository/mock"
	"github.com/savannahghi/profileutils"
	"gitlab.slade360emr.com/go/apiclient"
)

var fakeISCExt extMock.ISCClientExtension
var ediClient extension.ISCClientExtension = &fakeISCExt
var fakeRepo mockRepo.FakeOnboardingRepository
var r repository.OnboardingRepository = &fakeRepo
var fakeEDIsvs ediMock.FakeServiceEDI

func TestServiceEDIImpl_LinkCover(t *testing.T) {
	e := edi.NewEdiService(ediClient, r)

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
				fakeEDIsvs.GetSladerDataFn = func(ctx context.Context, phoneNumber string) (*[]apiclient.MarketingData, error) {
					return &[]apiclient.MarketingData{
						{
							MemberNumber:   uuid.New().String(),
							PayerSladeCode: "457",
						},
					}, nil
				}

				data := []apiclient.MarketingData{
					{
						FirstName:      gofakeit.Name(),
						LastName:       gofakeit.Name(),
						Email:          gofakeit.Email(),
						Phone:          gofakeit.PhoneFormatted(),
						PayerSladeCode: "32",
						MemberNumber:   "A100",
						Segment:        "One",
					},
				}

				b, _ := json.Marshal(data)

				fakeISCExt.MakeRequestFn = func(
					ctx context.Context,
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       ioutil.NopCloser(bytes.NewBuffer(b)),
					}, nil
				}

				fakeRepo.SaveCoverAutolinkingEventsFn = func(ctx context.Context, input *dto.CoverLinkingEvent) (*dto.CoverLinkingEvent, error) {
					return &dto.CoverLinkingEvent{ID: uuid.NewString()}, nil
				}
			}

			if tt.name == "Sad Case - Fail to link a cover" {
				fakeEDIsvs.GetSladerDataFn = func(ctx context.Context, phoneNumber string) (*[]apiclient.MarketingData, error) {
					return &[]apiclient.MarketingData{
						{
							MemberNumber:   uuid.New().String(),
							PayerSladeCode: "457",
						},
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

func TestServiceEDIImpl_GetSladerData(t *testing.T) {
	e := edi.NewEdiService(ediClient, r)
	ctx := context.Background()
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	payload := []apiclient.MarketingData{
		{
			FirstName:      "Test",
			LastName:       "User",
			Email:          "+254700000000@users.bewell.co.ke",
			Phone:          "+254700000000",
			Payor:          "Resolution Insurance",
			PayerSladeCode: "1234",
			MemberNumber:   "1234",
			Segment:        "test_segment",
		},
	}
	validRespPayload := `
	[
		{
			"firstname":"Test",
			"lastname":"User",
			"email":"+254700000000@users.bewell.co.ke",
			"phone":"+254700000000",
			"payor":"Resolution Insurance",
			"payer_slade_code":"1234",
			"member_number":"1234",
			"segment":"test_segment"
		}
	]
	`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	tests := []struct {
		name    string
		args    args
		want    *[]apiclient.MarketingData
		wantErr bool
	}{
		{
			name: "Happy Case -> Successfully Get a slader data",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
			},
			want:    &payload,
			wantErr: false,
		},
		{
			name: "Sad Case -> Fail to Get a slader data",
			args: args{
				ctx:         ctx,
				phoneNumber: interserviceclient.TestUserPhoneNumber,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case -> Successfully Get a slader data" {
				fakeISCExt.MakeRequestFn = func(
					ctx context.Context,
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "Sad Case -> Fail to Get a slader data" {
				fakeISCExt.MakeRequestFn = func(
					ctx context.Context,
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return nil, fmt.Errorf("failed to get slader data")
				}
			}

			got, err := e.GetSladerData(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEDIImpl.GetSladerData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceEDIImpl.GetSladerData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceEDIImpl_LinkEDIMemberCover(t *testing.T) {
	e := edi.NewEdiService(ediClient, r)
	ctx := context.Background()

	payload := dto.CoverInput{
		PayerSladeCode: 456,
		MemberNumber:   "123456",
		UID:            "Oq70MFGhY7fkoEXiQrRvqMm0BqB3",
		PushToken:      []string{"Oq70MFGhY7fkoEXiQrRvqMm0BqB3"},
	}

	marshalled, err := json.Marshal(payload)
	if err != nil {
		t.Errorf("failed to marshall payload: %w", err)
		return
	}

	bs := ioutil.NopCloser(bytes.NewReader([]byte(marshalled)))
	type args struct {
		ctx            context.Context
		phoneNumber    string
		membernumber   string
		payersladecode int
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
				ctx:            ctx,
				phoneNumber:    interserviceclient.TestUserPhoneNumber,
				membernumber:   "123456",
				payersladecode: 456,
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "Sad Case - Fail to link a cover",
			args: args{
				ctx:            ctx,
				phoneNumber:    interserviceclient.TestUserPhoneNumber,
				membernumber:   "123456",
				payersladecode: 456,
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case - Successfully link a cover" {
				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						PushTokens:   []string{"Oq70MFGhY7fkoEXiQrRvqMm0BqB3"},
						VerifiedUIDS: []string{"Oq70MFGhY7fkoEXiQrRvqMm0BqB3"},
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
						Body:       bs,
					}, nil
				}
			}

			if tt.name == "Sad Case - Fail to link a cover" {
				fakeRepo.GetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, suspended bool) (*profileutils.UserProfile, error) {
					return &profileutils.UserProfile{
						PushTokens:   []string{"Oq70MFGhY7fkoEXiQrRvqMm0BqB3"},
						VerifiedUIDS: []string{"Oq70MFGhY7fkoEXiQrRvqMm0BqB3"},
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
			got, err := e.LinkEDIMemberCover(tt.args.ctx, tt.args.phoneNumber, tt.args.membernumber, tt.args.payersladecode)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEDIImpl.LinkEDIMemberCover() error = %v, wantErr %v", err, tt.wantErr)
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
