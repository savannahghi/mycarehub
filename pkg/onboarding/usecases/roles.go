package usecases

import (
	"context"
	"fmt"

	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/exceptions"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/utils"
	"github.com/savannahghi/onboarding/pkg/onboarding/repository"
	"github.com/savannahghi/profileutils"
)

// RoleUseCase represent the business logic required for management of agents
type RoleUseCase interface {
	CreateRole(ctx context.Context, input dto.RoleInput) (*dto.RoleOutput, error)
	GetAllRoles(ctx context.Context) ([]*dto.RoleOutput, error)
	GetAllPermissions(ctx context.Context) ([]*profileutils.Permission, error)

	AddPermissionsToRole(
		ctx context.Context,
		input dto.RolePermissionInput,
	) (*dto.RoleOutput, error)

	GetUserRoles(ctx context.Context, UID string) (*[]profileutils.Role, error)

	GetRole(ctx context.Context, ID string) (*dto.RoleOutput, error)

	GetUserPermissions(ctx context.Context, UID string) ([]profileutils.Permission, error)
}

// RoleUseCaseImpl  represents usecase implementation object
type RoleUseCaseImpl struct {
	repo    repository.OnboardingRepository
	baseExt extension.BaseExtension
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

	// Check logged in user has the right permissions
	allowed, err := r.repo.CheckIfUserHasPermission(ctx, user.UID, profileutils.CanCreateRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !allowed {
		return nil, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to create role"),
		)
	}

	userProfile, err := r.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	role, err := r.repo.CreateRole(ctx, userProfile.ID, input)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	perms, err := role.Permissions(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	output := &dto.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Active:      role.Active,
		Scopes:      role.Scopes,
		Permissions: perms,
	}

	return output, nil
}

//GetAllRoles returns a list of all created roles
func (r *RoleUseCaseImpl) GetAllRoles(ctx context.Context) ([]*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "GetAllRoles")
	defer span.End()

	user, err := r.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	// Check logged in user has the right permissions
	allowed, err := r.repo.CheckIfUserHasPermission(ctx, user.UID, profileutils.CanViewRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !allowed {
		return nil, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to list roles"),
		)
	}

	roles, err := r.repo.GetAllRoles(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	roleOutput := []*dto.RoleOutput{}
	for _, role := range *roles {
		perms, err := role.Permissions(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}
		output := &dto.RoleOutput{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Active:      role.Active,
			Scopes:      role.Scopes,
			Permissions: perms,
		}
		roleOutput = append(roleOutput, output)
	}

	return roleOutput, nil
}

//GetAllPermissions returns a list of all permissions declared in the system
func (r *RoleUseCaseImpl) GetAllPermissions(
	ctx context.Context,
) ([]*profileutils.Permission, error) {
	ctx, span := tracer.Start(ctx, "GetAllPermissions")
	defer span.End()

	perms, err := profileutils.AllPermissions(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	output := []*profileutils.Permission{}
	for _, perm := range perms {

		p := &profileutils.Permission{
			Scope:       perm.Scope,
			Group:       perm.Group,
			Description: perm.Description,
		}
		output = append(output, p)
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
	allowed, err := r.repo.CheckIfUserHasPermission(ctx, user.UID, profileutils.CanEditRole)
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

	for _, scope := range input.Scopes {
		permission, err := profileutils.GetPermissionByScope(ctx, scope)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}

		if !role.HasPermission(ctx, permission.Scope) {
			role.Scopes = append(role.Scopes, permission.Scope)
		}
	}

	userProfile, err := r.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// save role to database
	updatedRole, err := r.repo.UpdateRoleDetails(ctx, userProfile.ID, *role)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// get permissions
	perms, err := updatedRole.Permissions(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	output := &dto.RoleOutput{
		ID:          updatedRole.ID,
		Name:        updatedRole.Name,
		Description: updatedRole.Description,
		Active:      updatedRole.Active,
		Scopes:      updatedRole.Scopes,
		Permissions: perms,
	}

	return output, nil
}

// GetRole gets a specific role and its permissions
func (r *RoleUseCaseImpl) GetRole(ctx context.Context, ID string) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "GetRole")
	defer span.End()

	// Check logged in user has permissions/role of employee
	user, err := r.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// Check logged in user has the right permissions
	allowed, err := r.repo.CheckIfUserHasPermission(ctx, user.UID, profileutils.CanViewRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	if !allowed {
		return nil, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to view role"),
		)
	}

	role, err := r.repo.GetRoleByID(ctx, ID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	// return all permissions but mark the allowed ones as true
	rolePerms, err := role.Permissions(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}
	output := dto.RoleOutput{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		Active:      role.Active,
		Scopes:      role.Scopes,
		Permissions: rolePerms,
	}
	return &output, nil
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

// GetUserPermissions ...
func (r *RoleUseCaseImpl) GetUserPermissions(
	ctx context.Context,
	UID string,
) ([]profileutils.Permission, error) {
	ctx, span := tracer.Start(ctx, "GetUserPermissions")
	defer span.End()

	roles, err := r.GetUserRoles(ctx, UID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	allPermmissions := []profileutils.Permission{}

	for _, role := range *roles {
		perms, err := role.Permissions(ctx)
		if err != nil {
			utils.RecordSpanError(span, err)
			return nil, err
		}
		allPermmissions = append(allPermmissions, perms...)
	}

	permissions := profileutils.GetUniquePermissions(ctx, allPermmissions)

	return permissions, nil
}
