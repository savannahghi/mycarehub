# Onboarding service

This service manages user onboarding process. 

## Description

The project implements the `Clean Architecture` advocated by 
Robert Martin ('Uncle Bob').

### Clean Architecture

A cleanly architected project should be:

- *Independent of Frameworks*: The architecture does not depend on the 
existence of some library of feature laden software. This allows you to use 
such frameworks as tools, rather than having to cram your system into their 
limited constraints.

- *Testable*: The business rules can be tested without the UI, Database, 
Web Server, or any other external element.

- *Independent of UI*: The UI can change easily, without changing the rest of
the system. A Web UI could be replaced with a console UI, for example,
without changing the business rules.

- *Independent of Database*: You can swap out Cloud Firestore or SQL Server,
for Mongo, Postgres, MySQL, or something else. Your business rules are not
bound to the database.

- *Independent of any external agency*: In fact your business rules simply
don’t know anything at all about the outside world.

## This project has 5 layers:

### Domain Layer

Here we have `business objects` or `entities` and should represent and 
encapsulate the fundamental business rules.

### Repository Layer

In the domain layer we should have no idea about any database nor any storage, 
so the repository is just an interface.

### Infrastructure Layer

These are the `ports` that allow the system to talk to 'outside things' which 
could be a `database` for persistence or a `web server` for the UI. None of 
the inner use cases or domain entities should know about the implementation of 
these layers and they may change over time because ... well, we used to store 
data in SQL, then document database and changing the storage should not change
the application or any of the business rules.

### Usecase Layer

The code in this layer contains application specific business rules. It 
encapsulates and implements all of the use cases of the system. These use cases
orchestrate the flow of data to and from the entities, and direct those
entities to use their enterprise wide business rules to achieve the goals of
the use case.

This represents the pure business logic of the application.
The rules of the application also shouldn't rely on the UI or the persistence
frameworks being used.

### Presentation Layer

This represents logic that consume the business logic from the `Usecase Layer`
and renders to the view. Here you can choose to render the view in e.g 
`graphql` or `rest`

For more information, see:

- [The Clean Architecture](https://blog.8thlight.com/uncle-bob/2012/08/13/the-clean-architecture.html) advocated by Robert Martin ('Uncle Bob')
- Ports & Adapters or [Hexagonal Architecture](http://alistair.cockburn.us/Hexagonal+architecture) by Alistair Cockburn
- [Onion Architecture](http://jeffreypalermo.com/blog/the-onion-architecture-part-1/) by Jeffrey Palermo
- [Implementing Domain-Driven Design](http://www.amazon.com/Implementing-Domain-Driven-Design-Vaughn-Vernon/dp/0321834577)

## Environment variables

For local development, you need to _export_ the following env vars:

```bash
# Google Cloud Settings
export GOOGLE_APPLICATION_CREDENTIALS="<a path to a Google service account JSON file>"
export GOOGLE_CLOUD_PROJECT="<the name of the project that the service account above belongs to>"
export FIREBASE_WEB_API_KEY="<an API key from the Firebase console for the project mentioned above>"

# Go private modules
export GOPRIVATE="gitlab.slade360emr.com/go/*,gitlab.slade360emr.com/optimalhealth/*"

# Charge Master API settings
export CHARGE_MASTER_API_HOST="<a charge master API host>"
export CHARGE_MASTER_API_SCHEME=https
export CHARGE_MASTER_TOKEN_URL="<an auth server token URL>"
export CHARGE_MASTER_CLIENT_ID="<an auth server client ID>"
export CHARGE_MASTER_CLIENT_SECRET="<an auth server client secret>"
export CHARGE_MASTER_USERNAME="<an auth server username>"
export CHARGE_MASTER_PASSWORD="<an auth server password>"
export CHARGE_MASTER_GRANT_TYPE="<an auth server grant type>"

```

The server deploys to Google Cloud Run. The environment variables defined above
should also be set on Google Cloud Run.

---------------------------------------------------

## ISC EndPoints

### `/user_profile`

 Expects a post request with a payload of context with format,

 ```golang
 type userContext struct {
     token *auth.Token `json:"token"`
 }
 ```

 token is a firebase auth Token.  
 Returns a json response with the users profile.  

```golang
    type UserProfile struct {
        UID              string           `json:"uid" firestore:"uid"`
        TermsAccepted    bool             `json:"termsAccepted" firestore:"termsAccepted"`
        IsApproved       bool             `json:"isApproved" firestore:"isApproved"`
        Msisdns          []string         `json:"msisdns" firestore:"msisdns"`
        Emails           []string         `json:"emails" firestore:"emails"`
        PhotoBase64      string           `json:"photoBase64" firestore:"photoBase64"`
        PhotoContentType base.ContentType `json:"photoContentType" firestore:"photoContentType"`
        Covers           []Cover          `json:"covers" firestore:"covers"`

        DateOfBirth *base.Date   `json:"dateOfBirth,omitempty" firestore:"dateOfBirth,omitempty"`
        Gender      *base.Gender `json:"gender,omitempty" firestore:"gender,omitempty"`
        PatientID   *string      `json:"patientID,omitempty" firestore:"patientID"`
        PushTokens  []string     `json:"pushTokens" firestore:"pushTokens"`

        Name                               *string `json:"name" firestore:"name"`
        Bio                                *string `json:"bio" firestore:"bio"`
        PractitionerApproved               *bool   `json:"practitionerApproved" firestore:"practitionerApproved"`
        PractitionerTermsOfServiceAccepted *bool   `json:"practitionerTermsOfServiceAccepted" firestore:"practitionerTermsOfServiceAccepted"`

        IsTester      bool          `json:"isTester" firestore:"isTester"`
        CanExperiment bool          `json:"canExperiment" firestore:"canExperiment"`
        Language      base.Language `json:"language" firestore:"language"`

        // used to determine whether to persist asking the user on the UI
        AskAgainToSetIsTester      bool             `json:"askAgainToSetIsTester" firestore:"askAgainToSetIsTester"`
        AskAgainToSetCanExperiment bool             `json:"askAgainToSetCanExperiment" firestore:"askAgainToSetCanExperiment"`
        VerifiedEmails             []VerifiedEmail  `json:"verifiedEmails" firestore:"verifiedEmails"`
        VerifiedPhones             []VerifiedMsisdn `json:"verifiedPhones" firestore:"verifiedPhones"`
        HasPin                     bool             `json:"hasPin" firestore:"hasPin"`
        HasSupplierAccount         bool             `json:"hasSupplierAccount" firestore:"hasSupplierAccount"`
        HasCustomerAccount         bool             `json:"hasCustomerAccount" firestore:"hasCustomerAccount"`
        PractitionerHasServices    bool             `json:"practitionerHasServices" firestore:"practitionerHasServices"`
}
```


### `customer`

requestPayload

```json
    {
        "uid": string,
        "token": auth.Token // firebase token
    }
```

`profileClient.MakeRequest("internal/customer", "POST", requestPayload)`  
responsePayload

```json
    {
        "customer_id": string,
        "receivables_account": {
            "id": string,
            "name": string,
            "is_active": bool,
            "number": string,
            "tag": string,
            "description" string
        },
        "profile": {
            "uid": string,
            "msisdns": [string],
            "emails": [emails],
            "gender": "male|female|other|unknown",
            "pushTokens": [string],
            "name": string,
            "bio": string
        },
        "customer_kyc": {
            "kra_pin": string,
            "occupation": string,
            "id_number": string,
            "address": string,
            "city": string,
            "beneficiary": [
                {
                    "name": string,
                    "msisdns": [string],
                    "emails": [string],
                    "relationship": "SPOUSE|CHILD",
                    "dateOfBirth": {
                        "Year": int,
                        "Month": int,
                        "Day": int
                    }
                }
            ]
        }
    }
```

### `/supplier`

requestPayload

```json
    {
        "uid": string,
        "token": auth.Token // firebase token
    }
```

responsePayload

```json
    {
        "supplier_id": string,
        "payables_account": {
            "id": string,
            "name": string,
            "is_active": bool,
            "number": string,
            "tag": string,
            "description": string
        },
        "profile": {
            "uid": string,
            "msisdns": [string],
            "emails": [emails],
            "gender": "male|female|other|unknown",
            "pushTokens": [string],
            "name": string,
            "bio": string
        }
    }
```

### `/contactdetails/{attribute}/`

`profileClient.MakeRequest("internal/contactdetails/112", "POST", requestPayload)`

requestPayload is

```json
    {
        "uids": [string]
    }
```

responsePayload is `map[string][]string`

```json
    {

    }
```

### `/retrieve_user_profile`

requestPayload

```json
    {
        "uid": string
    }
```

requestResponse `UserProfile`

### `/save_cover`

requestPayload

```json
    {
        "payerName": string,
        "memberName": string,
        "memberNumber": string,
        "payerSladeCode": int,
        "uid": string,
        "token": auth.Token
    }
```

responsePayload

```json
    {
        "successfullySaved": bool
    }
```

step 1 run 
go run github.com/99designs/gqlgen init
go generate ./...