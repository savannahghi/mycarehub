package authorization

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/authorization"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IAuthorization contains methods associated with authorization
type IAuthorization interface {
	CheckPermissions(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error)
	CheckAuthorization(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error)
	IsAuthorized(ctx context.Context, permission dto.PermissionInput) (bool, error)
	AddPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error)
	RemovePolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error)
	AddGroupingPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error)
	RemoveGroupingPolicy(ctx context.Context, subject string, input dto.PermissionInput) (bool, error)
}

// UsecaseAuthorization groups al the interfaces for the Authorization usecase
type UsecaseAuthorization interface {
	IAuthorization
}

// UsecaseAuthorizationImpl represents the Authorization implementation
type UsecaseAuthorizationImpl struct {
	Authorization authorization.Authorization
	ExternalExt   extension.ExternalMethodsExtension
	Query         infrastructure.Query
}

// NewUsecaseAuthorization is the controller function for the Authorization usecase
func NewUsecaseAuthorization(
	authorization authorization.Authorization,
	extension extension.ExternalMethodsExtension,
	query infrastructure.Query,

) *UsecaseAuthorizationImpl {
	return &UsecaseAuthorizationImpl{
		Authorization: authorization,
		ExternalExt:   extension,
		Query:         query,
	}
}

// CheckPermissions is used to check whether the permissions of a subject are set
func (a UsecaseAuthorizationImpl) CheckPermissions(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {

	ok, err := a.Authorization.Enforce(ctx, subject, permission)
	if err != nil {
		return false, fmt.Errorf("unable to check permissions %w", err)
	}
	if ok {
		return true, nil
	}
	return false, nil
}

// CheckAuthorization is used to check the user permissions
func (a UsecaseAuthorizationImpl) CheckAuthorization(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	isAuthorized, err := a.CheckPermissions(ctx, subject, permission)
	if err != nil {
		return false, fmt.Errorf("internal server error: can't authorize user: %w", err)
	}

	if !isAuthorized {
		return false, nil
	}

	return true, nil
}

// IsAuthorized checks if the subject identified by their email has permission to access the
// specified resource
// currently only known internal anonymous users and external API Integrations emails are checked, internal and default logged in users
// have access by default.
// for subjects identified by their phone number normalize the phone and omit the first (+) character
func (a UsecaseAuthorizationImpl) IsAuthorized(ctx context.Context, permission dto.PermissionInput) (bool, error) {
	loggedInUserID, err := a.ExternalExt.GetLoggedInUserUID(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.GetLoggedInUserUIDErr(err)
	}
	user, err := a.Query.GetUserProfileByUserID(ctx, loggedInUserID)
	if err != nil {
		return false, exceptions.ProfileNotFoundErr(err)
	}

	subject := loggedInUserID
	permission.OrganizationID = user.OrganizationID
	// TODO: change to the actual program ID
	permission.ProgramID = user.OrganizationID

	return a.CheckAuthorization(ctx, subject, permission)

}

func (a UsecaseAuthorizationImpl) AddPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return a.Authorization.AddPolicy(ctx, subject, permission)
}
func (a UsecaseAuthorizationImpl) RemovePolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return a.Authorization.RemovePolicy(ctx, subject, permission)
}
func (a UsecaseAuthorizationImpl) AddGroupingPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return a.Authorization.AddGroupingPolicy(ctx, subject, permission)
}

func (a UsecaseAuthorizationImpl) RemoveGroupingPolicy(ctx context.Context, subject string, permission dto.PermissionInput) (bool, error) {
	return a.Authorization.RemoveGroupingPolicy(ctx, subject, permission)
}
