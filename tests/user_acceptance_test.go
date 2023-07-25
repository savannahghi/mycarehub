package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func TestDeleteClientProfile(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	mutation deleteClientProfile($clientID: String!){
		deleteClientProfile(clientID: $clientID)
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
			name: "success: opt out as client",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID": testOPtOutClient,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: opt out as client with caregiver profile",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID": testOPtOutClientCaregiver,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "success: opt out as client with staff profile",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID": testOPtOutClientStaff,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "success: opt out as client with staff profile 2",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID": testOptOutStaffClient,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "success: opt out as client with 2 client profiles",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID": testOPtOutTwoClient,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "invalid: client id is not a valid",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID": "invalid",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid: client id is not passed",
			args: args{
				query: map[string]interface{}{
					"query":     graphqlQuery,
					"variables": map[string]interface{}{},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "success: opt out as client" {
				regPayload := &domain.MatrixUserRegistration{
					Username: "testoptoutclient",
					Password: testOPtOutClient,
					Admin:    false,
				}

				err := registerMatrixUser(ctx, regPayload)
				if err != nil {
					fmt.Print("the error is %w: ", err)
				}
			}

			if tt.name == "success: opt out as client with caregiver profile" {
				regPayload := &domain.MatrixUserRegistration{
					Username: "testoptoutclientcaregiver",
					Password: testOPtOutClientCaregiver,
					Admin:    false,
				}

				err := registerMatrixUser(ctx, regPayload)
				if err != nil {
					fmt.Print("the error is %w: ", err)
				}
			}

			if tt.name == "success: opt out as client with staff profile" {
				regPayload := &domain.MatrixUserRegistration{
					Username: "testoptoutclientstaff",
					Password: testOPtOutClientStaff,
					Admin:    true,
				}

				err := registerMatrixUser(ctx, regPayload)
				if err != nil {
					fmt.Print("the error is %w: ", err)
				}
			}

			if tt.name == "success: opt out as client with staff profile 2" {
				regPayload := &domain.MatrixUserRegistration{
					Username: "testoptoutstaffclient",
					Password: testOptOutStaffClient,
					Admin:    true,
				}

				err := registerMatrixUser(ctx, regPayload)
				if err != nil {
					fmt.Print("the error is %w: ", err)
				}
			}

			if tt.name == "success: opt out as client with 2 client profiles" {
				regPayload := &domain.MatrixUserRegistration{
					Username: "testoptouttwoclient",
					Password: testOPtOutTwoClient,
					Admin:    false,
				}

				err := registerMatrixUser(ctx, regPayload)
				if err != nil {
					fmt.Print("the error is %w: ", err)
				}
			}

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get GQL JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			for k, v := range headers {
				r.Header.Add(k, v)
			}
			client := http.Client{
				Timeout: time.Second * testHTTPClientTimeout,
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

			data := map[string]interface{}{}
			err = json.Unmarshal(dataResponse, &data)
			if err != nil {
				t.Errorf("bad data returned")
				return
			}

			if tt.wantErr {
				errMsg, ok := data["errors"]
				if !ok {
					t.Errorf("GraphQL error: %s", errMsg)
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected, got %v", data["errors"])
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned, expected %v, got %v", tt.wantStatus, resp.StatusCode)
				return
			}
		})
	}
}
