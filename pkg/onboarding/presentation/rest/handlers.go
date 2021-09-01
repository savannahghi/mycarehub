package rest

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"firebase.google.com/go/auth"
	"github.com/gorilla/mux"
	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/serverutils"
)

// HandlersInterfaces represents all the REST API logic
type HandlersInterfaces interface {
	VerifySignUpPhoneNumber() http.HandlerFunc
	CreateUserWithPhoneNumber() http.HandlerFunc
	UserRecoveryPhoneNumbers() http.HandlerFunc
	SetPrimaryPhoneNumber() http.HandlerFunc
	LoginByPhone() http.HandlerFunc
	LoginAnonymous() http.HandlerFunc
	RequestPINReset() http.HandlerFunc
	ResetPin() http.HandlerFunc
	SendOTP() http.HandlerFunc
	SendRetryOTP() http.HandlerFunc
	RefreshToken() http.HandlerFunc
	RemoveUserByPhoneNumber() http.HandlerFunc
	GetUserProfileByUID() http.HandlerFunc
	GetUserProfileByPhoneOrEmail() http.HandlerFunc
	ProfileAttributes() http.HandlerFunc
	RegisterPushToken() http.HandlerFunc
	AddAdminPermsToUser() http.HandlerFunc
	RemoveAdminPermsToUser() http.HandlerFunc
	AddRoleToUser() http.HandlerFunc
	RemoveRoleToUser() http.HandlerFunc
	UpdateUserProfile() http.HandlerFunc
	SwitchFlaggedFeaturesHandler() http.HandlerFunc
	PollServices() http.HandlerFunc
	CheckHasPermission() http.HandlerFunc

	CreateRole() http.HandlerFunc
	AssignRole() http.HandlerFunc
	RemoveRoleByName() http.HandlerFunc
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
func (h *HandlersInterfacesImpl) VerifySignUpPhoneNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p, err := decodeOTPPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		otpResp, err := h.interactor.Signup.VerifyPhoneNumber(
			ctx,
			*p.PhoneNumber,
			p.AppID,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("verify phone number OTP response", trace.WithAttributes(
			attribute.Any("response", otpResp),
		))

		serverutils.WriteJSONResponse(w, otpResp, http.StatusOK)
	}
}

// CreateUserWithPhoneNumber is an unauthenticated endpoint that is called to create
func (h *HandlersInterfacesImpl) CreateUserWithPhoneNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p := &dto.SignUpInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", p),
		))

		response, err := h.interactor.Signup.CreateUserByPhone(ctx, p)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("create user by phone", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// UserRecoveryPhoneNumbers fetches the phone numbers associated with a profile for the purpose of account recovery.
// The returned phone numbers slice should be masked. E.G +254700***123
func (h *HandlersInterfacesImpl) UserRecoveryPhoneNumbers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p, err := decodePhoneNumberPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Signup.GetUserRecoveryPhoneNumbers(
			ctx,
			*p.PhoneNumber,
		)

		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent(
			"retrieve user recovery phone numbers",
			trace.WithAttributes(
				attribute.Any("response", response),
			),
		)

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// SetPrimaryPhoneNumber sets the provided phone number as the primary phone of the profile associated with it
func (h *HandlersInterfacesImpl) SetPrimaryPhoneNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p := &dto.SetPrimaryPhoneNumberPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", p),
		))

		if p.PhoneNumber == nil || p.OTP == nil {
			err := fmt.Errorf("expected `phoneNumber` and `otp` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Signup.SetPhoneAsPrimary(
			ctx,
			*p.PhoneNumber,
			*p.OTP,
		)

		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("setting primary phone number", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(
			w,
			dto.NewOKResp(response),
			http.StatusOK,
		)
	}
}

// LoginByPhone is an unauthenticated endpoint that:
// Collects a phonenumber and pin from the user and checks if the phonenumber
// is an existing PRIMARY PHONENUMBER. If it does then it fetches the PIN that
// belongs to the profile and returns auth credentials to allow the user to login
func (h *HandlersInterfacesImpl) LoginByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p := &dto.LoginPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", p),
		))

		if p.PhoneNumber == nil || p.PIN == nil {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !p.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Login.LoginByPhone(
			ctx,
			*p.PhoneNumber,
			*p.PIN,
			p.Flavour,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}
		span.AddEvent("login by phone response", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// LoginAnonymous is an unauthenticated endpoint that returns only auth credentials for anonymous users
func (h *HandlersInterfacesImpl) LoginAnonymous() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p := &dto.LoginPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", p),
		))

		if p.Flavour.String() == "" {
			err := fmt.Errorf("expected `flavour` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !p.Flavour.IsValid() || p.Flavour != feedlib.FlavourConsumer {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Login.LoginAsAnonymous(ctx)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("log in as anonymous", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RequestPINReset is an unauthenticated request that takes in a phone number
// sends an otp to an msisdn that requests a PIN reset request during login
func (h *HandlersInterfacesImpl) RequestPINReset() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p, err := decodeOTPPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		otpResp, err := h.interactor.UserPIN.RequestPINReset(
			ctx,
			*p.PhoneNumber,
			p.AppID,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("request pin reset otp response", trace.WithAttributes(
			attribute.Any("response", otpResp),
		))

		serverutils.WriteJSONResponse(w, otpResp, http.StatusOK)
	}
}

// ResetPin used to change/update a user's PIN
func (h *HandlersInterfacesImpl) ResetPin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		pin := &dto.ChangePINRequest{}
		serverutils.DecodeJSONToTargetStruct(w, r, pin)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", pin),
		))

		if pin.PhoneNumber == "" || pin.PIN == "" || pin.OTP == "" {
			err := fmt.Errorf(
				"expected `phoneNumber`, `PIN` to be defined, `OTP` to be defined",
			)
			serverutils.WriteJSONResponse(
				w,
				errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				},
				http.StatusBadRequest,
			)
			return
		}

		response, err := h.interactor.UserPIN.ResetUserPIN(
			ctx,
			pin.PhoneNumber,
			pin.PIN,
			pin.OTP,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("reset user pin success", trace.WithAttributes(
			attribute.Bool("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// SendOTP is an unauthenticated request that takes in a phone number
// and generates an OTP and sends a valid OTP to the phone number. This API will mostly be used
// during account recovery workflow
func (h *HandlersInterfacesImpl) SendOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		payload, err := decodeOTPPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Engagement.GenerateAndSendOTP(
			ctx,
			*payload.PhoneNumber,
			payload.AppID,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("generate and send otp response", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// SendRetryOTP is an unauthenticated request that takes in a phone number
// and a retry step (1 for sending an OTP via WhatsApp and 2 for Twilio Messages)
// and generates and sends a valid OTP to the phone number
func (h *HandlersInterfacesImpl) SendRetryOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		retryPayload := &dto.SendRetryOTPPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, retryPayload)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", retryPayload),
		))

		if retryPayload.Phone == nil || retryPayload.RetryStep == nil {
			err := fmt.Errorf(
				"expected `phoneNumber`, `retryStep` to be defined",
			)
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Engagement.SendRetryOTP(
			ctx,
			*retryPayload.Phone,
			*retryPayload.RetryStep,
			retryPayload.AppID,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("send retry OTP", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RefreshToken is an unauthenticated endpoint that
// takes a custom Firebase refresh token and tries to fetch
// an ID token and returns auth credentials if successful
// Otherwise, an error is returned
func (h *HandlersInterfacesImpl) RefreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p := &dto.RefreshTokenPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)

		span.AddEvent("decode json payload to struct", trace.WithAttributes(
			attribute.Any("payload", p),
		))

		if p.RefreshToken == nil {
			err := fmt.Errorf("expected `refreshToken` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		response, err := h.interactor.Login.RefreshToken(ctx, *p.RefreshToken)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		span.AddEvent("new token", trace.WithAttributes(
			attribute.Any("response", response),
		))

		serverutils.WriteJSONResponse(w, response, http.StatusOK)
	}
}

// RemoveUserByPhoneNumber is an unauthenticated endpoint that removes a user
// whose phone number, either PRIMARY PHONE NUMBER or SECONDARY PHONE NUMBERS,matches the provided
// phone number in the request. This endpoint will ONLY be available under testing environment
func (h *HandlersInterfacesImpl) RemoveUserByPhoneNumber() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p, err := decodePhoneNumberPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		v, err := h.interactor.Onboarding.CheckPhoneExists(ctx, *p.PhoneNumber)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if v {
			if err := h.interactor.Signup.RemoveUserByPhoneNumber(ctx, *p.PhoneNumber); err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}
			serverutils.WriteJSONResponse(
				w,
				dto.OKResp{Status: "OK"},
				http.StatusOK,
			)
			return
		}
		err = fmt.Errorf(
			"`phoneNumber` does not exist and not associated with any user ",
		)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
	}
}

// GetUserProfileByUID fetches and returns a user profile via REST ISC
func (h *HandlersInterfacesImpl) GetUserProfileByUID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p := &dto.UIDPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)

		if p.UID == nil {
			err := fmt.Errorf("expected `UID` to be defined")
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		profile, err := h.interactor.Onboarding.GetUserProfileByUID(
			ctx,
			*p.UID,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, profile, http.StatusOK)
	}
}

// GetUserProfileByPhoneOrEmail fetches and returns a user profile via REST ISC
func (h *HandlersInterfacesImpl) GetUserProfileByPhoneOrEmail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.RetrieveUserProfileInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)

		if payload.Email == nil && payload.PhoneNumber == nil {
			err := fmt.Errorf("expected both `phonenumber` and `email` to be defined")
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		if payload.Email != nil && payload.PhoneNumber != nil {
			err := fmt.Errorf(
				"only one parameter can be used to retrieve user profile. use either email or phone",
			)
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		profile, err := h.interactor.Onboarding.GetUserProfileByPhoneOrEmail(
			ctx,
			payload,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, profile, http.StatusOK)
	}
}

// RegisterPushToken adds a new push token in the users profile if the push token does not exist
// via REST ISC
func (h *HandlersInterfacesImpl) RegisterPushToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		t := &dto.PushTokenPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, t)

		if t.PushToken == "" || t.UID == "" {
			err := fmt.Errorf("expected `PushToken` or `UID` to be defined")
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		authCred := &auth.Token{UID: t.UID}
		ctx = context.WithValue(
			ctx,
			firebasetools.AuthTokenContextKey,
			authCred,
		)

		profile, err := h.interactor.Signup.RegisterPushToken(
			ctx,
			t.PushToken,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, profile, http.StatusOK)
	}
}

// ProfileAttributes retrieves confirmed user profile attributes.
// These attributes include a user's verified phone numbers, verified emails
// and verified FCM push tokens
func (h *HandlersInterfacesImpl) ProfileAttributes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		attribute, found := vars["attribute"]
		if !found {
			err := fmt.Errorf("request does not have a path var named `%s`",
				attribute,
			)
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		p := &dto.UIDsPayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)
		if len(p.UIDs) == 0 {
			err := fmt.Errorf("expected a `UID` to be defined")
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		output, err := h.interactor.Onboarding.ProfileAttributes(
			ctx,
			p.UIDs,
			attribute,
		)
		if err != nil {
			errorcodeutil.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, output, http.StatusOK)
	}
}

// AddAdminPermsToUser is authenticated endpoint that adds admin permissions to a
// whose phone number, either PRIMARY PHONE NUMBER or SECONDARY PHONE NUMBERS,matches
// the provided phone number in the request.
func (h *HandlersInterfacesImpl) AddAdminPermsToUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		p, err := decodePhoneNumberPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		v, err := h.interactor.Onboarding.CheckPhoneExists(ctx, *p.PhoneNumber)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if v {
			if err := h.interactor.Onboarding.AddAdminPermsToUser(ctx, *p.PhoneNumber); err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}
			serverutils.WriteJSONResponse(
				w,
				dto.OKResp{Status: "OK"},
				http.StatusOK,
			)
			return
		}
		err = fmt.Errorf(
			"`phoneNumber` does not exist and not associated with any user ",
		)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
	}
}

// RemoveAdminPermsToUser is authenticated endpoint that removes admin permissions to a
// whose phone number, either PRIMARY PHONE NUMBER or SECONDARY PHONE NUMBERS,matches
// the provided phone number in the request.
func (h *HandlersInterfacesImpl) RemoveAdminPermsToUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)

		p, err := decodePhoneNumberPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		v, err := h.interactor.Onboarding.CheckPhoneExists(ctx, *p.PhoneNumber)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if v {
			if err := h.interactor.Onboarding.RemoveAdminPermsToUser(ctx, *p.PhoneNumber); err != nil {
				serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
					Err:     err,
					Message: err.Error(),
				}, http.StatusBadRequest)
				return
			}
			serverutils.WriteJSONResponse(
				w,
				dto.OKResp{Status: "OK"},
				http.StatusOK,
			)
			return
		}
		err = fmt.Errorf(
			"`phoneNumber` does not exist and not associated with any user ",
		)
		serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
			Err:     err,
			Message: err.Error(),
		}, http.StatusBadRequest)
	}
}

// AddRoleToUser is authenticated endpoint that adds role and role based permissions to a user
// whose PRIMARY PHONE NUMBER matches the provided phone number in the request.
func (h *HandlersInterfacesImpl) AddRoleToUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p := &dto.RolePayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, p)
		if p.PhoneNumber == nil {
			err := fmt.Errorf("expected `phoneNumber` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		if p.Role == nil {
			err := fmt.Errorf("expected `roles` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		err := h.interactor.Onboarding.AddRoleToUser(
			ctx,
			*p.PhoneNumber,
			*p.Role,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		serverutils.WriteJSONResponse(
			w,
			dto.OKResp{Status: "OK"},
			http.StatusOK,
		)
	}
}

// RemoveRoleToUser is authenticated endpoint that removes role and role based permissions to a user
// whose PRIMARY PHONE NUMBER matches the provided phone number in the request.
func (h *HandlersInterfacesImpl) RemoveRoleToUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		span := trace.SpanFromContext(ctx)

		p, err := decodePhoneNumberPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		err = h.interactor.Onboarding.RemoveRoleToUser(ctx, *p.PhoneNumber)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		serverutils.WriteJSONResponse(
			w,
			dto.OKResp{Status: "OK"},
			http.StatusOK,
		)
	}
}

// UpdateUserProfile is an unauthenticated REST endpoint to update a user's profile
func (h *HandlersInterfacesImpl) UpdateUserProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.UserProfilePayload{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.UID == nil {
			err := fmt.Errorf("expected `uid` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		authCred := &auth.Token{UID: *payload.UID}
		ctx = context.WithValue(
			ctx,
			firebasetools.AuthTokenContextKey,
			authCred,
		)

		userProfile, err := h.interactor.Signup.UpdateUserProfile(
			ctx,
			&dto.UserProfileInput{
				PhotoUploadID: payload.PhotoUploadID,
				DateOfBirth:   payload.DateOfBirth,
				FirstName:     payload.FirstName,
				LastName:      payload.LastName,
				Gender:        payload.Gender,
			},
		)
		if err != nil {
			errorcodeutil.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, userProfile, http.StatusOK)
	}
}

// SwitchFlaggedFeaturesHandler flips the user as opt-in or opt-out to flagged features
func (h *HandlersInterfacesImpl) SwitchFlaggedFeaturesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		span := trace.SpanFromContext(ctx)
		p, err := decodePhoneNumberPayload(w, r, span)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		okRes, err := h.interactor.Onboarding.SwitchUserFlaggedFeatures(
			ctx,
			*p.PhoneNumber,
		)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}
		serverutils.WriteJSONResponse(w, okRes, http.StatusOK)

	}
}

// PollServices polls registered services for their status
func (h *HandlersInterfacesImpl) PollServices() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		services, err := h.interactor.AdminSrv.PollMicroservicesStatus(ctx)
		if err != nil {
			serverutils.WriteJSONResponse(rw, err, http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(rw, services, http.StatusOK)
	}
}

// CheckHasPermission checks if the user has a permission
func (h *HandlersInterfacesImpl) CheckHasPermission() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		p := &dto.CheckPermissionPayload{}
		serverutils.DecodeJSONToTargetStruct(rw, r, p)

		if p.UID == nil || p.Permission == nil {
			serverutils.WriteJSONResponse(rw, nil, http.StatusBadRequest)
			return
		}

		authorized, err := h.interactor.Role.CheckPermission(ctx, *p.UID, *p.Permission)
		if err != nil {
			serverutils.WriteJSONResponse(rw, err, http.StatusInternalServerError)
			return
		}

		if !authorized {
			serverutils.WriteJSONResponse(rw, nil, http.StatusUnauthorized)
			return
		}

		serverutils.WriteJSONResponse(rw, nil, http.StatusOK)
	}
}

// CreateRole creates a new role given the required role creation input
func (h *HandlersInterfacesImpl) CreateRole() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input := &dto.RoleInput{}
		serverutils.DecodeJSONToTargetStruct(rw, r, input)

		if input.Name == "" {
			serverutils.WriteJSONResponse(rw, nil, http.StatusBadRequest)
			return
		}

		role, err := h.interactor.Role.CreateUnauthorizedRole(ctx, *input)
		if err != nil {
			serverutils.WriteJSONResponse(rw, err, http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(rw, role, http.StatusCreated)
	}
}

// AssignRole assigns a role to the provided user
func (h *HandlersInterfacesImpl) AssignRole() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input := &dto.AssignRolePayload{}
		serverutils.DecodeJSONToTargetStruct(rw, r, input)

		if input.UserID == "" || input.RoleID == "" {
			serverutils.WriteJSONResponse(rw, nil, http.StatusBadRequest)
			return
		}

		_, err := h.interactor.Role.AssignRole(ctx, input.UserID, input.RoleID)
		if err != nil {
			serverutils.WriteJSONResponse(rw, err, http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(rw, nil, http.StatusOK)
	}
}

// RemoveRoleByName removes a role given a valid name
func (h *HandlersInterfacesImpl) RemoveRoleByName() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		input := &dto.DeleteRolePayload{}
		serverutils.DecodeJSONToTargetStruct(rw, r, input)

		if input.Name == "" {
			serverutils.WriteJSONResponse(rw, nil, http.StatusBadRequest)
			return
		}

		role, err := h.interactor.Role.GetRoleByName(ctx, input.Name)
		if err != nil {
			serverutils.WriteJSONResponse(rw, err, http.StatusInternalServerError)
			return
		}

		_, err = h.interactor.Role.UnauthorizedDeleteRole(ctx, role.ID)
		if err != nil {
			serverutils.WriteJSONResponse(rw, err, http.StatusInternalServerError)
			return
		}

		serverutils.WriteJSONResponse(rw, nil, http.StatusOK)
	}
}
