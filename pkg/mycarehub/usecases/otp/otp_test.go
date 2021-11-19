package otp_test

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	externalExtension "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/presentation/interactor"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
	"github.com/savannahghi/profileutils"
	"github.com/savannahghi/serverutils"
	"github.com/segmentio/ksuid"
)

var termsID = 50005

func InitializeTestService(ctx context.Context) *interactor.Interactor {
	pg, err := gorm.NewPGInstance()
	if err != nil {
		return nil
	}

	db := postgres.NewMyCareHubDb(pg, pg, pg, pg)
	externalExt := externalExtension.NewExternalMethodsImpl()

	// Initialize facility usecase
	facilityUseCase := facility.NewFacilityUsecase(db, db, db, db)

	// Initialize user usecase
	userUsecase := user.NewUseCasesUserImpl(db, db, db, db, externalExt)

	termsUsecase := terms.NewUseCasesTermsOfService(db, db)

	securityQuestionsUsecase := securityquestions.NewSecurityQuestionsUsecase(db, db, externalExt)

	otpUseCase := otp.NewOTPUseCase(db, db, externalExt)

	i := interactor.NewMyCareHubInteractor(facilityUseCase, userUsecase, termsUsecase, securityQuestionsUsecase, otpUseCase)
	return i
}
func TestUseCaseOTPImpl_VerifyPhoneNumber_Integration(t *testing.T) {
	ctx := context.Background()

	i := InitializeTestService(ctx)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	flavour := feedlib.FlavourConsumer
	userID := uuid.New().String()
	inexistentNo := "+254700000520"
	invalidPhone := ksuid.New().String()

	// Setup test user
	userInput := &gorm.User{
		UserID:          &userID,
		Username:        uuid.New().String(),
		FirstName:       gofakeit.FirstName(),
		LastName:        gofakeit.LastName(),
		MiddleName:      gofakeit.FirstName(),
		UserType:        enums.ClientUser,
		Gender:          enumutils.GenderMale,
		Flavour:         flavour,
		AcceptedTermsID: &termsID,
		TermsAccepted:   true,
		OrganisationID:  serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
		Suspended:       false,
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	contactID := uuid.New().String()
	contact := &gorm.Contact{
		ContactID:      &contactID,
		ContactType:    "SMS",
		ContactValue:   "+254710000100",
		Active:         true,
		OptedIn:        true,
		UserID:         userInput.UserID,
		Flavour:        userInput.Flavour,
		OrganisationID: serverutils.MustGetEnvVar("DEFAULT_ORG_ID"),
	}

	err = pg.DB.Create(&contact).Error
	if err != nil {
		t.Errorf("failed to create contact: %v", err)
	}

	type args struct {
		ctx     context.Context
		phone   string
		flavour feedlib.Flavour
	}
	tests := []struct {
		name    string
		args    args
		want    *profileutils.OtpResponse
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				phone:   contact.ContactValue,
				flavour: contact.Flavour,
			},
			wantErr: false,
		},
		{
			name: "Sad case - invalid flavour",
			args: args{
				ctx:     ctx,
				phone:   contact.ContactValue,
				flavour: "contact.Flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - inexistent number and good flavour",
			args: args{
				ctx:     ctx,
				phone:   inexistentNo,
				flavour: contact.Flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - inexistent number and bad flavour",
			args: args{
				ctx:     ctx,
				phone:   inexistentNo,
				flavour: "gofakeit.HipsterSentence(100)",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid phone and valid flavour",
			args: args{
				ctx:     ctx,
				phone:   invalidPhone,
				flavour: contact.Flavour,
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid phone and invalid flavour",
			args: args{
				ctx:     ctx,
				phone:   invalidPhone,
				flavour: "contact.Flavour",
			},
			wantErr: true,
		},
		{
			name: "Sad case - invalid phone and empty flavour",
			args: args{
				ctx:     ctx,
				phone:   invalidPhone,
				flavour: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.OTPUsecase.VerifyPhoneNumber(tt.args.ctx, tt.args.phone, tt.args.flavour)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCaseOTPImpl.VerifyPhoneNumber() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", contact.ContactID).Unscoped().Delete(&gorm.Contact{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}
