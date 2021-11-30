package gorm

import (
	"testing"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/tj/assert"
)

func TestNewPGInstance(t *testing.T) {
	tests := []struct {
		name    string
		want    *PGInstance
		wantErr bool
	}{
		{
			name:    "Happy case",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPGInstance()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPGInstance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.NotNil(t, got)
		})
	}
}
