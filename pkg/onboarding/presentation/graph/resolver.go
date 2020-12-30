package graph

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/presentation/service"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	srv *service.Service
}

//go:generate go run github.com/99designs/gqlgen

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	service *service.Service,

) (*Resolver, error) {
	return &Resolver{
		srv: service,
	}, nil
}

func (r Resolver) checkPreconditions() {
	if r.srv.Onboarding == nil {
		log.Panicf("expected onboarding usecases to be defined resolver")
	}

	if r.srv.Signup == nil {
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
