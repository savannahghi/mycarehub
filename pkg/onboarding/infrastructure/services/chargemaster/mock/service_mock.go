package mock

import (
	"context"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/dto"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"

	"gitlab.slade360emr.com/go/base"
)

// FakeServiceChargeMaster is an `Chargemaster` service mock .
type FakeServiceChargeMaster struct {
	FetchChargeMasterClientFn func() *base.ServerClient
	FetchProviderByIDFn       func(ctx context.Context, id string) (*domain.BusinessPartner, error)
	FindProviderFn            func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
		sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error)
	FindBranchFn func(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
		sort []*dto.BranchSortInput) (*dto.BranchConnection, error)
}

// FetchChargeMasterClient ...
func (f *FakeServiceChargeMaster) FetchChargeMasterClient() *base.ServerClient {
	return f.FetchChargeMasterClientFn()
}

// FindProvider ...
func (f *FakeServiceChargeMaster) FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BusinessPartnerFilterInput,
	sort []*dto.BusinessPartnerSortInput) (*dto.BusinessPartnerConnection, error) {
	return f.FindProviderFn(ctx, pagination, filter, sort)
}

// FindBranch ...
func (f *FakeServiceChargeMaster) FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*dto.BranchFilterInput,
	sort []*dto.BranchSortInput) (*dto.BranchConnection, error) {
	return f.FindBranchFn(ctx, pagination, filter, sort)
}

// FetchProviderByID ...
func (f *FakeServiceChargeMaster) FetchProviderByID(ctx context.Context, id string) (*domain.BusinessPartner, error) {
	return f.FetchProviderByIDFn(ctx, id)
}
