package utils

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: read json file",
			args: args{
				path: "testData/test.json",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid path",
			args: args{
				path: "invalid.path",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
