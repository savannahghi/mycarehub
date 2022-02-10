package graph

import (
	"context"
	"log"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"

	"firebase.google.com/go/auth"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	mycarehub *usecases.MyCareHub
}

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	mycarehub usecases.MyCareHub,
) (*Resolver, error) {
	return &Resolver{
		mycarehub: &mycarehub,
	}, nil
}

func (r Resolver) checkPreconditions() {
	if r.mycarehub == nil {
		log.Panicf("expected mycarehub usecases to be defined resolver")
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
