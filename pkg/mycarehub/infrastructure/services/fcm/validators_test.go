package fcm_test

import (
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm"
)

func TestValidateFCMData(t *testing.T) {
	type args struct {
		data map[string]string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Nil data",
			args: args{
				data: nil,
			},
			wantErr: false,
		},
		{
			name: "Happy Case - Good data",
			args: args{
				data: map[string]string{
					"a": "1",
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Using reserved words",
			args: args{
				data: map[string]string{
					"from": "should not be used",
				},
			},
			wantErr: true,
		},
		{
			name: "Sad Case - Using illegal prefix",
			args: args{
				data: map[string]string{
					"gcmData": "gcm is an illegal prefix",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := fcm.ValidateFCMData(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ValidateFCMData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
