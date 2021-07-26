package fb_test

import (
	"context"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/savannahghi/interserviceclient"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure/database/fb"
	"github.com/stretchr/testify/assert"
)

func TestNewFirebaseClientExtensionImpl(t *testing.T) {
	fbdb := fb.NewFirebaseClientExtensionImpl()
	assert.NotNil(t, fbdb)

	// GetUserByPhoneNumber should fail
	assert.Panics(t, func() {
		_, _ = fbdb.GetUserByPhoneNumber(context.Background(), interserviceclient.TestUserPhoneNumber)
	})

	// CreateUser should fail
	assert.Panics(t, func() {
		_, _ = fbdb.CreateUser(context.Background(), &auth.UserToCreate{})
	})

	// DeleteUser should fail
	assert.Panics(t, func() {
		_ = fbdb.DeleteUser(context.Background(), uuid.New().String())
	})
}
