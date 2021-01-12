package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
)

func (r *mutationResolver) CompleteSignup(ctx context.Context, flavour base.Flavour) (bool, error) {
	return r.interactor.Signup.CompleteSignup(ctx, flavour)
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input resources.UserProfileInput) (*base.UserProfile, error) {
	return r.interactor.Signup.UpdateUserProfile(ctx, &input)
}

func (r *mutationResolver) UpdateUserPin(ctx context.Context, phone string, pin string) (*resources.PINOutput, error) {
	return r.interactor.UserPIN.ChangeUserPIN(ctx, phone, pin)
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	return r.interactor.Signup.RegisterPushToken(ctx, token)
}

func (r *mutationResolver) AddPartnerType(ctx context.Context, name string, partnerType domain.PartnerType) (bool, error) {
	return r.interactor.Supplier.AddPartnerType(ctx, &name, &partnerType)
}

func (r *mutationResolver) SuspendSupplier(ctx context.Context) (bool, error) {
	return r.interactor.Supplier.SuspendSupplier(ctx)
}

func (r *mutationResolver) SetUpSupplier(ctx context.Context, accountType domain.AccountType) (*domain.Supplier, error) {
	return r.interactor.Supplier.SetUpSupplier(ctx, accountType)
}

func (r *mutationResolver) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*resources.BranchConnection, error) {
	return r.interactor.Supplier.SupplierEDILogin(ctx, username, password, sladeCode)
}

func (r *mutationResolver) SupplierSetDefaultLocation(ctx context.Context, locatonID string) (bool, error) {
	return r.interactor.Supplier.SupplierSetDefaultLocation(ctx, locatonID)
}

func (r *mutationResolver) AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error) {
	return r.interactor.Supplier.AddIndividualRiderKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error) {
	return r.interactor.Supplier.AddOrganizationRiderKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error) {
	return r.interactor.Supplier.AddIndividualPractitionerKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error) {
	return r.interactor.Supplier.AddOrganizationPractitionerKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error) {
	return r.interactor.Supplier.AddOrganizationProviderKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error) {
	return r.interactor.Supplier.AddIndividualPharmaceuticalKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error) {
	return r.interactor.Supplier.AddOrganizationPharmaceuticalKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error) {
	return r.interactor.Supplier.AddIndividualCoachKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error) {
	return r.interactor.Supplier.AddOrganizationCoachKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error) {
	return r.interactor.Supplier.AddIndividualNutritionKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error) {
	return r.interactor.Supplier.AddOrganizationNutritionKyc(ctx, input)
}

func (r *mutationResolver) ProcessKYCRequest(ctx context.Context, id string, status domain.KYCProcessStatus, rejectionReason *string) (bool, error) {
	return r.interactor.Supplier.ProcessKYCRequest(ctx, id, status, rejectionReason)
}

func (r *mutationResolver) RecordPostVisitSurvey(ctx context.Context, input resources.PostVisitSurveyInput) (bool, error) {
	return r.interactor.Survey.RecordPostVisitSurvey(ctx, input)
}

func (r *queryResolver) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	return r.interactor.Onboarding.UserProfile(ctx)
}

func (r *queryResolver) SupplierProfile(ctx context.Context) (*domain.Supplier, error) {
	return r.interactor.Supplier.FindSupplierByUID(ctx)
}

func (r *queryResolver) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BusinessPartnerFilterInput, sort []*resources.BusinessPartnerSortInput) (*resources.BusinessPartnerConnection, error) {
	return r.interactor.ChargeMaster.FindProvider(ctx, pagination, filter, sort)
}

func (r *queryResolver) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BranchFilterInput, sort []*resources.BranchSortInput) (*resources.BranchConnection, error) {
	return r.interactor.ChargeMaster.FindBranch(ctx, pagination, filter, sort)
}

func (r *queryResolver) FetchSupplierAllowedLocations(ctx context.Context) (*resources.BranchConnection, error) {
	return r.interactor.Supplier.FetchSupplierAllowedLocations(ctx)
}

func (r *queryResolver) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	return r.interactor.Supplier.FetchKYCProcessingRequests(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
