package terms

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetCurrentTerms represents the interface to get all the terms of service
type IGetCurrentTerms interface {
	GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error)
}

// IAcceptTerms represents hold the accept terms method
type IAcceptTerms interface {
	AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error)
}

// UseCasesTerms groups all the logic related to getting the terms of service
type UseCasesTerms interface {
	IGetCurrentTerms
	IAcceptTerms
}

// ServiceTermsImpl represents terms implementation object
type ServiceTermsImpl struct {
	Query  infrastructure.Query
	Update infrastructure.Update
}

// NewUseCasesTermsOfService is the controler for the terms usecases
func NewUseCasesTermsOfService(
	query infrastructure.Query,
	update infrastructure.Update,
) *ServiceTermsImpl {
	return &ServiceTermsImpl{
		Query:  query,
		Update: update,
	}
}

//GetCurrentTerms get all the current terms of service
func (t *ServiceTermsImpl) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {

	return t.Query.GetCurrentTerms(ctx)
}

// AcceptTerms can be used to accept or review terms of service
func (t *ServiceTermsImpl) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	return t.Update.AcceptTerms(ctx, userID, termsID)
}
