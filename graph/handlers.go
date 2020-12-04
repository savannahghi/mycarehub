package graph

import (
	"compress/gzip"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"firebase.google.com/go/auth"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

const serverTimeoutSeconds = 120

var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type",
}

// Router sets up the ginContext router
func Router(ctx context.Context) (*mux.Router, error) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	srv := profile.NewService()

	r := mux.NewRouter() // gorilla mux
	r.Use(
		handlers.RecoveryHandler(
			handlers.PrintRecoveryStack(true),
			handlers.RecoveryLogger(log.StandardLogger()),
		),
	) // recover from panics by writing a HTTP error
	r.Use(base.RequestDebugMiddleware())

	// Unauthenticated routes
	r.Path("/ide").HandlerFunc(playground.Handler("GraphQL IDE", "/graphql"))
	r.Path("/msisdn_login").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(base.GetPhoneNumberLoginFunc(ctx, fc))
	r.Path("/request_pin_reset").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(RequestPINResetFunc(ctx))
	r.Path("/update_pin").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(UpdatePinHandler(ctx))
	r.Path("/send_retry_otp").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(SendRetryOTPHandler(ctx))

	// check server status.
	r.Path("/health").HandlerFunc(base.HealthStatusCheck)

	// Interservice Authenticated routes
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(base.InterServiceAuthenticationMiddleware())
	isc.Path("/customer").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(profile.FindCustomerByUIDHandler(ctx, srv))

	isc.Path("/supplier").Methods(
		http.MethodPost, http.MethodOptions,
	).HandlerFunc(profile.FindSupplierByUIDHandler(ctx, srv))

	isc.Path("/contactdetails/{attribute}/").Methods(
		http.MethodPost).HandlerFunc(
		GetProfileAttributesHandler(ctx),
	).Name("getProfileAttributes")

	isc.Path("/retrieve_user_profile_firebase_doc").Methods(
		http.MethodPost).HandlerFunc(RetrieveUserProfileFirebaseDocSnapshotHandler(ctx, srv))

	isc.Path("/save_cover").Methods(http.MethodPost).HandlerFunc(SaveMemberCoverToFirestoreHandler(ctx, srv))

	isc.Path("/is_underage").Methods(http.MethodPost).HandlerFunc(IsUnderAgeHandler(ctx, srv))
	isc.Path("/user_profile").Methods(http.MethodPost).HandlerFunc(UserProfileHandler(ctx, srv))

	// Authenticated routes
	gqlR := r.Path("/graphql").Subrouter()
	gqlR.Use(base.AuthenticationMiddleware(firebaseApp))
	gqlR.Methods(
		http.MethodPost, http.MethodGet, http.MethodOptions,
	).HandlerFunc(graphqlHandler())
	return r, nil

}

func graphqlHandler() http.HandlerFunc {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: NewResolver(),
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}

// PrepareServer starts up a server
func PrepareServer(ctx context.Context, port int, allowedOrigins []string) *http.Server {
	// start up the router
	r, err := Router(ctx)
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// start the server
	addr := fmt.Sprintf(":%d", port)
	h := handlers.CompressHandlerLevel(r, gzip.BestCompression)
	h = handlers.CORS(
		handlers.AllowedHeaders(allowedHeaders),
		handlers.AllowedOrigins(allowedOrigins),
		handlers.AllowCredentials(),
		handlers.AllowedMethods([]string{"OPTIONS", "GET", "POST"}),
	)(h)
	h = handlers.CombinedLoggingHandler(os.Stdout, h)
	h = handlers.ContentTypeHandler(h, "application/json")
	srv := &http.Server{
		Handler:      h,
		Addr:         addr,
		WriteTimeout: serverTimeoutSeconds * time.Second,
		ReadTimeout:  serverTimeoutSeconds * time.Second,
	}
	log.Infof("Server running at port %v", addr)
	return srv
}

// GetProfileAttributesHandler retreives confirmed user profile attributes
func GetProfileAttributesHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		validUids, err := profile.ValidateUserProfileUIDs(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		uids := validUids.UIDs
		response, err := profile.GetAttribute(ctx, r, uids)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		base.WriteJSONResponse(w, response, http.StatusOK)

	}
}

// RequestPINResetFunc returns a function that sends an otp to an msisdn that requests a
// PIN reset request during login
func RequestPINResetFunc(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := profile.NewService()
		validMsisdn, err := profile.ValidateMsisdn(w, r)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		msisdn := validMsisdn.MSISDN
		otp, err := s.RequestPINReset(ctx, msisdn)
		if err != nil {
			base.ReportErr(w, err, http.StatusBadRequest)
			return
		}

		otpResponse := profile.OtpResponse{
			OTP: otp,
		}

		base.WriteJSONResponse(w, otpResponse, http.StatusOK)
	}
}

// UpdatePinHandler used to update a user's PIN
func UpdatePinHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := profile.NewService()
		payload, validateErr := profile.ValidateUpdatePinPayload(w, r)
		if validateErr != nil {
			base.ReportErr(w, validateErr, http.StatusBadRequest)
			return
		}

		_, updateErr := s.UpdateUserPIN(ctx, payload.MSISDN, payload.PINNumber, payload.OTP)
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

// RetrieveUserProfileFirebaseDocSnapshotHandler process requests for ISC to RetrieveUserProfileFirebaseDocSnapshot
func RetrieveUserProfileFirebaseDocSnapshotHandler(ctx context.Context, srv *profile.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		profileUID := &profile.BusinessPartnerUID{}
		base.DecodeJSONToTargetStruct(rw, r, profileUID)
		if profileUID == nil {
			err := fmt.Errorf("invalid credetials")
			base.RespondWithError(rw, http.StatusBadRequest, err)
			return
		}

		if profileUID.Token == nil {
			base.RespondWithError(rw, http.StatusBadRequest, fmt.Errorf("bad token provided"))
			return
		}

		newContext := context.WithValue(ctx, base.AuthTokenContextKey, profileUID.Token)
		doc, err := srv.RetrieveUserProfileFirebaseDocSnapshot(newContext)

		if err != nil {
			base.RespondWithError(rw, http.StatusBadRequest, err)
			return
		}

		base.WriteJSONResponse(rw, doc, http.StatusOK)
	}
}

// SaveMemberCoverToFirestoreHandler process ISC requests to SaveMemberCoverToFirestore
func SaveMemberCoverToFirestoreHandler(ctx context.Context, srv *profile.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		type Payload struct {
			PayerName      string      `json:"payerName"`
			MemberName     string      `json:"memberName"`
			MemberNumber   string      `json:"memberNumber"`
			PayerSladeCode int         `json:"payerSladeCode"`
			UUID           string      `json:"uid"`
			Token          *auth.Token `json:"token"`
		}

		payload := &Payload{}

		base.DecodeJSONToTargetStruct(rw, r, payload)

		if payload.Token == nil {
			base.RespondWithError(rw, http.StatusBadRequest, fmt.Errorf("bad token provided"))
			return
		}

		newContext := context.WithValue(ctx, base.AuthTokenContextKey, payload.Token)

		err := srv.SaveMemberCoverToFirestore(newContext, payload.PayerName, payload.MemberNumber, payload.MemberName, payload.PayerSladeCode)

		if err != nil {
			base.RespondWithError(rw, http.StatusBadRequest, fmt.Errorf("failed to save cover"))
			return
		}

		type ResponsePayload struct {
			SuccessfulySaved bool `json:"successfullySaved"`
		}

		base.WriteJSONResponse(rw, ResponsePayload{SuccessfulySaved: true}, http.StatusOK)
	}
}

// IsUnderAgeHandler process ISC requests to IsUnderAge
func IsUnderAgeHandler(ctx context.Context, srv *profile.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		type UserContext struct {
			Token *auth.Token `json:"token"`
		}
		var userContext *UserContext
		base.DecodeJSONToTargetStruct(rw, r, &userContext)
		if userContext == nil {
			err := fmt.Errorf("invalid credetials")
			base.RespondWithError(rw, http.StatusBadRequest, err)
			return
		}

		if userContext.Token == nil {
			base.RespondWithError(rw, http.StatusBadRequest, fmt.Errorf("bad token provided"))
			return
		}

		newContext := context.WithValue(ctx, base.AuthTokenContextKey, userContext.Token)
		isUnderAge, err := srv.IsUnderAge(newContext)

		if err != nil {
			base.RespondWithError(rw, http.StatusInternalServerError, err)
		}

		type Payload struct {
			IsUnderAge bool `json:"isUnderAge"`
		}

		payload := Payload{
			IsUnderAge: isUnderAge,
		}

		base.WriteJSONResponse(rw, payload, http.StatusOK)
	}
}

// UserProfileHandler process ISC request to UserProfile
func UserProfileHandler(ctx context.Context, srv *profile.Service) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		type userContext struct {
			Token *auth.Token `json:"token"`
		}

		profileContext := &userContext{}
		base.DecodeJSONToTargetStruct(rw, r, profileContext)
		if profileContext == nil {
			err := fmt.Errorf("invalid credetials")
			base.RespondWithError(rw, http.StatusBadRequest, err)
			return
		}

		if profileContext.Token == nil {
			base.RespondWithError(rw, http.StatusBadRequest, fmt.Errorf("bad token provided"))
			return
		}

		newContext := context.WithValue(ctx, base.AuthTokenContextKey, profileContext.Token)

		userProfile, err := srv.UserProfile(newContext)
		if err != nil {
			base.RespondWithError(rw, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSONResponse(rw, userProfile, http.StatusOK)

	}
}

// SendRetryOTPHandler generates fallback OTPs when Africa is talking sms fails
func SendRetryOTPHandler(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s := profile.NewService()
		payload, validateErr := profile.ValidateSendRetryOTPPayload(w, r)
		if validateErr != nil {
			base.ReportErr(w, validateErr, http.StatusBadRequest)
			return
		}

		_, updateErr := s.SendRetryOTP(ctx, payload.Msisdn, payload.RetryStep)
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
