package enums

import (
	"fmt"
	"io"
	"strconv"
)

// PermissionType is a list of all the permissions.
type PermissionType string

const (
	// PermissionTypeCanInviteUser defines a can invite user permission
	PermissionTypeCanInviteUser PermissionType = "CAN_INVITE_USER"

	// PermissionTypeCanResetUserPassword defines a can reset user password permission
	PermissionTypeCanResetUserPassword PermissionType = "CAN_RESET_USER_PASSWORD"

	// PermissionTypeCanEditUserRole defines a can edit user permission
	PermissionTypeCanEditUserRole PermissionType = "CAN_EDIT_USER_ROLE"

	// PermissionTypeCanEditOwnRole defines a client can edit own permission
	PermissionTypeCanEditOwnRole PermissionType = "CAN_EDIT_OWN_ROLE"

	// PermissionTypeCanTriggerAnalyticsJobs defines a client can can trigger analytics permission
	PermissionTypeCanTriggerAnalyticsJobs PermissionType = "CAN_TRIGGER_ANALYTICS_JOBS"

	// PermissionTypeCanManageOpenMRSIntegration defines a client can manage open MRS integration permission
	PermissionTypeCanManageOpenMRSIntegration PermissionType = "CAN_MANAGE_OPENMRS_INTEGRATION"

	// PermissionTypeCanCreateGroup defines a client can update group permission
	PermissionTypeCanCreateGroup PermissionType = "CAN_CREATE_GROUP"

	// PermissionTypeCanUpdateGroup defines a client can update group permission
	PermissionTypeCanUpdateGroup PermissionType = "CAN_UPDATE_GROUP"

	// PermissionTypeCanModerateGroup defines a client can moderate group permission
	PermissionTypeCanModerateGroup PermissionType = "CAN_MODERATE_GROUP"

	// PermissionTypeCanInviteClientToGroup defines a client can invite client permission
	PermissionTypeCanInviteClientToGroup PermissionType = "CAN_INVITE_CLIENT_TO_GROUP"

	// PermissionTypeCanCreateContentInCMS defines a client can create content in CMS permission
	PermissionTypeCanCreateContentInCMS PermissionType = "CAN_CREATE_CONTENT_IN_CMS"

	// PermissionTypeCanManageContent defines a client can manage content permission
	PermissionTypeCanManageContent PermissionType = "CAN_MANAGE_CONTENT"

	// PermissionTypeCanInviteClient defines a client can invite client permission
	PermissionTypeCanInviteClient PermissionType = "CAN_INVITE_CLIENT"

	// PermissionTypeCanManageClient defines a client can manage client permission
	PermissionTypeCanManageClient PermissionType = "CAN_MANAGE_CLIENT"

	// PermissionTypeCanManageServiceRequest defines a client can manage service request permission
	PermissionTypeCanManageServiceRequest PermissionType = "CAN_MANAGE_SERVICE_REQUEST"

	// PermissionTypeCanViewClientHealthRecords defines a client can view client health records permission
	PermissionTypeCanViewClientHealthRecords PermissionType = "CAN_VIEW_CLIENT_HEALTH_RECORDS"
)

// SystemAdminPermissions is a set of a  valid and known staff admin permissions.
var SystemAdminPermissions = []PermissionType{
	PermissionTypeCanInviteUser,
	PermissionTypeCanResetUserPassword,
	PermissionTypeCanEditUserRole,
	PermissionTypeCanEditOwnRole,
	PermissionTypeCanTriggerAnalyticsJobs,
	PermissionTypeCanManageOpenMRSIntegration,
}

// CommunityManagementPermissions is a set of a  valid and known community management permissions.
var CommunityManagementPermissions = []PermissionType{
	PermissionTypeCanCreateGroup,
	PermissionTypeCanUpdateGroup,
	PermissionTypeCanModerateGroup,
	PermissionTypeCanInviteClientToGroup,
}

//ContentManagementPermissions is a set of a  valid and known content management permissions.
var ContentManagementPermissions = []PermissionType{
	PermissionTypeCanCreateContentInCMS,
	PermissionTypeCanManageContent,
}

//ClientManagementPermissions is a set of a  valid and known client management permissions.
var ClientManagementPermissions = []PermissionType{
	PermissionTypeCanInviteClient,
	PermissionTypeCanManageClient,
	PermissionTypeCanManageServiceRequest,
	PermissionTypeCanViewClientHealthRecords,
}

// IsValid returns true if a permission is valid
func (m PermissionType) IsValid() bool {
	switch m {
	case PermissionTypeCanInviteUser,
		PermissionTypeCanResetUserPassword,
		PermissionTypeCanEditUserRole,
		PermissionTypeCanEditOwnRole,
		PermissionTypeCanTriggerAnalyticsJobs,
		PermissionTypeCanManageOpenMRSIntegration,

		PermissionTypeCanCreateGroup,
		PermissionTypeCanUpdateGroup,
		PermissionTypeCanModerateGroup,
		PermissionTypeCanInviteClientToGroup,

		PermissionTypeCanCreateContentInCMS,
		PermissionTypeCanManageContent,

		PermissionTypeCanInviteClient,
		PermissionTypeCanManageClient,
		PermissionTypeCanManageServiceRequest,
		PermissionTypeCanViewClientHealthRecords:
		return true
	}
	return false
}

// String converts permission type to string
func (m PermissionType) String() string {
	return string(m)
}

// UnmarshalGQL converts the supplied value to a sort type.
func (m *PermissionType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*m = PermissionType(str)
	if !m.IsValid() {
		return fmt.Errorf("%s is not a valid PermissionType", str)
	}
	return nil
}

// MarshalGQL writes the sort type to the supplied
func (m PermissionType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(m.String()))
}
