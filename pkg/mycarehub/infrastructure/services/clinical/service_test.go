package clinical_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/clinical"
)

func TestServiceClinical_DeleteFHIRPatientByPhone(t *testing.T) {
	type args struct {
		ctx         context.Context
		phoneNumber string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: success remove patient",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to make request",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid status code",
			args: args{
				ctx:         context.Background(),
				phoneNumber: gofakeit.Phone(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExt := extensionMock.NewFakeExtension()
			c := clinical.NewServiceClinical(fakeExt)

			if tt.name == "sad case: fail to make request" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					return nil, fmt.Errorf("failed to make request")
				}
			}

			if tt.name == "sad case: invalid status code" {
				fakeExt.MockMakeRequestFn = func(ctx context.Context, method, path string, body interface{}) (*http.Response, error) {
					msg := struct {
						Message string `json:"message"`
					}{
						Message: "success",
					}

					payload, _ := json.Marshal(msg)

					return &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(bytes.NewBuffer(payload))}, nil
				}
			}
			if err := c.DeleteFHIRPatientByPhone(tt.args.ctx, tt.args.phoneNumber); (err != nil) != tt.wantErr {
				t.Errorf("ServiceClinical.DeleteFHIRPatientByPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
