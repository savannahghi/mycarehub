package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

//GetFacilities returns a slice of healthcare facilities in the platform.
func (d *MyCareHubDb) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	var facility []*domain.Facility
	facilities, err := d.query.GetFacilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	if len(facilities) == 0 {
		return facility, nil
	}
	for _, m := range facilities {
		active, err := strconv.ParseBool(m.Active)
		if err != nil {
			return nil, fmt.Errorf("failed to parse facility.Active to boolean")
		}
		singleFacility := domain.Facility{
			ID:          m.FacilityID,
			Name:        m.Name,
			Code:        m.Code,
			Active:      active,
			County:      m.County,
			Description: m.Description,
		}

		facility = append(facility, &singleFacility)
	}

	return facility, nil
}

// RetrieveFacility gets a facility by ID from the database
func (d *MyCareHubDb) RetrieveFacility(ctx context.Context, id *string, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacility(ctx, id, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve one facility: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// RetrieveFacilityByMFLCode gets a facility by ID from the database
func (d *MyCareHubDb) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	if MFLCode == "" {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve facility by MFLCode: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// ListFacilities gets facilities that are filtered from search and filter,
// the results are also paginated
func (d *MyCareHubDb) ListFacilities(
	ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationsInput *dto.PaginationsInput) (*domain.FacilityPage, error) {
	// if user did not provide current page, throw an error
	if err := paginationsInput.Validate(); err != nil {
		return nil, fmt.Errorf("pagination input validation failed: %v", err)
	}

	sortOutput := &domain.SortParam{
		Field:     paginationsInput.Sort.Field,
		Direction: paginationsInput.Sort.Direction,
	}
	paginationOutput := domain.FacilityPage{
		Pagination: domain.Pagination{
			Limit:       paginationsInput.Limit,
			CurrentPage: paginationsInput.CurrentPage,
			Sort:        sortOutput,
		},
	}
	filtersOutput := []*domain.FiltersParam{}
	for _, f := range filterInput {
		filter := &domain.FiltersParam{
			Name:     string(f.DataType),
			DataType: f.DataType,
			Value:    f.Value,
		}
		filtersOutput = append(filtersOutput, filter)
	}

	facilities, err := d.query.ListFacilities(ctx, searchTerm, filtersOutput, &paginationOutput)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}
	return facilities, nil
}

// GetUserProfileByPhoneNumber fetches and returns a userprofile using their phonenumber
func (d *MyCareHubDb) GetUserProfileByPhoneNumber(ctx context.Context, phoneNumber string) (*domain.User, error) {
	if phoneNumber == "" {
		return nil, fmt.Errorf("phone number should be provided")
	}

	user, err := d.query.GetUserProfileByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile by phonenumber: %v", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetUserPINByUserID fetches a user pin by the user ID
func (d *MyCareHubDb) GetUserPINByUserID(ctx context.Context, userID string) (*domain.UserPIN, error) {
	pinData, err := d.query.GetUserPINByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve user PIN data: %s", err)
	}

	return d.mapPINObjectToDomain(pinData), nil
}
