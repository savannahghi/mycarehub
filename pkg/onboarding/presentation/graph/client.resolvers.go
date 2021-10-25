package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/graph/generated"
)

func (r *mutationResolver) RegisterClientUser(ctx context.Context, userInput dto.UserInput, clientInput dto.ClientProfileInput) (*domain.ClientUserProfile, error) {
	return r.interactor.ClientUseCase.RegisterClient(ctx, &userInput, &clientInput)
}

func (r *mutationResolver) AddIdentifier(ctx context.Context, clientID string, idType enums.IdentifierType, idValue string, isPrimary bool) (*domain.Identifier, error) {
	return r.interactor.ClientUseCase.AddIdentifier(ctx, clientID, idType, idValue, isPrimary)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
