package pubsubmessaging_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	fakeFCM "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm/mock"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
	"github.com/savannahghi/scalarutils"
)

func TestServicePubSubMessaging_NotifyCreatePatient(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		client *dto.PatientCreationOutput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create patient topic",
			args: args{
				ctx: ctx,
				client: &dto.PatientCreationOutput{
					UserID:   uuid.New().String(),
					ClientID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreatePatient(tt.args.ctx, tt.args.client); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreatePatient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateVitals(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx    context.Context
		vitals *dto.PatientVitalSignOutput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create vitals topic",
			args: args{
				ctx: ctx,
				vitals: &dto.PatientVitalSignOutput{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "Vitals Test",
					ConceptID:      new(string),
					Value:          "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateVitals(tt.args.ctx, tt.args.vitals); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateVitals() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateAllergy(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx     context.Context
		allergy *dto.PatientAllergyOutput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create allergy topic",
			args: args{
				ctx: ctx,
				allergy: &dto.PatientAllergyOutput{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "Vitals Test",
					ConceptID:      new(string),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateAllergy(tt.args.ctx, tt.args.allergy); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateAllergy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateMedication(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx        context.Context
		medication *dto.PatientMedicationOutput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create medication topic",
			args: args{
				ctx: ctx,
				medication: &dto.PatientMedicationOutput{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "Vitals Test",
					ConceptID:      new(string),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateMedication(tt.args.ctx, tt.args.medication); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateMedication() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateTestOrder(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx       context.Context
		testOrder *dto.PatientTestOrderOutput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create test order topic",
			args: args{
				ctx: ctx,
				testOrder: &dto.PatientTestOrderOutput{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "Vitals Test",
					ConceptID:      new(string),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateTestOrder(tt.args.ctx, tt.args.testOrder); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateTestOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateTestResult(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx        context.Context
		testResult *dto.PatientTestResultOutput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create test result topic",
			args: args{
				ctx: ctx,
				testResult: &dto.PatientTestResultOutput{
					PatientID:      uuid.New().String(),
					OrganizationID: uuid.New().String(),
					Name:           "Vitals Test",
					ConceptID:      new(string),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateTestResult(tt.args.ctx, tt.args.testResult); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateTestResult() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateOrganization(t *testing.T) {
	ctx := context.Background()

	type args struct {
		ctx      context.Context
		facility *domain.Facility
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create organization topic",
			args: args{
				ctx: ctx,
				facility: &domain.Facility{
					ID:                 new(string),
					Name:               "Test Organization",
					Phone:              "0711111111",
					Active:             false,
					Country:            "Kenya",
					Description:        "",
					FHIROrganisationID: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateOrganization(tt.args.ctx, tt.args.facility); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateCMSClient(t *testing.T) {
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakeFCMService := fakeFCM.NewFCMServiceMock()

	ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
	type args struct {
		ctx  context.Context
		user *dto.PubsubCreateCMSClientPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create cms user topic",
			args: args{
				ctx: context.Background(),
				user: &dto.PubsubCreateCMSClientPayload{
					ClientID: uuid.New().String(),
					Name:     gofakeit.BeerAlcohol(),
					Gender:   "Male",
					DateOfBirth: scalarutils.Date{
						Year:  2022,
						Month: 10,
						Day:   02,
					},
					OrganisationID: uuid.New().String(),
					ProgramID:      uuid.New().String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ps.NotifyCreateCMSClient(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateCMSClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyDeleteCMSClient(t *testing.T) {
	fakeExtension := extensionMock.NewFakeExtension()
	fakeDB := pgMock.NewPostgresMock()
	fakeFCMService := fakeFCM.NewFCMServiceMock()

	ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)
	type args struct {
		ctx  context.Context
		user *dto.DeleteCMSUserPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to delete cms user topic",
			args: args{
				ctx: context.Background(),
				user: &dto.DeleteCMSUserPayload{
					UserID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ps.NotifyDeleteCMSClient(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyDeleteCMSClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateCMSProgram(t *testing.T) {
	type args struct {
		ctx     context.Context
		program *dto.CreateCMSProgramPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create cms program topic",
			args: args{
				ctx: context.Background(),
				program: &dto.CreateCMSProgramPayload{
					ProgramID:      uuid.New().String(),
					Name:           gofakeit.BeerAlcohol(),
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to publish to create cms program topic",
			args: args{
				ctx: context.Background(),
				program: &dto.CreateCMSProgramPayload{
					ProgramID:      uuid.New().String(),
					Name:           gofakeit.BeerAlcohol(),
					OrganisationID: uuid.New().String(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)

			if tt.name == "Sad Case - unable to publish to create cms program topic" {
				fakeExtension.MockPublishToPubsubFn = func(ctx context.Context, pubsubClient *pubsub.Client, topicID, environment, serviceName, version string, payload []byte) error {
					return fmt.Errorf("error")
				}
			}

			if err := ps.NotifyCreateCMSProgram(tt.args.ctx, tt.args.program); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateCMSProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateCMSFacility(t *testing.T) {
	type args struct {
		ctx      context.Context
		facility *dto.CreateCMSFacilityPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to create cms facility topic",
			args: args{
				ctx: context.Background(),
				facility: &dto.CreateCMSFacilityPayload{
					FacilityID: uuid.New().String(),
					Name:       gofakeit.BeerAlcohol(),
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - unable to publish to create cms facility topic",
			args: args{
				ctx: context.Background(),
				facility: &dto.CreateCMSFacilityPayload{
					FacilityID: uuid.New().String(),
					Name:       gofakeit.BeerAlcohol(),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)

			if tt.name == "Sad Case - unable to publish to create cms facility topic" {
				fakeExtension.MockPublishToPubsubFn = func(ctx context.Context, pubsubClient *pubsub.Client, topicID, environment, serviceName, version string, payload []byte) error {
					return fmt.Errorf("error")
				}
			}

			if err := ps.NotifyCreateCMSFacility(tt.args.ctx, tt.args.facility); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateCMSFacility() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCMSAddFacilityToProgram(t *testing.T) {
	type args struct {
		ctx     context.Context
		payload *dto.CMSLinkFacilityToProgramPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to add facility to program topic",
			args: args{
				ctx: nil,
				payload: &dto.CMSLinkFacilityToProgramPayload{
					FacilityID: []string{uuid.New().String()},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unable to publish to add facility to program topic",
			args: args{
				ctx: nil,
				payload: &dto.CMSLinkFacilityToProgramPayload{
					FacilityID: []string{uuid.New().String()},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)

			if tt.name == "Sad Case - Unable to publish to add facility to program topic" {
				fakeExtension.MockPublishToPubsubFn = func(ctx context.Context, pubsubClient *pubsub.Client, topicID, environment, serviceName, version string, payload []byte) error {
					return fmt.Errorf("error")
				}
			}
			if err := ps.NotifyCMSAddFacilityToProgram(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCMSAddFacilityToProgram() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyCreateClinicalTenant(t *testing.T) {
	type args struct {
		ctx    context.Context
		tenant *dto.ClinicalTenantPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy case: create tenant",
			args: args{
				ctx: context.Background(),
				tenant: &dto.ClinicalTenantPayload{
					Name:        "test",
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
					Identifiers: []dto.ClinicalTenantIdentifier{
						{
							Type:  "programID",
							Value: gofakeit.UUID(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Sad case: unable to create tenant",
			args: args{
				ctx: context.Background(),
				tenant: &dto.ClinicalTenantPayload{
					Name:        "test",
					PhoneNumber: interserviceclient.TestUserPhoneNumber,
					Identifiers: []dto.ClinicalTenantIdentifier{
						{
							Type:  "programID",
							Value: gofakeit.UUID(),
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeDB, fakeFCMService)

			if tt.name == "Sad case: unable to create tenant" {
				fakeExtension.MockPublishToPubsubFn = func(ctx context.Context, pubsubClient *pubsub.Client, topicID, environment, serviceName, version string, payload []byte) error {
					return errors.New("unable to create tenant")
				}
			}

			if err := ps.NotifyCreateClinicalTenant(tt.args.ctx, tt.args.tenant); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateClinicalTenant() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
