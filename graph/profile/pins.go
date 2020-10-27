package profile

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

// ValidateMsisdn checks that the msisdn supplied in the indicated request is valid
func ValidateMsisdn(w http.ResponseWriter, r *http.Request) (*PinRecovery, error) {
	data := &PinRecovery{}
	base.DecodeJSONToTargetStruct(w, r, data)
	if data.MSISDN == "" {
		err := fmt.Errorf("invalid credentials, expected a phone number")
		base.ReportErr(w, err, http.StatusBadRequest)
		return nil, err
	}
	return data, nil
}

// RequestPinResetFunc returns a function that sends an otp to an msisdn that requests a
// pin reset request during login
func RequestPinResetFunc(ctx context.Context, srv *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validMsisdn, err := ValidateMsisdn(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		msisdn := validMsisdn.MSISDN
		otp, err := srv.RequestPinReset(ctx, msisdn)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		otpResponse := OtpResponse{
			OTP: otp,
		}

		base.WriteJSONResponse(w, otpResponse, http.StatusOK)
	}
}
