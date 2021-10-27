package mock

import (
	"context"
	"time"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/dto"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
)

// UpdateMock ...
type UpdateMock struct {
	UpdateUserLastSuccessfulLoginFn func(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateUserLastFailedLoginFn     func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error
	UpdateUserFailedLoginCountFn    func(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error
	UpdateUserNextAllowedLoginFn    func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error
	TransferClientFn                func(
		ctx context.Context,
		clientID string,
		originFacilityID string,
		destinationFacilityID string,
		reason enums.TransferReason,
		notes string,
	) (bool, error)
	UpdateStaffUserProfileFn func(ctx context.Context, userID string, user *dto.UserInput, staff *dto.StaffProfileInput) (bool, error)
}

// NewUpdateMock initializes a new instance of `GormMock` then mocking the case of success.
func NewUpdateMock() *UpdateMock {
	return &UpdateMock{
		UpdateUserLastSuccessfulLoginFn: func(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateUserLastFailedLoginFn: func(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateUserFailedLoginCountFn: func(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateUserNextAllowedLoginFn: func(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error {
			return nil
		},

		UpdateStaffUserProfileFn: func(ctx context.Context, userID string, user *dto.UserInput, staff *dto.StaffProfileInput) (bool, error) {
			return true, nil
		},

		TransferClientFn: func(ctx context.Context, clientID, originFacilityID, destinationFacilityID string, reason enums.TransferReason, notes string) (bool, error) {
			return true, nil
		},
	}
}

//UpdateUserLastSuccessfulLogin ...
func (um *UpdateMock) UpdateUserLastSuccessfulLogin(ctx context.Context, userID string, lastLoginTime time.Time, flavour feedlib.Flavour) error {
	return um.UpdateUserLastSuccessfulLoginFn(ctx, userID, lastLoginTime, flavour)
}

// UpdateUserLastFailedLogin ...
func (um *UpdateMock) UpdateUserLastFailedLogin(ctx context.Context, userID string, lastFailedLoginTime time.Time, flavour feedlib.Flavour) error {
	return um.UpdateUserLastFailedLoginFn(ctx, userID, lastFailedLoginTime, flavour)
}

// UpdateUserFailedLoginCount ...
func (um *UpdateMock) UpdateUserFailedLoginCount(ctx context.Context, userID string, failedLoginCount string, flavour feedlib.Flavour) error {
	return um.UpdateUserFailedLoginCountFn(ctx, userID, failedLoginCount, flavour)
}

// UpdateUserNextAllowedLogin ...
func (um *UpdateMock) UpdateUserNextAllowedLogin(ctx context.Context, userID string, nextAllowedLoginTime time.Time, flavour feedlib.Flavour) error {
	return um.UpdateUserNextAllowedLoginFn(ctx, userID, nextAllowedLoginTime, flavour)
}

// UpdateStaffUserProfile mocks the implementation of  UpdateStaffUserProfile method.
func (um *UpdateMock) UpdateStaffUserProfile(ctx context.Context, userID string, user *dto.UserInput, staff *dto.StaffProfileInput) (bool, error) {
	return um.UpdateStaffUserProfileFn(ctx, userID, user, staff)
}

// TransferClient mocks the implementation of  TransferClient method
func (um *UpdateMock) TransferClient(
	ctx context.Context,
	clientID string,
	originFacilityID string,
	destinationFacilityID string,
	reason enums.TransferReason,
	notes string,
) (bool, error) {
	return um.TransferClientFn(ctx, clientID, originFacilityID, destinationFacilityID, reason, notes)
}
