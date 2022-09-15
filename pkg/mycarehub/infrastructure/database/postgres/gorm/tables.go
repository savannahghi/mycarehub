package gorm

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/serverutils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrganizationID assign a default organisation to a type
var OrganizationID = serverutils.MustGetEnvVar(common.OrganizationID)

// Base model contains defines common fields across tables
type Base struct {
	CreatedAt time.Time `gorm:"column:created;not null"`
	UpdatedAt time.Time `gorm:"column:updated;not null"`
	CreatedBy *string   `gorm:"column:created_by"`
	UpdatedBy *string   `gorm:"column:updated_by"`
}

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	Base

	FacilityID *string `gorm:"primaryKey;unique;column:id"`
	Name       string  `gorm:"column:name;unique;not null"`
	// MFL Code for Kenyan facilities, globally unique
	Code               int    `gorm:"unique;column:mfl_code;not null"`
	Active             bool   `gorm:"column:active;not null"`
	County             string `gorm:"column:county;not null"` // TODO: Controlled list of counties
	Phone              string `gorm:"column:phone"`
	Description        string `gorm:"column:description;not null"`
	FHIROrganisationID string `gorm:"column:fhir_organization_id"`

	// Foreign Keys
	OrganisationID string `gorm:"column:organisation_id;not null"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *Facility) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.FacilityID = &id
	f.OrganisationID = OrganizationID
	return
}

// TableName customizes how the table name is generated
func (Facility) TableName() string {
	return "common_facility"
}

// FacilityAttachment holds any document attached to a facility.
type FacilityAttachment struct {
	Base

	ID           *string   `gorm:"primaryKey;column:id"`
	Title        string    `gorm:"column:title;not null"`
	Active       bool      `gorm:"column:active;not null"`
	Data         string    `gorm:"column:data;not null"`
	Description  string    `gorm:"column:description"`
	CreationDate time.Time `gorm:"creation_date:title;not null"`
	AspectRatio  string    `gorm:"column:aspect_ratio"`
	Notes        string    `gorm:"column:notes"`

	FacilityID     string `gorm:"column:facility_id;not null"`
	OrganisationID string `gorm:"column:organisation_id;not null"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *FacilityAttachment) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.ID = &id
	f.OrganisationID = OrganizationID
	return nil
}

// TableName customizes how the table name is generated
func (FacilityAttachment) TableName() string {
	return "common_facilityattachment"
}

// AuditLog is used to record all sensitive changes e.g
// - changing a client's treatment buddy
// - changing a client's facility
// - deactivating a client
// - changing a client's assigned community health volunteer
// Rules of thumb: is there a need to find out what/when/why something
// occurred? Is a mistake potentially serious? Is there potential for fraud?
type AuditLog struct {
	Base

	ID         *string      `gorm:"primaryKey;column:id"`
	Active     bool         `gorm:"column:active;not null"`
	Timestamp  time.Time    `gorm:"column:timestamp;not null"`
	RecordType string       `gorm:"column:record_type;not null"`
	Notes      string       `gorm:"column:notes"`
	Payload    pgtype.JSONB `gorm:"column:payload"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.ID = &id
	f.OrganisationID = OrganizationID
	return nil
}

// TableName customizes how the table name is generated
func (AuditLog) TableName() string {
	return "common_auditlog"
}

// Address holds the details of an organization/facility
// Types include:- Postal, Physical or Both
type Address struct {
	Base

	ID          *string `gorm:"primaryKey;column:id"`
	Active      bool    `gorm:"column:active;not null"`
	AddressType string  `gorm:"column:address_type;not null"`
	Text        string  `gorm:"column:text;not null"`
	PostalCode  string  `gorm:"column:postal_code;not null"`
	Country     string  `gorm:"column:country;not null"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *Address) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.ID = &id
	f.OrganisationID = OrganizationID
	return nil
}

// TableName customizes how the table name is generated
func (Address) TableName() string {
	return "common_address"
}

// User represents the table data structure for a user
type User struct {
	// Base

	UserID *string `gorm:"primaryKey;unique;column:id"` // globally unique ID

	Username string `gorm:"column:username;unique;not null"` // @handle, also globally unique; nickname

	// DisplayName string `gorm:"column:display_name"` // user's preferred display name

	FirstName  string `gorm:"column:first_name;not null"` // given name
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name;not null"`

	UserType enums.UsersType `gorm:"column:user_type;not null"`

	Gender enumutils.Gender `gorm:"column:gender;not null"`

	Active bool `gorm:"column:is_active;not null"`

	Contacts Contact `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"` // TODO: validate, ensure

	// for the preferred language list, order matters
	// Languages pq.StringArray `gorm:"type:text[];column:languages;not null"` // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens pq.StringArray `gorm:"type:text[];column:push_tokens"`

	// when a user logs in successfully, set this
	LastSuccessfulLogin *time.Time `gorm:"type:time;column:last_successful_login"`

	// whenever there is a failed login (e.g bad PIN), set this
	// reset to null / blank when they succeed at logging in
	LastFailedLogin *time.Time `gorm:"type:time;column:last_failed_login"`

	// each time there is a failed login, **increment** this
	// set to zero after successful login
	FailedLoginCount int `gorm:"column:failed_login_count"`

	// calculated each time there is a failed login
	NextAllowedLogin *time.Time `gorm:"type:time;column:next_allowed_login"`

	TermsAccepted          bool            `gorm:"type:bool;column:terms_accepted;not null"`
	Flavour                feedlib.Flavour `gorm:"column:flavour;not null"`
	Avatar                 string          `gorm:"column:avatar"`
	IsSuspended            bool            `gorm:"column:is_suspended;not null"`
	PinChangeRequired      bool            `gorm:"column:pin_change_required"`
	HasSetPin              bool            `gorm:"column:has_set_pin"`
	HasSetSecurityQuestion bool            `gorm:"column:has_set_security_questions"`
	IsPhoneVerified        bool            `gorm:"column:is_phone_verified"`

	Password            string     `gorm:"column:password"`
	IsSuperuser         bool       `gorm:"column:is_superuser"`
	IsStaff             bool       `gorm:"column:is_staff"`
	Email               string     `gorm:"column:email"`
	DateJoined          string     `gorm:"column:date_joined"`
	Name                string     `gorm:"column:name"`
	IsApproved          bool       `gorm:"column:is_approved"`
	ApprovalNotified    bool       `gorm:"column:approval_notified"`
	Handle              string     `gorm:"column:handle"`
	DateOfBirth         *time.Time `gorm:"column:date_of_birth"`
	FailedSecurityCount int        `gorm:"column:failed_security_count"`
	PinUpdateRequired   bool       `gorm:"column:pin_update_required"`
	HasSetNickname      bool       `gorm:"column:has_set_nickname"`

	OrganisationID  string `gorm:"column:organisation_id"`
	AcceptedTermsID *int   `gorm:"column:accepted_terms_of_service_id"` // foreign key to version of terms they accepted
}

// BeforeCreate is a hook run before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	u.UserID = &id
	u.OrganisationID = OrganizationID
	u.IsSuperuser = false
	u.IsStaff = false
	u.DateJoined = time.Now().UTC().Format(time.RFC1123Z)
	u.IsApproved = false
	u.ApprovalNotified = false
	u.Handle = "@" + u.Username

	login := time.Now()

	u.NextAllowedLogin = &login
	u.FailedLoginCount = 0
	u.FailedSecurityCount = 0
	u.IsApproved = false
	u.ApprovalNotified = false
	u.TermsAccepted = false
	u.IsSuspended = false
	u.PinChangeRequired = true
	u.HasSetPin = false
	u.IsPhoneVerified = false
	u.HasSetSecurityQuestion = false
	u.PinUpdateRequired = false

	return
}

// BeforeDelete hook is run before deleting a user profile
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Unscoped().Where(&PINData{UserID: *u.UserID}).Delete(&PINData{})
	tx.Unscoped().Where(&SecurityQuestionResponse{UserID: *u.UserID}).Delete(&SecurityQuestionResponse{})
	tx.Unscoped().Where(&UserOTP{UserID: *u.UserID}).Delete(&UserOTP{})
	tx.Unscoped().Where(&Notification{UserID: u.UserID}).Delete(&Notification{})
	tx.Unscoped().Where(&UserSurvey{UserID: *u.UserID}).Delete(&UserSurvey{})
	tx.Unscoped().Where(&AuthorityRoleUser{UserID: u.UserID}).Delete(&AuthorityRoleUser{})
	tx.Unscoped().Where(&Metric{UserID: u.UserID}).Delete(&Metric{})

	return
}

// TableName customizes how the table name is generated
func (User) TableName() string {
	return "users_user"
}

// Contact hold contact information/details for users
type Contact struct {
	Base

	ContactID    *string `gorm:"primaryKey;unique;column:id"`
	ContactType  string  `gorm:"column:contact_type;not null"`         // TODO enum
	ContactValue string  `gorm:"unique;column:contact_value;not null"` // TODO Validate: phones are E164, emails are valid
	Active       bool    `gorm:"column:active;not null"`
	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool            `gorm:"column:opted_in;not null"`
	Flavour feedlib.Flavour `gorm:"column:flavour"`

	UserID         *string `gorm:"column:user_id;not null"`
	OrganisationID string  `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a new contact
func (c *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ContactID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName customizes how the table name is generated
func (Contact) TableName() string {
	return "common_contact"
}

// PINData is the PIN's gorm data model.
type PINData struct {
	Base

	PINDataID *int            `gorm:"primaryKey;unique;column:id;autoincrement"`
	HashedPIN string          `gorm:"column:hashed_pin;not null"`
	ValidFrom time.Time       `gorm:"column:valid_from;not null"`
	ValidTo   time.Time       `gorm:"column:valid_to;not null"`
	IsValid   bool            `gorm:"column:active;not null"`
	Flavour   feedlib.Flavour `gorm:"column:flavour;not null"`
	Salt      string          `gorm:"column:salt;not null"`

	UserID string `gorm:"column:user_id;not null"`
}

// TableName customizes how the table name is generated
func (PINData) TableName() string {
	return "users_userpin"
}

// TermsOfService is the gorms terms of service model
type TermsOfService struct {
	Base

	TermsID   *int            `gorm:"primaryKey;unique;column:id;autoincrement"`
	Text      *string         `gorm:"column:text;not null"`
	Flavour   feedlib.Flavour `gorm:"column:flavour;not null"`
	ValidFrom *time.Time      `gorm:"column:valid_from;not null"`
	ValidTo   *time.Time      `gorm:"column:valid_to;not null"`
	Active    bool            `gorm:"column:active;not null"`
}

// TableName customizes how the table name is generated
func (TermsOfService) TableName() string {
	return "users_termsofservice"
}

// Organisation maps the organization table. This will be useful when running integration
// tests since many models have an organization ID as a foreign key.
type Organisation struct {
	Base

	ID              *string `gorm:"primaryKey;unique;column:id"`
	Active          bool    `gorm:"column:active;not null"`
	OrgCode         string  `gorm:"column:org_code;not null;unique"`
	Code            int     `gorm:"column:code;not null;unique"`
	Name            string  `gorm:"column:organisation_name;not null;unique"`
	EmailAddress    string  `gorm:"column:email_address;not null"`
	PhoneNumber     string  `gorm:"column:phone_number;not null"`
	PostalAddress   string  `gorm:"column:postal_address;not null"`
	PhysicalAddress string  `gorm:"column:physical_address;not null"`
	DefaultCountry  string  `gorm:"column:default_country;not null"`
}

// BeforeCreate is a hook run before creating a new organisation
func (t *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = &OrganizationID
	return
}

// TableName customizes how the table name is generated
func (Organisation) TableName() string {
	return "common_organisation"
}

// SecurityQuestion is the gorms security question model
type SecurityQuestion struct {
	Base

	SecurityQuestionID *string                            `gorm:"column:id"`
	QuestionStem       string                             `gorm:"column:stem"`
	Description        string                             `gorm:"column:description"` // help text
	ResponseType       enums.SecurityQuestionResponseType `gorm:"column:response_type"`
	Flavour            feedlib.Flavour                    `gorm:"column:flavour"`
	Active             bool                               `gorm:"column:active"`
	Sequence           *int                               `gorm:"column:sequence"` // for sorting

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating security question
func (s *SecurityQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	s.SecurityQuestionID = &id
	s.OrganisationID = OrganizationID
	return
}

// TableName customizes how the table name is generated
func (SecurityQuestion) TableName() string {
	return "clients_securityquestion"
}

// UserOTP maps the schema for the table that stores the user's OTP
type UserOTP struct {
	Base

	OTPID       int             `gorm:"unique;column:id;autoincrement"`
	Valid       bool            `gorm:"column:is_valid"`
	GeneratedAt time.Time       `gorm:"column:generated_at"`
	ValidUntil  time.Time       `gorm:"column:valid_until"`
	Channel     string          `gorm:"column:channel"`
	Flavour     feedlib.Flavour `gorm:"column:flavour"`
	PhoneNumber string          `gorm:"column:phonenumber"`
	OTP         string          `gorm:"column:otp"`

	UserID string `gorm:"column:user_id"`
}

// TableName customizes how the table name is generated
func (UserOTP) TableName() string {
	return "users_userotp"
}

// SecurityQuestionResponse maps the schema for the table that stores the security question
// responses
type SecurityQuestionResponse struct {
	Base

	ResponseID string    `gorm:"column:id"`
	QuestionID string    `gorm:"column:question_id"`
	Active     bool      `gorm:"column:active"`
	Response   string    `gorm:"column:response"`
	Timestamp  time.Time `gorm:"column:timestamp"`
	IsCorrect  bool      `gorm:"column:is_correct"`

	UserID         string `gorm:"column:user_id"`
	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate function is called before creating a security question response
// It generates the organisation id and response ID
func (s *SecurityQuestionResponse) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	s.ResponseID = id
	s.OrganisationID = OrganizationID
	s.Timestamp = time.Now()
	// is_correct default to true since the user setting the security question responses will initially set
	// them correctly
	// this field will help during verification of security question responses whe a user is resetting the
	// pin. it will change to false if they answer any of the security questions wrong
	s.IsCorrect = true
	return
}

// TableName customizes how the table name is generated
func (SecurityQuestionResponse) TableName() string {
	return "clients_securityquestionresponse"
}

// Client holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
// It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type Client struct {
	Base

	ID *string `gorm:"primaryKey;unique;column:id"`

	Active bool `gorm:"column:active"`

	ClientTypes pq.StringArray `gorm:"type:text[];column:client_types"`

	TreatmentEnrollmentDate *time.Time `gorm:"column:enrollment_date"`

	FHIRPatientID *string `gorm:"column:fhir_patient_id"`

	HealthRecordID *string `gorm:"column:emr_health_record_id"`

	TreatmentBuddy string `gorm:"column:treatment_buddy"`

	ClientCounselled bool `gorm:"column:counselled"`

	FacilityID string `gorm:"column:current_facility_id"`

	CHVUserID *string `gorm:"column:chv_id"`

	CaregiverID *string `gorm:"column:caregiver_id"`

	OrganisationID string `gorm:"column:organisation_id"`

	UserID *string `gorm:"column:user_id;not null"`
	User   User    `gorm:"ForeignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
}

// BeforeCreate is a hook run before creating a client profile
func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = &id
	c.OrganisationID = OrganizationID
	return
}

// BeforeDelete hook is called when deleting a record
func (c *Client) BeforeDelete(tx *gorm.DB) (err error) {
	clientID := *c.ID

	var clientProfile Client
	var clientIdentifiers ClientIdentifiers
	var clientRelatedPerson ClientRelatedPerson
	tx.Model(&ClientIdentifiers{}).Where(&ClientIdentifiers{ClientID: &clientID}).Find(&clientIdentifiers)
	tx.Model(&Client{}).Preload(clause.Associations).Where(&Client{ID: &clientID}).Find(&clientProfile)
	tx.Model(&ClientRelatedPerson{}).Where(&ClientRelatedPerson{ClientID: &clientID}).Find(&clientRelatedPerson)
	userID := clientProfile.User.UserID

	tx.Unscoped().Select(clause.Associations).Where(&Contact{UserID: userID}).Delete(&Contact{})
	tx.Unscoped().Where(&ClientContacts{ClientID: &clientID}).Delete(&ClientContacts{})
	tx.Unscoped().Where(&ClientFacility{ClientID: clientID}).Delete(&ClientFacility{})
	tx.Unscoped().Where("identifier_id", clientIdentifiers.IdentifierID).Delete(&ClientIdentifiers{})
	tx.Unscoped().Where("id", clientIdentifiers.IdentifierID).Delete(&Identifier{})
	tx.Unscoped().Where(&ClientRelatedPerson{ClientID: &clientID}).Delete(&ClientRelatedPerson{})
	tx.Unscoped().Where(&RelatedPersonAddresses{RelatedPersonID: clientRelatedPerson.RelatedPersonID}).Delete(&RelatedPersonAddresses{})
	tx.Unscoped().Where(&RelatedPersonContacts{RelatedPersonID: clientRelatedPerson.RelatedPersonID}).Delete(&RelatedPersonContacts{})
	tx.Unscoped().Where(&ClientAddress{ClientID: clientID}).Delete(&ClientAddress{})
	tx.Unscoped().Where(&ClientHealthDiaryEntry{ClientID: clientID}).Delete(&ClientHealthDiaryEntry{})
	tx.Unscoped().Where(&ClientServiceRequest{ClientID: clientID}).Delete(&ClientServiceRequest{})
	tx.Unscoped().Where(&Appointment{ClientID: clientID}).Delete(&Appointment{})
	tx.Unscoped().Where(&ScreeningToolsResponse{ClientID: clientID}).Delete(&ScreeningToolsResponse{})

	return
}

// TableName references the table that we map data from
func (Client) TableName() string {
	return "clients_client"
}

// ClientFacility represents the client facility table
type ClientFacility struct {
	Base

	ID     *string `gorm:"column:id"`
	Active bool    `gorm:"column:active"`

	OrganisationID string `gorm:"column:organisation_id"`
	ClientID       string `gorm:"column:client_id"`
	FacilityID     string `gorm:"column:facility_id"`
}

// TableName represents the client facility table name
func (ClientFacility) TableName() string {
	return "clients_clientfacility"
}

// ClientAddress represents a through table that holds addresses that belong to a client
type ClientAddress struct {
	ID        *string `gorm:"column:id"`
	ClientID  string  `gorm:"column:client_id"`
	AddressID string  `gorm:"column:address_id"`
}

// TableName composes the table's name
func (ClientAddress) TableName() string {
	return "clients_client_addresses"
}

// StaffProfile represents the staff profile model
type StaffProfile struct {
	Base

	ID *string `gorm:"column:id"`

	Active bool `gorm:"column:active"`

	StaffNumber string `gorm:"column:staff_number"`

	Facilities []Facility `gorm:"ForeignKey:FacilityID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"` // TODO: needs at least one

	DefaultFacilityID string `gorm:"column:default_facility_id"` // TODO: required, FK to facility

	OrganisationID string `gorm:"column:organisation_id"`

	UserID      string `gorm:"column:user_id"` // foreign key to user
	UserProfile User   `gorm:"ForeignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
}

// BeforeCreate is a hook run before creating a staff profile
func (s *StaffProfile) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	s.ID = &id
	s.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (StaffProfile) TableName() string {
	return "staff_staff"
}

// ClientHealthDiaryEntry models a client's health diary entry
type ClientHealthDiaryEntry struct {
	Base

	ClientHealthDiaryEntryID *string    `gorm:"column:id"`
	Active                   bool       `gorm:"column:active"`
	Mood                     string     `gorm:"column:mood"`
	Note                     string     `gorm:"column:note"`
	EntryType                string     `gorm:"column:entry_type"`
	ShareWithHealthWorker    bool       `gorm:"column:share_with_health_worker"`
	SharedAt                 *time.Time `gorm:"column:shared_at"`

	ClientID       string `gorm:"column:client_id"`
	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a client Health Diary Entry
func (c *ClientHealthDiaryEntry) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ClientHealthDiaryEntryID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ClientHealthDiaryEntry) TableName() string {
	return "clients_healthdiaryentry"
}

// ClientServiceRequest maps the client service request table. It is used to
// store the tasks for the healthcare staff on the platform
type ClientServiceRequest struct {
	Base

	ID             *string    `gorm:"column:id"`
	Active         bool       `gorm:"column:active"`
	RequestType    string     `gorm:"column:request_type"`
	Request        string     `gorm:"column:request"`
	Status         string     `gorm:"column:status"`
	InProgressAt   *time.Time `gorm:"column:in_progress_at"`
	ResolvedAt     *time.Time `gorm:"column:resolved_at"`
	InProgressByID *string    `gorm:"column:in_progress_by_id"`
	Meta           string     `gorm:"column:meta"`

	OrganisationID string  `gorm:"column:organisation_id"`
	ResolvedByID   *string `gorm:"column:resolved_by_id"`
	FacilityID     string  `gorm:"column:facility_id"`
	ClientID       string  `gorm:"column:client_id"`
}

// BeforeCreate is a hook called before creating a service request.
func (c *ClientServiceRequest) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ClientServiceRequest) TableName() string {
	return "clients_servicerequest"
}

// StaffServiceRequest maps the staffs's service request table. It is used to
// store the tasks for the healthcare staff on the platform
type StaffServiceRequest struct {
	Base

	ID          *string    `gorm:"column:id"`
	Active      bool       `gorm:"column:active"`
	RequestType string     `gorm:"column:request_type"`
	Request     string     `gorm:"column:request"`
	Status      string     `gorm:"column:status"`
	ResolvedAt  *time.Time `gorm:"column:resolved_at"`
	Meta        string     `gorm:"column:meta"`

	StaffID           string  `gorm:"column:staff_id"`
	OrganisationID    string  `gorm:"column:organisation_id"`
	ResolvedByID      *string `gorm:"column:resolved_by_id"`
	DefaultFacilityID *string `gorm:"column:facility_id"`
}

// BeforeCreate is a hook called before creating a service request.
func (c *StaffServiceRequest) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (StaffServiceRequest) TableName() string {
	return "staff_servicerequest"
}

// ClientHealthDiaryQuote is the gorms client health diary quotes model
type ClientHealthDiaryQuote struct {
	Base

	ClientHealthDiaryQuoteID *string `gorm:"column:id"`
	Active                   bool    `gorm:"column:active"`
	Quote                    string  `gorm:"column:quote"`
	Author                   string  `gorm:"column:by"`

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating view count
func (c *ClientHealthDiaryQuote) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ClientHealthDiaryQuoteID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ClientHealthDiaryQuote) TableName() string {
	return "clients_healthdiaryquote"
}

// Caregiver is the gorms caregiver model
type Caregiver struct {
	Base
	CaregiverID   *string             `gorm:"column:id"`
	FirstName     string              `gorm:"column:first_name"`
	LastName      string              `gorm:"column:last_name"`
	PhoneNumber   string              `gorm:"column:phone_number"`
	CaregiverType enums.CaregiverType `gorm:"column:caregiver_type"`
	Active        bool                `gorm:"column:active"`

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating Caregiver content
func (c *Caregiver) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.CaregiverID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Caregiver) TableName() string {
	return "clients_caregiver"
}

// AuthorityRole is the gorms authority role model
type AuthorityRole struct {
	Base
	AuthorityRoleID *string `gorm:"column:id"`
	Name            string  `gorm:"column:name"`
	Active          bool    `gorm:"column:active"`

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating authority role
func (c *AuthorityRole) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.AuthorityRoleID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (AuthorityRole) TableName() string {
	return "authority_authorityrole"
}

// AuthorityPermission is the gorms authority permission model
type AuthorityPermission struct {
	Base
	AuthorityPermissionID *string `gorm:"column:id"`
	Name                  string  `gorm:"column:name"`

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating authority permission
func (c *AuthorityPermission) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.AuthorityPermissionID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (AuthorityPermission) TableName() string {
	return "authority_authoritypermission"
}

// AuthorityRoleUser is the gorms authority role user model
type AuthorityRoleUser struct {
	ID     int     `gorm:"primaryKey;column:id;autoincrement"`
	UserID *string `gorm:"column:user_id"`
	RoleID *string `gorm:"column:authorityrole_id"`
}

// TableName references the table that we map data from
func (AuthorityRoleUser) TableName() string {
	return "authority_authorityrole_users"
}

// Community defines the payload to create a channel
type Community struct {
	Base

	ID           string         `gorm:"primaryKey;column:id"`
	Name         string         `gorm:"column:name"`
	Description  string         `gorm:"column:description"`
	Active       bool           `gorm:"column:active"`
	MinimumAge   int            `gorm:"column:min_age"`
	MaximumAge   int            `gorm:"column:max_age"`
	Gender       pq.StringArray `gorm:"type:text[];column:gender"`
	ClientTypes  pq.StringArray `gorm:"type:text[];column:client_types"`
	InviteOnly   bool           `gorm:"column:invite_only"`
	Discoverable bool           `gorm:"column:discoverable"`

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a community
func (c *Community) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Community) TableName() string {
	return "communities_community"
}

// Identifier is the the table used to store a user's identifying documents
type Identifier struct {
	Base

	ID                  string    `gorm:"primaryKey;column:id;"`
	Active              bool      `gorm:"column:active;not null"`
	IdentifierType      string    `gorm:"column:identifier_type;not null"`
	IdentifierValue     string    `gorm:"column:identifier_value;not null"`
	IdentifierUse       string    `gorm:"column:identifier_use;not null"`
	Description         string    `gorm:"column:description;not null"`
	ValidFrom           time.Time `gorm:"column:valid_from;not null"`
	ValidTo             time.Time `gorm:"column:valid_to"`
	IsPrimaryIdentifier bool      `gorm:"column:is_primary_identifier"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
}

// TableName references the table that we map data from
func (i *Identifier) TableName() string {
	return "clients_identifier"
}

// BeforeCreate is a hook run before creating authority permission
func (i *Identifier) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()

	i.ID = id
	i.OrganisationID = OrganizationID

	return
}

// ClientIdentifiers links a client with their identifiers
type ClientIdentifiers struct {
	ID           int     `gorm:"primaryKey;column:id;autoincrement"`
	ClientID     *string `gorm:"column:client_id"`
	IdentifierID *string `gorm:"column:identifier_id"`
}

// TableName references the table that we map data from
func (c *ClientIdentifiers) TableName() string {
	return "clients_client_identifiers"
}

// ClientRelatedPerson links a client with their related person e.g next of kin
type ClientRelatedPerson struct {
	ID              int     `gorm:"primaryKey;column:id;autoincrement"`
	ClientID        *string `gorm:"column:client_id"`
	RelatedPersonID *string `gorm:"column:relatedperson_id"`
}

// TableName references the table that we map data from
func (c *ClientRelatedPerson) TableName() string {
	return "clients_client_related_persons"
}

// ClientContacts links a client with their contacts
type ClientContacts struct {
	ID        int     `gorm:"primaryKey;column:id;autoincrement"`
	ClientID  *string `gorm:"column:client_id"`
	ContactID *string `gorm:"column:contact_id"`
}

// TableName references the table that we map data from
func (c *ClientContacts) TableName() string {
	return "clients_client_contacts"
}

// RelatedPerson represents information for a person related to another user
type RelatedPerson struct {
	Base

	ID               string `gorm:"primaryKey;column:id;"`
	Active           bool   `gorm:"column:active;not null"`
	FirstName        string `gorm:"column:first_name"`
	LastName         string `gorm:"column:last_name"`
	OtherName        string `gorm:"column:other_name"`
	Gender           string `gorm:"column:gender"`
	RelationshipType string `gorm:"column:relationship_type"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
}

// TableName references the table that we map data from
func (r *RelatedPerson) TableName() string {
	return "clients_relatedperson"
}

// BeforeCreate is a hook run before creating a related person
func (r *RelatedPerson) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	r.ID = id
	r.OrganisationID = OrganizationID

	return
}

// RelatedPersonContacts links a related person with their contact
type RelatedPersonContacts struct {
	ID              int     `gorm:"primaryKey;column:id;autoincrement"`
	RelatedPersonID *string `gorm:"column:relatedperson_id"`
	ContactID       *string `gorm:"column:contact_id"`
}

// TableName references the table that we map data from
func (r *RelatedPersonContacts) TableName() string {
	return "clients_relatedperson_contacts"
}

// RelatedPersonAddresses links a related person with their addresses
type RelatedPersonAddresses struct {
	ID              int     `gorm:"primaryKey;column:id;autoincrement"`
	RelatedPersonID *string `gorm:"column:relatedperson_id"`
	AddressID       *string `gorm:"column:address_id"`
}

// TableName references the table that we map data from
func (r *RelatedPersonAddresses) TableName() string {
	return "clients_relatedperson_addresses"
}

// ScreeningToolQuestion defines the payload to create screening tools questions
type ScreeningToolQuestion struct {
	Base

	ID               string `gorm:"primaryKey;column:id"`
	Question         string `gorm:"column:question"`
	ToolType         string `gorm:"column:tool_type"`
	ResponseChoices  string `gorm:"column:response_choices"`
	ResponseType     string `gorm:"column:response_type"`
	ResponseCategory string `gorm:"column:response_category"`
	Sequence         int    `gorm:"column:sequence"`
	Active           bool   `gorm:"column:active"`
	Meta             string `gorm:"column:meta"`

	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a screening tools question
func (c *ScreeningToolQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ScreeningToolQuestion) TableName() string {
	return "screeningtools_screeningtoolsquestion"
}

// ScreeningToolsResponse defines the payload to create screening tools responses
type ScreeningToolsResponse struct {
	Base

	ID       string `gorm:"primaryKey;column:id"`
	Response string `gorm:"column:response"`
	Active   bool   `gorm:"column:active"`

	ClientID       string `gorm:"column:client_id"`
	QuestionID     string `gorm:"column:question_id"`
	OrganisationID string `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a screening tools response
func (c *ScreeningToolsResponse) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ScreeningToolsResponse) TableName() string {
	return "screeningtools_screeningtoolsresponse"
}

// Appointment represents a single appointment
type Appointment struct {
	Base

	ID                        string    `gorm:"primaryKey;column:id;"`
	Active                    bool      `gorm:"column:active;not null"`
	ExternalID                string    `gorm:"column:external_id"`
	Reason                    string    `gorm:"column:reason"`
	Provider                  string    `gorm:"column:provider"`
	Date                      time.Time `gorm:"column:date"`
	HasRescheduledAppointment bool      `gorm:"column:has_rescheduled_appointment"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
	ClientID       string `gorm:"column:client_id"`
	FacilityID     string `gorm:"column:facility_id"`
}

// BeforeCreate is a hook run before creating an appointment
func (a *Appointment) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	a.ID = id
	a.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Appointment) TableName() string {
	return "appointments_appointment"
}

// Notification represents a single notification
type Notification struct {
	Base

	ID      string          `gorm:"primaryKey;column:id;"`
	Active  bool            `gorm:"column:active;not null"`
	Title   string          `gorm:"column:title"`
	Body    string          `gorm:"column:body"`
	Type    string          `gorm:"column:notification_type"`
	Flavour feedlib.Flavour `gorm:"column:flavour"`
	IsRead  bool            `gorm:"column:is_read"`

	UserID         *string `gorm:"column:user_id"`
	FacilityID     *string `gorm:"column:facility_id"`
	OrganisationID string  `gorm:"column:organisation_id;not null"`
}

// BeforeCreate is a hook run before creating an appointment
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	n.ID = id
	n.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Notification) TableName() string {
	return "common_notification"
}

// StaffIdentifiers links a staff with their identifiers
type StaffIdentifiers struct {
	ID           int     `gorm:"primaryKey;column:id;autoincrement"`
	StaffID      *string `gorm:"column:staff_id"`
	IdentifierID *string `gorm:"column:identifier_id"`
}

// TableName references the table that we map data from
func (s *StaffIdentifiers) TableName() string {
	return "staff_staff_identifiers"
}

// StaffContacts links a staff with their contacts
type StaffContacts struct {
	ID        int     `gorm:"primaryKey;column:id;autoincrement"`
	StaffID   *string `gorm:"column:staff_id"`
	ContactID *string `gorm:"column:contact_id"`
}

// TableName references the table that we map data from
func (c *StaffContacts) TableName() string {
	return "staff_staff_contacts"
}

// StaffFacilities links a staff with their facilities
type StaffFacilities struct {
	ID         int     `gorm:"primaryKey;column:id;autoincrement"`
	StaffID    *string `gorm:"column:staff_id"`
	FacilityID *string `gorm:"column:facility_id"`
}

// TableName references the table that we map data from
func (s *StaffFacilities) TableName() string {
	return "staff_staff_facilities"
}

// UserSurvey represents a user's surveys database model
type UserSurvey struct {
	Base

	ID           string     `gorm:"id"`
	Active       bool       `gorm:"active"`
	Link         string     `gorm:"link"`
	Title        string     `gorm:"title"`
	Description  string     `gorm:"description"`
	HasSubmitted bool       `gorm:"submitted"`
	FormID       string     `gorm:"form_id"`
	ProjectID    int        `gorm:"project_id"`
	LinkID       int        `gorm:"link_id"`
	Token        string     `gorm:"token"`
	SubmittedAt  *time.Time `gorm:"submitted_at"`

	UserID         string `gorm:"user_id"`
	OrganisationID string `gorm:"organisation_id"`
}

// BeforeCreate is a hook run before creating a user survey model
func (u *UserSurvey) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	u.ID = id
	u.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (UserSurvey) TableName() string {
	return "common_usersurveys"
}

// Metric is a recording of an event that occurs within the platform
type Metric struct {
	Base

	ID        int              `gorm:"primaryKey;column:id;autoincrement"`
	Active    bool             `gorm:"column:active"`
	Type      enums.MetricType `gorm:"column:metric_type"`
	Payload   string           `gorm:"column:payload"`
	Timestamp time.Time        `gorm:"column:timestamp"`

	UserID *string `gorm:"column:user_id"`
}

// TableName references the table that we map data from
func (Metric) TableName() string {
	return "users_metric"
}

// Feedback defines the feedback database model
type Feedback struct {
	Base

	ID                string `gorm:"primaryKey;column:id"`
	Active            bool   `gorm:"column:active"`
	FeedbackType      string `gorm:"column:feedback_type"`
	SatisfactionLevel int    `gorm:"column:satisfaction_level"`
	ServiceName       string `gorm:"column:service_name"`
	Feedback          string `gorm:"column:feedback"`
	RequiresFollowUp  bool   `gorm:"column:requires_followup"`
	PhoneNumber       string `gorm:"column:phone_number"`

	OrganisationID string `gorm:"column:organisation_id"`
	UserID         string `gorm:"column:user_id"`
}

// BeforeCreate is a hook run before creating an appointment
func (f *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.ID = id
	f.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Feedback) TableName() string {
	return "common_feedback"
}

// Questionnaire defines the questionnaire database models
type Questionnaire struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID          string `gorm:"primaryKey;column:id"`
	Active      bool   `gorm:"column:active"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

// BeforeCreate is a hook run before creating a questionnaire
func (q *Questionnaire) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	q.ID = id
	q.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Questionnaire) TableName() string {
	return "questionnaires_questionnaire"
}

// ScreeningTool defines the screening tool database models
type ScreeningTool struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID              string         `gorm:"primaryKey;column:id"`
	Active          bool           `gorm:"column:active"`
	QuestionnaireID string         `gorm:"column:questionnaire_id"`
	Threshold       int            `gorm:"column:threshold"`
	ClientTypes     pq.StringArray `gorm:"type:text[];column:client_types"`
	Genders         pq.StringArray `gorm:"type:text[];column:genders"`
	MinimumAge      int            `gorm:"column:min_age"`
	MaximumAge      int            `gorm:"column:max_age"`
}

// BeforeCreate is a hook run before creating a screening tool
func (f *ScreeningTool) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	f.ID = id
	f.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ScreeningTool) TableName() string {
	return "questionnaires_screeningtool"
}

// Question defines the question database models
type Question struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID                string `gorm:"primaryKey;column:id"`
	Active            bool   `gorm:"column:active"`
	QuestionnaireID   string `gorm:"column:questionnaire_id"`
	Text              string `gorm:"column:text"`
	QuestionType      string `gorm:"column:question_type"`
	ResponseValueType string `gorm:"column:response_value_type"`
	SelectMultiple    bool   `gorm:"column:select_multiple"`
	Required          bool   `gorm:"column:required"`
	Sequence          int    `gorm:"column:sequence"`
}

// BeforeCreate is a hook run before creating a question
func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	q.ID = id
	q.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Question) TableName() string {
	return "questionnaires_question"
}

// QuestionInputChoice defines the question input choice database models
type QuestionInputChoice struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID         string `gorm:"primaryKey;column:id"`
	Active     bool   `gorm:"column:active"`
	QuestionID string `gorm:"column:question_id"`
	Choice     string `gorm:"column:choice"`
	Value      string `gorm:"column:value"`
	Score      int    `gorm:"column:score"`
}

// BeforeCreate is a hook run before creating a question input choice
func (q *QuestionInputChoice) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	q.ID = id
	q.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (QuestionInputChoice) TableName() string {
	return "questionnaires_questioninputchoice"
}

// ScreeningToolResponse defines the screening tool response database models
type ScreeningToolResponse struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID              string `gorm:"primaryKey;column:id"`
	Active          bool   `gorm:"column:active"`
	ScreeningToolID string `gorm:"column:screeningtool_id"`
	FacilityID      string `gorm:"column:facility_id"`
	ClientID        string `gorm:"column:client_id"`
	AggregateScore  int    `gorm:"column:aggregate_score"`
}

// BeforeCreate is a hook run before creating a screening tool response
func (s *ScreeningToolResponse) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	s.ID = id
	s.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ScreeningToolResponse) TableName() string {
	return "questionnaires_screeningtoolresponse"
}

// ScreeningToolQuestionResponse defines the screening tool question response database models
type ScreeningToolQuestionResponse struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID                      string `gorm:"primaryKey;column:id"`
	Active                  bool   `gorm:"column:active"`
	ScreeningToolResponseID string `gorm:"column:screeningtoolresponse_id"`
	QuestionID              string `gorm:"column:question_id"`
	Response                string `gorm:"column:response"`
	Score                   int    `gorm:"column:score"`
}

// BeforeCreate is a hook run before creating a screening tool question response
func (s *ScreeningToolQuestionResponse) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	s.ID = id
	s.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ScreeningToolQuestionResponse) TableName() string {
	return "questionnaires_screeningtoolquestionresponse"
}

// ClientFacilities links a client with their facilities
// it is a through table
type ClientFacilities struct {
	ID         int     `gorm:"primaryKey;column:id;autoincrement"`
	ClientID   *string `gorm:"column:client_id"`
	FacilityID *string `gorm:"column:facility_id"`
}

// TableName references the table that we map data from
func (s *ClientFacilities) TableName() string {
	return "clients_client_facilities"
}

// NCaregiver is the caregiver profile information for a user
// TODO: remove "N" when original caregiver is removed
type NCaregiver struct {
	Base

	ID              string `gorm:"primaryKey;column:id"`
	Active          bool   `gorm:"column:active"`
	CaregiverNumber string `gorm:"column:caregiver_number"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
	UserID         string `gorm:"column:user_id"`
}

// BeforeCreate is a hook run before creating a caregiver
func (c *NCaregiver) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = id
	c.OrganisationID = OrganizationID

	return nil
}

// TableName references the table name in the database
func (c *NCaregiver) TableName() string {
	return "caregivers_caregiver"
}

// CaregiverClients stores clients assigned to a caregiver and the associated relationship details
type CaregiverClients struct {
	Base

	CaregiverID        string     `gorm:"column:caregiver_id;not null"`
	ClientID           string     `gorm:"column:client_id;not null"`
	Active             bool       `gorm:"column:active"`
	RelationshipType   string     `gorm:"column:relationship_type;not null"`
	CaregiverConsent   *bool      `gorm:"column:caregiver_consent"`
	CaregiverConsentAt *time.Time `gorm:"column:caregiver_consent_at"`
	ClientConsent      *bool      `gorm:"column:client_consent"`
	ClientConsentAt    *time.Time `gorm:"column:client_consent_at"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
	AssignedBy     string `gorm:"column:assigned_by;not null"`
}

// BeforeCreate is a hook run before creating a caregiver client
func (c *CaregiverClients) BeforeCreate(tx *gorm.DB) (err error) {
	c.OrganisationID = OrganizationID

	return nil
}

// TableName references the table name in the database
func (c *CaregiverClients) TableName() string {
	return "caregivers_caregiver_client"
}
