package gorm

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"github.com/lib/pq"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Base model contains defines common fields across tables
type Base struct {
	CreatedAt time.Time  `gorm:"column:created;not null"`
	UpdatedAt time.Time  `gorm:"column:updated;not null"`
	CreatedBy *string    `gorm:"column:created_by"`
	UpdatedBy *string    `gorm:"column:updated_by"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	Base

	FacilityID         *string `gorm:"primaryKey;unique;column:id"`
	Name               string  `gorm:"column:name;unique;not null"`
	Active             bool    `gorm:"column:active;not null"`
	Country            string  `gorm:"column:country;not null"`
	Phone              string  `gorm:"column:phone"`
	Description        string  `gorm:"column:description;not null"`
	FHIROrganisationID string  `gorm:"column:fhir_organization_id"`
	Identifier         []*FacilityIdentifier
	Coordinates        *FacilityCoordinates
}

// BeforeCreate is a hook run before creating a new facility
func (f *Facility) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.CreatedBy = userID
	}

	id := uuid.New().String()
	f.FacilityID = &id

	return
}

// BeforeUpdate is a hook called before updating Facility.
func (f *Facility) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (Facility) TableName() string {
	return "common_facility"
}

// FacilityIdentifier stores a facilities identifiers
type FacilityIdentifier struct {
	Base

	ID     string `gorm:"primaryKey;unique;column:id"`
	Active bool   `gorm:"column:active;not null"`
	Type   string `gorm:"column:identifier_type;not null"`
	Value  string `gorm:"column:identifier_value;not null"`

	FacilityID string `gorm:"column:facility_id;not null"`
}

// BeforeCreate is a hook run before creating a new facility
func (f *FacilityIdentifier) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.CreatedBy = userID
	}
	id := uuid.New().String()
	f.ID = id

	return
}

// BeforeUpdate is a hook called before updating FacilityIdentifier.
func (f *FacilityIdentifier) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (FacilityIdentifier) TableName() string {
	return "common_facility_identifier"
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
func (a *AuditLog) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.CreatedBy = userID
	}
	id := uuid.New().String()
	a.ID = &id

	return nil
}

// BeforeUpdate is a hook called before updating AuditLog.
func (a *AuditLog) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.UpdatedBy = userID
	}
	return
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
func (a *Address) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.CreatedBy = userID
	}
	id := uuid.New().String()
	a.ID = &id
	return nil
}

// BeforeUpdate is a hook called before updating Address.
func (a *Address) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (Address) TableName() string {
	return "common_address"
}

// User represents the table data structure for a user
type User struct {
	Base

	UserID   *string          `gorm:"primaryKey;unique;column:id"`
	Username string           `gorm:"column:username;unique;not null"`
	Email    *string          `gorm:"column:email;unique"`
	Gender   enumutils.Gender `gorm:"column:gender;not null"`
	Active   bool             `gorm:"column:active;not null"`

	Contacts Contact `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"` // TODO: validate, ensure

	// for the preferred language list, order matters
	Languages pq.StringArray `gorm:"type:text[];column:languages;not null"` // TODO: turn this into a slice of enums, start small (en, sw)

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

	TermsAccepted          bool       `gorm:"type:bool;column:terms_accepted;not null"`
	Avatar                 string     `gorm:"column:avatar"`
	IsSuspended            bool       `gorm:"column:is_suspended;not null"`
	PinChangeRequired      bool       `gorm:"column:pin_change_required"`
	HasSetPin              bool       `gorm:"column:has_set_pin"`
	HasSetSecurityQuestion bool       `gorm:"column:has_set_security_questions"`
	HasSetUsername         bool       `gorm:"column:has_set_username"`
	IsPhoneVerified        bool       `gorm:"column:is_phone_verified"`
	IsSuperuser            bool       `gorm:"column:is_superuser"`
	Name                   string     `gorm:"column:name"`
	DateOfBirth            *time.Time `gorm:"column:date_of_birth"`
	FailedSecurityCount    int        `gorm:"column:failed_security_count"`
	PinUpdateRequired      bool       `gorm:"column:pin_update_required"`

	CurrentOrganisationID string `gorm:"column:current_organisation_id"`
	CurrentProgramID      string `gorm:"column:current_program_id"`
	CurrentUserType       string `gorm:"column:current_usertype"`
	AcceptedTermsID       *int   `gorm:"column:accepted_terms_of_service_id"`
}

// BeforeCreate is a hook run before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if err != nil {
		return err
	}
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		u.CreatedBy = userID
	}

	id := uuid.New().String()
	u.UserID = &id
	u.IsSuperuser = false

	login := time.Now()

	u.NextAllowedLogin = &login
	u.FailedLoginCount = 0
	u.FailedSecurityCount = 0
	u.TermsAccepted = false
	u.IsSuspended = false
	u.PinChangeRequired = true
	u.HasSetPin = false
	u.IsPhoneVerified = false
	u.HasSetSecurityQuestion = false
	u.PinUpdateRequired = false

	return
}

// BeforeUpdate is a hook called before updating User.
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		u.UpdatedBy = userID
	}
	return
}

// BeforeDelete hook is run before deleting a user profile
func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Model(&Notification{}).Where("created_by=?", u.UserID).Updates(map[string]interface{}{"created_by": nil, "updated_by": nil})

	tx.Unscoped().Where(&PINData{UserID: *u.UserID}).Delete(&PINData{})
	tx.Unscoped().Where(&SecurityQuestionResponse{UserID: *u.UserID}).Delete(&SecurityQuestionResponse{})
	tx.Unscoped().Where(&UserOTP{UserID: *u.UserID}).Delete(&UserOTP{})
	tx.Unscoped().Where(&Notification{UserID: u.UserID}).Delete(&Notification{})
	tx.Unscoped().Where(&UserSurvey{UserID: *u.UserID}).Delete(&UserSurvey{})
	tx.Unscoped().Where(&Metric{UserID: u.UserID}).Delete(&Metric{})
	tx.Unscoped().Where(&OrganisationUser{UserID: *u.UserID}).Delete(&OrganisationUser{})
	tx.Unscoped().Where(&Contact{UserID: u.UserID}).Delete(&Contact{})
	tx.Unscoped().Where(&Feedback{UserID: *u.UserID}).Delete(&Feedback{})
	tx.Unscoped().Where(&Metric{UserID: u.UserID}).Delete(&Metric{})
	tx.Unscoped().Where(&Caregiver{UserID: *u.UserID}).Delete(&Caregiver{})

	var sessions []*Session
	tx.Unscoped().Where(&Session{UserID: *u.UserID}).Find(&sessions)
	for _, session := range sessions {
		tx.Unscoped().Where(&AccessToken{SessionID: session.ID}).Delete(&AccessToken{})

		tx.Unscoped().Where(&RefreshToken{SessionID: session.ID}).Delete(&RefreshToken{})

		tx.Unscoped().Where(&Session{ID: session.ID}).Delete(&Session{})
	}

	return
}

// TableName customizes how the table name is generated
func (User) TableName() string {
	return "users_user"
}

// OrganisationUser models the relationship between a user and an organisation
type OrganisationUser struct {
	ID             int    `gorm:"column:id;primary_key;autoincrement"`
	OrganisationID string `gorm:"column:organisation_id"`
	UserID         string `gorm:"column:user_id"`
}

// TableName customizes how the table name is generated
func (OrganisationUser) TableName() string {
	return "users_user_organisation"
}

// ProgramFacility models the relationship between a program and a facility
type ProgramFacility struct {
	ID         int    `gorm:"primaryKey;column:id;autoincrement"`
	ProgramID  string `gorm:"column:program_id"`
	FacilityID string `gorm:"column:facility_id"`
}

// TableName customizes how the table name is generated
func (ProgramFacility) TableName() string {
	return "common_program_facility"
}

// Contact hold contact information/details for users
type Contact struct {
	Base

	ID     string `gorm:"primaryKey;unique;column:id"`
	Type   string `gorm:"column:contact_type;not null"`
	Value  string `gorm:"unique;column:contact_value;not null"`
	Active bool   `gorm:"column:active;not null"`
	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `gorm:"column:opted_in;not null"`

	UserID         *string `gorm:"column:user_id;not null"`
	OrganisationID string  `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a new contact
func (c *Contact) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}

	id := uuid.New().String()
	c.ID = id

	return
}

// BeforeUpdate is a hook called before updating Contact.
func (c *Contact) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (Contact) TableName() string {
	return "common_contact"
}

// PINData is the PIN's gorm data model.
type PINData struct {
	Base

	PINDataID *int      `gorm:"primaryKey;unique;column:id;autoincrement"`
	HashedPIN string    `gorm:"column:hashed_pin;not null"`
	ValidFrom time.Time `gorm:"column:valid_from;not null"`
	ValidTo   time.Time `gorm:"column:valid_to;not null"`
	IsValid   bool      `gorm:"column:active;not null"`
	Salt      string    `gorm:"column:salt;not null"`

	UserID string `gorm:"column:user_id;not null"`
}

// BeforeCreate is a hook run before creating a new PINData
func (p *PINData) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		p.CreatedBy = userID
	}
	return
}

// BeforeUpdate is a hook called before updating PINData.
func (p *PINData) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		p.UpdatedBy = userID
	}

	return
}

// TableName customizes how the table name is generated
func (PINData) TableName() string {
	return "users_userpin"
}

// TermsOfService is the gorms terms of service model
type TermsOfService struct {
	Base

	TermsID   *int       `gorm:"primaryKey;unique;column:id;autoincrement"`
	Text      *string    `gorm:"column:text;not null"`
	ValidFrom *time.Time `gorm:"column:valid_from;not null"`
	ValidTo   *time.Time `gorm:"column:valid_to;not null"`
	Active    bool       `gorm:"column:active;not null"`
}

// BeforeCreate is a hook run before creating a new TermsOfService
func (t *TermsOfService) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		t.CreatedBy = userID
	}
	return
}

// BeforeUpdate is a hook called before updating TermsOfService.
func (t *TermsOfService) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		t.UpdatedBy = userID
	}
	return
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
	Code            string  `gorm:"column:org_code;not null;unique"`
	Name            string  `gorm:"column:name;not null;unique"`
	Description     string  `gorm:"column:description"`
	EmailAddress    string  `gorm:"column:email_address;not null"`
	PhoneNumber     string  `gorm:"column:phone_number;not null"`
	PostalAddress   string  `gorm:"column:postal_address;not null"`
	PhysicalAddress string  `gorm:"column:physical_address;not null"`
	DefaultCountry  string  `gorm:"column:default_country;not null"`
}

// BeforeCreate is a hook run before creating a new organisation
func (o *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		o.CreatedBy = userID
	}

	UUID := uuid.New().String()
	o.ID = &UUID

	return
}

// BeforeUpdate is a hook called before updating Organisation.
func (o *Organisation) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		o.UpdatedBy = userID
	}

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
}

// BeforeCreate is a hook run before creating security question
func (s *SecurityQuestion) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}

	id := uuid.New().String()
	s.SecurityQuestionID = &id
	return
}

// BeforeUpdate is a hook called before updating SecurityQuestion.
func (s *SecurityQuestion) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (SecurityQuestion) TableName() string {
	return "common_securityquestion"
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

// BeforeCreate is a hook called before updating UserOTP.
func (u *UserOTP) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		u.CreatedBy = userID
	}
	return
}

// BeforeUpdate is a hook called before updating UserOTP.
func (u *UserOTP) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		u.UpdatedBy = userID
	}
	return
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

	UserID string `gorm:"column:user_id"`
}

// BeforeCreate function is called before creating a security question response
// It generates the organisation id and response ID
func (s *SecurityQuestionResponse) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}

	id := uuid.New().String()
	s.ResponseID = id
	s.Timestamp = time.Now()
	// is_correct default to true since the user setting the security question responses will initially set
	// them correctly
	// this field will help during verification of security question responses whe a user is resetting the
	// pin. it will change to false if they answer any of the security questions wrong
	s.IsCorrect = true
	return
}

// BeforeUpdate is a hook called before updating SecurityQuestionResponse.
func (s *SecurityQuestionResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (SecurityQuestionResponse) TableName() string {
	return "common_securityquestionresponse"
}

// Client holds the details of end users who are not using the system in
// a professional capacity e.g consumers, patients etc.
// It is a linkage model e.g to tie together all of a person's identifiers
// and their health record ID
type Client struct {
	Base

	ID                      *string        `gorm:"primaryKey;unique;column:id"`
	Active                  bool           `gorm:"column:active"`
	ClientTypes             pq.StringArray `gorm:"type:text[];column:client_types"`
	TreatmentEnrollmentDate *time.Time     `gorm:"column:enrollment_date"`
	FHIRPatientID           *string        `gorm:"column:fhir_patient_id"`
	HealthRecordID          *string        `gorm:"column:emr_health_record_id"`
	ClientCounselled        bool           `gorm:"column:counselled"`

	UserID         *string `gorm:"column:user_id;not null"`
	User           User    `gorm:"ForeignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	FacilityID     string  `gorm:"column:current_facility_id"`
	OrganisationID string  `gorm:"column:organisation_id"`
	ProgramID      string  `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a client profile
func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}

	id := uuid.New().String()
	c.ID = &id

	return
}

// BeforeUpdate is a hook called before updating Client.
func (c *Client) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}

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

	var screeningToolResponses []*ScreeningToolResponse
	tx.Model(&ScreeningToolResponse{}).Where(&ScreeningToolResponse{ClientID: clientID}).Find(&screeningToolResponses)
	for _, screeningToolResponse := range screeningToolResponses {
		tx.Unscoped().Where(&ScreeningToolQuestionResponse{ScreeningToolResponseID: screeningToolResponse.ID}).Delete(&ScreeningToolQuestionResponse{})
	}

	tx.Model(&Caregiver{}).Where(&Caregiver{CurrentClient: &clientID}).Updates(map[string]interface{}{"current_client": nil})

	tx.Unscoped().Select(clause.Associations).Where(&Contact{UserID: userID}).Delete(&Contact{})
	tx.Unscoped().Where(&ClientFacility{ClientID: clientID}).Delete(&ClientFacility{})
	tx.Unscoped().Where("identifier_id", clientIdentifiers.IdentifierID).Delete(&ClientIdentifiers{})
	tx.Unscoped().Where("id", clientIdentifiers.IdentifierID).Delete(&Identifier{})
	tx.Unscoped().Where(&ClientRelatedPerson{ClientID: &clientID}).Delete(&ClientRelatedPerson{})
	tx.Unscoped().Where(&RelatedPersonAddresses{RelatedPersonID: clientRelatedPerson.RelatedPersonID}).Delete(&RelatedPersonAddresses{})
	tx.Unscoped().Where(&RelatedPersonContacts{RelatedPersonID: clientRelatedPerson.RelatedPersonID}).Delete(&RelatedPersonContacts{})
	tx.Unscoped().Where(&ClientHealthDiaryEntry{ClientID: clientID}).Delete(&ClientHealthDiaryEntry{})
	tx.Unscoped().Where(&ClientServiceRequest{ClientID: clientID}).Delete(&ClientServiceRequest{})
	tx.Unscoped().Where(&Appointment{ClientID: clientID}).Delete(&Appointment{})
	tx.Unscoped().Where(&ScreeningToolResponse{ClientID: clientID}).Delete(&ScreeningToolResponse{})
	tx.Unscoped().Where(&ClientFacilities{ClientID: &clientID}).Delete(&ClientFacilities{})
	tx.Unscoped().Where(&CaregiverClient{ClientID: clientID}).Delete(&CaregiverClient{})
	tx.Unscoped().Where(&CommunityClient{ClientID: &clientID}).Delete(&CommunityClient{})
	tx.Unscoped().Where(&AuthorityRoleClient{ClientID: &clientID}).Delete(&AuthorityRoleClient{})
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
	ProgramID      string `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a client facility
func (c *ClientFacility) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}

	id := uuid.New().String()
	c.ID = &id
	return
}

// BeforeUpdate is a hook called before updating ClientFacility.
func (c *ClientFacility) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
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

	DefaultFacilityID string `gorm:"column:current_facility_id"` // TODO: required, FK to facility

	OrganisationID string `gorm:"column:organisation_id"`

	UserID      string `gorm:"column:user_id"` // foreign key to user
	UserProfile User   `gorm:"ForeignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`

	ProgramID           string `gorm:"column:program_id"` // foreign key to program
	IsOrganisationAdmin bool   `gorm:"column:is_organisation_admin"`
}

// BeforeCreate is a hook run before creating a staff profile
func (s *StaffProfile) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}

	id := uuid.New().String()
	s.ID = &id

	return
}

// BeforeUpdate is a hook called before updating StaffProfile.
func (s *StaffProfile) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}

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
	ProgramID                string     `gorm:"column:program_id"`
	ClientID                 string     `gorm:"column:client_id"`
	OrganisationID           string     `gorm:"column:organisation_id"`
	CaregiverID              *string    `gorm:"column:caregiver_id"`
}

// BeforeCreate is a hook run before creating a client Health Diary Entry
func (c *ClientHealthDiaryEntry) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}

	id := uuid.New().String()
	c.ClientHealthDiaryEntryID = &id
	return
}

// BeforeUpdate is a hook called before updating ClientHealthDiaryEntry.
func (c *ClientHealthDiaryEntry) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
}

// TableName references the table that we map data from
func (c *ClientHealthDiaryEntry) TableName() string {
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
	ProgramID      string     `gorm:"column:program_id"`
	OrganisationID string     `gorm:"column:organisation_id"`
	ResolvedByID   *string    `gorm:"column:resolved_by_id"`
	FacilityID     string     `gorm:"column:facility_id"`
	ClientID       string     `gorm:"column:client_id"`
	CaregiverID    *string    `gorm:"column:caregiver_id"`
}

// BeforeCreate is a hook called before creating a service request.
func (c *ClientServiceRequest) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}
	id := uuid.New().String()
	c.ID = &id

	return
}

// BeforeUpdate is a hook called before updating a client service request.
func (c *ClientServiceRequest) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
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
	ProgramID         string  `gorm:"column:program_id"`
}

// BeforeCreate is a hook called before creating a service request.
func (s *StaffServiceRequest) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}

	id := uuid.New().String()
	s.ID = &id

	return
}

// BeforeUpdate is a hook called before updating a service request.
func (s *StaffServiceRequest) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}
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
	ProgramID                string  `gorm:"column:program_id"`
	OrganisationID           string  `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating view count
func (c *ClientHealthDiaryQuote) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID

		var userProfile User
		err = tx.Model(User{UserID: userID}).Find(&userProfile).Error
		if err != nil {
			logrus.Println("could not get user profile")
		}
		c.ProgramID = userProfile.CurrentProgramID
	}

	id := uuid.New().String()
	c.ClientHealthDiaryQuoteID = &id

	return
}

// BeforeUpdate is a hook called before updating ClientHealthDiaryQuote.
func (c *ClientHealthDiaryQuote) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
}

// TableName references the table that we map data from
func (ClientHealthDiaryQuote) TableName() string {
	return "clients_healthdiaryquote"
}

// AuthorityRole is the gorms authority role model
type AuthorityRole struct {
	Base
	AuthorityRoleID *string `gorm:"column:id"`
	Name            string  `gorm:"column:name"`
	Active          bool    `gorm:"column:active"`

	OrganisationID string `gorm:"column:organisation_id"`
	ProgramID      string `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating authority role
func (a *AuthorityRole) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.CreatedBy = userID
	}

	id := uuid.New().String()
	a.AuthorityRoleID = &id

	return
}

// BeforeUpdate is a hook called before updating AuthorityRole.
func (a *AuthorityRole) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.UpdatedBy = userID
	}
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
	Description           string  `gorm:"column:description"`
	Category              string  `gorm:"column:category"`
	Scope                 string  `gorm:"column:scope"`
}

// BeforeCreate is a hook run before creating authority permission
func (a *AuthorityPermission) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.CreatedBy = userID
	}
	id := uuid.New().String()
	a.AuthorityPermissionID = &id

	return
}

// BeforeUpdate is a hook called before updating AuthorityPermission.
func (a *AuthorityPermission) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		a.UpdatedBy = userID
	}
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

// AuthorityRoleClient is the gorms authority role client model
type AuthorityRoleClient struct {
	ID       int     `gorm:"primaryKey;column:id;autoincrement"`
	ClientID *string `gorm:"column:client_id"`
	RoleID   *string `gorm:"column:authorityrole_id"`
}

// TableName references the table that we map data from
func (AuthorityRoleClient) TableName() string {
	return "authority_authorityrole_clients"
}

// AuthorityRolePermission is the gorms authority role permission model
type AuthorityRolePermission struct {
	ID           int     `gorm:"primaryKey;column:id;autoincrement"`
	PermissionID *string `gorm:"column:authoritypermission_id"`
	RoleID       *string `gorm:"column:authorityrole_id"`
	Active       bool    `gorm:"column:active"`
}

// TableName references the table that we map data from
func (AuthorityRolePermission) TableName() string {
	return "authority_authorityrole_permissions"
}

// Community defines the payload to create a community
type Community struct {
	Base

	ID             string         `gorm:"primaryKey;column:id"`
	RoomID         string         `json:"room_id"`
	Name           string         `gorm:"column:name"`
	Description    string         `gorm:"column:description"`
	Active         bool           `gorm:"column:active"`
	MinimumAge     int            `gorm:"column:min_age"`
	MaximumAge     int            `gorm:"column:max_age"`
	Gender         pq.StringArray `gorm:"type:text[];column:gender"`
	ClientTypes    pq.StringArray `gorm:"type:text[];column:client_types"`
	ProgramID      string         `gorm:"column:program_id"`
	OrganisationID string         `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating a community
func (c *Community) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}

	id := uuid.New().String()
	c.ID = id

	return
}

// BeforeUpdate is a hook called before updating community.
func (c *Community) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
}

// TableName references the table that we map data from
func (c *Community) TableName() string {
	return "communities_community"
}

// Identifier is the the table used to store a user's identifying documents
type Identifier struct {
	Base

	ID                  string    `gorm:"primaryKey;column:id;"`
	Active              bool      `gorm:"column:active;not null"`
	Type                string    `gorm:"column:identifier_type;not null"`
	Value               string    `gorm:"column:identifier_value;not null"`
	Use                 string    `gorm:"column:identifier_use;not null"`
	Description         string    `gorm:"column:description;not null"`
	ValidFrom           time.Time `gorm:"column:valid_from;not null"`
	ValidTo             time.Time `gorm:"column:valid_to"`
	IsPrimaryIdentifier bool      `gorm:"column:is_primary_identifier"`
	OrganisationID      string    `gorm:"column:organisation_id;not null"`
	ProgramID           string    `gorm:"column:program_id;not null"`
}

// TableName references the table that we map data from
func (i *Identifier) TableName() string {
	return "common_identifiers"
}

// BeforeCreate is a hook run before creating authority permission
func (i *Identifier) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		i.CreatedBy = userID
	}

	id := uuid.New().String()
	i.ID = id

	return
}

// BeforeUpdate is a hook called before updating Identifier.
func (i *Identifier) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		i.UpdatedBy = userID
	}
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
	ProgramID        string `gorm:"column:program_id"`
	OrganisationID   string `gorm:"column:organisation_id;not null"`
}

// TableName references the table that we map data from
func (r *RelatedPerson) TableName() string {
	return "clients_relatedperson"
}

// BeforeCreate is a hook run before creating a related person
func (r *RelatedPerson) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		r.CreatedBy = userID
	}

	id := uuid.New().String()
	r.ID = id

	return
}

// BeforeUpdate is a hook called before updating RelatedPerson.
func (r *RelatedPerson) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		r.UpdatedBy = userID
	}
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
	ProgramID                 string    `gorm:"column:program_id"`
	OrganisationID            string    `gorm:"column:organisation_id;not null"`
	ClientID                  string    `gorm:"column:client_id"`
	FacilityID                string    `gorm:"column:facility_id"`
	CaregiverID               *string   `gorm:"column:caregiver_id"`
}

// BeforeCreate is a hook run before creating an appointment
func (a *Appointment) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	a.ID = id

	return
}

// TableName references the table that we map data from
func (Appointment) TableName() string {
	return "appointments_appointment"
}

// Notification represents a single notification
type Notification struct {
	Base

	ID             string          `gorm:"primaryKey;column:id;"`
	Active         bool            `gorm:"column:active;not null"`
	Title          string          `gorm:"column:title"`
	Body           string          `gorm:"column:body"`
	Type           string          `gorm:"column:notification_type"`
	Flavour        feedlib.Flavour `gorm:"column:flavour"`
	IsRead         bool            `gorm:"column:is_read"`
	ProgramID      string          `gorm:"column:program_id"`
	UserID         *string         `gorm:"column:user_id"`
	FacilityID     *string         `gorm:"column:facility_id"`
	OrganisationID string          `gorm:"column:organisation_id;not null"`
}

// BeforeCreate is a hook run before creating an appointment
func (n *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		n.CreatedBy = userID
	}

	id := uuid.New().String()
	n.ID = id

	return
}

// BeforeUpdate is a hook called before updating Notification.
func (n *Notification) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		n.UpdatedBy = userID
	}
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

	ID             string     `gorm:"id"`
	Active         bool       `gorm:"active"`
	Link           string     `gorm:"link"`
	Title          string     `gorm:"title"`
	Description    string     `gorm:"description"`
	HasSubmitted   bool       `gorm:"submitted"`
	FormID         string     `gorm:"form_id"`
	ProjectID      int        `gorm:"project_id"`
	LinkID         int        `gorm:"link_id"`
	Token          string     `gorm:"token"`
	SubmittedAt    *time.Time `gorm:"submitted_at"`
	ProgramID      string     `gorm:"program_id"`
	UserID         string     `gorm:"user_id"`
	OrganisationID string     `gorm:"organisation_id"`
	CaregiverID    *string    `gorm:"column:caregiver_id"`
}

// BeforeCreate is a hook run before creating a user survey model
func (u *UserSurvey) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		u.CreatedBy = userID
	}

	id := uuid.New().String()
	u.ID = id

	return
}

// BeforeUpdate is a hook called before updating UserSurvey.
func (u *UserSurvey) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		u.UpdatedBy = userID
	}
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

// BeforeCreate is a hook called before updating Metric.
func (m *Metric) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		m.CreatedBy = userID
	}
	return
}

// BeforeUpdate is a hook called before updating Metric.
func (m *Metric) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		m.UpdatedBy = userID
	}
	return
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
	ProgramID         string `gorm:"column:program_id"`
	OrganisationID    string `gorm:"column:organisation_id"`
	UserID            string `gorm:"column:user_id"`
}

// BeforeCreate is a hook run before creating an appointment
func (f *Feedback) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.CreatedBy = userID
	}

	id := uuid.New().String()
	f.ID = id

	return
}

// BeforeUpdate is a hook called before updating Feedback.
func (f *Feedback) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.UpdatedBy = userID
	}
	return
}

// TableName references the table that we map data from
func (Feedback) TableName() string {
	return "common_feedback"
}

// CommunityClient is represents the relationship between a client and a community. It is basically a through table
type CommunityClient struct {
	ID          int     `gorm:"primaryKey;column:id;autoincrement"`
	CommunityID *string `gorm:"column:community_id"`
	ClientID    *string `gorm:"column:client_id"`
}

// TableName references the table that we map data from
func (CommunityClient) TableName() string {
	return "communities_community_clients"
}

// CommunityStaff is represents the relationship between a staff and a Community. It is basically a through table
type CommunityStaff struct {
	ID          int     `gorm:"primaryKey;column:id;autoincrement"`
	CommunityID *string `gorm:"column:community_id"`
	StaffID     *string `gorm:"column:staff_id"`
}

// TableName references the table that we map data from
func (CommunityStaff) TableName() string {
	return "communities_community_staff"
}

// Questionnaire defines the questionnaire database models
type Questionnaire struct {
	Base
	OrganisationID string `gorm:"column:organisation_id"`

	ID          string `gorm:"primaryKey;column:id"`
	Active      bool   `gorm:"column:active"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
	ProgramID   string `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a questionnaire
func (q *Questionnaire) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		q.CreatedBy = userID
	}
	id := uuid.New().String()
	q.ID = id

	return
}

// BeforeUpdate is a hook called before updating a Questionnaire.
func (q *Questionnaire) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		q.UpdatedBy = userID
	}
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
	ProgramID       string         `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a screening tool
func (s *ScreeningTool) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}
	id := uuid.New().String()
	s.ID = id

	return
}

// BeforeUpdate is a hook called before updating a ScreeningTool.
func (s *ScreeningTool) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}
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
	ProgramID         string `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a question
func (q *Question) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		q.CreatedBy = userID
	}
	id := uuid.New().String()
	q.ID = id

	return
}

// BeforeUpdate is a hook called before updating a Question.
func (q *Question) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		q.UpdatedBy = userID
	}
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
	ProgramID  string `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a question input choice
func (q *QuestionInputChoice) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		q.CreatedBy = userID
	}
	id := uuid.New().String()
	q.ID = id

	return
}

// BeforeUpdate is a hook called before updating a QuestionInputChoice.
func (q *QuestionInputChoice) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		q.UpdatedBy = userID
	}
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

	ID              string  `gorm:"primaryKey;column:id"`
	Active          bool    `gorm:"column:active"`
	ScreeningToolID string  `gorm:"column:screeningtool_id"`
	FacilityID      string  `gorm:"column:facility_id"`
	ClientID        string  `gorm:"column:client_id"`
	AggregateScore  int     `gorm:"column:aggregate_score"`
	ProgramID       string  `gorm:"column:program_id"`
	CaregiverID     *string `gorm:"column:caregiver_id"`
}

// BeforeCreate is a hook run before creating a screening tool response
func (s *ScreeningToolResponse) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}

	id := uuid.New().String()
	s.ID = id

	return
}

// BeforeUpdate is a hook called before updating a ScreeningToolResponse.
func (s *ScreeningToolResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}
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
	ProgramID               string `gorm:"column:program_id"`
	FacilityID              string `gorm:"column:facility_id"`
}

// BeforeCreate is a hook run before creating a screening tool question response
func (s *ScreeningToolQuestionResponse) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.CreatedBy = userID
	}

	id := uuid.New().String()
	s.ID = id

	return
}

// BeforeUpdate is a hook called before updating a screeningtool question response.
func (s *ScreeningToolQuestionResponse) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		s.UpdatedBy = userID
	}
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

// Caregiver is the caregiver profile information for a user
// TODO: remove "N" when original caregiver is removed
type Caregiver struct {
	Base

	ID              string `gorm:"primaryKey;column:id"`
	Active          bool   `gorm:"column:active"`
	CaregiverNumber string `gorm:"column:caregiver_number"`

	OrganisationID  string  `gorm:"column:organisation_id;not null"`
	UserID          string  `gorm:"column:user_id"`
	UserProfile     User    `gorm:"ForeignKey:user_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
	CurrentClient   *string `gorm:"column:current_client"`
	CurrentFacility *string `gorm:"column:current_facility"`
}

// BeforeCreate is a hook run before creating a caregiver
func (c *Caregiver) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}
	id := uuid.New().String()
	c.ID = id

	return nil
}

// BeforeUpdate is a hook called before updating a caregiver.
func (c *Caregiver) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
}

// TableName references the table name in the database
func (c *Caregiver) TableName() string {
	return "caregivers_caregiver"
}

// CaregiverClient stores clients assigned to a caregiver and the associated relationship details
type CaregiverClient struct {
	Base

	CaregiverID        string              `gorm:"column:caregiver_id;not null"`
	ClientID           string              `gorm:"column:client_id;not null"`
	Active             bool                `gorm:"column:active"`
	RelationshipType   enums.CaregiverType `gorm:"column:relationship_type;not null"`
	CaregiverConsent   enums.ConsentState  `gorm:"column:caregiver_consent"`
	CaregiverConsentAt *time.Time          `gorm:"column:caregiver_consent_at"`
	ClientConsent      enums.ConsentState  `gorm:"column:client_consent"`
	ClientConsentAt    *time.Time          `gorm:"column:client_consent_at"`

	OrganisationID string `gorm:"column:organisation_id;not null"`
	AssignedBy     string `gorm:"column:assigned_by;not null"`
	ProgramID      string `gorm:"column:program_id"`
}

// BeforeCreate is a hook run before creating a caregiver client
func (c *CaregiverClient) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.CreatedBy = userID
	}

	return nil
}

// BeforeUpdate is a hook called before updating a caregiver client.
func (c *CaregiverClient) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		c.UpdatedBy = userID
	}
	return
}

// TableName references the table name in the database
func (c *CaregiverClient) TableName() string {
	return "caregivers_caregiver_client"
}

// Program is the database model for a program
type Program struct {
	Base

	ID                 string `gorm:"primaryKey;column:id"`
	Active             bool   `gorm:"column:active"`
	Name               string `gorm:"column:name"`
	Description        string `gorm:"column:description"`
	OrganisationID     string `gorm:"column:organisation_id;not null"`
	FHIROrganisationID string `gorm:"column:fhir_organisation_id"`
}

// BeforeCreate is a hook run before creating a program
func (p *Program) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		p.CreatedBy = userID
	}
	id := uuid.New().String()
	p.ID = id

	return
}

// BeforeUpdate is a hook called before updating a client program.
func (p *Program) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		p.UpdatedBy = userID
	}
	return
}

// TableName references the table name in the database
func (p *Program) TableName() string {
	return "common_program"
}

type AccessToken struct {
	Base

	ID        string `gorm:"primarykey"`
	Active    bool   `gorm:"column:active"`
	Signature string `gorm:"unique;column:signature"`

	RequestedAt       time.Time      `gorm:"column:requested_at"`
	RequestedScopes   pq.StringArray `gorm:"type:varchar(256)[];column:requested_scopes"`
	GrantedScopes     pq.StringArray `gorm:"type:varchar(256)[];column:granted_scopes"`
	Form              pgtype.JSONB   `gorm:"type:jsonb;column:form;default:'{}'"`
	RequestedAudience pq.StringArray `gorm:"type:varchar(256)[];column:requested_audience"`
	GrantedAudience   pq.StringArray `gorm:"type:varchar(256)[];column:granted_audience"`

	ClientID  string `gorm:"column:client_id"`
	Client    OauthClient
	SessionID string `gorm:"column:session_id"`
	Session   Session
}

// TableName references the table name in the database
func (AccessToken) TableName() string {
	return "oauth_access_token"
}

// BeforeCreate is a hook run before creating
func (a *AccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

type AuthorizationCode struct {
	Base

	ID     string `gorm:"primarykey"`
	Active bool   `gorm:"column:active"`
	Code   string `gorm:"column:code"`

	RequestedAt       time.Time      `gorm:"column:requested_at"`
	RequestedScopes   pq.StringArray `gorm:"type:varchar(256)[];column:requested_scopes"`
	GrantedScopes     pq.StringArray `gorm:"type:varchar(256)[];column:granted_scopes"`
	Form              pgtype.JSONB   `gorm:"type:jsonb;column:form;default:'{}'"`
	RequestedAudience pq.StringArray `gorm:"type:varchar(256)[];column:requested_audience"`
	GrantedAudience   pq.StringArray `gorm:"type:varchar(256)[];column:granted_audience"`

	SessionID string `gorm:"column:session_id"`
	Session   Session
	ClientID  string `gorm:"column:client_id"`
	Client    OauthClient
}

// TableName references the table name in the database
func (AuthorizationCode) TableName() string {
	return "oauth_authorization_code"
}

// BeforeCreate is a hook run before creating
func (a *AuthorizationCode) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

type OauthClient struct {
	Base

	ID                      string         `gorm:"primarykey"`
	Name                    string         `gorm:"column:name"`
	Active                  bool           `gorm:"column:active"`
	Secret                  string         `gorm:"column:secret"`
	RotatedSecrets          pq.StringArray `gorm:"type:varchar(256)[];column:rotated_secrets"`
	Public                  bool           `gorm:"column:public"`
	RedirectURIs            pq.StringArray `gorm:"type:varchar(256)[];column:redirect_uris"`
	Scopes                  pq.StringArray `gorm:"type:varchar(256)[];column:scopes"`
	Audience                pq.StringArray `gorm:"type:varchar(256)[];column:audience"`
	Grants                  pq.StringArray `gorm:"type:varchar(256)[];column:grants"`
	ResponseTypes           pq.StringArray `gorm:"type:varchar(256)[];column:response_types"`
	TokenEndpointAuthMethod string         `gorm:"column:token_endpoint_auth_method"`
}

// TableName references the table name in the database
func (OauthClient) TableName() string {
	return "oauth_client"
}

// BeforeCreate is a hook run before creating
func (a *OauthClient) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

type OauthClientJWT struct {
	Base

	ID        string    `gorm:"primarykey"`
	Active    bool      `gorm:"column:active"`
	JTI       string    `gorm:"column:jti"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
}

// TableName references the table name in the database
func (OauthClientJWT) TableName() string {
	return "oauth_client_jwt"
}

// BeforeCreate is a hook run before creating
func (a *OauthClientJWT) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

type PKCE struct {
	Base

	ID        string `gorm:"primarykey"`
	Active    bool   `gorm:"column:active"`
	Signature string `gorm:"unique;column:signature"`

	RequestedAt       time.Time      `gorm:"column:requested_at"`
	RequestedScopes   pq.StringArray `gorm:"type:varchar(256)[];column:requested_scopes"`
	GrantedScopes     pq.StringArray `gorm:"type:varchar(256)[];column:granted_scopes"`
	Form              pgtype.JSONB   `gorm:"type:jsonb;column:form;default:'{}'"`
	RequestedAudience pq.StringArray `gorm:"type:varchar(256)[];column:requested_audience"`
	GrantedAudience   pq.StringArray `gorm:"type:varchar(256)[];column:granted_audience"`

	SessionID string `gorm:"column:session_id"`
	Session   Session
	ClientID  string `gorm:"column:client_id"`
	Client    OauthClient
}

// TableName references the table name in the database
func (PKCE) TableName() string {
	return "oauth_pkce"
}

// BeforeCreate is a hook run before creating
func (a *PKCE) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

type RefreshToken struct {
	Base

	ID        string `gorm:"primarykey"`
	Active    bool   `gorm:"column:active"`
	Signature string `gorm:"unique;column:signature"`

	RequestedAt       time.Time      `gorm:"column:requested_at"`
	RequestedScopes   pq.StringArray `gorm:"type:varchar(256)[];column:requested_scopes"`
	GrantedScopes     pq.StringArray `gorm:"type:varchar(256)[];column:granted_scopes"`
	Form              pgtype.JSONB   `gorm:"type:jsonb;column:form;default:'{}'"`
	RequestedAudience pq.StringArray `gorm:"type:varchar(256)[];column:requested_audience"`
	GrantedAudience   pq.StringArray `gorm:"type:varchar(256)[];column:granted_audience"`

	ClientID  string `gorm:"column:client_id"`
	Client    OauthClient
	SessionID string `gorm:"column:session_id"`
	Session   Session
}

// TableName references the table name in the database
func (RefreshToken) TableName() string {
	return "oauth_refresh_token"
}

// BeforeCreate is a hook run before creating
func (a *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

type Session struct {
	Base

	ID       string `gorm:"primarykey"`
	ClientID string `gorm:"column:client_id"`

	Username  string       `gorm:"column:username"`
	Subject   string       `gorm:"column:subject"`
	ExpiresAt pgtype.JSONB `gorm:"type:jsonb;column:expires_at;default:'{}'"`

	// Default
	Extra pgtype.JSONB `gorm:"type:jsonb;column:extra;default:'{}'"`

	UserID string `gorm:"column:user_id;default:null"`
	User   User
}

// TableName references the table name in the database
func (Session) TableName() string {
	return "oauth_session"
}

// BeforeCreate is a hook run before creating
func (a *Session) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}

	return nil
}

// FacilityCoordinates stores a facilities coordinates
type FacilityCoordinates struct {
	Base

	ID     string  `gorm:"primaryKey;unique;column:id"`
	Active bool    `gorm:"column:active;not null"`
	Lat    float64 `gorm:"column:lat;not null"`
	Lng    float64 `gorm:"column:lng;not null"`

	FacilityID string `gorm:"column:facility_id;not null"`
}

// BeforeCreate is a hook run before creating a new facility coordinates
func (f *FacilityCoordinates) BeforeCreate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.CreatedBy = userID
	}
	id := uuid.New().String()
	f.ID = id

	return
}

// BeforeUpdate is a hook called before updating facility coordinates.
func (f *FacilityCoordinates) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	if userID := utils.GetLoggedInUserID(ctx); userID != nil {
		f.UpdatedBy = userID
	}
	return
}

// TableName customizes how the table name is generated
func (FacilityCoordinates) TableName() string {
	return "facility_coordinates"
}
