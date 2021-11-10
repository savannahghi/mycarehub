package terms

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// IGetCurrentTerms represents the interface to get all the terms of service
type IGetCurrentTerms interface {
	GetCurrentTerms(ctx context.Context) (string, error)
}

// UseCasesTerms groups all the logic related to getting the terms of service
type UseCasesTerms interface {
	IGetCurrentTerms
}

// ServiceTermsImpl represents terms implementation object
type ServiceTermsImpl struct {
	Query infrastructure.Query
}

// NewUseCasesTermsOfService is the controler for the terms usecases
func NewUseCasesTermsOfService(
	query infrastructure.Query,
) *ServiceTermsImpl {
	return &ServiceTermsImpl{
		Query: query,
	}
}

//GetCurrentTerms get all the current terms of service
func (t *ServiceTermsImpl) GetCurrentTerms(ctx context.Context) (string, error) {

	return t.Query.GetCurrentTerms(ctx)
}
