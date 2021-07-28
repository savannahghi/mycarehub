package usecases

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/pubsubtools"
)

// RoleUseCase represent roles interface
type RoleUseCase interface {
	CreateRole(ctx context.Context, input dto.RoleInput) (*dto.RoleOutput, error)

	AddPermissionsToRole(
		ctx context.Context,
		input dto.RolePermissionInput,
	) (*dto.RoleOutput, error)

	CheckUserHasPermission(
		ctx context.Context,
		UID string,
		permission profileutils.Permission,
	) (bool, error)

	GetUserRoles(ctx context.Context, UID string) (*[]profileutils.Role, error)

	// GetRole(ctx context.Context, ID string) (*profileutils.Role, error)

	// GetUserPermissions(ctx context.Context, UID string) ([]profileutils.Permission, error)
}

// NewRoleUseCases returns a new a onboarding usecase
func NewRoleUseCases(
	r repository.OnboardingRepository,
	ext extension.BaseExtension,
) RoleUseCase {
	return &RoleUseCaseImpl{
		repo:    r,
		baseExt: ext,
	}
}

// RoleUseCaseImpl  represents usecase implementation object
type RoleUseCaseImpl struct {
	repo    repository.OnboardingRepository
	baseExt extension.BaseExtension
}

// CreateRole creates a new Role
func (r *RoleUseCaseImpl) CreateRole(
	ctx context.Context,
	input dto.RoleInput,
) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "CreateRole")
	defer span.End()

	user, err := r.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	userProfile, err := r.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// Check logged in user has the right permissions
	allowed, err := r.CheckUserHasPermission(ctx, user.UID, profileutils.CanCreateRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !allowed {
		return nil, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to create role"),
		)
	}

	timestamp := time.Now().In(pubsubtools.TimeLocation)
	// move some logic to repository
	// ID,timestamp,active
	role := profileutils.Role{
		ID:          uuid.New().String(),
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   userProfile.ID,
		Created:     timestamp,
		Active:      true,
		Scopes:      input.Scopes,
	}

	// save role to database
	newRole, err := r.repo.CreateRole(ctx, role)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// get permissions
	perms, err := role.Permissions(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	output := &dto.RoleOutput{
		ID:          newRole.ID,
		Name:        newRole.Name,
		Description: newRole.Description,
		Scopes:      newRole.Scopes,
		Permissions: perms,
	}

	return output, nil
}

// AddPermissionsToRole add permission operation to a role
func (r *RoleUseCaseImpl) AddPermissionsToRole(
	ctx context.Context,
	input dto.RolePermissionInput,
) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "AddPermissionsToRole")
	defer span.End()

	user, err := r.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// Check logged in user has the right permissions
	allowed, err := r.CheckUserHasPermission(ctx, user.UID, profileutils.CanEditRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !allowed {
		return nil, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to edit role"),
		)
	}

	role, err := r.repo.GetRoleByID(ctx, input.RoleID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	permission, err := profileutils.GetPermissionByScope(ctx, input.Scope)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	if !role.HasPermission(ctx, permission.Scope) {
		role.Scopes = append(role.Scopes, permission.Scope)
	}

	// save role to database
	err = r.repo.UpdateRoleDetails(ctx, *role)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// get permissions
	perms, err := role.Permissions(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	output := &dto.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Scopes:      role.Scopes,
		Permissions: perms,
	}

	return output, nil
}

// CheckUserHasPermission checks if the user has a permission
func (r *RoleUseCaseImpl) CheckUserHasPermission(
	ctx context.Context,
	UID string,
	requiredPermission profileutils.Permission,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "CheckUserHasPermission")
	defer span.End()

	roles, err := r.GetUserRoles(ctx, UID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	perms := []profileutils.Permission{}

	for _, role := range *roles {
		rolePerms, err := role.Permissions(ctx)
		if err != nil {
			perms = append(perms, rolePerms...)
		}
	}

	for _, permission := range perms {
		if permission == requiredPermission {
			return true, nil
		}
	}

	return false, nil
}

// GetUserRoles returns a list of user roles
func (r *RoleUseCaseImpl) GetUserRoles(
	ctx context.Context,
	UID string,
) (*[]profileutils.Role, error) {
	ctx, span := tracer.Start(ctx, "GetUserRoles")
	defer span.End()

	userProfile, err := r.repo.GetUserProfileByUID(ctx, UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	roles, err := r.repo.GetRolesByIDs(ctx, userProfile.Roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
