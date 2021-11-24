package testutils

import (
	"context"
	"fmt"

	"github.com/savannahghi/firebasetools"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
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

	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt, otpUseCase)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	contentUsecase := content.NewUsecaseContent(db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, db, externalExt)

	i := usecases.NewMyCareHubUseCase(userUsecase, termsUsecase, facilityUseCase, securityQuestionsUsecase, otpUseCase, contentUsecase)
	return i, nil
}
