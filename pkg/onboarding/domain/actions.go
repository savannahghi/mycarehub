package domain

import (
	"github.com/savannahghi/onboarding/pkg/onboarding/application/common"
	"github.com/savannahghi/profileutils"
)

const (
	//HomeGroup groups all actions under the home resource
	HomeGroup NavigationGroup = "home"

	//RoleGroup groups all actions under the role resource
	RoleGroup NavigationGroup = "role"

	//HelpGroup groups all actions under the help resource
	HelpGroup NavigationGroup = "help"

	//AgentGroup groups all actions under the agent resource
	AgentGroup NavigationGroup = "agents"

	//KYCGroup groups all actions under the kyc resource
	KYCGroup NavigationGroup = "kyc"

	//EmployeeGroup groups all actions under the employees resource
	EmployeeGroup NavigationGroup = "employees"

	//CoversGroup groups all actions under the covers resource
	CoversGroup NavigationGroup = "covers"

	//PatientGroup groups all actions under the patient resource
	PatientGroup NavigationGroup = "patient"

	//PartnerGroup groups all actions under the partner resource
	PartnerGroup NavigationGroup = "partner"

	//RolesGroup groups all actions under the role resource
	RolesGroup NavigationGroup = "role"

	//ConsumerGroup groups all actions under the consumer resource
	ConsumerGroup NavigationGroup = "consumer"
)

// Determines the sequence number of a navigation action
// Order of the constants matters!!
const (
	HomeNavActionSequence = iota + 1

	RoleNavActionSequence
	RoleCreationNavActionSequence
	RoleViewingNavActionSequence

	RequestsNavActionSequence

	PartnerNavactionSequence

	ConsumerNavactionSequence

	EmployeeNavActionSequence
	EmployeeSearchNavActionSequence
	EmployeeRegistrationActionSequence

	AgentNavActionSequence
	AgentSearchNavActionSequence
	AgentRegistrationActionSequence

	PatientNavActionSequence
	PatientSearchNavActionSequence
	PatientRegistrationNavActionSequence

	HelpNavActionSequence
)

// the structure and definition of all navigation actions
var (
	// HomeNavAction is the primary home button
	HomeNavAction = NavigationAction{
		Group:              HomeGroup,
		Title:              common.HomeNavActionTitle,
		OnTapRoute:         common.HomeRoute,
		Icon:               common.HomeNavActionIcon,
		RequiredPermission: nil,
		SequenceNumber:     HomeNavActionSequence,
	}

	// HelpNavAction navigation action to help and FAQs page
	HelpNavAction = NavigationAction{
		Group:              HelpGroup,
		Title:              common.HelpNavActionTitle,
		OnTapRoute:         common.GetHelpRouteRoute,
		Icon:               common.HelpNavActionIcon,
		RequiredPermission: nil,
		SequenceNumber:     HelpNavActionSequence,
	}
)

var (

	// KYCNavActions is the navigation acction to KYC processing
	KYCNavActions = NavigationAction{
		Group:              KYCGroup,
		Title:              common.RequestsNavActionTitle,
		OnTapRoute:         common.RequestsRoute,
		Icon:               common.RequestNavActionIcon,
		RequiredPermission: &profileutils.CanProcessKYC,
		SequenceNumber:     RequestsNavActionSequence,
	}
)

var (
	//PartnerNavActions is the navigation actions to partner management
	PartnerNavActions = NavigationAction{
		Group: PartnerGroup,
		Title: common.PartnerNavActionTitle,
		// Not provided yet
		OnTapRoute:         "",
		Icon:               common.PartnerNavActionIcon,
		RequiredPermission: &profileutils.CanViewPartner,
		SequenceNumber:     PartnerNavactionSequence,
	}
)

var (
	//ConsumerNavActions is the navigation actions to consumer management
	ConsumerNavActions = NavigationAction{
		Group: ConsumerGroup,
		Title: common.ConsumerNavActionTitle,
		// Not provided yet
		OnTapRoute:         "",
		Icon:               common.ConsumerNavActionIcon,
		RequiredPermission: &profileutils.CanViewConsumers,
		SequenceNumber:     ConsumerNavactionSequence,
	}
)

var (
	//RoleNavActions this is the parent navigation action for role resource
	// it has nested navigation actions below
	RoleNavActions = NavigationAction{
		Group:              RoleGroup,
		Title:              common.RoleNavActionTitle,
		Icon:               common.RoleNavActionIcon,
		RequiredPermission: &profileutils.CanViewRole,
		SequenceNumber:     RoleNavActionSequence,
	}

	//RoleCreationNavAction a child of the RoleNavActions
	RoleCreationNavAction = NavigationAction{
		Group:              RoleGroup,
		Title:              common.RoleCreationActionTitle,
		OnTapRoute:         common.RoleCreationRoute,
		RequiredPermission: &profileutils.CanCreateRole,
		HasParent:          true,
		SequenceNumber:     RoleCreationNavActionSequence,
	}

	//RoleViewNavAction a child of the RoleNavActions
	RoleViewNavAction = NavigationAction{
		Group:              RoleGroup,
		Title:              common.RoleViewActionTitle,
		OnTapRoute:         common.RoleViewRoute,
		RequiredPermission: &profileutils.CanViewRole,
		HasParent:          true,
		SequenceNumber:     RoleViewingNavActionSequence,
	}
)

var (
	//AgentNavActions this is the parent navigation action for agent resource
	// it has nested navigation actions below
	AgentNavActions = NavigationAction{
		Group:              AgentGroup,
		Title:              common.AgentNavActionTitle,
		Icon:               common.AgentNavActionIcon,
		RequiredPermission: &profileutils.CanViewAgent,
		SequenceNumber:     AgentNavActionSequence,
	}

	//AgentRegistrationNavAction a child of the AgentNavActions
	AgentRegistrationNavAction = NavigationAction{
		Group:              AgentGroup,
		Title:              common.AgentRegistrationActionTitle,
		OnTapRoute:         common.AgentRegistrationRoute,
		RequiredPermission: &profileutils.CanRegisterAgent,
		HasParent:          true,
		SequenceNumber:     AgentRegistrationActionSequence,
	}

	//AgentidentificationNavAction a child of the AgentNavActions
	AgentidentificationNavAction = NavigationAction{
		Group:              AgentGroup,
		Title:              common.AgentIdentificationActionTitle,
		OnTapRoute:         common.AgentIdentificationRoute,
		RequiredPermission: &profileutils.CanIdentifyAgent,
		HasParent:          true,
		SequenceNumber:     AgentSearchNavActionSequence,
	}
)

var (
	//EmployeeNavActions this is the parent navigation action for agent resource
	// it has nested navigation actions below
	EmployeeNavActions = NavigationAction{
		Group:              EmployeeGroup,
		Title:              common.EmployeeNavActionTitle,
		Icon:               common.EmployeeNavActionIcon,
		RequiredPermission: &profileutils.CanViewEmployee,
		SequenceNumber:     EmployeeNavActionSequence,
	}

	//EmployeeRegistrationNavAction a child of the EmployeeNavActions
	EmployeeRegistrationNavAction = NavigationAction{
		Group:              EmployeeGroup,
		Title:              common.EmployeeRegistrationActionTitle,
		OnTapRoute:         common.EmployeeRegistrationRoute,
		RequiredPermission: &profileutils.CanCreateEmployee,
		HasParent:          true,
		SequenceNumber:     EmployeeRegistrationActionSequence,
	}

	//EmployeeidentificationNavAction a child of the EmployeeNavActions
	EmployeeidentificationNavAction = NavigationAction{
		Group:              EmployeeGroup,
		Title:              common.EmployeeIdentificationActionTitle,
		OnTapRoute:         common.EmployeeIdentificationRoute,
		RequiredPermission: &profileutils.CanViewEmployee,
		HasParent:          true,
		SequenceNumber:     EmployeeSearchNavActionSequence,
	}
)

var (
	//PatientNavActions this is the parent navigation action for patient resource
	// it has nested navigation actions below
	PatientNavActions = NavigationAction{
		Group:              PatientGroup,
		Title:              common.PatientNavActionTitle,
		Icon:               common.PatientNavActionIcon,
		RequiredPermission: &profileutils.CanViewPatient,
		SequenceNumber:     PatientNavActionSequence,
	}

	//PatientRegistrationNavAction a child of the PatientNavActions
	PatientRegistrationNavAction = NavigationAction{
		Group:              PatientGroup,
		Title:              common.PatientRegistrationActionTitle,
		OnTapRoute:         common.PatientRegistrationRoute,
		RequiredPermission: &profileutils.CanCreatePatient,
		HasParent:          true,
		SequenceNumber:     PatientRegistrationNavActionSequence,
	}

	//PatientIdentificationNavAction a child of the PatientNavActions
	PatientIdentificationNavAction = NavigationAction{
		Group:              PatientGroup,
		Title:              common.PatientIdentificationActionTitle,
		OnTapRoute:         common.PatientIdentificationRoute,
		RequiredPermission: &profileutils.CanIdentifyPatient,
		HasParent:          true,
		SequenceNumber:     PatientSearchNavActionSequence,
	}
)

// AllNavigationActions is a grouping of all navigation actions
var AllNavigationActions = []NavigationAction{
	HomeNavAction, HelpNavAction,

	KYCNavActions, PartnerNavActions, ConsumerNavActions,

	AgentNavActions, AgentRegistrationNavAction, AgentidentificationNavAction,

	EmployeeNavActions, EmployeeRegistrationNavAction, EmployeeidentificationNavAction,

	PatientNavActions, PatientRegistrationNavAction, PatientIdentificationNavAction,

	RoleNavActions, RoleCreationNavAction, RoleViewNavAction,
}
