package rest

import (
	"fmt"
	"net/http"

	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"github.com/savannahghi/serverutils"
)

// MyCareHubHandlersInterfaces represents all the REST API logic
type MyCareHubHandlersInterfaces interface {
	VerifyOTP() http.HandlerFunc
	LoginByPhone() http.HandlerFunc
	VerifySecurityQuestions() http.HandlerFunc
	VerifyPhone() http.HandlerFunc
	SendOTP() http.HandlerFunc
	RequestPINReset() http.HandlerFunc
	SendRetryOTP() http.HandlerFunc
}

// MyCareHubHandlersInterfacesImpl represents the usecase implementation object
type MyCareHubHandlersInterfacesImpl struct {
	usecase usecases.MyCareHub
}

// NewMyCareHubHandlersInterfaces initializes a new rest handlers usecase
func NewMyCareHubHandlersInterfaces(usecase usecases.MyCareHub) MyCareHubHandlersInterfaces {
	return &MyCareHubHandlersInterfacesImpl{usecase}
}

// LoginByPhone is an unauthenticated endpoint that gets the phonenumber and pin
// from a user, checks whether they exist, if present, we fetch the pin and if they match,
// we return the user profile and auth credentials to allow the user to login
func (h *MyCareHubHandlersInterfacesImpl) LoginByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.LoginInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == nil || payload.PIN == nil {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, responseCode, err := h.usecase.User.Login(ctx, *payload.PhoneNumber, *payload.PIN, payload.Flavour)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
				Code:    responseCode,
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, resp, http.StatusOK)
	}
}

// VerifySecurityQuestions get the user ID, question ID and the security question response from the payload and
// looks up the saved responses to determine whether the answers match to what has been stored. All of them must match.
// This is a security layer that will be used when a user attempts to reset their pin
func (h *MyCareHubHandlersInterfacesImpl) VerifySecurityQuestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payloadData := &[]dto.VerifySecurityQuestionInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payloadData)

		for _, payload := range *payloadData {
			err := payload.Validate()
			if err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}
		}

		ok, err := h.usecase.SecurityQuestions.VerifySecurityQuestionResponses(ctx, payloadData)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, ok, http.StatusOK)
	}
}

// VerifyPhone is an unauthenticated endpoint that does a check on the provided phone and flavour and
// performs a check to ascertain whether the supplied phone number and flavour are associated with the user.
func (h *MyCareHubHandlersInterfacesImpl) VerifyPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.VerifyPhoneInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.PhoneNumber == "" {
			err := fmt.Errorf("expected a phone input")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		otpResponse, err := h.usecase.OTP.VerifyPhoneNumber(ctx, payload.PhoneNumber, payload.Flavour)
		if err != nil {
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, otpResponse, http.StatusOK)
	}
}

// VerifyOTP is an unauthenticated endpoint that gets the phonenumber and flavour
// from a user, checks whether the provided otp matches. If they match, return true, otherwise false.
func (h *MyCareHubHandlersInterfacesImpl) VerifyOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.VerifyOTPInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.OTP == "" || payload.PhoneNumber == "" {
			err := fmt.Errorf("expected `userID`, `otp` and phone to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.OTP.VerifyOTP(ctx, payload)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, resp, http.StatusOK)
	}
}

// SendOTP is an unauthenticated endpoint that gets the phonenumber and flavour
// from a user and sends an OTP
func (h *MyCareHubHandlersInterfacesImpl) SendOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.SendOTPInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == "" {
			err := fmt.Errorf("expected `phone number` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.OTP.GenerateAndSendOTP(ctx, payload.PhoneNumber, payload.Flavour)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, resp, http.StatusOK)
	}
}

// RequestPINReset exposes an endpoint that takes in a user phonenumber and the flavour. It then sends
// an OTP to the phone number that requests a PIN reset
func (h *MyCareHubHandlersInterfacesImpl) RequestPINReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.SendOTPInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == "" {
			err := fmt.Errorf("expected `phone number` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.User.RequestPINReset(ctx, payload.PhoneNumber, payload.Flavour)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		serverutils.WriteJSONResponse(w, resp, http.StatusOK)
	}
}

// SendRetryOTP is an unauthenticated request that takes in a phone number
// generates an OTP and sends the OTP to the phone number
func (h *MyCareHubHandlersInterfacesImpl) SendRetryOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		retryPayload := &dto.SendRetryOTPPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, retryPayload)

		if retryPayload.Phone == "" || !retryPayload.Flavour.IsValid() {
			err := fmt.Errorf(
				"expected `phoneNumber`, `flavour` to be defined",
			)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.usecase.OTP.GenerateRetryOTP(ctx, retryPayload)
		if err != nil {
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}
