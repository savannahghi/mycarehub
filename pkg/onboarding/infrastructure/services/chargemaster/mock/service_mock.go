package mock

import (
	"context"

	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/apiclient"
)

// FakeServiceChargeMaster is an `Chargemaster` service mock .
type FakeServiceChargeMaster struct {
	FetchChargeMasterClientFn func() *apiclient.ServerClient
	FetchProviderByIDFn       func(ctx context.Context, id string) (*domain.BusinessPartner, error)
	FindProviderFn            func(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
		sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error)
	FindBranchFn func(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.BranchFilterInput,
		sort []*dto.BranchSortInput) (*dto.BranchConnection, error)
}

// FetchChargeMasterClient ...
func (f *FakeServiceChargeMaster) FetchChargeMasterClient() *apiclient.ServerClient {
	return f.FetchChargeMasterClientFn()
}

// FindProvider ...
func (f *FakeServiceChargeMaster) FindProvider(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
	sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
	return f.FindProviderFn(ctx, pagination, filter, sort)
}

// FindBranch ...
func (f *FakeServiceChargeMaster) FindBranch(ctx context.Context, pagination *firebasetools.PaginationInput, filter []*dto.BranchFilterInput,
	sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
	return f.FindBranchFn(ctx, pagination, filter, sort)
}

// FetchProviderByID ...
func (f *FakeServiceChargeMaster) FetchProviderByID(ctx context.Context, id string) (*domain.BusinessPartner, error) {
	return f.FetchProviderByIDFn(ctx, id)
}
