// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/client"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	FacilityUsecase facility.UseCasesFacility
	ClientUseCase   client.UseCasesClientProfile
	UserUseCase     user.UseCasesUser
}

// NewMyCareHubInteractor returns a new onboarding interactor
func NewMyCareHubInteractor(
	facilityUseCase facility.UseCasesFacility,
	clientUseCase client.UseCasesClientProfile,
	userUseCase user.UseCasesUser,
) *Interactor {
	return &Interactor{
		FacilityUsecase: facilityUseCase,
		ClientUseCase:   clientUseCase,
		UserUseCase:     userUseCase,
	}
}
