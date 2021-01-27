package usecases

import (
	"context"
	"log"
	"testing"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/erp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/mailgun"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
)

const (
	otpService        = "otp"
	mailgunService    = "mailgun"
	engagementService = "engagement"
)

func TestParseKYCAsMap(t *testing.T) {
	ctx := context.Background()

	fc := base.FirebaseClient{}
	fa, err := fc.InitFirebase()
	if err != nil {
		log.Fatalf("unable to initialize Firestore for the Feed: %s", err)
	}

	fsc, err := fa.Firestore(ctx)
	if err != nil {
		log.Fatalf("unable to initialize Firestore: %s", err)
	}

	fbc, err := fa.Auth(ctx)
	if err != nil {
		log.Panicf("can't initialize Firebase auth when setting up profile service: %s", err)
	}
	firestoreExtension := database.NewFirestoreClientExtension(fsc)
	fr := database.NewFirebaseRepository(firestoreExtension, fbc)
	if err != nil {
		return
	}
	// Initialize ISC clients
	otpClient := utils.NewInterServiceClient(otpService)
	mailgunClient := utils.NewInterServiceClient(mailgunService)
	engagementClient := utils.NewInterServiceClient(engagementService)

	erp := erp.NewERPService()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	engage := engagement.NewServiceEngagementImpl(engagementClient)
	mg := mailgun.NewServiceMailgunImpl(mailgunClient)
	mes := messaging.NewServiceMessagingImpl()
	ext := extension.NewBaseExtensionImpl()
	otp := otp.NewOTPService(otpClient, ext)
	profile := NewProfileUseCase(fr, otp, ext, engage)

	supplier := SupplierUseCasesImpl{
		repo:         fr,
		profile:      profile,
		erp:          erp,
		chargemaster: chrg,
		engagement:   engage,
		mg:           mg,
		messaging:    mes,
	}

	tests := []struct {
		name string
	}{
		{
			name: "valid case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			data := domain.IndividualRider{
				IdentificationDoc: domain.Identification{
					IdentificationDocType:           domain.IdentificationDocTypeMilitary,
					IdentificationDocNumber:         "11111111",
					IdentificationDocNumberUploadID: "11111111",
				},
				KRAPIN:                         "krapin",
				KRAPINUploadID:                 "krapinuploadID",
				DrivingLicenseID:               "dlid",
				DrivingLicenseUploadID:         "dliduploadid",
				CertificateGoodConductUploadID: "cert",
				SupportingDocumentsUploadID:    []string{"someID", "anotherID"},
			}

			response, err := supplier.parseKYCAsMap(data)
			if err != nil {
				t.Errorf("failed to parse data, returned error: %v", err)
				return
			}

			identificationDoc, ok := response["identificationDoc"]
			if !ok {
				t.Errorf("identificationDoc is nil")
				return
			}
			if ok {
				identificationDoc := identificationDoc.(map[string]interface{})
				_, ok := identificationDoc["identificationDocType"]
				if !ok {
					t.Errorf("identificationDoc['identificationDocType'] is nil")
					return
				}
				_, ok = identificationDoc["identificationDocNumber"]
				if !ok {
					t.Errorf("identificationDoc['identificationDocNumber'] is nil")
					return
				}

				_, ok = identificationDoc["identificationDocNumberUploadID"]
				if !ok {
					t.Errorf("identificationDoc['identificationDocNumberUploadID'] is nil")
					return
				}
			}

			_, ok = response["KRAPIN"]
			if !ok {
				t.Errorf("KRAPIN is nil")
				return
			}

			_, ok = response["KRAPINUploadID"]
			if !ok {
				t.Errorf("KRAPINUploadID is nil")
				return
			}
			_, ok = response["drivingLicenseID"]
			if !ok {
				t.Errorf("drivingLicenseID is nil")
				return
			}

			_, ok = response["drivingLicenseUploadID"]
			if !ok {
				t.Errorf("drivingLicenseUploadID is nil")
				return
			}
			_, ok = response["certificateGoodConductUploadID"]
			if !ok {
				t.Errorf("certificateGoodConductUploadID is nil")
				return
			}

			supportingDocumentsUploadID, ok := response["supportingDocumentsUploadID"]
			if !ok {
				t.Errorf("supportingDocumentsUploadID is nil")
				return
			}

			if ok {
				supportingDocumentsUploadID := supportingDocumentsUploadID.([]interface{})
				if len(data.SupportingDocumentsUploadID) != len(supportingDocumentsUploadID) {
					t.Errorf("wanted: %v, got: %v", len(data.SupportingDocumentsUploadID), len(supportingDocumentsUploadID))
					return
				}
			}

		})
	}

}
