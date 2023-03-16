package utils

import (
	"context"

	"github.com/savannahghi/firebasetools"
)

func GetLoggedInUserID(ctx context.Context) *string {
	userID, err := firebasetools.GetLoggedInUserUID(ctx)
	if err != nil {
		return nil
	}
	return &userID
}
