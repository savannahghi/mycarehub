package appointment

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/extension"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	pubsubmessaging "github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/services/pubsub"
)

// ICreateAppointments defines method signatures for creating appointments
type ICreateAppointments interface {
	CreateKenyaEMRAppointments(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error)
}

// ICreateHealthRecords defines method signatures for creating health records
type ICreateHealthRecords interface {
	AddPatientsRecords(ctx context.Context, input dto.PatientsRecordsPayload) error
	AddPatientRecord(ctx context.Context, input dto.PatientRecordPayload) error
}

// IUpdateAppointments defines method signatures for updating appointments
type IUpdateAppointments interface {
	UpdateKenyaEMRAppointments(ctx context.Context, payload dto.FacilityAppointmentsPayload) (*dto.FacilityAppointmentsResponse, error)
	RescheduleClientAppointment(ctx context.Context, appointmentID string) (bool, error)
}

// IListAppointments defines method signatures for listing appointments
type IListAppointments interface {
	FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filters []*firebasetools.FilterParam) (*domain.AppointmentsPage, error)
	GetAppointmentServiceRequests(ctx context.Context, payload dto.AppointmentServiceRequestInput) (*dto.AppointmentServiceRequestOutput, error)
}

// UseCasesAppointments holds all interfaces required to implement the appointments features
type UseCasesAppointments interface {
	ICreateHealthRecords
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
	Pubsub      pubsubmessaging.ServicePubsub
}

// NewUseCaseAppointmentsImpl initializes a new appointments usecase
func NewUseCaseAppointmentsImpl(
	ext extension.ExternalMethodsExtension,
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
	pubsub pubsubmessaging.ServicePubsub,
) *UseCasesAppointmentsImpl {
	return &UseCasesAppointmentsImpl{
		Create:      create,
		ExternalExt: ext,
		Query:       query,
		Update:      update,
		Pubsub:      pubsub,
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
			Status: ap.Status,
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
		if err != nil {
			return nil, fmt.Errorf("error parsing appointment date")
		}
		appointmentToUpdate := map[string]interface{}{
			"appointment_type": ap.AppointmentType,
			"status":           ap.Status,
			"date":             ap.AppointmentDate.AsTime(),
			"start_time":       *ap.StartTime(),
			"end_time":         *ap.EndTime(),
			"facility_id":      *facility.ID,
		}

		// get client profile using the ccc number
		_, err := a.Query.GetClientProfileByCCCNumber(ctx, ap.CCCNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get client profile by CCC number")
		}

		appointment, err := a.Query.GetAppointmentByAppointmentUUID(ctx, ap.AppointmentUUID)
		if err != nil {
			return nil, fmt.Errorf("failed to get appointment by UUID")
		}

		appointmentInput := &domain.Appointment{ID: appointment.ID}

		_, err = a.Update.UpdateAppointment(ctx, appointmentInput, appointmentToUpdate)
		if err != nil {
			return nil, err
		}

		response.Appointments = append(response.Appointments, dto.AppointmentResponse(ap))
	}

	return &response, nil
}

// FetchClientAppointments fetches appointments for a client
func (a *UseCasesAppointmentsImpl) FetchClientAppointments(ctx context.Context, clientID string, paginationInput dto.PaginationsInput, filters []*firebasetools.FilterParam) (*domain.AppointmentsPage, error) {

	// if user did not provide current page, throw an error
	if err := paginationInput.Validate(); err != nil {
		helpers.ReportErrorToSentry(err)
		return nil, fmt.Errorf("pagination input validation failed: %v", err)
	}

	page := &domain.Pagination{
		Limit:       paginationInput.Limit,
		CurrentPage: paginationInput.CurrentPage,
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

// AddPatientsRecords adds records for multiple clients and is especially useful when performing a bulk creation from KenyaEMR
func (a *UseCasesAppointmentsImpl) AddPatientsRecords(ctx context.Context, input dto.PatientsRecordsPayload) error {

	MFLCode, err := strconv.Atoi(input.MFLCode)
	if err != nil {
		return err
	}

	exists, err := a.Query.CheckFacilityExistsByMFLCode(ctx, MFLCode)
	if err != nil {
		return fmt.Errorf("error checking for facility")
	}
	if !exists {
		return fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", MFLCode)
	}

	var errs error
	for _, record := range input.Records {
		record.MFLCode = MFLCode
		err = a.AddPatientRecord(ctx, record)
		if err != nil {
			// accumulate errors rather than failing early
			errs = multierror.Append(errs, err)
		}
	}

	if errs != nil {
		return err
	}

	return nil
}

// AddPatientRecord adds records for a single client. It is used for push updates for a particular client
func (a *UseCasesAppointmentsImpl) AddPatientRecord(ctx context.Context, input dto.PatientRecordPayload) error {
	if input.CCCNumber == "" {
		return fmt.Errorf("ccc number is required")
	}

	_, err := a.Query.RetrieveFacilityByMFLCode(ctx, input.MFLCode, true)
	if err != nil {
		return fmt.Errorf("error retrieving facility with mfl code: %v", input.MFLCode)
	}

	client, err := a.Query.GetClientProfileByCCCNumber(ctx, input.CCCNumber)
	if err != nil {
		return fmt.Errorf("error retrieving client with ccc number: %v", input.CCCNumber)
	}

	for _, vital := range input.VitalSigns {
		payload := dto.PatientVitalSignOutput{
			PatientID:      *client.FHIRPatientID,
			OrganizationID: "", //TODO: FHIR organization ID
			Name:           vital.Name,
			ConceptID:      vital.ConceptID,
			Value:          vital.Value,
			Date:           vital.Date,
		}
		err = a.Pubsub.NotifyCreateVitals(ctx, &payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("failed to publish to create patient topic: %v", err)
		}
	}
	for _, allergy := range input.Allergies {
		payload := dto.PatientAllergyOutput{
			PatientID:      *client.FHIRPatientID,
			OrganizationID: "", //TODO: FHIR organization ID
			Name:           allergy.Name,
			ConceptID:      allergy.AllergyConceptID,
			Date:           allergy.Date,
			Reaction: dto.AllergyReaction{
				Name:      allergy.Reaction,
				ConceptID: allergy.ReactionConceptID,
			},
			Severity: dto.AllergySeverity{
				Name:      allergy.Severity,
				ConceptID: allergy.SeverityConceptID,
			},
		}
		err = a.Pubsub.NotifyCreateAllergy(ctx, &payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("failed to publish to create allergy topic: %v", err)
		}
	}
	for _, medication := range input.Medications {
		payload := dto.PatientMedicationOutput{
			PatientID:      *client.FHIRPatientID,
			OrganizationID: "", //TODO: FHIR organization ID
			Name:           medication.Name,
			ConceptID:      medication.MedicationConceptID,
			Date:           medication.Date,
			Value:          medication.Value,
		}

		if medication.DrugConceptID != nil {
			payload.Drug = &dto.MedicationDrug{
				ConceptID: medication.DrugConceptID,
			}
		}

		err = a.Pubsub.NotifyCreateMedication(ctx, &payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("failed to publish to create medication topic: %v", err)
		}
	}
	for _, result := range input.TestResults {
		payload := dto.PatientTestResultOutput{
			PatientID:      *client.FHIRPatientID,
			OrganizationID: "", //TODO: FHIR organization ID
			Name:           result.Name,
			ConceptID:      result.TestConceptID,
			Date:           result.Date,
			Result: dto.TestResult{
				Name:      result.Result,
				ConceptID: result.ResultConceptID,
			},
		}
		err = a.Pubsub.NotifyCreateTestResult(ctx, &payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("failed to publish to create test result topic: %v", err)
		}
	}
	for _, order := range input.TestOrders {
		payload := dto.PatientTestOrderOutput{
			PatientID:      *client.FHIRPatientID,
			OrganizationID: "", //TODO: FHIR organization ID
			Name:           order.Name,
			ConceptID:      order.ConceptID,
			Date:           order.Date,
		}
		err := a.Pubsub.NotifyCreateTestOrder(ctx, &payload)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			log.Printf("failed to publish to create test order topic: %v", err)
		}
	}

	return nil
}

// GetAppointmentServiceRequests returns a list of appointment service requests
func (a *UseCasesAppointmentsImpl) GetAppointmentServiceRequests(ctx context.Context, payload dto.AppointmentServiceRequestInput) (*dto.AppointmentServiceRequestOutput, error) {

	exists, err := a.Query.CheckFacilityExistsByMFLCode(ctx, payload.MFLCode)
	if err != nil {
		return nil, fmt.Errorf("error checking for facility")
	}

	if !exists {
		return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", payload.MFLCode)
	}

	mflCode := strconv.Itoa(payload.MFLCode)

	appointmentServiceRequests, err := a.Query.GetAppointmentServiceRequests(ctx, *payload.LastSyncTime, mflCode)
	if err != nil {
		return nil, fmt.Errorf("error getting appointment service requests")
	}

	if appointmentServiceRequests == nil {
		return &dto.AppointmentServiceRequestOutput{
			AppointmentServiceRequests: []domain.AppointmentServiceRequests{},
		}, nil
	}

	return &dto.AppointmentServiceRequestOutput{
		AppointmentServiceRequests: appointmentServiceRequests,
	}, nil
}

// RescheduleClientAppointment creates a service request to reschedule a client appointment
func (a *UseCasesAppointmentsImpl) RescheduleClientAppointment(ctx context.Context, appointmentID string) (bool, error) {
	if appointmentID == "" {
		return false, fmt.Errorf("invalid input provided")
	}

	appointment, err := a.Query.GetClientAppointmentByID(ctx, appointmentID)
	if err != nil {
		return false, fmt.Errorf("error getting client appointment: %v", err)
	}

	client, err := a.Query.GetClientProfileByClientID(ctx, appointment.ClientID)
	if err != nil {
		return false, fmt.Errorf("error getting client profile")
	}

	userProfile, err := a.Query.GetUserProfileByUserID(ctx, client.UserID)
	if err != nil {
		return false, fmt.Errorf("error getting user profile")
	}

	// update appointment has rescheduled field to true
	_, err = a.Update.UpdateAppointment(ctx, &domain.Appointment{ID: appointment.ID}, map[string]interface{}{"has_rescheduled_appointment": true})
	if err != nil {
		return false, fmt.Errorf("error updating appointment")
	}

	appointmentDate := appointment.Date.AsTime().Format("02-Jan-2006")

	serviceRequest := &dto.ServiceRequestInput{
		Active:      true,
		RequestType: enums.ServiceRequestTypeAppointments.String(),
		Request:     fmt.Sprintf(`%v has requested to reschedule a %v appointment from %v to a later date`, userProfile.Name, appointment.Type, appointmentDate),
		Status:      enums.ServiceRequestStatusPending.String(),
		ClientID:    appointment.ClientID,
		FacilityID:  appointment.FacilityID,
		Meta:        map[string]interface{}{"appointmentID": appointmentID, "appointmentUUID": appointment.AppointmentUUID},
	}

	err = a.Create.CreateServiceRequest(ctx, serviceRequest)
	if err != nil {
		return false, fmt.Errorf("error creating service request")
	}

	return true, nil
}
