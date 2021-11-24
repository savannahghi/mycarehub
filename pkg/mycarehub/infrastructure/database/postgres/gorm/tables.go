package gorm

import (
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

// OrganizationID assign a default organisation to a type
var OrganizationID = serverutils.MustGetEnvVar(common.OrganizationID)

// Base model contains defines commin fields across tables
type Base struct {
	CreatedAt time.Time `gorm:"column:created"`
	UpdatedAt time.Time `gorm:"column:updated"`
	// OrganisationID string    `gorm:"column:organisation_id"`
	//DeletedAt      time.Time `gorm:"column:deleted_at"`
}

// Facility models the details of healthcare facilities that are on the platform.
//
// e.g CCC clinics, Pharmacies.
type Facility struct {
	Base
	//globally unique when set
	FacilityID *string `gorm:"primaryKey;unique;column:id"`
	// unique within this structure
	Name string `gorm:"column:name;unique;not null"`
	// MFL Code for Kenyan facilities, globally unique
	Code           int    `gorm:"unique;column:mfl_code;not null"`
	Active         bool   `gorm:"column:active;not null"`
	County         string `gorm:"column:county;not null"` // TODO: Controlled list of counties
	Description    string `gorm:"column:description;not null"`
	OrganisationID string `gorm:"column:organisation_id"`
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
	return common.FacilityTableName
}

// User represents the table data structure for a user
type User struct {
	// Base

	UserID *string `gorm:"primaryKey;unique;column:id"` // globally unique ID

	Username string `gorm:"column:username;unique;not null"` // @handle, also globally unique; nickname

	// DisplayName string `gorm:"column:display_name"` // user's preferred display name

	// TODO Consider making the names optional in DB; validation in frontends
	FirstName  string `gorm:"column:first_name;not null"` // given name
	MiddleName string `gorm:"column:middle_name"`
	LastName   string `gorm:"column:last_name;not null"`

	UserType enums.UsersType `gorm:"column:user_type;not null"` // TODO enum; e.g client, health care worker

	Gender enumutils.Gender `gorm:"column:gender;not null"` // TODO enum; genders; keep it simple

	Active bool `gorm:"column:is_active;not null"`

	Contacts Contact `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"` // TODO: validate, ensure

	// // for the preferred language list, order matters
	// Languages pq.StringArray `gorm:"type:text[];column:languages;not null"` // TODO: turn this into a slice of enums, start small (en, sw)

	PushTokens []string `gorm:"type:text[];column:push_tokens"`

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
	AcceptedTermsID        *int            `gorm:"column:accepted_terms_of_service_id"` // foreign key to version of terms they accepted
	Flavour                feedlib.Flavour `gorm:"column:flavour;not null"`
	Avatar                 string          `gorm:"column:avatar"`
	IsSuspended            bool            `gorm:"column:is_suspended;not null"`
	PinChangeRequired      bool            `gorm:"column:pin_change_required"`
	HasSetPin              bool            `gorm:"column:has_set_pin"`
	HasSetSecurityQuestion bool            `gorm:"column:has_set_security_questions"`
	IsPhoneVerified        bool            `gorm:"column:is_phone_verified"`

	// Django required fields
	OrganisationID   string `gorm:"column:organisation_id"`
	Password         string `gorm:"column:password"`
	IsSuperuser      bool   `gorm:"column:is_superuser"`
	IsStaff          bool   `gorm:"column:is_staff"`
	Email            string `gorm:"column:email"`
	DateJoined       string `gorm:"column:date_joined"`
	Name             string `gorm:"column:name"`
	IsApproved       bool   `gorm:"column:is_approved"`
	ApprovalNotified bool   `gorm:"column:approval_notified"`
	Handle           string `gorm:"column:handle"`
}

// BeforeCreate is a hook run before creating a new user
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	u.UserID = &id
	u.OrganisationID = OrganizationID
	salt, _ := extension.NewExternalMethodsImpl().EncryptPIN(ksuid.New().String(), nil)
	bytePass := []byte(salt)
	u.Password = string(bytePass[0:127])
	u.IsSuperuser = false
	u.IsStaff = false
	u.Email = gofakeit.Email()
	u.DateJoined = time.Now().UTC().Format(time.RFC1123Z)
	u.Name = gofakeit.Name()
	u.IsApproved = false
	u.ApprovalNotified = false
	u.Handle = "@" + u.Username

	return
}

// TableName customizes how the table name is generated
func (User) TableName() string {
	return "users_user"
}

// Contact hold contact information/details for users
type Contact struct {
	Base

	ContactID *string `gorm:"primaryKey;unique;column:id"`

	ContactType string `gorm:"column:contact_type;not null"` // TODO enum

	ContactValue string `gorm:"unique;column:contact_value;not null"` // TODO Validate: phones are E164, emails are valid

	Active bool `gorm:"column:active;not null"`

	// a user may opt not to be contacted via this contact
	// e.g if it's a shared phone owned by a teenager
	OptedIn bool `gorm:"column:opted_in;not null"`

	UserID *string `gorm:"column:user_id;not null"`

	Flavour feedlib.Flavour `gorm:"column:flavour"`

	OrganisationID string `gorm:"column:organisation_id"`
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
	UserID    string          `gorm:"column:user_id;not null"`
	HashedPIN string          `gorm:"column:hashed_pin;not null"`
	ValidFrom time.Time       `gorm:"column:valid_from;not null"`
	ValidTo   time.Time       `gorm:"column:valid_to;not null"`
	IsValid   bool            `gorm:"column:active;not null"`
	Flavour   feedlib.Flavour `gorm:"column:flavour;not null"`
	Salt      string          `gorm:"column:salt;not null"`
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

// TableName customizes how the table name is generated
func (TermsOfService) TableName() string {
	return "users_termsofservice"
}

// Organisation maps the organization table. This will be useful when running integration
// tests since many models have an organization ID as a foreign key.
type Organisation struct {
	Base

	OrganisationID   *string `gorm:"primaryKey;unique;column:id"`
	Active           bool    `gorm:"column:active;not null"`
	OrgCode          string  `gorm:"column:org_code"`
	Code             int     `gorm:"column:code"`
	OrganisationName string  `gorm:"column:organisation_name"`
	EmailAddress     string  `gorm:"column:email_address"`
	PhoneNumber      string  `gorm:"column:phone_number"`
	PostalAddress    string  `gorm:"column:postal_address"`
	PhysicalAddress  string  `gorm:"column:physical_address"`
	DefaultCountry   string  `gorm:"column:default_country"`
}

// BeforeCreate is a hook run before creating a new organisation
func (t *Organisation) BeforeCreate(tx *gorm.DB) (err error) {
	t.OrganisationID = &OrganizationID
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
	Description        string                             `gorm:"column:description"`   // help text
	ResponseType       enums.SecurityQuestionResponseType `gorm:"column:response_type"` // TODO: Enum
	Flavour            feedlib.Flavour                    `gorm:"column:flavour"`       // TODO: Enum
	Active             bool                               `gorm:"column:active"`
	Sequence           *int                               `gorm:"column:sequence"` // for sorting

	// Django reqired fields
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
	UserID      string          `gorm:"column:user_id"`
	Valid       bool            `gorm:"column:is_valid"`
	GeneratedAt time.Time       `gorm:"column:generated_at"`
	ValidUntil  time.Time       `gorm:"column:valid_until"`
	Channel     string          `gorm:"column:channel"`
	Flavour     feedlib.Flavour `gorm:"column:flavour"`
	PhoneNumber string          `gorm:"column:phonenumber"`
	OTP         string          `gorm:"column:otp"`
}

// TableName customizes how the table name is generated
func (UserOTP) TableName() string {
	return "users_userotp"
}

// SecurityQuestionResponse maps the schema for the table that stores the security question
// responses
type SecurityQuestionResponse struct {
	Base

	ResponseID     string    `gorm:"column:id"`
	QuestionID     string    `gorm:"column:question_id"`
	UserID         string    `gorm:"column:user_id"`
	Active         bool      `gorm:"column:active"`
	Response       string    `gorm:"column:response"`
	OrganisationID string    `gorm:"column:organisation_id"`
	Timestamp      time.Time `gorm:"column:timestamp"`
	IsCorrect      bool      `gorm:"column:is_correct"`
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

	ClientType string `gorm:"column:client_type"`

	UserProfile User `gorm:"ForeignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`

	TreatmentEnrollmentDate *time.Time `gorm:"column:enrollment_date"`

	FHIRPatientID string `gorm:"column:fhir_patient_id"`

	HealthRecordID *string `gorm:"column:emr_health_record_id"`

	TreatmentBuddy string `gorm:"column:treatment_buddy"` // TODO: optional, free text OR FK to user?

	ClientCounselled bool `gorm:"column:counselled"`

	OrganisationID string `gorm:"column:organisation_id"`

	FacilityID string `gorm:"column:current_facility_id"`

	CHVUserID string `gorm:"column:chv_id"`

	UserID *string `gorm:"column:user_id;not null"`
}

// BeforeCreate is a hook run before creating a client profile
func (c *Client) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (Client) TableName() string {
	return "clients_client"
}

// ContentItemCategory maps the schema for the table that stores the content item category
type ContentItemCategory struct {
	ID   int          `gorm:"unique;column:id;autoincrement"`
	Name string       `gorm:"column:name"`
	Icon WagtailImage `gorm:"ForeignKey:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;not null"`
}

// TableName customizes how the table name is generated
func (ContentItemCategory) TableName() string {
	return "content_contentitemcategory"
}

// WagtailImage maps the schema for the table that stores the wagtail images
type WagtailImage struct {
	ID   int    `gorm:"primaryKey;column:id;"`
	File string `gorm:"column:file"`
}

// TableName customizes how the table name is generated
func (WagtailImage) TableName() string {
	return "wagtailimages_image"
}

// ContentAuthor is the gorms content author model
type ContentAuthor struct {
	Base
	ContentAuthorID *string `gorm:"column:id"`
	Active          bool    `gorm:"column:active"`
	Name            string  `gorm:"column:name"`
	OrganisationID  string  `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating an author
func (c *ContentAuthor) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ContentAuthorID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ContentAuthor) TableName() string {
	return "content_author"
}

// ContentItem is the gorms content item model
type ContentItem struct {
	PagePtrID           int       `gorm:"column:page_ptr_id"`
	Date                time.Time `gorm:"column:date"`
	Intro               string    `gorm:"column:intro"`
	ItemType            string    `gorm:"column:item_type"`
	TimeEstimateSeconds int       `gorm:"column:time_estimate_seconds"`
	Body                string    `gorm:"column:body"`
	LikeCount           int       `gorm:"column:like_count"`
	BookmarkCount       int       `gorm:"column:bookmark_count"`
	ShareCount          int       `gorm:"column:share_count"`
	ViewCount           int       `gorm:"column:view_count"`
	AuthorID            string    `gorm:"column:author_id"`
	HeroImageID         *string   `gorm:"column:hero_image_id"`
}

// TableName references the table that we map data from
func (ContentItem) TableName() string {
	return "content_contentitem"
}

// ContentShare is the gorms content contentshare model
type ContentShare struct {
	Base
	ContentShareID *string `gorm:"column:id"`
	Active         bool    `gorm:"column:active"`
	ContentID      int     `gorm:"column:content_item_id"`
	UserID         string  `gorm:"column:user_id"`
	OrganisationID string  `gorm:"column:organisation_id"`
}

// BeforeCreate is a hook run before creating count share
func (c *ContentShare) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.New().String()
	c.ContentShareID = &id
	c.OrganisationID = OrganizationID
	return
}

// TableName references the table that we map data from
func (ContentShare) TableName() string {
	return "public.content_contentshare"
}

// WagtailCorePage models the details of core wagtail fields
type WagtailCorePage struct {
	WagtailCorePageID     int    `gorm:"unique;column:id;autoincrement"`
	Path                  string `gorm:"column:path"`
	Depth                 int    `gorm:"column:depth"`
	Numchild              int    `gorm:"column:numchild"`
	Title                 string `gorm:"column:title"`
	Slug                  string `gorm:"column:slug"`
	Live                  bool   `gorm:"column:live"`
	HasUnpublishedChanges bool   `gorm:"column:has_unpublished_changes"`
	URLPath               string `gorm:"column:url_path"`
	SEOTitle              string `gorm:"column:seo_title"`
	ShowInMenus           bool   `gorm:"column:show_in_menus"`
	SearchDescription     string `gorm:"column:search_description"`
	Expired               bool   `gorm:"column:expired"`
	ContentTypeID         int    `gorm:"column:content_type_id"` // default to 1 => wagtailcore page
	Locked                bool   `gorm:"column:locked"`
	DraftTitle            string `gorm:"column:draft_title"`
	TranslationKey        string `gorm:"column:translation_key"`
	LocaleID              int    `gorm:"column:locale_id"` // default to 1 => en
}

// TableName references the table that we map data from
func (WagtailCorePage) TableName() string {
	return "wagtailcore_page"
}
