package rest

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// HandlersInterfaces represents all the REST API logic
type HandlersInterfaces interface {
	VerifySignUpPhoneNumber(ctx context.Context) http.HandlerFunc
	CreateUserWithPhoneNumber(ctx context.Context) http.HandlerFunc
	UserRecoveryPhoneNumbers(ctx context.Context) http.HandlerFunc
	LoginByPhone(ctx context.Context) http.HandlerFunc
	RequestPINReset(ctx context.Context) http.HandlerFunc
	ChangePin(ctx context.Context) http.HandlerFunc
	SendRetryOTP(ctx context.Context) http.HandlerFunc
	RefreshToken(ctx context.Context) http.HandlerFunc
	FindSupplierByUID(ctx context.Context) http.HandlerFunc
}

// HandlersInterfacesImpl represents the usecase implementation object
type HandlersInterfacesImpl struct {
	interactor *interactor.Interactor
}

// NewHandlersInterfaces initializes a new rest handlers usecase
func NewHandlersInterfaces(i *interactor.Interactor) HandlersInterfaces {
	return &HandlersInterfacesImpl{i}
}

// VerifySignUpPhoneNumber is an unauthenticated endpoint that does a
// check on the supplied phone number asserting whether the phone is associated with
// a user profile. It check both the PRIMARY PHONE and SECONDARY PHONE NUMBER.
// If the phone number does not exist, it sends the OTP to the phone number
func (h *HandlersInterfacesImpl) VerifySignUpPhoneNumber(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := &resources.PhoneNumberPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		v, err := h.interactor.Signup.CheckPhoneExists(ctx, *p.PhoneNumber)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		if v {
			base.ReportErr(w, fmt.Errorf("%v", base.PhoneNumberInUse), http.StatusBadRequest)
			return
		}

		// send otp to the phone number
		otp, err := h.interactor.Otp.GenerateAndSendOTP(ctx, *p.PhoneNumber)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, resources.OtpResponse{OTP: otp}, http.StatusOK)
	}
}

// CreateUserWithPhoneNumber is an unauthenticated endpoint that is called to create
func (h *HandlersInterfacesImpl) CreateUserWithPhoneNumber(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := &resources.SignUpPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil && p.PIN == nil {
			err := fmt.Errorf("expected `phoneNumber` and `pin` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		if !p.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Signup.CreateUserByPhone(
			ctx,
			*p.PhoneNumber,
			*p.PIN,
			p.Flavour,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// UserRecoveryPhoneNumbers fetches the phone numbers associated with a profile for the purpose of account recovery.
// The returned phone numbers slice should be masked. E.G +254700***123
func (h *HandlersInterfacesImpl) UserRecoveryPhoneNumbers(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := &resources.PhoneNumberPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Signup.GetUserRecoveryPhoneNumbers(
			ctx,
			*p.PhoneNumber,
		)

		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// LoginByPhone is an unauthenticated endpoint that:
// Collects a phonenumber and pin from the user and checks if the phonenumber
// is an existing PRIMARY PHONENUMBER. If it does then it fetches the PIN that
// belongs to the profile and returns auth credentials to allow the user to login
func (h *HandlersInterfacesImpl) LoginByPhone(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &resources.LoginPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil || p.PIN == nil {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		if !p.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Login.LoginByPhone(
			ctx,
			*p.PhoneNumber,
			*p.PIN,
			p.Flavour,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RequestPINReset is an unauthenticated request that takes in a phone number
// sends an otp to an msisdn that requests a PIN reset request during login
func (h *HandlersInterfacesImpl) RequestPINReset(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &resources.PhoneNumberPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		otp, err := h.interactor.UserPIN.RequestPINReset(ctx, *p.PhoneNumber)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(
			w,
			resources.OtpResponse{OTP: otp},
			http.StatusOK,
		)
	}
}

// ChangePin used to change/update a user's PIN
func (h *HandlersInterfacesImpl) ChangePin(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pin := &resources.SetPINRequest{}
		base.DecodeJSONToTargetStruct(w, r, pin)
		if pin.PhoneNumber == "" || pin.PIN == "" {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.UserPIN.ChangeUserPIN(
			ctx,
			pin.PhoneNumber,
			pin.PIN,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// SendRetryOTP is an unauthenticated request that takes in a phone number
// and a retry step (1 for sending an OTP via WhatsApp and 2 for Twilio Messages)
// and generates and sends a valid OTP to the phone number
func (h *HandlersInterfacesImpl) SendRetryOTP(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retryPayload := &resources.SendRetryOTPPayload{}
		base.DecodeJSONToTargetStruct(w, r, retryPayload)
		if retryPayload.Phone == nil || retryPayload.RetryStep == nil {
			err := fmt.Errorf("expected `phoneNumber`, `retryStep` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Otp.SendRetryOTP(
			ctx,
			*retryPayload.Phone,
			*retryPayload.RetryStep,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(
			w,
			resources.OtpResponse{OTP: response},
			http.StatusOK,
		)
	}
}

// RefreshToken is an unauthenticated endpoint that
// takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (h *HandlersInterfacesImpl) RefreshToken(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &resources.RefreshTokenPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.RefreshToken == nil {
			err := fmt.Errorf("expected `refreshToken` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Login.RefreshToken(*p.RefreshToken)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(
			w,
			resources.AuthCredentialResponse{
				RefreshToken: response.RefreshToken,
				ExpiresIn:    response.ExpiresIn,
				IDToken:      response.IDToken,
			},
			http.StatusOK,
		)
	}
}

// FindSupplierByUID fetch supplier profile via REST
func (h *HandlersInterfacesImpl) FindSupplierByUID(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := utils.ValidateUID(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		if s.UID == nil {
			err := fmt.Errorf("expected `uid` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		var supplier *domain.Supplier

		newContext := context.WithValue(ctx, base.AuthTokenContextKey, s.UID)
		supplier, err = h.interactor.Supplier.FindSupplierByUID(newContext)

		if supplier == nil || err != nil {
			err := fmt.Errorf("supplier profile not found")
			base.ReportErr(w, err, http.StatusNotFound)
			return
		}

		base.WriteJSONResponse(w, supplier, http.StatusOK)
	}
}
