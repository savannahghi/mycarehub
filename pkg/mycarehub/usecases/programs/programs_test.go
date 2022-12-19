package programs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/programs"
)

func TestUsecaseProgramsImpl_CreateProgram(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		input *dto.ProgramInput
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case: create program",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerHop(),
					OrganisationID: uuid.NewString(),
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case: invalid input",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					OrganisationID: uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to check organization exists",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerHop(),
					OrganisationID: uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: organization does not exists",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerHop(),
					OrganisationID: uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to check if organization has program",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerHop(),
					OrganisationID: uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: organization already has a program",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerHop(),
					OrganisationID: uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case: failed to create program",
			args: args{
				ctx: ctx,
				input: &dto.ProgramInput{
					Name:           gofakeit.BeerHop(),
					OrganisationID: uuid.NewString(),
				},
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			u := programs.NewUsecasePrograms(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "Sad case: failed to check organization exists" {
				fakeDB.MockCheckOrganisationExistsFn = func(ctx context.Context, organisationID string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}
			if tt.name == "Sad case: organization does not exists" {
				fakeDB.MockCheckOrganisationExistsFn = func(ctx context.Context, organisationID string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "Sad case: failed to check if organization has program" {
				fakeDB.MockCheckIfProgramNameExistsFn = func(ctx context.Context, organisationID string, programName string) (bool, error) {
					return false, fmt.Errorf("an error occurred")
				}
			}

			if tt.name == "Sad case: organization already has a program" {
				fakeDB.MockCheckIfProgramNameExistsFn = func(ctx context.Context, organisationID string, programName string) (bool, error) {
					return true, nil
				}
			}
			if tt.name == "Sad case: failed to create program" {
				fakeDB.MockCreateProgramFn = func(ctx context.Context, program *dto.ProgramInput) error {
					return fmt.Errorf("an error occurred")
				}
			}

			got, err := u.CreateProgram(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseProgramsImpl.CreateProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseProgramsImpl.CreateProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseProgramsImpl_SetCurrentProgram(t *testing.T) {

	type args struct {
		ctx       context.Context
		programID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "sad case: fail to get logged in user",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to get user profile",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "sad case: fail to update user profile",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "happy case: set current program",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			u := programs.NewUsecasePrograms(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: fail to get logged in user" {
				fakeExtension.MockGetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			if tt.name == "sad case: fail to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad case: fail to update user profile" {
				fakeDB.MockUpdateUserFn = func(ctx context.Context, user *domain.User, updateData map[string]interface{}) error {
					return fmt.Errorf("failed to update user")
				}
			}

			got, err := u.SetCurrentProgram(tt.args.ctx, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseProgramsImpl.SetCurrentProgram() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UsecaseProgramsImpl.SetCurrentProgram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUsecaseProgramsImpl_ListUserPrograms(t *testing.T) {

	type args struct {
		ctx     context.Context
		userID  string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad case: fail to get user profile",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get user programs, pro",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
			},
			wantErr: true,
		},
		{
			name: "sad case: fail to get user programs, consumer",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: true,
		},
		{
			name: "happy case: get user programs, pro",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourPro,
			},
			wantErr: false,
		},
		{
			name: "happy case: get user programs, consumer",
			args: args{
				ctx:     context.Background(),
				userID:  gofakeit.UUID(),
				flavour: feedlib.FlavourConsumer,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			u := programs.NewUsecasePrograms(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: fail to get user profile" {
				fakeDB.MockGetUserProfileByUserIDFn = func(ctx context.Context, userID string) (*domain.User, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad case: fail to get user programs, pro" {
				fakeDB.MockGetStaffUserProgramsFn = func(ctx context.Context, userID string) ([]*domain.Program, error) {
					return nil, fmt.Errorf("failed to get user programs")
				}
			}
			if tt.name == "sad case: fail to get user programs, consumer" {
				fakeDB.MockGetClientUserProgramsFn = func(ctx context.Context, userID string) ([]*domain.Program, error) {
					return nil, fmt.Errorf("failed to get user programs")
				}
			}

			_, err := u.ListUserPrograms(tt.args.ctx, tt.args.userID, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseProgramsImpl.ListUserPrograms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUsecaseProgramsImpl_GetProgramFacilities(t *testing.T) {
	type args struct {
		ctx       context.Context
		programID string
	}
	tests := []struct {
		name    string
		args    args
		want    []*domain.Facility
		wantErr bool
	}{
		{
			name: "happy case: get program facilities",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			want:    []*domain.Facility{},
			wantErr: false,
		},
		{
			name: "sad case: unable to get program facilities",
			args: args{
				ctx:       context.Background(),
				programID: gofakeit.UUID(),
			},
			want:    []*domain.Facility{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			u := programs.NewUsecasePrograms(fakeDB, fakeDB, fakeDB, fakeExtension)

			if tt.name == "sad case: unable to get program facilities" {
				fakeDB.MockGetProgramFacilitiesFn = func(ctx context.Context, programID string) ([]*domain.Facility, error) {
					return nil, fmt.Errorf("failed to get program facilities")
				}
			}

			got, err := u.GetProgramFacilities(tt.args.ctx, tt.args.programID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UsecaseProgramsImpl.GetProgramFacilities() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("did nox expect error, got %v", got)
			}
		})
	}
}
