package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/imroc/req"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/savannahghi/serverutils"
)

const (
	testHTTPClientTimeout = 180
)

var (
	srv           *http.Server
	baseURL       string
	serverErr     error
	matrixBaseURL = serverutils.MustGetEnvVar("MATRIX_BASE_URL")
	token         string
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

	originalIntrospectURL := os.Getenv("MYCAREHUB_INTROSPECT_URL")
	originalTokenURL := os.Getenv("MYCAREHUB_TOKEN_URL")
	os.Setenv("MYCAREHUB_INTROSPECT_URL", fmt.Sprintf("%s/oauth/introspect", baseURL))
	os.Setenv("MYCAREHUB_TOKEN_URL", fmt.Sprintf("%s/oauth/token", baseURL))

	regPayload := &domain.MatrixUserRegistration{
		Username: "thisisatestuser",
		Password: userID,
		Admin:    true,
	}
	err := registerMatrixUser(ctx, regPayload)
	if err != nil {
		fmt.Print("the error is %w: ", err)
	}

	tk, err := GetBearerTokenHeader(ctx)
	if err != nil {
		log.Printf("unable to start test server: %s", err)
	}

	token = tk
	// run tests
	log.Printf("Running tests ...")
	code := m.Run()

	// restore envs
	os.Setenv("ENVIRONMENT", initialEnv)
	os.Setenv("MYCAREHUB_INTROSPECT_URL", originalIntrospectURL)
	os.Setenv("MYCAREHUB_TOKEN_URL", originalTokenURL)

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

// CommunityUserRegistration defines the structure of the input to be used when registering a Matrix user
type UserRdegistration struct {
	Auth     *Auth  `json:"auth"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Identifier represents the matrix identifier to be used while logging in
type Identifier struct {
	Type string `json:"type"`
	User string `json:"user"`
}

// Auth is defines the type of authentication to be used when registering a new user
type Auth struct {
	Type string `json:"type"`
}

// RequestHelperPayload is the payload that is used to make requests to matrix client
type RequestHelperPayload struct {
	Method string
	Path   string
	Body   interface{}
}

func loginMatrixUser(ctx context.Context, username string, password string) (string, error) {
	client := http.Client{}

	loginPayload := struct {
		Identifier *Identifier `json:"identifier"`
		Type       string      `json:"type"`
		Password   string      `json:"password"`
	}{
		Identifier: &Identifier{
			Type: "m.id.user",
			User: username,
		},
		Type:     "m.login.password",
		Password: password,
	}

	matrixLoginURL := fmt.Sprintf("%s/_matrix/client/v3/login", matrixBaseURL)
	payload := RequestHelperPayload{
		Method: http.MethodPost,
		Path:   matrixLoginURL,
		Body:   loginPayload,
	}

	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return "", err
	}

	p := bytes.NewBuffer(encoded)
	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	data := struct {
		AccessToken string `json:"access_token"`
	}{}
	if err := json.Unmarshal(respBytes, &data); err != nil {
		return "", err
	}

	return data.AccessToken, nil
}

func registerMatrixUser(ctx context.Context, registrationPayload *domain.MatrixUserRegistration) error {
	tkn, err := loginMatrixUser(ctx, os.Getenv("MCH_MATRIX_USER"), os.Getenv("MCH_MATRIX_PASSWORD"))
	if err != nil {
		fmt.Print("the error is: %w", err)
	}

	client := http.Client{}
	matrixUser := &domain.MatrixUserRegistration{
		Username: registrationPayload.Username,
		Password: registrationPayload.Password,
		Admin:    registrationPayload.Admin,
	}

	matrixRegisterURL := fmt.Sprintf("%s/_synapse/admin/v2/users/@%s:prohealth360.org", matrixBaseURL, matrixUser.Username)
	payload := RequestHelperPayload{
		Method: http.MethodPut,
		Path:   matrixRegisterURL,
		Body:   matrixUser,
	}

	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return err
	}

	p := bytes.NewBuffer(encoded)
	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tkn)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode > http.StatusAccepted {
		return fmt.Errorf("matrix registration failed with status code %d", resp.StatusCode)
	}

	return nil
}

// GetGraphQLHeaders gets relevant GraphQLHeaders
func GetGraphQLHeaders(ctx context.Context) (map[string]string, error) {
	return req.Header{
		"Accept":        "application/json",
		"Content-Type":  "application/json",
		"Authorization": token,
	}, nil
}

// GetBearerTokenHeader gets bearer Token Header
func GetBearerTokenHeader(ctx context.Context) (string, error) {
	input := dto.LoginInput{
		Username: "thisisatestuser",
		PIN:      "0000",
		Flavour:  feedlib.FlavourConsumer,
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/login_by_phone", baseURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	var response = dto.LoginResponse{}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(bs, &response); err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", response.Response.AuthCredentials.IDToken), nil
}

// GetBearerTokenHeader gets bearer Token Header
func GetFirebaseBearerTokenHeader(ctx context.Context) (string, error) {
	customToken, err := firebasetools.CreateFirebaseCustomTokenWithClaims(ctx, userID, nil)
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
