package presentation_test

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

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation"

	"gitlab.slade360emr.com/go/base"
)

const (
	testHTTPClientTimeout = 180
)

/// these are set up once in TestMain and used by all the acceptance tests in
// this package
var srv *http.Server
var baseURL string
var serverErr error

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestMain(m *testing.M) {
	// setup
	os.Setenv("ENVIRONMENT", "staging")
	os.Setenv("ROOT_COLLECTION_SUFFIX", "onboarding_testing")

	ctx := context.Background()
	srv, baseURL, serverErr = base.StartTestServer(
		ctx,
		presentation.PrepareServer,
		presentation.AllowedOrigins,
	) // set the globals
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	// run the tests
	log.Printf("about to run tests")
	code := m.Run()
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

func TestRouter(t *testing.T) {
	ctx := context.Background()
	router, err := presentation.Router(ctx)
	if err != nil {
		t.Errorf("can't initialize router: %v", err)
		return
	}

	if router == nil {
		t.Errorf("nil router")
		return
	}
}

func TestHealthStatusCheck(t *testing.T) {
	client := http.DefaultClient

	type args struct {
		url        string
		httpMethod string
		body       io.Reader
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "successful health check",
			args: args{
				url: fmt.Sprintf(
					"%s/health",
					baseURL,
				),
				httpMethod: http.MethodPost,
				body:       nil,
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := http.NewRequest(
				tt.args.httpMethod,
				tt.args.url,
				tt.args.body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range base.GetDefaultHeaders(t, baseURL, "profile") {
				r.Header.Add(k, v)
			}

			resp, err := client.Do(r)
			if err != nil {
				t.Errorf("HTTP error: %v", err)
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}

			if tt.wantErr {
				return
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("can't read response body: %v", err)
				return
			}

			if data == nil {
				t.Errorf("nil response body data")
				return
			}

			if tt.wantStatus != resp.StatusCode {
				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
				return
			}

			if !tt.wantErr && resp == nil {
				t.Errorf("unexpected nil response (did not expect an error)")
				return
			}
		})
	}
}

// ComposeGraphqlServerRequest creates a graphql server client
// func ComposeGraphqlServerRequest(ctx context.Context, query map[string]interface{}) ([]byte, error) {

// 	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
// 	headers, err := base.GetGraphQLHeaders(ctx)
// 	if err != nil {
// 		return []byte{}, fmt.Errorf("error in getting headers: %w", err)
// 	}

// 	body, err := mapToJSONReader(query)

// 	if err != nil {
// 		return []byte{}, fmt.Errorf("unable to get GQL JSON io Reader: %s", err)
// 	}

// 	r, err := http.NewRequest(
// 		http.MethodPost,
// 		graphQLURL,
// 		body,
// 	)

// 	if err != nil {
// 		return []byte{}, fmt.Errorf("unable to compose request: %s", err)
// 	}

// 	if r == nil {
// 		return []byte{}, fmt.Errorf("nil request")
// 	}

// 	for k, v := range headers {
// 		r.Header.Add(k, v)
// 	}

// 	client := http.Client{
// 		Timeout: time.Second * testHTTPClientTimeout,
// 	}
// 	resp, err := client.Do(r)
// 	if err != nil {
// 		return []byte{}, fmt.Errorf("request error: %s", err)
// 	}
// 	dataResponse, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return []byte{}, fmt.Errorf("can't read request body: %s", err)
// 	}

// 	return dataResponse, nil
// }
