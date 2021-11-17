package postgres

import (
	"context"
	"fmt"

	"github.com/savannahghi/feedlib"
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
		singleFacility := domain.Facility{
			ID:          m.FacilityID,
			Name:        m.Name,
			Code:        m.Code,
			Active:      m.Active,
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
func (d *MyCareHubDb) RetrieveFacilityByMFLCode(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
	if MFLCode == 0 {
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

// GetCurrentTerms fetches the current terms service
func (d *MyCareHubDb) GetCurrentTerms(ctx context.Context) (*domain.TermsOfService, error) {
	terms, err := d.query.GetCurrentTerms(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get current terms of service: %v", err)
	}

	return d.mapTermsOfServiceObjectToDomain(terms), nil
}

// GetUserProfileByUserID fetches and returns a userprofile using their user ID
func (d *MyCareHubDb) GetUserProfileByUserID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID should be provided")
	}

	user, err := d.query.GetUserProfileByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile by user ID: %v", err)
	}

	return d.mapProfileObjectToDomain(user), nil
}

// GetSecurityQuestions fetches all the security questions
func (d *MyCareHubDb) GetSecurityQuestions(ctx context.Context, flavour feedlib.Flavour) ([]*domain.SecurityQuestion, error) {
	var securityQuestion []*domain.SecurityQuestion

	allSecurityQuestions, err := d.query.GetSecurityQuestions(ctx, flavour)
	if err != nil {
		return nil, fmt.Errorf("unable to get security questions: %v", err)
	}

	if len(allSecurityQuestions) == 0 {
		return securityQuestion, nil
	}

	for _, sq := range allSecurityQuestions {
		singleSecurityQuestion := &domain.SecurityQuestion{
			SecurityQuestionID: *sq.SecurityQuestionID,
			QuestionStem:       sq.QuestionStem,
			Description:        sq.Description,
			Flavour:            sq.Flavour,
			Active:             sq.Active,
			ResponseType:       sq.ResponseType,
		}

		securityQuestion = append(securityQuestion, singleSecurityQuestion)
	}

	return securityQuestion, nil
}

// GetSecurityQuestionByID fetches a security question by ID
func (d *MyCareHubDb) GetSecurityQuestionByID(ctx context.Context, securityQuestionID *string) (*domain.SecurityQuestion, error) {
	securityQuestion, err := d.query.GetSecurityQuestionByID(ctx, securityQuestionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get security question by ID: %v", err)
	}

	return d.mapSecurityQuestionObjectToDomain(securityQuestion), nil
}
