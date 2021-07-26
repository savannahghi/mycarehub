package chargemaster

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/onboarding/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding/pkg/onboarding/domain"
	"github.com/stretchr/testify/assert"
)

func TestServiceChargeMasterImpl_FindProvider(t *testing.T) {
	cm := NewChargeMasterUseCasesImpl()
	assert.NotNil(t, cm)
	type args struct {
		ctx        context.Context
		pagination *firebasetools.PaginationInput
		filter     []*dto.BusinessPartnerFilterInput
		sort       []*dto.BusinessPartnerSortInput
	}
	first := 10
	after := "0"
	last := 10
	before := "20"
	testSladeCode := "PRO-50"
	ascSort := enumutils.SortOrderAsc
	invalidPage := "invalidpage"

	tests := []struct {
		name                   string
		args                   args
		wantErr                bool
		expectNonNilConnection bool
		expectedErr            error
	}{
		{
			name:                   "happy case - query params only no pagination filter or sort params",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx:        context.Background(),
				pagination: &firebasetools.PaginationInput{},
				filter:     []*dto.BusinessPartnerFilterInput{},
				sort:       []*dto.BusinessPartnerSortInput{},
			},
		},
		{
			name:                   "happy case - with forward pagination",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx: context.Background(),
				pagination: &firebasetools.PaginationInput{
					First: first,
					After: after,
				},
				filter: []*dto.BusinessPartnerFilterInput{},
				sort:   []*dto.BusinessPartnerSortInput{},
			},
		},
		{
			name:                   "happy case - with backward pagination",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx: context.Background(),
				pagination: &firebasetools.PaginationInput{
					Last:   last,
					Before: before,
				},
				filter: []*dto.BusinessPartnerFilterInput{},
				sort:   []*dto.BusinessPartnerSortInput{},
			},
		},
		{
			name:                   "happy case - with filter",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx:        context.Background(),
				pagination: &firebasetools.PaginationInput{},
				filter: []*dto.BusinessPartnerFilterInput{
					{
						SladeCode: &testSladeCode,
					},
				},
				sort: []*dto.BusinessPartnerSortInput{},
			},
		},
		{
			name:                   "happy case - with sort",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx:        context.Background(),
				pagination: &firebasetools.PaginationInput{},
				filter:     []*dto.BusinessPartnerFilterInput{},
				sort: []*dto.BusinessPartnerSortInput{
					{
						Name:      &ascSort,
						SladeCode: &ascSort,
					},
				},
			},
		},
		{
			name:                   "sad case - with invalid pagination",
			expectNonNilConnection: false,
			expectedErr:            errors.New("expected `after` to be parseable as an int; got invalidpage"),
			wantErr:                true,
			args: args{
				ctx: context.Background(),
				pagination: &firebasetools.PaginationInput{
					After: invalidPage,
				},
				filter: []*dto.BusinessPartnerFilterInput{},
				sort:   []*dto.BusinessPartnerSortInput{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cm.FindProvider(tt.args.ctx, tt.args.pagination, tt.args.filter, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceChargeMasterImpl.FindProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expectNonNilConnection {
				assert.NotNil(t, got)
			}
			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func TestServiceChargeMasterImpl_FindBranch(t *testing.T) {
	cm := NewChargeMasterUseCasesImpl()
	assert.NotNil(t, cm)
	type args struct {
		ctx        context.Context
		pagination *firebasetools.PaginationInput
		filter     []*dto.BranchFilterInput
		sort       []*dto.BranchSortInput
	}
	first := 10
	after := "0"
	last := 10
	before := "20"
	testSladeCode := "PRO-50"
	ascSort := enumutils.SortOrderAsc
	invalidPage := "invalidpage"

	tests := []struct {
		name                   string
		args                   args
		wantErr                bool
		expectNonNilConnection bool
		expectedErr            error
	}{
		{
			name:                   "happy case - query params only no pagination filter or sort params",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx:        context.Background(),
				pagination: &firebasetools.PaginationInput{},
				filter:     []*dto.BranchFilterInput{},
				sort:       []*dto.BranchSortInput{},
			},
		},
		{
			name:                   "happy case - with forward pagination",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx: context.Background(),
				pagination: &firebasetools.PaginationInput{
					First: first,
					After: after,
				},
				filter: []*dto.BranchFilterInput{},
				sort:   []*dto.BranchSortInput{},
			},
		},
		{
			name:                   "happy case - with backward pagination",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx: context.Background(),
				pagination: &firebasetools.PaginationInput{
					Last:   last,
					Before: before,
				},
				filter: []*dto.BranchFilterInput{},
				sort:   []*dto.BranchSortInput{},
			},
		},
		{
			name:                   "happy case -with filter",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx:        context.Background(),
				pagination: &firebasetools.PaginationInput{},
				filter: []*dto.BranchFilterInput{
					{
						SladeCode: &testSladeCode,
					},
				},
				sort: []*dto.BranchSortInput{},
			},
		},
		{
			name:                   "happy case - with sort",
			expectNonNilConnection: true,
			expectedErr:            nil,
			wantErr:                false,
			args: args{
				ctx:        context.Background(),
				pagination: &firebasetools.PaginationInput{},
				filter:     []*dto.BranchFilterInput{},
				sort: []*dto.BranchSortInput{
					{
						Name:      &ascSort,
						SladeCode: &ascSort,
					},
				},
			},
		},
		{
			name:                   "sad case - with invalid pagination",
			expectNonNilConnection: false,
			expectedErr:            errors.New("expected `after` to be parseable as an int; got invalidpage"),
			wantErr:                true,
			args: args{
				ctx: context.Background(),
				pagination: &firebasetools.PaginationInput{
					After: invalidPage,
				},
				filter: []*dto.BranchFilterInput{},
				sort:   []*dto.BranchSortInput{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cm.FindBranch(tt.args.ctx, tt.args.pagination, tt.args.filter, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceChargeMasterImpl.FindBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.expectNonNilConnection {
				assert.NotNil(t, got)
			}
			if tt.expectedErr != nil {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			}
		})
	}
}

func Test_parentOrgSladeCodeFromBranch(t *testing.T) {
	type args struct {
		branch *domain.BusinessPartner
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				branch: &domain.BusinessPartner{
					SladeCode: "BRA-PRO-4313-1",
				},
			},
			want:    "PRO-4313",
			wantErr: false,
		},
		{
			name: "no BRA prefix",
			args: args{
				branch: &domain.BusinessPartner{
					SladeCode: "PRO-4313-1",
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "sad case (long branch slade code)",
			args: args{
				branch: &domain.BusinessPartner{
					SladeCode: "BRA-PRO-4313-1-9393030",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parentOrgSladeCodeFromBranch(tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("parentOrgSladeCodeFromBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parentOrgSladeCodeFromBranch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServiceChargeMasterImpl_FetchProviderByID(t *testing.T) {
	ctx := context.Background()
	cm := NewChargeMasterUseCasesImpl()

	pagination := &firebasetools.PaginationInput{}
	filter := []*dto.BusinessPartnerFilterInput{}
	sort := []*dto.BusinessPartnerSortInput{}

	partners, err := cm.FindProvider(ctx, pagination, filter, sort)
	if err != nil {
		t.Errorf("can't find provider: %w", err)
		return
	}

	partner := partners.Edges[0].Node

	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.BusinessPartner
		wantErr bool
	}{
		{
			name: "happy case: valid",
			args: args{
				ctx: ctx,
				id:  partner.ID,
			},
			want:    partner,
			wantErr: false,
		},
		{
			name: "sad case: invalid ID",
			args: args{
				ctx: ctx,
				id:  "InvalidID",
			},
			wantErr: true,
		},
		{
			name: "sad case: empty ID",
			args: args{
				ctx: ctx,
				id:  "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cm.FetchProviderByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceChargeMasterImpl.FetchProviderByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ServiceChargeMasterImpl.FetchProviderByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
