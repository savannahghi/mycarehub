package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
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
				identifiers {
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
				identifiers {
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

func TestCreateFacilities(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation createFacilities($input: [FacilityInput!]!) {
		createFacilities(input: $input) {
		  name
		  phone
		  active
		  country
		  county
		  address
		  description
		  coordinates {
			lat
			lng
		  }
		  identifiers {
			active
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
		// TODO: Restore this once a solution that averts actual creation of facility in health CRM
		// {
		// 	name: "Happy case: create facilities, multiple facilities",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"input": []map[string]interface{}{
		// 					{
		// 						"name":        "One Padmore Facility",
		// 						"phone":       "+254762474022",
		// 						"active":      true,
		// 						"country":     "KE",
		// 						"county":      "Nairobi",
		// 						"address":     "100-Nairobi",
		// 						"description": "One Padmore Place",
		// 						"identifier": map[string]interface{}{
		// 							"type":  "MFL_CODE",
		// 							"value": "12390",
		// 						},
		// 						"coordinates": map[string]interface{}{
		// 							"lat": "48.40338",
		// 							"lng": "15.17403",
		// 						},
		// 					},
		// 					{
		// 						"name":        "One Padmore Facility Two",
		// 						"phone":       "+254767474022",
		// 						"active":      true,
		// 						"country":     "KE",
		// 						"county":      "Nairobi",
		// 						"address":     "100-Nairobi",
		// 						"description": "One Padmore Place Two",
		// 						"identifier": map[string]interface{}{
		// 							"type":  "MFL_CODE",
		// 							"value": "12313",
		// 						},
		// 						"coordinates": map[string]interface{}{
		// 							"lat": "48.40338",
		// 							"lng": "1.17403",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		// {
		// 	name: "Happy case: create facilities, single facility",
		// 	args: args{
		// 		query: map[string]interface{}{
		// 			"query": graphqlMutation,
		// 			"variables": map[string]interface{}{
		// 				"input": []map[string]interface{}{
		// 					{
		// 						"name":        "One Padmore Facility Three",
		// 						"phone":       "+254767474022",
		// 						"active":      true,
		// 						"country":     "KE",
		// 						"county":      "Nairobi",
		// 						"address":     "100-Nairobi",
		// 						"description": "One Padmore Place Three",
		// 						"identifier": map[string]interface{}{
		// 							"type":  "MFL_CODE",
		// 							"value": "12319",
		// 						},
		// 						"coordinates": map[string]interface{}{
		// 							"lat": "48.40338",
		// 							"lng": "1.17403",
		// 						},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// 	wantStatus: http.StatusOK,
		// 	wantErr:    false,
		// },
		{
			name: "Sad case: missing identifier",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": []map[string]interface{}{
							{
								"name":        gofakeit.Name(),
								"phone":       gofakeit.Phone(),
								"active":      true,
								"country":     gofakeit.Country(),
								"description": gofakeit.BS(),
							},
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "Sad case: empty input",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": []map[string]interface{}{},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "Sad case: no input",
			args: args{
				query: map[string]interface{}{
					"query":     graphqlMutation,
					"variables": map[string]interface{}{},
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

func TestGetNearbyFacilities(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	query getNearbyFacilities($locationInput: LocationInput, $paginationInput: PaginationsInput!) {
		getNearbyFacilities(locationInput: $locationInput, paginationInput: $paginationInput){
		  facilities {
			id
			name
			phone
			active
			country
			county
			address
			description
			fhirOrganisationID
			coordinates{
			  lat
			  lng
			}
			identifiers {
			  id
			  active
			  type
			  value
			}
			workStationDetails{
			  notifications
			  surveys
			  messages
			  serviceRequests
			}
			services{
			  id
			  name
			  description
			  identifiers{
				id
				identifierType
				identifierValue
				serviceID
			  }
			}
			businessHours{
			  id
			  day
			  openingTime
			  closingTime
			  facilityID
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
			name: "Happy case: get nearby facility(ies)",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"locationInput": map[string]interface{}{
							"lat":    -1.23456789,
							"lng":    56.4567890,
							"radius": 10.3,
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
			name: "Sad case: unable to get nearby facility(ies)",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"locationInput": map[string]interface{}{
							"lat": -1.23456789,
							"lng": 56.4567890,
						},
						"paginationInput": map[string]interface{}{
							"limit": 1,
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

func TestGetServices(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	query getServices($paginationInput: PaginationsInput!) {
		getServices(paginationInput: $paginationInput) {
		  results {
			id
			name
			description
			identifiers {
			  id
			  identifierType
			  identifierValue
			  serviceID
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
			name: "Happy case: get services",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"paginationInput": map[string]interface{}{
							"limit":       10,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case: unable to get service(s)",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"paginationInput": map[string]interface{}{
							"limit": 10,
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

func Test_BookAService(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	mutation bookService($facilityID: ID!, $serviceIDs: [ID!]!, $time: DateTime!){
		bookService(facilityID: $facilityID, serviceIDs: $serviceIDs, time: $time){
		  id
		  services
		  facility{
			id
		  }
		  client{
			id
		  }
		  organisationID
		  programID
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
			name: "Sad case: unable to book a service",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"serviceIDs": []string{""},
						"time":       "2023-10-29T14:30:00.000Z",
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

func Test_VerifyBookingCode(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	mutation verifyBookingCode($bookingID: ID!, $code: String!, $programID: ID!){
		verifyBookingCode(bookingID: $bookingID, code: $code, programID: $programID)
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
			name: "Happy case: verify booking code",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"bookingID": bookingID,
						"code":      "1234",
						"programID": programID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case: unable to verify booking code",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"bookingID": bookingID,
						"code":      "1234",
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

func Test_ListClientBookings(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	query listBookings($clientID: ID!, $bookingStatus: BookingStatus!, $pagination: PaginationsInput!){
		listBookings(clientID: $clientID, bookingStatus: $bookingStatus, pagination: $pagination){
		  results{
			id
			active
			date
			services{
			  id
			  name
			  description
			  identifiers{
				id
				identifierType
				identifierValue
				serviceID
			  }
			}
			facility{
			  id
			  name
			  phone
			  country
			  county
			  address
			  description
			  coordinates{
				lat
				lng
			  }
			  fhirOrganisationID
			  identifiers{
				id
				active
				type
				value
			  }
			  services{
				id
				name
				description
				identifiers{
				  id
				  identifierType
				  identifierValue
				  serviceID
				}
			  }
			  businessHours{
				id
				day
				openingTime
				closingTime
				facilityID
			  }
			}
			client{
			  id
			  user{
				id
				username
				name
			  }
			}
			organisationID
			programID
			verificationCode
			verificationCodeStatus
			bookingStatus
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
			name: "Happy case: list client bookings",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"clientID":      clientID,
						"bookingStatus": enums.Pending,
						"pagination": map[string]interface{}{
							"limit":       2,
							"currentPage": 1,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad case: unable to list client bookings",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"bookingStatus": enums.Pending,
						"pagination": map[string]interface{}{
							"limit":       2,
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
