package utils_test

import (
	"fmt"
	"testing"

	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

var baseExt mock.FakeBaseExtensionImpl

func TestISCClient(t *testing.T) {

	var baseExt = &baseExt

	tests := []struct {
		name    string
		want    *interserviceclient.InterServiceClient
		wantErr bool
	}{
		{
			name:    "should_load_yaml_and_start_isc_client",
			want:    &interserviceclient.InterServiceClient{},
			wantErr: false,
		},
		{
			name:    "should_fail_to_load_yaml",
			want:    &interserviceclient.InterServiceClient{},
			wantErr: true,
		},
		{
			name:    "should_fail_to_start_isc_service",
			want:    &interserviceclient.InterServiceClient{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "should_load_yaml_and_start_isc_client" {
				baseExt.LoadDepsFromYAMLFn = func() (*interserviceclient.DepsConfig, error) {
					return &interserviceclient.DepsConfig{}, nil
				}
				baseExt.SetupISCclientFn = func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
					return &interserviceclient.InterServiceClient{}, nil
				}
			}

			if tt.name == "should_fail_to_load_yaml" {
				baseExt.LoadDepsFromYAMLFn = func() (*interserviceclient.DepsConfig, error) {
					return nil, fmt.Errorf("error")
				}
			}

			if tt.name == "should_fail_to_start_isc_service" {
				baseExt.LoadDepsFromYAMLFn = func() (*interserviceclient.DepsConfig, error) {
					return &interserviceclient.DepsConfig{}, nil
				}
				baseExt.SetupISCclientFn = func(config interserviceclient.DepsConfig, serviceName string) (*interserviceclient.InterServiceClient, error) {
					return nil, fmt.Errorf("error")
				}
			}

			if !tt.wantErr {
				resp := utils.NewInterServiceClient("servicename", baseExt)
				assert.NotNil(t, resp)
			}
			if tt.wantErr {
				assert.Panics(t, func() {
					_ = utils.NewInterServiceClient("servicename", baseExt)
				})
			}
		})
	}

}
