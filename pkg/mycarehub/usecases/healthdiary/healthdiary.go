package healthdiary

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/exceptions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
)

// The healthdiary is used for engagement with clients on a day-by-day basis.
// The idea behind it is to track sustained changes in client's moods. The moods range
// from VERY_HAPPY, HAPPY, NEUTRAL, SAD, VERY_SAD. When a client fills the mood board, a health diary
// entry is recorded in the database. In cases where the client is VERY_SAD, the client is asked if they
// want to report it to a healthcare worker and if they do, a service request is created. The service request
// is a task for the healthcare worker on the platform. All this should happen within a 24 hour time window. If
// a health diary was filled within the past 24 hours, the client is shown an inspirational post on the frontend
// and if it hasn't been filled, we show them the health diary.

// ICreateHealthDiaryEntry is an interface that holds the method signature for creating a health diary entry
type ICreateHealthDiaryEntry interface {
	CreateHealthDiaryEntry(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool, caregiverID *string) (bool, error)
}

// ICanRecordHealthDiary contains methods that check whether a client can record a health diary entry
type ICanRecordHealthDiary interface {
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
}

// IGetRandomQuote defines a method signature that returns a single quote to the frontend. This will be used in place
// of the healthdiary (after it has been filled)
type IGetRandomQuote interface {
	GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error)
}

// IGetClientHealthDiaryEntry defines a method signature that is used to fetch a client's health diary records
type IGetClientHealthDiaryEntry interface {
	GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error)
	GetFacilityHealthDiaryEntries(ctx context.Context, input dto.FetchHealthDiaryEntries) (*dto.HealthDiaryEntriesResponse, error)
	GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error)
	GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error)
}

// IShareHealthDiaryEntry contains the methods to share the health diary with the health care worker
type IShareHealthDiaryEntry interface {
	ShareHealthDiaryEntry(ctx context.Context, healthDiaryEntryID string, shareEntireHealthDiary bool) (bool, error)
}

// UseCasesHealthDiary holds all the interfaces that represents the business logic to implement the health diary
type UseCasesHealthDiary interface {
	ICanRecordHealthDiary
	ICreateHealthDiaryEntry
	IGetRandomQuote
	IGetClientHealthDiaryEntry
	IShareHealthDiaryEntry
}

// UseCasesHealthDiaryImpl embeds the healthdiary logic defined on the domain
type UseCasesHealthDiaryImpl struct {
	Create         infrastructure.Create
	Query          infrastructure.Query
	Update         infrastructure.Update
	ServiceRequest servicerequest.UseCaseServiceRequest
}

// NewUseCaseHealthDiaryImpl creates a new instance of health diary
func NewUseCaseHealthDiaryImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	update infrastructure.Update,
	servicerequest servicerequest.UseCaseServiceRequest,
) *UseCasesHealthDiaryImpl {
	return &UseCasesHealthDiaryImpl{
		Create:         create,
		Query:          query,
		Update:         update,
		ServiceRequest: servicerequest,
	}
}

// CreateHealthDiaryEntry captures a client's mood and creates a health diary entry. This will be used to
// track the client's moods on a day-to-day basis
func (h UseCasesHealthDiaryImpl) CreateHealthDiaryEntry(
	ctx context.Context,
	clientID string,
	note *string,
	mood string,
	reportToStaff bool,
	caregiverID *string,
) (bool, error) {
	currentTime := time.Now()
	clientProfile, err := h.Query.GetClientProfileByClientID(ctx, clientID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, fmt.Errorf("error querying client profile: %w", err)
	}

	switch mood {
	case enums.MoodVerySad.String():
		healthDiaryEntry := &domain.ClientHealthDiaryEntry{
			Active:                true,
			Mood:                  mood,
			Note:                  *note,
			EntryType:             enums.ServiceRequestTypeHomePageHealthDiary.String(),
			ShareWithHealthWorker: reportToStaff,
			ClientID:              clientID,
			CreatedAt:             currentTime,
			ProgramID:             clientProfile.User.CurrentProgramID,
			OrganisationID:        clientProfile.User.CurrentOrganizationID,
			CaregiverID:           caregiverID,
		}

		if reportToStaff {
			healthDiaryEntry.SharedAt = &currentTime
		}

		healthDiaryEntry, err := h.Create.CreateHealthDiaryEntry(ctx, healthDiaryEntry)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to save health diary entry")
		}

		serviceRequestInput := &dto.ServiceRequestInput{
			ClientID:    clientID,
			Flavour:     feedlib.FlavourConsumer,
			RequestType: enums.ServiceRequestTypeRedFlag.String(),
			Request:     fmt.Sprintf("%s is feeling very sad. Please reach out and help them to feel better.", clientProfile.User.Name),
			FacilityID:  *clientProfile.DefaultFacility.ID,
			ClientName:  &clientProfile.User.Name,
			Meta: map[string]interface{}{
				"note":               &note,
				"healthDiaryEntryID": healthDiaryEntry.ID,
			},
			ProgramID:      clientProfile.User.CurrentProgramID,
			OrganisationID: clientProfile.User.CurrentOrganizationID,
			CaregiverID:    caregiverID,
		}

		_, err = h.ServiceRequest.CreateServiceRequest(
			ctx,
			serviceRequestInput,
		)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create service request: %v", err)
		}

	default:
		healthDiaryEntry := &domain.ClientHealthDiaryEntry{
			Active:                true,
			Mood:                  mood,
			Note:                  *note,
			EntryType:             enums.ServiceRequestTypeHomePageHealthDiary.String(),
			ClientID:              clientID,
			ShareWithHealthWorker: reportToStaff,
			ProgramID:             clientProfile.User.CurrentProgramID,
			OrganisationID:        clientProfile.User.CurrentOrganizationID,
			CaregiverID:           caregiverID,
		}
		if reportToStaff {
			healthDiaryEntry.SharedAt = &currentTime
		}
		_, err := h.Create.CreateHealthDiaryEntry(ctx, healthDiaryEntry)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to save health diary entry")
		}
	}
	return true, nil
}

// CanRecordHeathDiary implements check for eligibility of a health diary to be shown to a user
func (h UseCasesHealthDiaryImpl) CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error) {
	if clientID == "" {
		return false, exceptions.EmptyInputErr(fmt.Errorf("empty client ID value passed in input"))
	}
	return h.Query.CanRecordHeathDiary(ctx, clientID)
}

// GetClientHealthDiaryQuote gets a quote from the database to display on the UI. This happens after a client has already
// filled in their health diary.
func (h UseCasesHealthDiaryImpl) GetClientHealthDiaryQuote(ctx context.Context, limit int) ([]*domain.ClientHealthDiaryQuote, error) {
	return h.Query.GetClientHealthDiaryQuote(ctx, limit)
}

// GetClientHealthDiaryEntries retrieves all health diary entries that belong to a specific user/client
func (h UseCasesHealthDiaryImpl) GetClientHealthDiaryEntries(ctx context.Context, clientID string, moodType *enums.Mood, shared *bool) ([]*domain.ClientHealthDiaryEntry, error) {
	if clientID == "" {
		return nil, exceptions.EmptyInputErr(fmt.Errorf("missing client ID"))
	}
	return h.Query.GetClientHealthDiaryEntries(ctx, clientID, moodType, shared)
}

// GetFacilityHealthDiaryEntries retrieves all the health diary entries that have been recorded by clients
// from a specified facility and have not yet been synced to KenyaEMR.
// This will be used by the KenyaEMR module to retrieve the health diaries and save them into KenyaEMR database
func (h UseCasesHealthDiaryImpl) GetFacilityHealthDiaryEntries(ctx context.Context, input dto.FetchHealthDiaryEntries) (*dto.HealthDiaryEntriesResponse, error) {
	exists, err := h.Query.CheckFacilityExistsByIdentifier(ctx, &dto.FacilityIdentifierInput{
		Type:  enums.FacilityIdentifierTypeMFLCode,
		Value: strconv.Itoa(input.MFLCode),
	})
	if err != nil {
		return nil, fmt.Errorf("error checking for facility")
	}

	if !exists {
		return nil, fmt.Errorf("facility with provided MFL code doesn't exist, code: %v", input.MFLCode)
	}

	facility, err := h.Query.RetrieveFacilityByIdentifier(ctx, &dto.FacilityIdentifierInput{
		Type:  enums.FacilityIdentifierTypeMFLCode,
		Value: strconv.Itoa(input.MFLCode),
	}, true)
	if err != nil {
		return nil, fmt.Errorf("failed to get facility: %v", err)
	}

	clients, err := h.Query.GetClientsInAFacility(ctx, *facility.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users in %v facility: %v", facility.Name, err)
	}

	response := dto.HealthDiaryEntriesResponse{
		MFLCode:            input.MFLCode,
		HealthDiaryEntries: []*domain.ClientHealthDiaryEntry{},
	}

	for _, client := range clients {
		healthDiaryEntry, err := h.GetRecentHealthDiaryEntries(ctx, *input.LastSyncTime, client)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch client health diary entries: %v", err)
		}
		response.HealthDiaryEntries = append(response.HealthDiaryEntries, healthDiaryEntry...)
	}

	return &response, nil
}

// GetRecentHealthDiaryEntries fetches the most recent health diary entries. It returns the new entries
// that were added after the last synced time. This will help KenyEMR module fetch for new health diary entries
func (h UseCasesHealthDiaryImpl) GetRecentHealthDiaryEntries(ctx context.Context, lastSyncTime time.Time, client *domain.ClientProfile) ([]*domain.ClientHealthDiaryEntry, error) {
	return h.Query.GetRecentHealthDiaryEntries(ctx, lastSyncTime, client)
}

// ShareHealthDiaryEntry create a service request when the client opts to share their service request
func (h UseCasesHealthDiaryImpl) ShareHealthDiaryEntry(ctx context.Context, healthDiaryEntryID string, shareEntireHealthDiary bool) (bool, error) {
	healthDiaryEntry, err := h.Query.GetHealthDiaryEntryByID(ctx, healthDiaryEntryID)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	var payload *domain.ClientHealthDiaryEntry
	// If shareEntireHealthDiary is true, update the health diary records with the client's ID.
	// Else, update the record whose records match both the client's ID and the health diary entry ID.
	if shareEntireHealthDiary {
		payload = &domain.ClientHealthDiaryEntry{
			ClientID: healthDiaryEntry.ClientID,
		}
	} else {
		payload = &domain.ClientHealthDiaryEntry{
			ID:       healthDiaryEntry.ID,
			ClientID: healthDiaryEntry.ClientID,
		}
	}

	updateData := map[string]interface{}{
		"share_with_health_worker": true,
		"shared_at":                time.Now(),
	}

	err = h.Update.UpdateHealthDiary(ctx, payload, updateData)
	if err != nil {
		helpers.ReportErrorToSentry(err)
		return false, err
	}

	if healthDiaryEntry.Mood == enums.MoodVerySad.String() || healthDiaryEntry.Mood == enums.MoodSad.String() {
		serviceRequestInput := &dto.ServiceRequestInput{
			RequestType:    healthDiaryEntry.EntryType,
			Status:         enums.ServiceRequestStatusPending.String(),
			Request:        healthDiaryEntry.Note,
			ClientID:       healthDiaryEntry.ClientID,
			Flavour:        feedlib.FlavourConsumer,
			ProgramID:      healthDiaryEntry.ProgramID,
			OrganisationID: healthDiaryEntry.OrganisationID,
		}

		_, err := h.ServiceRequest.CreateServiceRequest(ctx, serviceRequestInput)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to create service request: %v", err)
		}
	}

	return true, nil
}

// GetSharedHealthDiaryEntries fetches the most recent health diary(ies) shared by the client with the health worker
func (h UseCasesHealthDiaryImpl) GetSharedHealthDiaryEntries(ctx context.Context, clientID string, facilityID string) ([]*domain.ClientHealthDiaryEntry, error) {
	if facilityID == "" || clientID == "" {
		return nil, fmt.Errorf("neither facility id nor client id can be empty")
	}

	return h.Query.GetSharedHealthDiaryEntries(ctx, clientID, facilityID)
}
