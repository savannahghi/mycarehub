package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

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

func (r *mutationResolver) UpdateUserPin(ctx context.Context, phone string, pin string) (bool, error) {
	return r.interactor.UserPIN.ChangeUserPIN(ctx, phone, pin)
}

func (r *mutationResolver) SetPrimaryPhoneNumber(ctx context.Context, phone string, otp string) (bool, error) {
	if err := r.interactor.Onboarding.SetPrimaryPhoneNumber(ctx, phone, otp, true); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) SetPrimaryEmailAddress(ctx context.Context, email string, otp string) (bool, error) {
	if err := r.interactor.Onboarding.SetPrimaryEmailAddress(ctx, email, otp); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) AddSecondaryPhoneNumber(ctx context.Context, phone []string) (bool, error) {
	if err := r.interactor.Onboarding.UpdateSecondaryPhoneNumbers(ctx, phone); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RetireSecondaryPhoneNumbers(ctx context.Context, phones []string) (bool, error) {
	return r.interactor.Onboarding.RetireSecondaryPhoneNumbers(ctx, phones)
}

func (r *mutationResolver) AddSecondaryEmailAddress(ctx context.Context, email []string) (bool, error) {
	if err := r.interactor.Onboarding.UpdateSecondaryEmailAddresses(ctx, email); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RetireSecondaryEmailAddresses(ctx context.Context, emails []string) (bool, error) {
	return r.interactor.Onboarding.RetireSecondaryEmailAddress(ctx, emails)
}

func (r *mutationResolver) UpdateUserName(ctx context.Context, username string) (bool, error) {
	if err := r.interactor.Onboarding.UpdateUserName(ctx, username); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	return r.interactor.Signup.RegisterPushToken(ctx, token)
}

func (r *mutationResolver) AddPartnerType(ctx context.Context, name string, partnerType base.PartnerType) (bool, error) {
	return r.interactor.Supplier.AddPartnerType(ctx, &name, &partnerType)
}

func (r *mutationResolver) SuspendSupplier(ctx context.Context) (bool, error) {
	return r.interactor.Supplier.SuspendSupplier(ctx)
}

func (r *mutationResolver) SetUpSupplier(ctx context.Context, accountType base.AccountType) (*base.Supplier, error) {
	return r.interactor.Supplier.SetUpSupplier(ctx, accountType)
}

func (r *mutationResolver) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*resources.SupplierLogin, error) {
	return r.interactor.Supplier.SupplierEDILogin(ctx, username, password, sladeCode)
}

func (r *mutationResolver) SupplierSetDefaultLocation(ctx context.Context, locatonID string) (*base.Supplier, error) {
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

func (r *mutationResolver) RetireKYCProcessingRequest(ctx context.Context) (bool, error) {
	if err := r.interactor.Supplier.RetireKYCRequest(ctx); err != nil {
		return false, err
	}
	return true, nil
}

func (r *mutationResolver) SetupAsExperimentParticipant(ctx context.Context, participate *bool) (bool, error) {
	return r.interactor.Onboarding.SetupAsExperimentParticipant(ctx, participate)
}

func (r *mutationResolver) AddNHIFDetails(ctx context.Context, input resources.NHIFDetailsInput) (*domain.NHIFDetails, error) {
	return r.interactor.NHIF.AddNHIFDetails(ctx, input)
}

func (r *mutationResolver) AddAddress(ctx context.Context, input resources.UserAddressInput, addressType base.AddressType) (*base.Address, error) {
	return r.interactor.Onboarding.AddAddress(
		ctx,
		input,
		addressType,
	)
}

func (r *mutationResolver) SetUserCommunicationsSettings(ctx context.Context, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
	return r.interactor.Onboarding.SetUserCommunicationsSettings(ctx, allowWhatsApp, allowTextSms, allowPush, allowEmail)
}

func (r *queryResolver) UserProfile(ctx context.Context) (*base.UserProfile, error) {
	startTime := time.Now()

	userProfile, err := r.interactor.Onboarding.UserProfile(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "userProfile", err)

	return userProfile, err
}

func (r *queryResolver) SupplierProfile(ctx context.Context) (*base.Supplier, error) {
	startTime := time.Now()

	supplier, err := r.interactor.Supplier.FindSupplierByUID(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "supplierProfile", err)

	return supplier, err
}

func (r *queryResolver) ResumeWithPin(ctx context.Context, pin string) (bool, error) {
	return r.interactor.Login.ResumeWithPin(ctx, pin)
}

func (r *queryResolver) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BusinessPartnerFilterInput, sort []*resources.BusinessPartnerSortInput) (*resources.BusinessPartnerConnection, error) {
	startTime := time.Now()

	provider, err := r.interactor.ChargeMaster.FindProvider(ctx, pagination, filter, sort)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "findProvider", err)

	return provider, err
}

func (r *queryResolver) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BranchFilterInput, sort []*resources.BranchSortInput) (*resources.BranchConnection, error) {
	startTime := time.Now()

	branch, err := r.interactor.ChargeMaster.FindBranch(ctx, pagination, filter, sort)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "findBranch", err)

	return branch, err
}

func (r *queryResolver) FetchSupplierAllowedLocations(ctx context.Context) (*resources.BranchConnection, error) {
	startTime := time.Now()

	supplierAllowedLocations, err := r.interactor.Supplier.FetchSupplierAllowedLocations(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "fetchSupplierAllowedLocations", err)

	return supplierAllowedLocations, err
}

func (r *queryResolver) FetchKYCProcessingRequests(ctx context.Context) ([]*domain.KYCRequest, error) {
	startTime := time.Now()

	kycProcessingRequests, err := r.interactor.Supplier.FetchKYCProcessingRequests(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "fetchKYCProcessingRequests", err)

	return kycProcessingRequests, err
}

func (r *queryResolver) GetAddresses(ctx context.Context) (*domain.UserAddresses, error) {
	startTime := time.Now()

	addresses, err := r.interactor.Onboarding.GetAddresses(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "getAddresses", err)

	return addresses, err
}

func (r *queryResolver) NHIFDetails(ctx context.Context) (*domain.NHIFDetails, error) {
	return r.interactor.NHIF.NHIFDetails(ctx)
}

func (r *queryResolver) GetUserCommunicationsSettings(ctx context.Context) (*base.UserCommunicationsSetting, error) {
	return r.interactor.Onboarding.GetUserCommunicationsSettings(ctx)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
