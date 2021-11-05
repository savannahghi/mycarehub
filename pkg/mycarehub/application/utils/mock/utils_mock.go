package mock

import (
	"time"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

// UtilsMock mocks the implementation of utils methods.
type UtilsMock struct {
	MockCheckPINExpiryFn func(currentTime time.Time, pinData *domain.UserPIN) bool
}

// NewUtilsMock creates in itializes create type mocks
func NewUtilsMock() *UtilsMock {
	return &UtilsMock{
		MockCheckPINExpiryFn: func(currentTime time.Time, pinData *domain.UserPIN) bool {
			return true
		},
	}
}

// CheckPINExpiry ...
func (u *UtilsMock) CheckPINExpiry(currentTime time.Time, pinData *domain.UserPIN) bool {
	return u.MockCheckPINExpiryFn(currentTime, pinData)
}
