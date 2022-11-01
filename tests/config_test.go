package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/imroc/req"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/savannahghi/serverutils"
)

const (
	testHTTPClientTimeout = 180
)

var (
	srv       *http.Server
	baseURL   string
	serverErr error
)

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestMain(m *testing.M) {
	log.Printf("Setting tests up ...")

	initialEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")

	setupFixtures()

	ctx := context.Background()

	srv, baseURL, serverErr = serverutils.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.AllowedOrigins,
	)
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	// run tests
	log.Printf("Running tests ...")
	code := m.Run()

	// restore envs
	os.Setenv("ENVIRONMENT", initialEnv)

	log.Printf("finished running tests")

	// cleanup here
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()

	os.Exit(code)
}

// GetGraphQLHeaders gets relevant GraphQLHeaders
func GetGraphQLHeaders(ctx context.Context) (map[string]string, error) {
	authorization, err := GetBearerTokenHeader(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't Generate Bearer Token: %s", err)
	}
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": authorization,
	}, nil
}

// GetBearerTokenHeader gets bearer Token Header
func GetBearerTokenHeader(ctx context.Context) (string, error) {
	claims := map[string]interface{}{
		"organisationID": orgID,
	}
	customToken, err := firebasetools.CreateFirebaseCustomTokenWithClaims(ctx, userID, claims)
	if err != nil {
		return "", fmt.Errorf("can't create custom token: %s", err)
	}

	if customToken == "" {
		return "", fmt.Errorf("blank custom token: %s", err)
	}

	idTokens, err := firebasetools.AuthenticateCustomFirebaseToken(customToken)
	if err != nil {
		return "", fmt.Errorf("can't authenticate custom token: %s", err)
	}
	if idTokens == nil {
		return "", fmt.Errorf("nil idTokens")
	}

	return fmt.Sprintf("Bearer %s", idTokens.IDToken), nil
}
