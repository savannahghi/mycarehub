package surveys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

var (
	surveysSystemEmail    = serverutils.MustGetEnvVar("SURVEYS_SYSTEM_EMAIL")
	surveysSystemPassword = serverutils.MustGetEnvVar("SURVEYS_SYSTEM_PASSWORD")
)

// Surveys is the interface that defines the methods that are required to access the surveys client
type Surveys interface {
	ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	MakeRequest(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error)
}

// Impl implements the Surveys interface
type Impl struct {
	client domain.SurveysClient
}

// NewSurveysImpl returns a new Impl
func NewSurveysImpl(client domain.SurveysClient) Surveys {
	return &Impl{
		client: client,
	}
}

// ListSurveyForms returns a list of survey forms
func (s *Impl) ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {

	payload := domain.RequestHelperPayload{
		Method: http.MethodGet,
		Path:   s.client.BaseURL + "/projects/" + strconv.Itoa(projectID) + "/forms",
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("surveys: ListSurveyForms error: status code: %d", resp.StatusCode)
	}

	var surveyForms []*dto.SurveyForm
	err = json.NewDecoder(resp.Body).Decode(&surveyForms)
	if err != nil {
		return nil, err
	}

	var surveyFormsDomain []*domain.SurveyForm
	for _, surveyForm := range surveyForms {
		surveyFormsDomain = append(surveyFormsDomain, &domain.SurveyForm{
			ProjectID: surveyForm.ProjectID,
			Name:      surveyForm.Name,
		})
	}

	return surveyFormsDomain, nil
}

// MakeRequest performs a http request and returns a response
func (s *Impl) MakeRequest(ctx context.Context, payload domain.RequestHelperPayload) (*http.Response, error) {
	client := s.client.HTTPClient

	// A GET or DELETE request should not send data when doing a request. We should use query parameters
	// instead of having a request body. In some cases where a GET request has an empty body {},
	// it might result in status code 400 with the error:
	//  `Your client has issued a malformed or illegal request. That’s all we know.`
	if payload.Method == http.MethodGet || payload.Method == http.MethodDelete {
		req, reqErr := http.NewRequestWithContext(ctx, payload.Method, payload.Path, nil)
		if reqErr != nil {
			return nil, reqErr
		}

		req.SetBasicAuth(surveysSystemEmail, surveysSystemPassword)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		return client.Do(req)
	}

	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return nil, err
	}

	p := bytes.NewBuffer(encoded)
	req, reqErr := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if reqErr != nil {
		return nil, reqErr
	}

	req.SetBasicAuth(surveysSystemEmail, surveysSystemPassword)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}
