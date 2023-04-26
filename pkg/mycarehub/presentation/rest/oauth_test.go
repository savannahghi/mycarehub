package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/ory/fosite"
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

func TestMyCareHubHandlersInterfacesImpl_handleLoginPage(t *testing.T) {
	tests := []struct {
		name                     string
		ctx                      context.Context
		method                   string
		url                      string
		formValues               url.Values
		expectedStatusCode       int
		expectedPageTitle        string
		expectedAvailableProgram []string
	}{
		{
			name:               "Happy case: login and redirect to select program",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/login",
			formValues:         url.Values{"username": {"test"}, "pin": {"1234"}},
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "Programs",
			expectedAvailableProgram: []string{
				"test program",
			},
		},
		{
			name:                     "Sad case: login params not passed when submitting",
			ctx:                      context.Background(),
			method:                   "POST",
			url:                      "/login",
			expectedStatusCode:       http.StatusOK,
			expectedPageTitle:        "Login",
			expectedAvailableProgram: nil,
		},
		{
			name:                     "Sad case: invalid login credentials",
			ctx:                      context.Background(),
			method:                   "POST",
			url:                      "/login",
			formValues:               url.Values{"username": {"test"}, "pin": {"900"}},
			expectedStatusCode:       http.StatusOK,
			expectedPageTitle:        "Login",
			expectedAvailableProgram: nil,
		},

		{
			name:                     "Sad case: failed to get user profile",
			ctx:                      context.Background(),
			method:                   "POST",
			url:                      "/login",
			formValues:               url.Values{"username": {"test"}, "pin": {"900"}},
			expectedStatusCode:       http.StatusOK,
			expectedPageTitle:        "",
			expectedAvailableProgram: nil,
		},

		{
			name:                     "Sad case: failed to list user programs",
			ctx:                      context.Background(),
			method:                   "POST",
			url:                      "/login",
			formValues:               url.Values{"username": {"test"}, "pin": {"900"}},
			expectedStatusCode:       http.StatusOK,
			expectedPageTitle:        "",
			expectedAvailableProgram: nil,
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

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// call the handler
				// ctx := r.Context()
				ar := fosite.NewAuthorizeRequest()
				authorizationSession := &AuthorizationSession{}

				if tt.name == "Sad case: invalid login credentials" {
					userUsecase.MockLoginFn = func(ctx context.Context, input *dto.LoginInput) (*dto.LoginResponse, bool) {
						return &dto.LoginResponse{Message: "invalid credentials"}, false
					}
				}

				if tt.name == "Sad case: failed to get user profile" {
					userUsecase.MockGetUserProfileFn = func(ctx context.Context, userID string) (*domain.User, error) {
						return nil, fmt.Errorf("an error occurred")
					}
				}

				if tt.name == "Sad case: failed to list user programs" {
					programsUsecase.MockListUserProgramsFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error) {
						return nil, fmt.Errorf("an error occurred")
					}
				}

				h.handleLoginPage(w, r, ar, authorizationSession)
			}))
			defer ts.Close()

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, strings.NewReader(tt.formValues.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}

			// parse the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// check the page title
			if !strings.Contains(string(body), tt.expectedPageTitle) {
				t.Errorf("Expected page title '%s', but got '%s'", tt.expectedPageTitle, string(body))
			}

			// check the available programs
			for _, expectedProgram := range tt.expectedAvailableProgram {
				if !strings.Contains(string(body), expectedProgram) {
					t.Errorf("Expected program '%s' to be available, but it was not found in the response body", expectedProgram)
				}
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_handleChooseProgramPage(t *testing.T) {
	tests := []struct {
		name                        string
		ctx                         context.Context
		method                      string
		url                         string
		formValues                  url.Values
		expectedStatusCode          int
		expectedPageTitle           string
		expectedAvailableFacilities []string
	}{
		{
			name:               "Happy case: select program and redirect to select facility",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-program",
			formValues:         url.Values{"program": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "Facilities",
			expectedAvailableFacilities: []string{
				"test facility",
			},
		},
		{
			name:                        "Sad case: program not selected when submitting",
			ctx:                         context.Background(),
			method:                      "POST",
			url:                         "/choose-program",
			expectedStatusCode:          http.StatusOK,
			expectedPageTitle:           "Programs",
			expectedAvailableFacilities: nil,
		},
		{
			name:                        "Sad case: failed to list programs",
			ctx:                         context.Background(),
			method:                      "POST",
			url:                         "/choose-program",
			expectedStatusCode:          http.StatusOK,
			expectedPageTitle:           "",
			expectedAvailableFacilities: nil,
		},
		{
			name:                        "Sad case: failed to get program by id",
			ctx:                         context.Background(),
			method:                      "POST",
			url:                         "/choose-program",
			formValues:                  url.Values{"program": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			expectedStatusCode:          http.StatusOK,
			expectedPageTitle:           "",
			expectedAvailableFacilities: nil,
		},

		{
			name:                        "Sad case: failed to list program facilities",
			ctx:                         context.Background(),
			method:                      "POST",
			url:                         "/choose-program",
			formValues:                  url.Values{"program": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			expectedStatusCode:          http.StatusOK,
			expectedPageTitle:           "",
			expectedAvailableFacilities: nil,
		},
	}
	for _, tt := range tests {
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

		// create a test server
		h := &MyCareHubHandlersInterfacesImpl{
			provider:       provider,
			usecase:        *fakeUsecases,
			sessionManager: sessionManager,
		}
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// call the handler
			// ctx := r.Context()
			ar := fosite.NewAuthorizeRequest()
			uuid := gofakeit.UUID()
			authorizationSession := &AuthorizationSession{
				Page: "chooseProgram",
				User: domain.User{ID: &uuid},
			}

			if tt.name == "Sad case: failed to list programs" {
				programsUsecase.MockListUserProgramsFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get program by id" {
				programsUsecase.MockGetProgramByIDFn = func(ctx context.Context, programID string) (*domain.Program, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to list program facilities" {
				facilityUseCase.MockListProgramFacilitiesFn = func(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			h.handleChooseProgramPage(w, r, ar, authorizationSession)
		}))
		defer ts.Close()

		// create the request
		req, err := http.NewRequest(tt.method, ts.URL+tt.url, strings.NewReader(tt.formValues.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// make the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		// check the response
		if resp.StatusCode != tt.expectedStatusCode {
			t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
		}

		// parse the response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		// check the page title
		if !strings.Contains(string(body), tt.expectedPageTitle) {
			t.Errorf("Expected page title '%s', but got '%s'", tt.expectedPageTitle, string(body))
		}

		// check the available facilities
		for _, expectedFacility := range tt.expectedAvailableFacilities {
			if !strings.Contains(string(body), expectedFacility) {
				t.Errorf("Expected facilities '%s' to be available, but it was not found in the response body", expectedFacility)
			}
		}
	}
}

func TestMyCareHubHandlersInterfacesImpl_handleChooseFacilityPage(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		method             string
		url                string
		formValues         url.Values
		expectedStatusCode int
		expectedPageTitle  string
	}{
		{
			name:               "Happy case: select facility and redirect to the home page",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},
		{
			name:               "Sad case: no facility selected when submitting",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "Facilities",
		},
		{
			name:               "Sad case: failed to retrieve facility",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
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

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// call the handler
				ar := fosite.NewAuthorizeRequest()
				uuid := gofakeit.UUID()
				authorizationSession := &AuthorizationSession{
					Page:    "chooseFacility",
					User:    domain.User{ID: &uuid},
					Program: domain.Program{ID: uuid},
				}
				if tt.name == "Sad case: failed to retrieve facility" {
					facilityUseCase.MockRetrieveFacilityFn = func(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
						return nil, fmt.Errorf("an error occurred")
					}
				}

				h.handleChooseFacilityPage(w, r, ar, authorizationSession)
			}))
			defer ts.Close()

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, strings.NewReader(tt.formValues.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}

			// parse the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// check the page title
			if !strings.Contains(string(body), tt.expectedPageTitle) {
				t.Errorf("Expected page title '%s', but got '%s'", tt.expectedPageTitle, string(body))
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_AuthorizeHandler(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		method             string
		url                string
		formValues         url.Values
		page               string
		expectedStatusCode int
		expectedPageTitle  string
	}{
		{
			name:               "Happy case: login and redirect to select program",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/login",
			formValues:         url.Values{"username": {"test"}, "pin": {"1234"}},
			page:               "login",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "Programs",
		},
		{
			name:               "Happy case: select program and redirect to select facility",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-program",
			formValues:         url.Values{"program": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "chooseProgram",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "Facilities",
		},

		{
			name:               "Happy case: select facility and redirect to homepage",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},

		{
			name:               "Sad case: failed to initialize authorize requester",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},

		{
			name:               "Sad case: failed to set current program",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},

		{
			name:               "Sad case: failed to get staff profile",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},

		{
			name:               "Sad case: failed to set staff default facility",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},

		{
			name:               "Sad case: failed to initialize authorize response",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
		},

		{
			name:               "Sad case: failed to cleanup session",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/choose-facility",
			formValues:         url.Values{"facility": {"5e30451c-9672-4c08-a34c-85b036294362"}},
			page:               "",
			expectedStatusCode: http.StatusOK,
			expectedPageTitle:  "",
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

			// set the session data
			UUID := gofakeit.UUID()
			sessionManager.MockGetBytesFn = func(ctx context.Context, key string) []byte {
				bs, _ := json.Marshal(&AuthorizationSession{
					Page:     tt.page,
					Program:  domain.Program{ID: UUID},
					User:     domain.User{ID: &UUID},
					Facility: domain.Facility{ID: &UUID},
				},
				)
				return bs
			}

			if tt.name == "Happy case: login and redirect to select program" {
				sessionManager.MockExistsFn = func(ctx context.Context, key string) bool {
					return false
				}
			}

			if tt.name == "Sad case: failed to initialize authorize requester" {
				provider.MockNewAuthorizeRequestFn = func(ctx context.Context, req *http.Request) (fosite.AuthorizeRequester, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to set current program" {
				programsUsecase.MockSetCurrentProgramFn = func(ctx context.Context, programID string) (bool, error) {
					return false, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to get staff profile" {
				userUsecase.MockGetStaffProfileFn = func(ctx context.Context, userID, programID string) (*domain.StaffProfile, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to set staff default facility" {
				userUsecase.MockSetStaffDefaultFacilityFn = func(ctx context.Context, staffID, facilityID string) (*domain.Facility, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to initialize authorize response" {
				provider.MockNewAuthorizeResponseFn = func(ctx context.Context, requester fosite.AuthorizeRequester, session fosite.Session) (fosite.AuthorizeResponder, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to cleanup session" {
				sessionManager.MockDestroyFn = func(ctx context.Context) error {
					return errors.New("an error occurred")
				}
			}

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}

			ts := httptest.NewServer(h.AuthorizeHandler())
			defer ts.Close()

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, strings.NewReader(tt.formValues.Encode()))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}

			// parse the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			// check the page title
			if !strings.Contains(string(body), tt.expectedPageTitle) {
				t.Errorf("Expected page title '%s', but got '%s'", tt.expectedPageTitle, string(body))
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_TokenHandler(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Happy case: get token",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/get-token",
			expectedStatusCode: http.StatusOK,
		},

		{
			name:               "Sad case: failed to initialize access request",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/get-token",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Sad case: failed to initialize access response",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/get-token",
			expectedStatusCode: http.StatusOK,
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

			if tt.name == "Sad case: failed to initialize access request" {
				provider.MockNewAccessRequestFn = func(ctx context.Context, req *http.Request, session fosite.Session) (fosite.AccessRequester, error) {
					return nil, errors.New("an error occurred")
				}
			}

			if tt.name == "Sad case: failed to initialize access response" {
				provider.MockNewAccessResponseFn = func(ctx context.Context, requester fosite.AccessRequester) (fosite.AccessResponder, error) {
					return nil, errors.New("an error occurred")
				}
			}

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}

			ts := httptest.NewServer(h.TokenHandler())
			defer ts.Close()

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_RevokeHandler(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Happy case: revoke token",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/revoke-token",
			expectedStatusCode: http.StatusOK,
		},

		{
			name:               "Sad case: failed to initialize revoke request",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/revoke-token",
			expectedStatusCode: http.StatusOK,
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

			if tt.name == "Sad case: failed to initialize revoke request" {
				provider.MockNewRevocationRequestFn = func(ctx context.Context, r *http.Request) error {
					return errors.New("an error occurred")
				}
			}

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}

			ts := httptest.NewServer(h.RevokeHandler())
			defer ts.Close()

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}
}

func TestMyCareHubHandlersInterfacesImpl_IntrospectionHandler(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Happy case: introspect token",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/introspect-token",
			expectedStatusCode: http.StatusOK,
		},

		{
			name:               "Sad case: failed to initialize introspect request",
			ctx:                context.Background(),
			method:             "POST",
			url:                "/introspect-token",
			expectedStatusCode: http.StatusOK,
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

			if tt.name == "Sad case: failed to initialize introspect request" {
				provider.MockNewIntrospectionRequestFn = func(ctx context.Context, r *http.Request, session fosite.Session) (fosite.IntrospectionResponder, error) {
					return nil, errors.New("an error occurred")
				}
			}

			// create a test server
			h := &MyCareHubHandlersInterfacesImpl{
				provider:       provider,
				usecase:        *fakeUsecases,
				sessionManager: sessionManager,
			}

			ts := httptest.NewServer(h.IntrospectionHandler())
			defer ts.Close()

			// create the request
			req, err := http.NewRequest(tt.method, ts.URL+tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// make the request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			// check the response
			if resp.StatusCode != tt.expectedStatusCode {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatusCode, resp.StatusCode)
			}
		})
	}

}
