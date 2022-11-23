package exceptions_test

import (
	"testing"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/stretchr/testify/assert"
)

func TestModelHasCustomError(t *testing.T) {
	exceptions := exceptions.CustomError{}
	cr := exceptions.Error()
	assert.NotNil(t, cr)
}

func TestGetErrorCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "TestGetErrorCode",
			args: args{
				err: nil,
			},
			want: int(exceptions.Internal),
		},
		{
			name: "TestGetErrorCode",
			args: args{
				err: &exceptions.CustomError{
					Err:     nil,
					Message: "",
					Code:    int(exceptions.Internal),
				},
			},
			want: int(exceptions.Internal),
		},
		{
			name: "TestGetErrorCode",
			args: args{
				err: &exceptions.CustomError{
					Err:     nil,
					Message: "",
					Code:    int(exceptions.PINError),
				},
			},
			want: int(exceptions.PINError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exceptions.GetErrorCode(tt.args.err); got != tt.want {
				t.Errorf("GetErrorCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestGetError",
			args: args{
				err: nil,
			},
			wantErr: false,
		},
		{
			name: "TestGetError",
			args: args{
				err: &exceptions.CustomError{
					Err:     nil,
					Message: "",
					Code:    int(exceptions.Internal),
				},
			},
			wantErr: false,
		},
		{
			name: "TestGetError",
			args: args{
				err: &exceptions.CustomError{
					Err:     nil,
					Message: "",
					Code:    int(exceptions.PINError),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := exceptions.GetError(tt.args.err); (err != nil) != tt.wantErr {
				t.Errorf("GetError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
