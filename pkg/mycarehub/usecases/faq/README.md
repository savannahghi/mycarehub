# FAQs API specification

This file gives an overview of the `Frequently Asked Questions (FAQ)` API(s).

FAQs are a list of questions and answers relating to a particular subject, especially one giving basic information for users of a website or mobile app.

The API schemas is defined in `GraphQL`.

## EndPoint Definitions

Base URLs:
- https://mycarehub-testing.savannahghi.org/ide

- https://mycarehub-staging.savannahghi.org/ide

- https://mycarehub-prod.savannahghi.org/ide


## Query Definitions

### Queries
```
extend type Query {
    getFAQContent(flavour: Flavour! limit: Int): [FAQ!]!
}
```

### 2. Queries
#### 2.1. Get FAQs
This API fetches all the frequently asked question from the users.
```
query getFAQContent($flavour: Flavour!, $limit: Int){
  getFAQContent(flavour: $flavour, limit: $limit){
    ID
    Active
    Title
    Description
    Body
    Flavour
  }
}
```

Variables:
```
{
  "flavour": "PRO",
  "limit": 5
}
```
