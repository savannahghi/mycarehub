package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func TestGetContent(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlQuery := `
	query getContent($categoryIDs: [Int!], $categoryNames: [String!], $limit: String!, $clientID: String) {
		getContent(
		  categoryIDs: $categoryIDs
		  categoryNames: $categoryNames
		  limit: $limit
		  clientID: $clientID
		) {
		  items{
			id
			title
			date
			meta{
			  contentType
			  contentDetailURL
			}
			intro
			authorName
			itemType
			timeEstimateSeconds
			body
			heroImage{
			  id
			  meta{
				type
				imageDetailUrl
				imageDownloadUrl
			  }
			  title
			}
			heroImage{
			  id
			  meta{
				type
				imageDetailUrl
				imageDownloadUrl
			  }
			}
		  }
		  meta{
			totalCount
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
			name: "success: get content",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"categoryIDs": []int{1},
						"limit":       10,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: get content as caregiver",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"categoryID": []int{1},
						"limit":      10,
						"clientID":   clientID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: missing parameters",
			args: args{
				query: map[string]interface{}{
					"query": graphqlQuery,
					"variables": map[string]interface{}{
						"categoryID": []int{1},
						"clientID":   clientID,
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
				t.Errorf("Bad status response returned, expected %v, got %v, %v", tt.wantStatus, resp.StatusCode, data["errors"])
				return
			}
		})
	}
}

func Test_FetchContent(t *testing.T) {
	fetchContentURL := fmt.Sprintf("%s/%s", baseURL, "content")

	type args struct {
		limit string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success: fetch content",
			args: args{
				limit: "15",
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to fetch content",
			args: args{
				limit: "20",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs, err := json.Marshal(tt.args.limit)
			if err != nil {
				t.Errorf("unable to marshal test item to JSON: %s", err)
			}

			payload := bytes.NewBuffer(bs)

			r, err := http.NewRequest(
				http.MethodGet,
				fetchContentURL,
				payload,
			)
			if err != nil {
				t.Errorf("unable to compose request: %s", err)
				return
			}

			if r == nil {
				t.Errorf("nil request")
				return
			}

			r.Header.Add("Content-Type", "application/json")

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

			if tt.wantErr {
				errorMap := map[string]interface{}{}
				err = json.Unmarshal(dataResponse, &errorMap)
				if err != nil {
					t.Errorf("unable to unmarshal response: %s", err)
					return
				}
				if errorMap["error"] == nil {
					t.Errorf("expected an error but got nil")
					return
				}

			}
			if !tt.wantErr {
				data := &domain.Content{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}

}
