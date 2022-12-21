package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.22

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// CreateProgram is the resolver for the createProgram field.
func (r *mutationResolver) CreateProgram(ctx context.Context, input dto.ProgramInput) (bool, error) {
	return r.mycarehub.Programs.CreateProgram(ctx, &input)
}

// SetStaffProgram is the resolver for the setStaffProgram field.
func (r *mutationResolver) SetStaffProgram(ctx context.Context, programID string) (*domain.StaffProfile, error) {
	r.checkPreconditions()

	return r.mycarehub.Programs.SetStaffProgram(ctx, programID)
}

// SetClientProgram is the resolver for the setClientProgram field.
func (r *mutationResolver) SetClientProgram(ctx context.Context, programID string) (*domain.ClientProfile, error) {
	r.checkPreconditions()

	return r.mycarehub.Programs.SetClientProgram(ctx, programID)
}

// ListUserPrograms is the resolver for the listUserPrograms field.
func (r *queryResolver) ListUserPrograms(ctx context.Context, userID string, flavour feedlib.Flavour) (*dto.ProgramOutput, error) {
	return r.mycarehub.Programs.ListUserPrograms(ctx, userID, flavour)
}

// GetProgramFacilities is the resolver for the getProgramFacilities field.
func (r *queryResolver) GetProgramFacilities(ctx context.Context, programID string) ([]*domain.Facility, error) {
	return r.mycarehub.Programs.GetProgramFacilities(ctx, programID)
}
