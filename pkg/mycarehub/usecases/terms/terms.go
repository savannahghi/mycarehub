package terms

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
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

// GetCurrentTerms get all the current terms of service
func (t *ServiceTermsImpl) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	termsOfService, err := t.Query.GetCurrentTerms(ctx)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, exceptions.ItemNotFoundErr(fmt.Errorf("failed to get current terms of service: %v", err))
	}

	return termsOfService, nil
}

// AcceptTerms can be used to accept or review terms of service
func (t *ServiceTermsImpl) AcceptTerms(ctx context.Context, userID *string, termsID *int) (bool, error) {
	ok, err := t.Update.AcceptTerms(ctx, userID, termsID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, exceptions.FailedToUpdateItemErr(fmt.Errorf("failed to accept terms: %v", err))
	}
	return ok, err
}
