package usecases

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/content"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/faq"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/feedback"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/healthdiary"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/servicerequest"
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
	FAQ               faq.UsecaseFAQ
	ServiceRequest    servicerequest.UseCaseServiceRequest
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
	faq faq.UsecaseFAQ,
	servicerequest servicerequest.UseCaseServiceRequest,
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
		FAQ:               faq,
		ServiceRequest:    servicerequest,
	}
}
