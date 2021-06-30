package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"time"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/graph/generated"
)

func (r *mutationResolver) CompleteSignup(ctx context.Context, flavour base.Flavour) (bool, error) {
	startTime := time.Now()

	completeSignup, err := r.interactor.Signup.CompleteSignup(ctx, flavour)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "completeSignup", err)

	return completeSignup, err
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input dto.UserProfileInput) (*base.UserProfile, error) {
	startTime := time.Now()

	updateUserProfile, err := r.interactor.Signup.UpdateUserProfile(ctx, &input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "updateUserProfile", err)

	return updateUserProfile, err
}

func (r *mutationResolver) UpdateUserPin(ctx context.Context, phone string, pin string) (bool, error) {
	startTime := time.Now()

	updateUserPIN, err := r.interactor.UserPIN.ChangeUserPIN(ctx, phone, pin)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "updateUserPIN", err)

	return updateUserPIN, err
}

func (r *mutationResolver) SetPrimaryPhoneNumber(ctx context.Context, phone string, otp string) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Onboarding.SetPrimaryPhoneNumber(ctx, phone, otp, true)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "setPrimaryPhoneNumber", err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetPrimaryEmailAddress(ctx context.Context, email string, otp string) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Onboarding.SetPrimaryEmailAddress(ctx, email, otp)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "setPrimaryEmailAddress", err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) SetOptOut(ctx context.Context, option string, phoneNumber string) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Onboarding.SetOptOut(ctx, option, phoneNumber)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "setOptOut", err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) AddSecondaryPhoneNumber(ctx context.Context, phone []string) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Onboarding.UpdateSecondaryPhoneNumbers(ctx, phone)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addSecondaryPhoneNumber", err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RetireSecondaryPhoneNumbers(ctx context.Context, phones []string) (bool, error) {
	startTime := time.Now()

	retireSecondaryPhoneNumbers, err := r.interactor.Onboarding.RetireSecondaryPhoneNumbers(ctx, phones)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "retireSecondaryPhoneNumbers", err)

	return retireSecondaryPhoneNumbers, err
}

func (r *mutationResolver) AddSecondaryEmailAddress(ctx context.Context, email []string) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Onboarding.UpdateSecondaryEmailAddresses(ctx, email)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addSecondaryEmailAddress", err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RetireSecondaryEmailAddresses(ctx context.Context, emails []string) (bool, error) {
	startTime := time.Now()

	retireSecondaryEmailAddresses, err := r.interactor.Onboarding.RetireSecondaryEmailAddress(ctx, emails)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "retireSecondaryEmailAddresses", err)

	return retireSecondaryEmailAddresses, err
}

func (r *mutationResolver) UpdateUserName(ctx context.Context, username string) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Onboarding.UpdateUserName(ctx, username)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "updateUserName", err)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	startTime := time.Now()

	registerPushToken, err := r.interactor.Signup.RegisterPushToken(ctx, token)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "registerPushToken", err)

	return registerPushToken, err
}

func (r *mutationResolver) AddPartnerType(ctx context.Context, name string, partnerType base.PartnerType) (bool, error) {
	startTime := time.Now()

	addPartnerType, err := r.interactor.Supplier.AddPartnerType(ctx, &name, &partnerType)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addPartnerType", err)

	return addPartnerType, err
}

func (r *mutationResolver) SuspendSupplier(ctx context.Context, suspensionReason *string) (bool, error) {
	startTime := time.Now()

	suspendSupplier, err := r.interactor.Supplier.SuspendSupplier(ctx, suspensionReason)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "suspendSupplier", err)

	return suspendSupplier, err
}

func (r *mutationResolver) SetUpSupplier(ctx context.Context, accountType base.AccountType) (*base.Supplier, error) {
	startTime := time.Now()

	supplier, err := r.interactor.Supplier.SetUpSupplier(ctx, accountType)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "setUpSupplier", err)

	return supplier, err
}

func (r *mutationResolver) SupplierEDILogin(ctx context.Context, username string, password string, sladeCode string) (*dto.SupplierLogin, error) {
	startTime := time.Now()

	supplierEDILogin, err := r.interactor.Supplier.SupplierEDILogin(ctx, username, password, sladeCode)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "supplierEDILogin", err)

	return supplierEDILogin, err
}

func (r *mutationResolver) SupplierSetDefaultLocation(ctx context.Context, locationID string) (*base.Supplier, error) {
	startTime := time.Now()

	supplier, err := r.interactor.Supplier.SupplierSetDefaultLocation(ctx, locationID)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "supplierSetDefaultLocation", err)

	return supplier, err
}

func (r *mutationResolver) AddIndividualRiderKyc(ctx context.Context, input domain.IndividualRider) (*domain.IndividualRider, error) {
	startTime := time.Now()

	individualRider, err := r.interactor.Supplier.AddIndividualRiderKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addIndividualRiderKYC", err)

	return individualRider, err
}

func (r *mutationResolver) AddOrganizationRiderKyc(ctx context.Context, input domain.OrganizationRider) (*domain.OrganizationRider, error) {
	startTime := time.Now()

	organizationRider, err := r.interactor.Supplier.AddOrganizationRiderKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addOrganizationRiderKYC", err)

	return organizationRider, err
}

func (r *mutationResolver) AddIndividualPractitionerKyc(ctx context.Context, input domain.IndividualPractitioner) (*domain.IndividualPractitioner, error) {
	startTime := time.Now()

	individualPractitioner, err := r.interactor.Supplier.AddIndividualPractitionerKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addIndividualPractitionerKYC", err)

	return individualPractitioner, err
}

func (r *mutationResolver) AddOrganizationPractitionerKyc(ctx context.Context, input domain.OrganizationPractitioner) (*domain.OrganizationPractitioner, error) {
	startTime := time.Now()

	organizationPractitioner, err := r.interactor.Supplier.AddOrganizationPractitionerKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addOrganizationPractitionerKYC", err)

	return organizationPractitioner, err
}

func (r *mutationResolver) AddOrganizationProviderKyc(ctx context.Context, input domain.OrganizationProvider) (*domain.OrganizationProvider, error) {
	startTime := time.Now()

	organizationProvider, err := r.interactor.Supplier.AddOrganizationProviderKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addOrganizationProviderKYC", err)

	return organizationProvider, err
}

func (r *mutationResolver) AddIndividualPharmaceuticalKyc(ctx context.Context, input domain.IndividualPharmaceutical) (*domain.IndividualPharmaceutical, error) {
	startTime := time.Now()

	individualPharmaceutical, err := r.interactor.Supplier.AddIndividualPharmaceuticalKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addIndividualPharmaceuticalKYC", err)

	return individualPharmaceutical, err
}

func (r *mutationResolver) AddOrganizationPharmaceuticalKyc(ctx context.Context, input domain.OrganizationPharmaceutical) (*domain.OrganizationPharmaceutical, error) {
	startTime := time.Now()

	organizationPharmaceutical, err := r.interactor.Supplier.AddOrganizationPharmaceuticalKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addOrganizationPharmaceuticalKYC", err)

	return organizationPharmaceutical, err
}

func (r *mutationResolver) AddIndividualCoachKyc(ctx context.Context, input domain.IndividualCoach) (*domain.IndividualCoach, error) {
	startTime := time.Now()

	individualCoach, err := r.interactor.Supplier.AddIndividualCoachKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addIndividualCoachKYC", err)

	return individualCoach, err
}

func (r *mutationResolver) AddOrganizationCoachKyc(ctx context.Context, input domain.OrganizationCoach) (*domain.OrganizationCoach, error) {
	startTime := time.Now()

	organizationCoach, err := r.interactor.Supplier.AddOrganizationCoachKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addOrganizationCoachKYC", err)

	return organizationCoach, err
}

func (r *mutationResolver) AddIndividualNutritionKyc(ctx context.Context, input domain.IndividualNutrition) (*domain.IndividualNutrition, error) {
	startTime := time.Now()

	individualNutrition, err := r.interactor.Supplier.AddIndividualNutritionKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addIndividualNutritionKYC", err)

	return individualNutrition, err
}

func (r *mutationResolver) AddOrganizationNutritionKyc(ctx context.Context, input domain.OrganizationNutrition) (*domain.OrganizationNutrition, error) {
	startTime := time.Now()

	organizationNutrition, err := r.interactor.Supplier.AddOrganizationNutritionKyc(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addOrganizationNutritionKYC", err)

	return organizationNutrition, err
}

func (r *mutationResolver) ProcessKYCRequest(ctx context.Context, id string, status domain.KYCProcessStatus, rejectionReason *string) (bool, error) {
	startTime := time.Now()

	processKYCRequest, err := r.interactor.Supplier.ProcessKYCRequest(ctx, id, status, rejectionReason)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "processKYCRequest", err)

	return processKYCRequest, err
}

func (r *mutationResolver) RecordPostVisitSurvey(ctx context.Context, input dto.PostVisitSurveyInput) (bool, error) {
	startTime := time.Now()

	recordPostVisitSurvey, err := r.interactor.Survey.RecordPostVisitSurvey(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "recordPostVisitSurvey", err)

	return recordPostVisitSurvey, err
}

func (r *mutationResolver) RetireKYCProcessingRequest(ctx context.Context) (bool, error) {
	startTime := time.Now()

	err := r.interactor.Supplier.RetireKYCRequest(ctx)

	if err != nil {
		return false, err
	}

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "retireKYCProcessingRequest", err)

	return true, nil
}

func (r *mutationResolver) SetupAsExperimentParticipant(ctx context.Context, participate *bool) (bool, error) {
	startTime := time.Now()

	setupAsExperimentParticipant, err := r.interactor.Onboarding.SetupAsExperimentParticipant(ctx, participate)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "setupAsExperimentParticipant", err)

	return setupAsExperimentParticipant, err
}

func (r *mutationResolver) AddNHIFDetails(ctx context.Context, input dto.NHIFDetailsInput) (*domain.NHIFDetails, error) {
	startTime := time.Now()

	addNHIFDetails, err := r.interactor.NHIF.AddNHIFDetails(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addNHIFDetails", err)

	return addNHIFDetails, err
}

func (r *mutationResolver) AddAddress(ctx context.Context, input dto.UserAddressInput, addressType base.AddressType) (*base.Address, error) {
	startTime := time.Now()

	addAddress, err := r.interactor.Onboarding.AddAddress(
		ctx,
		input,
		addressType,
	)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "addAddress", err)

	return addAddress, err
}

func (r *mutationResolver) SetUserCommunicationsSettings(ctx context.Context, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*base.UserCommunicationsSetting, error) {
	startTime := time.Now()

	setUserCommunicationsSettings, err := r.interactor.Onboarding.SetUserCommunicationsSettings(ctx, allowWhatsApp, allowTextSms, allowPush, allowEmail)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "setUserCommunicationsSettings", err)

	return setUserCommunicationsSettings, err
}

func (r *mutationResolver) RegisterAgent(ctx context.Context, input dto.RegisterAgentInput) (*base.UserProfile, error) {
	startTime := time.Now()

	userProfile, err := r.interactor.Agent.RegisterAgent(ctx, input)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "registerAgent", err)

	return userProfile, err
}

func (r *mutationResolver) ActivateAgent(ctx context.Context, agentID string) (bool, error) {
	startTime := time.Now()

	success, err := r.interactor.Agent.ActivateAgent(ctx, agentID)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "activateAgent", err)

	return success, err
}

func (r *mutationResolver) DeactivateAgent(ctx context.Context, agentID string) (bool, error) {
	startTime := time.Now()

	success, err := r.interactor.Agent.DeactivateAgent(ctx, agentID)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "deactivateAgent", err)

	return success, err
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
	startTime := time.Now()

	resumeWithPin, err := r.interactor.Login.ResumeWithPin(ctx, pin)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "resumeWithPin", err)

	return resumeWithPin, err
}

func (r *queryResolver) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput, sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
	startTime := time.Now()

	provider, err := r.interactor.ChargeMaster.FindProvider(ctx, pagination, filter, sort)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "findProvider", err)

	return provider, err
}

func (r *queryResolver) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput, sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
	startTime := time.Now()

	branch, err := r.interactor.ChargeMaster.FindBranch(ctx, pagination, filter, sort)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "findBranch", err)

	return branch, err
}

func (r *queryResolver) FetchSupplierAllowedLocations(ctx context.Context) (*dto.BranchConnection, error) {
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
	startTime := time.Now()

	NHIFDetails, err := r.interactor.NHIF.NHIFDetails(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "NHIFDetails", err)

	return NHIFDetails, err
}

func (r *queryResolver) GetUserCommunicationsSettings(ctx context.Context) (*base.UserCommunicationsSetting, error) {
	startTime := time.Now()

	userCommunicationsSettings, err := r.interactor.Onboarding.GetUserCommunicationsSettings(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "getUserCommunicationsSettings", err)

	return userCommunicationsSettings, err
}

func (r *queryResolver) CheckSupplierKYCSubmitted(ctx context.Context) (bool, error) {
	startTime := time.Now()

	checkSupplierKYCSubmitted, err := r.interactor.Supplier.CheckSupplierKYCSubmitted(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "checkSupplierKYCSubmitted", err)

	return checkSupplierKYCSubmitted, err
}

func (r *queryResolver) FetchAgents(ctx context.Context) ([]*dto.Agent, error) {
	startTime := time.Now()

	agents, err := r.interactor.Agent.FetchAgents(ctx)

	defer base.RecordGraphqlResolverMetrics(ctx, startTime, "fetchAgents", err)

	return agents, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
