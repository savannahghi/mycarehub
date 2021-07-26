# Onboarding service

[![pipeline status](https://github.com/savannahghi/onboarding/badges/develop/pipeline.svg)](https://github.com/savannahghi/onboarding/-/commits/develop)
[![coverage report](https://github.com/savannahghi/onboarding/badges/develop/coverage.svg)](https://github.com/savannahghi/onboarding/-/commits/develop)

This service manages user onboarding process.

## Description

The project implements the `Clean Architecture` advocated by
Robert Martin ('Uncle Bob').

## Documentation

API documentation are available at: https://profile-service-docs-uyajqt434q-ew.a.run.app/

We should strive to always make the documentation reflect the true state of the service

### Clean Architecture

A cleanly architected project should be:

- _Independent of Frameworks_: The architecture does not depend on the
  existence of some library of feature laden software. This allows you to use
  such frameworks as tools, rather than having to cram your system into their
  limited constraints.

- _Testable_: The business rules can be tested without the UI, Database,
  Web Server, or any other external element.

- _Independent of UI_: The UI can change easily, without changing the rest of
  the system. A Web UI could be replaced with a console UI, for example,
  without changing the business rules.

- _Independent of Database_: You can swap out Cloud Firestore or SQL Server,
  for Mongo, Postgres, MySQL, or something else. Your business rules are not
  bound to the database.

- _Independent of any external agency_: In fact your business rules simply
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

### Points to note

- Interfaces let Go programmers describe what their package provides–not how it does it. This is all just another way of saying “decoupling”, which is indeed the goal, because software that is loosely coupled is software that is easier to change.
- Design your public API/ports to keep secrets(Hide implementation details)
  abstract information that you present so that you can change your implementation behind your public API without changing the contract of exchanging information with other services.

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

export REPOSITORY="firebase" # when we switch to PG the value will be `postgres`

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

# login

[![pipeline status](https://gitlab.slade360emr.com/go/login/badges/develop/pipeline.svg)](https://gitlab.slade360emr.com/go/login/-/commits/develop)
[![coverage report](https://gitlab.slade360emr.com/go/login/badges/develop/coverage.svg)](https://gitlab.slade360emr.com/go/login/-/commits/develop)

A login, logout and refresh micro-service.

# Endpoints

The production service is deployed at https://login-prod.healthcloud.co.ke .

The test service is deployed at https://login-test.healthcloud.co.ke .

There's a service wired to _EDI Core_'s auth server at https://login-core.healthcloud.co.ke/ .

There's another login service wired to _EDI Portal_'s auth server at https://login-portal.healthcloud.co.ke/ .

There's a different login service wired to _Multi-tenant_ auth server at https://login-multitenant.healthcloud.co.ke/ .

# Logging in from the command line

In order to explore the API, you need to log in.

Please adapt the `curl` command below:

```
curl -v -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"username": "<a username>", "password": "<a password>"}' https://login-prod.healthcloud.co.ke/login
```

After logging in, you should look for the `id_token` - that's the value that you need to supply
to other APIs as part of the `Authorization: Bearer <token>` header.

NB:

1. For the test API, send to https://login-test.healthcloud.co.ke/login
2. For staging API, send to https://login-staging.healthcloud.co.ke/login
3. If you want to know which _specific_ auth server is getting called, please take a look at the
   environment variables at https://console.cloud.google.com/run/detail/europe-west1/login/revisions?project=bewell-app-testing or
   https://console.cloud.google.com/run/detail/europe-west1/login/revisions?folder=&organizationId=&project=bewell-app .

# Logging in from a web browser (Slade application)

Get an auth server `accessToken` from the regular Slade web app login process.

Exchange it for an ID token as indicated below:

```
curl -v -i -H "Accept: application/json" -H "Content-Type: application/json" -d '{"accessToken": "<an access token from a valid login session>"}' https://login-prod.healthcloud.co.ke/verify_access_token
```

NB: send to the test endpoint when testing

The `accessToken` should come from a valid Slade 360 EDI login.

## Important note about the auth server

The access token you use should correspond to the correct EDI frontend, login service and auth server.

For example, the `test` login service is, at the time of writing, connected to accounts-healthcloud.multitenant.slade360.co.ke .
You can check it yourself at https://console.cloud.google.com/run/detail/europe-west1/login/revisions?project=bewell-app-testing .
A suitable "release" front-end for experimentation would therefore be https://healthcloud.multitenant.slade360.co.ke/eligibility .

Two `prod` login services will need to be deployed after testing: one that is set up with the portal auth server
and another that is set up with the EDI core auth server. Providers paying on portal will have their requests
verified against the portal auth server; our staff on core will have their requests verified against the core
auth server.

# Switch user to opt-in or opt-out to flagged features

```sh
http https://profile-testing.healthcloud.co.ke/switch_flagged_features phoneNumber="<phone-number-of-user>"

```

Replace `http` if using `curl`. 
Replace `https://profile-testing.healthcloud.co.ke` with `https://profile-prod.healthcloud.co.ke` if running in PROD environment
