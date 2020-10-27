package profile

import (
	"context"
	"net/http"

	"gitlab.slade360emr.com/go/base"
)

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

// UpdatePinHandler used to update a user's PIN
func UpdatePinHandler(ctx context.Context, srv *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, validateErr := ValidateUpdatePinPayload(w, r)
		if validateErr != nil {
			base.ReportErr(w, validateErr, http.StatusBadRequest)
			return
		}

		_, updateErr := srv.UpdateUserPin(ctx, payload.MSISDN, payload.PIN, payload.OTP)
		if updateErr != nil {
			base.ReportErr(w, updateErr, http.StatusBadRequest)
			return
		}

		type okResp struct {
			Status string `json:"status"`
		}

		base.WriteJSONResponse(w, okResp{Status: "ok"}, http.StatusOK)

	}
}
