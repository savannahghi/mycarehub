package rest

import (
	"fmt"
	"net/http"

	"github.com/savannahghi/errorcodeutil"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/serverutils"
)

// MyCareHubHandlersInterfaces represents all the REST API logic
type MyCareHubHandlersInterfaces interface {
	LoginByPhone() http.HandlerFunc
}

// MyCareHubHandlersInterfacesImpl represents the usecase implementation object
type MyCareHubHandlersInterfacesImpl struct {
	interactor interactor.Interactor
}

// NewMyCareHubHandlersInterfaces initializes a new rest handlers usecase
func NewMyCareHubHandlersInterfaces(interactor interactor.Interactor) MyCareHubHandlersInterfaces {
	return &MyCareHubHandlersInterfacesImpl{interactor}
}

// LoginByPhone is an unauthenticated endpoint that gets the phonenumber and pin
// from a user, checks whether they exist, if present, we fetch the pin and if they match,
// we return the user profile and auth credentials to allow the user to login
func (h *MyCareHubHandlersInterfacesImpl) LoginByPhone() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		payload := &dto.LoginInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, payload)
		if payload.PhoneNumber == nil || payload.PIN == nil {
			err := fmt.Errorf("expected `phoneNumber`, `pin` to be defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		if !payload.Flavour.IsValid() {
			err := fmt.Errorf("an invalid `flavour` defined")
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Err:     err,
				Message: err.Error(),
			}, http.StatusBadRequest)
			return
		}

		resp, responseCode, err := h.interactor.UserUsecase.Login(ctx, *payload.PhoneNumber, *payload.PIN, payload.Flavour)
		if err != nil {
			serverutils.WriteJSONResponse(w, errorcodeutil.CustomError{
				Message: err.Error(),
				Code:    responseCode,
			}, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, resp, http.StatusOK)
	}
}
