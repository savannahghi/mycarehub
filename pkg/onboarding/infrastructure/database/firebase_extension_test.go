package database_test

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/database"
)

func TestNewFirebaseClientExtensionImpl(t *testing.T) {
	fb := database.NewFirebaseClientExtensionImpl()
	assert.NotNil(t, fb)

	// GetUserByPhoneNumber should fail
	assert.Panics(t, func() {
		_, _ = fb.GetUserByPhoneNumber(context.Background(), base.TestUserPhoneNumber)
	})

	// CreateUser should fail
	assert.Panics(t, func() {
		_, _ = fb.CreateUser(context.Background(), &auth.UserToCreate{})
	})

	// DeleteUser should fail
	assert.Panics(t, func() {
		_ = fb.DeleteUser(context.Background(), uuid.New().String())
	})
}
