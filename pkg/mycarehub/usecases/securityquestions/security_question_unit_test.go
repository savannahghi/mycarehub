package securityquestions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions/mock"
	"github.com/segmentio/ksuid"
)

func TestUseCaseSecurityQuestionsImpl_GetSecurityQuestions(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.SecurityQuestion
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
		{
			name: "Sad case",
			args: args{
				ctx:     ctx,
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid flavor",
			args: args{
				ctx:     ctx,
				flavour: "invalid-flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - nil flavor",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewSecurityQuestionsUseCaseMock()

			fakeExtension := extensionMock.NewFakeExtension()
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case" {
				fakeDB.MockGetSecurityQuestionsFn = func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - invalid flavor" {
				fakeDB.MockGetSecurityQuestionsFn = func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case - nil flavor" {
				fakeDB.MockGetSecurityQuestionsFn = func(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("an error occurred")
				}
			}

			got, err := s.GetSecurityQuestions(tt.args.ctx, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseSecurityQuestionsImpl.GetSecurityQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseSecurityQuestionsImpl_RecordSecurityQuestionResponses(t *testing.T) {
	type args struct {
		ctx   context.Context
		input []*dto.SecurityQuestionResponseInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: valid number response",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             ksuid.New().String(),
						SecurityQuestionID: ksuid.New().String(),
						Response:           "20",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: invalid valid number response",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             ksuid.New().String(),
						SecurityQuestionID: ksuid.New().String(),
						Response:           "invalid",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing user id",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						SecurityQuestionID: ksuid.New().String(),
						Response:           "20",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing security question",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						UserID:   ksuid.New().String(),
						Response: "20",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: missing response",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             ksuid.New().String(),
						SecurityQuestionID: ksuid.New().String(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed to get security question by id",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             ksuid.New().String(),
						SecurityQuestionID: ksuid.New().String(),
						Response:           "20",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Sad case: failed save security question response",
			args: args{
				ctx: context.Background(),
				input: []*dto.SecurityQuestionResponseInput{
					{
						UserID:             ksuid.New().String(),
						SecurityQuestionID: ksuid.New().String(),
						Response:           "20",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			_ = mock.NewSecurityQuestionsUseCaseMock()

			fakeExtension := extensionMock.NewFakeExtension()
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case: failed to get security question by id" {
				fakeDB.MockGetSecurityQuestionByIDFn = func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("failed to get security questions")
				}
			}

			if tt.name == "Sad case: failed save security question response" {
				fakeDB.MockSaveSecurityQuestionResponseFn = func(ctx context.Context, securityQuestionResponse *dto.SecurityQuestionResponseInput) error {
					return fmt.Errorf("failed to save security question response")
				}
			}

			got, err := s.RecordSecurityQuestionResponses(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseSecurityQuestionsImpl.RecordSecurityQuestionResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCaseSecurityQuestionsImpl_VerifySecurityQuestionResponses(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx       context.Context
		responses *[]dto.VerifySecurityQuestionInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully verify security question response",
			args: args{
				ctx: ctx,
				responses: &[]dto.VerifySecurityQuestionInput{
					{
						QuestionID: "1234",
						Flavour:    feedlib.FlavourConsumer,
						Response:   "",
						UserID:     uuid.New().String(),
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - Fail to get security question by ID",
			args: args{
				ctx: ctx,
				responses: &[]dto.VerifySecurityQuestionInput{
					{
						QuestionID: "1234",
						Flavour:    feedlib.FlavourConsumer,
						Response:   "Nairobi",
						UserID:     uuid.New().String(),
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - response mismatch",
			args: args{
				ctx: ctx,
				responses: &[]dto.VerifySecurityQuestionInput{
					{
						QuestionID: "1234",
						Flavour:    feedlib.FlavourConsumer,
						Response:   "Nakuru",
						UserID:     uuid.New().String(),
					},
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeSecurity := mock.NewSecurityQuestionsUseCaseMock()

			fakeExtension := extensionMock.NewFakeExtension()
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad Case - Fail to get security question by ID" {
				fakeDB.MockGetSecurityQuestionResponseByIDFn = func(ctx context.Context, questionID string) (*domain.SecurityQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get security question response")
				}
			}

			if tt.name == "Sad Case - response mismatch" {
				fakeSecurity.MockVerifySecurityQuestionResponsesFn = func(
					ctx context.Context,
					responses *[]dto.VerifySecurityQuestionInput,
				) (bool, error) {
					return false, fmt.Errorf("the response does not match")
				}
			}

			got, err := s.VerifySecurityQuestionResponses(tt.args.ctx, tt.args.responses)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseSecurityQuestionsImpl.VerifySecurityQuestionResponses() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCaseSecurityQuestionsImpl.VerifySecurityQuestionResponses() = %v, want %v", got, tt.want)
			}
		})
	}
}
