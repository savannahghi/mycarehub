# Creating an account

The following APIs correspond to the following agreed upon workflows;

- [Consumer](https://drive.google.com/file/d/1zRjZse_dyyjVhqVw48BnLRt7U1ovgd0D/view?usp=sharing)

- [Pro](https://drive.google.com/file/d/1LnvGj_jOFMgpi0HNRw8A_zspjH0-SamQ/view?usp=sharing)

## APIs

### Check if phone number exists

!!! note
    {BASE_URL}/verify_phone"

This first step to creating an account, is to creating whether a the provided phone number has been registered to another user or not.

That is what this API does. If the phone number is unique (unknown), it will send an OTP to the provided phone number

It check both PRIMARY PHONE NUMBERS and SECONDARY PHONE NUMBERS

Example (phone number exists);

```sh
http https://profile-testing.healthcloud.co.ke/verify_phone phoneNumber="+254718376163"


HTTP/1.1 400 Bad Request
Cache-Control: private
Content-Encoding: gzip
Content-Type: text/html
Date: Tue, 05 Jan 2021 12:50:17 GMT
Server: Google Frontend
Transfer-Encoding: chunked
X-Cloud-Trace-Context: e0f23bcbd33ad69e1a87f66b286cd91e;o=1
vary: Accept-Encoding

{
    "code": 4,
    "message": "provided phone number is already in use"
}
```

Example (phone number does not exist)

```sh
http https://profile-testing.healthcloud.co.ke/verify_phone phoneNumber="+254715825862"

HTTP/1.1 200 OK
Cache-Control: private
Content-Encoding: gzip
Content-Type: text/html
Date: Tue, 05 Jan 2021 12:55:45 GMT
Server: Google Frontend
Transfer-Encoding: chunked
X-Cloud-Trace-Context: eef2c7afb07866afbcebaff0342c0b51
vary: Accept-Encoding

{
    "otp": "{\"otp\":\"779124\"}"
}
```

!!! note
	Go to `Error codes` page for details regarding error codes.



### Create account

!!! note
    {BASE_URL}/create_user_by_phone


This should be called after verifying the phone and receiving an OTP

Payload:

```json

{
    "phoneNumber": "+254712345678",
    "pin": "1234",
    "flavour" : "CONSUMER", // PRO
}
```

Example;

```sh
http https://profile-testing.healthcloud.co.ke/create_user_by_phone phoneNumber="+254712345678" pin="1234" flavour="CONSUMER"


HTTP/1.1 201 Created
Cache-Control: private
Content-Encoding: gzip
Content-Type: text/html
Date: Wed, 06 Jan 2021 10:39:33 GMT
Server: Google Frontend
Transfer-Encoding: chunked
X-Cloud-Trace-Context: 225754d0ed47e549ba0ec3916e649034;o=1
vary: Accept-Encoding

{
    "auth": {
        "customToken": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiI4NDE5NDc3NTQ4NDctY29tcHV0ZUBkZXZlbG9wZXIuZ3NlcnZpY2VhY2NvdW50LmNvbSIsImF1ZCI6Imh0dHBzOi8vaWRlbnRpdHl0b29sa2l0Lmdvb2dsZWFwaXMuY29tL2dvb2dsZS5pZGVudGl0eS5pZGVudGl0eXRvb2xraXQudjEuSWRlbnRpdHlUb29sa2l0IiwiZXhwIjoxNjA5OTMzMTcxLCJpYXQiOjE2MDk5Mjk1NzEsInN1YiI6Ijg0MTk0Nzc1NDg0Ny1jb21wdXRlQGRldmVsb3Blci5nc2VydmljZWFjY291bnQuY29tIiwidWlkIjoieHBvRE55cDJpSFhuWUFsMHJuT09Id2E1R3lzMSJ9.Fvehk46bwe2xZWrFVZZJHsmIzENjTrvoXT3N3dQMkYigkefZEGxC-I4r5TvUKbeyNUcMIEsfVuw_T_pOsbvpnAfwWInRIWaXF2pT8LJ22Uau7XgLTAOxXn7N5sfYUhlqzXZs86A293wJpieuTlWK3ymTfqNFPrRfnyzLartZJeW9ieD2i6HeWIqrgWDIf1w9hDP_cEqkC1ngoVvAaT0li8qVQASB8meTbqfhzS_tsgcRWn2pJc8xM_DYzfyqiacIbJOqGRV0ighAOPOUZyDsTv06EaKoQ4XHd0MF8i6T815ZRIjARXE5ompu4w1aLK-hIRIgsAcInbYHv8II-Eu6_A",
        "expires_in": "3600",
        "id_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImUwOGI0NzM0YjYxNmE0MWFhZmE5MmNlZTVjYzg3Yjc2MmRmNjRmYTIiLCJ0eXAiOiJKV1QifQ.eyJpc3MiOiJodHRwczovL3NlY3VyZXRva2VuLmdvb2dsZS5jb20vYmV3ZWxsLWFwcCIsImF1ZCI6ImJld2VsbC1hcHAiLCJhdXRoX3RpbWUiOjE2MDk5Mjk1NzIsInVzZXJfaWQiOiJ4cG9ETnlwMmlIWG5ZQWwwcm5PT0h3YTVHeXMxIiwic3ViIjoieHBvRE55cDJpSFhuWUFsMHJuT09Id2E1R3lzMSIsImlhdCI6MTYwOTkyOTU3MiwiZXhwIjoxNjA5OTMzMTcyLCJwaG9uZV9udW1iZXIiOiIrMjU0NzEyMzQ1Njc4IiwiZmlyZWJhc2UiOnsiaWRlbnRpdGllcyI6eyJwaG9uZSI6WyIrMjU0NzEyMzQ1Njc4Il19LCJzaWduX2luX3Byb3ZpZGVyIjoiY3VzdG9tIn19.eSkV1wke9pxJn4GJducehUubtQ18UYkwsXqCS8mrFDccwObI-l6tBpZT3-vHiECqhhXgggf4e8redjznOHKti_bU6bhnx7egVxqRpja47pZVQlFdKXM56C09gA4tJgk5xnpxRk3QcgrhOJM9TJXJFjLhqxYPKURCukNkqFfBfXeWBGja8gtepUENwyEwQb-FIT8iCzW6zUW3zt553xkKW2R358WIa7eEG97bx5MmrDABOnTP6KiXDwCPs83tiRUI8EpHfREjZ5HtR9_3b4W92FnCpMT3zg41SFoBUo1bej5ZEfPPHLdUuMBN1Cj32zB61UGfacyvToHfR-Wl5gfe3g",
        "is_admin": false,
        "is_anonymous": false,
        "refresh_token": "AG8BCne3t20pAaaMnqMywyhtrMo53D34IUPpsdu-njTYUfktP2xLsNaSKDjvQY5Nd3D7XroMAUdYKnQN76pJpkG7wXCTim8G7AXUOzdq2KRhwIJ6KJLu8oNejindfLKkxBCXb9Vfl8PupPV6oTtMkaYUvnFmCzQAiTRAJPpOf3gtqkLjJBOT4FM",
        "uid": "xpoDNyp2iHXnYAl0rnOOHwa5Gys1"
    },
    "customerProfile": {
        "active": false,
        "id": "d0ead8a8-d9d2-4f16-a2d7-3e0f5c32bf51",
        "receivablesAccount": {
            "description": "",
            "id": "",
            "isActive": false,
            "name": "",
            "number": "",
            "tag": ""
        }
    },
    "profile": {
        "id": "b8124bc0-7321-45cb-a435-63a805c57323",
        "primaryEmailAddress": "",
        "primaryPhone": "+254712345678",
        "secondaryEmailAddresses ": null,
        "secondaryPhoneNumbers": null,
        "suspended": false,
        "terms_accepted": true,
        "userBioData": {
            "dateOfBirth": null,
            "firstName": "",
            "gender": "",
            "lastName": ""
        },
        "userName": "@stoic_euler17167254",
        "verifiedIdentifiers": [
            {
                "loginProvider": "PHONE",
                "timeStamp": "2021-01-06T10:39:31.628298Z",
                "uid": "xpoDNyp2iHXnYAl0rnOOHwa5Gys1"
            }
        ],
        "verifiedUIDS": [
            "xpoDNyp2iHXnYAl0rnOOHwa5Gys1"
        ]
    },
    "supplierProfile": {
        "accountType": "",
        "active": false,
        "ediuserprofile": null,
        "id": "17815826-3678-43e5-875f-cebb8a082cdd",
        "isOrganizationVerified": false,
        "kycSubmitted": false,
        "parentOrganizationID": "",
        "partnerSetupComplete": false,
        "partnerType": "",
        "payablesAccount": null,
        "profileID": "b8124bc0-7321-45cb-a435-63a805c57323",
        "sladeCode": "",
        "supplierID": "",
        "supplierKYC": null,
        "supplierName": "",
        "underOrganization": false
    }
}

```


### Update Bio data


Input:

```json

{
  "input": {
    "firstName": "ule",
    "lastName": "makmende"
  }
}
```

Mutation: 

```graphql

mutation updateProfile($input:UserProfileInput!){
  updateUserProfile(input:$input){
    id
    userBioData{
      firstName
      lastName
    }
  }
}
```