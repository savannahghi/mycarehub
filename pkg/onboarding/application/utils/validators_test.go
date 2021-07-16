package utils_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
)

func TestValidateSignUpInput(t *testing.T) {
	phone := interserviceclient.TestUserPhoneNumber
	pin := interserviceclient.TestUserPin
	flavour := feedlib.FlavourConsumer
	otp := "12345"

	alphanumericPhone := "+254-not-valid-123"
	badPhone := "+254712"
	shortPin := "123"
	longPin := "1234567"
	alphabeticalPin := "abcd"

	type args struct {
		input *dto.SignUpInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: return a valid output",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &pin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: false,
		},
		{
			name: "failure: bad phone number provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &badPhone,
					PIN:         &pin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: alphanumeric phone number provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &alphanumericPhone,
					PIN:         &pin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: short pin number provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &shortPin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: long pin number provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &longPin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: alphabetical pin number provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &alphabeticalPin,
					Flavour:     flavour,
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: bad flavour provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &pin,
					Flavour:     "not-a-flavour",
					OTP:         &otp,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: no OTP provided",
			args: args{
				input: &dto.SignUpInput{
					PhoneNumber: &phone,
					PIN:         &pin,
					Flavour:     flavour,
					OTP:         nil,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validInput, err := utils.ValidateSignUpInput(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSignUpInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && validInput != nil {
				t.Errorf("expected a nil valid input since an error :%v occurred", err)
			}

			if err == nil && validInput == nil {
				t.Errorf("expected a valid input %v since no error occurred", validInput)
			}
		})
	}
}

func TestValidateUID(t *testing.T) {
	tests := []struct {
		name    string
		args    map[string]interface{}
		wantErr bool
	}{
		{
			name: "valid",
			args: map[string]interface{}{
				"uid": uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "invalid",
			args: map[string]interface{}{
				"uuid": uuid.New().String(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.args)
			if err != nil {
				t.Errorf("failed to marshal body: %v", err)
				return
			}
			// Create a request to pass to our handler.
			req, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer(body))
			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}
			rw := httptest.NewRecorder()
			resp, err := utils.ValidateUID(rw, req)
			if tt.wantErr {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
			}
			if !tt.wantErr {
				assert.Nil(t, err)
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.UID)
			}

		})
	}
}

func TestValidateSMSData(t *testing.T) {
	validLinkID := uuid.New().String()
	text := "Test Covers"
	to := "3601"
	id := "60119"
	from := "+254705385894"
	date := "2021-05-17T13:20:04.490Z"

	// valid payload
	validSMSData := &dto.AfricasTalkingMessage{
		LinkID: validLinkID,
		Text:   text,
		To:     to,
		ID:     id,
		Date:   date,
		From:   from,
	}

	invalidData := &dto.AfricasTalkingMessage{
		LinkID: " ",
		Text:   text,
		To:     to,
		ID:     id,
		Date:   date,
		From:   from,
	}

	type args struct {
		input *dto.AfricasTalkingMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case :) Return valid input",
			args: args{
				input: validSMSData,
			},
			wantErr: false,
		},
		{
			name: "Sad case :( Return invalid input",
			args: args{
				input: invalidData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validatedData, err := utils.ValidateAficasTalkingSMSData(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSMSData() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && validatedData != nil {
				t.Errorf("the error was not expected")
				return
			}

			if !tt.wantErr && validatedData == nil {
				t.Errorf("an error was expected: %v", err)
				return
			}
		})
	}
}

func TestValidateUSSDDetails(t *testing.T) {
	phone := "+254711223344"
	sessionId := "1235678"
	text := ""

	alphanumericPhone := "+254-not-valid-123"
	emptySessionId := ""
	badPhone := "+254712"

	type args struct {
		input *dto.SessionDetails
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: return a valid output",
			args: args{
				input: &dto.SessionDetails{
					PhoneNumber: &phone,
					SessionID:   sessionId,
					Text:        text,
				},
			},
			wantErr: false,
		},
		{
			name: "failure: bad phone number provided",
			args: args{
				input: &dto.SessionDetails{
					PhoneNumber: &badPhone,
					SessionID:   sessionId,
					Text:        text,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: alphanumeric phone number provided",
			args: args{
				input: &dto.SessionDetails{
					PhoneNumber: &alphanumericPhone,
					SessionID:   sessionId,
					Text:        text,
				},
			},
			wantErr: true,
		},
		{
			name: "failure: empty Ussd SessionId",
			args: args{
				input: &dto.SessionDetails{
					PhoneNumber: &phone,
					SessionID:   emptySessionId,
					Text:        text,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validInput, err := utils.ValidateUSSDDetails(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUSSDDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && validInput != nil {
				t.Errorf("expected a nil valid input since an error :%v occurred", err)
			}

			if err == nil && validInput == nil {
				t.Errorf("expected a valid input %v since no error occurred", validInput)
			}
		})
	}
}
