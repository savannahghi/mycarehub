package graph

//go:generate go run github.com/99designs/gqlgen

import (
	"context"
	"log"

	"firebase.google.com/go/auth"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/graph/profile"
)

// NewResolver sets up a properly initialized resolver
func NewResolver() *Resolver {
	return &Resolver{
		profileService: profile.NewService(),
	}
}

// Resolver sets up dependencies needed by the query and mutation resolvers
type Resolver struct {
	profileService *profile.Service
}

// CheckUserTokenInContext ensures that the context has a valid Firebase auth token
func (r *Resolver) CheckUserTokenInContext(ctx context.Context) *auth.Token {
	token, err := base.GetUserTokenFromContext(ctx)
	if err != nil {
		log.Panicf("graph.Resolver: context user token is nil")
	}
	return token
}

// CheckDependencies ensures that the resolver has what it needs in order to work
func (r *Resolver) CheckDependencies() {
	if r.profileService == nil {
		log.Panicf("graph.Resolver: profileService is nil")
	}
}
