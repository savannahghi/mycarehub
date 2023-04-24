package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

const kenyaEMRSubRoute = "kenya-emr"

func Test_GetClientHealthDiaryEntries(t *testing.T) {
	ctx := context.Background()
	getKenyaEMRHealthDairyEntries := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "health_diary")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	syncTime := time.Now().UTC()

	type args struct {
		input *dto.FetchHealthDiaryEntries
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success: get health diary entries",
			args: args{
				input: &dto.FetchHealthDiaryEntries{
					MFLCode:      mflCode,
					LastSyncTime: &syncTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get health diary entries; invalid mfl code",
			args: args{
				input: &dto.FetchHealthDiaryEntries{
					MFLCode:      -1,
					LastSyncTime: &syncTime,
				},
			},
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
				http.MethodGet,
				getKenyaEMRHealthDairyEntries,
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
				data := &dto.HealthDiaryEntriesResponse{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}

}

func Test_GetServiceRequestForKenyaEMR(t *testing.T) {
	ctx := context.Background()
	getKenyaEMRServiceReq := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "service_request")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	syncTime := time.Now().UTC()

	mflcode, err := strconv.Atoi(mflIdentifier)
	if err != nil {
		t.Errorf("failed to convert string to int: %v", err)
	}

	type args struct {
		input *dto.FetchHealthDiaryEntries
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success: get service requests",
			args: args{
				input: &dto.FetchHealthDiaryEntries{
					MFLCode:      mflcode,
					LastSyncTime: &syncTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get service requests; invalid mfl code",
			args: args{
				input: &dto.FetchHealthDiaryEntries{
					MFLCode:      -1,
					LastSyncTime: &syncTime,
				},
			},
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
				http.MethodGet,
				getKenyaEMRServiceReq,
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
				data := &dto.RedFlagServiceRequestResponse{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}

}

func Test_UpdateServiceRequestfromKenyaEMR(t *testing.T) {
	ctx := context.Background()
	updateKenyaEMRServiceReq := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "service_request")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	type args struct {
		input *dto.UpdateServiceRequestsPayload
	}

	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStatus int
	}{
		{
			name: "Success: update service requests",
			args: args{
				input: &dto.UpdateServiceRequestsPayload{
					ServiceRequests: []dto.UpdateServiceRequestPayload{
						{
							ID:           clientsServiceRequestID,
							RequestType:  enums.ServiceRequestTypeAppointments.String(),
							Status:       enums.ServiceRequestStatusInProgress.String(),
							InProgressAt: time.Now().UTC(),
							InProgressBy: "04d892cb-463e-4fb1-92bd-2e8f91295dce",
							ResolvedAt:   time.Now().UTC(),
							ResolvedBy:   "04d892cb-463e-4fb1-92bd-2e8f91295dce",
						},
					},
				},
			},
			wantErr:    false,
			wantStatus: http.StatusOK,
		},
		{
			name: "Fail: unable to update service requests; invalid ID",
			args: args{
				input: &dto.UpdateServiceRequestsPayload{
					ServiceRequests: []dto.UpdateServiceRequestPayload{
						{
							ID:           gofakeit.UUID(),
							RequestType:  enums.ServiceRequestTypeAppointments.String(),
							Status:       enums.ServiceRequestStatusResolved.String(),
							InProgressAt: time.Now().UTC(),
							InProgressBy: "04d892cb-463e-4fb1-92bd-2e8f91295dce",
							ResolvedAt:   time.Now().UTC(),
							ResolvedBy:   "04d892cb-463e-4fb1-92bd-2e8f91295dce",
						},
					},
				},
			},
			wantErr:    true,
			wantStatus: http.StatusInternalServerError,
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
				updateKenyaEMRServiceReq,
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
				data := &domain.UpdateServiceRequestsPayload{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}

}

func Test_GetRegisteredFacilityPatientsForKenyaEMR(t *testing.T) {
	ctx := context.Background()
	getRegisteredPatientsURL := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "patients")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	syncTime := time.Now().UTC()

	type args struct {
		input *dto.PatientSyncPayload
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success: update service requests",
			args: args{
				input: &dto.PatientSyncPayload{
					MFLCode:  mflCode,
					SyncTime: &syncTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update service requests; invalid mfl code",
			args: args{
				input: &dto.PatientSyncPayload{
					MFLCode:  -1,
					SyncTime: &syncTime,
				},
			},
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
				http.MethodGet,
				getRegisteredPatientsURL,
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
				data := &dto.PatientSyncResponse{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}
}

func Test_GetAppointmentsServiceRequests(t *testing.T) {
	ctx := context.Background()
	getAppointmentURL := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "appointment-service-request")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	syncTime := time.Now().UTC()

	type args struct {
		input *dto.AppointmentServiceRequestInput
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Success: update service requests",
			args: args{
				input: &dto.AppointmentServiceRequestInput{
					MFLCode:      mflCode,
					LastSyncTime: &syncTime,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to update service requests; invalid mfl code",
			args: args{
				input: &dto.AppointmentServiceRequestInput{
					MFLCode:      -1,
					LastSyncTime: &syncTime,
				},
			},
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
				http.MethodGet,
				getAppointmentURL,
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
				data := &dto.AppointmentServiceRequestsOutput{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}
}

func Test_UpdateAppointmentsServiceRequests(t *testing.T) {
	ctx := context.Background()
	updateAppointmentURL := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "appointment-service-request")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	type args struct {
		input *dto.UpdateServiceRequestsPayload
	}

	tests := []struct {
		name       string
		args       args
		wantStatus int
		wantErr    bool
	}{
		{
			name: "Success: update appointment service requests",
			args: args{
				input: &dto.UpdateServiceRequestsPayload{
					ServiceRequests: []dto.UpdateServiceRequestPayload{
						{
							ID:           clientsServiceRequestID,
							RequestType:  "PIN_RESET",
							Status:       enums.ServiceRequestStatusInProgress.String(),
							InProgressAt: time.Now(),
							InProgressBy: "04d892cb-463e-4fb1-92bd-2e8f91295dce",
							ResolvedAt:   time.Now(),
							ResolvedBy:   "04d892cb-463e-4fb1-92bd-2e8f91295dce",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail: unable to update appointment service requests; invalid ID",
			args: args{
				input: &dto.UpdateServiceRequestsPayload{
					ServiceRequests: []dto.UpdateServiceRequestPayload{
						{
							ID:           gofakeit.BS(),
							RequestType:  "PIN_RESET",
							Status:       enums.ServiceRequestStatusInProgress.String(),
							InProgressAt: time.Now(),
							InProgressBy: "04d892cb-463e-4fb1-92bd-2e8f91295dce",
							ResolvedAt:   time.Now(),
							ResolvedBy:   "04d892cb-463e-4fb1-92bd-2e8f91295dce",
						},
					},
				},
			},
			wantErr: false,
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
				updateAppointmentURL,
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

		})
	}
}

func Test_CreateOrUpdateKenyaEMRAppointments(t *testing.T) {
	ctx := context.Background()
	createOrUpdateAppointmentURL := fmt.Sprintf("%s/%s/%s", baseURL, kenyaEMRSubRoute, "appointments")

	token, err := GetFirebaseBearerTokenHeader(ctx)
	if err != nil {
		t.Errorf("failed to get bearer token header: %v", err)
		return
	}

	type args struct {
		input *dto.FacilityAppointmentsPayload
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success: create or update service requests",
			args: args{
				input: &dto.FacilityAppointmentsPayload{
					MFLCode: strconv.Itoa(mflCode),
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:  cccNumber,
							ExternalID: "5",
							AppointmentDate: scalarutils.Date{
								Year:  2022,
								Month: 10,
								Day:   10,
							},
							AppointmentReason: "Pharmacy Visit",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail: unable to create or update service requests; invalid mfl code",
			args: args{
				input: &dto.FacilityAppointmentsPayload{
					MFLCode: strconv.Itoa(-11),
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:  gofakeit.BS(),
							ExternalID: "5",
							AppointmentDate: scalarutils.Date{
								Year:  2022,
								Month: 10,
								Day:   10,
							},
							AppointmentReason: "Pharmacy Visit",
						},
					},
				},
			},
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
				createOrUpdateAppointmentURL,
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
				data := &dto.FacilityAppointmentsResponse{}
				err = json.Unmarshal(dataResponse, &data)
				if err != nil {
					t.Errorf("bad data returned")
					return
				}
			}

		})
	}
}
