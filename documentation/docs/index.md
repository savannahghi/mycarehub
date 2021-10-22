# Profile service

The profile service is responsible of creating and managing users who can and should transact on the Be.Well platform.

As such, it is the entry point for all activities when Be.Well.

## User profile

The `userProfile` is the model that is used represent a user within Be.Well.

Its has a number attributes that are used internal to uniquely and correctly identify a specific user.

Such attributes are as follows;

- PRIMARY PHONE NUMBER - this is the single most important attribute in a user profile. Every user proile must have this

- PRIMARY EMAIL ADDRESS - this comes second in importance when seeking to uniquely identify a user, It is not available by defauly, though it can be provided via social login

- SECONDARY PHONE NUMBERS - a user can have multiple phone numbers. Provided they belong and we have prove they belong to the user, they will added under this attribute.

- SECONDARY EMAIL ADDRESS - a user can have multiple email addresses. Once proff of ownership has been established, they will added under this attribute.

- USERNAME - this a name that the user use whill transacting on Be.Well. Since it's unique by nature, it should be treated a strong identification factor

- VERIFIED IDENTIFIERS - these are system generated tokens of identification used internally to identify a specific user.

The structure of the `userProfile` model is like below

```go

// UserProfile serializes the profile of the logged in user.
type UserProfile struct {
	// globally unique identifier for a profile
	ID string `json:"id" firestore:"id"`

	// unique user name. Synonymous to a handle
	// e.g @juliusowino
	// this will be auto-generated on first login, meaning a user must have a username
	Username string `json:"userName" firestore:"userName"`

	// VerifiedIdentifiers represent various ways the user has been able to login
	// and these providers point to the same user
	VerifiedIdentifiers []VerifiedIdentifier `json:"verifiedIdentifiers" firestore:"verifiedIdentifiers"`

	// uids associated with a profile. Theses UIDS should match those in the verfiedIdentifiers.
	// the purpose of having verifiedUIDS is enbale ease querying of the profile using firebase query constructs.
	// when we migrate to postgres, this will be retired
	// the length of verfiedIdentifiers and verifiedUIDS should match
	VerifiedUIDS []string `json:"verifiedUIDS" firestore:"verifiedUIDS"`

	// this is the first class unique attribute of a user profile.  A user profile MUST HAVE A PRIMARY PHONE NUMBER
	PrimaryPhone string `json:"primaryPhone" firestore:"primaryPhone"`

	// this is the second class unique attribute of a user profile. This can be updated as the user desires
	PrimaryEmailAddress string `json:"primaryEmailAddress" firestore:"primaryEmailAddress"`

	// these are all phone numbers associated with a user. These phone numbers can be promoted to PRIMARY PHONE NUMBER
	// and/or used for account recovery
	SecondaryPhoneNumbers []string `json:"secondaryPhoneNumbers" firestore:"secondaryPhoneNumbers"`

	SecondaryEmailAddresses []string `json:"secondaryEmailAddresses " firestore:"secondaryEmailAddresses"`

	PushTokens []string `json:"pushTokens,omitempty" firestore:"pushTokens"`

	// what the user is allowed to do. Only valid for admins
	Permissions []PermissionType `json:"permissions,omitempty" firestore:"permissions"`

	// we determine if a user is "live" by examining fields on their profile
	TermsAccepted bool `json:"terms_accepted,omitempty" firestore:"termsAccepted"`

	// determines whether a specific will be visible in query results. If the `true`, means the profile in not
	// in active state and the user should not be allowed to login
	Suspended bool `json:"suspended" firestore:"suspended"`

	// a user's profile photo can be stored as base 64 encoded PNG
	PhotoUploadID string `json:"photoUploadID,omitempty" firestore:"photoUploadID"`

	// a user can have zero or more insurance covers
	Covers []Cover `json:"covers,omitempty" firestore:"covers"`

	// a user's biodata is stored on the profile
	UserBioData BioData `json:"userBioData,omitempty" firestore:"userBioData"`
}
```

## Validity of a user profile

Since there are a finite number of pre-defined MUST have attributes in a `userProfile`, any `userProfile` validator must assert the following;

- A valid `userProfile` have user provided PRIMARY PHONE NUMBER.

- A valid `userProfile` have a generated or user USERNAME

- A valid `userProfile` have at least one VERIFIED IDENTIFIER

- A valid `userProfile` have at least one PUSH TOKEN

- A valid `userProfile` have at a FIRSTNAME and LASTNAME in its bio data information

## Endpoints

- Staging

  https://profile-staging.healthcloud.co.ke

- Testing

  https://profile-testing.healthcloud.co.ke

- Production

  https://profile-prod.healthcloud.co.ke
