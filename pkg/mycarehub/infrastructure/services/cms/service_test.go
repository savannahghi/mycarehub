package cms_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	cmsService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/cms"
)

func TestServiceCMSImpl_RegisterClient(t *testing.T) {
	fakeExtension := extensionMock.NewFakeExtension()
	cmsService := cmsService.NewServiceCMS(fakeExtension, fakeExtension)

	type args struct {
		ctx    context.Context
		client *dto.PubSubCMSClientInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully register client",
			args: args{
				ctx: context.Background(),
				client: &dto.PubSubCMSClientInput{
					UserID:         uuid.New().String(),
					Name:           "test",
					Gender:         enumutils.GenderFemale,
					UserType:       "STAFF",
					PhoneNumber:    "0812345678",
					Handle:         "@test",
					Flavour:        feedlib.FlavourConsumer,
					DateOfBirth:    time.Time{},
					ClientID:       uuid.New().String(),
					ClientTypes:    []enums.ClientType{},
					EnrollmentDate: time.Time{},
					FacilityID:     uuid.New().String(),
					FacilityName:   "test",
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to make http request to register client",
			args: args{
				ctx: context.Background(),
				client: &dto.PubSubCMSClientInput{
					UserID:         uuid.New().String(),
					Name:           "test",
					Gender:         enumutils.GenderFemale,
					UserType:       "STAFF",
					PhoneNumber:    "0812345678",
					Handle:         "@test",
					Flavour:        feedlib.FlavourConsumer,
					DateOfBirth:    time.Time{},
					ClientID:       uuid.New().String(),
					ClientTypes:    []enums.ClientType{},
					EnrollmentDate: time.Time{},
					FacilityID:     uuid.New().String(),
					FacilityName:   "test",
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Bad status code returned",
			args: args{
				ctx: context.Background(),
				client: &dto.PubSubCMSClientInput{
					UserID:         uuid.New().String(),
					Name:           "test",
					Gender:         enumutils.GenderFemale,
					UserType:       "STAFF",
					PhoneNumber:    "0812345678",
					Handle:         "@test",
					Flavour:        feedlib.FlavourConsumer,
					DateOfBirth:    time.Time{},
					ClientID:       uuid.New().String(),
					ClientTypes:    []enums.ClientType{},
					EnrollmentDate: time.Time{},
					FacilityID:     uuid.New().String(),
					FacilityName:   "test",
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Happy Case - Successfully register client" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					input := dto.PubSubCMSClientInput{
						PhoneNumber: interserviceclient.TestUserPhoneNumber,
					}

					payload, err := json.Marshal(input)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}
			if tt.name == "Sad Case - unable to make http request to register client" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {

					return nil, fmt.Errorf("unable to make http request")
				}
			}
			if tt.name == "Sad Case - Bad status code returned" {
				fakeExtension.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					input := dto.PubSubCMSClientInput{
						PhoneNumber: interserviceclient.TestUserPhoneNumber,
					}

					payload, err := json.Marshal(input)
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusBadRequest,
						Status:     "OK",
						Body:       io.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}
			if err := cmsService.RegisterClient(tt.args.ctx, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("ServiceCMSImpl.RegisterClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
