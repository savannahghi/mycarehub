package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestUpdateUserProfile(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	// update the user profile that was created
	dateOfBirth := base.Date{
		Day:   1,
		Year:  2019,
		Month: 4,
	}
	firstName := "kamau"
	lastName := "mwas"
	uploadID := "photo-upload-id"
	up := resources.UserProfileInput{
		PhotoUploadID: &uploadID,
		DateOfBirth:   &dateOfBirth,
		FirstName:     &firstName,
		LastName:      &lastName,
	}

	graphqlMutation := `
	mutation updateUserProfile($input:UserProfileInput!){
		updateUserProfile(input: $input){
			id
		}
	}`

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
			name: "success: update profile with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"photoUploadID": up.PhotoUploadID,
							"dateOfBirth":   up.DateOfBirth,
							"firstName":     up.FirstName,
							"lastName":      up.LastName,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid : wrong dateOfBirth type ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"photoUploadID": "",
							"dateOfBirth":   "",
							"firstName":     "",
							"lastName":      "",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "valid : should not update when inputs are empty ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"photoUploadID": "",
							"firstName":     "",
							"lastName":      "",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "valid : should update firstname only ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"firstName": "new-firstName",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "valid : should update lastname only ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"lastName": "new-lastname",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "valid : should update photoUploadID only ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"photoUploadID": "new-photoUploadID",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "valid : should update dateOfBirth only ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"dateOfBirth": "1997-12-10", // YYYYMMDD
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid : wrong dateOfBirth format ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"dateOfBirth": "10-12-1997",
						},
					},
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

			res, _ := httputil.DumpResponse(resp, true)
			logrus.Printf("%v", string(res))

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
				v, ok := data["errors"]
				if ok {
					t.Errorf("error not expected %v", v)
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned. Expected %v, got %v", tt.wantStatus, resp.StatusCode)
				return
			}
		})
	}
	// perform tear down; remove user
	_, err := RemoveTestUserByPhone(t, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}

func TestUserProfile(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	graphqlMutation := `
	query getUserProfile {
		userProfile {
		  id
		  userName
		  verifiedIdentifiers {
			uid
			loginProvider
			timestamp
		  }
		  primaryPhone
		  primaryEmailAddress
		  secondaryPhoneNumbers
		  secondaryEmailAddresses
		  pushTokens
		  termsAccepted
		  suspended
		  photoUploadID
		  userBioData {
			firstName
			lastName
			gender
			dateOfBirth
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
			name: "success: retrieve user profile for a registred user",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
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
	// perform tear down; remove user
	_, err := RemoveTestUserByPhone(t, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}

func TestSupplierProfile(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	graphqlMutation := `
	query getSUpplierProfile{
		supplierProfile{
		  id
		  profileID
		  supplierId
		  sladeCode
		  hasBranches
		  partnerSetupComplete
		  isOrganizationVerified
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
			name: "success: retrieve user profile for a registred user",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
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
	// perform tear down; remove user
	_, err := RemoveTestUserByPhone(t, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user: %s", err)
	}
}
