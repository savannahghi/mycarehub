package acceptancetests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
)

const kenyaEMRSubRoute = "kenya-emr"

func Test_RegisterPatient(t *testing.T) {
	ctx := context.Background()
	registerPatientURL := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "register_patient")

	fmt.Println("db==", os.Getenv("POSTGRES_DB"))

	token, err := GetBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	registerPatientParamsMap := map[string]interface{}{
		"patients": []map[string]interface{}{
			{
				"MFLCODE":            strconv.Itoa(mflCode),
				"cccNumber":          "10001",
				"name":               gofakeit.Name(),
				"dateOfBirth":        "1999-01-01",
				"clientType":         enums.ClientTypeHvl.String(),
				"phoneNumber":        "+254888888888",
				"enrollmentDate":     "2000-01-01",
				"birthDateEstimated": false,
				"gender":             enumutils.GenderFemale.String(),
				"counselled":         true,
				"nextOfKin": map[string]interface{}{
					"name":         gofakeit.Name(),
					"relationship": enums.CaregiverTypeHealthCareProfessional.String(),
				},
			},
		},
	}
	registerPatientParamsMissingFacilityMap := map[string]interface{}{
		"patients": []map[string]interface{}{
			{
				"cccNumber":          "10002",
				"name":               gofakeit.Name(),
				"dateOfBirth":        "1999-01-01",
				"clientType":         enums.ClientTypeHvl.String(),
				"phoneNumber":        "+254888888881",
				"enrollmentDate":     "2000-01-01",
				"birthDateEstimated": false,
				"gender":             enumutils.GenderFemale.String(),
				"counselled":         true,
				"nextOfKin": map[string]interface{}{
					"name":         gofakeit.Name(),
					"relationship": enums.CaregiverTypeHealthCareProfessional.String(),
				},
			},
		},
	}

	type args struct {
		input map[string]interface{}
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				input: registerPatientParamsMap,
			},
		},
		{
			name: "sad case: missing facility",
			args: args{
				input: registerPatientParamsMissingFacilityMap,
			},
			wantErr: true,
		},
		{
			name:    "sad case: missing input",
			args:    args{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := mapToJSONReader(tt.args.input)
			if err != nil {
				t.Errorf("unable to get JSON io Reader: %s", err)
				return
			}

			r, err := http.NewRequest(
				http.MethodPost,
				registerPatientURL,
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

			r.Header.Add("Authorization", token)
			r.Header.Add("Content-Type", "application/json")

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
				data := []*dto.PatientRegistrationPayload{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}

}
