# Service Request API specification

This file gives an overview of the Service Request API(s) required by the application to perform service request
activities.

Service Request is notification sent to the staff when a user faces an issue or when the application detects an issue.
It might be a red flag, health diary entry, profile update, contact request, pin reset request

The API(s) schemas are defined in `GraphQL`.

## EndPoint Definitions

Base URLs:
- https://mycarehub-testing.savannahghi.org/ide

- https://mycarehub-staging.savannahghi.org/ide

- https://mycarehub-prod.savannahghi.org/ide

### Schema Inputs


### Schema Types
```
type ServiceRequest{
  ID: String!
  RequestType: String!
  Request: String!
  Status: String! 
  ClientID: String!
  InProgressAt: Time!
  InProgressBy: String!
  ResolvedAt: Time!
  ResolvedBy: String!
}
```

## Query Definitions

### Mutations
```
extend type Mutation {
  createServiceRequest(
    clientID: String!
    requestType: String!
    request: String
): Boolean!
}
```

### Queries
```
extend type Query {
    getServiceRequests(requestType: String, requestStatus: String): [ServiceRequest]
}
```

### 1. Mutations
#### 1.1. createServiceRequest
create service request allows the user or the system to create a service request.
```
mutation  createServiceRequest($clientID: String!, $requestType: String!, $request: String){
  createServiceRequest(clientID: $clientID, requestType: $requestType, request: $request)
}
```
Variables:
```
{
  "clientID": : "bf3ed095-607c-4c08-a79d-8a82897adb0f",
  "requestType": "RED_FLAG",,
  "request": "red flag request"
}
}
```

#### 1.2.  setInProgressBy
setInProgressBy allows a healthcare practitioner to indicate they have started looking into a client's service request
```
 mutation setInProgressBy($serviceRequestID: String!, $StaffID: String!){
  setInProgressBy(
    serviceRequestID: $serviceRequestID,
    StaffID: $StaffID
  )
}
```

Variables

```
  {
    "serviceRequestID": "8ecbbc80-24c8-421a-9f1a-e14e12678ef1",
    "StaffID": "c1d58335-f441-4112-9ae4-c30e9d630e8c"
  }
```

#### 1.2.  resolveServiceRequest
resolveServiceRequest allows a healthcare practitioner to indicate they have finished looking into a client's service request 
```
mutation resolveServiceRequest($staffID: String!, $requestID: String!){
  resolveServiceRequest(
    staffID: $staffID,
    requestID: $requestID
  )
}
```

Variables

```
  {
    "staffID": "8ecbbc80-24c8-421a-9f1a-e14e12678ef1",
    "requestID": "c1d58335-f441-4112-9ae4-c30e9d630e8c"
  }
```

### 2. Queries
#### 2.1. getServiceRequests
Get service request gets all  service requests, if no params are passed, it will get all service requests
```
query getServiceRequests{
  getServiceRequests{
    ID
    ClientID
    Request
    RequestType
    Status
    ResolvedAt
    ResolvedBy
		InProgressAt
    InProgressBy
  }
}
```

you can pass optional variables to get specific service requests; type, status or both.
```
query getServiceRequests($type: String, $status: String){
  getServiceRequests(requestType: $type, requestStatus: $status){
    ID
    ClientID
    Request
    RequestType
    Status
    ResolvedAt
    ResolvedBy
		InProgressAt
    InProgressBy
  }
}
```

VARIABLES
```
{
  "type": "RED_FLAG",
  "status": "PENDING"
}
```
