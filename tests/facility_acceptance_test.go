package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestCreateFacility(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation createFacility($facility: FacilityInput!, $identifier:FacilityIdentifierInput!) {
		createFacility (facility: $facility, identifier: $identifier) {
		  facility {
			name
			phone
			active
			county
			description
		  }
		  identifier {
			type
			value
		  }
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
		// {
		// 	name: "success: create a facility with valid payload",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"facility": map[string]interface{}{
		// 					"name":        facilityName,
		// 					"phone":       phone,
		// 					"active":      true,
		// 					"county":      county,
		// 					"description": description,
		// 				},
		// 				"identifier": map[string]interface{}{
		// 					"type":  mflIdentifierType,
		// 					"value": "893298329",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		{
			name: "invalid: missing name param",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facility": map[string]interface{}{
							"active":      true,
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
						},
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": "4343445",
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: missing code param",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facility": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"active":      true,
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
						},
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": "545345343",
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: missing active param",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facility": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
						},
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": "566498082232",
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: missing county param",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facility": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"active":      true,
							"description": "located at Giddo plaza building town",
						},
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": "988967822434643",
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: missing description param",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facility": map[string]interface{}{
							"name":   "Mediheal Hospital (Nakuru) Annex",
							"active": true,
							"county": "Nakuru",
						},
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": "65645487878",
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: invalid value for active",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facility": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"active":      "invalid",
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
						},
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": "454545454",
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
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestInactivateFacility(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation inactivateFacility($identifier: FacilityIdentifierInput!) {
		inactivateFacility (identifier: $identifier)
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
		// {
		// 	name: "Happy case",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"identifier": map[string]interface{}{
		// 					"type":  mflIdentifierType,
		// 					"value": inactiveFacilityIdentifier,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		{
			name: "Sad case - nil MFL Code",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"identifier": nil,
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
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}
		})
	}
}

func TestReactivateFacility(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation reactivateFacility($identifier: FacilityIdentifierInput!) {
		reactivateFacility (identifier: $identifier)
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
		// {
		// 	name: "Happy case",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"identifier": map[string]interface{}{
		// 					"type":  mflIdentifierType,
		// 					"value": facilityIdentifierToInactivate,
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		{
			name: "Sad case - nil MFL Code",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"mflCode": nil,
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
					t.Errorf("error not expected")
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned")
				return
			}

		})
	}
}
