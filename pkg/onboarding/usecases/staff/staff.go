package staff

import (
	"context"
	"fmt"
	"strings"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

// IRegisterStaffUser contains staff registration methods
type IRegisterStaffUser interface {
	// TODO: ensure default facility is set
	//		validation: ensure the staff profile has at least one facility
	//		ensure that the default facility is one of these
	// TODO: ensure the user exists...userID in profile
	GetOrCreateStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error)
}

// IUpdateStaffUser contains staff update methods
type IUpdateStaffUser interface {
	// TODO: ensure default facility is set
	//		validation: ensure the staff profile has at least one facility
	//		ensure that the default facility is one of these
	// TODO: ensure the user exists...userID in profile
	UpdateStaffUserProfile(ctx context.Context, userID string, user *dto.UserInput, staff *dto.StaffProfileInput) (bool, error)
}

// // IAddRoles contains add staff role methods
// type IAddRoles interface {
// 	AddRoles(userID string, roles []string) (bool, error)
// }

// // IRemoveRole contains remove role methods for staff
// type IRemoveRole interface {
// 	RemoveRole(userID string, role string) (bool, error)
// }

// // IUpdateDefaultFacility contains update default facility methods for staff
// type IUpdateDefaultFacility interface {
// 	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
// 	UpdateDefaultFacility(userID string, facilityID string) (bool, error)
// }

// UsecasesStaffProfile contains all the staff profile usecases
type UsecasesStaffProfile interface {
	IRegisterStaffUser
	IUpdateStaffUser
	// IAddRoles
	// IRemoveRole
	// IUpdateDefaultFacility
}

// UsecasesStaffProfileImpl represents user implementation object
type UsecasesStaffProfileImpl struct {
	Infrastructure infrastructure.Interactor
}

// NewUsecasesStaffProfileImpl returns a new staff profile service
func NewUsecasesStaffProfileImpl(infra infrastructure.Interactor) *UsecasesStaffProfileImpl {
	return &UsecasesStaffProfileImpl{
		Infrastructure: infra,
	}
}

// GetOrCreateStaffUser returns a staff profile
func (u *UsecasesStaffProfileImpl) GetOrCreateStaffUser(ctx context.Context, user *dto.UserInput, staff *dto.StaffProfileInput) (*domain.StaffUserProfile, error) {
	// Try creating a facility
	staffProfile, err := u.Infrastructure.GetOrCreateStaffUser(ctx, user, staff)
	if err != nil {
		contactKeyConstraintError := strings.Contains(err.Error(), "duplicate key value violates unique constraint \"contact_contact_key\"")
		staffNumberConsraintError := strings.Contains(err.Error(), "duplicate key value violates unique constraint \"staffprofile_staff_number_key\"")
		// if we find a duplicate staff number get the staff
		if staffNumberConsraintError || contactKeyConstraintError {
			staffProfileSession, err := u.Infrastructure.GetStaffProfileByStaffNumber(ctx, staff.StaffNumber)
			if err != nil {
				return nil, fmt.Errorf("failed query staff by staff number: %v", err)
			}
			return staffProfileSession, nil
		}

		return nil, fmt.Errorf("failed to get or create staff user profile: %v", err)
	}

	return staffProfile, nil
}

// UpdateStaffUserProfile updates a staff profile
func (u *UsecasesStaffProfileImpl) UpdateStaffUserProfile(ctx context.Context, userID string, user *dto.UserInput, staff *dto.StaffProfileInput) (bool, error) {
	return u.Infrastructure.UpdateStaffUserProfile(ctx, userID, user, staff)
}
