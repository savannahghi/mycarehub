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
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func Test_CreateServiceRequest(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLMutation := `
		mutation createServiceRequest($input: ServiceRequestInput!) {
			createServiceRequest(input: $input)
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
			name: "success: create service request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"active":       true,
							"requestType":  enums.ServiceRequestTypeRedFlag,
							"status":       enums.ServiceRequestStatusPending,
							"request":      "TEST",
							"clientID":     clientID,
							"inProgressBy": staffID,
							"resolvedBy":   staffID,
							"facilityID":   facilityID,
							"clientName":   gofakeit.BeerName(),
							"flavour":      feedlib.FlavourConsumer,
							"meta": map[string]interface{}{
								"test": "test",
							},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: create service request - by a caregiver",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"active":       true,
							"requestType":  enums.ServiceRequestTypeRedFlag,
							"status":       enums.ServiceRequestStatusPending,
							"request":      "TEST",
							"clientID":     clientID,
							"inProgressBy": staffID,
							"resolvedBy":   staffID,
							"facilityID":   facilityID,
							"clientName":   gofakeit.BeerName(),
							"flavour":      feedlib.FlavourConsumer,
							"meta": map[string]interface{}{
								"test": "test",
							},
							"caregiverID": testCaregiverID,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad: unable to create service request without client ID",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"active":       true,
							"requestType":  enums.ServiceRequestTypeRedFlag,
							"status":       enums.ServiceRequestStatusPending,
							"request":      "TEST",
							"inProgressBy": staffID,
							"resolvedBy":   staffID,
							"facilityID":   facilityID,
							"clientName":   gofakeit.BeerName(),
							"flavour":      feedlib.FlavourConsumer,
							"meta": map[string]interface{}{
								"test": "test",
							},
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

func Test_ResolveServiceRequest(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLMutation := `
	mutation resolveServiceRequest($staffID: String!, $requestID: String!, $action: [String!]!, $comment: String) {
		resolveServiceRequest(staffID: $staffID, requestID: $requestID, action: $action, comment: $comment)
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
			name: "success: resolve service request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"staffID":   staffID,
						"requestID": serviceRequestID,
						"action":    []string{"resolve"},
						"comment":   "resolved",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad: unable to resolve service request with invalid staff ID",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"requestID": serviceRequestID,
						"action":    []string{"resolve"},
						"comment":   "resolved",
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

func Test_GetServiceRequests(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLQuery := `
	query getServiceRequests($type: String, $status: String, $facilityID: String!, $flavour: Flavour!, $pagination: PaginationsInput!) {
		getServiceRequests(
		  requestType: $type
		  requestStatus: $status
		  facilityID: $facilityID
		  flavour: $flavour
		  pagination: $pagination
		) {
		  results {
			id
			requestType
			request
			clientID
			facilityID
			clientName
			clientContact
			staffContact
			staffID
			staffName
			username
			status
			createdAt
			meta
			services {
			  id
			  name
			  description
			}
		  }
		  pagination {
			limit
			currentPage
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
			name: "success: get service request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQuery,
					"variables": map[string]interface{}{
						"type":       "RED_FLAG",
						"status":     enums.ServiceRequestStatusResolved,
						"facilityID": facilityID,
						"flavour":    feedlib.FlavourConsumer,
						"pagination": map[string]interface{}{
							"currentPage": 1,
							"limit":       5,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: unable to get service request; no flavour defined",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQuery,
					"variables": map[string]interface{}{
						"requestType":   "RED_FLAG",
						"requestStatus": enums.ServiceRequestStatusResolved,
						"facilityID":    facilityID,
						"pagination": map[string]interface{}{
							"currentPage": 1,
							"limit":       5,
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

func Test_GetPendingServiceRequestsCount(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLQuery := `
	query getPendingServiceRequestsCount{
		getPendingServiceRequestsCount{
		  clientsServiceRequestCount{
			requestsTypeCount{
			  requestType
			  total
			}
		  }
		  staffServiceRequestCount{
			requestsTypeCount{
			  requestType
			  total
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
			name: "success: get pending service request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQuery,
					"variables": map[string]interface{}{
						"facilityID": facilityID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
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

func Test_SearchServiceRequests(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLQuery := `
	query searchServiceRequests($searchTerm: String!, $flavour: Flavour!, $requestType: String!, $facilityID: String!){
		searchServiceRequests(searchTerm: $searchTerm, flavour: $flavour, requestType: $requestType, facilityID: $facilityID){
		  id
		  requestType
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
			name: "success: search service request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQuery,
					"variables": map[string]interface{}{
						"searchTerm":  "test",
						"flavour":     feedlib.FlavourConsumer,
						"requestType": "RED_FLAG",
						"facilityID":  facilityID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: unable to search service request without facility ID",
			args: args{
				query: map[string]interface{}{
					"query": graphQLQuery,
					"variables": map[string]interface{}{
						"searchTerm":  "test",
						"requestType": "RED_FLAG",
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

func Test_CompleteVisit(t *testing.T) {
	ctx := context.Background()

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("unable to get graphql headers: %s", err)
		return
	}

	graphQLMutation := `
	mutation completeVisit($staffID: ID!, $serviceRequestID: String!, $bookingID: String!, $notes: String) {
		completeVisit(
		  staffID: $staffID
		  serviceRequestID: $serviceRequestID
		  bookingID: $bookingID
		  notes: $notes
		)
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
			name: "success: create service request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"staffID":          staffID,
						"serviceRequestID": "26b20a42-cbb8-4553-aedb-c539602d04fc",
						"status":           enums.ServiceRequestStatusResolved,
						"notes":            "TEST",
						"bookingID":        bookingID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad: unable to create service request without client ID",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"staffID":          "staffID",
						"serviceRequestID": serviceRequestID,
						"status":           enums.ServiceRequestStatusResolved,
						"notes":            "TEST",
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
