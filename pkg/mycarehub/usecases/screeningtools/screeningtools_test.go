package screeningtools

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
)

func TestServiceScreeningToolsImpl_GetScreeningToolsQuestions(t *testing.T) {
	tbQuestionType := enums.ScreeningToolTypeTB.String()
	invalidToolType := enums.ScreeningToolType("invalid").String()
	type args struct {
		ctx      context.Context
		toolType *string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningToolQuestion
		wantErr bool
	}{
		{
			name: "happy case: get all screening tools questions",
			args: args{
				ctx:      context.Background(),
				toolType: &tbQuestionType,
			},
			wantErr: false,
		},
		{
			name: "happy case: get all screening tools questions, no params",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
		{
			name: "sad case: get screening tools questions, invalid question type",
			args: args{
				ctx:      context.Background(),
				toolType: &invalidToolType,
			},
			wantErr: true,
		},
		{
			name: "sad case: get screening tools questions, error",
			args: args{
				ctx:      context.Background(),
				toolType: &tbQuestionType,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			tr := NewUseCasesScreeningTools(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: get screening tools questions, error" {
				fakeDB.MockGetScreeningToolsQuestionsFn = func(ctx context.Context, toolType string) ([]*domain.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("failed to get screening tools questions")
				}
			}

			got, err := tr.GetScreeningToolQuestions(tt.args.ctx, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.GetScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("ServiceScreeningToolsImpl.GetScreeningToolQuestions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceScreeningToolsImpl_AnswerScreeningToolQuestions(t *testing.T) {
	type args struct {
		ctx                    context.Context
		screeningToolResponses []*dto.ScreeningToolQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case: answer screening tool questions",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "Sad case: empty response",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "invalid response",
					},
				},
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "Sad case: empty response",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
					},
				},
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "sad case: empty screening tool responses",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get screening tools question by id",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get client by id",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid screening tool response type, int not valid",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "invalid",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid screening tool response type, key not valid",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "2",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to invalidate previous responses",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to save screening tool response",
			args: args{
				ctx: context.Background(),
				screeningToolResponses: []*dto.ScreeningToolQuestionResponseInput{
					{
						ClientID:   uuid.New().String(),
						QuestionID: uuid.New().String(),
						Response:   "0",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			tr := NewUseCasesScreeningTools(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: failed to get screening tools question by id" {
				fakeDB.MockGetScreeningToolQuestionByQuestionIDFn = func(ctx context.Context, id string) (*domain.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("failed to get screening tool question by id")
				}
			}

			if tt.name == "sad case: failed to save screening tool response" {
				fakeDB.MockAnswerScreeningToolQuestionsFn = func(ctx context.Context, screeningToolResponses []*dto.ScreeningToolQuestionResponseInput) error {
					return fmt.Errorf("failed to save screening tool response")
				}
			}
			if tt.name == "sad case: failed to invalidate previous responses" {
				fakeDB.MockInvalidateScreeningToolResponseFn = func(ctx context.Context, clientID string, questionID string) error {
					return fmt.Errorf("failed to invalidate previous responses")
				}
			}

			if tt.name == "sad case: failed to get client by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client by id")
				}
			}

			got, err := tr.AnswerScreeningToolQuestions(tt.args.ctx, tt.args.screeningToolResponses)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.AnswerScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceScreeningToolsImpl.AnswerScreeningToolQuestions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceScreeningToolsImpl_GetAvailableScreeningToolQuestions(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	tr := NewUseCasesScreeningTools(fakeDB, fakeDB, fakeDB, fakeExtension)
	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get available screening tool questions",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
		},
		{
			name: "sad case: failed to get client by id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get active screening tool responses",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get system generated service request",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get screening tool question by id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: failed to get client by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client by id")
				}
			}

			if tt.name == "sad case: failed to get active screening tool responses" {
				fakeDB.MockGetActiveScreeningToolResponsesFn = func(ctx context.Context, clientID string) ([]*domain.ScreeningToolQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get active screening tool responses")
				}
			}

			if tt.name == "sad case: failed to get system generated service request" {
				fakeDB.MockGetClientServiceRequestsFn = func(ctx context.Context, toolType string, status string, clientID string) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get system generated service request")
				}
			}

			if tt.name == "sad case: failed to get screening tool question by id" {
				fakeDB.MockGetScreeningToolQuestionByQuestionIDFn = func(ctx context.Context, id string) (*domain.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("failed to get screening tool question by id")
				}
			}

			got, err := tr.GetAvailableScreeningToolQuestions(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.GetAvailableScreeningToolQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("expected to get screening tool questions: %v", got)
			}
		})
	}
}

func TestServiceScreeningToolsImpl_GetAvailableFacilityScreeningTools(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	tr := NewUseCasesScreeningTools(fakeDB, fakeDB, fakeDB, fakeExtension)
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
			name: "happy case: get available facility screening tools",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
			},
			wantErr: false,
		},
		{
			name: "Sad case: failed to get service requests",
			args: args{
				ctx:        context.Background(),
				facilityID: uuid.New().String(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "Sad case: failed to get service requests" {
				fakeDB.MockGetServiceRequestsFn = func(ctx context.Context, requestType *string, requestStatus *string, facilityID string, flavour feedlib.Flavour) ([]*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get service requests")
				}
			}

			got, err := tr.GetAvailableFacilityScreeningTools(tt.args.ctx, tt.args.facilityID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.GetAvailableFacilityScreeningTools() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) == 0 {
				t.Errorf("expected to get screening tool questions: %v", got)
			}
		})
	}
}

func TestServiceScreeningToolsImpl_GetAssessmentResponses(t *testing.T) {
	ctx := context.Background()

	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	tr := NewUseCasesScreeningTools(fakeDB, fakeDB, fakeDB, fakeExtension)

	type args struct {
		ctx        context.Context
		facilityID string
		toolType   string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.ScreeningToolAssessmentResponse
		wantErr bool
	}{
		{

			name: "Happy case:  get assessment responses",
			args: args{
				ctx:        ctx,
				facilityID: uuid.New().String(),
				toolType:   "VIOLENCE_ASSESSMENT",
			},
			wantErr: false,
		},
		{
			name: "Sad case:  unable to get assessment responses",
			args: args{
				ctx:        ctx,
				facilityID: "facilityID",
				toolType:   "GBV",
			},
			wantErr: true,
		},
		{
			name: "Sad case:  empty facility",
			args: args{
				ctx:        ctx,
				facilityID: "",
				toolType:   "VIOLENCE_ASSESSMENT",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Sad case:  unable to get assessment responses" {
				fakeDB.MockGetAssessmentResponsesFn = func(ctx context.Context, facilityID, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
					return nil, fmt.Errorf("failed to get assessment responses")
				}
			}
			if tt.name == "Sad case:  empty facility" {
				fakeDB.MockGetAssessmentResponsesFn = func(ctx context.Context, facilityID, toolType string) ([]*domain.ScreeningToolAssessmentResponse, error) {
					return nil, fmt.Errorf("failed to get assessment responses")
				}
			}

			got, err := tr.GetAssessmentResponses(tt.args.ctx, tt.args.facilityID, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.GetAssessmentResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected response to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected response not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestServiceScreeningToolsImpl_GetScreeningToolServiceRequestResponses(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	tr := NewUseCasesScreeningTools(fakeDB, fakeDB, fakeDB, fakeExtension)
	type args struct {
		ctx      context.Context
		clientID string
		toolType enums.ScreeningToolType
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case: get screening tool responses",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolTypeGBV,
			},
		},
		{
			name: "sad case: empty client ID",
			args: args{
				ctx:      context.Background(),
				toolType: enums.ScreeningToolTypeGBV,
			},
			wantErr: true,
		},
		{
			name: "sad case: invalid tool type",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolType("invalid"),
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get screening tool responses by tool type",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolTypeGBV,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get screening tool question by id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolTypeGBV,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get user profile by client ID",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolTypeGBV,
			},
			wantErr: true,
		},
		{
			name: "sad case: failed to get client service request by question id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolTypeGBV,
			},
			wantErr: true,
		},
		{
			name: "sad case: client profile by client id",
			args: args{
				ctx:      context.Background(),
				clientID: uuid.New().String(),
				toolType: enums.ScreeningToolTypeGBV,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "sad case: failed to get screening tool responses by tool type" {
				fakeDB.MockGetClientScreeningToolResponsesByToolTypeFn = func(ctx context.Context, clientID string, toolType string, active bool) ([]*domain.ScreeningToolQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get active screening tool responses")
				}
			}

			if tt.name == "sad case: failed to get screening tool question by id" {
				fakeDB.MockGetScreeningToolQuestionByQuestionIDFn = func(ctx context.Context, id string) (*domain.ScreeningToolQuestion, error) {
					return nil, fmt.Errorf("failed to get screening tool question by id")
				}
			}

			if tt.name == "sad case: failed to get user profile by client ID" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client by id")
				}
			}
			if tt.name == "sad case: failed to get client service request by question id" {
				fakeDB.MockGetClientScreeningToolServiceRequestByToolTypeFn = func(ctx context.Context, clientID string, toolType string, status string) (*domain.ServiceRequest, error) {
					return nil, fmt.Errorf("failed to get client service request by question id")
				}
			}

			if tt.name == "sad case: client profile by client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to get client by id")
				}
			}
			got, err := tr.GetScreeningToolServiceRequestResponses(tt.args.ctx, tt.args.clientID, tt.args.toolType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceScreeningToolsImpl.GetScreeningToolServiceRequestResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected to get responses: %v", got)
			}
		})
	}
}
