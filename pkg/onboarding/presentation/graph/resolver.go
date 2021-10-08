package graph

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
)

//go:generate go run github.com/99designs/gqlgen

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver sets up a GraphQL resolver with all necessary dependencies
type Resolver struct {
	interractor *interactor.Interactor
}

// NewResolver sets up the dependencies needed for query and mutation resolvers to work
func NewResolver(
	ctx context.Context,
	interractor *interactor.Interactor,
) (*Resolver, error) {
	return &Resolver{
		interractor,
	}, nil
}
