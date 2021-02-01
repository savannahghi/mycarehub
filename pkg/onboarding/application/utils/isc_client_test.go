package utils_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension/mock"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

var baseExt mock.FakeBaseExtensionImpl

func TestISCClient(t *testing.T) {

	var baseExt = &baseExt

	tests := []struct {
		name    string
		want    *base.InterServiceClient
		wantErr bool
	}{
		{
			name:    "should_load_yaml_and_start_isc_client",
			want:    &base.InterServiceClient{},
			wantErr: false,
		},
		{
			name:    "should_fail_to_load_yaml",
			want:    &base.InterServiceClient{},
			wantErr: true,
		},
		{
			name:    "should_fail_to_start_isc_service",
			want:    &base.InterServiceClient{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "should_load_yaml_and_start_isc_client" {
				baseExt.LoadDepsFromYAMLFn = func() (*base.DepsConfig, error) {
					return &base.DepsConfig{}, nil
				}
				baseExt.SetupISCclientFn = func(config base.DepsConfig, serviceName string) (*base.InterServiceClient, error) {
					return &base.InterServiceClient{}, nil
				}
			}

			if tt.name == "should_fail_to_load_yaml" {
				baseExt.LoadDepsFromYAMLFn = func() (*base.DepsConfig, error) {
					return nil, fmt.Errorf("error")
				}
			}

			if tt.name == "should_fail_to_start_isc_service" {
				baseExt.LoadDepsFromYAMLFn = func() (*base.DepsConfig, error) {
					return &base.DepsConfig{}, nil
				}
				baseExt.SetupISCclientFn = func(config base.DepsConfig, serviceName string) (*base.InterServiceClient, error) {
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
