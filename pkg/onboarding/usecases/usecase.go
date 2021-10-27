package usecases

import (
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/extension"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/client"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/metric"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/staff"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/user"
	engagementSvc "github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/services/engagement"
)

// Interactor is an implementation of the usecases interface
type Interactor struct {
	*engagementSvc.ServiceEngagementImpl
	*facility.UseCaseFacilityImpl
	*metric.UsecaseMetricsImpl
	*user.UseCasesUserImpl
	*client.UseCasesClientImpl
	*staff.UsecasesStaffProfileImpl
}

// NewUsecasesInteractor initializes a new usecases interactor
func NewUsecasesInteractor(infrastructure infrastructure.Interactor) Interactor {
	var engagement *engagementSvc.ServiceEngagementImpl
	onboardingExt := extension.NewOnboardingLibImpl()
	facility := facility.NewFacilityUsecase(infrastructure)
	metrics := metric.NewMetricUsecase(infrastructure)
	user := user.NewUseCasesUserImpl(infrastructure, onboardingExt, engagement)
	client := client.NewUseCasesClientImpl(infrastructure)
	staff := staff.NewUsecasesStaffProfileImpl(infrastructure)

	return Interactor{
		engagement,
		facility,
		metrics,
		user,
		client,
		staff,
	}
}
