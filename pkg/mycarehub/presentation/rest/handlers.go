package rest

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
)

// OnboardingHandlersInterfaces represents all the REST API logic
type OnboardingHandlersInterfaces interface {
	//Collect metrics handler
}

// OnboardingHandlersInterfacesImpl represents the usecase implementation object
type OnboardingHandlersInterfacesImpl struct {
	infrastructure infrastructure.Interactor
	interactor     interactor.Interactor
}

// NewOnboardingHandlersInterfaces initializes a new rest handlers usecase
func NewOnboardingHandlersInterfaces(infrastructure infrastructure.Interactor, interactor interactor.Interactor) OnboardingHandlersInterfaces {
	return &OnboardingHandlersInterfacesImpl{infrastructure, interactor}
}
