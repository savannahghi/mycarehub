package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

func (r *mutationResolver) ConfirmEmail(ctx context.Context, email string) (*base.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ConfirmEmail(ctx, email)
}

func (r *mutationResolver) AcceptTermsAndConditions(ctx context.Context, accept bool) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AcceptTermsAndConditions(ctx, accept)
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input profile.UserProfileInput) (*base.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateUserProfile(ctx, input)
}

func (r *mutationResolver) PractitionerSignUp(ctx context.Context, input profile.PractitionerSignupInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.PractitionerSignUp(ctx, input)
}

func (r *mutationResolver) UpdateBiodata(ctx context.Context, input profile.BiodataInput) (*base.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateBiodata(ctx, input)
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RegisterPushToken(ctx, token)
}

func (r *mutationResolver) CompleteSignup(ctx context.Context) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CompleteSignup(ctx)
}

func (r *mutationResolver) RecordPostVisitSurvey(ctx context.Context, input profile.PostVisitSurveyInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RecordPostVisitSurvey(ctx, input)
}

func (r *mutationResolver) AddTester(ctx context.Context, email string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddTester(ctx, email)
}

func (r *mutationResolver) RemoveTester(ctx context.Context, email string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RemoveTester(ctx, email)
}

func (r *mutationResolver) SetUserPin(ctx context.Context, msisdn string, pin string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SetUserPIN(ctx, msisdn, pin)
}

func (r *mutationResolver) UpdateUserPin(ctx context.Context, msisdn string, pin string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateUserPIN(ctx, msisdn, pin)
}

func (r *mutationResolver) SetLanguagePreference(ctx context.Context, language base.Language) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SetLanguagePreference(ctx, language)
}

func (r *mutationResolver) VerifyEmailOtp(ctx context.Context, email string, otp string, flavour base.Flavour) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.VerifyEmailOtp(ctx, email, otp, flavour)
}

func (r *mutationResolver) CreateSignUpMethod(ctx context.Context, signUpMethod profile.SignUpMethod) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CreateSignUpMethod(ctx, signUpMethod)
}

func (r *mutationResolver) AddPractitionerServices(ctx context.Context, services profile.PractitionerServiceInput, otherServices *profile.OtherPractitionerServiceInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddPractitionerServices(ctx, services, otherServices)
}

func (r *mutationResolver) AddPartnerType(ctx context.Context, name string, partnerType profile.PartnerType) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddPartnerType(ctx, &name, &partnerType)
}

func (r *mutationResolver) SuspendSupplier(ctx context.Context, uid string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SuspendSupplier(ctx, uid)
}

func (r *mutationResolver) SetUpSupplier(ctx context.Context, accountType profile.AccountType) (*profile.Supplier, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SetUpSupplier(ctx, accountType)
}

func (r *mutationResolver) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*profile.BranchConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SupplierEDILogin(ctx, username, password, sladeCode)
}

func (r *mutationResolver) SupplierSetDefaultLocation(ctx context.Context, locatonID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SupplierSetDefaultLocation(ctx, locatonID)
}

func (r *mutationResolver) AddIndividualRiderKyc(ctx context.Context, input profile.IndividualRider) (*profile.IndividualRider, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddIndividualRiderKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationRiderKyc(ctx context.Context, input profile.OrganizationRider) (*profile.OrganizationRider, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddOrganizationRiderKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualPractitionerKyc(ctx context.Context, input profile.IndividualPractitioner) (*profile.IndividualPractitioner, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddIndividualPractitionerKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationPractitionerKyc(ctx context.Context, input profile.OrganizationPractitioner) (*profile.OrganizationPractitioner, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddOrganizationPractitionerKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationProviderKyc(ctx context.Context, input profile.OrganizationProvider) (*profile.OrganizationProvider, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddOrganizationProviderKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualPharmaceuticalKyc(ctx context.Context, input profile.IndividualPharmaceutical) (*profile.IndividualPharmaceutical, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddIndividualPharmaceuticalKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationPharmaceuticalKyc(ctx context.Context, input profile.OrganizationPharmaceutical) (*profile.OrganizationPharmaceutical, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddOrganizationPharmaceuticalKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualCoachKyc(ctx context.Context, input profile.IndividualCoach) (*profile.IndividualCoach, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddIndividualCoachKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationCoachKyc(ctx context.Context, input profile.OrganizationCoach) (*profile.OrganizationCoach, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddOrganizationCoachKyc(ctx, input)
}

func (r *mutationResolver) AddIndividualNutritionKyc(ctx context.Context, input profile.IndividualNutrition) (*profile.IndividualNutrition, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddIndividualNutritionKyc(ctx, input)
}

func (r *mutationResolver) AddOrganizationNutritionKyc(ctx context.Context, input profile.OrganizationNutrition) (*profile.OrganizationNutrition, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddOrganizationNutritionKyc(ctx, input)
}

func (r *mutationResolver) ProcessKYCRequest(ctx context.Context, id string, status profile.KYCProcessStatus, rejectionReason *string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ProcessKYCRequest(ctx, id, status, rejectionReason)
}

func (r *queryResolver) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UserProfile(ctx)
}

func (r *queryResolver) ListTesters(ctx context.Context) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ListTesters(ctx)
}

func (r *queryResolver) RequestPinReset(ctx context.Context, msisdn string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RequestPINReset(ctx, msisdn)
}

func (r *queryResolver) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*profile.BusinessPartnerFilterInput, sort []*profile.BusinessPartnerSortInput) (*profile.BusinessPartnerConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.FindProvider(ctx, pagination, filter, sort)
}

func (r *queryResolver) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*profile.BranchFilterInput, sort []*profile.BranchSortInput) (*profile.BranchConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.FindBranch(ctx, pagination, filter, sort)
}

func (r *queryResolver) FetchSupplierAllowedLocations(ctx context.Context) (*profile.BranchConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.FetchSupplierAllowedLocations(ctx)
}

func (r *queryResolver) FetchKYCProcessingRequests(ctx context.Context) ([]*profile.KYCRequest, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.FetchKYCProcessingRequests(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *queryResolver) SupplierProfile(ctx context.Context, uid string) (*profile.Supplier, error) {
	panic(fmt.Errorf("not implemented"))
}
