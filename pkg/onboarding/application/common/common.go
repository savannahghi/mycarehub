package common

// AuthorizedEmails represent emails to check whether they have access to certain dto
var AuthorizedEmails = []string{"apa-dev@healthcloud.co.ke", "apa-prod@healthcloud.co.ke"}

// AuthorizedPhones represent phonenumbers to check whether they have access to certain dto
var AuthorizedPhones = []string{"+254700000000"}

// Icon links for navactions
const (
	// StaticBase is the default path at which static assets are hosted
	StaticBase = "https://assets.healthcloud.co.ke"

	RoleNavActionIcon     = StaticBase + "/actions/roles_navaction.png"
	AgentNavActionIcon    = StaticBase + "/actions/agent_navaction.png"
	EmployeeNavActionIcon = StaticBase + "/actions/employee_navaction.png"
	ConsumerNavActionIcon = StaticBase + "/actions/consumer_navaction.png"
	HelpNavActionIcon     = StaticBase + "/actions/help_navaction.png"
	HomeNavActionIcon     = StaticBase + "/actions/home_navaction.png"
	KYCNavActionIcon      = StaticBase + "/actions/kyc_navaction.png"
	PartnerNavActionIcon  = StaticBase + "/actions/partner_navaction.png"
	PatientNavActionIcon  = StaticBase + "/actions/patient_navaction.png"
	RequestNavActionIcon  = StaticBase + "/actions/request_navaction.png"
)

// On Tap Routes
const (
	HomeRoute                  = "/home"
	PatientRegistrationRoute   = "/addPatient"
	PatientIdentificationRoute = "/patients"
	GetHelpRouteRoute          = "/helpCenter"

	// Has KYC and Covers
	RequestsRoute = "/admin"

	RoleViewRoute     = "/viewCreatedRolesPage"
	RoleCreationRoute = "/createRoleStepOne"

	AgentRegistrationRoute   = "/agentRegistration"
	AgentIdentificationRoute = "/agentIdentification"

	EmployeeRegistrationRoute   = "/employeeRegistration"
	EmployeeIdentificationRoute = "/employeeIdentification"
)

// Navigation actions
const (
	HomeNavActionTitle       = "Home"
	HomeNavActionDescription = "Home Navigation action"

	HelpNavActionTitle       = "Help"
	HelpNavActionDescription = "Help Navigation action"

	RoleNavActionTitle      = "Role Management"
	RoleViewActionTitle     = "View Roles"
	RoleCreationActionTitle = "Create Role"

	PatientNavActionTitle            = "Patients"
	PatientNavActionDescription      = "Patient Navigation action"
	PatientRegistrationActionTitle   = "Register Patient"
	PatientIdentificationActionTitle = "Search Patient"

	RequestsNavActionTitle       = "Requests"
	RequestsNavActionDescription = "Requests Navigation action"

	AgentNavActionTitle            = "Agents"
	AgentNavActionDescription      = "Agent Navigation action"
	AgentRegistrationActionTitle   = "Register Agent"
	AgentIdentificationActionTitle = "View Agents"

	EmployeeNavActionTitle            = "Employees"
	EmployeeNavActionDescription      = "Employee Navigation action"
	EmployeeRegistrationActionTitle   = "Register Employee"
	EmployeeIdentificationActionTitle = "View Employees"

	ConsumerNavActionTitle       = "Consumers"
	ConsumerNavActionDescription = "Consumer Navigation action"

	PartnerNavActionTitle       = "Partners"
	PartnerNavActionDescription = "Partner Navigation action"
)

// PubSub topic names
const (
	// CreateCRMContact is the TopicID for CRM contact creation
	CreateCRMContact = "crm.contact.create"

	// UpdateCRMContact is the topicID for CRM contact updates
	UpdateCRMContact = "crm.contact.update"

	// LinkCoverTopic is the topicID for cover linking topic
	LinkCoverTopic = "covers.link"

	// LinkEDIMemberCoverTopic is the topic ID for cover linking topic of an EDI member who has
	// received a message with the link to download bewell
	LinkEDIMemberCoverTopic = "edi.covers.link"
)
