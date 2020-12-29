package repository

import (
	"context"

	"gitlab.slade360emr.com/go/base"
)

// OnboardingRepository interface that provide access to all persistent storage operations
type OnboardingRepository interface {
	base.UserProfileRepository
	GetUserProfile(ctx context.Context, uid string) (*base.UserProfile, error)
	GetUserProfileByID(ctx context.Context, id string) (*base.UserProfile, error)
}
