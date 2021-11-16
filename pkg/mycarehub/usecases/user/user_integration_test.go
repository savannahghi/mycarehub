package user_test

import (
	"context"

	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
)

func InitializeTestService(ctx context.Context) *interactor.Interactor {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil
	}

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)
	externalExt := externalExtension.NewExternalMethodsImpl()

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db)

	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, externalExt)

	otpUseCase := otp.NewOTPUseCase(db, db, externalExt)

	i := interactor.NewMyCareHubInteractor(facilityUseCase, userUsecase, termsUsecase, securityQuestionsUsecase, otpUseCase)
	return i
}
