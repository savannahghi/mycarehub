package testutils

import (
	"context"
	"fmt"

	"github.com/savannahghi/firebasetools"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	streamService "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	appointment "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/appointments"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/faq"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/screeningtools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
)

// InitializeTestService sets up the structure that will be used by the usecase layer for
// integration tests
func InitializeTestService(ctx context.Context) (*usecases.MyCareHub, error) {
	fc := &firebasetools.FirebaseClient{}
	_, err := fc.InitFirebase()
	if err != nil {
		return nil, err
	}

	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil, fmt.Errorf("can't instantiate test repository: %v", err)
	}

	// add organization
	// createOrganization(pg)

	externalExt := externalExtension.NewExternalMethodsImpl()

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db)

	otpUseCase := otp.NewOTPUseCase(db, db, externalExt)
	getStream := streamService.NewServiceGetStream()
	authorityUseCase := authority.NewUsecaseAuthority(db, db, externalExt)

	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt, otpUseCase, authorityUseCase, getStream)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, db, externalExt)
	contentUseCase := content.NewUseCasesContentImplementation(db, db, externalExt)
	feedbackUsecase := feedback.NewUsecaseFeedback(db, externalExt)

	faq := faq.NewUsecaseFAQ(db)
	serviceRequestUseCase := servicerequest.NewUseCaseServiceRequestImpl(db, db, db, externalExt, userUsecase)
	healthDiaryUseCase := healthdiary.NewUseCaseHealthDiaryImpl(db, db, serviceRequestUseCase)
	appointmentUsecase := appointment.NewUseCaseAppointmentsImpl(externalExt, db, db, db)
	communityUsecase := communities.NewUseCaseCommunitiesImpl(getStream, externalExt, db, db)

	screeningToolsUsecases := screeningtools.NewUseCasesScreeningTools(db, db, db)

	i := usecases.NewMyCareHubUseCase(
		userUsecase, termsUsecase, facilityUseCase,
		securityQuestionsUsecase, otpUseCase, contentUseCase, feedbackUsecase, healthDiaryUseCase,
		faq, serviceRequestUseCase, authorityUseCase, communityUsecase, screeningToolsUsecases,
		appointmentUsecase,
	)
	return i, nil
}
