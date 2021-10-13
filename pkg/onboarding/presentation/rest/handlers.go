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
