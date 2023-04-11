package oauth

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// UseCasesOauthImpl represents oauth implementation
type UseCasesOauthImpl struct {
	Update infrastructure.Update
	Query  infrastructure.Query
	Create infrastructure.Create
}

// NewUseCasesOauthImplementation initializes an implementation of the fosite storage
func NewUseCasesOauthImplementation(create infrastructure.Create, update infrastructure.Update, query infrastructure.Query) *UseCasesOauthImpl {
	return &UseCasesOauthImpl{
		Update: update,
		Query:  query,
		Create: create,
	}
}
