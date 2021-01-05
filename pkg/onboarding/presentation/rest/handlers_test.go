package rest

import (
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestValidateSignUpPayload(t *testing.T) {
	testNumber := string(base.TestUserPhoneNumber)
	testPin := "3456"
	testInvalidPhone := ""
	invalidTestPin := "11"

	invalidCase1 := &resources.SignUpPayload{
		PhoneNumber: &testInvalidPhone,
		PIN:         &testPin,
		Flavour:     base.FlavourPro,
	}

	tests := []struct {
		name    string
		args    *resources.SignUpPayload
		want    *resources.SignUpPayload
		wantErr bool
	}{
		{
			name: "valid case",
			args: &resources.SignUpPayload{
				PhoneNumber: &testNumber,
				PIN:         &testPin,
				Flavour:     base.FlavourPro,
			},

			want: &resources.SignUpPayload{
				PhoneNumber: &testNumber,
				PIN:         &testPin,
				Flavour:     base.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "invalid phone number",
			args: invalidCase1,
			want: &resources.SignUpPayload{
				PhoneNumber: &testNumber,
				PIN:         &testPin,
				Flavour:     base.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "invalid pin",
			args: &resources.SignUpPayload{
				PhoneNumber: &testNumber,
				PIN:         &invalidTestPin,
				Flavour:     base.FlavourPro,
			},
			want: &resources.SignUpPayload{
				PhoneNumber: &testNumber,
				PIN:         &testPin,
				Flavour:     base.FlavourPro,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := ValidateSignUpPayload(tt.args)

			if tt.wantErr && err == nil {
				t.Errorf("error was expected")
				return
			}

			if (err != nil) && !tt.wantErr {
				t.Errorf("No error was expected, got err: %v", err)
				return
			}

			// check output is the expected
			if !tt.wantErr {
				if output.Flavour != tt.want.Flavour {
					t.Errorf("wanted %v, got : %v", tt.want.Flavour, output.Flavour)
				}
				if *(output.PhoneNumber) != *(tt.want.PhoneNumber) {
					t.Errorf("wanted %v, got : %v", *(tt.want.PhoneNumber), *(output.PhoneNumber))
				}
				// check output is the expected
				if *(output.PIN) != *(tt.want.PIN) {
					t.Errorf("wanted %v, got : %v", *(tt.want.PIN), *(output.PIN))
				}
			}

		})
	}
}
