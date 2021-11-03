package usecase

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type MyCareHubUseCase interface {
	facility.UseCasesFacility
	// client.UseCasesClientProfile
	// user.UseCasesUser
}

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type MyCareHub struct {
	Facility facility.UseCasesFacility
	// Client   client.UseCasesClientProfile
	// User     user.UseCasesUser
}

// NewMyCareHubInteractor returns a new onboarding interactor
func NewMyCareHubUseCase(
	facilityUseCase facility.UseCasesFacility,
	// clientUseCase client.UseCasesClientProfile,
	// userUseCase user.UseCasesUser,
) *MyCareHub {
	return &MyCareHub{
		Facility: facilityUseCase,
		// Client:   clientUseCase,
		// User:     userUseCase,
	}
}
