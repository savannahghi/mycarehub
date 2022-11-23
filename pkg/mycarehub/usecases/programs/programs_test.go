package programs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
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
			u := programs.NewUsecasePrograms(fakeDB, fakeDB)

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
