package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/shopspring/decimal"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

func (r *mutationResolver) ConfirmEmail(ctx context.Context, email string) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ConfirmEmail(ctx, email)
}

func (r *mutationResolver) AcceptTermsAndConditions(ctx context.Context, accept bool) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AcceptTermsAndConditions(ctx, accept)
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input profile.UserProfileInput) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateUserProfile(ctx, input)
}

func (r *mutationResolver) PractitionerSignUp(ctx context.Context, input profile.PractitionerSignupInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.PractitionerSignUp(ctx, input)
}

func (r *mutationResolver) UpdateBiodata(ctx context.Context, input profile.BiodataInput) (*profile.UserProfile, error) {
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

func (r *mutationResolver) ApprovePractitionerSignup(ctx context.Context, practitionerID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ApprovePractitionerSignup(ctx)
}

func (r *mutationResolver) RejectPractitionerSignup(ctx context.Context, practitionerID string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RejectPractitionerSignup(ctx)
}

func (r *mutationResolver) SetUserPin(ctx context.Context, msisdn string, pin string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SetUserPin(ctx, msisdn, pin)
}

func (r *mutationResolver) UpdateUserPin(ctx context.Context, msisdn string, pin string, otp string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateUserPin(ctx, msisdn, pin, otp)
}

func (r *mutationResolver) SetLanguagePreference(ctx context.Context, language base.Language) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SetLanguagePreference(ctx, language)
}

func (r *mutationResolver) VerifyEmailOtp(ctx context.Context, email string, otp string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.VerifyEmailOtp(ctx, email, otp)
}

func (r *mutationResolver) CreateSignUpMethod(ctx context.Context, signUpMethod profile.SignUpMethod) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CreateSignUpMethod(ctx, signUpMethod)
}

func (r *mutationResolver) AddCustomer(ctx context.Context) (*profile.Customer, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddCustomer(ctx)
}

func (r *mutationResolver) AddCustomerKyc(ctx context.Context, input profile.CustomerKYCInput) (*profile.CustomerKYC, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddCustomerKYC(ctx, input)
}

func (r *mutationResolver) UpdateCustomer(ctx context.Context, input profile.CustomerKYCInput) (*profile.Customer, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateCustomer(ctx, input)
}

func (r *mutationResolver) AddPractitionerServices(ctx context.Context, services profile.PractitionerServiceInput, otherServices *profile.OtherPractitionerServiceInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AddPractitionerServices(ctx, services, otherServices)
}

func (r *queryResolver) UserProfile(ctx context.Context) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UserProfile(ctx)
}

func (r *queryResolver) HealthcashBalance(ctx context.Context) (*base.Decimal, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	expectedBalance := base.Decimal(decimal.NewFromFloat(0))
	return &expectedBalance, nil
}

func (r *queryResolver) GetProfile(ctx context.Context, uid string) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.GetProfile(ctx, uid)
}

func (r *queryResolver) ListTesters(ctx context.Context) ([]string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ListTesters(ctx)
}

func (r *queryResolver) ListKMPDURegisteredPractitioners(ctx context.Context, pagination *base.PaginationInput, filter *base.FilterInput, sort *base.SortInput) (*profile.KMPDUPractitionerConnection, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ListKMPDURegisteredPractitioners(ctx, pagination, filter, sort)
}

func (r *queryResolver) GetKMPDURegisteredPractitioner(ctx context.Context, regno string) (*profile.KMPDUPractitioner, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.GetRegisteredPractitionerByLicense(ctx, regno)
}

func (r *queryResolver) IsUnderAge(ctx context.Context) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.IsUnderAge(ctx)
}

func (r *queryResolver) VerifyMSISDNandPin(ctx context.Context, msisdn string, pin string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.VerifyMSISDNandPin(ctx, msisdn, pin)
}

func (r *queryResolver) RequestPinReset(ctx context.Context, msisdn string) (string, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RequestPinReset(ctx, msisdn)
}

func (r *queryResolver) CheckUserWithMsisdn(ctx context.Context, msisdn string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CheckUserWithMsisdn(ctx, msisdn)
}

func (r *queryResolver) CheckEmailVerified(ctx context.Context) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CheckEmailVerified(ctx)
}

func (r *queryResolver) CheckPhoneNumberVerified(ctx context.Context) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CheckPhoneNumberVerified(ctx)
}

func (r *queryResolver) GetSignUpMethod(ctx context.Context, id string) (profile.SignUpMethod, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.GetSignUpMethod(ctx, id)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
