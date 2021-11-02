package rest

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
)

// MyCareHubHandlersInterfaces represents all the REST API logic
type MyCareHubHandlersInterfaces interface {
	//Collect metrics handler
}

// MyCareHubHandlersInterfacesImpl represents the usecase implementation object
type MyCareHubHandlersInterfacesImpl struct {
	interactor interactor.Interactor
}

// NewMyCareHubHandlersInterfaces initializes a new rest handlers usecase
func NewMyCareHubHandlersInterfaces(interactor interactor.Interactor) MyCareHubHandlersInterfaces {
	return &MyCareHubHandlersInterfacesImpl{interactor}
}
