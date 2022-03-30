package rest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
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
	GetClientHealthDiaryEntries() http.HandlerFunc
	RegisteredFacilityPatients() http.HandlerFunc
	ServiceRequests() http.HandlerFunc
	CreateOrUpdateKenyaEMRAppointments() http.HandlerFunc
	CreatePinResetServiceRequest() http.HandlerFunc
	OptIn() http.HandlerFunc
	GetUserProfile() http.HandlerFunc
	AddClientFHIRID() http.HandlerFunc
	AddPatientsRecords() http.HandlerFunc
	GetAppointmentServiceRequests() http.HandlerFunc
}

type okResp struct {
	Status bool `json:"status"`
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

		response, err := h.usecase.User.Login(ctx, *payload.PhoneNumber, *payload.PIN, payload.Flavour)
		if err != nil {
			resp := &domain.CustomResponse{
				Message: err.Error(),
				Code:    response.Code,
			}

			if response.RetryTime != 0 {
				resp.RetryTime = response.RetryTime
			}

			if response.Attempts != 0 {
				resp.FailedLoginCount = response.Attempts
			}

			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, resp, http.StatusBadRequest)
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
					Err:     exceptions.InternalErr(err),
					Message: err.Error(),
					Code:    int(exceptions.Internal),
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
				Code:    exceptions.GetErrorCode(err),
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

		payload := &dto.PatientsPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if len(payload.Patients) == 0 {
			err := fmt.Errorf("expected at least one patient")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		response, err := h.usecase.User.RegisterKenyaEMRPatients(ctx, payload.Patients)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// GetClientHealthDiaryEntries fetches and returns the health diary entries that were recorded
// in the specified facility.
func (h *MyCareHubHandlersInterfacesImpl) GetClientHealthDiaryEntries() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		MFLCode, err := strconv.Atoi(r.URL.Query().Get("MFLCODE"))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		syncTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("lastSyncTime"))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		payload := &dto.FetchHealthDiaryEntries{
			MFLCode:      MFLCode,
			LastSyncTime: &syncTime,
		}

		if payload.MFLCode == 0 || payload.LastSyncTime == nil {
			err := fmt.Errorf("expected `MFLCODE` and `lastSyncTime` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.usecase.HealthDiary.GetFacilityHealthDiaryEntries(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RegisteredFacilityPatients handler for syncing newly registered patients for a facility
func (h *MyCareHubHandlersInterfacesImpl) RegisteredFacilityPatients() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		MFLCode, err := strconv.Atoi(r.URL.Query().Get("MFLCODE"))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		syncTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("lastSyncTime"))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		payload := &dto.PatientSyncPayload{
			MFLCode:  MFLCode,
			SyncTime: &syncTime,
		}

		if payload.MFLCode == 0 {
			err := fmt.Errorf("expected `MFLCODE` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		response, err := h.usecase.User.RegisteredFacilityPatients(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// ServiceRequests is the endpoint used to sync service requests from MyCareHub to KenyaEMR
func (h *MyCareHubHandlersInterfacesImpl) ServiceRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if r.Method == http.MethodGet {
			err := h.GetServiceRequestsForKenyaEMR(ctx, r, w)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return
			}
		}

		if r.Method == http.MethodPost {
			err := h.UpdateServiceRequests(ctx, w, r)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				return
			}
		}

	}
}

// GetServiceRequestsForKenyaEMR gets all the service requests from MyCareHub
func (h *MyCareHubHandlersInterfacesImpl) GetServiceRequestsForKenyaEMR(ctx context.Context, r *http.Request, w http.ResponseWriter) error {
	MFLCode, err := strconv.Atoi(r.URL.Query().Get("MFLCODE"))
	if err != nil {
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return err
	}
	syncTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("lastSyncTime"))
	if err != nil {
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return err
	}
	payload := &dto.ServiceRequestPayload{
		MFLCode:      MFLCode,
		LastSyncTime: &syncTime,
	}
	if payload.MFLCode == 0 || payload.LastSyncTime == nil {
		err := fmt.Errorf("expected `MFLCODE` and `lastSyncTime` to be defined")
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return err
	}

	serviceRequests, err := h.usecase.ServiceRequest.GetServiceRequestsForKenyaEMR(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
		return err
	}

	serverutils.WriteJSONResponse(w, serviceRequests, http.StatusOK)
	return nil
}

//UpdateServiceRequests is an endpoint used to update service requests from KenyaEMR to MyCareHub
func (h *MyCareHubHandlersInterfacesImpl) UpdateServiceRequests(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	payload := &dto.UpdateServiceRequestsPayload{}
	serverutils.DecodeJSONToTargetStruct(w, r, payload)

	if len(payload.ServiceRequests) == 0 {
		err := fmt.Errorf("no service requests payload defined")
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return err
	}

	serviceRequests, err := h.usecase.ServiceRequest.UpdateServiceRequestsFromKenyaEMR(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
	}

	serverutils.WriteJSONResponse(w, serviceRequests, http.StatusOK)
	return nil
}

// CreateOrUpdateKenyaEMRAppointments is tha handler used to sync appointmens from Kenya EMR
// The appointment can be a POST, handled as a create or PUT handled as an update to existing appointment
func (h *MyCareHubHandlersInterfacesImpl) CreateOrUpdateKenyaEMRAppointments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		payload := &dto.FacilityAppointmentsPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.MFLCode == "" {
			err := fmt.Errorf("expected an MFL code to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		if len(payload.Appointments) == 0 {
			err := fmt.Errorf("expected at least one appointment to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		if r.Method == http.MethodPost {
			response, err := h.usecase.Appointment.CreateKenyaEMRAppointments(ctx, *payload)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
				return
			}

			serverutils.WriteJSONResponse(w, response, http.StatusCreated)
			return
		}

		if r.Method == http.MethodPatch {
			response, err := h.usecase.Appointment.CreateKenyaEMRAppointments(ctx, *payload)
			if err != nil {
				helpers.ReportErrorToSentry(err)
				serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
				return
			}

			serverutils.WriteJSONResponse(w, response, http.StatusOK)
			return
		}
	}
}

// CreatePinResetServiceRequest is used to create a "PIN_RESET" service request. This is trigerred
// when a user has failed loggin in to the app and requests for help. The service request will be viewed
// by the healthcare worker and either approved/rejected
func (h *MyCareHubHandlersInterfacesImpl) CreatePinResetServiceRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.PinResetServiceRequestPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if !payload.Flavour.IsValid() || payload.PhoneNumber == "" {
			err := fmt.Errorf("expected a valid `flavour` or `phoneNumber` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.usecase.ServiceRequest.CreatePinResetServiceRequest(ctx, payload.PhoneNumber, payload.CCCNumber, payload.Flavour)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, okResp{Status: response}, http.StatusOK)
	}
}

// OptIn will be used by users to take an affirmative action to offer their consent to the app
func (h *MyCareHubHandlersInterfacesImpl) OptIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.OptInPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.Flavour == "" || payload.PhoneNumber == "" {
			err := fmt.Errorf("expected both `flavour` and `phoneNumber` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.usecase.User.Consent(ctx, payload.PhoneNumber, payload.Flavour, true)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, okResp{Status: response}, http.StatusOK)
	}
}

// GetUserProfile returns a user profile given the user ID
func (h *MyCareHubHandlersInterfacesImpl) GetUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID := vars["id"]

		ctx := r.Context()

		user, err := h.usecase.User.GetUserProfile(ctx, userID)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, user, http.StatusOK)
	}
}

// AddPatientsRecords handles bulk creation of patient records
func (h *MyCareHubHandlersInterfacesImpl) AddPatientsRecords() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := &dto.PatientsRecordsPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.MFLCode == "" {
			err := fmt.Errorf("expected an MFL code to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if len(payload.Records) == 0 {
			err := fmt.Errorf("expected at least one record to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err := h.usecase.Appointment.AddPatientsRecords(r.Context(), *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, okResp{Status: true}, http.StatusCreated)
	}
}

// AddClientFHIRID adds the created fhir ID to a client profile
func (h *MyCareHubHandlersInterfacesImpl) AddClientFHIRID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.ClientFHIRPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.ClientID == "" || payload.FHIRID == "" {
			err := fmt.Errorf("expected both `client ID` and `fhir ID` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err := h.usecase.User.AddClientFHIRID(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, okResp{Status: true}, http.StatusOK)

	}
}

// GetAppointmentServiceRequests handler for syncing red-flags from the my carehub endpoint to Kenya EMR for display
func (h *MyCareHubHandlersInterfacesImpl) GetAppointmentServiceRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		MFLCode, err := strconv.Atoi(r.URL.Query().Get("MFLCODE"))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		lastSyncTime, err := time.Parse(time.RFC3339, r.URL.Query().Get("lastSyncTime"))
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		payload := &dto.AppointmentServiceRequestInput{
			MFLCode:      MFLCode,
			LastSyncTime: &lastSyncTime,
		}

		if payload.MFLCode == 0 {
			err := fmt.Errorf("expected `MFLCODE` to be defined")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		response, err := h.usecase.Appointment.GetAppointmentServiceRequests(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}
