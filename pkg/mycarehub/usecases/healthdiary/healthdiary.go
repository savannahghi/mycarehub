package healthdiary

import (
	"context"
	"fmt"
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/helpers"
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
	CreateHealthDiaryEntry(ctx context.Context, clientID string, note *string, mood string, reportToStaff bool) (bool, error)
}

// ICanRecordHealthDiary contains methods that check whether a client can record a health diary entry
type ICanRecordHealthDiary interface {
	CanRecordHeathDiary(ctx context.Context, clientID string) (bool, error)
}

// IGetRandomQuote defines a method signature that returns a single quote to the frontend. This will be used in place
// of the healthdiary (after it has been filled)
type IGetRandomQuote interface {
	GetClientHealthDiaryQuote(ctx context.Context) (*domain.ClientHealthDiaryQuote, error)
}

// IGetClientHealthDiaryEntry defines a method signature that is used to fetch a client's health diary records
type IGetClientHealthDiaryEntry interface {
	GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error)
}

// UseCasesHealthDiary holds all the interfaces that represents the business logic to implement the health diary
type UseCasesHealthDiary interface {
	ICanRecordHealthDiary
	ICreateHealthDiaryEntry
	IGetRandomQuote
	IGetClientHealthDiaryEntry
}

// UseCasesHealthDiaryImpl embeds the healthdiary logic defined on the domain
type UseCasesHealthDiaryImpl struct {
	Create         infrastructure.Create
	Query          infrastructure.Query
	ServiceRequest servicerequest.UseCaseServiceRequest
}

// NewUseCaseHealthDiaryImpl creates a new instance of health diary
func NewUseCaseHealthDiaryImpl(
	create infrastructure.Create,
	query infrastructure.Query,
	servicerequest servicerequest.UseCaseServiceRequest,
) *UseCasesHealthDiaryImpl {
	return &UseCasesHealthDiaryImpl{
		Create:         create,
		Query:          query,
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
) (bool, error) {
	switch mood {
	case string(enums.MoodVerySad):
		currentTime := time.Now()
		healthDiaryEntry := &domain.ClientHealthDiaryEntry{
			Active:                true,
			Mood:                  mood,
			Note:                  *note,
			EntryType:             "HOME_PAGE_HEALTH_DIARY_ENTRY", //TODO: Make this an enum
			ShareWithHealthWorker: reportToStaff,
			ClientID:              clientID,
			SharedAt:              currentTime,
			CreatedAt:             currentTime,
		}

		err := h.Create.CreateHealthDiaryEntry(ctx, healthDiaryEntry)
		if err != nil {
			helpers.ReportErrorToSentry(err)
			return false, fmt.Errorf("failed to save health diary entry")
		}

		_, err = h.ServiceRequest.CreateServiceRequest(
			ctx,
			clientID,
			string(enums.ServiceRequestTypeRedFlag),
			"",
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
			EntryType:             "HOME_PAGE_HEALTH_DIARY_ENTRY", //TODO: Make this an enum
			ShareWithHealthWorker: false,
			ClientID:              clientID,
			SharedAt:              time.Now(),
		}
		err := h.Create.CreateHealthDiaryEntry(ctx, healthDiaryEntry)
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
func (h UseCasesHealthDiaryImpl) GetClientHealthDiaryQuote(ctx context.Context) (*domain.ClientHealthDiaryQuote, error) {
	return h.Query.GetClientHealthDiaryQuote(ctx)
}

// GetClientHealthDiaryEntries retrieves all health diary entries that belong to a specific user/client
func (h UseCasesHealthDiaryImpl) GetClientHealthDiaryEntries(ctx context.Context, clientID string) ([]*domain.ClientHealthDiaryEntry, error) {
	if clientID == "" {
		return nil, exceptions.EmptyInputErr(fmt.Errorf("missing client ID"))
	}
	return h.Query.GetClientHealthDiaryEntries(ctx, clientID)
}
