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
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires"
)

func TestUseCaseQuestionnaireImpl_CreateScreeningTool(t *testing.T) {
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
					ProgramID:         uuid.NewString(),
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &closedEndedChoice2,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
					ProgramID:         uuid.NewString(),
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &closedEndedChoice2,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
					ProgramID:         uuid.NewString(),
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice:    &openEndedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &openEndedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
					ProgramID:         uuid.NewString(),
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &closedEndedChoice2,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
					ProgramID:         uuid.NewString(),
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &closedEndedChoice2,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
					ProgramID:         uuid.NewString(),
					Choices: []dto.QuestionInputChoiceInput{
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
						},
						{
							Choice:    &closedEndedChoice,
							Value:     "YES",
							Score:     1,
							ProgramID: uuid.NewString(),
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
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)
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
	fakeExtension := extensionMock.NewFakeExtension()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "yes",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: uuid.NewString(),
							Response:   "0",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: uuid.NewString(),
							Response:   "0",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "7",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
							ProgramID:  uuid.NewString(),
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
					ProgramID:       uuid.NewString(),
					QuestionResponses: []*dto.QuestionnaireScreeningToolQuestionResponseInput{
						{
							QuestionID: UUID,
							Response:   "0",
							ProgramID:  uuid.NewString(),
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
						DefaultFacility: &domain.Facility{
							ID:   &UUID,
							Name: gofakeit.Name(),
						},
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
							ID:               new(string),
							Name:             gofakeit.Name(),
							CurrentProgramID: UUID,
						},
						DefaultFacility: &domain.Facility{
							ID:   &UUID,
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:              UUID,
						QuestionnaireID: UUID,
						Questionnaire: domain.Questionnaire{
							ID:        UUID,
							Name:      gofakeit.BeerAlcohol(),
							ProgramID: UUID,
							Questions: []domain.Question{
								{
									ID:        UUID,
									ProgramID: UUID,
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
						User: &domain.User{
							CurrentProgramID: UUID,
						},
						DefaultFacility: &domain.Facility{
							ID:   &UUID,
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:        UUID,
						ProgramID: UUID,
					}, nil
				}
			}

			if tt.name == "Sad case: failed to validate response" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{
						ID: &UUID,
						User: &domain.User{
							ID:               new(string),
							Name:             gofakeit.Name(),
							CurrentProgramID: UUID,
						},
						DefaultFacility: &domain.Facility{
							ID:   &UUID,
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:              UUID,
						QuestionnaireID: UUID,
						Questionnaire: domain.Questionnaire{
							ID:        UUID,
							Name:      gofakeit.BeerAlcohol(),
							ProgramID: UUID,
							Questions: []domain.Question{
								{
									ID:                UUID,
									Active:            false,
									ResponseValueType: enums.QuestionResponseValueTypeBoolean,
									ProgramID:         UUID,
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
							ID:               new(string),
							Name:             gofakeit.Name(),
							CurrentProgramID: UUID,
						},
						DefaultFacility: &domain.Facility{
							ID:   &UUID,
							Name: gofakeit.Name(),
						},
					}, nil
				}
				fakeDB.MockGetScreeningToolByIDFn = func(ctx context.Context, id string) (*domain.ScreeningTool, error) {
					return &domain.ScreeningTool{
						ID:              UUID,
						QuestionnaireID: UUID,
						Questionnaire: domain.Questionnaire{
							ID:        UUID,
							Name:      gofakeit.BeerAlcohol(),
							ProgramID: UUID,
							Questions: []domain.Question{
								{
									ID:        UUID,
									ProgramID: UUID,
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

	type args struct {
		ctx context.Context
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
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get available screening tools",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get logged in user",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get user profile",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case: unable to get available screening tools" {
				fakeDB.MockGetAvailableScreeningToolsFn = func(ctx context.Context, clientID string, screeningTool domain.ScreeningTool, screeningToolIDs []string) ([]*domain.ScreeningTool, error) {
					return nil, errors.New("unable to get available screening tools")
				}
			}
			if tt.name == "Sad case: failed to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("an error occurred")
				}
			}
			_, err := q.GetAvailableScreeningTools(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.GetAvailableScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUseCaseQuestionnaireImpl_GetScreeningToolByID(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)
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
	type args struct {
		ctx             context.Context
		facilityID      string
		paginationInput *dto.PaginationsInput
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
				paginationInput: &dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       10,
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get facility responded screening tools",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get logged in user",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get user profile",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					CurrentPage: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)
			if tt.name == "Sad case: unable to get facility responded screening tools" {
				fakeDB.MockGetFacilityRespondedScreeningToolsFn = func(ctx context.Context, facilityID, programID string, pagination *domain.Pagination) ([]*domain.ScreeningTool, *domain.Pagination, error) {
					return nil, nil, errors.New("unable to get facility responded screening tools")
				}
			}

			if tt.name == "Sad case: failed to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("an error occurred")
				}
			}

			got, err := q.GetFacilityRespondedScreeningTools(tt.args.ctx, tt.args.facilityID, tt.args.paginationInput)
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
	term := "term"
	type args struct {
		ctx             context.Context
		facilityID      string
		screeningToolID string
		searchTerm      *string
		paginationInput *dto.PaginationsInput
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
				paginationInput: &dto.PaginationsInput{
					Limit:       1,
					CurrentPage: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "Happy case: get screening tool respondents",
			args: args{
				ctx:             context.Background(),
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
				paginationInput: &dto.PaginationsInput{
					Limit:       1,
					CurrentPage: 1,
				},
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
		{
			name: "Sad case: failed to get logged in user",
			args: args{
				ctx:             context.Background(),
				facilityID:      uuid.New().String(),
				screeningToolID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get user profile",
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
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)
			if tt.name == "Sad case: unable to get screening tool respondents" {
				fakeDB.MockGetScreeningToolRespondentsFn = func(ctx context.Context, facilityID, ProgramID string, screeningToolID string, searchTerm string, paginationInput *dto.PaginationsInput) ([]*domain.ScreeningToolRespondent, *domain.Pagination, error) {
					return nil, nil, errors.New("failed to get screening tool respondents")
				}
			}
			if tt.name == "Sad case: failed to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", errors.New("an error occurred")
				}
			}
			if tt.name == "Sad case: failed to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, errors.New("an error occurred")
				}
			}
			got, err := q.GetScreeningToolRespondents(tt.args.ctx, tt.args.facilityID, tt.args.screeningToolID, tt.args.searchTerm, tt.args.paginationInput)
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

func TestUseCaseQuestionnaireImpl_GetScreeningToolResponse(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	q := questionnaires.NewUseCaseQuestionnaire(fakeDB, fakeDB, fakeDB, fakeDB, fakeExtension)
	UUID := uuid.New().String()
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.QuestionnaireScreeningToolResponse
		wantErr bool
	}{
		{
			name: "Happy case: get screening tool response",
			args: args{
				ctx: context.Background(),
				id:  UUID,
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to get screening tool response",
			args: args{
				ctx: context.Background(),
				id:  UUID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case: unable to get screening tool response" {
				fakeDB.MockGetScreeningToolResponseByIDFn = func(ctx context.Context, id string) (*domain.QuestionnaireScreeningToolResponse, error) {
					return nil, errors.New("failed to get screening tool response")
				}
			}
			got, err := q.GetScreeningToolResponse(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseQuestionnaireImpl.GetScreeningToolResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCaseQuestionnaireImpl.GetScreeningToolResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
