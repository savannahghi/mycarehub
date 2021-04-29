package engagement_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub/mock"
)

var fakeISCExt extMock.ISCClientExtension
var engClient extension.ISCClientExtension = &fakeISCExt
var fakeBaseExt extMock.FakeBaseExtensionImpl
var baseExt extension.BaseExtension = &fakeBaseExt
var fakePubSub mock.FakeServicePubSub
var ps pubsubmessaging.ServicePubSub = &fakePubSub

const (
	futureHours = 878400
)

func TestServiceEngagementImpl_ResolveDefaultNudgeByTitle(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

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
			name: "invalid:_error_occurred_when_sending_the_request",
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
					}, fmt.Errorf("error occurred")
				}
			}

			if tt.name == "invalid:_error_occurred_when_sending_the_request" {
				fakeISCExt.MakeRequestFn = func(
					method string,
					path string,
					body interface{},
				) (*http.Response, error) {
					return nil, fmt.Errorf("error occurred")
				}
			}

			err := e.ResolveDefaultNudgeByTitle(
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
			}
		})
	}
}

func TestServiceEngagementImpl_PublishKYCFeedItem(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	payload := base.Item{
		ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
		SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Persistent:     true,
		Status:         base.StatusPending,
		Visibility:     base.VisibilityShow,
		Author:         "Be.Well Team",
		Label:          "KYC",
		Tagline:        "Process incoming KYC",
		Text:           "Review KYC for the partner and either approve or reject",
		TextType:       base.TextTypeMarkdown,
		Icon: base.Link{
			ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
			URL:         base.LogoURL,
			LinkType:    base.LinkTypePngImage,
			Title:       "KYC Review",
			Description: "Review KYC for the partner and either approve or reject",
			Thumbnail:   base.LogoURL,
		},
		Timestamp: time.Now(),
		Actions: []base.Action{
			{
				ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
				SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
				Name:           "Review KYC details",
				Icon: base.Link{
					ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
					URL:         base.LogoURL,
					LinkType:    base.LinkTypePngImage,
					Title:       "Review KYC details",
					Description: "Review and approve or reject KYC details for the supplier",
					Thumbnail:   base.LogoURL,
				},
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
			},
		},
		Links: []base.Link{
			{
				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC process request",
				Description: "Process KYC request",
				Thumbnail:   base.LogoURL,
			},
		},
	}
	type args struct {
		uid     string
		payload base.Item
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "valid:publish_kyc_feed_item",
			args: args{
				uid:     uuid.New().String(),
				payload: payload,
			},
			want: &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_publish_kyc_feed_item",
			args: args{
				uid:     uuid.New().String(),
				payload: payload,
			},
			want: &http.Response{
				Status:     "BAD REQUEST",
				StatusCode: http.StatusBadRequest,
				Body:       nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:publish_kyc_feed_item" {
				fakeBaseExt.GetPubSubTopicFn = func(m *base.PubSubPayload) (string, error) {
					return "some.topic", nil
				}

				fakePubSub.PublishToPubsubFn = func(
					ctx context.Context,
					topicID string,
					payload []byte,
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_publish_kyc_feed_item" {
				fakeBaseExt.GetPubSubTopicFn = func(m *base.PubSubPayload) (string, error) {
					return "some.topic", nil
				}

				fakePubSub.PublishToPubsubFn = func(
					ctx context.Context,
					topicID string,
					payload []byte,
				) error {
					return fmt.Errorf("failed to publish feed item")
				}
			}

			err := e.PublishKYCFeedItem(tt.args.uid, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEngagementImpl.PublishKYCFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

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
			}
		})
	}
}

func TestServiceEngagementImpl_PublishKYCNudge(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	payload := base.Nudge{
		ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
		SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
		Visibility:     "SHOW",
		Status:         "PENDING",
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Title:          "Complete your KYC",
		Text:           "Fill in your Be.Well business KYC in order to start transacting",
		Links: []base.Link{
			{
				ID:          strconv.Itoa(int(time.Now().Unix()) + 30), // add 30 to make it unique,
				URL:         base.LogoURL,
				LinkType:    base.LinkTypePngImage,
				Title:       "KYC",
				Description: fmt.Sprintf("KYC for %v", base.PartnerTypeRider),
				Thumbnail:   base.LogoURL,
			},
		},
		Actions: []base.Action{
			{
				ID:             strconv.Itoa(int(time.Now().Unix()) + 40), // add 40 to make it unique
				SequenceNumber: int(time.Now().Unix()) + 50,               // add 50 to make it unique
				Name:           strings.ToUpper(fmt.Sprintf("COMPLETE_%v_%v_KYC", base.AccountTypeIndividual, base.PartnerTypeRider)),
				ActionType:     base.ActionTypePrimary,
				Handling:       base.HandlingFullPage,
				AllowAnonymous: false,
				Icon: base.Link{
					ID:          strconv.Itoa(int(time.Now().Unix()) + 60), // add 60 to make it unique
					URL:         base.LogoURL,
					LinkType:    base.LinkTypePngImage,
					Title:       fmt.Sprintf("Complete your %v KYC", base.PartnerTypeRider),
					Description: "Fill in your Be.Well business KYC in order to start transacting",
					Thumbnail:   base.LogoURL,
				},
			},
		},
		Users:                []string{uuid.New().String()},
		Groups:               []string{uuid.New().String()},
		NotificationChannels: []base.Channel{base.ChannelEmail, base.ChannelFcm},
	}

	type args struct {
		uid     string
		payload base.Nudge
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "valid:successfully_publish_kyc_nudge",
			args: args{
				uid:     uuid.New().String(),
				payload: payload,
			},
			want: &http.Response{
				Status:     "OK",
				StatusCode: http.StatusOK,
				Body:       nil,
			},
			wantErr: false,
		},
		{
			name: "invalid:fail_to_publish_kyc_nudge",
			args: args{
				uid:     uuid.New().String(),
				payload: payload,
			},
			want: &http.Response{
				Status:     "BAD REQUEST",
				StatusCode: http.StatusBadRequest,
				Body:       nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_publish_kyc_nudge" {
				fakeBaseExt.GetPubSubTopicFn = func(m *base.PubSubPayload) (string, error) {
					return "nudges.publish", nil
				}

				fakePubSub.PublishToPubsubFn = func(
					ctx context.Context,
					topicID string,
					payload []byte,
				) error {
					return nil
				}
			}

			if tt.name == "invalid:fail_to_publish_kyc_nudge" {
				fakeBaseExt.GetPubSubTopicFn = func(m *base.PubSubPayload) (string, error) {
					return "nudges.publish", nil
				}

				fakePubSub.PublishToPubsubFn = func(
					ctx context.Context,
					topicID string,
					payload []byte,
				) error {
					return fmt.Errorf("failed to publish kyc nudge")
				}
			}

			err := e.PublishKYCNudge(tt.args.uid, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEngagementImpl.PublishKYCNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

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

			}
		})
	}
}

func TestServiceEngagementImpl_SendMail(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	type args struct {
		email   string
		message string
		subject string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid:successfully_send_email",
			args: args{
				email:   "johndoe@gmail.com",
				message: "This is an update of how things are",
				subject: "update",
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid:use_an_invalid_email",
			args: args{
				email:   "12345",
				message: "This is an update of how things are",
				subject: "update",
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid:error_while_sending_request",
			args: args{
				email:   "johndoe",
				message: "This is an update of how things are",
				subject: "update",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:successfully_send_email" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       nil,
					}, nil
				}
			}

			if tt.name == "invalid:use_an_invalid_email" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: 400,
						Body:       nil,
					}, fmt.Errorf("an error occured! Invalid email address")
				}
			}

			if tt.name == "invalid:error_while_sending_request" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("an error occured!")
				}
			}
			err := e.SendMail(tt.args.email, tt.args.message, tt.args.subject)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEngagementImpl.SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}

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
			}
		})
	}
}

func TestServiceOTPImpl_VerifyOTP(t *testing.T) {
	ctx := context.Background()
	p := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	validRespPayload := `{"IsVerified":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `{""}`
	respReader1 := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))

	type args struct {
		ctx   context.Context
		phone string
		otp   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:verify_phone_with_valid_phone",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: false,
		},
		{
			name: "invalid:normalize_phone_fails",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
		{
			name: "invalid:make_http_request_fails",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
		{
			name: "invalid:make_http_request_returns_unexpected_status_code",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
		{
			name: "invalid:unmarshalling_of_respose_fails",
			args: args{
				ctx:   ctx,
				phone: "0721526895",
				otp:   "225025",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:verify_phone_with_valid_phone" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:normalize_phone_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize phone")
				}
			}

			if tt.name == "invalid:make_http_request_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make http request")
				}
			}

			if tt.name == "invalid:make_http_request_returns_unexpected_status_code" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 400,
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:unmarshalling_of_respose_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}

				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       respReader1,
					}, nil
				}
			}

			resp, err := p.VerifyOTP(tt.args.ctx, tt.args.phone, tt.args.otp)

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

				if resp != true {
					t.Errorf("response should be true got %v", resp)
					return
				}
			}

		})
	}
}

func TestServiceOTPImpl_GenerateAndSendOTP(t *testing.T) {
	ctx := context.Background()
	p := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	validRespPayload := `"234234"`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `"otp":"234234"`
	invalidRespReader := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))

	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *base.OtpResponse
		wantErr bool
	}{
		{
			name: "valid:_successfully_generate_and_send_otp",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want: &base.OtpResponse{
				OTP: "234234",
			},
			wantErr: false,
		},
		{
			name: "invalid:_make_request_fails",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_invalid_HTTP_response",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_unmarshall",
			args: args{
				ctx:   ctx,
				phone: "+2547345678",
			},
			want: &base.OtpResponse{
				OTP: "234234",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_successfully_generate_and_send_otp" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_make_request_fails" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make a request")
				}
			}

			if tt.name == "invalid:_invalid_HTTP_response" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
						Status:     "",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_unmarshall" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       invalidRespReader,
					}, nil
				}
			}

			resp, err := p.GenerateAndSendOTP(tt.args.ctx, tt.args.phone)

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

				if resp.OTP != tt.want.OTP {
					t.Errorf("expected %v, got %v", tt.want.OTP, resp.OTP)
					return
				}
			}

		})
	}
}

func TestServiceOTPImpl_SendRetryOTP(t *testing.T) {
	ctx := context.Background()
	p := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	validRespPayload := `"123123"`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `"otp":"123123"`
	invalidRespReader := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))

	type args struct {
		ctx       context.Context
		msisdn    string
		retryStep int
	}
	tests := []struct {
		name    string
		args    args
		want    *base.OtpResponse
		wantErr bool
	}{
		{
			name: "valid:_successfully_send_retry_otp",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want: &base.OtpResponse{
				OTP: "123123",
			},
			wantErr: false,
		},
		{
			name: "invalid:_make_request_fails",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_invalid_HTTP_response",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_unmarshall",
			args: args{
				ctx:       ctx,
				msisdn:    "+2547345678",
				retryStep: 1,
			},
			want: &base.OtpResponse{
				OTP: "234234",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_normalize_msisdn",
			args: args{
				ctx:       ctx,
				msisdn:    "+asc719ASD678",
				retryStep: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_successfully_send_retry_otp" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_make_request_fails" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make a request")
				}
			}

			if tt.name == "invalid:_invalid_HTTP_response" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}

				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
						Status:     "",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_unmarshall" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+2547345678"
					return &phone, nil
				}

				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       invalidRespReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_normalize_msisdn" {
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					return nil, fmt.Errorf("unable to normalize msisdn")
				}
			}

			resp, err := p.SendRetryOTP(tt.args.ctx, tt.args.msisdn, tt.args.retryStep)

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

				if resp.OTP != tt.want.OTP {
					t.Errorf("expected %v, got %v", tt.want.OTP, resp.OTP)
					return
				}
			}
		})
	}
}

func TestServiceOTPImpl_VerifyEmailOTP(t *testing.T) {
	ctx := context.Background()
	p := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)

	validRespPayload := `{"IsVerified":true}`
	respReader := ioutil.NopCloser(bytes.NewReader([]byte(validRespPayload)))

	inValidRespPayload := `{""}`
	invalidRespReader := ioutil.NopCloser(bytes.NewReader([]byte(inValidRespPayload)))
	type args struct {
		ctx   context.Context
		email string
		otp   string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid:_successfully_verify_email_otp",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "invalid:_make_request_fails",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_invalid_HTTP_response",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_unmarshall",
			args: args{
				ctx:   ctx,
				email: "johndoe@gmail.com",
				otp:   "345345",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_successfully_verify_email_otp" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_make_request_fails" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("unable to make a request")
				}
			}

			if tt.name == "invalid:_invalid_HTTP_response" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusUnprocessableEntity,
						Status:     "",
						Body:       respReader,
					}, nil
				}
			}

			if tt.name == "invalid:_unable_to_unmarshall" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       invalidRespReader,
					}, nil
				}
			}

			resp, err := p.VerifyEmailOTP(tt.args.ctx, tt.args.email, tt.args.otp)

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

				if resp != tt.want {
					t.Errorf("expected %v, got %v", tt.want, resp)
					return
				}
			}
		})
	}
}

func TestServiceEngagementImpl_NotifySupplierOnSuspension(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient, baseExt, ps)
	type args struct {
		input resources.EmailNotificationPayload
	}
	suspensionReason := `
	"This email is to inform you that as a result of your actions on April 12th, 2021, you have been issued a suspension for 1 week (7 days)"
	`
	supplierName := "Akaku Danger"
	subjectTitle := "Suspension from Be.Well"
	emailBody := suspensionReason
	emailAddress := base.TestUserEmail
	primaryPhone := base.TestUserPhoneNumber
	validInput := resources.EmailNotificationPayload{
		SupplierName: supplierName,
		SubjectTitle: subjectTitle,
		EmailBody:    emailBody,
		EmailAddress: emailAddress,
		PrimaryPhone: primaryPhone,
	}
	invalidEmailAddress := "12345"
	invalidInput := resources.EmailNotificationPayload{
		SupplierName: supplierName,
		SubjectTitle: subjectTitle,
		EmailBody:    emailBody,
		EmailAddress: invalidEmailAddress,
		PrimaryPhone: primaryPhone,
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "valid:successfully_send_email",
			args: args{
				input: validInput,
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid:wrong_email_address",
			args: args{
				input: invalidInput,
			},
			wantErr:    true,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid:error_while_sending_request",
			args: args{
				input: invalidInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:successfully_send_email" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: 200,
						Body:       nil,
					}, nil
				}
			}
			if tt.name == "invalid:wrong_email_address" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: 400,
						Body:       nil,
					}, fmt.Errorf("an error occured! Invalid email address")
				}
			}
			if tt.name == "invalid:error_while_sending_request" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("an error occured!")
				}
			}
			err := e.NotifySupplierOnSuspension(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEngagementImpl.SendMail() error = %v, wantErr %v", err, tt.wantErr)
			}

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
			}
		})
	}
}
