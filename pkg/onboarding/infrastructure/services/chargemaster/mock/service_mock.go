package mock

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceChargeMaster is an `Chargemaster` service mock .
type FakeServiceChargeMaster struct {
	FetchChargeMasterClientFn func() *base.ServerClient
	FetchProviderByIDFn       func(ctx context.Context, id string) (*domain.BusinessPartner, error)
	FindProviderFn            func(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BusinessPartnerFilterInput,
		sort []*resources.BusinessPartnerSortInput) (*resources.BusinessPartnerConnection, error)
	FindBranchFn func(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BranchFilterInput,
		sort []*resources.BranchSortInput) (*resources.BranchConnection, error)
}

// FetchChargeMasterClient ...
func (f *FakeServiceChargeMaster) FetchChargeMasterClient() *base.ServerClient {
	return f.FetchChargeMasterClientFn()
}

// FindProvider ...
func (f *FakeServiceChargeMaster) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BusinessPartnerFilterInput,
	sort []*resources.BusinessPartnerSortInput) (*resources.BusinessPartnerConnection, error) {
	return f.FindProviderFn(ctx, pagination, filter, sort)
}

// FindBranch ...
func (f *FakeServiceChargeMaster) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*resources.BranchFilterInput,
	sort []*resources.BranchSortInput) (*resources.BranchConnection, error) {
	return f.FindBranchFn(ctx, pagination, filter, sort)
}

// FetchProviderByID ...
func (f *FakeServiceChargeMaster) FetchProviderByID(ctx context.Context, id string) (*domain.BusinessPartner, error) {
	return f.FetchProviderByIDFn(ctx, id)
}
