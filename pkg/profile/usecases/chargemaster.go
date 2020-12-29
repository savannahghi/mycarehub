package usecases

import (
	"context"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/profile/domain"
)

// ChargeMasterUseCases represents logic required to communicate with chargemaster
type ChargeMasterUseCases interface {
	FindProvider(ctx context.Context, pagination *base.PaginationInput, filter []*domain.BusinessPartnerFilterInput, sort []*domain.BusinessPartnerSortInput) (*domain.BusinessPartnerConnection, error)
	FindBranch(ctx context.Context, pagination *base.PaginationInput, filter []*domain.BranchFilterInput, sort []*domain.BranchSortInput) (*domain.BranchConnection, error)
}
