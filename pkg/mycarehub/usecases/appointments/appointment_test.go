package appointment

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	extensionMock "github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension/mock"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	pgMock "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/mock"
	"github.com/savannahghi/scalarutils"
)

func TestUseCasesAppointmentsImpl_CreateKenyaEMRAppointments(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()

	a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB)

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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: error checking facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("error retrieving facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "sad case: error creating appointment" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCreateAppointment = func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID string) error {
					return fmt.Errorf("cannot create appointment")
				}
			}

			if tt.name == "happy case: success creating appointment" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockCreateAppointment = func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID string) error {
					return nil
				}
			}

			got, err := a.CreateKenyaEMRAppointments(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.CreateKenyaEMRAppointments() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesAppointmentsImpl_UpdateKenyaEMRAppointments(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()

	a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB)

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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
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
							CCCNumber:       "1234",
							AppointmentUUID: gofakeit.UUID(),
							AppointmentType: "Dental",
							Status:          enums.AppointmentStatusCompleted,
							AppointmentDate: scalarutils.Date{
								Year:  2020,
								Month: 12,
								Day:   12,
							},
							TimeSlot: "8:00 - 12:00",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: error checking facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, fmt.Errorf("error retrieving facility")
				}
			}

			if tt.name == "sad case: facility doesn't exist" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "sad case: error retrieving facility" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}
				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					return nil, fmt.Errorf("cannot retrieve facility by mfl code")
				}
			}

			if tt.name == "sad case: error updating appointment" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockUpdateAppointment = func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID string) error {
					return fmt.Errorf("error updating appointment")
				}
			}

			if tt.name == "happy case: success creating appointment" {
				fakeDB.MockCheckFacilityExistsByMFLCode = func(ctx context.Context, MFLCode int) (bool, error) {
					return true, nil
				}

				fakeDB.MockRetrieveFacilityByMFLCodeFn = func(ctx context.Context, MFLCode int, isActive bool) (*domain.Facility, error) {
					id := gofakeit.UUID()
					return &domain.Facility{ID: &id}, nil
				}

				fakeDB.MockUpdateAppointment = func(ctx context.Context, appointment domain.Appointment, appointmentUUID, clientID string) error {
					return nil
				}
			}

			got, err := a.UpdateKenyaEMRAppointments(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("UseCasesAppointmentsImpl.UpdateKenyaEMRAppointments() error = %v, wantErr %v", err, tt.wantErr)
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

func TestUseCasesAppointmentsImpl_FetchClientAppointments(t *testing.T) {
	fakeDB := pgMock.NewPostgresMock()
	fakeExtension := extensionMock.NewFakeExtension()

	a := NewUseCaseAppointmentsImpl(fakeExtension, fakeDB, fakeDB, fakeDB)

	type args struct {
		ctx             context.Context
		clientID        string
		paginationInput dto.PaginationsInput
		filterInput     []*dto.FiltersInput
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
				filterInput:     []*dto.FiltersInput{},
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
				filterInput: []*dto.FiltersInput{},
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
				filterInput: []*dto.FiltersInput{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "sad case: error listing appointments" {
				fakeDB.MockListAppointments = func(ctx context.Context, params *domain.Appointment, filter []*domain.FiltersParam, pagination *domain.Pagination) ([]*domain.Appointment, *domain.Pagination, error) {
					return nil, nil, fmt.Errorf("error listing appointments")
				}
			}
			got, err := a.FetchClientAppointments(tt.args.ctx, tt.args.clientID, tt.args.paginationInput, tt.args.filterInput)
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
