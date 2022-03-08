package rest

import (
	"fmt"
	"net/http"

	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
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
	GetUserRespondedSecurityQuestions() http.HandlerFunc
	ResetPIN() http.HandlerFunc
	RefreshToken() http.HandlerFunc
	RefreshGetStreamToken() http.HandlerFunc
	RegisterKenyaEMRPatients() http.HandlerFunc
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
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, responseCode, err := h.usecase.User.Login(ctx, *payload.PhoneNumber, *payload.PIN, payload.Flavour)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
				Code:    responseCode,
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// VerifySecurityQuestions get the user ID, question ID and the security question response from the payload and
// looks up the saved responses to determine whether the answers match to what has been stored. All of them must match.
// This is a security layer that will be used when a user attempts to reset their pin
func (h *MyCareHubHandlersInterfacesImpl) VerifySecurityQuestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payloadData := &dto.VerifySecurityQuestionsPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payloadData)

		for _, payload := range payloadData.SecurityQuestionsInput {
			err := payload.Validate()
			if err != nil {
				helpers.ReportErrorToSentry(err)
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}
		}

		ok, err := h.usecase.SecurityQuestions.VerifySecurityQuestionResponses(ctx, payloadData)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("verifySecurityQuestionResponses", ok)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
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
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		otpResponse, err := h.usecase.OTP.VerifyPhoneNumber(ctx, payload.PhoneNumber, payload.Flavour)
		if err != nil {
			helpers.ReportErrorToSentry(err)
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
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.OTP.VerifyOTP(ctx, payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("verifyOTP", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
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
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.OTP.GenerateAndSendOTP(ctx, payload.PhoneNumber, payload.Flavour)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("sendOTP", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
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
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.User.RequestPINReset(ctx, payload.PhoneNumber, payload.Flavour)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("requestPINReset", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
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
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.OTP.GenerateRetryOTP(ctx, retryPayload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("sendRetryOTP", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// GetUserRespondedSecurityQuestions is an unauthenticated endpoint that gets the user id and returns the security questions
// associated with the user.
func (h *MyCareHubHandlersInterfacesImpl) GetUserRespondedSecurityQuestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.GetUserRespondedSecurityQuestionsInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		err := payload.Validate()
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.SecurityQuestions.GetUserRespondedSecurityQuestions(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("getUserRespondedSecurityQuestions", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// ResetPIN is an unauthenticated endpoint that resets the user's PIN
func (h *MyCareHubHandlersInterfacesImpl) ResetPIN() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.UserResetPinInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		err := payload.Validate()
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, err := h.usecase.User.ResetPIN(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response := helpers.RestAPIResponseHelper("resetPIN", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RefreshToken is an unauthenticated endpoint that
// takes a user ID and creates a custom Firebase refresh token. It then tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (h *MyCareHubHandlersInterfacesImpl) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		payload := &dto.RefreshTokenPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.UserID == nil {
			err := fmt.Errorf("expected `userID` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.usecase.User.RefreshToken(ctx, *payload.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RefreshGetStreamToken takes a userID and returns a valid getstream token
func (h *MyCareHubHandlersInterfacesImpl) RefreshGetStreamToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		payload := &dto.RefreshTokenPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.UserID == nil {
			err := fmt.Errorf("expected `userID` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.usecase.User.RefreshGetStreamToken(ctx, *payload.UserID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RegisterKenyaEMRPatients is the handler for registering patients from KenyaEMR as clients
// It accepts multiple record for registration.
func (h *MyCareHubHandlersInterfacesImpl) RegisterKenyaEMRPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var payload []*dto.PatientRegistrationPayload

		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		// error decoding json
		if payload == nil {
			return
		}

		response, err := h.usecase.User.RegisterKenyaEMRPatients(ctx, payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
	}
}
