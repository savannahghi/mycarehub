package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/scalarutils"
)

func (r *mutationResolver) CompleteSignup(ctx context.Context, flavour feedlib.Flavour) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, input dto.UserProfileInput) (*profileutils.UserProfile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUserPin(ctx context.Context, phone string, pin string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPrimaryPhoneNumber(ctx context.Context, phone string, otp string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetPrimaryEmailAddress(ctx context.Context, email string, otp string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddSecondaryPhoneNumber(ctx context.Context, phone []string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RetireSecondaryPhoneNumbers(ctx context.Context, phones []string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddSecondaryEmailAddress(ctx context.Context, email []string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RetireSecondaryEmailAddresses(ctx context.Context, emails []string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateUserName(ctx context.Context, username string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RecordPostVisitSurvey(ctx context.Context, input dto.PostVisitSurveyInput) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetupAsExperimentParticipant(ctx context.Context, participate *bool) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddAddress(ctx context.Context, input dto.UserAddressInput, addressType enumutils.AddressType) (*profileutils.Address, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SetUserCommunicationsSettings(ctx context.Context, allowWhatsApp *bool, allowTextSms *bool, allowPush *bool, allowEmail *bool) (*profileutils.UserCommunicationsSetting, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) SaveFavoriteNavAction(ctx context.Context, title string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteFavoriteNavAction(ctx context.Context, title string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RegisterMicroservice(ctx context.Context, input domain.Microservice) (*domain.Microservice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeregisterMicroservice(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeregisterAllMicroservices(ctx context.Context) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateRole(ctx context.Context, input dto.RoleInput) (*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteRole(ctx context.Context, roleID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AddPermissionsToRole(ctx context.Context, input dto.RolePermissionInput) (*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RevokeRolePermission(ctx context.Context, input dto.RolePermissionInput) (*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateRolePermissions(ctx context.Context, input dto.RolePermissionInput) (*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignRole(ctx context.Context, userID string, roleID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) AssignMultipleRoles(ctx context.Context, userID string, roleIDs []string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RevokeRole(ctx context.Context, userID string, roleID string, reason string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) ActivateRole(ctx context.Context, roleID string) (*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeactivateRole(ctx context.Context, roleID string) (*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) DummyQuery(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) UserProfile(ctx context.Context) (*profileutils.UserProfile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ResumeWithPin(ctx context.Context, pin string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAddresses(ctx context.Context) (*domain.UserAddresses, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetUserCommunicationsSettings(ctx context.Context) (*profileutils.UserCommunicationsSetting, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FetchUserNavigationActions(ctx context.Context) (*profileutils.NavigationActions, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) ListMicroservices(ctx context.Context) ([]*domain.Microservice, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAllRoles(ctx context.Context) ([]*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindRoleByName(ctx context.Context, roleName *string) ([]*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetAllPermissions(ctx context.Context) ([]*profileutils.Permission, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindUserByPhone(ctx context.Context, phoneNumber string) (*profileutils.UserProfile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) FindUsersByPhone(ctx context.Context, phoneNumber string) ([]*profileutils.UserProfile, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) GetNavigationActions(ctx context.Context) (*dto.GroupedNavigationActions, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *userProfileResolver) RoleDetails(ctx context.Context, obj *profileutils.UserProfile) ([]*dto.RoleOutput, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *verifiedIdentifierResolver) Timestamp(ctx context.Context, obj *profileutils.VerifiedIdentifier) (*scalarutils.Date, error) {
	panic(fmt.Errorf("not implemented"))
}

// UserProfile returns generated.UserProfileResolver implementation.
func (r *Resolver) UserProfile() generated.UserProfileResolver { return &userProfileResolver{r} }

// VerifiedIdentifier returns generated.VerifiedIdentifierResolver implementation.
func (r *Resolver) VerifiedIdentifier() generated.VerifiedIdentifierResolver {
	return &verifiedIdentifierResolver{r}
}

type userProfileResolver struct{ *Resolver }
type verifiedIdentifierResolver struct{ *Resolver }
