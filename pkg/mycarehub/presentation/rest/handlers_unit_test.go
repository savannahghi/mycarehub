package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	restMock "github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/rest/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	appointmentMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/appointments/mock"
	authorityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority/mock"
	communitiesMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities/mock"
	contentMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content/mock"
	facilityMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility/mock"
	feedbackMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback/mock"
	healthdiaryMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary/mock"
	metricsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/metrics/mock"
	notificationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
	oauthMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/oauth/mock"
	organisationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation/mock"
	otpMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp/mock"
	programsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/programs/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/pubsub/mock"
	questionnairesMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires/mock"
	securityquestionsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions/mock"
	servicerequestMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest/mock"
	surveysMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/surveys/mock"
	termsMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms/mock"
	userMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user/mock"
)

func mapToJSONReader(m map[string]interface{}) (io.Reader, error) {
	bs, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal map to JSON: %w", err)
	}

	buf := bytes.NewBuffer(bs)
	return buf, nil
}

func TestUnit_GetUserRespondedSecurityQuestions(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		method             string
		body               map[string]interface{}
		url                string
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name:   "Happy case: get responded security questions",
			ctx:    context.Background(),
			method: "POST",
			body: map[string]interface{}{
				"username": "test",
				"flavour":  feedlib.FlavourConsumer,
				"otp":      "1234",
			},

			url:                "/get_user_responded_security_questions",
			expectedStatusCode: http.StatusOK,
			wantErr:            false,
		},

		{
			name:   "Sad case: missing input, username",
			ctx:    context.Background(),
			method: "POST",
			body: map[string]interface{}{
				"flavour": feedlib.FlavourConsumer,
				"otp":     "1234",
			},

			url:                "/get_user_responded_security_questions",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},

		{
			name:   "Sad case: failed to get user responded security questions",
			ctx:    context.Background(),
			method: "POST",
			body: map[string]interface{}{
				"username": "test",
				"flavour":  feedlib.FlavourConsumer,
				"otp":      "1234",
			},

			url:                "/get_user_responded_security_questions",
			expectedStatusCode: http.StatusBadRequest,
			wantErr:            true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			facilityUseCase := facilityMock.NewFacilityUsecaseMock()
			notificationUseCase := notificationMock.NewServiceNotificationMock()
			authorityUseCase := authorityMock.NewAuthorityUseCaseMock()
			userUsecase := userMock.NewUserUseCaseMock()
			termsUsecase := termsMock.NewTermsUseCaseMock()
			securityQuestionsUsecase := securityquestionsMock.NewSecurityQuestionsUseCaseMock()
			contentUseCase := contentMock.NewContentUsecaseMock()
			feedbackUsecase := feedbackMock.NewFeedbackUsecaseMock()
			serviceRequestUseCase := servicerequestMock.NewServiceRequestUseCaseMock()
			appointmentUsecase := appointmentMock.NewAppointmentsUseCaseMock()
			healthDiaryUseCase := healthdiaryMock.NewHealthDiaryUseCaseMock()
			surveysUsecase := surveysMock.NewSurveysMock()
			metricsUsecase := metricsMock.NewMetricsUseCaseMock()
			questionnaireUsecase := questionnairesMock.NewServiceRequestUseCaseMock()
			programsUsecase := programsMock.NewProgramsUseCaseMock()
			organisationUsecase := organisationMock.NewOrganisationUseCaseMock()
			otpUseCase := otpMock.NewOTPUseCaseMock()
			pubSubUseCase := pubsubMock.NewServicePubSubMock()
			communityUsecase := communitiesMock.NewCommunityUsecaseMock()
			oauthUsecases := oauthMock.NewOauthUseCaseMock()
			fakeUsecases := usecases.NewMyCareHubUseCase(
				userUsecase, termsUsecase, facilityUseCase,
				securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
				serviceRequestUseCase, authorityUseCase,
				appointmentUsecase, notificationUseCase, surveysUsecase, metricsUsecase, questionnaireUsecase,
				programsUsecase,
				organisationUsecase, pubSubUseCase, communityUsecase, oauthUsecases,
			)
			sessionManager := restMock.NewSCSSessionManagerMock()
			provider := restMock.NewFositeOAuth2Mock()

			if tt.name == "Sad case: failed to get user responded security questions" {
				securityQuestionsUsecase.MockGetUserRespondedSecurityQuestionsFn = func(ctx context.Context, input dto.GetUserRespondedSecurityQuestionsInput) ([]*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}

			ts := httptest.NewServer(h.GetUserRespondedSecurityQuestions())
			defer ts.Close()

			body, err := mapToJSONReader(tt.body)
			if err != nil {
				t.Errorf("invalid response body: %v", err)
			}

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, body)
			if err != nil {
				t.Fatal(err)
			}

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

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

			if !tt.wantErr {
				_, ok := data["error"]
				if ok {
					t.Errorf("error not expected, got %v", data["error"])
					return
				}
			}

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}

}
