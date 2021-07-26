package graph

import (
	"context"
	"log"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/presentation/interactor"

	"firebase.google.com/go/auth"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	interactor *interactor.Interactor
}

//go:generate go run github.com/99designs/gqlgen

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	interactor *interactor.Interactor,

) (*Resolver, error) {
	return &Resolver{
		interactor: interactor,
	}, nil
}

func (r Resolver) checkPreconditions() {
	if r.interactor.Onboarding == nil {
		log.Panicf("expected onboarding usecases to be defined resolver")
	}

	if r.interactor.Signup == nil {
		log.Panicf("expected signup usecases to be define in resolver ")
	}
}

// CheckUserTokenInContext ensures that the context has a valid Firebase auth token
func (r *Resolver) CheckUserTokenInContext(ctx context.Context) *auth.Token {
	token, err := firebasetools.GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("graph.Resolver: context user token is nil")
	}
	return token
}
