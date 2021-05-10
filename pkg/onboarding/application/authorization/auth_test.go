package authorization

import (
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestCheckPemissions(t *testing.T) {
	type args struct {
		subject string
		input   resources.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: permission is set and subject has permission",
			args: args{
				subject: "254711223344",
				input: resources.PermissionInput{
					Resource: "update_primary_phone",
					Action:   "edit",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid: unknown subject with unkown resource",
			args: args{
				subject: "mail@example.com",
				input: resources.PermissionInput{
					Resource: "unknown_resource",
					Action:   "edit",
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckPemissions(tt.args.subject, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPemissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckPemissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckAuthorization(t *testing.T) {
	type args struct {
		subject    string
		permission resources.PermissionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid: permission is set and subject has permission",
			args: args{
				subject: "254711223344",
				permission: resources.PermissionInput{
					Resource: "update_primary_phone",
					Action:   "edit",
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid: unknown subject with unkown resource",
			args: args{
				subject: "mail@example.com",
				permission: resources.PermissionInput{
					Resource: "unknown_resource",
					Action:   "edit",
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAuthorization(tt.args.subject, tt.args.permission)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAuthorization() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckAuthorization() = %v, want %v", got, tt.want)
			}
		})
	}
}
