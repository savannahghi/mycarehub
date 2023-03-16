package terms

import (
	"context"
	"fmt"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ICreateTerms represents the interface to create terms
type ICreateTerms interface {
	CreateTermsOfService(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error)
}

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
	ICreateTerms
}

// ServiceTermsImpl represents terms implementation object
type ServiceTermsImpl struct {
	Query  infrastructure.Query
	Update infrastructure.Update
	Create infrastructure.Create
}

// NewUseCasesTermsOfService is the controler for the terms usecases
func NewUseCasesTermsOfService(
	query infrastructure.Query,
	update infrastructure.Update,
	create infrastructure.Create,
) *ServiceTermsImpl {
	return &ServiceTermsImpl{
		Query:  query,
		Update: update,
		Create: create,
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

// CreateTermsOfService creates a new record of terms in  the database
func (t *ServiceTermsImpl) CreateTermsOfService(ctx context.Context, termsOfService *domain.TermsOfService) (*domain.TermsOfService, error) {
	return t.Create.CreateTermsOfService(ctx, termsOfService)
}
