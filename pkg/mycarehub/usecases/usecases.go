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
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/screeningtools"
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
	Community         communities.UseCasesCommunities
	ScreeningTools    screeningtools.UseCasesScreeningTools
	Appointment       appointment.UseCasesAppointments
	Notification      notification.UseCaseNotification
	Surveys           surveys.UsecaseSurveys
	Metrics           metrics.UsecaseMetrics
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
	community communities.UseCasesCommunities,
	screeningTools screeningtools.UseCasesScreeningTools,
	appointment appointment.UseCasesAppointments,
	notification notification.UseCaseNotification,
	surveys surveys.UsecaseSurveys,
	metrics metrics.UsecaseMetrics,
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
		Community:         community,
		ScreeningTools:    screeningTools,
		Appointment:       appointment,
		Notification:      notification,
		Surveys:           surveys,
		Metrics:           metrics,
	}
}
