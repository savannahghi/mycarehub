package appointment

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	pubsubMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub/mock"
	notificationMock "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification/mock"
	"github.com/savannahghi/scalarutils"
	"gorm.io/gorm"
)

func TestUseCasesAppointmentsImpl_CreateKenyaEMRAppointments(t *testing.T) {

	type args struct {
		ctx   context.Context
		input dto.FacilityAppointmentsPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.FacilityAppointmentsResponse
		wantErr bool
	}{
		{
			name: "sad case: error checking facility",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: facility doesn't exist",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving facility",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving client profile",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error creating appointment",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "happy case: success creating appointment",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to notify user",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: error checking facility" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error retrieving facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "sad case: error retrieving client profile" {
				fakeDB.MockCheckAppointmentExistsByExternalIDFn = func(ctx context.Context, externalID string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("failed to retrieve client profile")
				}
			}

			if tt.name == "sad case: error creating appointment" {
				fakeDB.MockCheckAppointmentExistsByExternalIDFn = func(ctx context.Context, externalID string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{
						ID: &id,
						User: &domain.User{
							CurrentProgramID: id,
						},
					}, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCreateAppointment = func(ctx context.Context, appointment domain.Appointment) error {
					return fmt.Errorf("cannot create appointment")
				}
			}

			if tt.name == "happy case: success creating appointment" {
				fakeDB.MockCheckAppointmentExistsByExternalIDFn = func(ctx context.Context, externalID string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{
						ID: &id,
						User: &domain.User{
							CurrentProgramID: id,
						},
					}, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCreateAppointment = func(ctx context.Context, appointment domain.Appointment) error {
					return nil
				}
			}

			if tt.name == "sad case: fail to notify user" {
				fakeDB.MockCheckAppointmentExistsByExternalIDFn = func(ctx context.Context, externalID string) (bool, error) {
					return false, nil
				}

				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					id := gofakeit.UUID()
					return &domain.ClientProfile{
						ID: &id,
						User: &domain.User{
							CurrentProgramID: id,
						},
					}, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCreateAppointment = func(ctx context.Context, appointment domain.Appointment) error {
					return nil
				}

				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify patient")
				}
			}

			got, err := a.CreateOrUpdateKenyaEMRAppointments(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.CreateKenyaEMRAppointments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected appointments not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_UpdateKenyaEMRAppointments(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()
	fakePubsub := pubsubMock.NewPubsubServiceMock()
	fakeNotification := notificationMock.NewServiceNotificationMock()

	a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

	type args struct {
		ctx   context.Context
		input dto.FacilityAppointmentsPayload
	}
	tests := []struct {
		name    string
		args    args
		want    *dto.FacilityAppointmentsResponse
		wantErr bool
	}{
		{
			name: "sad case: error checking facility",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: facility doesn't exist",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving facility",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving client profile",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error updating appointment",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "happy case: success updating appointment",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: fail to notify user",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: failed to get appointment by UUID",
			args: args{
				ctx: context.Background(),
				input: dto.FacilityAppointmentsPayload{
					MFLCode: "1234",
					Appointments: []dto.AppointmentPayload{
						{
							CCCNumber:         "1234",
							ExternalID:        gofakeit.UUID(),
							AppointmentReason: "Dental",
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: error checking facility" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error retrieving facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "sad case: error retrieving client profile" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}
				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("cannot retrieve client profile")
				}
			}

			if tt.name == "sad case: error updating appointment" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{}, nil
				}

				fakeDB.MockUpdateAppointmentFn = func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
					return nil, fmt.Errorf("error updating appointment")
				}
			}

			if tt.name == "happy case: success updating appointment" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{}, nil
				}

				fakeDB.MockUpdateAppointmentFn = func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
					return appointment, nil
				}

			}

			if tt.name == "sad case: fail to notify user" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{}, nil
				}

				fakeDB.MockUpdateAppointmentFn = func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
					return appointment, nil
				}

				fakeNotification.MockNotifyUserFn = func(ctx context.Context, userProfile *domain.User, notificationPayload *domain.Notification) error {
					return fmt.Errorf("failed to notify user")
				}
			}
			if tt.name == "sad case: failed to get appointment by UUID" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{}, nil
				}

				fakeDB.MockUpdateAppointmentFn = func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
					return &domain.Appointment{
						ID:                        "",
						ExternalID:                "uuid",
						Reason:                    "",
						Date:                      scalarutils.Date{},
						ClientID:                  "",
						FacilityID:                "",
						Provider:                  "",
						HasRescheduledAppointment: false,
					}, nil
				}

				fakeDB.MockGetAppointmentFn = func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
					return nil, fmt.Errorf("error retrieving appointment by UUID")
				}
			}

			got, err := a.CreateOrUpdateKenyaEMRAppointments(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.UpdateKenyaEMRAppointments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected appointments not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_FetchClientAppointments(t *testing.T) {
	type args struct {
		ctx             context.Context
		clientID        string
		paginationInput dto.PaginationsInput
		filters         []*firebasetools.FilterParam
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.AppointmentsPage
		wantErr bool
	}{
		{
			name: "sad case: current page not provided",
			args: args{
				ctx:             context.Background(),
				clientID:        "client-id",
				paginationInput: dto.PaginationsInput{},
				filters:         []*firebasetools.FilterParam{},
			},
			wantErr: true,
		},
		{
			name: "sad case: error listing appointments",
			args: args{
				ctx:      context.Background(),
				clientID: "client-id",
				paginationInput: dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       5,
				},
				filters: []*firebasetools.FilterParam{},
			},
			wantErr: true,
		},
		{
			name: "happy case: success listing appointments",
			args: args{
				ctx:      context.Background(),
				clientID: "client-id",
				paginationInput: dto.PaginationsInput{
					CurrentPage: 1,
					Limit:       5,
				},
				filters: []*firebasetools.FilterParam{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: error listing appointments" {
				fakeDB.MockListAppointments = func(ctx context.Context, params *domain.Appointment, filters []*firebasetools.FilterParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("error listing appointments")
				}
			}
			got, err := a.FetchClientAppointments(tt.args.ctx, tt.args.clientID, tt.args.paginationInput, tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.FetchClientAppointments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected appointments to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("expected appointments not to be nil for %v", tt.name)
				return
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_AddPatientsRecords(t *testing.T) {
	conceptID := gofakeit.UUID()

	type args struct {
		ctx   context.Context
		input dto.PatientsRecordsPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad case: invalid mfl code",
			args: args{
				ctx: context.Background(),
				input: dto.PatientsRecordsPayload{
					MFLCode: "invalid",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error checking facility exist",
			args: args{
				ctx: context.Background(),
				input: dto.PatientsRecordsPayload{
					MFLCode: "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: facility doesn't exist",
			args: args{
				ctx: context.Background(),
				input: dto.PatientsRecordsPayload{
					MFLCode: "1234",
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving facility",
			args: args{
				ctx: context.Background(),
				input: dto.PatientsRecordsPayload{
					MFLCode: "1234",
					Records: []dto.PatientRecordPayload{
						{
							CCCNumber: "1234",
						},
						{
							CCCNumber: "12345",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "happy case: success add observations",
			args: args{
				ctx: context.Background(),
				input: dto.PatientsRecordsPayload{
					MFLCode: "1234",
					Records: []dto.PatientRecordPayload{
						{
							CCCNumber: "1234",
							MFLCode:   1234,
							VitalSigns: []*dto.VitalSignPayload{
								{
									Name:      "Vitals",
									ConceptID: &conceptID,
									Value:     "23",
									Date:      time.Now(),
								},
							},
							TestOrders: []*dto.TestOrderPayload{
								{
									Name:      "Test order",
									ConceptID: &conceptID,
									Date:      time.Now(),
								},
							},
							TestResults: []*dto.TestResultPayload{
								{
									Name:            "Result",
									TestConceptID:   &conceptID,
									Date:            time.Now(),
									Result:          "Good",
									ResultConceptID: &conceptID,
								},
							},
							Medications: []*dto.MedicationPayload{
								{
									Name:                "Medication",
									MedicationConceptID: &conceptID,
									Date:                time.Now(),
									Value:               "ARV",
									DrugConceptID:       &conceptID,
								},
							},
							Allergies: []*dto.AllergyPayload{
								{
									Name:              "Allergy",
									AllergyConceptID:  &conceptID,
									Reaction:          "Bad",
									ReactionConceptID: &conceptID,
									Severity:          "High",
									SeverityConceptID: &conceptID,
									Date:              time.Now(),
								},
							},
						},
						{
							CCCNumber: "12345",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: error checking facility exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error checking facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("error retrieving facility")
				}
			}

			if err := a.AddPatientsRecords(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.AddPatientsRecords() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_AddPatientRecord(t *testing.T) {
	conceptID := gofakeit.UUID()

	labTestConceptID := labTestConceptID
	counsellingConceptID := counsellingConceptID
	pharmacyRefillConceptID := pharmacyRefillConceptID
	otherConceptID := otherConceptID
	followUpConceptID := followUpConceptID
	returnVisitConceptID := returnVisitConceptID

	type args struct {
		ctx   context.Context
		input dto.PatientRecordPayload
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sad case: missing ccc number",
			args: args{
				ctx:   context.Background(),
				input: dto.PatientRecordPayload{CCCNumber: ""},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving mfl code",
			args: args{
				ctx:   context.Background(),
				input: dto.PatientRecordPayload{CCCNumber: "1234", MFLCode: 1234},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving client profile",
			args: args{
				ctx:   context.Background(),
				input: dto.PatientRecordPayload{CCCNumber: "1234", MFLCode: 1234},
			},
			wantErr: true,
		},
		{
			name: "sad case: missing fhir patient id",
			args: args{
				ctx:   context.Background(),
				input: dto.PatientRecordPayload{CCCNumber: "1234", MFLCode: 1234},
			},
			wantErr: true,
		},
		{
			name: "sad case: error publishing to vitals topic",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					VitalSigns: []*dto.VitalSignPayload{
						{
							Name:      "Vitals",
							ConceptID: &conceptID,
							Value:     "23",
							Date:      time.Now(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error publishing to allergy topic",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					Allergies: []*dto.AllergyPayload{
						{
							Name:              "Allergy",
							AllergyConceptID:  &conceptID,
							Reaction:          "Bad",
							ReactionConceptID: &conceptID,
							Severity:          "High",
							SeverityConceptID: &conceptID,
							Date:              time.Now(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error publishing to medication topic",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					Medications: []*dto.MedicationPayload{
						{
							Name:                "Medication",
							MedicationConceptID: &conceptID,
							Date:                time.Now(),
							Value:               "ARV",
							DrugConceptID:       &conceptID,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error publishing to result topic",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					TestResults: []*dto.TestResultPayload{
						{
							Name:            "Result",
							TestConceptID:   &conceptID,
							Date:            time.Now(),
							Result:          "Good",
							ResultConceptID: &conceptID,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error publishing to order topic",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					TestOrders: []*dto.TestOrderPayload{
						{
							Name:      "Test order",
							ConceptID: &conceptID,
							Date:      time.Now(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: successfully update patient observations",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					VitalSigns: []*dto.VitalSignPayload{
						{
							Name:      "Vitals",
							ConceptID: &conceptID,
							Value:     "23",
							Date:      time.Now(),
						},
					},
					TestOrders: []*dto.TestOrderPayload{
						{
							Name:      "Test order",
							ConceptID: &conceptID,
							Date:      time.Now(),
						},
					},
					TestResults: []*dto.TestResultPayload{
						{
							Name:            "Result",
							TestConceptID:   &conceptID,
							Date:            time.Now(),
							Result:          "Good",
							ResultConceptID: &conceptID,
						},
					},
					Medications: []*dto.MedicationPayload{
						{
							Name:                "Medication",
							MedicationConceptID: &conceptID,
							Date:                time.Now(),
							Value:               "ARV",
							DrugConceptID:       &conceptID,
						},
					},
					Allergies: []*dto.AllergyPayload{
						{
							Name:              "Allergy",
							AllergyConceptID:  &conceptID,
							Reaction:          "Bad",
							ReactionConceptID: &conceptID,
							Severity:          "High",
							SeverityConceptID: &conceptID,
							Date:              time.Now(),
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "happy case: skip appointment vitals",
			args: args{
				ctx: context.Background(),
				input: dto.PatientRecordPayload{
					CCCNumber: "1234",
					MFLCode:   1234,
					VitalSigns: []*dto.VitalSignPayload{
						{
							Name:      "Vitals",
							ConceptID: &labTestConceptID,
							Value:     "value",
							Date:      time.Now(),
						},
						{
							Name:      "Vitals",
							ConceptID: &counsellingConceptID,
							Value:     "value",
							Date:      time.Now(),
						},
						{
							Name:      "Vitals",
							ConceptID: &pharmacyRefillConceptID,
							Value:     "value",
							Date:      time.Now(),
						},
						{
							Name:      "Vitals",
							ConceptID: &otherConceptID,
							Value:     "value",
							Date:      time.Now(),
						},
						{
							Name:      "Vitals",
							ConceptID: &returnVisitConceptID,
							Value:     "value",
							Date:      time.Now(),
						},
						{
							Name:      "Vitals",
							ConceptID: &followUpConceptID,
							Value:     "value",
							Date:      time.Now(),
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: error retrieving mfl code" {
				fakeDB.MockRetrieveFacilityByIdentifierFn = func(ctx context.Context, identifier *dto.FacilityIdentifierInput, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "sad case: error retrieving client profile" {
				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error retrieving client by mfl code")
				}
			}

			if tt.name == "sad case: missing fhir patient id" {
				fakeDB.MockGetClientProfileByCCCNumberFn = func(ctx context.Context, CCCNumber string) (*domain.ClientProfile, error) {
					return &domain.ClientProfile{FHIRPatientID: nil}, nil
				}
			}

			if tt.name == "sad case: error publishing to vitals topic" {
				fakePubsub.MockNotifyCreateVitalsFn = func(ctx context.Context, vitals *dto.PatientVitalSignOutput) error {
					return fmt.Errorf("error notifying topic")
				}
			}

			if tt.name == "sad case: error publishing to allergy topic" {
				fakePubsub.MockNotifyCreateAllergyFn = func(ctx context.Context, allergy *dto.PatientAllergyOutput) error {
					return fmt.Errorf("error notifying topic")
				}
			}

			if tt.name == "sad case: error publishing to medication topic" {
				fakePubsub.MockNotifyCreateMedicationFn = func(ctx context.Context, medication *dto.PatientMedicationOutput) error {
					return fmt.Errorf("error notifying topic")
				}
			}

			if tt.name == "sad case: error publishing to result topic" {
				fakePubsub.MockNotifyCreateTestResultFn = func(ctx context.Context, testResult *dto.PatientTestResultOutput) error {
					return fmt.Errorf("error notifying topic")
				}
			}

			if tt.name == "sad case: error publishing to order topic" {
				fakePubsub.MockNotifyCreateTestOrderFn = func(ctx context.Context, testOrder *dto.PatientTestOrderOutput) error {
					return fmt.Errorf("error notifying topic")
				}
			}

			if err := a.AddPatientRecord(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.AddPatientRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_GetAppointmentServiceRequests(t *testing.T) {
	now := time.Now()
	type args struct {
		ctx     context.Context
		payload dto.AppointmentServiceRequestInput
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx: context.Background(),
				payload: dto.AppointmentServiceRequestInput{
					LastSyncTime: &now,
					MFLCode:      123,
				},
			},
			wantErr: false,
		},
		{
			name: "sad case: error facility with provided mfl code not found",
			args: args{
				ctx: context.Background(),
				payload: dto.AppointmentServiceRequestInput{
					LastSyncTime: &now,
					MFLCode:      123,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving facility by mfl code",
			args: args{
				ctx: context.Background(),
				payload: dto.AppointmentServiceRequestInput{
					LastSyncTime: &now,
					MFLCode:      123,
				},
			},
			wantErr: true,
		},
		{
			name: "sad case: error retrieving appointment service requests",
			args: args{
				ctx: context.Background(),
				payload: dto.AppointmentServiceRequestInput{
					LastSyncTime: &now,
					MFLCode:      123,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: error facility with provided mfl code not found" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, nil
				}

			}
			if tt.name == "sad case: error retrieving facility by mfl code" {
				fakeDB.MockCheckFacilityExistsByIdentifier = func(ctx context.Context, identifier *dto.FacilityIdentifierInput) (bool, error) {
					return false, fmt.Errorf("error retrieving facility by mfl code")
				}
			}

			if tt.name == "sad case: error retrieving appointment service requests" {
				fakeDB.MockGetAppointmentServiceRequestsFn = func(ctx context.Context, lastSyncTime time.Time, mflCode string) ([]domain.AppointmentServiceRequests, error) {
					return nil, fmt.Errorf("error retrieving appointment service requests")
				}
			}
			got, err := a.GetAppointmentServiceRequests(tt.args.ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.GetAppointmentServiceRequests() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Errorf("UseCasesAppointmentsImpl.GetAppointmentServiceRequests() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_RescheduleClientAppointment(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	futureDate, err := scalarutils.NewDate(futureTime.Day(), int(futureTime.Month()), futureTime.Year())
	if err != nil {
		t.Errorf("unable to create future date error: %v", err)
		return
	}

	type args struct {
		ctx           context.Context
		appointmentID string
		date          scalarutils.Date
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "happy case",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
				date:          *futureDate,
			},
			wantErr: false,
			want:    true,
		},
		{
			name: "sad case: empty appointment id",
			args: args{
				ctx:           context.Background(),
				appointmentID: "",
				date:          *futureDate,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "sad case: failed to get client by id",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
				date:          *futureDate,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "sad case: failed to get appointment by id",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
				date:          *futureDate,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "sad case: failed to create service request",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
				date:          *futureDate,
			},
			wantErr: true,
			want:    false,
		},
		{
			name: "sad case: failed to update appointment",
			args: args{
				ctx:           context.Background(),
				appointmentID: uuid.New().String(),
				date:          *futureDate,
			},
			wantErr: true,
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()
			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)

			if tt.name == "sad case: failed to get client by id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("error retrieving client by id")
				}
			}

			if tt.name == "sad case: failed to get appointment by id" {
				fakeDB.MockGetAppointmentFn = func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
					return nil, fmt.Errorf("error retrieving appointment by id")
				}
			}
			if tt.name == "sad case: failed to create service request" {
				fakeDB.MockCreateServiceRequestFn = func(ctx context.Context, serviceRequestInput *dto.ServiceRequestInput) error {
					return fmt.Errorf("failed to create service request")
				}
			}
			if tt.name == "sad case: failed to update appointment" {
				fakeDB.MockUpdateAppointmentFn = func(ctx context.Context, appointment *domain.Appointment, updateData map[string]interface{}) (*domain.Appointment, error) {
					return nil, fmt.Errorf("error updating appointment")
				}
			}

			got, err := a.RescheduleClientAppointment(tt.args.ctx, tt.args.appointmentID, tt.args.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.RescheduleClientAppointment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UseCasesAppointmentsImpl.RescheduleClientAppointment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUseCasesAppointmentsImpl_NextRefill(t *testing.T) {

	type args struct {
		ctx      context.Context
		clientID string
	}
	tests := []struct {
		name    string
		args    args
		want    *scalarutils.Date
		wantErr bool
		wantNil bool
	}{
		{
			name: "Happy case: has next refill date",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "Happy case: has no refill date",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: false,
			wantNil: true,
		},
		{
			name: "Sad case: invalid client id",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
			wantNil: true,
		},
		{
			name: "Sad case: error fetching appointment",
			args: args{
				ctx:      context.Background(),
				clientID: gofakeit.UUID(),
			},
			wantErr: true,
			wantNil: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeDB := pgMock.NewPostgresMock()
			fakeExtension := extensionMock.NewFakeExtension()
			fakePubsub := pubsubMock.NewPubsubServiceMock()

			if tt.name == "Happy case: has no refill date" {
				fakeDB.MockGetAppointmentFn = func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
					return nil, gorm.ErrRecordNotFound
				}
			}

			if tt.name == "Sad case: invalid client id" {
				fakeDB.MockGetClientProfileByClientIDFn = func(ctx context.Context, clientID string) (*domain.ClientProfile, error) {
					return nil, fmt.Errorf("client does not exist")
				}
			}

			if tt.name == "Sad case: error fetching appointment" {
				fakeDB.MockGetAppointmentFn = func(ctx context.Context, params domain.Appointment) (*domain.Appointment, error) {
					return nil, fmt.Errorf("db transaction error")
				}
			}

			fakeNotification := notificationMock.NewServiceNotificationMock()

			a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB, fakePubsub, fakeNotification)
			got, err := a.NextRefill(tt.args.ctx, tt.args.clientID)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.NextRefill() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected next refill to be nil for %v", tt.name)
				return
			}
			if !tt.wantErr && !tt.wantNil && got == nil {
				t.Errorf("expected next refill not to be nil for %v", tt.name)
				return
			}
		})
	}
}
