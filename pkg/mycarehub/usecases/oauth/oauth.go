package oauth

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// UseCasesCommunities holds all interfaces required to implement the communities feature
type UseCasesOauth interface {
	CreateOauthClient(ctx context.Context, input dto.OauthClientInput) (*domain.OauthClient, error)
}

// UseCasesOauthImpl represents oauth implementation
type UseCasesOauthImpl struct {
	Update infrastructure.Update
	Query  infrastructure.Query
	Create infrastructure.Create
	Delete infrastructure.Delete
}

// NewUseCasesOauthImplementation initializes an implementation of the fosite storage
func NewUseCasesOauthImplementation(create infrastructure.Create, update infrastructure.Update, query infrastructure.Query, delete infrastructure.Delete) *UseCasesOauthImpl {
	return &UseCasesOauthImpl{
		Update: update,
		Query:  query,
		Create: create,
		Delete: delete,
	}
}
