package utils_test

import (
	"strings"
	"testing"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

func TestResponseMenu(t *testing.T) {

	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success: text with input",
			args: args{
				input: "1",
			},
			want: "END Thank you",
		},
		{
			name: "success: text with invalid input",
			args: args{
				input: "7",
			},
			want: "CON Invalid choice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := utils.ResponseMenu(tt.args.input)
			if strings.Contains(resp, tt.want) != true {
				t.Errorf("expected %v to be in  %v  ", tt.want, resp)
				return
			}
		})
	}
}

func TestDefaultResponseMenu(t *testing.T) {

	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success: text with input",
			args: args{
				input: "1",
			},
			want: "END Thank you",
		},
		{
			name: "success: text with invalid input",
			args: args{
				input: "7",
			},
			want: "CON Invalid choice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := utils.ResponseMenu(tt.args.input)
			if strings.Contains(resp, tt.want) != true {
				t.Errorf("expected %v to be in  %v  ", tt.want, resp)
				return
			}
		})
	}
}

func TestDefaultMenu(t *testing.T) {

	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success: text length is 1 ",
			args: args{
				input: "1",
			},
			want: "CON Invalid choice",
		},
		{
			name: "success: text length is more than one",
			args: args{
				input: "7*1",
			},
			want: "END Thank you",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := utils.DefaultMenu(tt.args.input)
			if strings.Contains(resp, tt.want) != true {
				t.Errorf("expected %v to be in  %v  ", tt.want, resp)
				return
			}
		})
	}
}
