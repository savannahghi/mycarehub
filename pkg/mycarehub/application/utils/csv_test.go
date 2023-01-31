package utils

import (
	"testing"
)

func TestReadCSVFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: read csv file",
			args: args{
				path: "testData/facility.csv",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid Path",
			args: args{
				path: "invalidPath.csv",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadCSVFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadCSVFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got: %v", got)
			}
		})
	}
}

func TestParseFacilitiesFromCSV(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case: Parse facilities from csv",
			args: args{
				path: "testData/facility.csv",
			},
			wantErr: false,
		},
		{
			name: "Sad Case: invalid path",
			args: args{
				path: "testData/invalidIdentifier.csv",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFacilitiesFromCSV(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFacilitiesFromCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected value, got: %v", got)
			}
		})
	}
}
