package profile

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
)

func TestFindProvider(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)

	ctx := context.Background()
	first := 10
	after := "0"
	last := 10
	before := "20"
	testSladeCode := "PRO-50"
	ascSort := base.SortOrderAsc
	invalidPage := "invalidpage"

	tests := map[string]struct {
		expectNonNilConnection bool
		expectedErr            error

		pagination *base.PaginationInput
		filter     []*BusinessPartnerFilterInput
		sort       []*BusinessPartnerSortInput

		testSingleFetches bool
	}{
		"query_params_only_no_pagination_filter_or_sort_params": {
			expectNonNilConnection: true,
			expectedErr:            nil,
		},
		"with_forward_pagination": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			pagination: &base.PaginationInput{
				First: first,
				After: after,
			},
			testSingleFetches: true,
		},
		"with_backward_pagination": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			pagination: &base.PaginationInput{
				Last:   last,
				Before: before,
			},
		},
		"with_filter": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			filter: []*BusinessPartnerFilterInput{
				{
					SladeCode: &testSladeCode,
				},
			},
		},
		"with_sort": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			sort: []*BusinessPartnerSortInput{
				{
					Name:      &ascSort,
					SladeCode: &ascSort,
				},
			},
		},
		"with_invalid_pagination": {
			expectNonNilConnection: false,
			expectedErr:            errors.New("expected `after` to be parseable as an int; got invalidpage"),
			pagination: &base.PaginationInput{
				After: invalidPage,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conn, err := service.FindProvider(ctx, tc.pagination, tc.filter, tc.sort)
			if tc.expectNonNilConnection {
				assert.NotNil(t, conn)
			}
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			}
			if err != nil {
				log.Printf("Error: %#v", err)
			}
		})
	}
}

func TestFindBranch(t *testing.T) {
	service := NewService()
	assert.NotNil(t, service)

	ctx := context.Background()
	first := 10
	after := "0"
	last := 10
	before := "20"
	testSladeCode := "PRO-50"
	ascSort := base.SortOrderAsc
	invalidPage := "invalidpage"

	tests := map[string]struct {
		expectNonNilConnection bool
		expectedErr            error

		pagination *base.PaginationInput
		filter     []*BranchFilterInput
		sort       []*BranchSortInput

		testSingleFetches bool
	}{
		"query_params_only_no_pagination_filter_or_sort_params": {
			expectNonNilConnection: true,
			expectedErr:            nil,
		},
		"with_forward_pagination": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			pagination: &base.PaginationInput{
				First: first,
				After: after,
			},
			testSingleFetches: true,
		},
		"with_backward_pagination": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			pagination: &base.PaginationInput{
				Last:   last,
				Before: before,
			},
		},
		"with_filter": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			filter: []*BranchFilterInput{
				{
					SladeCode: &testSladeCode,
				},
			},
		},
		"with_sort": {
			expectNonNilConnection: true,
			expectedErr:            nil,
			sort: []*BranchSortInput{
				{
					Name:      &ascSort,
					SladeCode: &ascSort,
				},
			},
		},
		"with_invalid_pagination": {
			expectNonNilConnection: false,
			expectedErr:            errors.New("expected `after` to be parseable as an int; got invalidpage"),
			pagination: &base.PaginationInput{
				After: invalidPage,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			conn, err := service.FindBranch(ctx, tc.pagination, tc.filter, tc.sort)
			if tc.expectNonNilConnection {
				assert.NotNil(t, conn)
			}
			if tc.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tc.expectedErr.Error())
			}
			if err != nil {
				log.Printf("Error: %#v", err)
			}
		})
	}
}
