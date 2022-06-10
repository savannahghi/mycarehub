package pubsubmessaging_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	fakeFCM "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/fcm/mock"
	getStreamMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/getstream/mock"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
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
					ID:     uuid.New().String(),
					UserID: uuid.New().String(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
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
					Code:               0,
					Phone:              "0711111111",
					Active:             false,
					County:             "Nairobi",
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
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
			if err := ps.NotifyCreateOrganization(tt.args.ctx, tt.args.facility); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyCreateOrganization() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestServicePubSubMessaging_NotifyGetStreamEvent(t *testing.T) {
	ctx := context.Background()
	type args struct {
		ctx   context.Context
		event *dto.GetStreamEvent
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully publish to getstream event topic",
			args: args{
				ctx: ctx,
				event: &dto.GetStreamEvent{
					CID:  uuid.New().String(),
					Type: "messaging.new",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeExtension := extensionMock.NewFakeExtension()
			fakeGetStream := getStreamMock.NewGetStreamServiceMock()
			fakeDB := pgMock.NewPostgresMock()
			fakeFCMService := fakeFCM.NewFCMServiceMock()

			ps, _ := pubsubmessaging.NewServicePubSubMessaging(fakeExtension, fakeGetStream, fakeDB, fakeFCMService)
			if err := ps.NotifyGetStreamEvent(tt.args.ctx, tt.args.event); (err != nil) != tt.wantErr {
				t.Errorf("ServicePubSubMessaging.NotifyGetStreamEvent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
