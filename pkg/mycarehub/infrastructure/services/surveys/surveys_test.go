package surveys_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	xj "github.com/basgys/goxml2json"
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

func TestImpl_GetSubmissionXML(t *testing.T) {
	submission := `<data xmlns:jr="http://openrosa.org/javarosa" xmlns:orx="http://openrosa.org/xforms" id="akmCQQxf4LaFjAWDbg29pj (1)" version="2 (2022-05-11 06:59:42)"><start>2022-08-12T13:11:23.965+03:00</start><end>2022-08-12T13:12:01.646+03:00</end><Over_the_last_2_week_sure_in_doing_things>0</Over_the_last_2_week_sure_in_doing_things><Over_the_last_2_week_epressed_or_hopeless>0</Over_the_last_2_week_epressed_or_hopeless><Over_the_last_2_week_or_sleeping_too_much>0</Over_the_last_2_week_or_sleeping_too_much><Over_the_last_2_week_having_little_energy>0</Over_the_last_2_week_having_little_energy><Over_the_last_2_week_petite_or_overeating>0</Over_the_last_2_week_petite_or_overeating><Over_the_last_2_week_rself_or_family_down>0</Over_the_last_2_week_rself_or_family_down><Over_the_last_2_week_watching_television>0</Over_the_last_2_week_watching_television><Over_the_last_2_week_lot_more_than_usual>0</Over_the_last_2_week_lot_more_than_usual><Over_the_last_2_week_yourself_in_some_way>0</Over_the_last_2_week_yourself_in_some_way><If_you_checked_off_a_ng_with_other_people>not_difficult_at_all</If_you_checked_off_a_ng_with_other_people><__version__>vesMF8UKLW5gnZgBnBmzd9</__version__><meta><instanceID>uuid:808431e7-e2ed-4065-b19d-9fd780ce7f9c</instanceID></meta></data>`

	j, err := xj.Convert(strings.NewReader(submission))
	if err != nil {
		t.Errorf("failed to convert submission")
		return
	}

	parsedSubmission := make(map[string]interface{})
	err = json.Unmarshal(j.Bytes(), &parsedSubmission)
	if err != nil {
		t.Errorf("failed to unmarshal submission")
		return
	}

	type args struct {
		ctx        context.Context
		projectID  int
		formID     string
		instanceID string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "happy case: retrieve submission xml",
			args: args{
				ctx:        context.Background(),
				projectID:  1,
				formID:     "test",
				instanceID: "instance",
			},
			want:    parsedSubmission,
			wantErr: false,
		},
		{
			name: "sad case: bad status code",
			args: args{
				ctx:        context.Background(),
				projectID:  1,
				formID:     "test",
				instanceID: "instance",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: invalid xml",
			args: args{
				ctx:        context.Background(),
				projectID:  1,
				formID:     "test",
				instanceID: "instance",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			surveysClient := surveys.ODKClient{
				BaseURL:    "https://example.com",
				HTTPClient: &http.Client{},
			}
			s := surveys.NewSurveysImpl(surveysClient)

			if tt.name == "happy case: retrieve submission xml" {
				httpmock.RegisterResponder("GET", "/v1/projects/1/forms/test/submissions/instance.xml",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(
							http.StatusOK,
							submission,
						)

						return resp, nil
					},
				)
			}

			if tt.name == "sad case: bad status code" {
				httpmock.RegisterResponder("GET", "/v1/projects/1/forms/test/submissions/instance.xml",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(
							http.StatusBadRequest,
							"",
						)

						return resp, nil
					},
				)
			}

			if tt.name == "sad case: invalid xml" {
				httpmock.RegisterResponder("GET", "/v1/projects/1/forms/test/submissions/instance.xml",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(
							http.StatusOK,
							"invalid xml>#$^&*<",
						)

						return resp, nil
					},
				)
			}

			got, err := s.GetSubmissionXML(tt.args.ctx, tt.args.projectID, tt.args.formID, tt.args.instanceID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetSubmissionXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.GetSubmissionXML() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImpl_GetFormXMLs(t *testing.T) {
	form := `<?xml version="1.0"?><h:html xmlns="http://www.w3.org/2002/xforms" xmlns:ev="http://www.w3.org/2001/xml-events" xmlns:h="http://www.w3.org/1999/xhtml" xmlns:jr="http://openrosa.org/javarosa" xmlns:odk="http://www.opendatakit.org/xforms" xmlns:orx="http://openrosa.org/xforms" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><h:head><h:title>Patient Health Questionnaire-9 (PH Q-9)</h:title><model odk:xforms-version="1.0.0"><instance><data id="akmCQQxf4LaFjAWDbg29pj (1)" version="2 (2022-05-11 06:59:42)"><start/><end/><Over_the_last_2_week_sure_in_doing_things/><Over_the_last_2_week_epressed_or_hopeless/><Over_the_last_2_week_or_sleeping_too_much/><Over_the_last_2_week_having_little_energy/><Over_the_last_2_week_petite_or_overeating/><Over_the_last_2_week_rself_or_family_down/><Over_the_last_2_week_watching_television/><Over_the_last_2_week_lot_more_than_usual/><Over_the_last_2_week_yourself_in_some_way/><If_you_checked_off_a_ng_with_other_people/><__version__/><meta><instanceID/></meta></data></instance><bind jr:preload="timestamp" jr:preloadParams="start" nodeset="/data/start" type="dateTime"/><bind jr:preload="timestamp" jr:preloadParams="end" nodeset="/data/end" type="dateTime"/><bind nodeset="/data/Over_the_last_2_week_sure_in_doing_things" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_epressed_or_hopeless" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_or_sleeping_too_much" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_having_little_energy" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_petite_or_overeating" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_rself_or_family_down" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_watching_television" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_lot_more_than_usual" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_yourself_in_some_way" required="true()" type="string"/><bind nodeset="/data/If_you_checked_off_a_ng_with_other_people" required="true()" type="string"/><bind calculate="'vesMF8UKLW5gnZgBnBmzd9'" nodeset="/data/__version__" type="string"/><bind jr:preload="uid" nodeset="/data/meta/instanceID" readonly="true()" type="string"/></model></h:head><h:body><select1 ref="/data/Over_the_last_2_week_sure_in_doing_things"><label>Over the last 2 weeks how often have you been bothered by little interest of pleasure in doing things?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_epressed_or_hopeless"><label>Over the last 2 weeks how often have you been bothered by feeling down, depressed or hopeless?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_or_sleeping_too_much"><label>Over the last 2 weeks how often have you been bothered by trouble falling asleep or sleeping too much?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_having_little_energy"><label>Over the last 2 weeks how often have you been bothered by feeling tired or having little energy?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_petite_or_overeating"><label>Over the last 2 weeks how often have you been bothered by poor appetite or overeating?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_rself_or_family_down"><label>Over the last 2 weeks how often have you been bothered by feeling bad about yourself -- or that you are a failure or have let yourself or family down/</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_watching_television"><label>Over the last 2 weeks how often have you been bothered by trouble concentrating on things, such as reading the newspaper or watching television?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_lot_more_than_usual"><label>Over the last 2 weeks how often have you been bothered by moving or speaking so slowly that other people could have noticed? Or the opposite — being so fidgety or restless that you have been moving around a lot more than usual?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_yourself_in_some_way"><label>Over the last 2 weeks how often have you been bothered thoughts that you would be better off dead or of hurting yourself in some way?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/If_you_checked_off_a_ng_with_other_people"><label>If you checked off any problems, how difficult have these problems made it for you to do your work, take care of things at home, or get along with other people?</label><item><label>Not Difficult at all</label><value>not_difficult_at_all</value></item><item><label>Somewhat difficult</label><value>somewhat_difficult</value></item><item><label>Very difficult</label><value>very_difficult</value></item><item><label>Extremely difficult</label><value>extremely_difficult</value></item></select1></h:body></h:html>`

	j, err := xj.Convert(strings.NewReader(form))
	if err != nil {
		t.Errorf("failed to convert form")
		return
	}

	parsedForm := make(map[string]interface{})
	err = json.Unmarshal(j.Bytes(), &parsedForm)
	if err != nil {
		t.Errorf("failed to unmarshal form")
		return
	}

	type args struct {
		ctx       context.Context
		projectID int
		formID    string
		version   string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "happy case: retrieve form xml",
			args: args{
				ctx:       context.Background(),
				projectID: 1,
				formID:    "test",
				version:   "test",
			},
			want:    parsedForm,
			wantErr: false,
		},
		{
			name: "sad case: bad status code",
			args: args{
				ctx:       context.Background(),
				projectID: 1,
				formID:    "test",
				version:   "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "sad case: invalid xml",
			args: args{
				ctx:       context.Background(),
				projectID: 1,
				formID:    "test",
				version:   "test",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			surveysClient := surveys.ODKClient{
				BaseURL:    "https://example.com",
				HTTPClient: &http.Client{},
			}
			s := surveys.NewSurveysImpl(surveysClient)

			if tt.name == "happy case: retrieve form xml" {

				httpmock.RegisterResponder("GET", "/v1/projects/1/forms/test/versions/test.xml",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(
							http.StatusOK,
							form,
						)
						return resp, nil
					},
				)
			}

			if tt.name == "sad case: bad status code" {
				httpmock.RegisterResponder("GET", "/v1/projects/1/forms/test/versions/test.xml",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(
							http.StatusBadRequest,
							"",
						)
						return resp, nil
					},
				)
			}

			if tt.name == "sad case: invalid xml" {
				httpmock.RegisterResponder("GET", "/v1/projects/1/forms/test/versions/test.xml",
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(
							http.StatusOK,
							"invalid xml>#$^&*<",
						)
						return resp, nil
					},
				)
			}

			got, err := s.GetFormXML(tt.args.ctx, tt.args.projectID, tt.args.formID, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Impl.GetFormXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Impl.GetFormXML() = %v, want %v", got, tt.want)
			}
		})
	}
}
