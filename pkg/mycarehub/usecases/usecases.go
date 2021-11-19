package usecases

import (
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/facility"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/otp"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/usecases/securityquestions"
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
}

// NewMyCareHubUseCase initializes a new my care hub instance
func NewMyCareHubUseCase(
	user user.UseCasesUser,
	terms terms.UseCasesTerms,
	facility facility.UseCasesFacility,
	securityQuestions securityquestions.UseCaseSecurityQuestion,
	OTP otp.UsecaseOTP,
) *MyCareHub {
	return &MyCareHub{
		User:              user,
		Terms:             terms,
		Facility:          facility,
		SecurityQuestions: securityQuestions,
		OTP:               OTP,
	}
}
