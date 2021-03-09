package engagement_test

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	extMock "gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
)

var fakeISCExt extMock.ISCClientExtension
var engClient extension.ISCClientExtension = &fakeISCExt

const (
	futureHours = 878400
)

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
	e := engagement.NewServiceEngagementImpl(engClient)

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
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: http.StatusOK,
						Body:       nil,
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_publish_kyc_feed_item" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: http.StatusBadRequest,
						Body:       nil,
					}, fmt.Errorf("fail to publish kyc feed item")
				}
			}

			resp, err := e.PublishKYCFeedItem(tt.args.uid, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEngagementImpl.PublishKYCFeedItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("ServiceEngagementImpl.PublishKYCFeedItem() = %v, want %v", resp, tt.want)
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

				if resp.StatusCode != tt.want.StatusCode {
					t.Errorf("expected status code 200 but got %v", resp.StatusCode)
					return
				}
			}
		})
	}
}

func TestServiceEngagementImpl_PublishKYCNudge(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient)

	payload := base.Nudge{
		ID:             strconv.Itoa(int(time.Now().Unix()) + 10), // add 10 to make it unique
		SequenceNumber: int(time.Now().Unix()) + 20,               // add 20 to make it unique
		Visibility:     "SHOW",
		Status:         "PENDING",
		Expiry:         time.Now().Add(time.Hour * futureHours),
		Title:          fmt.Sprintf("Complete your %v KYC", base.PartnerTypeRider),
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
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "OK",
						StatusCode: http.StatusOK,
						Body:       nil,
					}, nil
				}
			}

			if tt.name == "invalid:fail_to_publish_kyc_nudge" {
				fakeISCExt.MakeRequestFn = func(method string, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						Status:     "BAD REQUEST",
						StatusCode: http.StatusBadRequest,
						Body:       nil,
					}, fmt.Errorf("fail to publish kyc feed item")
				}
			}

			resp, err := e.PublishKYCNudge(tt.args.uid, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceEngagementImpl.PublishKYCNudge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(resp, tt.want) {
				t.Errorf("ServiceEngagementImpl.PublishKYCNudge() = %v, want %v", resp, tt.want)
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

				if resp.StatusCode != tt.want.StatusCode {
					t.Errorf("expected status code 200 but got %v", resp.StatusCode)
					return
				}
			}
		})
	}
}

func TestServiceEngagementImpl_SendMail(t *testing.T) {
	e := engagement.NewServiceEngagementImpl(engClient)

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
