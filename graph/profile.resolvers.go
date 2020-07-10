package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/generated"
	"gitlab.slade360emr.com/go/profile/graph/profile"
	"gitlab.slade360emr.com/go/uploads/graph/uploads"
)

func (r *queryResolver) UserProfile(ctx context.Context) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UserProfile(ctx)
}

func (r *queryResolver) HealthcashBalance(ctx context.Context) (*base.Decimal, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.HealthcashBalance(ctx)
}

func (r *queryResolver) GetPresence(ctx context.Context) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.GetPresence(ctx)
}

func (r *mutationResolver) AcceptTermsAndConditions(
	ctx context.Context, accept bool) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.AcceptTermsAndConditions(ctx, accept)
}

func (r *mutationResolver) PractitionerSignUp(
	ctx context.Context, input profile.PractitionerSignupInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.PractitionerSignUp(ctx, input)
}

func (r *mutationResolver) UpdateUserProfile(
	ctx context.Context, input profile.UserProfileInput) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateUserProfile(ctx, input)
}

func (r *mutationResolver) UpdateBiodata(ctx context.Context, input profile.BiodataInput) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.UpdateBiodata(ctx, input)
}

func (r *mutationResolver) CompleteSignup(ctx context.Context) (*base.Decimal, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.CompleteSignup(ctx)
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RegisterPushToken(ctx, token)
}

func (r *mutationResolver) RecordPostVisitSurvey(ctx context.Context, input profile.PostVisitSurveyInput) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.RecordPostVisitSurvey(ctx, input)
}

func (r *mutationResolver) SetPresence(ctx context.Context, presence bool) (bool, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.SetPresence(ctx, presence)
}

func (r *mutationResolver) ConfirmEmail(ctx context.Context, email string) (*profile.UserProfile, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.profileService.ConfirmEmail(ctx, email)
}

func (r *mutationResolver) Upload(ctx context.Context, input *uploads.UploadInput) (*uploads.Upload, error) {
	r.CheckUserTokenInContext(ctx)
	r.CheckDependencies()
	return r.uploadService.Upload(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
