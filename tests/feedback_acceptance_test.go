package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func Test_SendFeedback(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLMutation := `
	mutation sendFeedback($input: FeedbackResponseInput!){
		sendFeedback(input: $input)
	  }
	`

	type args struct {
		query map[string]interface{}
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "send feedback",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"userID":            userID,
							"feedbackType":      "GENERAL_FEEDBACK",
							"satisfactionLevel": 4,
							"serviceName":       "Test",
							"feedback":          "test",
							"requiresFollowUp":  true,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "FailL unable to send feedback; no user ID",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"feedbackType":      "GENERAL_FEEDBACK",
							"satisfactionLevel": 4,
							"serviceName":       "Test",
							"feedback":          "test",
							"requiresFollowUp":  true,
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to marshal query: %s", err)
				return
			}

			req, err := http.NewRequest("POST", graphQLURL, body)
			if err != nil {
				t.Errorf("unable to create request: %s", err)
				return
			}
			if req == nil {
				t.Errorf("request is nil")
				return
			}

			for k, v := range headers {
				req.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("unable to make request: %s", err)
				return
			}

			dataResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("unable to read response body: %s", err)
				return
			}
			if dataResp == nil {
				t.Errorf("response body is nil")
				return
			}

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResp, &data)
			if err != nil {
				t.Errorf("unable to unmarshal response body: %s", err)
				return
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("unexpected error: %s", data["errors"])
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("unexpected status code: %d", resp.StatusCode)
				return
			}
		})
	}

}
