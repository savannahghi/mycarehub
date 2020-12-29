package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
)

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input domain.UserProfileInput) (*base.UserProfile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddPartnerType(ctx context.Context, name string, partnerType domain.PartnerType) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SuspendSupplier(ctx context.Context, uid string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetUpSupplier(ctx context.Context, accountType domain.AccountType) (*domain.Supplier, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*domain.BranchConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SupplierSetDefaultLocation(ctx context.Context, locatonID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ProcessKYCRequest(ctx context.Context, id string, status domain.KYCProcessStatus, rejectionReason *string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*domain.BusinessPartnerFilterInput, sort []*domain.BusinessPartnerSortInput) (*domain.BusinessPartnerConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*domain.BranchFilterInput, sort []*domain.BranchSortInput) (*domain.BranchConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FetchSupplierAllowedLocations(ctx context.Context) (*domain.BranchConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
