package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/scalarutils"
)

const kenyaEMRSubRoute = "kenya-emr"

func Test_RegisterKenyaEMRPatients(t *testing.T) {
	ctx := context.Background()
	registerPatientURL := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "register_patient")

	token, err := GetBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	type args struct {
		input *dto.PatientsPayload
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success: register patient",
			args: args{
				input: &dto.PatientsPayload{
					Patients: []*dto.PatientRegistrationPayload{
						{
							MFLCode:   strconv.Itoa(mflCode),
							CCCNumber: cccNumber,
							Name:      gofakeit.BeerName(),
							DateOfBirth: scalarutils.Date{
								Year:  2000,
								Month: 10,
								Day:   10,
							},
							ClientType:  enums.ClientTypePmtct,
							PhoneNumber: "+254888888888",
							EnrollmentDate: scalarutils.Date{
								Year:  2000,
								Month: 10,
								Day:   20,
							},
							BirthDateEstimated: false,
							Gender:             enumutils.GenderFemale.String(),
							Counselled:         true,
							NextOfKin: dto.NextOfKinPayload{
								Name:         gofakeit.Name(),
								Relationship: enums.CaregiverTypeFather.String(),
							},
						},
					},
				},
			},
		},
		{
			name: "sad case: missing facility code",
			args: args{
				input: &dto.PatientsPayload{
					Patients: []*dto.PatientRegistrationPayload{
						{
							CCCNumber: cccNumber,
							Name:      gofakeit.BeerName(),
							DateOfBirth: scalarutils.Date{
								Year:  2000,
								Month: 10,
								Day:   10,
							},
							ClientType:  enums.ClientTypePmtct,
							PhoneNumber: "+254888888888",
							EnrollmentDate: scalarutils.Date{
								Year:  2000,
								Month: 10,
								Day:   20,
							},
							BirthDateEstimated: false,
							Gender:             enumutils.GenderFemale.String(),
							Counselled:         true,
							NextOfKin: dto.NextOfKinPayload{
								Name:         gofakeit.Name(),
								Relationship: enums.CaregiverTypeFather.String(),
							},
						},
					},
				},
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
			bs, err := json.Marshal(tt.args.input)
			if err != nil {
				t.Errorf("unable to marshal test item to JSON: %s", err)
			}
			payload := bytes.NewBuffer(bs)

			r, err := http.NewRequest(
				http.MethodPost,
				registerPatientURL,
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
