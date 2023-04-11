package authorization

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// PermissionCategory groups permissions into categories
type PermissionCategory string

const (
	PermissionCategoryAppointment      PermissionCategory = "Appointment"
	PermissionCategoryAuthorization    PermissionCategory = "Authorization"
	PermissionCategoryContent          PermissionCategory = "Content"
	PermissionCategoryFacility         PermissionCategory = "Facility"
	PermissionCategoryFeedback         PermissionCategory = "Feedback"
	PermissionCategoryHealthDiary      PermissionCategory = "HealthDiary"
	PermissionCategoryNotification     PermissionCategory = "Notification"
	PermissionCategoryOrganisation     PermissionCategory = "Organisation"
	PermissionCategoryOTP              PermissionCategory = "OTP"
	PermissionCategoryProgram          PermissionCategory = "Program"
	PermissionCategoryScreeningTool    PermissionCategory = "ScreeningTool"
	PermissionCategorySecurityQuestion PermissionCategory = "SecurityQuestion"
	PermissionCategoryServiceRequest   PermissionCategory = "ServiceRequest"
	PermissionCategorySurvey           PermissionCategory = "Survey"
	PermissionCategoryUser             PermissionCategory = "User"
)

// IsValid checks if a string of type PermissionCategory is of the valid type
func (p PermissionCategory) IsValid() bool {
	switch p {
	case PermissionCategoryAppointment,
		PermissionCategoryAuthorization,
		PermissionCategoryContent,
		PermissionCategoryFacility,
		PermissionCategoryFeedback,
		PermissionCategoryHealthDiary,
		PermissionCategoryNotification,
		PermissionCategoryOrganisation,
		PermissionCategoryOTP,
		PermissionCategoryProgram,
		PermissionCategoryScreeningTool,
		PermissionCategorySecurityQuestion,
		PermissionCategoryServiceRequest,
		PermissionCategorySurvey,
		PermissionCategoryUser:
		return true
	}
	return false
}

// String converts PermissionCategory type to type string
func (p PermissionCategory) String() string {
	return string(p)
}

// UnmarshalGQL converts the supplied value to a PermissionCategory type.
func (p *PermissionCategory) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("value must be of type string")
	}

	*p = PermissionCategory(str)
	if !p.IsValid() {
		return fmt.Errorf("%s is not a valid PermissionCategory", str)
	}
	return nil
}

// MarshalGQL writes the PermissionCategory type to the supplied
func (p PermissionCategory) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(p.String()))
}

// Appointment Permissions
var (
	canReadClientAppointment = domain.AuthorityPermission{
		Name:        "Read client appointment",
		Description: "Can read client appointment",
		Category:    PermissionCategoryAppointment.String(),
		Scope:       "client.appointment.read",
	}
	canUpdateClientAppointment = domain.AuthorityPermission{
		Name:        "Update client appointment",
		Description: "Can update client appointment",
		Category:    PermissionCategoryAppointment.String(),
		Scope:       "client.appointment.update",
	}
)

// Authorization Permissions
var (
	canReadSystemRole = domain.AuthorityPermission{
		Name:        "Read system roles",
		Description: "Can read system roles",
		Category:    PermissionCategoryAuthorization.String(),
		Scope:       "role.read",
	}
)

// Content Permissions
var (
	canReadContent = domain.AuthorityPermission{
		Name:        "Read content",
		Description: "Can read content",
		Category:    PermissionCategoryContent.String(),
		Scope:       "content.read",
	}
)

// Facility Permissions
var (
	canDeleteFacility = domain.AuthorityPermission{
		Name:        "Delete facility",
		Description: "Can delete facility",
		Category:    PermissionCategoryFacility.String(),
		Scope:       "facility.delete",
	}
	canUpdateFacility = domain.AuthorityPermission{
		Name:        "Update Facility",
		Description: "Can update facility",
		Category:    PermissionCategoryFacility.String(),
		Scope:       "facility.update",
	}
	canReadFacility = domain.AuthorityPermission{
		Name:        "Read Facility",
		Description: "Can read facility",
		Category:    PermissionCategoryFacility.String(),
		Scope:       "facility.read",
	}
	canCreateProgramFacility = domain.AuthorityPermission{
		Name:        "Create program facility",
		Description: "Can create a facility in a program",
		Category:    PermissionCategoryFacility.String(),
		Scope:       "program.facility.create",
	}
)

// Feedback Permissions
var (
	canCreateFeedback = domain.AuthorityPermission{
		Name:        "Create feedback",
		Description: "Can create feedback",
		Category:    PermissionCategoryFeedback.String(),
		Scope:       "feedback.create",
	}
)

// HealthDiary Permissions
var (
	canCreateHealthDiary = domain.AuthorityPermission{
		Name:        "Create health diary",
		Description: "Can create health diary",
		Category:    PermissionCategoryHealthDiary.String(),
		Scope:       "healthdiary.create",
	}
	canReadHealthDiary = domain.AuthorityPermission{
		Name:        "Read health diary",
		Description: "Can read health diary",
		Category:    PermissionCategoryHealthDiary.String(),
		Scope:       "healthdiary.read",
	}
	canReadClientHealthDiary = domain.AuthorityPermission{
		Name:        "Read Client health diary",
		Description: "Can read client health diary",
		Category:    PermissionCategoryHealthDiary.String(),
		Scope:       "client.healthdiary.read",
	}
)

// Notification Permissions
var (
	canReadNotification = domain.AuthorityPermission{
		Name:        "Read notification",
		Description: "Can read notification",
		Category:    PermissionCategoryNotification.String(),
		Scope:       "notification.read",
	}
)

// Organisation Permissions
var (
	canReadOrganisation = domain.AuthorityPermission{
		Name:        "Read organisation",
		Description: "Can read organisation",
		Category:    PermissionCategoryOrganisation.String(),
		Scope:       "organisation.read",
	}
	canCreateOrganisation = domain.AuthorityPermission{
		Name:        "Create organisation",
		Description: "Can create organisation",
		Category:    PermissionCategoryOrganisation.String(),
		Scope:       "organisation.create",
	}
	canDeleteOrganisation = domain.AuthorityPermission{
		Name:        "Delete organisation",
		Description: "Can delete organisation",
		Category:    PermissionCategoryOrganisation.String(),
		Scope:       "organisation.delete",
	}
)

// OTP Permissions
var (
	canCreateOTP = domain.AuthorityPermission{
		Name:        "Create OTP",
		Description: "Can create OTP",
		Category:    PermissionCategoryOTP.String(),
		Scope:       "otp.create",
	}
)

// Program Permissions
var (
	canReadProgram = domain.AuthorityPermission{
		Name:        "Read Program",
		Description: "Can read program",
		Category:    PermissionCategoryProgram.String(),
		Scope:       "program.read",
	}
	canCreateProgram = domain.AuthorityPermission{
		Name:        "Create Program",
		Description: "Can create program",
		Category:    PermissionCategoryProgram.String(),
		Scope:       "program.create",
	}
	canUpdateProgram = domain.AuthorityPermission{
		Name:        "Update Program",
		Description: "Can update program",
		Category:    PermissionCategoryProgram.String(),
		Scope:       "program.update",
	}
)

// ScreeningTool Permissions
var (
	canReadScreeningTool = domain.AuthorityPermission{
		Name:        "Read screening tool",
		Description: "Can read screening tool",
		Category:    PermissionCategoryScreeningTool.String(),
		Scope:       "screeningtool.read",
	}
	canCreateScreeningTool = domain.AuthorityPermission{
		Name:        "Create screening tool",
		Description: "Can create screening tool",
		Category:    PermissionCategoryScreeningTool.String(),
		Scope:       "screeningtool.create",
	}
	canReadScreeningToolResponse = domain.AuthorityPermission{
		Name:        "Read screening tool response",
		Description: "Can read screening tool response",
		Category:    PermissionCategoryScreeningTool.String(),
		Scope:       "screeningtool.response.read",
	}
	canCreateScreeningToolResponse = domain.AuthorityPermission{
		Name:        "Create screening tool response",
		Description: "Can create screening tool response",
		Category:    PermissionCategoryScreeningTool.String(),
		Scope:       "screeningtool.response.create",
	}
	canReadScreeningToolRespondent = domain.AuthorityPermission{
		Name:        "Read screening tool respondent",
		Description: "Can read screening tool respondent",
		Category:    PermissionCategoryScreeningTool.String(),
		Scope:       "screeningtool.respondent.read",
	}
)

// SecurityQuestion Permissions
var (
	canReadSecurityQuestion = domain.AuthorityPermission{
		Name:        "Read security question",
		Description: "Can read security question",
		Category:    PermissionCategorySecurityQuestion.String(),
		Scope:       "securityquestion.read",
	}
	canCreateSecurityQuestion = domain.AuthorityPermission{
		Name:        "Create security question",
		Description: "Can create security question",
		Category:    PermissionCategorySecurityQuestion.String(),
		Scope:       "securityquestion.create",
	}
)

// ServiceRequest Permissions
var (
	canReadServiceRequest = domain.AuthorityPermission{
		Name:        "Read service request",
		Description: "Can read service request",
		Category:    PermissionCategoryServiceRequest.String(),
		Scope:       "servicerequest.read",
	}
	canCreateServiceRequest = domain.AuthorityPermission{
		Name:        "Create service request",
		Description: "Can create service request",
		Category:    PermissionCategoryServiceRequest.String(),
		Scope:       "servicerequest.create",
	}
	canUpdateServiceRequest = domain.AuthorityPermission{
		Name:        "Update service request",
		Description: "Can update service request",
		Category:    PermissionCategoryServiceRequest.String(),
		Scope:       "servicerequest.update",
	}
	canUpdateClientServiceRequest = domain.AuthorityPermission{
		Name:        "Update client service request",
		Description: "Can update client service request",
		Category:    PermissionCategoryServiceRequest.String(),
		Scope:       "client.servicerequest.update",
	}
	canUpdateStaffServiceRequest = domain.AuthorityPermission{
		Name:        "Update staff service request",
		Description: "Can update staff service request",
		Category:    PermissionCategoryServiceRequest.String(),
		Scope:       "staff.servicerequest.update",
	}
)

// Survey Permissions
var (
	canReadSurvey = domain.AuthorityPermission{
		Name:        "Read survey",
		Description: "Can read survey",
		Category:    PermissionCategorySurvey.String(),
		Scope:       "survey.read",
	}
	canReadSurveyRespondent = domain.AuthorityPermission{
		Name:        "Read survey respondent",
		Description: "Can read survey respondent",
		Category:    PermissionCategorySurvey.String(),
		Scope:       "survey.respondent.read",
	}
	canReadClientWithServiceRequest = domain.AuthorityPermission{
		Name:        "Read client with survey service request",
		Description: "Can read client with service request from the survey",
		Category:    PermissionCategorySurvey.String(),
		Scope:       "client.servicerequest.survey.read",
	}
	canReadSurveyResponse = domain.AuthorityPermission{
		Name:        "Read survey response",
		Description: "Can read survey response",
		Category:    PermissionCategorySurvey.String(),
		Scope:       "survey.response.read",
	}
	canCreateSurveyLink = domain.AuthorityPermission{
		Name:        "Create survey link",
		Description: "Can create survey link",
		Category:    PermissionCategorySurvey.String(),
		Scope:       "survey.link.create",
	}
)

// User Permissions
var (
	canReadTerms = domain.AuthorityPermission{
		Name:        "Read terms",
		Description: "Can read terms",
		Category:    PermissionCategoryUser.String(),
		Scope:       "terms.read",
	}
	canReadPin = domain.AuthorityPermission{
		Name:        "Read PIN",
		Description: "Can read PIN",
		Category:    PermissionCategoryUser.String(),
		Scope:       "pin.read",
	}
	canReadClient = domain.AuthorityPermission{
		Name:        "Read client",
		Description: "Can read client",
		Category:    PermissionCategoryUser.String(),
		Scope:       "client.read",
	}
	canReadStaff = domain.AuthorityPermission{
		Name:        "Read staff",
		Description: "Can read staff",
		Category:    PermissionCategoryUser.String(),
		Scope:       "staff.read",
	}
	canReadCaregiver = domain.AuthorityPermission{
		Name:        "Read caregiver",
		Description: "Can read caregiver",
		Category:    PermissionCategoryUser.String(),
		Scope:       "caregiver.read",
	}
	canReadClientOfCaregiver = domain.AuthorityPermission{
		Name:        "Read clients of a caregiver",
		Description: "Can read clients of a caregiver",
		Category:    PermissionCategoryUser.String(),
		Scope:       "client.caregiver.read",
	}
	canReadCaregiverOfClient = domain.AuthorityPermission{
		Name:        "Read caregivers of a client",
		Description: "Can read caregivers of a client",
		Category:    PermissionCategoryUser.String(),
		Scope:       "caregiver.client.read",
	}
	canReadStaffFacility = domain.AuthorityPermission{
		Name:        "Read staff's facility",
		Description: "Can read staff's facility",
		Category:    PermissionCategoryUser.String(),
		Scope:       "staff.facility.read",
	}
	canReadClientFacility = domain.AuthorityPermission{
		Name:        "Read client's facility",
		Description: "Can read client's facility",
		Category:    PermissionCategoryUser.String(),
		Scope:       "client.facility.read",
	}
	canReadClientIdentifier = domain.AuthorityPermission{
		Name:        "Read client's identifier",
		Description: "Can read client's identifier",
		Category:    PermissionCategoryUser.String(),
		Scope:       "client.identifier.read",
	}
	canUpdateUser = domain.AuthorityPermission{
		Name:        "Create user",
		Description: "Can create user",
		Category:    PermissionCategoryUser.String(),
		Scope:       "user.create",
	}
	canCreateClient = domain.AuthorityPermission{
		Name:        "Create client",
		Description: "Can create client",
		Category:    PermissionCategoryUser.String(),
		Scope:       "client.create",
	}
	canCreateStaff = domain.AuthorityPermission{
		Name:        "Create staff",
		Description: "Can create staff",
		Category:    PermissionCategoryUser.String(),
		Scope:       "staff.create",
	}
	canCreateCaregiver = domain.AuthorityPermission{
		Name:        "Create caregiver",
		Description: "Can create caregiver",
		Category:    PermissionCategoryUser.String(),
		Scope:       "caregiver.create",
	}
	canDeleteUser = domain.AuthorityPermission{
		Name:        "Delete user",
		Description: "Can delete user",
		Category:    PermissionCategoryUser.String(),
		Scope:       "user.delete",
	}
	//canCreateUserInvite = domain.AuthorityPermission{
	//	Name:        "Create user invite",
	//	Description: "Can create user invite",
	//	Category:    PermissionCategoryUser.String(),
	//	Scope:       "user.invite.create",
	//}
)

// AllPermissions returns all the defined permissions
func AllPermissions(ctx context.Context) []domain.AuthorityPermission {
	return []domain.AuthorityPermission{
		// Appointment Permissions
		canReadClientAppointment,
		canUpdateClientAppointment,

		// Authorization Permissions
		canReadSystemRole,

		// Content Permissions
		canReadContent,

		// Facility Permissions
		canDeleteFacility,
		canUpdateFacility,
		canReadFacility,
		canCreateProgramFacility,

		// Feedback Permissions
		canCreateFeedback,

		// HealthDiary Permissions
		canCreateHealthDiary,
		canReadHealthDiary,
		canReadClientHealthDiary,

		// Notification Permissions
		canReadNotification,

		// Organisation Permissions
		canReadOrganisation,
		canCreateOrganisation,
		canDeleteOrganisation,

		// OTP Permissions
		canCreateOTP,

		// Program Permissions
		canReadProgram,
		canCreateProgram,
		canUpdateProgram,

		// ScreeningTool Permissions
		canReadScreeningTool,
		canCreateScreeningTool,
		canReadScreeningToolResponse,
		canCreateScreeningToolResponse,
		canReadScreeningToolRespondent,

		// SecurityQuestion Permissions
		canReadSecurityQuestion,
		canCreateSecurityQuestion,

		// ServiceRequest Permissions
		canReadServiceRequest,
		canCreateServiceRequest,
		canUpdateServiceRequest,
		canUpdateClientServiceRequest,
		canUpdateStaffServiceRequest,

		// Survey Permissions
		canReadSurvey,
		canReadSurveyRespondent,
		canReadClientWithServiceRequest,
		canReadSurveyResponse,
		canCreateSurveyLink,

		// User Permissions
		canReadTerms,
		canReadPin,
		canReadClient,
		canReadStaff,
		canReadCaregiver,
		canReadClientOfCaregiver,
		canReadCaregiverOfClient,
		canReadStaffFacility,
		canReadClientFacility,
		canReadClientIdentifier,
		canUpdateUser,
		canCreateClient,
		canCreateStaff,
		canCreateCaregiver,
		canDeleteUser,
	}
}

// DefaultClientPermissions return default client permissions
func DefaultClientPermissions(ctx context.Context) []domain.AuthorityPermission {
	return []domain.AuthorityPermission{
		// Appointment Permissions
		canReadClientAppointment,
		canUpdateClientAppointment,

		// Content Permissions
		canReadContent,

		// Facility Permissions
		canReadFacility,

		// Feedback Permissions
		canCreateFeedback,

		// HealthDiary Permissions
		canCreateHealthDiary,
		canReadHealthDiary,

		// Notification Permissions
		canReadNotification,

		// Organisation Permissions
		canReadOrganisation,

		// OTP Permissions
		canCreateOTP,

		// Program Permissions
		canReadProgram,

		// ScreeningTool Permissions
		canReadScreeningTool,
		canCreateScreeningToolResponse,

		// SecurityQuestion Permissions
		canReadSecurityQuestion,
		canCreateSecurityQuestion,

		// ServiceRequest Permissions
		canCreateServiceRequest,

		// Survey Permissions
		canReadSurvey,

		// User Permissions
		canReadTerms,
		canReadPin,
		canReadClient,
		canReadCaregiverOfClient,
		canReadClientFacility,
		canUpdateUser,
	}
}

// DefaultCaregiverPermissions return default caregiver permissions
func DefaultCaregiverPermissions(ctx context.Context) []domain.AuthorityPermission {
	return []domain.AuthorityPermission{
		// Appointment Permissions
		canReadClientAppointment,
		canUpdateClientAppointment,

		// Facility Permissions
		canReadFacility,

		// Feedback Permissions
		canCreateFeedback,

		// HealthDiary Permissions
		canReadHealthDiary,

		// Notification Permissions
		canReadNotification,

		// Organisation Permissions
		canReadOrganisation,

		// OTP Permissions
		canCreateOTP,

		// Program Permissions
		canReadProgram,

		// SecurityQuestion Permissions
		canReadSecurityQuestion,
		canCreateSecurityQuestion,

		// ServiceRequest Permissions
		canCreateServiceRequest,

		// User Permissions
		canReadTerms,
		canReadPin,
		canReadClient,
		canReadClientOfCaregiver,
		canReadClientFacility,
		canUpdateUser,
	}
}
