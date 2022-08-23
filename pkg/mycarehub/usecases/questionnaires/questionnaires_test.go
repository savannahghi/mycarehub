package questionnaires_test

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires"
)

func TestUseCaseQuestionnaireImpl_CreateScreeningTool(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)
	closedEndedChoice := "1"
	closedEndedChoice2 := "2"
	openEndedEndedChoice := "YES"

	questionnare := &dto.ScreeningToolInput{
		Questionnaire: dto.QuestionnaireInput{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.Sentence(20),
			Questions: []*dto.QuestionInput{
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &closedEndedChoice2,
							Value:  "YES",
							Score:  1,
						},
					},
				},
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          2,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &closedEndedChoice2,
							Value:  "YES",
							Score:  1,
						},
					},
				},
			},
		},
		Threshold:   3,
		ClientTypes: []enums.ClientType{enums.ClientTypePmtct, enums.ClientTypePmtct},
		Genders:     []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: dto.AgeRangeInput{
			LowerBound: 10,
			UpperBound: 20,
		},
	}
	invalidQuestionnare := &dto.ScreeningToolInput{
		Questionnaire: dto.QuestionnaireInput{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.Sentence(20),
			Questions: []*dto.QuestionInput{
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeOpenEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &openEndedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &openEndedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
					},
				},
			},
		},
		Threshold:   3,
		ClientTypes: []enums.ClientType{enums.ClientTypePmtct, enums.ClientTypePmtct},
		Genders:     []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: dto.AgeRangeInput{
			LowerBound: 10,
			UpperBound: 20,
		},
	}
	duplicateSequenceQuestionnaire := &dto.ScreeningToolInput{
		Questionnaire: dto.QuestionnaireInput{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.Sentence(20),
			Questions: []*dto.QuestionInput{
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &closedEndedChoice2,
							Value:  "YES",
							Score:  1,
						},
					},
				},
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &closedEndedChoice2,
							Value:  "YES",
							Score:  1,
						},
					},
				},
			},
		},
		Threshold:   3,
		ClientTypes: []enums.ClientType{enums.ClientTypePmtct, enums.ClientTypePmtct},
		Genders:     []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: dto.AgeRangeInput{
			LowerBound: 10,
			UpperBound: 20,
		},
	}
	duplicateChoiceQuestionnaire := &dto.ScreeningToolInput{
		Questionnaire: dto.QuestionnaireInput{
			Name:        gofakeit.BeerName(),
			Description: gofakeit.Sentence(20),
			Questions: []*dto.QuestionInput{
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          1,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
					},
				},
				{
					Text:              gofakeit.BeerAlcohol(),
					QuestionType:      enums.QuestionTypeCloseEnded,
					ResponseValueType: enums.QuestionResponseValueTypeString,
					Required:          true,
					SelectMultiple:    true,
					Sequence:          2,
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
						{
							Choice: &closedEndedChoice,
							Value:  "YES",
							Score:  1,
						},
					},
				},
			},
		},
		Threshold:   3,
		ClientTypes: []enums.ClientType{enums.ClientTypePmtct, enums.ClientTypePmtct},
		Genders:     []enumutils.Gender{enumutils.GenderFemale},
		AgeRange: dto.AgeRangeInput{
			LowerBound: 10,
			UpperBound: 20,
		},
	}
	type args struct {
		ctx   context.Context
		input dto.ScreeningToolInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: Create screening tool",
			args: args{
				ctx:   context.Background(),
				input: *questionnare,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: unable to create screening tool",
			args: args{
				ctx:   context.Background(),
				input: *questionnare,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: unable to validate question",
			args: args{
				ctx:   context.Background(),
				input: *invalidQuestionnare,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: duplicate sequence",
			args: args{
				ctx:   context.Background(),
				input: *duplicateSequenceQuestionnaire,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: duplicate closed ended choice",
			args: args{
				ctx:   context.Background(),
				input: *duplicateChoiceQuestionnaire,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to create screening tool" {
				fakeDB.MockCreateScreeningToolFn = func(ctx context.Context, input *domain.ScreeningTool) error {
					return errors.New("unable to create screening tool")
				}
			}
			got, err := q.CreateScreeningTool(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.CreateScreeningTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseQuestionnaireImpl.CreateScreeningTool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseQuestionnaireImpl_RespondToScreeningTool(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)
	UUID := "f3f8f8f8-f3f8-f3f8-f3f8-f3f8f8f8f8f8"
	type args struct {
		ctx   context.Context
		input dto.QuestionnaireScreeningToolResponseInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: Respond to screening tool",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
						},
					},
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case: Invalid input",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "yes",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: Invalid input",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get client profile by client id",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: uuid.NewString(),
					ClientID:        uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: uuid.NewString(),
							Response:   "0",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to get screening tool by id",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to get question by id",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: uuid.NewString(),
							Response:   "0",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to validate response",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "7",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to create screening tool response",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to create service request",
			args: args{
				ctx: context.Background(),
				input: dto.QuestionnaireScreeningToolResponseInput{
					ScreeningToolID: UUID,
					ClientID:        UUID,
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case: failed to get client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, errors.New("failed to get client profile by client id")
				}
			}

			if tt.name == "Sad case: failed to get screening tool by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &UUID,
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return nil, errors.New("failed to get screening tool by id")
				}
			}

			if tt.name == "Sad case: failed to create service request" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &UUID,
						User: &domain.User{
							ID:   new(string),
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:              UUID,
						QuestionnaireID: UUID,
						Questionnaire: domain.Questionnaire{
							ID:   UUID,
							Name: gofakeit.BeerAlcohol(),
							Questions: []domain.Question{
								{
									ID: UUID,
								},
							},
						},
					}, nil
				}
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
					return errors.New("failed to create service request")
				}
			}

			if tt.name == "Sad case: failed to get question by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &UUID,
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID: UUID,
					}, nil
				}
			}

			if tt.name == "Sad case: failed to validate response" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &UUID,
						User: &domain.User{
							ID:   new(string),
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:              UUID,
						QuestionnaireID: UUID,
						Questionnaire: domain.Questionnaire{
							ID:   UUID,
							Name: gofakeit.BeerAlcohol(),
							Questions: []domain.Question{
								{
									ID:                UUID,
									Active:            false,
									ResponseValueType: enums.QuestionResponseValueTypeBoolean,
								},
							},
						},
					}, nil
				}
			}

			if tt.name == "Sad case: failed to create screening tool response" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &UUID,
						User: &domain.User{
							ID:   new(string),
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:              UUID,
						QuestionnaireID: UUID,
						Questionnaire: domain.Questionnaire{
							ID:   UUID,
							Name: gofakeit.BeerAlcohol(),
							Questions: []domain.Question{
								{
									ID: UUID,
								},
							},
						},
					}, nil
				}
				fakeDB.MockCreateScreeningToolResponseFn = func(ctx context.Context, input *domain.QuestionnaireScreeningToolResponse) (*string, error) {
					return nil, errors.New("failed to create screening tool response")
				}
			}

			got, err := q.RespondToScreeningTool(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.RespondToScreeningTool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseQuestionnaireImpl.RespondToScreeningTool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCaseQuestionnaireImpl_GetAvailableScreeningTools(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)

	type args struct {
		ctx        context.Context
		clientID   string
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningTool
		wantErr bool
	}{
		{
			name: "Happy case: Get available screening tools",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.New().String(),
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get available screening tools",
			args: args{
				ctx:        context.Background(),
				clientID:   uuid.New().String(),
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get available screening tools" {
				fakeDB.MockGetAvailableScreeningToolsFn = func(ctx context.Context, clientID string, facilityID string) ([]*domain.ScreeningTool, error) {
					return nil, errors.New("unable to get available screening tools")
				}
			}
			_, err := q.GetAvailableScreeningTools(tt.args.ctx, tt.args.clientID, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.GetAvailableScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseQuestionnaireImpl_GetScreeningToolByID(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)
	UUID := uuid.New().String()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool by id",
			args: args{
				ctx: context.Background(),
				id:  UUID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get screening tool by id",
			args: args{
				ctx: context.Background(),
				id:  UUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get screening tool by id" {
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return nil, errors.New("failed to get screening tool by id")
				}
			}
			_, err := q.GetScreeningToolByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.GetScreeningToolByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseQuestionnaireImpl_GetFacilityRespondedScreeningTools(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)

	type args struct {
		ctx        context.Context
		facilityID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: Get facility responded screening tools",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get facility responded screening tools",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get facility responded screening tools" {
				fakeDB.MockGetFacilityRespondedScreeningToolsFn = func(ctx context.Context, facilityID string) ([]*domain.ScreeningTool, error) {
					return nil, errors.New("unable to get facility responded screening tools")
				}
			}
			got, err := q.GetFacilityRespondedScreeningTools(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.GetFacilityRespondedScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected %v", got)
			}
		})
	}
}

func TestUseCaseQuestionnaireImpl_GetScreeningToolRespondents(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB)
	term := "term"
	type args struct {
		ctx             context.Context
		facilityID      string
		screeningToolID string
		searchTerm      *string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ClientProfile
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool respondents",
			args: args{
				ctx:             context.Background(),
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
				searchTerm:      &term,
			},
			wantErr: false,
		},
		{
			name: "Happy case: get screening tool respondents",
			args: args{
				ctx:             context.Background(),
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get screening tool respondents",
			args: args{
				ctx:             context.Background(),
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get screening tool respondents" {
				fakeDB.MockGetScreeningToolRespondentsFn = func(ctx context.Context, facilityID string, screeningToolID string, searchTerm string) ([]*domain.ScreeningToolRespondent, error) {
					return nil, errors.New("failed to get screening tool respondents")
				}
			}
			got, err := q.GetScreeningToolRespondents(tt.args.ctx, tt.args.facilityID, tt.args.screeningToolID, tt.args.searchTerm)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.GetScreeningToolRespondents() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCaseQuestionnaireImpl.GetScreeningToolRespondents() = %v, want %v", got, tt.want)
			}
		})
	}
}
