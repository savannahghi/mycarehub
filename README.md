# Profile micro-service

[![pipeline status](https://gitlab.slade360emr.com/go/profile/badges/develop/pipeline.svg)](https://gitlab.slade360emr.com/go/profile/-/commits/develop)
[![coverage report](https://gitlab.slade360emr.com/go/profile/badges/develop/coverage.svg)](https://gitlab.slade360emr.com/go/profile/-/commits/develop)

## Environment variables

For local development, you need to _export_ the following env vars:

```bash
# Google Cloud Settings
export GOOGLE_APPLICATION_CREDENTIALS="<a path to a Google service account JSON file>"
export GOOGLE_CLOUD_PROJECT="<the name of the project that the service account above belongs to>"
export FIREBASE_WEB_API_KEY="<an API key from the Firebase console for the project mentioned above>"

# Mailgun settings
export MAILGUN_API_KEY=key="<an API key>"
export MAILGUN_DOMAIN=app.healthcloud.co.ke
export MAILGUN_FROM=hello@app.healthcloud.co.ke

# AfricasTalking SMS API settings
export AIT_API_KEY="<an API key>"
export AIT_USERNAME=sandbox
export AIT_SENDER_ID=HealthCloud
export AIT_ENVIRONMENT=sandbox

# Go private modules
export GOPRIVATE="gitlab.slade360emr.com/go/*,gitlab.slade360emr.com/optimalhealth/*"
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

### `/is_underage`  

Process post request with a payload with a context

```golang
    type UserContext struct {
        Token *auth.Token `json:"token"`
    }
```

returns a response with json

```golang
    type Payload struct {
        IsUnderAge bool `json:"isUnderAge"`
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