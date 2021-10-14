package staff

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
)

// IRegisterStaffUser contains staff registration methods
type IRegisterStaffUser interface {
	// TODO: ensure default facility is set
	//		validation: ensure the staff profile has at least one facility
	//		ensure that the default facility is one of these
	// TODO: ensure the user exists...userID in profile
	RegisterStaffUser(user domain.User, profile domain.StaffProfile) (*domain.User, *domain.StaffProfile, error)
}

// IAddRoles contains add staff role methods
type IAddRoles interface {
	AddRoles(userID string, roles []string) (bool, error)
}

// IRemoveRole contains remove role methods for staff
type IRemoveRole interface {
	RemoveRole(userID string, role string) (bool, error)
}

// IUpdateDefaultFacility contains update default facility methods for staff
type IUpdateDefaultFacility interface {
	// TODO: the list of facilities to switch between is strictly those that the user is assigned to
	UpdateDefaultFacility(userID string, facilityID string) (bool, error)
}

// UsecasesStaffProfile contains all the staff profile usecases
type UsecasesStaffProfile interface {
	IRegisterStaffUser
	IAddRoles
	IRemoveRole
	IUpdateDefaultFacility
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
