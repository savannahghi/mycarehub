package main

import (
	"compress/gzip"
	"context"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"gitlab.slade360emr.com/go/accounting/graph"
	"gitlab.slade360emr.com/go/accounting/graph/generated"
	"gitlab.slade360emr.com/go/base"
)

const serverTimeoutSeconds = 120

var allowedOrigins = []string{
	"https://healthcloud.co.ke",
	"https://bewell.healthcloud.co.ke",
	"http://localhost:5000",
	"https://api-gateway-test.healthcloud.co.ke",
	"https://api-gateway-prod.healthcloud.co.ke",
	"https://accounting-lermby4k3q-ew.a.run.app",
}
var allowedHeaders = []string{
	"accounting", "Accept", "Accept-Charset", "Accept-Language",
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
	r, err := Router()
	if err != nil {
		base.LogStartupError(ctx, err)
	}

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
func Router() (*mux.Router, error) {
	fc := &base.FirebaseClient{}
	firebaseApp, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}
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
				Resolvers: graph.NewResolver(),
			},
		),
	)
	return func(w http.ResponseWriter, r *http.Request) {
		srv.ServeHTTP(w, r)
	}
}
