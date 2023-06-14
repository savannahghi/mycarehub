package mock

import (
	"context"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/scalarutils"
)

// AppointmentsUseCaseMock mocks the implementation of Appointments usecase methods.
type AppointmentsUseCaseMock struct {
	MockCreateKenyaEMRAppointmentsFn         func(ctx context.Context, facility *domain.Facility, payload dto.AppointmentPayload) ([]*dto.AppointmentPayload, error)
	MockCreateOrUpdateKenyaEMRAppointmentsFn func(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error)
	MockAddPatientsRecordsFn                 func(ctx context.Context, input dto.PatientsRecordsPayload) error
	MockAddPatientRecordFn                   func(ctx context.Context, input dto.PatientRecordPayload) error
	MockUpdateKenyaEMRAppointmentsFn         func(ctx context.Context, facility *domain.Facility, payload dto.AppointmentPayload) (*dto.AppointmentPayload, error)
	MockRescheduleClientAppointmentFn        func(ctx context.Context, appointmentID string, date scalarutils.Date, caregiverID *string) (bool, error)
	MockFetchClientAppointmentsFn            func(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filters []*firebasetools.FilterParam) (*domain.AppointmentsPage, error)
	MockGetAppointmentServiceRequestsFn      func(ctx context.Context, payload dto.AppointmentServiceRequestInput) (*dto.AppointmentServiceRequestsOutput, error)
	MockNextRefillFn                         func(ctx context.Context, clientID string) (*scalarutils.Date, error)
}

// NewAppointmentsUseCaseMock creates in itializes create type mocks
func NewAppointmentsUseCaseMock() *AppointmentsUseCaseMock {
	now := time.Now()
	UUID := gofakeit.UUID()
	contact := "0999999999"
	appointmentsPayload := dto.AppointmentPayload{
		CCCNumber:  "2121221",
		ExternalID: uuid.NewString(),
		AppointmentDate: scalarutils.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
		AppointmentReason: gofakeit.BS(),
	}

	appointmentResponse := dto.AppointmentResponse{
		CCCNumber:  "4783473834",
		ExternalID: gofakeit.UUID(),
		AppointmentDate: scalarutils.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
		AppointmentReason: gofakeit.BS(),
	}

	facilityAppointmemntsResponse := dto.FacilityAppointmentsResponse{
		MFLCode:      "3232323232",
		Appointments: []dto.AppointmentResponse{appointmentResponse},
	}

	domainAppointment := domain.Appointment{
		ID:         gofakeit.UUID(),
		ExternalID: gofakeit.UUID(),
		Reason:     gofakeit.BS(),
		Date: scalarutils.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
		ClientID:                  gofakeit.UUID(),
		FacilityID:                gofakeit.UUID(),
		Provider:                  gofakeit.UUID(),
		HasRescheduledAppointment: false,
		ProgramID:                 gofakeit.UUID(),
	}

	pagination := domain.Pagination{
		Limit:       1,
		CurrentPage: 1,
	}
	appointmentsPage := domain.AppointmentsPage{
		Appointments: []*domain.Appointment{&domainAppointment},
		Pagination:   pagination,
	}
	appointmentServiceRequest := domain.AppointmentServiceRequests{
		ID:         UUID,
		ExternalID: gofakeit.UUID(),
		Reason:     gofakeit.BS(),
		Date: scalarutils.Date{
			Year:  2000,
			Month: 1,
			Day:   1,
		},
		Status:        string(enums.ServiceRequestStatusInProgress),
		InProgressAt:  &now,
		InProgressBy:  &UUID,
		ClientName:    &UUID,
		ClientContact: &contact,
		CCCNumber:     "32323232",
		MFLCODE:       "3232322323",
	}
	appointmentServiceRequestOutput := dto.AppointmentServiceRequestsOutput{
		AppointmentServiceRequests: []domain.AppointmentServiceRequests{appointmentServiceRequest},
	}

	date := scalarutils.Date{
		Year:  2000,
		Month: 1,
		Day:   1,
	}
	return &AppointmentsUseCaseMock{
		MockCreateKenyaEMRAppointmentsFn: func(ctx context.Context, facility *domain.Facility, payload dto.AppointmentPayload) ([]*dto.AppointmentPayload, error) {
			return []*dto.AppointmentPayload{&appointmentsPayload}, nil
		},
		MockCreateOrUpdateKenyaEMRAppointmentsFn: func(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error) {
			return &facilityAppointmemntsResponse, nil
		},
		MockAddPatientsRecordsFn: func(ctx context.Context, input dto.PatientsRecordsPayload) error {
			return nil
		},
		MockAddPatientRecordFn: func(ctx context.Context, input dto.PatientRecordPayload) error {
			return nil
		},
		MockUpdateKenyaEMRAppointmentsFn: func(ctx context.Context, facility *domain.Facility, payload dto.AppointmentPayload) (*dto.AppointmentPayload, error) {
			return &appointmentsPayload, nil
		},
		MockRescheduleClientAppointmentFn: func(ctx context.Context, appointmentID string, date scalarutils.Date, caregiverID *string) (bool, error) {
			return false, nil
		},
		MockFetchClientAppointmentsFn: func(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filters []*firebasetools.FilterParam) (*domain.AppointmentsPage, error) {
			return &appointmentsPage, nil
		},
		MockGetAppointmentServiceRequestsFn: func(ctx context.Context, payload dto.AppointmentServiceRequestInput) (*dto.AppointmentServiceRequestsOutput, error) {
			return &appointmentServiceRequestOutput, nil
		},
		MockNextRefillFn: func(ctx context.Context, clientID string) (*scalarutils.Date, error) {
			return &date, nil
		},
	}
}

// CreateKenyaEMRAppointments mocks the implementation of creating kenya emr appointments
func (gm *AppointmentsUseCaseMock) CreateKenyaEMRAppointments(ctx context.Context, facility *domain.Facility, payload dto.AppointmentPayload) ([]*dto.AppointmentPayload, error) {
	return gm.MockCreateKenyaEMRAppointmentsFn(ctx, facility, payload)
}

// CreateOrUpdateKenyaEMRAppointments mocks the implementation of Create Or UpdateKenya EMR Appointments
func (gm *AppointmentsUseCaseMock) CreateOrUpdateKenyaEMRAppointments(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error) {
	return gm.MockCreateOrUpdateKenyaEMRAppointmentsFn(ctx, payload)
}

// AddPatientsRecords mocks the implementation of adding patients records
func (gm *AppointmentsUseCaseMock) AddPatientsRecords(ctx context.Context, input dto.PatientsRecordsPayload) error {
	return gm.MockAddPatientsRecordsFn(ctx, input)
}

// AddPatientRecord mocks the implementation of adding a new patient record
func (gm *AppointmentsUseCaseMock) AddPatientRecord(ctx context.Context, input dto.PatientRecordPayload) error {
	return gm.MockAddPatientRecordFn(ctx, input)
}

// UpdateKenyaEMRAppointments mocks the implementation of updating kenya EMR appointments
func (gm *AppointmentsUseCaseMock) UpdateKenyaEMRAppointments(ctx context.Context, facility *domain.Facility, payload dto.AppointmentPayload) (*dto.AppointmentPayload, error) {
	return gm.MockUpdateKenyaEMRAppointmentsFn(ctx, facility, payload)
}

// RescheduleClientAppointment mocks the implementation of rescheduling client appointments
func (gm *AppointmentsUseCaseMock) RescheduleClientAppointment(ctx context.Context, appointmentID string, date scalarutils.Date, caregiverID *string) (bool, error) {
	return gm.MockRescheduleClientAppointmentFn(ctx, appointmentID, date, caregiverID)
}

// FetchClientAppointments mocks the implementation of the FetchClientAppointments
func (gm *AppointmentsUseCaseMock) FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filters []*firebasetools.FilterParam) (*domain.AppointmentsPage, error) {
	return gm.MockFetchClientAppointmentsFn(ctx, clientID, paginationInput, filters)
}

// GetAppointmentServiceRequests mocks the implementation of GetAppointments service requests
func (gm *AppointmentsUseCaseMock) GetAppointmentServiceRequests(ctx context.Context, payload dto.AppointmentServiceRequestInput) (*dto.AppointmentServiceRequestsOutput, error) {
	return gm.MockGetAppointmentServiceRequestsFn(ctx, payload)
}

// NextRefill mocks the implementation of the nextRefill
func (gm *AppointmentsUseCaseMock) NextRefill(ctx context.Context, clientID string) (*scalarutils.Date, error) {
	return gm.MockNextRefillFn(ctx, clientID)
}
