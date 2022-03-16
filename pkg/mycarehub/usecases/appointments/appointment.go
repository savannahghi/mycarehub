package appointment

import (
	"context"
	"fmt"
	"strconv"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
)

// ICreateAppointments defines method signatures for creating appointments
type ICreateAppointments interface {
	CreateKenyaEMRAppointments(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error)
}

// IUpdateAppointments defines method signatures for updating appointments
type IUpdateAppointments interface {
	UpdateKenyaEMRAppointments(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error)
}

// IListAppointments defines method signatures for listing appointments
type IListAppointments interface {
	FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filterInput []*dto.FiltersInput) (*domain.AppointmentsPage, error)
}

// UseCasesAppointments holds all interfaces required to implement the appointments features
type UseCasesAppointments interface {
	ICreateAppointments
	IUpdateAppointments
	IListAppointments
}

// UseCasesAppointmentsImpl represents appointments implementation
type UseCasesAppointmentsImpl struct {
	Create      infrastructure.Create
	ExternalExt extension.ExternalMethodsExtension
	Query       infrastructure.Query
	Update      infrastructure.Update
}

// NewUseCaseAppointmentsImpl initializes a new appointments usecase
func NewUseCaseAppointmentsImpl(
	ext extension.ExternalMethodsExtension,
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
) *UseCasesAppointmentsImpl {
	return &UseCasesAppointmentsImpl{
		Create:      create,
		ExternalExt: ext,
		Query:       query,
		Update:      update,
	}
}

// CreateKenyaEMRAppointments creates appointments from Kenya EMR
func (a *UseCasesAppointmentsImpl) CreateKenyaEMRAppointments(ctx context.Context, input dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error) {

	MFLCode, err := strconv.Atoi(input.MFLCode)
	if err != nil {
		return nil, err
	}

	exists, err := a.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
	if err != nil {
		return nil, fmt.Errorf("error checking for facility")
	}
	if !exists {
		return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", MFLCode)
	}

	facility, err := a.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving facility: %v", err)
	}

	response := dto.FacilityAppointmentsResponse{MFLCode: input.MFLCode}

	for _, ap := range input.Appointments {
		appointment := domain.Appointment{
			Type:   ap.AppointmentType,
			Status: ap.Status.String(),
			Date:   ap.AppointmentDate,
			Start:  *ap.StartTime(),
			End:    *ap.EndTime(),

			FacilityID: *facility.ID,
		}

		// get client profile using the ccc number
		clientProfile, err := a.Query.GetClientProfileByCCCNumber(ctx, ap.CCCNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get client profile by CCC number")
		}

		clientID := clientProfile.ID

		err = a.Create.CreateAppointment(ctx, appointment, ap.AppointmentUUID, *clientID)
		if err != nil {
			return nil, err
		}

		response.Appointments = append(response.Appointments, dto.AppointmentResponse(ap))
	}

	return &response, nil
}

// UpdateKenyaEMRAppointments updates an appointment with changes from Kenya EMR
func (a *UseCasesAppointmentsImpl) UpdateKenyaEMRAppointments(ctx context.Context, input dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error) {

	MFLCode, err := strconv.Atoi(input.MFLCode)
	if err != nil {

		return nil, err
	}

	exists, err := a.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
	if err != nil {
		return nil, fmt.Errorf("error checking for facility")
	}
	if !exists {
		return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", MFLCode)
	}

	facility, err := a.Query.RetrieveFacilityByMFLCode(ctx, MFLCode, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving facility: %v", err)
	}

	response := dto.FacilityAppointmentsResponse{MFLCode: input.MFLCode}

	for _, ap := range input.Appointments {
		appointment := domain.Appointment{
			Type:       ap.AppointmentType,
			Status:     ap.Status.String(),
			Date:       ap.AppointmentDate,
			Start:      *ap.StartTime(),
			End:        *ap.EndTime(),
			FacilityID: *facility.ID,
		}

		// get client profile using the ccc number
		clientProfile, err := a.Query.GetClientProfileByCCCNumber(ctx, ap.CCCNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get client profile by CCC number")
		}

		clientID := clientProfile.ID

		err = a.Update.UpdateAppointment(ctx, appointment, ap.AppointmentUUID, *clientID)
		if err != nil {
			return nil, err
		}

		response.Appointments = append(response.Appointments, dto.AppointmentResponse(ap))
	}

	return &response, nil
}

// FetchClientAppointments fetches appointments for a client
func (a *UseCasesAppointmentsImpl) FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filterInput []*dto.FiltersInput) (*domain.AppointmentsPage, error) {

	// if user did not provide current page, throw an error
	if err := paginationInput.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("pagination input validation failed: %v", err)
	}

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
	}

	filters := []*domain.FiltersParam{}
	for _, f := range filterInput {
		filter := &domain.FiltersParam{
			Name:     string(f.DataType),
			DataType: f.DataType,
			Value:    f.Value,
		}
		filters = append(filters, filter)
	}

	appointments, pageInfo, err := a.Query.ListAppointments(ctx, &domain.Appointment{ClientID: clientID}, filters, page)
	if err != nil {
		return nil, err
	}

	response := &domain.AppointmentsPage{
		Appointments: appointments,
		Pagination:   *pageInfo,
	}

	return response, nil
}
