# Feedback API specification

This file gives an overview of the `Feedback` API(s).

Generally, feedback is the reaction to something and can be used as a basis for improvement.
In this context, this API sends `Clients` feedback email to the admin expressing their reaction
towards a particular subject matter.

The API schemas is defined in `GraphQL`.

## EndPoint Definitions

Base URLs:
- https://mycarehub-testing.savannahghi.org/ide

- https://mycarehub-staging.savannahghi.org/ide

- https://mycarehub-prod.savannahghi.org/ide


## Query Definitions

### Mutation
```
extend type Mutation{
    sendFeedback(input: FeedbackResponseInput!): Boolean!
}
```

### 1. Mutation
#### 1.1. Send Feedback
This API sends `clients` feedback via email to the admin.
```
mutation sendFeedback($input: FeedbackResponseInput!){
  sendFeedback(input: $input)
}
```

Variables:
```
{
  "input": {
    "userID": "userID Here",
    "message": "Your message",
    "requiresFollowUp": true
  }
}
```