package surveys

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	xj "github.com/basgys/goxml2json"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/serverutils"
)

var (
	surveysSystemEmail    = serverutils.MustGetEnvVar("SURVEYS_SYSTEM_EMAIL")
	surveysSystemPassword = serverutils.MustGetEnvVar("SURVEYS_SYSTEM_PASSWORD")
)

// ODKClient defines the fields required to access the surveys client
type ODKClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// RequestHelperPayload is the payload that is sent to the surveys client
type RequestHelperPayload struct {
	Method string
	Path   string
	Body   interface{}
}

// Surveys is the interface that defines the methods that are required to access the surveys client
type Surveys interface {
	ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error)

	GeneratePublicAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error)
	DeletePublicAccessLink(ctx context.Context, input dto.VerifySurveySubmissionInput) error
	ListPublicAccessLinks(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error)

	GetSubmissions(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error)
	ListSubmitters(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error)

	GetFormXML(ctx context.Context, projectID int, formID, version string) (map[string]interface{}, error)
	GetSubmissionXML(ctx context.Context, projectID int, formID, instanceID string) (map[string]interface{}, error)
}

// Impl implements the Surveys interface
type Impl struct {
	client ODKClient
}

// NewSurveysImpl returns a new Impl
func NewSurveysImpl(client ODKClient) Surveys {
	return &Impl{
		client: client,
	}
}

// MakeRequest performs a http request and returns a response
func (s *Impl) MakeRequest(ctx context.Context, payload RequestHelperPayload) (*http.Response, error) {
	client := s.client.HTTPClient

	// A GET or DELETE request should not send data when doing a request. We should use query parameters
	// instead of having a request body. In some cases where a GET request has an empty body {},
	// it might result in status code 400 with the error:
	//  `Your client has issued a malformed or illegal request. That’s all we know.`
	if payload.Method == http.MethodGet || payload.Method == http.MethodDelete {
		req, err := http.NewRequestWithContext(ctx, payload.Method, payload.Path, nil)
		if err != nil {
			return nil, err
		}

		req.SetBasicAuth(surveysSystemEmail, surveysSystemPassword)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Extended-Metadata", "true")
		return client.Do(req)
	}

	encoded, err := json.Marshal(payload.Body)
	if err != nil {
		return nil, err
	}

	p := bytes.NewBuffer(encoded)
	req, err := http.NewRequestWithContext(ctx, payload.Method, payload.Path, p)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(surveysSystemEmail, surveysSystemPassword)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

// ListSurveyForms returns a list of survey forms
func (s *Impl) ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {

	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s/v1/projects/%s/forms", s.client.BaseURL, strconv.Itoa(projectID)),
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
			XMLFormID: surveyForm.XMLFormID,
			Name:      surveyForm.Name,
			EnketoID:  surveyForm.EnketoID,
			Version:   surveyForm.Version,
		})
	}

	return surveyFormsDomain, nil
}

// GetSurveyForm returns a survey form
func (s *Impl) GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {

	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   fmt.Sprintf("%s/v1/projects/%s/forms/%s", s.client.BaseURL, strconv.Itoa(projectID), formID),
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("surveys: GetSurveyForm error: status code: %d", resp.StatusCode)
	}

	var surveyForm dto.SurveyForm
	err = json.NewDecoder(resp.Body).Decode(&surveyForm)

	if err != nil {
		return nil, err
	}

	return &domain.SurveyForm{
		ProjectID: surveyForm.ProjectID,
		XMLFormID: surveyForm.XMLFormID,
		Name:      surveyForm.Name,
		EnketoID:  surveyForm.EnketoID,
	}, nil
}

// GeneratePublicAccessLink returns a survey public link
func (s *Impl) GeneratePublicAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
	payload := RequestHelperPayload{
		Method: http.MethodPost,
		Path:   fmt.Sprintf("%s/v1/projects/%s/forms/%s/public-links", s.client.BaseURL, strconv.Itoa(input.ProjectID), input.FormID),
		Body:   map[string]interface{}{"once": input.OnceOnly, "displayName": input.DisplayName},
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("surveys: GeneratePublicAccessLink error: status code: %d", resp.StatusCode)
	}

	var surveyPublicLink dto.SurveyPublicLink
	err = json.NewDecoder(resp.Body).Decode(&surveyPublicLink)
	if err != nil {
		return nil, err
	}

	return &surveyPublicLink, nil
}

// GetSubmissions returns a list of all survey submissions
func (s *Impl) GetSubmissions(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/submissions", s.client.BaseURL, input.ProjectID, input.FormID)
	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}
	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var submissions []domain.Submission
	err = json.Unmarshal(body, &submissions)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal submissions: %w", err)
	}

	return submissions, nil
}

// DeletePublicAccessLink deletes the survey public link
func (s *Impl) DeletePublicAccessLink(ctx context.Context, input dto.VerifySurveySubmissionInput) error {
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/public-links/%v", s.client.BaseURL, input.ProjectID, input.FormID, input.SubmitterID)
	payload := RequestHelperPayload{
		Method: http.MethodDelete,
		Path:   url,
	}
	_, err := s.MakeRequest(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return err
	}

	return nil
}

// ListSubmitters returns a a listing of all known submitting actors to a given Form. Each Actor that has submitted to the given Form will be returned once.
func (s *Impl) ListSubmitters(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/submissions/submitters", s.client.BaseURL, projectID, formID)
	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}
	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var submitters []domain.Submitter
	err = json.Unmarshal(body, &submitters)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal submitters: %w", err)
	}

	return submitters, nil
}

// ListPublicAccessLinks returns a list of all public access links created for a particular form
func (s *Impl) ListPublicAccessLinks(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error) {
	// {ODK Base URL}/v1/projects/{projectId}/forms/{xmlFormId}/public-links
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/public-links", s.client.BaseURL, projectID, formID)
	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}
	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	var links []*dto.SurveyPublicLink
	err = json.Unmarshal(body, &links)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("unable to unmarshal submitters: %w", err)
	}

	return links, nil
}

// GetSubmissionXML retrieves a submission's XML definition using the instance id
func (s *Impl) GetSubmissionXML(ctx context.Context, projectID int, formID, instanceID string) (map[string]interface{}, error) {
	// {ODK Base URL}/v1/projects/projectId/forms/xmlFormId/submissions/instanceId.xml
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/submissions/%s.xml", s.client.BaseURL, projectID, formID, instanceID)
	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getSubmissionXML: invalid http response, got: %s", resp.Status)
	}

	parsed, err := xj.Convert(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to convert submission xml: %w", err)
	}

	submission := make(map[string]interface{})
	if err := json.Unmarshal(parsed.Bytes(), &submission); err != nil {
		return nil, fmt.Errorf("failed to unmarshal submission: %w", err)
	}

	return submission, nil
}

// GetFormXML retrieves a form's XML definition
func (s *Impl) GetFormXML(ctx context.Context, projectID int, formID, version string) (map[string]interface{}, error) {
	// {ODK Base URL}/v1/projects/projectId/forms/xmlFormId/versions/version.xml
	url := fmt.Sprintf("%s/v1/projects/%v/forms/%s/versions/%s.xml", s.client.BaseURL, projectID, formID, version)
	payload := RequestHelperPayload{
		Method: http.MethodGet,
		Path:   url,
	}

	resp, err := s.MakeRequest(ctx, payload)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			helpers.ReportErrorToSentry(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("getFormXML: invalid http response, got: %s", resp.Status)
	}

	parsed, err := xj.Convert(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to convert form xml: %w", err)
	}

	form := make(map[string]interface{})
	if err := json.Unmarshal(parsed.Bytes(), &form); err != nil {
		return nil, fmt.Errorf("failed to unmarshal form: %w", err)
	}

	return form, nil
}
