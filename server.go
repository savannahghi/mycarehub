package main

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph"
	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

const serverTimeoutSeconds = 120

var allowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
	"https://api-gateway-test.healthcloud.co.ke",
	"https://api-gateway-prod.healthcloud.co.ke",
	"https://profile-testing-uyajqt434q-ew.a.run.app",
	"https://profile-prod-uyajqt434q-ew.a.run.app",
}
var allowedHeaders = []string{
	"Authorization", "Accept", "Accept-Charset", "Accept-Language",
	"Accept-Encoding", "Origin", "Host", "User-Agent", "Content-Length",
	"Content-Type",
}

func main() {
	ctx := context.Background()

	err := base.Sentry()
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// start up the router
	r, err := Router(ctx)
	if err != nil {
		base.LogStartupError(ctx, err)
	}

	// check if the root colletion env variable exists
	// expects the server to die if this not explictly set
	base.MustGetEnvVar("ROOT_COLLECTION_SUFFIX")

	// start the server
	addr := ":" + base.MustGetEnvVar("PORT")
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
	log.Fatal(srv.ListenAndServe())
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
		http.MethodPost, http.MethodOptions).HandlerFunc(profile.RequestPinResetFunc(ctx, srv))
	r.Path("/update_pin").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(profile.UpdatePinHandler(ctx, srv))

	// Interservice Authenticated routes
	isc := r.PathPrefix("/internal").Subrouter()
	isc.Use(base.InterServiceAuthenticationMiddleware())
	isc.Path("/customer").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(profile.FindCustomerByUIDHandler(ctx, srv))
	isc.Path("/supplier").Methods(
		http.MethodPost, http.MethodOptions).HandlerFunc(profile.FindSupplierByUIDHandler(ctx, srv))

	// check server status.
	r.Path("/health").HandlerFunc(HealthStatusCheck)

	// Authenticated routes
	gqlR := r.Path("/graphql").Subrouter()
	gqlR.Use(base.AuthenticationMiddleware(firebaseApp))
	gqlR.Methods(
		http.MethodPost, http.MethodGet, http.MethodOptions,
	).HandlerFunc(graphqlHandler())
	return r, nil

}

//HealthStatusCheck endpoint to check if the server is working.
func HealthStatusCheck(w http.ResponseWriter, r *http.Request) {

	err := json.NewEncoder(w).Encode(true)
	if err != nil {
		log.Fatal(err)
	}

}

func graphqlHandler() http.HandlerFunc {
	srv := handler.NewDefaultServer(
		generated.NewExecutableSchema(
			generated.Config{
				Resolvers: graph.NewResolver(),
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
