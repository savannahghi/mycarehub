package usecases

import (
	"context"
	"fmt"

	"github.com/savannahghi/firebasetools"
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
	DeleteRole(ctx context.Context, roleID string) (bool, error)
	GetAllRoles(ctx context.Context, filter *firebasetools.FilterInput) ([]*dto.RoleOutput, error)
	GetAllPermissions(ctx context.Context) ([]*profileutils.Permission, error)

	// AddPermissionsToRole adds new scopes to a role
	AddPermissionsToRole(
		ctx context.Context,
		input dto.RolePermissionInput,
	) (*dto.RoleOutput, error)

	// RevokeRolePermission removes the specified scopes from a prole
	RevokeRolePermission(
		ctx context.Context,
		input dto.RolePermissionInput,
	) (*dto.RoleOutput, error)

	// UpdateRolePermissions replaces the scopes in a role with new updated scopes
	UpdateRolePermissions(ctx context.Context, input dto.RolePermissionInput) (*dto.RoleOutput, error)

	GetUserRoles(ctx context.Context, UID string) (*[]profileutils.Role, error)

	GetRole(ctx context.Context, ID string) (*dto.RoleOutput, error)

	GetUserPermissions(ctx context.Context, UID string) ([]profileutils.Permission, error)

	// AssignRole assigns a role to a user
	AssignRole(ctx context.Context, userID string, roleID string) (bool, error)

	// RevokeRole removes a role from a user
	RevokeRole(ctx context.Context, userID string, roleID string) (bool, error)

	// ActivateRole marks a role as active
	ActivateRole(ctx context.Context, roleID string) (*dto.RoleOutput, error)

	// DeactivateRole marks a role as inactive and cannot be used
	DeactivateRole(ctx context.Context, roleID string) (*dto.RoleOutput, error)
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
func (r *RoleUseCaseImpl) GetAllRoles(
	ctx context.Context,
	filter *firebasetools.FilterInput,
) ([]*dto.RoleOutput, error) {
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

	roles, err := r.repo.GetAllRoles(ctx, filter)
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

//DeleteRole removes a role from the database permanently
func (r *RoleUseCaseImpl) DeleteRole(ctx context.Context, roleID string) (bool, error) {
	ctx, span := tracer.Start(ctx, "DeleteRole")
	defer span.End()

	user, err := r.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	// Check logged in user has the right permissions
	allowed, err := r.repo.CheckIfUserHasPermission(ctx, user.UID, profileutils.CanEditRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	if !allowed {
		return false, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to delete roles"),
		)
	}

	success, err := r.repo.DeleteRole(ctx, roleID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	return success, nil
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

// RevokeRolePermission removes a permission from a role
func (r *RoleUseCaseImpl) RevokeRolePermission(
	ctx context.Context,
	input dto.RolePermissionInput,
) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "RevokeRolePermission")
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

	newScopes := []string{}

	for _, roleScope := range role.Scopes {
		for _, scope := range input.Scopes {
			if roleScope != scope {
				newScopes = append(newScopes, roleScope)
			}
		}
	}

	role.Scopes = newScopes

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

// AssignRole assigns a user a particular role
func (r *RoleUseCaseImpl) AssignRole(
	ctx context.Context,
	userID string,
	roleID string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "AssignRole")
	defer span.End()

	role, err := r.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	profile, err := r.repo.GetUserProfileByID(ctx, userID, false)
	if err != nil {
		return false, err
	}

	for _, r := range profile.Roles {
		// check if role exists first
		if r == role.ID {
			err := fmt.Errorf("role already exists: %v", role.Name)
			return false, err
		}
	}

	updated := append(profile.Roles, roleID)

	err = r.repo.UpdateUserRoleIDs(ctx, profile.ID, updated)
	if err != nil {
		return false, err
	}

	return true, nil
}

// RevokeRole removes a role from the user
func (r *RoleUseCaseImpl) RevokeRole(
	ctx context.Context,
	userID string,
	roleID string,
) (bool, error) {
	ctx, span := tracer.Start(ctx, "RevokeRole")
	defer span.End()

	// Check logged in user has permissions/role of employee
	user, err := r.baseExt.GetLoggedInUser(ctx)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	// Check logged in user has the right permissions
	allowed, err := r.repo.CheckIfUserHasPermission(ctx, user.UID, profileutils.CanAssignRole)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}
	if !allowed {
		return false, exceptions.RoleNotValid(
			fmt.Errorf("error: logged in user does not have permissions to view roles"),
		)
	}

	role, err := r.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return false, err
	}

	profile, err := r.repo.GetUserProfileByID(ctx, userID, false)
	if err != nil {
		return false, err
	}

	var exist bool
	for _, r := range profile.Roles {
		// check if role exists first
		if r == role.ID {
			exist = true
			break
		}
	}

	if !exist {
		err := fmt.Errorf("user doesn't have role: %v", role.Name)
		utils.RecordSpanError(span, err)
		return false, err
	}

	// roles copy
	updated := []string{}
	for _, r := range profile.Roles {
		if r != role.ID {
			updated = append(updated, r)
		}
	}

	err = r.repo.UpdateUserRoleIDs(ctx, profile.ID, updated)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ActivateRole marks a deactivated role as active and usable
func (r *RoleUseCaseImpl) ActivateRole(ctx context.Context, roleID string) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "ActivateRole")
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

	role, err := r.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	userProfile, err := r.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	role.Active = true

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

// DeactivateRole marks a role as inactive and cannot be used
func (r *RoleUseCaseImpl) DeactivateRole(ctx context.Context, roleID string) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "DeactivateRole")
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

	role, err := r.repo.GetRoleByID(ctx, roleID)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	userProfile, err := r.repo.GetUserProfileByUID(ctx, user.UID, false)
	if err != nil {
		utils.RecordSpanError(span, err)
		return nil, err
	}

	role.Active = false

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

// UpdateRolePermissions replaces the scopes in a role with new updated scopes
func (r *RoleUseCaseImpl) UpdateRolePermissions(ctx context.Context, input dto.RolePermissionInput) (*dto.RoleOutput, error) {
	ctx, span := tracer.Start(ctx, "UpdateRolePermissions")
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

	// new scopes
	role.Scopes = input.Scopes

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
