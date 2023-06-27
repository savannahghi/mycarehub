package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

func TestRegisterClient(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation registerClient($input: ClientRegistrationInput){
		registerClient(input: $input){
		  id
		  active
		  userID
		  currentFacilityID
		  clientTypes
		  enrollmentDate
		  fhirPatientID
		  emrHealthRecordID
		  treatmentBuddy
		  counselled
		  organisation
		  caregiver
		  chv
		  currentFacilityID
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
			name: "success: register client",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"username":       strings.ToLower(gofakeit.Username()),
							"facility":       mflIdentifier,
							"clientTypes":    []enums.ClientType{enums.ClientTypeDreams},
							"clientName":     gofakeit.Name(),
							"gender":         enumutils.GenderMale,
							"dateOfBirth":    "2000-12-20",
							"phoneNumber":    "+254711880993",
							"enrollmentDate": "2000-02-20",
							"cccNumber":      "202022",
							"counselled":     true,
							"inviteClient":   false,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: facility does not exist",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"username":       gofakeit.Username(),
							"facility":       "93990232",
							"clientTypes":    []enums.ClientType{enums.ClientTypeDreams},
							"clientName":     gofakeit.Name(),
							"gender":         enumutils.GenderMale,
							"dateOfBirth":    "2000-12-20",
							"phoneNumber":    "+254711880993",
							"enrollmentDate": "2000-02-20",
							"cccNumber":      "202022",
							"counselled":     true,
							"inviteClient":   false,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid: missing facility",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"clientTypes":    []enums.ClientType{enums.ClientTypeDreams},
							"clientName":     gofakeit.Name(),
							"gender":         enumutils.GenderMale,
							"dateOfBirth":    "2000-12-20",
							"phoneNumber":    "+254711880993",
							"enrollmentDate": "2000-02-20",
							"cccNumber":      "202022",
							"counselled":     true,
							"inviteClient":   false,
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

func TestSetClientDefaultFacility(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation setClientDefaultFacility ($clientID: ID!, $facilityID: ID!){
		setClientDefaultFacility(clientID: $clientID, facilityID: $facilityID){
				id
				name
				phone
				active
				country
				description
				fhirOrganisationID
				workStationDetails {
					notifications
					surveys
					articles
					messages
					serviceRequests
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
			name: "success: set default facility",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":   clientID,
						"facilityID": facilityID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid: facility does not exist",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":   clientID,
						"facilityID": uuid.NewString(),
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid: user does not exist",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":   uuid.NewString(),
						"facilityID": facilityID,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "invalid: missing facility",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID": clientID,
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "invalid: missing user",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"facilityID": facilityID,
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

func TestGetClientFacilities(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	query getClientFacilities($clientID: ID!, $paginationInput: PaginationsInput!){
		getClientFacilities(clientID: $clientID, paginationInput: $paginationInput){
		  facilities{
			id
			name
			phone
			active
			country
			description
			fhirOrganisationID
			workStationDetails{
			  notifications
			  surveys
			  articles
			  messages
			  serviceRequests
			}
		  }
		  pagination{
			limit
			currentPage
			count
			totalPages
			nextPage
			previousPage
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
			name: "success: get client facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID": clientID,
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
			name: "fail: unable to get client facilities",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID": "clientID",
						"paginationInput": map[string]interface{}{
							"limit":       10,
							"currentPage": 1,
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

func TestRegisterExistingUserAsClient(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation registerExistingUserAsClient($input: ExistingUserClientInput!){
		registerExistingUserAsClient(input: $input){
		  id
		  active
		  clientTypes
		  enrollmentDate
		  fhirPatientID
		  emrHealthRecordID
		  treatmentBuddy
		  counselled
		  organisation
		  userID
		  currentFacilityID
		  chv
		  caregiver
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
			name: "success: register an existing staff as a client in the same program",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"userID":         staffUserToAddAsClient,
							"programID":      programID,
							"facilityID":     facilityID2,
							"clientTypes":    []string{"PMTCT"},
							"enrollmentDate": "2022-02-20",
							"cccNumber":      "12345",
							"counselled":     true,
							"inviteClient":   true,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: register an existing client as a client in another program",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"userID":         staffUserToAddAsClient,
							"programID":      programID2,
							"facilityID":     facilityID2,
							"clientTypes":    []string{"PMTCT"},
							"enrollmentDate": "2022-02-20",
							"cccNumber":      "12345",
							"counselled":     true,
							"inviteClient":   true,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: User already has a client profile in the program,",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"userID":         userID,
							"programID":      programID,
							"facilityID":     facilityToAddExistingStaff,
							"clientTypes":    []string{"PMTCT"},
							"enrollmentDate": "2022-02-20",
							"cccNumber":      "12345",
							"counselled":     true,
							"inviteClient":   true,
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "fail: unable to register existing user as client",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"userID":         staffUserToAddAsClient,
							"programID":      "programID",
							"facilityID":     facilityToAddExistingStaff,
							"clientTypes":    []string{"PMTCT"},
							"enrollmentDate": "2022-02-20",
							"cccNumber":      "12345",
							"counselled":     true,
							"inviteClient":   true,
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

func TestRegisterExistingUserAsCaregiver(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation registerExistingUserAsCaregiver($userID: ID!, $caregiverNumber: String!){
		registerExistingUserAsCaregiver(userID: $userID, caregiverNumber: $caregiverNumber){
		  id
		  user{
			id
			username
			name
			gender
			active
			contacts{
			  id
			  contactType
			  contactValue
			  active
			  optedIn
			}
		  }
		  caregiverNumber
		  isClient
		  consent{
			consentStatus
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
			name: "success: get register existing user as caregiver",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"userID":          userIDtoAssignStaff,
						"caregiverNumber": "12345",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: unable to register existing user as caregiver with invalid userID",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"userID":          "userIDtoAssignStaff",
						"caregiverNumber": "12345",
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

func TestCheckIdentifierExists(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	query checkIdentifierExists($identifierType: UserIdentifierType!, $identifierValue: String!){
		checkIdentifierExists(identifierType: $identifierType, identifierValue: $identifierValue)
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
			name: "success: check ccc identifier exists",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"identifierType":  enums.UserIdentifierTypeCCC,
						"identifierValue": "123456",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: check national id identifier exists",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"identifierType":  enums.UserIdentifierTypeNationalID,
						"identifierValue": "12345678",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: unable to check identifier exists, invalid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"identifierType": enums.UserIdentifierTypeCCC,
						"invalid":        "123456",
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

func TestCheckIfPhoneExists(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
		query checkIfPhoneExists($phoneNumber: String!){
			checkIfPhoneExists(phoneNumber: $phoneNumber)
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
			name: "success: check if phone exists, pro",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phoneNumber": gofakeit.Phone(),
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "success: check if phone exists, consumer",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phoneNumber": gofakeit.Phone(),
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: invalid phone",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"phoneNumber": "invalid",
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

func TestVerifyClientPinResetServiceRequest(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation verifyClientPinResetServiceRequest($serviceRequestID: String!, $status: PINResetVerificationStatus!, $physicalIdentityVerified: Boolean!){
		verifyClientPinResetServiceRequest(serviceRequestID: $serviceRequestID, status: $status, physicalIdentityVerified: $physicalIdentityVerified)
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
			name: "success: verify client pin reset service request",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"serviceRequestID":         clientsServiceRequestID,
						"status":                   enums.PINResetVerificationStatusApproved,
						"physicalIdentityVerified": true,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "failure: invalid status",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"serviceRequestID":         clientsServiceRequestID,
						"status":                   "invalid",
						"physicalIdentityVerified": true,
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

func Test_GetClientProfileByCCCNumber(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
		query getClientProfileByCCCNumber($CCCNumber: String!){
			getClientProfileByCCCNumber(CCCNumber: $CCCNumber){
				id
				active
				clientTypes
				treatmentEnrollmentDate
				fhirPatientID
				healthRecordID
				treatmentBuddy
				clientCounselled
				chvUserID
				chvUserName
				caregiverID
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
			name: "success: get client profile by ccc number",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"CCCNumber": "123456",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: ccc does not exist",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"CCCNumber": "4343434343",
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

func Test_ConsentToAClientCaregiver(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation consentToAClientCaregiver($clientID: ID!, $caregiverID: ID!, $consent: ConsentState!) {
		consentToAClientCaregiver(
		  clientID: $clientID
		  caregiverID: $caregiverID
		  consent: $consent
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
			name: "success: consent to a caregiver",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":    testClientToUpdateConsent,
						"caregiverID": testCaregiverID3,
						"consent":     enums.ConsentStateAccepted,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: unable to consent to a caregiver",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":    "123456",
						"caregiverID": testCaregiverID3,
						"consent":     enums.ConsentStateAccepted,
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

func Test_ConsentToManagingClient(t *testing.T) {
	ctx := context.Background()
	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	headers, err := GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("failed to get GraphQL headers: %v", err)
		return
	}

	graphqlMutation := `
	mutation consentToManagingClient($caregiverID: ID!, $clientID: ID!, $consent: ConsentState!) {
		consentToManagingClient(
		  caregiverID: $caregiverID
		  clientID: $clientID
		  consent: $consent
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
			name: "success: consent to manage a client",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":    testClientToUpdateConsent,
						"caregiverID": testCaregiverID3,
						"consent":     enums.ConsentStateAccepted,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "fail: unable to consent to consent to managing a client",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"clientID":    "123456",
						"caregiverID": testCaregiverID3,
						"consent":     enums.ConsentStateAccepted,
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
