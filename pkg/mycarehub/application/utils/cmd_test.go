package utils

import (
	"testing"
)

func TestParseChoice(t *testing.T) {
	type args struct {
		choices []map[string]interface{}
		choice  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: valid Index number",

			args: args{
				choices: []map[string]interface{}{
					{"1": "item 1"},
					{"2": "item 2"},
				},
				choice: "0",
			},
			wantErr: false,
		},
		{
			name: "Sad case: index out of range",

			args: args{
				choices: []map[string]interface{}{
					{"1": "item 1"},
					{"2": "item 2"},
				},
				choice: "2",
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing Index number",
			args: args{
				choices: []map[string]interface{}{
					{"1": "item 1"},
					{"2": "item 2"},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: invalid input",

			args: args{
				choices: []map[string]interface{}{
					{"1": "item 1"},
					{"2": "item 2"},
				},
				choice: "zero",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseChoice(tt.args.choices, tt.args.choice)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseChoice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ParseIndex() = %v", got)
			}
		})
	}
}
