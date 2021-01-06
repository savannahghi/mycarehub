package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestGraphQLFindProvider(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	type args struct {
		query map[string]interface{}
	}

	variables := map[string]interface{}{
		"filters": map[string]string{
			"search": "khan",
		},
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "valid query request",
			args: args{
				query: map[string]interface{}{
					"query": `
						query findProvider($pagination: PaginationInput,$filter:[BusinessPartnerFilterInput],$sort:[BusinessPartnerSortInput]) {
							findProvider(pagination:$pagination,filter:$filter,sort:$sort){
								edges {
									cursor
									node {
									id
									name
									sladeCode
									}
								}
								pageInfo {
									hasNextPage
									hasPreviousPage
									startCursor
									endCursor
								}
							}
						}`,
					"variables": variables,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},

		{
			name: "invalid query request -- sladecode not same as sladeCode",
			args: args{
				query: map[string]interface{}{
					"query": `
						query findProvider($pagination: PaginationInput,$filter:[BusinessPartnerFilterInput],$sort:[BusinessPartnerSortInput]) {
							findProvider(pagination:$pagination,filter:$filter,sort:$sort){
								edges {
									cursor
									node {
									id
									name
									sladecode
									}
								}
								pageInfo {
									hasNextPage
									hasPreviousPage
									startCursor
									endCursor
								}
							}
						}`,
					"variables": variables,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			body, err := mapToJSONReader(tt.args.query)
			if err != nil {
				t.Errorf("unable to get request JSON io Reader: %s", err)
				return
			}
			r, err := http.NewRequest(
				http.MethodPost,
				graphQLURL,
				body,
			)

			if err != nil {
				t.Errorf("can't create new request: %v", err)
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

//todo(dexter) investigate why this is failing specifically on CI in my next PR
// func TestGraphQLFindBranch(t *testing.T) {
// 	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
// 	headers := setUpLoggedInTestUserGraphHeaders(t)

// 	type args struct {
// 		query map[string]interface{}
// 	}

// 	variables := map[string]interface{}{
// 		"filters": map[string]string{
// 			"search": "khan",
// 		},
// 	}

// 	tests := []struct {
// 		name       string
// 		args       args
// 		wantStatus int
// 		wantErr    bool
// 	}{
// 		{
// 			name: "valid query request",
// 			args: args{
// 				query: map[string]interface{}{
// 					"query": `
// 						query findBranch($pagination: PaginationInput,$filter:[BranchFilterInput],$sort:[BranchSortInput]) {
// 							findBranch(pagination:$pagination,filter:$filter,sort:$sort){
// 							edges {
// 								cursor
// 								node {
// 								id
// 								name
// 								organizationSladeCode
// 								branchSladeCode
// 								}
// 							}
// 							pageInfo {
// 								hasNextPage
// 								hasPreviousPage
// 								startCursor
// 								endCursor
// 							}
// 							}
// 						}`,
// 					"variables": variables,
// 				},
// 			},
// 			wantStatus: http.StatusOK,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "invalid query request - branchSladecode not same as branchSladeCode",
// 			args: args{
// 				query: map[string]interface{}{
// 					"query": `
// 						query findBranch($pagination: PaginationInput,$filter:[BranchFilterInput],$sort:[BranchSortInput]) {
// 							findBranch(pagination:$pagination,filter:$filter,sort:$sort){
// 							edges {
// 								cursor
// 								node {
// 								id
// 								name
// 								organizationSladeCode
// 								branchSladecode
// 								}
// 							}
// 							pageInfo {
// 								hasNextPage
// 								hasPreviousPage
// 								startCursor
// 								endCursor
// 							}
// 							}
// 						}`,
// 					"variables": variables,
// 				},
// 			},
// 			wantStatus: http.StatusOK,
// 			wantErr:    true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {

// 			body, err := mapToJSONReader(tt.args.query)
// 			if err != nil {
// 				t.Errorf("unable to get request JSON io Reader: %s", err)
// 				return
// 			}
// 			r, err := http.NewRequest(
// 				http.MethodPost,
// 				graphQLURL,
// 				body,
// 			)

// 			if err != nil {
// 				t.Errorf("can't create new request: %v", err)
// 				return
// 			}

// 			if r == nil {
// 				t.Errorf("nil request")
// 				return
// 			}

// 			for k, v := range headers {
// 				r.Header.Add(k, v)
// 			}
// 			client := http.Client{
// 				Timeout: time.Second * testHTTPClientTimeout,
// 			}
// 			resp, err := client.Do(r)
// 			if err != nil {
// 				t.Errorf("HTTP error: %v", err)
// 				return
// 			}

// 			if !tt.wantErr && resp == nil {
// 				t.Errorf("unexpected nil response (did not expect an error)")
// 				return
// 			}

// 			if tt.wantErr {
// 				return
// 			}

// 			data, err := ioutil.ReadAll(resp.Body)
// 			if err != nil {
// 				t.Errorf("can't read response body: %v", err)
// 				return
// 			}
// 			if data == nil {
// 				t.Errorf("nil response body data")
// 				return
// 			}

// 			if tt.wantStatus != resp.StatusCode {
// 				t.Errorf("expected status %d, got %d and response %s", tt.wantStatus, resp.StatusCode, string(data))
// 				return
// 			}

// 			if !tt.wantErr && resp == nil {
// 				t.Errorf("unexpected nil response (did not expect an error)")
// 				return
// 			}
// 		})
// 	}
// }
