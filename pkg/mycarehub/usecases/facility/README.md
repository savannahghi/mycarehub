# Facilities API Spec

This `spec` gives an overview of the API(s) required by the app to create, update, delete and get facilities

The API(s) are organized around `GraphQL`.

Facilities are health institutions that will be assigned to users here-in referred to as `Staff` who must have at least one facility assigned to them. 

A `Client`, who is also a user, must be assigned to a facility from which they will be getting health services. Client can also change their assigned facility for a specified reason e.g. relocation. 

In Kenya, all facilities are uniquely identified with `Master Facility Code (MFL)`. More details about facilities can be found at http://kmhfl.health.go.ke/#/home. In this usecase, we will maintain our list of facilities since we wont be making API calls to KMHFL.

## EndPoint Definitions

Base URLs:
- https://mycarehub-testing.savannahghi.org/ide

- https://mycarehub-prod.savannahghi.org/ide


## 1.  Create/Delete Facilities

Note: Creating or deleting a facility is Idempotent. 

If we create a facility with the same `MFL Code`, it will return the existing facility. 

If we delete a facility for the first time, the facility will be deleted and a boolean value `true` is returned. Deleting the facility for the second time with the same `MFL Code`, it will just return `true`

For this API a facility is created with `FacilityInput` and deleted using `MFL Code`.

### Inputs

```
input FacilityInput {
  name: String!
  code: String!
  active: Boolean!
  county: CountyType!
  description: String!
}
```

### Types

```
type Facility {
  ID: String!
  name: String!
  code: String!
  active: Boolean!
  county: CountyType!
  description: String!
}
```

### Mutation

```
extend type Mutation {
  createFacility(input: FacilityInput!): Facility!
  deleteFacility(mflCode: String!): Boolean!
}
```

### Query Definition

#### Create Facility

```
mutation createFacility($input: FacilityInput!) {
  createFacility (input: $input) {
    ID
    name
    code
    active
    county
    description
  }
}s
```

#### Delete Facility (by MFL Code)

```
mutation deleteFacilityByMFLCode($mflCode: String!) {
  deleteFacility(mflCode: $mflCode)
}
```

### Variables

#### Create Facility

```
{
    "input":{
      "ID":"2063f818-2c24-41ea-8b77-338f5e5366a"
      "name":"Mediheal Hospital (Nakuru) Annex",
      "code":"24929",
      "active": true,
      "county": "Nakuru",
      "description":"located at Giddo plaza building town"
  }
}

```

#### Delete Facility

```

{
  "mflCode": "24929",
}
```

### Output

#### Create Facility

```
{
  "data": {
    "createFacility": {
      "name": "Mediheal Hospital (Nakuru) Annex",
      "code": "24929",
      "active": true,
      "county": "Nakuru",
      "description": "located at Giddo plaza building town"
    }
  }
}
```

#### Delete Facility

```

{
  "data": {
    "deleteFacility": true
  }
}
```

## 2. Get Facility/Facilities

This contains retrieval operations where the user can query all facilities. A facility can be retrieved by `ID` or `mflCode`.

There is also a provision for searching a facility with a search term which looks up for any matching facility.

Filtering can also be done using `name`, `active`, `code`,  and `county` fields.

Pagination is also supported where the `CurrentPage` must be provided, i.e the page they want to view. `Limit` value is optional. If not provided, it defaults to `10` facilities per page

### Input

```
input FiltersInput {
  DataType: FilterDataType
  Value: String
}

input PaginationsInput {
  Limit: Int
  CurrentPage: Int!          
}
}
```

### Types

```
type Facility {
  ID: String!
  name: String!
  code: String!
  active: Boolean!
  county: CountyType!
  description: String!
}

type Pagination {
  Limit: Int!
	CurrentPage: Int!       
	Count:  Int        
	TotalPages: Int
	NextPage:      Int 
  PreviousPage: Int
}

type FacilityPage {
  Pagination: Pagination!
  Facilities: [Facility]!

}
```

### Queries

```
extend type Query {
  fetchFacilities: [Facility]
  retrieveFacility(id: String!, active: Boolean!): Facility
  retrieveFacilityByMFLCode(mflCode: String!, isActive: Boolean!): Facility!
  listFacilities(searchTerm: String, filterInput: [FiltersInput], paginationInput:PaginationsInput!):FacilityPage
}
```

### Query Definition

#### Fetch facilities

```
query fetchFacilities {
  fetchFacilities{
    ID
    name
    code
    active
    county
    description
  }
}
```

#### Retrieve Facility by ID

```
query retrieveFacilityByID($id: String!, $active: Boolean!){
  retrieveFacility(id: $id, active: $active){
    ID
    name
    code
    active
    county
    description
  }
}
```

#### Retrieve Facility by MFL code

```
query retrieveFacilityByMFLCode($mflCode: String!, $isActive: Boolean!){
  retrieveFacilityByMFLCode(mflCode: $mflCode, isActive: $isActive){
    ID
    name
    code
    active
    county
    description
  }
}
```

#### List Facilities (with search, filter and pagination)

```
query listFacilities($searchTerm:String, $filterInput:[FiltersInput], $paginationInput: PaginationsInput!){
  listFacilities(searchTerm: $searchTerm,filterInput:$filterInput, paginationInput: $paginationInput )
  {
    Pagination{ 
      Count
      Limit
      TotalPages
      PreviousPage
      CurrentPage
      NextPage
     
    }
    Facilities{
      ID
      name
      active
      code
      county
      description
    }
  }
}
```


### Variables

#### Retrieve Facility by ID

```
{
  "id": "2063f818-2c24-41ea-8b77-338f5e5366a3",
  "active": true
}
```

#### Retrieve Facility by MFL code

```
{
  "mflCode": "24929",
  "isActive": true
}
```

#### List Facilities (with search, filter,, and pagination)

```
{
  "searchTerm": "Nakuru",
  "filterInput": [
    {"DataType": "name" ,"Value":"Mediheal Hospital (Nakuru) Annex"},
    {"DataType": "active", "Value" :"true"},
    {"DataType": "county" ,"Value":"Nakuru"},
		{"DataType": "mfl_code", "Value":"24929"}
  ],
  "paginationInput": {
    "Limit":10,
    "CurrentPage": 1
  }
}
```

### Output

#### Fetch facilities

```

{
  "data": {
    "fetchFacilities": [
      {
        "ID": "2063f818-2c24-41ea-8b77-338f5e5366a3",
        "name": "Mediheal Hospital (Nakuru) Annex",
        "code": "24929",
        "active": true,
        "county": "Nakuru",
        "description": "located at Giddo plazza building town"
      }
    ]
  }
}
```

#### Retrieve Facility by ID

```
{
  "data": {
    "retrieveFacility": {
      "ID": "2063f818-2c24-41ea-8b77-338f5e5366a3",
      "name": "Mediheal Hospital (Nakuru) Annex",
      "code": "24929",
      "active": true,
      "county": "Nakuru",
      "description": "located at Giddo plaza building town"
    }
  }
}
```

#### Retrieve Facility by MFL code

```
{
  "data": {
    "retrieveFacilityByMFLCode": {
      "ID": "2063f818-2c24-41ea-8b77-338f5e5366a3",
      "name": "Mediheal Hospital (Nakuru) Annex",
      "code": "24929",
      "active": true,
      "county": "Nakuru",
      "description": "located at Giddo plaza building town"
    }
  }
}
```

#### List Facilities (with search, filter, and pagination)

```
{
  "data": {
    "listFacilities": {
      "Pagination": {
        "Count": 1,
        "Limit": 10,
        "TotalPages": 1,
        "PreviousPage": null,
        "CurrentPage": 1,
        "NextPage": null
      },
      "Facilities": [
        {
          "ID": "2063f818-2c24-41ea-8b77-338f5e5366a3",
          "name": "Mediheal Hospital (Nakuru) Annex",
          "active": true,
          "code": "24929",
          "county": "Nakuru",
          "description": "located at Giddo plazza building town"
        }
      ]
    }
  }
}
```

## 3. Inactivate / Re-activate Facilities
These set of APIs are responsible for activating an inactivated facility and re-activating an inactivated facility.
Both of these APIs use `MFL Code` to carry out the inactivation and re-activation operations.

### Inactivate Facility
#### Mutation
```
  mutation inactivateFacility($mflCode: String!) {
		inactivateFacility (mflCode: $mflCode)
	}
```
#### Variables
```
  {
    "mflCode": "<facility MFL Code>"
  }
```

### Reactivate Facility
#### Mutation
```
  mutation reactivateFacility($mflCode: String!) {
		reactivateFacility (mflCode: $mflCode)
	}
```
#### Variables
```
  {
    "mflCode": "<facility MFL Code>"
  }
```
## 3. UpdateFacilities
This APIs is responsible for updating a given facility identified with its `ID `.
This API is responsible for updating the following:
- Name
- Code
- Active
- County
- Description

#### Mutation
```
  mutation updateFacility($id: String! $facilityInput: FacilityInput!) {
  updateFacility (
    id: $id,
    facilityInput: $facilityInput
  ) 
}
```
#### Variables
```
  {
  "id": "d4e63fd7-ec60-4306-ab91-e7274caac1b2",
    "facilityInput":{
        "name":"Mediheal Hospital (Nakuru) Annex",
        "code":"24929",
        "active": true,
        "county": "Nakuru",
        "description":"located at Giddo plaza building town"
  }
}
```
