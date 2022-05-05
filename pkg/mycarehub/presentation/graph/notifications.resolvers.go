package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/firebasetools"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) SendFCMNotification(ctx context.Context, registrationTokens []string, data map[string]interface{}, notification firebasetools.FirebaseSimpleNotificationInput) (bool, error) {
	return r.mycarehub.Notification.SendNotification(ctx, registrationTokens, data, &notification)
}

func (r *queryResolver) FetchNotifications(ctx context.Context, userID string, flavour feedlib.Flavour, paginationInput dto.PaginationsInput) (*domain.NotificationsPage, error) {
	return r.mycarehub.Notification.FetchNotifications(ctx, userID, flavour, paginationInput)
}
