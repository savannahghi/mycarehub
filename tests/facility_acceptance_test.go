package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestListProgramFacilities(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	query listProgramFacilities($searchTerm: String, $filterInput: [FiltersInput],$paginationInput :PaginationsInput!) {
		listProgramFacilities (searchTerm: $searchTerm, filterInput: $filterInput, paginationInput: $paginationInput){
			pagination{
				limit
				currentPage
				count
				totalPages
				nextPage
				previousPage
			}
  			facilities {
				id
				name
				phone
				active
				country
				description
				fhirOrganisationID
				identifier {
					id
					active
					type
					value
				}
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
		{
			name: "Happy case: list facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: list facilities by another program",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"programID": programID,
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: filter facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"filterInput": []map[string]interface{}{
							{
								"dataType": enums.FilterSortDataTypeActive,
								"value":    "true",
							},
						},
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: search facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"searchTerm": "Nairobi",
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: search and filter facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"searchTerm": "Nairobi",
						"filterInput": []map[string]interface{}{
							{
								"dataType": enums.FilterSortDataTypeActive,
								"value":    "true",
							},
						},
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case: invalid filter",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"searchTerm": "Kenya",
						"filterInput": []map[string]interface{}{
							{
								"dataType": enums.FilterSortDataType("invalid"),
								"value":    "true",
							},
						},
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
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
					t.Errorf("error not expected, got: %v", data["errors"])
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned, expected: %v, got: %v", resp.StatusCode, tt.wantStatus)
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
		{
			name: "Happy case: inactivate facility",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": inactiveFacilityIdentifier,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case - nil identifier",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"identifier": map[string]interface{}{
						"type":  mflIdentifierType,
						"value": nil,
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
		{
			name: "Happy case: reactivate facility",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"identifier": map[string]interface{}{
							"type":  mflIdentifierType,
							"value": facilityIdentifierToInactivate,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case - nil identifier",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"identifier": map[string]interface{}{
						"type":  mflIdentifierType,
						"value": nil,
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

func TestListFacilities(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	query listFacilities($searchTerm: String, $filterInput: [FiltersInput],$paginationInput :PaginationsInput!) {
		listFacilities (searchTerm: $searchTerm, filterInput: $filterInput, paginationInput: $paginationInput){
			pagination{
				limit
				currentPage
				count
				totalPages
				nextPage
				previousPage
			}
  			facilities {
				id
				name
				phone
				active
				country
				description
				fhirOrganisationID
				identifier {
					id
					active
					type
					value
				}
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
		{
			name: "Happy case: list facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: filter facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"filterInput": []map[string]interface{}{
							{
								"dataType": enums.FilterSortDataTypeActive,
								"value":    "true",
							},
						},
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: search facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"searchTerm": "Nairobi",
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Happy case: search and filter facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"searchTerm": "Nairobi",
						"filterInput": []map[string]interface{}{
							{
								"dataType": enums.FilterSortDataTypeActive,
								"value":    "true",
							},
						},
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case: invalid filter",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"searchTerm": "Kenya",
						"filterInput": []map[string]interface{}{
							{
								"dataType": enums.FilterSortDataType("invalid"),
								"value":    "true",
							},
						},
						"paginationInput": map[string]interface{}{
							"limit":       1,
							"currentPage": 1,
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
