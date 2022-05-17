package clinical

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	extMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
)

func TestServiceClinicalImpl_DeleteFHIRPatientByPhone(t *testing.T) {
	ctx := context.Background()

	fakeISC := extMock.NewFakeISCClientExtension()
	fakeExt := extMock.NewFakeExtension()
	c := NewServiceClinical(fakeISC, fakeExt)

	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				ctx:         ctx,
				phoneNumber: "0712345678",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				ctx:         ctx,
				phoneNumber: "44",
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					type PhoneNumberPayload struct {
						PhoneNumber string `json:"phoneNumber"`
					}

					payload, err := json.Marshal(&PhoneNumberPayload{PhoneNumber: tt.args.phoneNumber})
					if err != nil {
						t.Errorf("unable to marshal test item: %s", err)
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     "OK",
						Body:       ioutil.NopCloser(bytes.NewBuffer(payload)),
					}, nil
				}
			}
			if tt.name == "failure" {
				fakeISC.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusBadRequest,
					}, fmt.Errorf("failed to delete fhir patient")
				}
			}
			got, err := c.DeleteFHIRPatientByPhone(tt.args.ctx, tt.args.phoneNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceClinicalImpl.DeleteFHIRPatientByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceClinicalImpl.DeleteFHIRPatientByPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}
