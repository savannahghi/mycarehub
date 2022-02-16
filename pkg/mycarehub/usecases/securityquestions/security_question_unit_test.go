package securityquestions_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/interserviceclient"
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
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeDB, fakeExtension)

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
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case: failed to get security question by id" {
				fakeDB.MockGetSecurityQuestionByIDFn = func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("failed to get security questions")
				}
			}

			if tt.name == "Sad case: failed save security question response" {
				fakeDB.MockSaveSecurityQuestionResponseFn = func(ctx context.Context, securityQuestionResponse []*dto.SecurityQuestionResponseInput) error {
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
		responses *dto.VerifySecurityQuestionsPayload
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
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{
						{
							QuestionID:  "1234",
							Flavour:     feedlib.FlavourConsumer,
							Response:    "",
							PhoneNumber: interserviceclient.TestUserPhoneNumber,
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad Case - no input provided",
			args: args{
				ctx: ctx,
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - Fail to get security question by ID",
			args: args{
				ctx: ctx,
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{
						{
							QuestionID:  "1234",
							Flavour:     feedlib.FlavourConsumer,
							Response:    "Nairobi",
							PhoneNumber: interserviceclient.TestUserPhoneNumber,
						},
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
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{
						{
							QuestionID:  "1234",
							Flavour:     feedlib.FlavourConsumer,
							Response:    "Nakuru",
							PhoneNumber: interserviceclient.TestUserPhoneNumber,
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - fail to get user profile by phone number",
			args: args{
				ctx: ctx,
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{
						{
							QuestionID:  "1234",
							Flavour:     feedlib.FlavourConsumer,
							Response:    "Nakuru",
							PhoneNumber: interserviceclient.TestUserPhoneNumber,
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad Case - fail if phone number is empty",
			args: args{
				ctx: ctx,
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{
						{
							QuestionID:  "1234",
							Flavour:     feedlib.FlavourConsumer,
							Response:    "Nakuru",
							PhoneNumber: interserviceclient.TestUserPhoneNumber,
						},
					},
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid: failed to verify security question response",
			args: args{
				ctx: ctx,
				responses: &dto.VerifySecurityQuestionsPayload{
					SecurityQuestionsInput: []*dto.VerifySecurityQuestionInput{
						{
							QuestionID:  "1234",
							Flavour:     feedlib.FlavourConsumer,
							Response:    "Nakuru",
							PhoneNumber: interserviceclient.TestUserPhoneNumber,
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
			fakeDB := pgMock.NewPostgresMock()
			fakeSecurity := mock.NewSecurityQuestionsUseCaseMock()

			fakeExtension := extensionMock.NewFakeExtension()
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeDB, fakeExtension)

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
			if tt.name == "invalid: failed to verify security question response" {
				fakeDB.MockUpdateIsCorrectSecurityQuestionResponseFn = func(ctx context.Context, userID string, isCorrectSecurityQuestionResponse bool) (bool, error) {
					return false, fmt.Errorf("the failed to verify security question response does not match")
				}
			}

			if tt.name == "Sad Case - fail to get user profile by phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone")
				}
			}

			if tt.name == "Sad Case - fail if phone number is empty" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile by phone")
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

func TestUseCaseSecurityQuestionsImpl_GetUserRespondedSecurityQuestions(t *testing.T) {

	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input dto.GetUserRespondedSecurityQuestionsInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully get user responded security questions",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
		},

		{
			name: "Invalid: missing phone number",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					Flavour: feedlib.FlavourConsumer,
					OTP:     "123456",
				},
			},
			wantErr: true,
		},

		{
			name: "Invalid: missing flavour",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					OTP:         "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid: invalid phone number",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: "invalid",
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
			wantErr: true,
		},

		{
			name: "Invalid: invalid flavour",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     "invalid",
					OTP:         "123456",
				},
			},
			wantErr: true,
		},

		{
			name: "Invalid: failed to get user profile by phone number",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid: failed to verify OTP",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid: failed to get security question responses",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
			wantErr: true,
		},

		{
			name: "Invalid: security question responses is less than 3",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
			wantErr: true,
		},

		{
			name: "Invalid: failed to get a security question",
			args: args{
				ctx: ctx,
				input: dto.GetUserRespondedSecurityQuestionsInput{
					PhoneNumber: gofakeit.Phone(),
					Flavour:     feedlib.FlavourConsumer,
					OTP:         "123456",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()

			fakeExtension := extensionMock.NewFakeExtension()
			s := securityquestions.NewSecurityQuestionsUsecase(fakeDB, fakeDB, fakeDB, fakeExtension)

			fakeSecurityQuestions := mock.NewSecurityQuestionsUseCaseMock()

			if tt.name == "Invalid: missing phone number" {
				fakeSecurityQuestions.MockGetUserRespondedSecurityQuestionsFn = func(ctx context.Context, userID string, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("missing phone number")
				}
			}

			if tt.name == "Invalid: failed to get user profile by phone number" {
				fakeDB.MockGetUserProfileByPhoneNumberFn = func(ctx context.Context, phoneNumber string, flavour feedlib.Flavour) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "Invalid: failed to verify OTP" {
				fakeDB.MockVerifyOTPFn = func(ctx context.Context, payload *dto.VerifyOTPInput) (bool, error) {
					return false, fmt.Errorf("failed to verify OTP")
				}
			}

			if tt.name == "Invalid: failed to get security question responses" {

				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return nil, fmt.Errorf("failed to get responded security questions")
				}
			}

			if tt.name == "Invalid: security question responses is less than 3" {
				fakeDB.MockGetUserSecurityQuestionsResponsesFn = func(ctx context.Context, userID string) ([]*domain.SecurityQuestionResponse, error) {
					return []*domain.SecurityQuestionResponse{
						{
							ResponseID: "1234",
							QuestionID: "1234",
							Active:     true,
							Response:   "Yes",
							IsCorrect:  true,
						},
					}, nil
				}
			}

			if tt.name == "Invalid: failed to get a security question" {
				fakeDB.MockGetSecurityQuestionByIDFn = func(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
					return nil, fmt.Errorf("failed to get a security question")
				}
			}

			got, err := s.GetUserRespondedSecurityQuestions(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseSecurityQuestionsImpl.GetUserRespondedSecurityQuestions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}
}
