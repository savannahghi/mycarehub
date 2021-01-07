package tests

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

// CreatedUserGraphQLHeaders updates the authorization header with the
// bearer(ID) token of the created user during test
// TODO:(muchogo)  create a reusable function in base that accepts a UID
// 				or modify the base.GetGraphQLHeaders(ctx) extra UID argument
func CreatedUserGraphQLHeaders(idToken *string) (map[string]string, error) {
	ctx := context.Background()

	authHeaderBearerToken := fmt.Sprintf("Bearer %v", *idToken)

	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in getting headers: %w", err)
	}

	headers["Authorization"] = authHeaderBearerToken

	return headers, nil
}
func TestAddPartnerType(t *testing.T) {
	ctx := base.GetAuthenticatedContext(t)

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")
	headers, err := base.GetGraphQLHeaders(ctx)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphqlMutation := `
	mutation addPartnerType($name:String!, $partnerType:PartnerType!){
		addPartnerType(name: $name, partnerType:$partnerType)
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
			name: "success: add partner type with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"name":        "juha kalulu",
						"partnerType": "RIDER",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true, // TODO fixme: the logged in user must have a registred profile
		},
		{
			name: "failure: add partner type with non existent partner type",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"name":        "juha",
						"partnerType": "NOT FOUND",
					},
				},
			},
			wantStatus: http.StatusUnprocessableEntity,
			wantErr:    true,
		},
		{
			name: "failure: add partner type with bogus payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"name":        "*",
						"partnerType": "*",
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}

}

func TestSetUpSupplier_acceptance(t *testing.T) {

	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphqlMutation := `
	mutation setUpSupplier($input: AccountType!) {
		setUpSupplier(accountType: $input) {
		  accountType
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
			name: "individual supplier setup with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": "INDIVIDUAL",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "organisation supplier setup with valid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": "ORGANISATION",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid case - supplier setup with invalid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": "",
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestSuspendSupplier_acceptance(t *testing.T) {

	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphqlMutation := `mutation{
		suspendSupplier
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
			name: "valid - Suspend existing supplier",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid - Suspend supplier using an invalid payload",
			args: args{
				query: map[string]interface{}{
					"query": "invalid mutation",
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestSupplierEDILogin(t *testing.T) {
	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	sladeCode := "1"

	if err != nil {
		t.Errorf("error getting headers: %w", err)
		return
	}

	graphQLMutationPayload := `
	mutation supplierEDILogin($username: String!, $password: String!, $sladeCode: String!){
		supplierEDILogin(username: $username, password:$password, sladeCode: $sladeCode){
		  edges{
			cursor
			node{
			  id
			  name
			  branchSladeCode
			  organizationSladeCode
			}
			
		  }
		  pageInfo{
			hasNextPage
			hasPreviousPage
			startCursor
			endCursor
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
			name: "valid edi portal login mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"username":  "avenue-4190@healthcloud.co.ke",
						"password":  "test provider",
						"sladeCode": "BRA-PRO-4190-4",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid edi portal login mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"username":  "avenue-4190@healthcloud.co.ke",
						"password":  "test provider",
						"sladeCode": "WRONG SLADE CODE",
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "valid edi core login mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"username":  "bewell@slade360.co.ke",
						"password":  "please change me",
						"sladeCode": sladeCode,
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid edi core login mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"username":  "bewell@slade360.co.ke",
						"password":  "please change me",
						"sladeCode": "BOGUS SLADE CODE",
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestAddIndividualPractitionerKYC(t *testing.T) {

	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphQLMutationPayload := `
	mutation addIndividualPractitionerKYC($input:IndividualPractitionerInput!){
		addIndividualPractitionerKYC(input:$input) {
		identificationDoc {
		  identificationDocType
		  identificationDocNumber
		  identificationDocNumberUploadID
	}
		registrationNumber
		KRAPIN
		KRAPINUploadID
		practiceServices
		practiceLicenseUploadID
		supportingDocumentsUploadID
		cadre
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
			name: "valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "12345",
								"identificationDocNumberUploadID": "12345",
							},
							"registrationNumber":          "12345",
							"KRAPIN":                      "12345",
							"KRAPINUploadID":              "12345",
							"practiceLicenseID":           "12345",
							"practiceServices":            []string{"OUTPATIENT_SERVICES"},
							"practiceLicenseUploadID":     "12345",
							"cadre":                       "DOCTOR",
							"supportingDocumentsUploadID": []string{"123456"},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request - wrong input",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutationPayload,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "12345",
								"identificationDocNumberUploadID": "12345",
							},
							"registrationNumber":          "12345",
							"KRAPIN":                      "12345",
							"KRAPINUploadID":              "12345",
							"practiceLicenseID":           "12345",
							"practiceServices":            []string{"OUTPATIENT_SERVICES"},
							"practiceLicenseUploadID":     12345,
							"cadre":                       "DOCTOR",
							"supportingDocumentsUploadID": []string{"123456"},
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}
func TestAddOrganizationProviderKYC(t *testing.T) {
	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphqlMutation := `
	mutation   addOrganizationProviderKYC($input:OrganizationProviderInput!){
		addOrganizationProviderKYC(input:$input) {
		organizationTypeName
		certificateOfIncorporation
		certificateOfInCorporationUploadID
		registrationNumber
		KRAPIN
		KRAPINUploadID
		practiceServices
		practiceLicenseUploadID
		practiceLicenseID
		supportingDocumentsUploadID
		directorIdentifications{
		  identificationDocType
				identificationDocNumber
				identificationDocNumberUploadID
		}
		cadre
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
			name: "valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"directorIdentifications": []map[string]interface{}{
								{
									"identificationDocType":           "NATIONALID",
									"identificationDocNumber":         "12345678",
									"identificationDocNumberUploadID": "12345678",
								},
							},
							"organizationTypeName":               "LIMITED_COMPANY",
							"certificateOfIncorporation":         "CERT-123456",
							"certificateOfInCorporationUploadID": "CERT-UPLOAD-123456",
							"registrationNumber":                 "REG-123456",
							"KRAPIN":                             "KRA-123456789",
							"KRAPINUploadID":                     "KRA-UPLOAD-123456789",
							"practiceServices":                   []string{"OUTPATIENT_SERVICES"},
							"practiceLicenseID":                  "PRAC-123456",
							"practiceLicenseUploadID":            "PRAC-UPLOAD-123456",
							"supportingDocumentsUploadID":        []string{"SUPP-UPLOAD-123456"},
							"cadre":                              "DOCTOR",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"directorIdentifications": []map[string]interface{}{
								{
									"identificationDocType":           "NATIONALID",
									"identificationDocNumber":         "12345678",
									"identificationDocNumberUploadID": "12345678",
								},
							},
							"organizationTypeName":               "LIMITED_COMPANY",
							"certificateOfIncorporation":         "CERT-123456",
							"certificateOfInCorporationUploadID": "CERT-UPLOAD-123456",
							"registrationNumber":                 "REG-123456",
							"KRAPIN":                             123456789,
							"KRAPINUploadID":                     "KRA-UPLOAD-123456789",
							"practiceServices":                   []string{"OUTPATIENT_SERVICES"},
							"practiceLicenseID":                  "PRAC-123456",
							"practiceLicenseUploadID":            "PRAC-123456",
							"supportingDocumentsUploadID":        []string{"SUPP-UPLOAD-123456"},
							"cadre":                              "DOCTOR",
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
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got error: %w", err)
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

func TestAddIndividualPharmaceuticalKYC(t *testing.T) {
	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphQLMutation := `
	mutation addIndividualPharmaceuticalKYC($input:IndividualPharmaceuticalInput!){
		addIndividualPharmaceuticalKYC(input:$input) {
		identificationDoc{
		  identificationDocType
		  identificationDocNumber
		  identificationDocNumberUploadID
	}
		registrationNumber
		KRAPIN
		KRAPINUploadID
		practiceLicenseID
		practiceLicenseUploadID
		supportingDocumentsUploadID
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
			name: "valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "ID-12345",
								"identificationDocNumberUploadID": "ID-UPLOAD-12345",
							},
							"registrationNumber":          "REG-12345",
							"KRAPIN":                      "KRA-12345",
							"KRAPINUploadID":              "KRA-UPLOAD-12345",
							"practiceLicenseUploadID":     "PRA-UPLOAD-12345",
							"practiceLicenseID":           "PRA-12345",
							"supportingDocumentsUploadID": []string{"SUPP-UPLOAD-123456"},
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request - wrong input",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "ID-12345",
								"identificationDocNumberUploadID": "ID-12345",
							},
							"registrationNumber":          12345,
							"KRAPIN":                      "KRA-12345",
							"KRAPINUploadID":              "KRA-UPLOAD-12345",
							"practiceLicenseUploadID":     "PRA-UPLOAD-12345",
							"practiceLicenseID":           "PRA-12345",
							"supportingDocumentsUploadID": []string{"SUPP-UPLOAD-123456"},
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
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got error: %w", err)
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestAddOrganizationPharmaceuticalKYC(t *testing.T) {
	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphQLMutation := `
	mutation addOrganizationPharmaceuticalKYC($input:OrganizationPharmaceuticalInput!){
		addOrganizationPharmaceuticalKYC(input:$input) {
		  organizationTypeName
		  certificateOfIncorporation
		  certificateOfInCorporationUploadID
		  directorIdentifications{
			identificationDocType
			identificationDocNumber
			identificationDocNumberUploadID
		  }
		  organizationCertificate
		  KRAPIN
		  KRAPINUploadID
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
			name: "valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"directorIdentifications": []map[string]interface{}{
								{
									"identificationDocType":           "NATIONALID",
									"identificationDocNumber":         "ID-12345678",
									"identificationDocNumberUploadID": "ID-UPLOAD-12345678",
								},
							},
							"organizationTypeName":               "LIMITED_COMPANY",
							"certificateOfIncorporation":         "CERT-12345",
							"certificateOfInCorporationUploadID": "CERT-UPLOAD-12345",
							"organizationCertificate":            "ORG-12345",
							"KRAPIN":                             "KRA-12345",
							"KRAPINUploadID":                     "KRA-UPLOAD-12345",
							"registrationNumber":                 "REG-12345678",
							"practiceLicenseID":                  "PRA-12345678",
							"practiceLicenseUploadID":            "PRA-UPLOAD-12345678",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request - wrong input type",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"directorIdentifications": []map[string]interface{}{
								{
									"identificationDocType":           "NATIONALID",
									"identificationDocNumber":         "ID-12345678",
									"identificationDocNumberUploadID": "ID-UPLOAD-12345678",
								},
							},
							"organizationTypeName":               "LIMITED_COMPANY",
							"certificateOfIncorporation":         12345,
							"certificateOfInCorporationUploadID": 12345,
							"organizationCertificate":            12345,
							"KRAPIN":                             "KRA-12345",
							"KRAPINUploadID":                     "KRA-UPLOAD-12345",
							"registrationNumber":                 "REG-12345678",
							"practiceLicenseID":                  "PRA-12345678",
							"practiceLicenseUploadID":            "PRA-UPLOAD-12345678",
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
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got error: %w", err)
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestAddIndividualCoachKYC(t *testing.T) {
	// create a user and their profile
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphQLMutation := `
	mutation addIndividualCoachKYC($input:IndividualCoachInput!){
		addIndividualCoachKYC(input:$input) {
		  identificationDoc {
			identificationDocType
			identificationDocNumber
			identificationDocNumberUploadID
	  }
		  KRAPIN
		  KRAPINUploadID
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
			name: "valid mutation request",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "ID-12345678",
								"identificationDocNumberUploadID": "ID-UPLOAD-12345678",
							},
							"KRAPIN":            "KRA-12345678",
							"KRAPINUploadID":    "KRA-UPLOAD-12345678",
							"practiceLicenseID": "PRA-12345678",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "invalid mutation request - wrong input type",
			args: args{
				query: map[string]interface{}{
					"query": graphQLMutation,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "ID-12345678",
								"identificationDocNumberUploadID": "ID-UPLOAD-12345678",
							},
							"KRAPIN":            12345678,
							"KRAPINUploadID":    12345678,
							"practiceLicenseID": "PRA-12345678",
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
				_, ok := data["errors"]
				if !ok {
					t.Errorf("expected an error")
					return
				}
			}

			if !tt.wantErr {
				_, ok := data["errors"]
				if ok {
					t.Errorf("error not expected got error: %w", err)
					return
				}
			}
			if tt.wantStatus != resp.StatusCode {
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}
}

func TestAddIndividualRiderKYC_acceptance(t *testing.T) {
	user, err := CreateTestUserByPhone(t)
	if err != nil {
		log.Printf("unable to create a test user: %s", err)
		return
	}

	idToken := user.Auth.IDToken
	headers, err := CreatedUserGraphQLHeaders(idToken)
	if err != nil {
		t.Errorf("error in getting headers: %w", err)
		return
	}

	graphQLURL := fmt.Sprintf("%s/%s", baseURL, "graphql")

	graphqlMutationPayload := `
	mutation addIndividualRiderKYC($input: IndividualRiderInput!){
		addIndividualRiderKYC(input:$input){
		  identificationDoc{
			identificationDocType
			identificationDocNumber
			identificationDocNumberUploadID
		  }
		  KRAPIN
		  KRAPINUploadID
		  drivingLicenseID
		  drivingLicenseUploadID
		  certificateGoodConductUploadID
		  supportingDocumentsUploadID
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
			name: "Happy Case - Successfully Add individual rider kyc",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutationPayload,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "NATIONALID",
								"identificationDocNumber":         "12345678",
								"identificationDocNumberUploadID": "12345678",
							},
							"KRAPIN":                         "12345678",
							"KRAPINUploadID":                 "12345678",
							"drivingLicenseID":               "12345678",
							"certificateGoodConductUploadID": "123456",
						},
					},
				},
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "Sad Case - Add individual rider kyc using invalid payload",
			args: args{
				query: map[string]interface{}{
					"query": graphqlMutationPayload,
					"variables": map[string]interface{}{
						"input": map[string]interface{}{
							"identificationDoc": map[string]interface{}{
								"identificationDocType":           "PASSPORT",
								"identificationDocNumber":         "12345678",
								"identificationDocNumberUploadID": "12345678",
							},
							"KRAPIN":                         123456789,
							"KRAPINUploadID":                 123456789,
							"drivingLicenseID":               "678910",
							"certificateGoodConductUploadID": "3458139",
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
				t.Errorf("Bad status reponse returned")
				return
			}

		})
	}

}
