# MyCareHub service

![Linting and Tests](https://github.com/savannahghi/mycarehub/actions/workflows/ci.yml/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/savannahghi/mycarehub/badge.svg?branch=develop)](https://coveralls.io/github/savannahghi/mycarehub?branch=develop)

This service contains the implementation of the mycarehub project.

## Description

The project implements the `Clean Architecture` advocated by
Robert Martin ('Uncle Bob').

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
export REPOSITORY="firebase" # when we switch to PG the value will be `postgres`

```

## Deployment

This application is deployed via Google Cloud Build ( <https://cloud.google.com/build> ) to Google Cloud Run ( <https://cloud.google.com/run> ).
There's a `cloudbuild.yaml` file in the home folder. Secrets (e.g production settings) are managed with Google Secret Manager ( <https://cloud.google.com/secret-manager> ).
