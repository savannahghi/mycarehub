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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"github.com/savannahghi/serverutils"
	"gopkg.in/go-playground/validator.v9"
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
	RegisterKenyaEMRPatients() http.HandlerFunc
	GetClientHealthDiaryEntries() http.HandlerFunc
	RegisteredFacilityPatients() http.HandlerFunc
	ServiceRequests() http.HandlerFunc
	CreateOrUpdateKenyaEMRAppointments() http.HandlerFunc
	CreatePinResetServiceRequest() http.HandlerFunc
	GetUserProfile() http.HandlerFunc
	AddClientFHIRID() http.HandlerFunc
	AddPatientsRecords() http.HandlerFunc
	SyncFacilities() http.HandlerFunc
	AddFacilityFHIRID() http.HandlerFunc
	AppointmentsServiceRequests() http.HandlerFunc
	DeleteUser() http.HandlerFunc
	FetchContactOrganisations() http.HandlerFunc
	Organisations() http.HandlerFunc
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

		err := payload.Validate()
		if err != nil {
			fields := ""
			for _, i := range err.(validator.ValidationErrors) {
				fields += fmt.Sprintf("%s, ", i.Field())
			}

			err := fmt.Errorf("expected %s to be defined", fields)
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

		response, successful := h.usecase.User.Login(ctx, payload)
		if !successful {
			serverutils.WriteJSONResponse(w, response, http.StatusBadRequest)
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

// SyncFacilities is an unauthenticated endpoint that returns a list of facilities
// that do not have an FHIR organisation ID
func (h *MyCareHubHandlersInterfacesImpl) SyncFacilities() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		err := h.usecase.Facility.SyncFacilities(ctx)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
				Code:    exceptions.GetErrorCode(err),
			}, http.StatusBadRequest)
			return
		}

		ok := okResp{
			Status: true,
		}

		serverutils.WriteJSONResponse(w, ok, http.StatusOK)
	}
}

// AddFacilityFHIRID is an authenticated endpoint used to update facility(also known as organization in FHIR), details from clinical service to MyCareHub
func (h *MyCareHubHandlersInterfacesImpl) AddFacilityFHIRID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		payload := &dto.UpdateFacilityPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.FacilityID == "" || payload.FHIROrganisationID == "" {
			err := fmt.Errorf("neither facility ID nor FHIR organisation ID can be empty")
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err := h.usecase.Facility.UpdateFacility(ctx, payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		ok := okResp{
			Status: true,
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

// SendOTP is an unauthenticated endpoint that gets the username and flavour
// from a user and sends an OTP
func (h *MyCareHubHandlersInterfacesImpl) SendOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.SendOTPInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.Username == "" {
			err := fmt.Errorf("expected `username` to be defined")
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

		resp, err := h.usecase.OTP.GenerateAndSendOTP(ctx, payload.Username, payload.Flavour)
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

// RequestPINReset exposes an endpoint that takes in a user username and the flavour. It then sends
// an OTP to the phone number that requests a PIN reset
func (h *MyCareHubHandlersInterfacesImpl) RequestPINReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.SendOTPInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.Username == "" {
			err := fmt.Errorf("expected `username` to be defined")
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

		resp, err := h.usecase.User.RequestPINReset(ctx, payload.Username, payload.Flavour)
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

// SendRetryOTP is an unauthenticated request that takes in a username
// generates an OTP and sends the OTP to the phone number of the user
func (h *MyCareHubHandlersInterfacesImpl) SendRetryOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		retryPayload := &dto.SendRetryOTPPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, retryPayload)
		if retryPayload.Username == "" || !retryPayload.Flavour.IsValid() {
			err := fmt.Errorf(
				"expected `username`, `flavour` to be defined",
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

		switch r.Method {
		case http.MethodGet:
			h.GetServiceRequestsForKenyaEMR(ctx, r, w)
		case http.MethodPost:
			h.UpdateServiceRequests(ctx, w, r)
		default:
			serverutils.WriteJSONResponse(w,
				serverutils.ErrorMap(fmt.Errorf("unsupported method")),
				http.StatusMethodNotAllowed,
			)
		}
	}
}

// DeleteUser is an unauthenticated endpoint that deletes a user from the system.
func (h *MyCareHubHandlersInterfacesImpl) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.PhoneInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == "" || payload.Flavour == "" {
			err := fmt.Errorf("expected `phoneNumber` and `flavour` to be defined")
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

		resp, err := h.usecase.User.DeleteUser(ctx, payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		response := helpers.RestAPIResponseHelper("deleteUser", resp)
		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// GetServiceRequestsForKenyaEMR gets all the service requests from MyCareHub
func (h *MyCareHubHandlersInterfacesImpl) GetServiceRequestsForKenyaEMR(ctx context.Context, r *http.Request, w http.ResponseWriter) {
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
		return
	}

	serviceRequests, err := h.usecase.ServiceRequest.GetServiceRequestsForKenyaEMR(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
		return
	}

	serverutils.WriteJSONResponse(w, serviceRequests, http.StatusOK)
}

// UpdateServiceRequests is an endpoint used to update service requests from KenyaEMR to MyCareHub
func (h *MyCareHubHandlersInterfacesImpl) UpdateServiceRequests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	payload := &dto.UpdateServiceRequestsPayload{}
	serverutils.DecodeJSONToTargetStruct(w, r, payload)

	if len(payload.ServiceRequests) == 0 {
		err := fmt.Errorf("no service requests payload defined")
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	_, err := h.usecase.ServiceRequest.UpdateServiceRequestsFromKenyaEMR(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
		return
	}

	serverutils.WriteJSONResponse(w, okResp{Status: true}, http.StatusOK)
}

// CreateOrUpdateKenyaEMRAppointments is tha handler used to sync appointments from Kenya EMR
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

		response, err := h.usecase.Appointment.CreateOrUpdateKenyaEMRAppointments(ctx, *payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
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

// AppointmentsServiceRequests is used to check for the oncoming request and routes it to the correct handler.
// If the method is POST, we update the appointment service request and if it's a GET, we return all the
// appointment service requests from the last time syncing occurred between the two platforms
func (h *MyCareHubHandlersInterfacesImpl) AppointmentsServiceRequests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		switch r.Method {
		case http.MethodGet:
			h.GetAppointmentServiceRequests(ctx, w, r)
		case http.MethodPost:
			h.UpdateServiceRequests(ctx, w, r)
		default:
			err := fmt.Errorf("wrong method passed")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
		}
	}
}

// GetAppointmentServiceRequests handler for syncing red-flags from the my carehub endpoint to Kenya EMR for display
func (h *MyCareHubHandlersInterfacesImpl) GetAppointmentServiceRequests(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

// FetchContactOrganisations fetches organisations associated with the provided contact
func (h *MyCareHubHandlersInterfacesImpl) FetchContactOrganisations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		phoneNumber := r.URL.Query().Get("phoneNumber")
		if phoneNumber == "" {
			err := fmt.Errorf("phone number is required")
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		organisations, err := h.usecase.User.FetchContactOrganisations(r.Context(), phoneNumber)
		if err != nil {
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		response := dto.OrganisationsOutput{
			Count:         0,
			Organisations: []dto.Organisation{},
		}

		for _, organisation := range organisations {
			org := dto.Organisation{
				ID:          organisation.ID,
				Name:        organisation.Name,
				Description: organisation.Description,
			}

			response.Count++
			response.Organisations = append(response.Organisations, org)
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// Organisations lists all organisations
func (h *MyCareHubHandlersInterfacesImpl) Organisations() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		organisationsOutput, err := h.usecase.Organisation.ListOrganisations(r.Context(), nil)
		if err != nil {
			serverutils.WriteJSONResponse(w, serverutils.ErrorMap(err), http.StatusBadRequest)
			return
		}

		response := dto.OrganisationsOutput{
			Count:         0,
			Organisations: []dto.Organisation{},
		}

		for _, organisation := range organisationsOutput.Organisations {
			org := dto.Organisation{
				ID:          organisation.ID,
				Name:        organisation.Name,
				Description: organisation.Description,
			}

			response.Count++
			response.Organisations = append(response.Organisations, org)
		}

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}
