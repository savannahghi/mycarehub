package surveys_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
)

func TestSurveysImpl_ListSurveyForms(t *testing.T) {
	projectID := 2
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/v1/projects/2/forms",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
				{
					"projectId":    2,
					"xmlFormId":    gofakeit.UUID(),
					"state":        gofakeit.UUID(),
					"enketoId":     gofakeit.UUID(),
					"enketoOnceId": gofakeit.UUID(),
					"createdAt":    "2022-04-28T08:19:17.473Z",
					"updatedAt":    "2022-04-28T08:20:43.591Z",
					"keyId":        nil,
					"version":      "1 (2022-04-28 08:18:20)",
					"hash":         gofakeit.UUID(),
					"sha":          gofakeit.UUID(),
					"sha256":       gofakeit.UUID(),
					"draftToken":   nil,
					"publishedAt":  "2022-04-28T08:20:27.643Z",
					"name":         "test",
				},
			})
			return resp, err
		},
	)

	surveysClient := surveys.ODKClient{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}
	s := surveys.NewSurveysImpl(surveysClient)
	type args struct {
		ctx       context.Context
		projectID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: ListSurveyForms successful",
			args: args{
				ctx:       context.Background(),
				projectID: projectID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListSurveyForms(tt.args.ctx, tt.args.projectID)
			if (err != nil) != tt.wantErr {
				t.Errorf("SurveysImpl.ListSurveyForms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("SurveysImpl.ListSurveyForms() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestImpl_GetSubmissions(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/v1/projects/2/forms/akmCQQxf4LaFjAWDbg29pj (1)/submissions",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
				{},
			})
			return resp, err
		},
	)

	surveysClient := surveys.ODKClient{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}
	s := surveys.NewSurveysImpl(surveysClient)

	type args struct {
		ctx   context.Context
		input dto.VerifySurveySubmissionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: GetSubmissions successful",
			args: args{
				ctx: context.Background(),
				input: dto.VerifySurveySubmissionInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 86,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetSubmissions(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetSubmissions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SurveysImpl.ListSurveyForms() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestImpl_DeletePublicAccessLink(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("DELETE", "/v1/projects/2/forms/akmCQQxf4LaFjAWDbg29pj (1)/public-links/86",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
				{},
			})
			return resp, err
		},
	)

	surveysClient := surveys.ODKClient{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}
	s := surveys.NewSurveysImpl(surveysClient)

	type args struct {
		ctx   context.Context
		input dto.VerifySurveySubmissionInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: GetSubmissions successful",
			args: args{
				ctx: context.Background(),
				input: dto.VerifySurveySubmissionInput{
					ProjectID:   2,
					FormID:      "akmCQQxf4LaFjAWDbg29pj (1)",
					SubmitterID: 86,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.DeletePublicAccessLink(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("Impl.DeletePublicAccessLink() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestImpl_ListSubmitters(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/v1/projects/2/forms/akmCQQxf4LaFjAWDbg29pj (1)/submissions/submitters",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
				{},
			})
			return resp, err
		},
	)

	surveysClient := surveys.ODKClient{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}
	s := surveys.NewSurveysImpl(surveysClient)

	type args struct {
		ctx       context.Context
		projectID int
		formID    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: ListSubmitters successful",
			args: args{
				ctx:       context.Background(),
				projectID: 2,
				formID:    "akmCQQxf4LaFjAWDbg29pj (1)",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ListSubmitters(tt.args.ctx, tt.args.projectID, tt.args.formID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.ListSubmitters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SurveysImpl.ListSubmitters() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestImpl_GeneratePublicAccessLink(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "/v1/projects/2/forms/akmCQQxf4LaFjAWDbg29pj (1)/public-links",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
				{},
			})
			return resp, err
		},
	)

	surveysClient := surveys.ODKClient{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}
	s := surveys.NewSurveysImpl(surveysClient)

	type args struct {
		ctx   context.Context
		input dto.SurveyLinkInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Sad case",
			args: args{
				ctx: context.Background(),
				input: dto.SurveyLinkInput{
					ProjectID: 2,
					FormID:    "akmCQQxf4LaFjAWDbg29pj (1)",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GeneratePublicAccessLink(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GeneratePublicAccessLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SurveysImpl.GeneratePublicAccessLink() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestImpl_GetSurveyForm(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", "/v1/projects/2/forms/akmCQQxf4LaFjAWDbg29pj (1)",
		func(req *http.Request) (*http.Response, error) {
			resp, err := httpmock.NewJsonResponse(200, []map[string]interface{}{
				{},
			})
			return resp, err
		},
	)

	surveysClient := surveys.ODKClient{
		BaseURL:    "https://example.com",
		HTTPClient: &http.Client{},
	}
	s := surveys.NewSurveysImpl(surveysClient)

	type args struct {
		ctx       context.Context
		projectID int
		formID    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad case: GetSurveyForm successful",
			args: args{
				ctx:       context.Background(),
				projectID: 2,
				formID:    "akmCQQxf4LaFjAWDbg29pj (1)",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetSurveyForm(tt.args.ctx, tt.args.projectID, tt.args.formID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetSurveyForm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("SurveysImpl.GetSurveyForm() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
