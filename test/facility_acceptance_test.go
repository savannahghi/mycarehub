package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
)

func TestCreateAFacility(t *testing.T) {
	token, err := loginByPhone(t, testPhone)
	if err != nil {
		t.Errorf("Error when loggin in by phone: %v", err)
	}
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers := getGraphHeaders(token.IDToken)

	mflcode := gofakeit.Name()

	graphqlMutation := `
	mutation createFacility($input: FacilityInput!) {
		createFacility (input: $input) {
		  name
		  code
		  active
		  county
		  description
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
			name: "success: create a facility with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"code":        mflcode,
							"active":      true,
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: missing name param",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"code":        mflcode,
							"active":      true,
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
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
						"input": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"active":      true,
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
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
						"input": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"code":        mflcode,
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
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
						"input": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"code":        mflcode,
							"active":      true,
							"description": "located at Giddo plaza building town",
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
						"input": map[string]interface{}{
							"name":   "Mediheal Hospital (Nakuru) Annex",
							"code":   mflcode,
							"active": true,
							"county": "Nakuru",
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
						"input": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"code":        mflcode,
							"active":      "invalid",
							"county":      "Nakuru",
							"description": "located at Giddo plaza building town",
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: invalid value for county",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"name":        "Mediheal Hospital (Nakuru) Annex",
							"code":        mflcode,
							"active":      true,
							"county":      "Kanairo",
							"description": "located at Giddo plaza building town",
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

			dataResponse, err := ioutil.ReadAll(resp.Body)
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
	pg, pgError := gorm.NewPGInstance()
	if pgError != nil {
		t.Errorf("can't instantiate test repository: %v", pgError)
	}

	// Teardown
	pgError = pg.DB.Unscoped().Where("mfl_code", mflcode).Delete(&gorm.Facility{}).Error
	if pgError != nil {
		t.Errorf("Error deleting facility: %v", pgError)
		return
	}
}
