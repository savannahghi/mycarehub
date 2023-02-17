package pubsubmessaging_test

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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation"
	"github.com/savannahghi/pubsubtools"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"google.golang.org/api/idtoken"
)

var (
	srv       *http.Server
	baseURL   string
	serverErr error
)

func TestMain(m *testing.M) {
	initialEnv := os.Getenv("ENVIRONMENT")
	os.Setenv("ENVIRONMENT", "staging")

	ctx := context.Background()
	srv, baseURL, serverErr = serverutils.StartTestServer(ctx, presentation.PrepareServer, presentation.AllowedOrigins)
	if serverErr != nil {
		log.Printf("unable to start test server: %s", serverErr)
	}

	code := m.Run()

	// restore envs
	os.Setenv("ENVIRONMENT", initialEnv)
	defer func() {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.Printf("test server shutdown error: %s", err)
		}
	}()

	os.Exit(code)
}

func composeInvalidPubsubTestPayload(t *testing.T, topic string) (*bytes.Buffer, error) {
	// Compose the payload
	pubsubPayload := pubsubtools.PubSubPayload{
		Message: pubsubtools.PubSubMessage{
			Data: nil,
			Attributes: map[string]string{
				"invalid": "invalid",
			},
		},
		Subscription: ksuid.New().String(),
	}

	payload, err := json.Marshal(pubsubPayload)
	if err != nil {
		return nil, err
	}
	bs := bytes.NewBuffer(payload)

	return bs, nil
}

func TestPubsub(t *testing.T) {
	ctx := context.Background()

	invalidPayload, err := composeInvalidPubsubTestPayload(t, "invalidTopic")
	if err != nil {
		t.Errorf("failed to compose invalid payload")
		return
	}

	header := req.Header{
		"Content-Type": "application/json",
	}
	type args struct {
		url        string
		httpMethod string
		body       io.Reader
		headers    map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Sad Case - Invalid payload",
			args: args{
				url:        fmt.Sprintf("%v/pubsub", baseURL),
				httpMethod: http.MethodPost,
				body:       invalidPayload,
				headers:    header,
			},
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
	}
	t.Parallel()
	for _, tt := range tests {
		r, err := http.NewRequest(
			tt.args.httpMethod,
			tt.args.url,
			tt.args.body,
		)
		if err != nil {
			t.Errorf("unable to compose request: %s", err)
			return
		}

		if r == nil {
			t.Errorf("nil request")
			return
		}

		for k, v := range tt.args.headers {
			r.Header.Add(k, v)
		}

		client, err := idtoken.NewClient(ctx, pubsubtools.Aud)
		if err != nil {
			t.Errorf("can't initialize client: %s", err)
			return
		}
		resp, err := client.Do(r)
		if err != nil {
			t.Errorf("request error: %s", err)
			return
		}

		dataResponse, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("can't read request body: %s", err)
			return
		}
		if dataResponse == nil {
			t.Errorf("nil response data")
			return
		}

		if tt.wantStatus != resp.StatusCode {
			t.Errorf(
				"expected status %d, got %d and response %s",
				tt.wantStatus,
				resp.StatusCode,
				string(dataResponse),
			)
			return
		}
	}
}
