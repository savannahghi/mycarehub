package graph

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/usecases"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	onboarding *usecases.OnboardingUseCaseImpl
	signup     *usecases.SignUpUseCasesImpl
}

//go:generate go run github.com/99designs/gqlgen

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	ob *usecases.OnboardingUseCaseImpl,
	su *usecases.SignUpUseCasesImpl,

) (*Resolver, error) {
	return &Resolver{
		onboarding: ob,
	}, nil
}

func (r Resolver) checkPreconditions() {
	if r.onboarding == nil {
		log.Panicf("expected onboarding usecases to be defined resolver")
	}

	if r.signup == nil {
		log.Panicf("expected signup usecases to be define in resolver ")
	}
}

// CheckUserTokenInContext ensures that the context has a valid Firebase auth token
func (r *Resolver) CheckUserTokenInContext(ctx context.Context) *auth.Token {
	token, err := base.GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("graph.Resolver: context user token is nil")
	}
	return token
}
