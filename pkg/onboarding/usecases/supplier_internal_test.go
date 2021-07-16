package usecases

import (
	"context"
	"log"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/serverutils"
	"gitlab.slade360emr.com/go/commontools/crm/pkg/infrastructure/services/hubspot"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"

	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database/fb"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/chargemaster"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/edi"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/engagement"

	erp "gitlab.slade360emr.com/go/commontools/accounting/pkg/usecases"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/messaging"
	pubsubmessaging "gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/pubsub"
)

const (
	otpService        = "otp"
	engagementService = "engagement"
	ediService        = "edi"
)

func TestParseKYCAsMap(t *testing.T) {
	ctx := context.Background()

	fc := firebasetools.FirebaseClient{}
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

	var repo repository.OnboardingRepository

	if serverutils.MustGetEnvVar(domain.Repo) == domain.FirebaseRepository {
		firestoreExtension := fb.NewFirestoreClientExtension(fsc)
		repo = fb.NewFirebaseRepository(firestoreExtension, fbc)
	}
	projectID, err := serverutils.GetEnvVar(serverutils.GoogleCloudProjectIDEnvVarName)
	if err != nil {
		t.Errorf("can't get projectID from env var `%s`: %w",
			serverutils.GoogleCloudProjectIDEnvVarName,
			err)
		return
	}
	pubSubClient, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		t.Errorf("unable to initialize pubsub client: %w", err)
		return
	}

	ext := extension.NewBaseExtensionImpl(&firebasetools.FirebaseClient{})
	// Initialize ISC clients
	engagementClient := utils.NewInterServiceClient(engagementService, ext)
	ediClient := utils.NewInterServiceClient(ediService, ext)
	edi := edi.NewEdiService(ediClient, repo)
	erp := erp.NewAccounting()
	chrg := chargemaster.NewChargeMasterUseCasesImpl()
	crm := hubspot.NewHubSpotService()
	ps, err := pubsubmessaging.NewServicePubSubMessaging(
		pubSubClient,
		ext,
		erp,
		crm,
		edi,
		repo,
	)
	if err != nil {
		t.Errorf("unable to initialize new pubsub messaging service: %w", err)
		return
	}
	engage := engagement.NewServiceEngagementImpl(engagementClient, ext, ps)
	mes := messaging.NewServiceMessagingImpl(ext)
	profile := NewProfileUseCase(repo, ext, engage, ps)

	supplier := SupplierUseCasesImpl{
		repo:         repo,
		profile:      profile,
		erp:          erp,
		chargemaster: chrg,
		engagement:   engage,
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
					IdentificationDocType:           enumutils.IdentificationDocTypeMilitary,
					IdentificationDocNumber:         "11111111",
					IdentificationDocNumberUploadID: "11111111",
				},
				KRAPIN:                         "krapin",
				KRAPINUploadID:                 "krapinuploadID",
				DrivingLicenseID:               "dlid",
				DrivingLicenseUploadID:         "dliduploadid",
				CertificateGoodConductUploadID: "cert",
				SupportingDocuments: []domain.SupportingDocument{
					{
						SupportingDocumentTitle:       "support-title",
						SupportingDocumentDescription: "support-description",
						SupportingDocumentUpload:      "support-upload-id",
					},
				},
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

			supportingDocumentsResp, ok := response["supportingDocuments"]
			if !ok {
				t.Errorf("supportingDocuments is nil")
				return
			}

			if ok {
				supportingDocuments := supportingDocumentsResp.([]interface{})
				if len(data.SupportingDocuments) != len(supportingDocuments) {
					t.Errorf("wanted: %v, got: %v", len(data.SupportingDocuments), len(supportingDocuments))
					return
				}
			}

		})
	}

}
