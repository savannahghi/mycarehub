package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
)

func TestRegisterAgent(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	up := dto.RegisterAgentInput{
		FirstName:   "Test",
		LastName:    "AgentTest",
		Gender:      "male",
		PhoneNumber: "0700011122",
		Email:       "test.agent@test.com",
		DateOfBirth: base.Date{Year: 1990, Month: 1, Day: 1},
	}

	graphqlMutation := `
	mutation registerAgent($input: RegisterAgentInput!) {
		registerAgent(input: $input) {
		  primaryPhone
		  termsAccepted
		  suspended
		  userBioData {
			firstName
			lastName
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
			name: "success: create agent profile",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"firstName":   up.FirstName,
							"lastName":    up.LastName,
							"gender":      up.Gender,
							"phoneNumber": up.PhoneNumber,
							"email":       up.Email,
							"dateOfBirth": up.DateOfBirth,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid:wrong variable type ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"firstName":   120,
							"lastName":    up.LastName,
							"gender":      up.Gender,
							"phoneNumber": up.PhoneNumber,
							"email":       up.Email,
							"dateOfBirth": up.DateOfBirth,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid:should not create agents when input is empty",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"firstName":   "",
							"lastName":    "",
							"gender":      "",
							"phoneNumber": "",
							"email":       "",
							"dateOfBirth": nil,
						},
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid:invalid phone number",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"firstName":   up.FirstName,
							"lastName":    up.LastName,
							"gender":      up.Gender,
							"phoneNumber": "0712345",
							"email":       up.Email,
							"dateOfBirth": up.DateOfBirth,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid:unknown gender type provided",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"firstName":   up.FirstName,
							"lastName":    up.LastName,
							"gender":      "cow",
							"phoneNumber": up.PhoneNumber,
							"email":       up.Email,
							"dateOfBirth": up.DateOfBirth,
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
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned. Expected %v, got %v", tt.wantStatus, resp.StatusCode)
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
				v, ok := data["errors"]
				if ok {
					t.Errorf("error not expected %v", v)
					return
				}
				// check/assert the returned data/response
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					for nestedKey := range nestedMap {
						if nestedKey == "registerAgent" {
							output, ok := nestedMap[nestedKey].(map[string]interface{})
							if !ok {
								t.Errorf("can't cast nestedKey to map[string]interface{}")
								return
							}
							_, present := output["userBioData"].(map[string]interface{})
							if !present {
								t.Errorf("Biodata not present in output")
								return
							}
						}
					}
				}
			}

		})
	}
	// perform tear down; remove user
	_, err := RemoveTestUserByPhone(t, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user employee: %s", err)
	}
}

func TestActivateAgent(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	err := setRoleForUserWithPhone(base.TestUserPhoneNumber, base.RoleTypeEmployee, headers)
	if err != nil {
		t.Errorf("failed to set employee role for logged in user: %s", err)
		return
	}

	agentInput := dto.RegisterAgentInput{
		FirstName:   "Test",
		LastName:    "Agent",
		Gender:      "male",
		PhoneNumber: "0704711122",
		Email:       "test.agent@test.com",
		DateOfBirth: base.Date{Year: 1990, Month: 1, Day: 1},
	}
	agent, err := registerTestAgent(agentInput, headers)
	if err != nil {
		t.Errorf("failed to create test agent: %s", err)
		return
	}

	graphqlMutation := `
	mutation activateAgent($agentID: String!) {
		activateAgent(agentID: $agentID)
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
			name: "success: activate agent profile",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"agentID": agent.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid : wrong agent id",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"agentID": "agarbagevalue",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid : should not update when inputs are empty ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"agentID": "",
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
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned. Expected %v, got %v", tt.wantStatus, resp.StatusCode)
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
				v, ok := data["errors"]
				if ok {
					t.Errorf("error not expected %v", v)
					return
				}
				// check/assert the returned data/response
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					for nestedKey := range nestedMap {
						if nestedKey == "activateAgent" {
							output, ok := nestedMap[nestedKey].(bool)
							if !ok {
								t.Errorf("can't cast nestedKey to map[string]interface{}")
								return
							}
							if !output {
								t.Errorf("expected the \"activateAgent\" property to be set to \"true\"")
								return
							}
						}
					}
				}
			}

		})
	}
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user employee: %s", err)
	}

	_, err = RemoveTestUserByPhone(t, agentInput.PhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user agent: %s", err)
	}
}

func TestDeactivateAgent(t *testing.T) {
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers := setUpLoggedInTestUserGraphHeaders(t)

	err := setRoleForUserWithPhone(base.TestUserPhoneNumber, base.RoleTypeEmployee, headers)
	if err != nil {
		t.Errorf("failed to set employee role for logged in user: %s", err)
		return
	}

	agentInput := dto.RegisterAgentInput{
		FirstName:   "Test",
		LastName:    "Agent",
		Gender:      "male",
		PhoneNumber: "0704711122",
		Email:       "test.agent@test.com",
		DateOfBirth: base.Date{Year: 1990, Month: 1, Day: 1},
	}
	agent, err := registerTestAgent(agentInput, headers)
	if err != nil {
		t.Errorf("failed to create test agent: %s", err)
		return
	}

	graphqlMutation := `
	mutation deactivateAgent($agentID: String!) {
		deactivateAgent(agentID: $agentID)
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
			name: "success: deactivate agent profile",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"agentID": agent.ID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid : non-existent agent id",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"agentID": "agarbagevalue",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid : should not update when inputs are empty ",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"agentID": "",
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
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status response returned. Expected %v, got %v", tt.wantStatus, resp.StatusCode)
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
				v, ok := data["errors"]
				if ok {
					t.Errorf("error not expected %v", v)
					return
				}
				// check/assert the returned data/response
				for key := range data {
					nestedMap, ok := data[key].(map[string]interface{})
					if !ok {
						t.Errorf("cannot cast key value of %v to type map[string]interface{}", key)
						return
					}
					for nestedKey := range nestedMap {
						if nestedKey == "deactivateAgent" {
							output, ok := nestedMap[nestedKey].(bool)
							if !ok {
								t.Errorf("can't cast nestedKey to map[string]interface{}")
								return
							}
							if !output {
								t.Errorf("expected the \"deactivateAgent\" property to be set to \"true\"")
								return
							}
						}
					}
				}
			}

		})
	}
	// perform tear down; remove user
	_, err = RemoveTestUserByPhone(t, base.TestUserPhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user employee: %s", err)
	}

	_, err = RemoveTestUserByPhone(t, agentInput.PhoneNumber)
	if err != nil {
		t.Errorf("unable to remove test user agent: %s", err)
	}
}

func registerTestAgent(agentInput dto.RegisterAgentInput, headers map[string]string) (*base.UserProfile, error) {
	url := fmt.Sprintf("%s/graphql", baseURL)

	graphqlMutation := `
	mutation registerAgent($input: RegisterAgentInput!) {
		registerAgent(input: $input) {
			id
			userName
		  	primaryPhone
		  	termsAccepted
		  	suspended
		  	userBioData {
				firstName
				lastName
				gender
		  	}
		}
	}`
	query := map[string]interface{}{
		"query": graphqlMutation,
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"firstName":   agentInput.FirstName,
				"lastName":    agentInput.LastName,
				"gender":      agentInput.Gender,
				"phoneNumber": agentInput.PhoneNumber,
				"email":       agentInput.Email,
				"dateOfBirth": agentInput.DateOfBirth,
			},
		},
	}

	body, err := mapToJSONReader(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query to GQL JSON io.Reader: %s", err)
	}

	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, fmt.Errorf("unable to compose request: %s", err)
	}
	if request == nil {
		return nil, fmt.Errorf("nil request")
	}

	for header, value := range headers {
		request.Header.Add(header, value)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("request error: %s", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to register agent: expected status to be %v got %v ", http.StatusOK, response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %s", err)
	}

	// Process response
	var responseContent map[string]interface{}
	err = json.Unmarshal(data, &responseContent)
	if err != nil {
		return nil, fmt.Errorf("unable to marshall response: %v", err)
	}
	responseData, ok := responseContent["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("expected response to contain a \"data\" object")
	}
	agentData, ok := responseData["registerAgent"]
	if !ok {
		return nil, fmt.Errorf("expected response to contain a \"data.registerAgent\" object")
	}
	agentDataAsJSON, err := json.Marshal(agentData)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal agent data to JSON: %s", err)
	}

	var agent base.UserProfile
	err = json.Unmarshal(agentDataAsJSON, &agent)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshall response: %v", err)
	}

	return &agent, nil
}
