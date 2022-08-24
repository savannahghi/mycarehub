package mock

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	xj "github.com/basgys/goxml2json"
	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/surveys"
)

// SurveysMock mocks the surveys service
type SurveysMock struct {
	MockMakeRequestFn              func(ctx context.Context, payload surveys.RequestHelperPayload) (*http.Response, error)
	MockListSurveyFormsFn          func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error)
	MockGetSurveyFormFn            func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error)
	MockGeneratePublicAccessLinkFn func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error)
	MockGetSubmissionsFn           func(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error)
	MockDeletePublicAccessLinkFn   func(ctx context.Context, input dto.VerifySurveySubmissionInput) error
	MockListSubmittersFn           func(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error)
	MockListPublicAccessLinksFn    func(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error)
	MockGetFormXMLFn               func(ctx context.Context, projectID int, formID, version string) (map[string]interface{}, error)
	MockGetSubmissionXMLFn         func(ctx context.Context, projectID int, formID, instanceID string) (map[string]interface{}, error)
}

// NewSurveysMock initializes the surveys mock service
func NewSurveysMock() *SurveysMock {
	return &SurveysMock{

		MockMakeRequestFn: func(ctx context.Context, payload surveys.RequestHelperPayload) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       nil,
			}, nil
		},
		MockGetFormXMLFn: func(ctx context.Context, projectID int, formID, version string) (map[string]interface{}, error) {
			parsedForm := make(map[string]interface{})

			form := `<?xml version="1.0"?><h:html xmlns="http://www.w3.org/2002/xforms" xmlns:ev="http://www.w3.org/2001/xml-events" xmlns:h="http://www.w3.org/1999/xhtml" xmlns:jr="http://openrosa.org/javarosa" xmlns:odk="http://www.opendatakit.org/xforms" xmlns:orx="http://openrosa.org/xforms" xmlns:xsd="http://www.w3.org/2001/XMLSchema"><h:head><h:title>Patient Health Questionnaire-9 (PH Q-9)</h:title><model odk:xforms-version="1.0.0"><instance><data id="akmCQQxf4LaFjAWDbg29pj (1)" version="2 (2022-05-11 06:59:42)"><start/><end/><Over_the_last_2_week_sure_in_doing_things/><Over_the_last_2_week_epressed_or_hopeless/><Over_the_last_2_week_or_sleeping_too_much/><Over_the_last_2_week_having_little_energy/><Over_the_last_2_week_petite_or_overeating/><Over_the_last_2_week_rself_or_family_down/><Over_the_last_2_week_watching_television/><Over_the_last_2_week_lot_more_than_usual/><Over_the_last_2_week_yourself_in_some_way/><If_you_checked_off_a_ng_with_other_people/><__version__/><meta><instanceID/></meta></data></instance><bind jr:preload="timestamp" jr:preloadParams="start" nodeset="/data/start" type="dateTime"/><bind jr:preload="timestamp" jr:preloadParams="end" nodeset="/data/end" type="dateTime"/><bind nodeset="/data/Over_the_last_2_week_sure_in_doing_things" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_epressed_or_hopeless" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_or_sleeping_too_much" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_having_little_energy" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_petite_or_overeating" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_rself_or_family_down" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_watching_television" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_lot_more_than_usual" required="true()" type="string"/><bind nodeset="/data/Over_the_last_2_week_yourself_in_some_way" required="true()" type="string"/><bind nodeset="/data/If_you_checked_off_a_ng_with_other_people" required="true()" type="string"/><bind calculate="'vesMF8UKLW5gnZgBnBmzd9'" nodeset="/data/__version__" type="string"/><bind jr:preload="uid" nodeset="/data/meta/instanceID" readonly="true()" type="string"/></model></h:head><h:body><select1 ref="/data/Over_the_last_2_week_sure_in_doing_things"><label>Over the last 2 weeks how often have you been bothered by little interest of pleasure in doing things?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_epressed_or_hopeless"><label>Over the last 2 weeks how often have you been bothered by feeling down, depressed or hopeless?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_or_sleeping_too_much"><label>Over the last 2 weeks how often have you been bothered by trouble falling asleep or sleeping too much?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_having_little_energy"><label>Over the last 2 weeks how often have you been bothered by feeling tired or having little energy?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_petite_or_overeating"><label>Over the last 2 weeks how often have you been bothered by poor appetite or overeating?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_rself_or_family_down"><label>Over the last 2 weeks how often have you been bothered by feeling bad about yourself -- or that you are a failure or have let yourself or family down/</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_watching_television"><label>Over the last 2 weeks how often have you been bothered by trouble concentrating on things, such as reading the newspaper or watching television?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_lot_more_than_usual"><label>Over the last 2 weeks how often have you been bothered by moving or speaking so slowly that other people could have noticed? Or the opposite — being so fidgety or restless that you have been moving around a lot more than usual?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/Over_the_last_2_week_yourself_in_some_way"><label>Over the last 2 weeks how often have you been bothered thoughts that you would be better off dead or of hurting yourself in some way?</label><item><label>Not at all</label><value>0</value></item><item><label>Several Days</label><value>1</value></item><item><label>More than half a day</label><value>2</value></item><item><label>Nearly every day</label><value>3</value></item></select1><select1 ref="/data/If_you_checked_off_a_ng_with_other_people"><label>If you checked off any problems, how difficult have these problems made it for you to do your work, take care of things at home, or get along with other people?</label><item><label>Not Difficult at all</label><value>not_difficult_at_all</value></item><item><label>Somewhat difficult</label><value>somewhat_difficult</value></item><item><label>Very difficult</label><value>very_difficult</value></item><item><label>Extremely difficult</label><value>extremely_difficult</value></item></select1></h:body></h:html>`

			j, err := xj.Convert(strings.NewReader(form))
			if err != nil {
				return parsedForm, err
			}

			err = json.Unmarshal(j.Bytes(), &parsedForm)
			if err != nil {
				return parsedForm, err
			}

			return parsedForm, nil

		},
		MockGetSubmissionXMLFn: func(ctx context.Context, projectID int, formID, instanceID string) (map[string]interface{}, error) {
			parsedSubmission := make(map[string]interface{})
			submission := `
				<?xml version="1.0" encoding="UTF-8"?>
				<data xmlns:jr="http://openrosa.org/javarosa" xmlns:orx="http://openrosa.org/xforms" id="akmCQQxf4LaFjAWDbg29pj (1)" version="2 (2022-05-11 06:59:42)">
					<start>2022-08-12T13:11:23.965+03:00</start>
					<end>2022-08-12T13:12:01.646+03:00</end>
					<Over_the_last_2_week_sure_in_doing_things>0</Over_the_last_2_week_sure_in_doing_things>
					<Over_the_last_2_week_epressed_or_hopeless>0</Over_the_last_2_week_epressed_or_hopeless>
					<Over_the_last_2_week_or_sleeping_too_much>0</Over_the_last_2_week_or_sleeping_too_much>
					<Over_the_last_2_week_having_little_energy>0</Over_the_last_2_week_having_little_energy>
					<Over_the_last_2_week_petite_or_overeating>0</Over_the_last_2_week_petite_or_overeating>
					<Over_the_last_2_week_rself_or_family_down>0</Over_the_last_2_week_rself_or_family_down>
					<Over_the_last_2_week_watching_television>0</Over_the_last_2_week_watching_television>
					<Over_the_last_2_week_lot_more_than_usual>0</Over_the_last_2_week_lot_more_than_usual>
					<Over_the_last_2_week_yourself_in_some_way>0</Over_the_last_2_week_yourself_in_some_way>
					<If_you_checked_off_a_ng_with_other_people>not_difficult_at_all</If_you_checked_off_a_ng_with_other_people>
					<__version__>vesMF8UKLW5gnZgBnBmzd9</__version__>
					<meta>
						<instanceID>uuid:808431e7-e2ed-4065-b19d-9fd780ce7f9c</instanceID>
					</meta>
					<send_alert>true</send_alert>
				</data>
			`

			j, err := xj.Convert(strings.NewReader(submission))
			if err != nil {
				return parsedSubmission, err
			}

			err = json.Unmarshal(j.Bytes(), &parsedSubmission)
			if err != nil {
				return parsedSubmission, err
			}

			return parsedSubmission, nil
		},
		MockListSurveyFormsFn: func(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
			return []*domain.SurveyForm{
				{
					ProjectID: 2,
					Name:      gofakeit.Name(),
					EnketoID:  gofakeit.UUID(),
				},
			}, nil
		},
		MockGetSurveyFormFn: func(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
			return &domain.SurveyForm{
				ProjectID: 2,
				Name:      gofakeit.Name(),
				EnketoID:  gofakeit.UUID(),
			}, nil
		},
		MockGeneratePublicAccessLinkFn: func(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
			return &dto.SurveyPublicLink{
				Once:        true,
				ID:          2,
				DisplayName: gofakeit.Name(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				DeletedAt:   nil,
				Token:       gofakeit.UUID(),
			}, nil
		},
		MockListPublicAccessLinksFn: func(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error) {
			return []*dto.SurveyPublicLink{
				{
					Once:        true,
					ID:          2,
					DisplayName: gofakeit.Name(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
					DeletedAt:   nil,
					Token:       gofakeit.UUID(),
				},
			}, nil
		},
		MockGetSubmissionsFn: func(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
			return []domain.Submission{
				{
					InstanceID:  gofakeit.UUID(),
					SubmitterID: 1096,
					DeviceID:    gofakeit.UUID(),
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
					ReviewState: gofakeit.UUID(),
					Submitter: domain.Submitter{
						ID:          1096,
						Type:        gofakeit.BeerAlcohol(),
						DisplayName: gofakeit.BeerBlg(),
						CreatedAt:   time.Time{},
						UpdatedAt:   time.Time{},
						DeletedAt:   time.Time{},
					},
				},
			}, nil
		},
		MockDeletePublicAccessLinkFn: func(ctx context.Context, input dto.VerifySurveySubmissionInput) error {
			return nil
		},
		MockListSubmittersFn: func(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
			return []domain.Submitter{
				{
					ID:          10,
					Type:        "test",
					DisplayName: "test",
					CreatedAt:   time.Time{},
					UpdatedAt:   time.Time{},
					DeletedAt:   time.Time{},
				},
			}, nil
		},
	}
}

// ListSurveyForms lists the survey forms for the given project
func (s *SurveysMock) ListSurveyForms(ctx context.Context, projectID int) ([]*domain.SurveyForm, error) {
	return s.MockListSurveyFormsFn(ctx, projectID)
}

// MakeRequest makes a request to the surveys service
func (s *SurveysMock) MakeRequest(ctx context.Context, payload surveys.RequestHelperPayload) (*http.Response, error) {
	return s.MockMakeRequestFn(ctx, payload)
}

// GetSurveyForm gets the survey form for the given project and form ID
func (s *SurveysMock) GetSurveyForm(ctx context.Context, projectID int, formID string) (*domain.SurveyForm, error) {
	return s.MockGetSurveyFormFn(ctx, projectID, formID)
}

// GeneratePublicAccessLink generates a public access link for the given survey
func (s *SurveysMock) GeneratePublicAccessLink(ctx context.Context, input dto.SurveyLinkInput) (*dto.SurveyPublicLink, error) {
	return s.MockGeneratePublicAccessLinkFn(ctx, input)
}

// GetSubmissions mocks the action of getting the submissions for the given survey
func (s *SurveysMock) GetSubmissions(ctx context.Context, input dto.VerifySurveySubmissionInput) ([]domain.Submission, error) {
	return s.MockGetSubmissionsFn(ctx, input)
}

// DeletePublicAccessLink mocks the implementation of deleting the public access link for the given survey
func (s *SurveysMock) DeletePublicAccessLink(ctx context.Context, input dto.VerifySurveySubmissionInput) error {
	return s.MockDeletePublicAccessLinkFn(ctx, input)
}

// ListSubmitters mocks the action of listing all the submitters of a given survey
func (s *SurveysMock) ListSubmitters(ctx context.Context, projectID int, formID string) ([]domain.Submitter, error) {
	return s.MockListSubmittersFn(ctx, projectID, formID)
}

// ListPublicAccessLinks returns a list of all public access links created for a particular form
func (s *SurveysMock) ListPublicAccessLinks(ctx context.Context, projectID int, formID string) ([]*dto.SurveyPublicLink, error) {
	return s.MockListPublicAccessLinksFn(ctx, projectID, formID)
}

//GetFormXML retrieves a form's XML definition
func (s *SurveysMock) GetFormXML(ctx context.Context, projectID int, formID, version string) (map[string]interface{}, error) {
	return s.MockGetFormXMLFn(ctx, projectID, formID, version)
}

//GetSubmissionXML retrieves a submission's XML definition using the instance id
func (s *SurveysMock) GetSubmissionXML(ctx context.Context, projectID int, formID, instanceID string) (map[string]interface{}, error) {
	return s.MockGetSubmissionXMLFn(ctx, projectID, formID, instanceID)
}
