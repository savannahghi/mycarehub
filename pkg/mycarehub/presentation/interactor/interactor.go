// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	FacilityUsecase facility.UseCasesFacility
	UserUsecase     user.UseCasesUser
}

// NewMyCareHubInteractor returns a new onboarding interactor
func NewMyCareHubInteractor(
	facilityUseCase facility.UseCasesFacility,
	userUseCase user.UseCasesUser,
) *Interactor {
	return &Interactor{
		facilityUseCase,
		userUseCase,
	}
}
