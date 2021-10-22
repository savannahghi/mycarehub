package rest

import (
	"net/http"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/presentation/interactor"
	"github.com/savannahghi/serverutils"
)

// OnboardingHandlersInterfaces represents all the REST API logic
type OnboardingHandlersInterfaces interface {
	//Collect metrics handler
	CollectMetricsHandler() http.HandlerFunc
	LoginHandler() http.HandlerFunc
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

// CollectMetricsHandler is an unauthenticated endpoint that is called to collect various metrics
func (h *OnboardingHandlersInterfacesImpl) CollectMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		metric := &dto.MetricInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, metric)

		response, err := h.interactor.MetricUsecase.CollectMetrics(ctx, metric)
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
	}
}

// LoginHandler is an unauthenticated endpoint that is called to login user
func (h *OnboardingHandlersInterfacesImpl) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		loginPayload := &dto.LoginInput{}
		serverutils.DecodeJSONToTargetStruct(w, r, loginPayload)

		response, _, err := h.interactor.UserUsecase.Login(ctx, loginPayload.UserID, loginPayload.PIN, loginPayload.Flavour.String())
		if err != nil {
			serverutils.WriteJSONResponse(w, err, http.StatusBadRequest)
			return
		}

		serverutils.WriteJSONResponse(w, response, http.StatusCreated)
	}
}
