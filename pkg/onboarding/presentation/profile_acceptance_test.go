package presentation_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/base"
)

func composeGraphqlServerRequest(ctx context.Context, query map[string]interface{}) ([]byte, error) {

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		return []byte{}, fmt.Errorf("error in getting headers: %w", err)
	}

	body, err := mapToJSONReader(query)

	if err != nil {
		return []byte{}, fmt.Errorf("unable to get GQL JSON io Reader: %s", err)
	}

	r, err := http.NewRequest(
		http.MethodPost,
		graphQLURL,
		body,
	)

	if err != nil {
		return []byte{}, fmt.Errorf("unable to compose request: %s", err)
	}

	if r == nil {
		return []byte{}, fmt.Errorf("nil request")
	}

	for k, v := range headers {
		r.Header.Add(k, v)
	}

	client := http.Client{
		Timeout: time.Second * testHTTPClientTimeout,
	}
	resp, err := client.Do(r)
	if err != nil {
		return []byte{}, fmt.Errorf("request error: %s", err)
	}
	dataResponse, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("can't read request body: %s", err)
	}

	return dataResponse, nil
}

func updateUserProfile(ctx context.Context) error {
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
	gql := map[string]interface{}{
		"query": graphqlMutation,
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"photoUploadID": "15050000",
				"dateOfBirth":   "2019-01-01",
				"firstName":     "test user",
				"lastName":      "juha",
			},
		},
	}

	dataResp, err := composeGraphqlServerRequest(ctx, gql)
	if err != nil {
		return fmt.Errorf("unable to compose a successful graphql server request: %s", err)
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(dataResp, &data)
	if err != nil {
		return fmt.Errorf("bad data returned")
	}

	return nil
}

func TestUpdateUserProfile(t *testing.T) {

	ctx := base.GetAuthenticatedContext(t)
	// Set Up - update user profile instance on firebase
	err := updateUserProfile(ctx)
	if err != nil {
		log.Printf("error during occurred: %s", err)

	}
}
