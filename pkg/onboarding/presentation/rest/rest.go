package rest

import (
	"context"
	"fmt"
	"net/http"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/interactor"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

// VerifySignUpPhoneNumber is an unauthenticated endpoint that does a
// check on the supplied phone number asserting whether the phone is associated with
// a user profile. It check both the PRIMARY PHONE and SECONDARY PHONE NUMBER.
// If the phone number does not exist, it sends the OTP to the phone number
func VerifySignUpPhoneNumber(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := &domain.PhoneNumberPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}
		v, err := i.Signup.CheckPhoneExists(ctx, *p.PhoneNumber)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		if !v {
			base.ReportErr(w, fmt.Errorf("%v", base.PhoneNumberInUse), http.StatusBadRequest)
			return
		}

		// send otp to the phone number
		otp, err := i.Otp.GenerateAndSendOTP(ctx, *p.PhoneNumber)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, domain.OtpResponse{OTP: otp}, http.StatusOK)
	}
}

// CreateUserWithPhoneNumber is an unauthenticated endpoint that is called to create
func CreateUserWithPhoneNumber(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := &domain.SignUpPayload{}
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

		response, err := i.Signup.CreateUserByPhone(ctx, *p.PhoneNumber, *p.PIN, p.Flavour)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// UserRecoveryPhoneNumbers fetches the phone numbers associated with a profile for the purpose of account recovery.
// The returned phone numbers slice should be masked. E.G +254700***123
func UserRecoveryPhoneNumbers(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		p := &domain.PhoneNumberPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := i.Signup.GetUserRecoveryPhoneNumbers(ctx, *p.PhoneNumber)

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
func LoginByPhone(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &domain.LoginPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil || p.PIN == nil {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := i.Login.LoginByPhone(
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

// SetUserPIN is an unauthenticated  endpoint that saves user Pin
func SetUserPIN(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pin := &domain.SetPINRequest{}
		base.DecodeJSONToTargetStruct(w, r, pin)
		if pin.PhoneNumber == "" || pin.PIN == "" {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := i.UserPIN.SetUserPIN(
			ctx,
			pin.PIN,
			pin.PhoneNumber,
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
func RequestPINReset(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &domain.PhoneNumberPayload{}
		base.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		otp, err := i.UserPIN.RequestPINReset(ctx, *p.PhoneNumber)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(
			w,
			domain.OtpResponse{OTP: otp},
			http.StatusOK,
		)
	}
}

// ChangePin used to change/update a user's PIN
func ChangePin(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pin := &domain.SetPINRequest{}
		base.DecodeJSONToTargetStruct(w, r, pin)
		if pin.PhoneNumber == "" || pin.PIN == "" {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := i.UserPIN.ChangeUserPIN(
			ctx,
			pin.PIN,
			pin.PhoneNumber,
		)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// SendRetryOTPHandler is an unauthenticated request that takes in a phone number
// and a retry step (1 for sending an OTP via WhatsApp and 2 for Twilio Messages)
// and generates and sends a valid OTP to the phone number
func SendRetryOTPHandler(
	ctx context.Context,
	i *interactor.Interactor,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		retryPayload := &domain.SendRetryOTPPayload{}
		base.DecodeJSONToTargetStruct(w, r, retryPayload)
		if retryPayload.Phone == nil || retryPayload.RetryStep == nil {
			err := fmt.Errorf("expected `phoneNumber`, `retryStep` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := i.Otp.SendRetryOTP(
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
			domain.OtpResponse{OTP: response},
			http.StatusOK,
		)
	}
}

// ExchangeRefreshTokenForIDToken is an unauthenticated endpoint that
// takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func ExchangeRefreshTokenForIDToken(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &domain.RefreshToken{}
		if p.RefreshToken == nil {
			err := fmt.Errorf("expected `refreshToken` to be defined")
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		response, err := i.Login.RefreshToken(*p.RefreshToken)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(
			w,
			domain.AuthCredentialResponse{
				RefreshToken: response.RefreshToken,
				ExpiresIn:    response.ExpiresIn,
				IDToken:      response.IDToken,
			},
			http.StatusOK,
		)
	}
}

// FindSupplierByUIDHandler fetch supplier profile via REST
func FindSupplierByUIDHandler(ctx context.Context, i *interactor.Interactor) http.HandlerFunc {
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
		supplier, err = i.Supplier.FindSupplierByUID(newContext)

		if supplier == nil || err != nil {
			err := fmt.Errorf("supplier profile not found")
			base.ReportErr(w, err, http.StatusNotFound)
			return
		}

		base.WriteJSONResponse(w, supplier, http.StatusOK)
	}
}
