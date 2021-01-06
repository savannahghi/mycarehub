# Profiles

### User profile

Query:

```graphql

query{
  userProfile{
    id
    userName
    PrimaryPhone
    verifiedIdentifiers{
      uid      
      loginProvider
      timestamp
    }
    termsAccepted
    suspended
    photoUploadID
    userBioData{
      firstName
      lastName
    }
  }
}

````

Response :

```json


```