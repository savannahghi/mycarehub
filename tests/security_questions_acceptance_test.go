package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/savannahghi/feedlib"
)

func Test_GetSecurityQuestions(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLMutation := `
		query getSecurityQuestions($flavour: Flavour!){
			getSecurityQuestions(flavour: $flavour){
			securityQuestionID
			questionStem
			description
			active
			responseType
			}
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
			name: "success: get security questions",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"flavour": feedlib.FlavourConsumer,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: invalid flavour",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"flavour": feedlib.Flavour("invalid"),
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

			dataResp, err := io.ReadAll(resp.Body)
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

func Test_RecordSecurityQuestionResponses(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLMutation := `
		mutation recordSecurityQuestionResponses($input: [SecurityQuestionResponseInput!]!){
			recordSecurityQuestionResponses(input: $input){
			securityQuestionID
			isCorrect
			}
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
			name: "success: record security questions",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": []map[string]interface{}{
							{
								"userID":             userID,
								"securityQuestionID": securityQuestionID,
								"response":           "1917",
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: invalid input",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": "invalid",
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

			dataResp, err := io.ReadAll(resp.Body)
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
