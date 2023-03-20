package usecases

import (
	appointment "github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/appointments"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/authority"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/communities"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/metrics"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/notification"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/organisation"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/programs"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/pubsub"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/questionnaires"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/surveys"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/terms"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/user"
)

// MyCareHub represents the my care hub core domain logic
type MyCareHub struct {
	User              user.UseCasesUser
	Terms             terms.UseCasesTerms
	Facility          facility.UseCasesFacility
	SecurityQuestions securityquestions.UseCaseSecurityQuestion
	OTP               otp.UsecaseOTP
	Content           content.UseCasesContent
	Feedback          feedback.UsecaseFeedback
	HealthDiary       healthdiary.UseCasesHealthDiary
	ServiceRequest    servicerequest.UseCaseServiceRequest
	Authority         authority.UsecaseAuthority
	Appointment       appointment.UseCasesAppointments
	Notification      notification.UseCaseNotification
	Surveys           surveys.UsecaseSurveys
	Metrics           metrics.UsecaseMetrics
	Questionnaires    questionnaires.UseCaseQuestionnaire
	Programs          programs.UsecasePrograms
	Organisation      organisation.UseCaseOrganisation
	Pubsub            pubsub.UseCasePubSub
	Community         communities.UseCasesCommunities
}

// NewMyCareHubUseCase initializes a new my care hub instance
func NewMyCareHubUseCase(
	user user.UseCasesUser,
	terms terms.UseCasesTerms,
	facility facility.UseCasesFacility,
	securityQuestions securityquestions.UseCaseSecurityQuestion,
	OTP otp.UsecaseOTP,
	content content.UseCasesContent,
	feedback feedback.UsecaseFeedback,
	healthDiary healthdiary.UseCasesHealthDiary,
	servicerequest servicerequest.UseCaseServiceRequest,
	authority authority.UsecaseAuthority,
	appointment appointment.UseCasesAppointments,
	notification notification.UseCaseNotification,
	surveys surveys.UsecaseSurveys,
	metrics metrics.UsecaseMetrics,
	questionnaires questionnaires.UseCaseQuestionnaire,
	programs programs.UsecasePrograms,
	organisation organisation.UseCaseOrganisation,
	pubsub pubsub.UseCasePubSub,
	communities communities.UseCasesCommunities,
) *MyCareHub {
	return &MyCareHub{
		User:              user,
		Terms:             terms,
		Facility:          facility,
		SecurityQuestions: securityQuestions,
		OTP:               OTP,
		Content:           content,
		Feedback:          feedback,
		HealthDiary:       healthDiary,
		ServiceRequest:    servicerequest,
		Authority:         authority,
		Appointment:       appointment,
		Notification:      notification,
		Surveys:           surveys,
		Metrics:           metrics,
		Questionnaires:    questionnaires,
		Programs:          programs,
		Organisation:      organisation,
		Pubsub:            pubsub,
		Community:         communities,
	}
}
