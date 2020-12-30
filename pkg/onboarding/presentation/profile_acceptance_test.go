package presentation_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

// func updateTestUserProfile(ctx context.Context, input domain.UserProfileInput) error {
// 	graphqlMutation := `
// 	mutation updateUserProfile($input:UserProfileInput!){
// 		updateUserProfile(input: $input){
// 			userName
// 			verifiedIdentifiers{
// 				uid
// 				timestamp
// 				loginProvider
// 			}
// 			PrimaryPhone
// 			PrimaryEmailAddress
// 			pushTokens
// 			userBioData{
// 				firstName
// 				lastName
// 				dateOfBirth
// 				gender
// 			}

// 		}
// 	}`
// 	gql := map[string]interface{}{
// 		"query": graphqlMutation,
// 		"variables": map[string]interface{}{
// 			"input": map[string]interface{}{
// 				"photoUploadID": input.PhotoUploadID,
// 				"dateOfBirth":   input.DateOfBirth,
// 				"firstName":     input.FirstName,
// 				"lastName":      input.LastName,
// 			},
// 		},
// 	}

// 	dataResp, err := ComposeGraphqlServerRequest(ctx, gql)
// 	if err != nil {
// 		return fmt.Errorf("unable to compose a successful graphql server request: %s", err)
// 	}

// 	data := map[string]interface{}{}
// 	err = json.Unmarshal(dataResp, &data)
// 	if err != nil {
// 		return fmt.Errorf("bad data returned")
// 	}
// 	return nil
// }

func TestUpdateUserProfile(t *testing.T) {
	// create a user and thier profile
	_, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		// return
	}
	ctx := base.GetAuthenticatedContext(t)

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	// update the user profile that was created
	dateOfBirth := base.Date{
		Day:   1,
		Year:  2019,
		Month: 4,
	}
	firstName := "kamau"
	lastName := "mwas"
	up := domain.UserProfileInput{
		PhotoUploadID: "12345",
		DateOfBirth:   &dateOfBirth,
		FirstName:     &firstName,
		LastName:      &lastName,
	}

	graphqlMutation := `
	mutation updateUserProfile($input:UserProfileInput!){
		updateUserProfile(input: $input){
			userName
			verifiedIdentifiers{
				uid
				timestamp
				loginProvider
			}
			PrimaryPhone
			PrimaryEmailAddress
			pushTokens
			userBioData{
				firstName
				lastName
				dateOfBirth
				gender
			}
			
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
			wantErr:    true,
		},
		{
			name: "failure: update profile with valid empty payload",
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
			wantStatus: http.StatusBadRequest,
			wantErr:    true,
		},
		{
			name: "failure: update profile with invalid inputs",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"photoUploadID": "1234",
							"dateOfBirth":   "2019-01-01",
							"firstName":     "mwas",
							"lastName":      "sss",
						},
					},
				},
			},
			wantStatus: http.StatusBadRequest,
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
				// for key := range data {
				// 	nestedMap, ok := data[key].(map[string]interface{})
				// 	if !ok {
				// 		t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
				// 		return
				// 	}
				// 	if key == "data" {
				// 		_, ok := nestedMap[""]
				// 		if !ok {
				// 			t.Errorf("no  key found")
				// 			return
				// 		}
				// 	}
				// }
			}
			//TODO mathenge uncomment below after fixing adding of verified identifiers
			// if tt.wantStatus != resp.StatusCode {
			// 	t.Errorf("Bad status reponse returned")
			// 	return
			// }

		})
	}

}
