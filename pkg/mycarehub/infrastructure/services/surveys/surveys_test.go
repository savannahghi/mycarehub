package surveys_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/jarcoal/httpmock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
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

	surveysClient := domain.SurveysClient{
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
